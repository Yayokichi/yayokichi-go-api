// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand" // ★ 追加: 乱数を生成するために必要
	"net/http"
	"os"
	"sync"
	"time"
)

// --- 構造体定義 (変更なし) ---
type StatusResponseData struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	Language  string `json:"language"`
	Timestamp string `json:"timestamp"`
}
type SpeedTestResponseData struct {
	SerialExecutionTime   string `json:"serialExecutionTime"`
	ParallelExecutionTime string `json:"parallelExecutionTime"`
	TaskCount             int    `json:"taskCount"`
	TaskDuration          string `json:"taskDuration"`
	Message               string `json:"message"`
}

// ★ 修正: 処理時間にランダムな揺らぎを持たせる
func dummyTask() {
	// 950msをベースに、0〜100msのランダムな時間を加える (合計 950ms 〜 1049ms)
	baseDuration := 950 * time.Millisecond
	randomOffset := time.Duration(rand.Intn(100)) * time.Millisecond // 0 to 99
	totalDuration := baseDuration + randomOffset
	time.Sleep(totalDuration)
}

// --- statusHandler (変更なし) ---
func statusHandler(w http.ResponseWriter, r *http.Request) {
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

// ★ 修正: speedTestHandler
func speedTestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	const taskCount = 3

	// --- 1. 直列実行 ---
	startSerial := time.Now()
	for i := 0; i < taskCount; i++ {
		dummyTask() // 引数をなくす
	}
	serialDuration := time.Since(startSerial)

	// --- 2. Goroutineを使った並行実行 ---
	var wg sync.WaitGroup
	wg.Add(taskCount)
	startParallel := time.Now()
	for i := 0; i < taskCount; i++ {
		go func() {
			defer wg.Done()
			dummyTask() // 引数をなくす
		}()
	}
	wg.Wait()
	parallelDuration := time.Since(startParallel)

	// フロントエンドに返すレスポンスを作成
	data := SpeedTestResponseData{
		SerialExecutionTime:   fmt.Sprintf("%.2f秒", serialDuration.Seconds()),
		ParallelExecutionTime: fmt.Sprintf("%.2f秒", parallelDuration.Seconds()),
		TaskCount:             taskCount,
		TaskDuration:          "約1秒",                                            // 表示を固定の文字列に変更
		Message:               fmt.Sprintf("%d個の「約1秒」かかる処理を実行しました。", taskCount), // メッセージも変更
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// ★ 追加: 乱数のシードを初期化。これにより毎回異なる乱数が生成される
	rand.New(rand.NewSource(time.Now().UnixNano()))

	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/speed-test", speedTestHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Go API Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
