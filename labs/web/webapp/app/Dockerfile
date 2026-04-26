FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod .
COPY *.go .
RUN go mod tidy && CGO_ENABLED=0 go build -o /app .

FROM debian:bookworm-slim@sha256:734a72b3d6142a71993b4bc580ca01d8ca687f77a2c422d96a961d3376b179bc
COPY --from=build /app /app
# パストラバーサル演習用の公開文書ディレクトリ
RUN mkdir -p /srv/docs && echo 'これは公開文書です。' > /srv/docs/readme.txt
EXPOSE 8080
CMD ["/app"]
