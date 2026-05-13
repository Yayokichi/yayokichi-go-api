// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// 返却するJSONデータの構造体
type ResponseData struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	Language  string `json:"language"`
	Timestamp string `json:"timestamp"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// CORS対応: どこからでも(Next.jsから)呼び出せるようにする
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// 構造体にデータをセット
	data := ResponseData{
		Message:   "Hello from Go Microservice!",
		Status:    "Active",
		Language:  "Go (Golang)",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// JSONにエンコードしてレスポンスを返す
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// ルーティングの設定
	http.HandleFunc("/api/status", apiHandler)

	// RenderなどのPaaSは環境変数PORTでポートを指定してくるため、それを取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル開発用のデフォルトポート
	}

	fmt.Printf("Go API Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}