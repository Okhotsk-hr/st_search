package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
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

func timetableHandler(w http.ResponseWriter, r *http.Request) {
	targetBin := r.URL.Query().Get("bin")
	if targetBin == "" {
		http.Error(w, "bin parameter is required", http.StatusBadRequest)
		return
	}

	stopTimesPath := "data/stop_times.json"
	stopsPath := "data/stops.json"

	// stop_times.json 読み込み
	stopTimesData, err := os.ReadFile(stopTimesPath)
	if err != nil {
		http.Error(w, "stop_times.json 読み込みエラー", 500)
		return
	}

	var stopTimes []StopTime
	if err := json.Unmarshal(stopTimesData, &stopTimes); err != nil {
		http.Error(w, "stop_times.json パースエラー", 500)
		return
	}

	// stops.json 読み込み
	stopsData, err := os.ReadFile(stopsPath)
	if err != nil {
		http.Error(w, "stops.json 読み込みエラー", 500)
		return
	}

	var stops []Stop
	if err := json.Unmarshal(stopsData, &stops); err != nil {
		http.Error(w, "stops.json パースエラー", 500)
		return
	}

	// stop_id → stop_name
	stopNameMap := make(map[string]string)
	for _, s := range stops {
		stopNameMap[s.StopID] = s.StopName
	}

	// 対象便のみ抽出
	var result []StopTime
	for _, st := range stopTimes {
		parts := strings.Split(st.TripID, "+")
		if len(parts) >= 3 && parts[len(parts)-1] == targetBin {
			result = append(result, st)
		}
	}

	if len(result) == 0 {
		http.Error(w, "指定された便番号のデータがありません", http.StatusNotFound)
		return
	}

	// stop_sequence で並び替え
	sort.Slice(result, func(i, j int) bool {
		a, _ := strconv.Atoi(result[i].StopSequence)
		b, _ := strconv.Atoi(result[j].StopSequence)
		return a < b
	})

	// テキストとして返す
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintf(w, "便番号 %s の時刻一覧\n", targetBin)
	fmt.Fprintln(w, "------------------------------------------")

	for _, st := range result {
		stopName := stopNameMap[st.StopID]
		if stopName == "" {
			stopName = "(不明)"
		}

		fmt.Fprintf(
			w,
			"順序:%s 停留所:%s 到着:%s 出発:%s\n",
			st.StopSequence,
			stopName,
			st.ArrivalTime,
			st.DepartureTime,
		)
	}
}

func main() {
	http.HandleFunc("/timetable", timetableHandler)
	http.ListenAndServe(":8080", nil)
}
