package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type StopTime struct {
	ArrivalTime   string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
	StopID        string `json:"stop_id"`
	StopSequence  string `json:"stop_sequence"`
	TripID        string `json:"trip_id"`
}

func main() {
	filePath := "data/stop_times.json"
	targetBin := "9300656" // 便番号

	// ファイル読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("ファイル読み込みエラー:", err)
		return
	}

	// JSONパース
	var stops []StopTime
	if err := json.Unmarshal(data, &stops); err != nil {
		fmt.Println("JSONパースエラー:", err)
		return
	}

	fmt.Println("便番号", targetBin, "の時刻一覧")
	fmt.Println("---------------------------")

	for _, s := range stops {
		// trip_id を + で分割
		parts := strings.Split(s.TripID, "+")
		if len(parts) < 3 {
			continue
		}

		binNo := parts[len(parts)-1]
		if binNo == targetBin {
			fmt.Printf(
				"順序:%s 停留所:%s 到着:%s 出発:%s\n",
				s.StopSequence,
				s.StopID,
				s.ArrivalTime,
				s.DepartureTime,
			)
		}
	}
}
