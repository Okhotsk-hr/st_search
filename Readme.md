# 札幌市交通局のGTFSデータで遊んでみる
### 各データ説明
#### agency
会社情報

#### calendar_dates

#### calendar

#### fare_attributes

#### fare_rules

#### feed_info

#### route_jp

#### routes

#### stop_times
**"trip_id"** <br>
初め6ケタの数字が"route_id"<br>
後ろ7ケタの数字が便番号？<br>

#### stops

#### translations

#### trips

### API説明
#### times.go
便の番号を入れたら各停留所と着発時刻を返す
"http://localhost:8080/timetable?bin=9300656"
