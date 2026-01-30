[![Go test](https://github.com/kaodekaoru2/video-converter-go/actions/workflows/test.yml/badge.svg)](https://github.com/kaodekaoru2/video-converter-go/actions/workflows/test.yml)

# 🎥 Go Video Converter

Go言語とFFmpegを使用した、Webベースの動画変換ツールです。
Goの開発練習として作成しました。

## 機能
- **MP4動画のアップロード**
- **多彩な変換モード**
  - GIFアニメーション生成
  - MP3音声抽出
  - JPEGサムネイル生成
- **高機能バックエンド**
  - 処理時間の計測ログ
  - 自動ファイル削除（お掃除機能）
  - Goroutineによる非同期処理

## 技術スタック
- Go (net/http, os/exec, goroutine)
- FFmpeg
- HTML/CSS/JavaScript

## 使い方
1. GoとFFmpegをインストール
2. `go run main.go`
3. `http://localhost:8080` にアクセス