package main

import (
	"encoding/json"
	"net/http"
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

type Route struct {
	RouteID       string `json:"route_id"`
	RouteLongName string `json:"route_long_name"`
}

func timetableHandler(w http.ResponseWriter, r *http.Request) {
	stopName := r.URL.Query().Get("stop_name")
	if stopName == "" {
		http.Error(w, "stop_name is required", 400)
		return
	}

	// ===== stops.json =====
	var stops []Stop
	json.NewDecoder(mustOpen("data/stops.json")).Decode(&stops)

	stopID := ""
	for _, s := range stops {
		if s.StopName == stopName {
			stopID = strings.TrimSpace(s.StopID)
			break
		}
	}
	if stopID == "" {
		http.Error(w, "stop not found", 404)
		return
	}

	// ===== routes.json =====
	var routes []Route
	json.NewDecoder(mustOpen("data/routes.json")).Decode(&routes)

	routeMap := make(map[string]string)
	for _, r := range routes {
		routeMap[r.RouteID] = r.RouteLongName
	}

	// ===== stop_times.json =====
	var times []StopTime
	json.NewDecoder(mustOpen("data/stop_times.json")).Decode(&times)

	// ===== grouping =====
	type key struct {
		Route string
		Day   string
	}

	group := make(map[key][]string)
	order := []key{}

	for _, t := range times {
		if strings.TrimSpace(t.StopID) != stopID {
			continue
		}

		parts := strings.Split(t.TripID, "+")
		if len(parts) < 2 {
			continue
		}

		routeID := parts[0]
		day := parts[1]

		routeName := routeMap[routeID]
		k := key{routeName, day}

		if _, ok := group[k]; !ok {
			order = append(order, k)
		}
		group[k] = append(group[k], t.DepartureTime)
	}

	// ===== 3次元配列生成 =====
	result := [][][]string{}

	for _, k := range order {
		block := [][]string{
			{k.Route, k.Day},
			group[k],
		}
		result = append(result, block)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func mustOpen(path string) *os.File {
	f, _ := os.Open(path)
	return f
}

func main() {
	http.HandleFunc("/timetable", timetableHandler)
	http.ListenAndServe(":8080", nil)
}
