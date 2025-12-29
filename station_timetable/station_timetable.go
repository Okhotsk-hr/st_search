package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Stop struct {
	StopID   string `json:"stop_id"`
	StopName string `json:"stop_name"`
}

type StopTime struct {
	StopID        string `json:"stop_id"`
	DepartureTime string `json:"departure_time"`
}

func main() {
	// ========= 駅名入力 =========
	fmt.Print("駅名を入力してください: ")
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	stationName := strings.TrimSpace(in.Text())

	// ========= stops.json 読み込み =========
	stopsFile, err := os.Open("data/stops.json")
	if err != nil {
		fmt.Println("stops.json を開けません:", err)
		return
	}
	defer stopsFile.Close()

	var stops []Stop
	if err := json.NewDecoder(stopsFile).Decode(&stops); err != nil {
		fmt.Println("stops.json の読み込み失敗:", err)
		return
	}

	// 駅名 → stop_id
	stopID := ""
	for _, s := range stops {
		if s.StopName == stationName {
			stopID = strings.TrimSpace(s.StopID)
			break
		}
	}

	if stopID == "" {
		fmt.Println("駅名が見つかりません")
		return
	}

	fmt.Println("取得した stop_id:", stopID)

	// ========= stop_times.json 読み込み =========
	timesFile, err := os.Open("data/stop_times.json")
	if err != nil {
		fmt.Println("stop_times.json を開けません:", err)
		return
	}
	defer timesFile.Close()

	var stopTimes []StopTime
	if err := json.NewDecoder(timesFile).Decode(&stopTimes); err != nil {
		fmt.Println("stop_times.json の読み込み失敗:", err)
		return
	}

	// stop_id → departure_time 一覧
	found := false
	fmt.Println("departure_time 一覧:")

	for _, st := range stopTimes {
		if strings.TrimSpace(st.StopID) == stopID {
			fmt.Println(st.DepartureTime)
			found = true
		}
	}

	if !found {
		fmt.Println("⚠ 該当する出発時刻がありません")
	}
}
