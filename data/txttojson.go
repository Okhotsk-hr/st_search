package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	// 引数チェック
	if len(os.Args) != 3 {
		fmt.Println("使い方:")
		fmt.Println("  go run gtfs_to_json.go 入力.txt 出力.json")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// ファイル読み込み
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("入力ファイル読み込みエラー:", err)
		return
	}

	// ★ ファイル先頭の BOM 除去
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})

	// CSV として読み込み
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("CSV解析エラー:", err)
		return
	}

	if len(records) < 2 {
		fmt.Println("データがありません")
		return
	}

	// ヘッダ行
	headers := records[0]
	headers[0] = strings.TrimPrefix(headers[0], "\uFEFF") // 念のため

	// データ変換
	var result []map[string]string

	for _, row := range records[1:] {
		m := map[string]string{}
		for i, h := range headers {
			if i < len(row) {
				m[h] = row[i]
			}
		}
		result = append(result, m)
	}

	// JSON出力（BOMなし）
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("JSON生成エラー:", err)
		return
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("出力ファイル書き込みエラー:", err)
		return
	}

	fmt.Println("変換完了:", outputFile)
}
