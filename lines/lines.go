package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// routes.json の1件分を map で受ける（BOM対策）
type Route map[string]string

func main() {
	// ファイル読み込み
	data, err := os.ReadFile("data/routes.json")
	if err != nil {
		fmt.Println("ファイル読み込みエラー:", err)
		return
	}

	// ★ JSON中に混入したBOMをすべて除去
	bom := []byte{0xEF, 0xBB, 0xBF}
	data = bytes.ReplaceAll(data, bom, []byte{})

	// JSON解析
	var routes []Route
	if err := json.Unmarshal(data, &routes); err != nil {
		fmt.Println("JSON解析エラー:", err)
		return
	}

	// 全データをリストアップ
	for i, r := range routes {
		fmt.Printf("===== Route %d =====\n", i+1)

		for key, value := range r {
			fmt.Printf("%s: %s\n", key, value)
		}
		fmt.Println()
	}
}
