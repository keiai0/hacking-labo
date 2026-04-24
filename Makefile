# 使い方: make <target> LAB=network/01-recon
LAB ?=
PROJECT = hl-$(subst /,-,$(LAB))
COMPOSE = docker compose -f labs/$(LAB)/compose.yml -p $(PROJECT)

.PHONY: help attacker list up down shell ps logs capture capture-stop clean _require-lab

help: ## コマンド一覧
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-14s %s\n", $$1, $$2}'

attacker: ## 攻撃用コンテナをビルド
	docker build -t hacklabo/attacker:latest common/attacker

list: ## ラボ一覧
	@find labs -name compose.yml | sed 's|labs/||; s|/compose.yml||' | sort

up: _require-lab ## ラボを起動
	$(COMPOSE) up -d --build

down: _require-lab ## ラボを停止・削除
	$(COMPOSE) --profile capture down -v

shell: _require-lab ## 攻撃用コンテナに入る
	$(COMPOSE) exec attacker bash

ps: _require-lab ## ラボのコンテナ一覧
	$(COMPOSE) ps

logs: _require-lab ## ラボのログ
	$(COMPOSE) logs -f

capture: _require-lab ## ネットワーク全体のキャプチャ開始（captures/ に保存）
	$(COMPOSE) --profile capture up -d capture

capture-stop: _require-lab ## キャプチャ停止
	$(COMPOSE) --profile capture stop capture

clean: ## 未使用のイメージ等を削除
	docker system prune -f

_require-lab:
	@test -n "$(LAB)" || (echo "LAB を指定してください（例: make up LAB=network/01-recon）"; make -s list; exit 1)
