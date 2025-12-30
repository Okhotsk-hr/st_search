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
	TripID        string `json:"trip_id"`
}

func main() {
	// ========= 駅名入力 =========
	fmt.Print("駅名を入力してください: ")
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	stationName := strings.TrimSpace(in.Text())

	// ========= stops.json =========
	stopsFile, err := os.Open("data/stops.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stopsFile.Close()

	var stops []Stop
	json.NewDecoder(stopsFile).Decode(&stops)

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

	// ========= stop_times.json =========
	timesFile, err := os.Open("data/stop_times.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer timesFile.Close()

	var stopTimes []StopTime
	json.NewDecoder(timesFile).Decode(&stopTimes)

	fmt.Println("departure_time 一覧:")

	var prevHead string
	var prevMid string
	first := true

	for _, st := range stopTimes {
		if st.StopID != stopID {
			continue
		}

		parts := strings.Split(st.TripID, "+")
		if len(parts) < 2 {
			continue
		}

		head := parts[0] // 先頭6桁
		mid := parts[1]  // 真ん中文字列

		// 区切り判定
		if first || head != prevHead || mid != prevMid {
			fmt.Println("----", head, mid, "----")
			prevHead = head
			prevMid = mid
			first = false
		}

		fmt.Println(st.DepartureTime)
	}
}
