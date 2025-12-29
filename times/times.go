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

type Stop struct {
	StopID   string `json:"stop_id"`
	StopName string `json:"stop_name"`
}

func main() {
	stopTimesPath := "data/stop_times.json"
	stopsPath := "data/stops.json"
	targetBin := "9300656"

	// stop_times.json 読み込み
	stopTimesData, err := os.ReadFile(stopTimesPath)
	if err != nil {
		fmt.Println("stop_times.json 読み込みエラー:", err)
		return
	}

	var stopTimes []StopTime
	if err := json.Unmarshal(stopTimesData, &stopTimes); err != nil {
		fmt.Println("stop_times.json パースエラー:", err)
		return
	}

	// stops.json 読み込み
	stopsData, err := os.ReadFile(stopsPath)
	if err != nil {
		fmt.Println("stops.json 読み込みエラー:", err)
		return
	}

	var stops []Stop
	if err := json.Unmarshal(stopsData, &stops); err != nil {
		fmt.Println("stops.json パースエラー:", err)
		return
	}

	// stop_id → stop_name のマップを作成
	stopNameMap := make(map[string]string)
	for _, s := range stops {
		stopNameMap[s.StopID] = s.StopName
	}

	fmt.Println("便番号", targetBin, "の時刻一覧")
	fmt.Println("------------------------------------------")

	for _, st := range stopTimes {
		parts := strings.Split(st.TripID, "+")
		if len(parts) < 3 {
			continue
		}

		binNo := parts[len(parts)-1]
		if binNo != targetBin {
			continue
		}

		stopName := stopNameMap[st.StopID]
		if stopName == "" {
			stopName = "(不明)"
		}

		// ★ 停留所番号を出さず、名前のみ表示
		fmt.Printf(
			"順序:%s 停留所:%s 到着:%s 出発:%s\n",
			st.StopSequence,
			stopName,
			st.ArrivalTime,
			st.DepartureTime,
		)
	}
}
