package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Stop struct {
	LocationType string `json:"location_type"`
	StopID       string `json:"stop_id"`
	StopLat      string `json:"stop_lat"`
	StopLon      string `json:"stop_lon"`
	StopName     string `json:"stop_name"`
	ZoneID       string `json:"zone_id"`
}

func main() {
	// JSONファイルを開く
	file, err := os.Open("data/stops.json")
	if err != nil {
		fmt.Println("ファイルを開けません:", err)
		return
	}
	defer file.Close()

	// JSONを読み込む
	var stops []Stop
	if err := json.NewDecoder(file).Decode(&stops); err != nil {
		fmt.Println("JSONの読み込みに失敗:", err)
		return
	}

	// 検索する停留所名を入力
	fmt.Print("停留所名を入力してください: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputName := scanner.Text()

	// 検索
	found := false
	for _, stop := range stops {
		if stop.StopName == inputName {
			fmt.Println("stop_id:", stop.StopID)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("該当する停留所が見つかりません")
	}
}
