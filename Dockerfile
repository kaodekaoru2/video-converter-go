# 1. ビルド用ステージ（Go言語の環境）
FROM golang:1.23-alpine AS builder

WORKDIR /app

# すべてのファイルをコピー
COPY . .

# Goの準備（go.modがなくてもエラーにならないようにする魔法）
RUN go mod init video-converter || true
RUN go mod tidy
# アプリをビルドして「main」という実行ファイルを作る
RUN go build -o main main.go

# 2. 実行用ステージ（本番で動く軽量Linux）
FROM alpine:latest

# ここでFFmpegをインストール！（これがDockerの強み）
RUN apk add --no-cache ffmpeg

WORKDIR /root/

# さっき作った実行ファイルとHTMLをコピー
COPY --from=builder /app/main .
COPY --from=builder /app/index.html .

# 保存用フォルダを作成
RUN mkdir uploads

# アプリ起動！
CMD ["./main"]