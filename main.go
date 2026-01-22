package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// 1. トップページ (index.htmlを表示するだけ)
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// index.html が同じフォルダにある前提です
	http.ServeFile(w, r, "index.html")
}

// 2. 変換処理のメインロジック
func convertHandler(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外は拒否
	if r.Method != "POST" {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	// --- Step A: ユーザーの入力を受け取る ---

	// 変換モードを取得 (gif, mp3, jpg)
	format := r.FormValue("format")

	// 動画ファイルを受け取る
	file, header, err := r.FormFile("videoFile")
	if err != nil {
		http.Error(w, "ファイル読み込みエラー", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// --- Step B: ファイルを保存する ---

	// uploadsフォルダのパス (なければエラーになるので注意)
	inputPath := filepath.Join("uploads", header.Filename)

	// 出力ファイル名を作成 (例: video.mp4.mp3)
	outputFilename := header.Filename + "." + format
	outputPath := filepath.Join("uploads", outputFilename)

	// 保存用の空ファイルを作る
	dst, err := os.Create(inputPath)
	if err != nil {
		http.Error(w, "保存エラー (uploadsフォルダはありますか？)", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 中身をコピー
	io.Copy(dst, file)

	// --- Step C: FFmpegコマンドを実行する (分岐) ---

	var cmd *exec.Cmd

	switch format {
	case "mp3":
		// 音声抽出: -vn (Video None)
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-vn", "-y", outputPath)
	case "jpg":
		// サムネイル: 開始1秒(-ss) の1フレーム(-vframes) を切り出す
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-ss", "00:00:01", "-vframes", "1", "-y", outputPath)
	default:
		// GIF変換 (デフォルト)
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-y", outputPath)
	}

	// 実行ログをターミナルに出す
	fmt.Printf("コマンド実行中: %v\n", cmd.Args)

	err = cmd.Run()
	if err != nil {
		http.Error(w, "変換失敗... FFmpegのエラーです: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// --- Step D: 結果画面を表示する ---

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, "<h1>変換成功！ (%s)</h1>", format)

	// ダウンロード用のリンク
	downloadLink := "/uploads/" + outputFilename

	// 画像(jpg, gif)なら、その場でプレビュー表示してみる
	if format == "jpg" || format == "gif" {
		fmt.Fprintf(w, `<img src="%s" width="400" style="border: 2px solid #333;"><br><br>`, downloadLink)
	}

	fmt.Fprintf(w, `<a href="%s" download style="font-size:20px;">[ ダウンロードはこちら ]</a><br><br>`, downloadLink)
	fmt.Fprintf(w, `<a href="/">← もう一度やる</a>`)
}

func main() {
	// uploadsフォルダをブラウザから見れるようにする設定
	fs := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// ルーティング
	http.HandleFunc("/", uploadHandler)
	http.HandleFunc("/convert", convertHandler)

	fmt.Println("サーバー起動中... http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
