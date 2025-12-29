package main

import (
	"encoding/json"
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

// 返却する1行分
type TimetableItem struct {
	Sequence      int    `json:"sequence"`
	StopName      string `json:"stop_name"`
	ArrivalTime   string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
}

func timetableHandler(w http.ResponseWriter, r *http.Request) {
	targetBin := r.URL.Query().Get("bin")
	if targetBin == "" {
		http.Error(w, "bin parameter is required", http.StatusBadRequest)
		return
	}

	// stop_times.json 読み込み
	stopTimesData, err := os.ReadFile("data/stop_times.json")
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
	stopsData, err := os.ReadFile("data/stops.json")
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

	// 対象便を抽出
	var items []TimetableItem
	for _, st := range stopTimes {
		parts := strings.Split(st.TripID, "+")
		if len(parts) >= 3 && parts[len(parts)-1] == targetBin {

			seq, _ := strconv.Atoi(st.StopSequence)
			name := stopNameMap[st.StopID]
			if name == "" {
				name = "(不明)"
			}

			items = append(items, TimetableItem{
				Sequence:      seq,
				StopName:      name,
				ArrivalTime:   st.ArrivalTime,
				DepartureTime: st.DepartureTime,
			})
		}
	}

	if len(items) == 0 {
		http.Error(w, "指定された便番号のデータがありません", http.StatusNotFound)
		return
	}

	// 順序でソート
	sort.Slice(items, func(i, j int) bool {
		return items[i].Sequence < items[j].Sequence
	})

	// JSON配列として返す
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(items)
}

func main() {
	http.HandleFunc("/timetable", timetableHandler)
	http.ListenAndServe(":8080", nil)
}
