<?php
// XXE演習（PHP/libxml）。外部実体を解決するかどうか。
$mode = ($_GET['mode'] ?? 'vuln') === 'fixed' ? 'fixed' : 'vuln';
$other = $mode === 'vuln' ? 'fixed' : 'vuln';
$sample = "<?xml version=\"1.0\"?>\n<!DOCTYPE r [<!ENTITY x SYSTEM \"file:///etc/passwd\">]>\n<order><item>&x;</item></order>";
$result = null;

if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $xml = file_get_contents('php://input');
    $doc = new DOMDocument();
    if ($mode === 'vuln') {
        // 脆弱: DTD読込＋実体展開を有効化 → 外部ファイルを読み込む
        $doc->loadXML($xml, LIBXML_DTDLOAD | LIBXML_NOENT);
    } else {
        // 対策: 既定のまま（外部実体を展開しない）
        $doc->loadXML($xml);
    }
    $node = $doc->getElementsByTagName('item')->item(0);
    $result = $node ? $node->textContent : '(解析不可)';
}

header('Content-Type: text/html; charset=utf-8');
echo '<!doctype html><meta charset="utf-8"><title>XXE</title>';
echo '<body style="font-family:sans-serif;max-width:760px;margin:2rem auto">';
echo '<nav style="font-size:.9em"><a href="/">&laquo; メニュー</a></nav><hr><h1>XXE（XML外部実体）</h1>';
echo "<p>現在: <b>" . ($mode === 'vuln' ? '脆弱版' : '対策版') . "</b>（<a href=\"/xxe?mode=$other\">切替</a>）</p>";
echo '<form method="post" action="/xxe?mode=' . $mode . '">';
echo '<p>注文XML:</p><textarea name="_ignored" rows=5 cols=70>' . htmlspecialchars($sample) . '</textarea><br>';
echo '<button formmethod="post">送信（上記サンプルXMLを本文として送る）</button></form>';
echo '<p style="color:#666">※ curl例: <code>curl --data-binary @- .../xxe?mode=vuln</code> にサンプルXMLを渡す</p>';
if ($result !== null) {
    echo '<h2>解析結果（item要素）</h2><pre>' . htmlspecialchars($result) . '</pre>';
}
echo '<hr><p style="color:#666">脆弱版は &amp;x; が /etc/passwd に展開される／対策版は展開されない。</p></body>';
