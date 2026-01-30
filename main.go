package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time" // ★追加: 時間を測るためのライブラリ
)

// 1. トップページ (index.htmlを表示するだけ)
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// 2. 変換処理のメインロジック
func convertHandler(w http.ResponseWriter, r *http.Request) {
	// ★計測スタート！
	start := time.Now()

	if r.Method != "POST" {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	// --- Step A: 入力受け取り ---
	format := r.FormValue("format")
	file, header, err := r.FormFile("videoFile")
	if err != nil {
		http.Error(w, "ファイル読み込みエラー", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// --- Step B: 保存 ---
	inputPath := filepath.Join("uploads", header.Filename)
	outputFilename := header.Filename + "." + format
	outputPath := filepath.Join("uploads", outputFilename)

	dst, err := os.Create(inputPath)
	if err != nil {
		http.Error(w, "保存エラー", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// --- Step C: FFmpeg実行 ---
	var cmd *exec.Cmd

	switch format {
	case "mp3":
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-vn", "-y", outputPath)
	case "jpg":
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-ss", "00:00:01", "-vframes", "1", "-y", outputPath)
	default: // gif
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-y", outputPath)
	}

	fmt.Printf("コマンド実行中: %v\n", cmd.Args)

	err = cmd.Run()
	if err != nil {
		http.Error(w, "変換失敗: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ★計測ストップ＆表示！
	duration := time.Since(start)
	fmt.Printf("✅ 完了！かかった時間: %v\n", duration)

	// --- Step D: 結果表示 ---
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>変換成功！ (%s)</h1>", format)

	// 処理時間を画面にも出してみる（オプション）
	fmt.Fprintf(w, "<p>処理時間: %v</p>", duration)

	downloadLink := "/uploads/" + outputFilename

	if format == "jpg" || format == "gif" {
		fmt.Fprintf(w, `<img src="%s" width="400" style="border: 2px solid #333;"><br><br>`, downloadLink)
	}

	fmt.Fprintf(w, `<a href="%s" download style="font-size:20px;">[ ダウンロードはこちら ]</a><br><br>`, downloadLink)
	fmt.Fprintf(w, `<a href="/">← もう一度やる</a>`)

	go func() {
		// 1. 5分待つ（テスト用に短くしたいなら time.Second * 10 とかにしてもOK）
		time.Sleep(5 * time.Minute)

		// 2. ファイルを削除する
		fmt.Println("🧹 お掃除開始: 古いファイルを削除します...")
		os.Remove(inputPath)  // 元動画を消す
		os.Remove(outputPath) // 変換後ファイルを消す
		fmt.Println("✨ お掃除完了: ", inputPath, outputPath)
	}()
}

func main() {
	// uploadsフォルダの公開設定
	fs := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", fs))

	// ルーティング
	http.HandleFunc("/", uploadHandler)
	http.HandleFunc("/convert", convertHandler)

	fmt.Println("🚀 サーバー起動中... http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Addは2つの整数を足す関数（テスト練習用）
func Add(a int, b int) int {
	return a + b
}
