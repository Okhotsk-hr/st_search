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

type Route struct {
	RouteID       string `json:"route_id"`
	RouteLongName string `json:"route_long_name"`
}

func main() {
	// ========= 駅名入力 =========
	fmt.Print("駅名を入力してください: ")
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	stationName := strings.TrimSpace(in.Text())

	// ========= stops.json =========
	stopsFile, _ := os.Open("data/stops.json")
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

	// ========= routes.json =========
	routesFile, _ := os.Open("data/routes.json")
	defer routesFile.Close()

	var routes []Route
	json.NewDecoder(routesFile).Decode(&routes)

	routeMap := make(map[string]string)
	for _, r := range routes {
		routeMap[r.RouteID] = r.RouteLongName
	}

	// ========= stop_times.json =========
	timesFile, _ := os.Open("data/stop_times.json")
	defer timesFile.Close()

	var stopTimes []StopTime
	json.NewDecoder(timesFile).Decode(&stopTimes)

	fmt.Println("departure_time 一覧:")

	var prevRoute string
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

		routeID := parts[0]
		mid := parts[1]

		routeName, ok := routeMap[routeID]
		if !ok {
			routeName = routeID // 見つからなければ番号表示
		}

		if first || routeID != prevRoute || mid != prevMid {
			fmt.Println("----", routeName, mid, "----")
			prevRoute = routeID
			prevMid = mid
			first = false
		}

		fmt.Println(st.DepartureTime)
	}
}
