// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync" // ★ 追加: Goroutineの完了を待つために必要
	"time"
)

// --- 既存の /api/status 用の構造体 ---
type StatusResponseData struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	Language  string `json:"language"`
	Timestamp string `json:"timestamp"`
}

// ★ 追加: /api/speed-test 用の新しい構造体 ---
type SpeedTestResponseData struct {
	SerialExecutionTime   string `json:"serialExecutionTime"`
	ParallelExecutionTime string `json:"parallelExecutionTime"`
	TaskCount             int    `json:"taskCount"`
	TaskDuration          string `json:"taskDuration"`
	Message               string `json:"message"`
}

// 1秒かかるダミーの重い処理
func dummyTask(duration time.Duration) {
	time.Sleep(duration)
}

// --- 既存の /api/status 用のハンドラ ---
func statusHandler(w http.ResponseWriter, r *http.Request) {
	// CORS対応
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	data := StatusResponseData{
		Message:   "Hello from Go Microservice!",
		Status:    "Active",
		Language:  "Go (Golang)",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ★ 追加: /api/speed-test 用の新しいハンドラ ---
func speedTestHandler(w http.ResponseWriter, r *http.Request) {
	// CORS対応
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	const taskCount = 3                  // 実行するタスクの数
	const taskDuration = 1 * time.Second // 1つのタスクにかかる時間

	// --- 1. 直列実行 ---
	startSerial := time.Now()
	for i := 0; i < taskCount; i++ {
		dummyTask(taskDuration)
	}
	serialDuration := time.Since(startSerial)

	// --- 2. Goroutineを使った並行実行 ---
	var wg sync.WaitGroup
	wg.Add(taskCount) // WaitGroupにタスクの数をセット

	startParallel := time.Now()
	for i := 0; i < taskCount; i++ {
		go func() {
			defer wg.Done() // Goroutineが完了したらカウンタをデクリメント
			dummyTask(taskDuration)
		}()
	}
	wg.Wait() // すべてのGoroutineが完了するまで待機
	parallelDuration := time.Since(startParallel)

	// フロントエンドに返すレスポンスを作成
	data := SpeedTestResponseData{
		SerialExecutionTime:   fmt.Sprintf("%.2f秒", serialDuration.Seconds()),
		ParallelExecutionTime: fmt.Sprintf("%.2f秒", parallelDuration.Seconds()),
		TaskCount:             taskCount,
		TaskDuration:          taskDuration.String(),
		Message:               fmt.Sprintf("%d個の「%s」かかる処理を実行しました。", taskCount, taskDuration.String()),
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// ルーティングの設定
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/speed-test", speedTestHandler) // ★ 追加: 新しいルート

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Go API Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
