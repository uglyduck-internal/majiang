package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"io/ioutil"
	"net/http"
	"time"
)

var TOKEN string

func init() {
	TOKEN = "Bearer eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ3eGFwcGxldHRqb285MHI4dHJlYSIsImNyZWF0ZWQiOjE3MTgxNjk1MDc0ODAsImV4cCI6MTgzODE2OTUwN30.pA7sIS8fCAr4UXaCaSQ9KkxeXPFZY848MvAGnyfDpFKLi4I6PGWqWg2CZUzOJay2TKS7Kd2lV2R1sBOSzo0BFw"
}

func getStores() interface{} {
	req, err := http.NewRequest("GET", "https://portal.quedou.cn/api/qsg/store/list", nil)
	if err != nil {
		fmt.Println("ERROR: Failed to create the HTTP request")
		return nil
	}

	req.Header.Set("Authorization", TOKEN)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 MicroMessenger/6.8.0(0x16080000) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080710) XWEB/1191")
	req.Header.Set("Referer", "https://servicewechat.com/wx697f0b89354ff12e/26/page-frame.html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Failed to send the HTTP request")
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Failed to read the response body")
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("ERROR: [getStores] Failed to parse the response body as JSON")
		return nil
	}

	data, ok := result["data"]
	if !ok {
		fmt.Println("ERROR: [getStores] No data field in the response body")
		return nil
	}

	return data
}

func getRooms(store map[string]interface{}) interface{} {
	storeID := int(store["id"].(float64))

	url := fmt.Sprintf("https://portal.quedou.cn/api/dms/device/list?storeId=%d&memberId=76965", storeID)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", TOKEN)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 MicroMessenger/6.8.0(0x16080000) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080710) XWEB/1191")
	req.Header.Set("Referer", "https://servicewechat.com/wx697f0b89354ff12e/26/page-frame.html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("ERROR: Failed to send the HTTP request, url: ", url)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: Failed to read the response body")
		return nil
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("ERROR: Failed to parse the response body as JSON")
		return nil
	}

	data, ok := result["data"]
	if !ok {
		fmt.Println("ERROR: No data field in the response body")
		return nil
	}

	return data
}
func initDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		"192.168.1.144", "postgres", "pgpassword", "majiang")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// 创建表
	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS rooms ("time" TEXT, 
"room_id" TEXT,
"room_name" TEXT,
"store_name" TEXT, 
"store_id" TEXT,
"store_address" TEXT,
"price" TEXT,
"status" TEXT,
UNIQUE(time, room_id)
                                 );`)
	if err != nil {
		fmt.Println(err)
	}
	var result sql.Result
	result, err = statement.Exec()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	return db
}

func insertRoom(db *sql.DB, timeh string, store interface{}, room interface{}) {
	stored := store.(map[string]interface{})
	roomd := room.(map[string]interface{})

	roomID := fmt.Sprintf("%d", int(roomd["id"].(float64)))
	roomName := fmt.Sprintf("%s", roomd["machineName"].(string))
	storeName := fmt.Sprintf("%s", stored["name"].(string))
	storeID := fmt.Sprintf("%d", int(stored["id"].(float64)))
	storeAddress := fmt.Sprintf("%s", stored["addressLat"].(string))
	price := fmt.Sprintf("%f", roomd["price"].(float64))
	status := fmt.Sprintf("%d", int(roomd["room"].(map[string]interface{})["status"].(float64)))

	fmt.Printf("time: %s, room_id: %s, room_name: %s, store_name: %s, store_id: %s, store_address: %s, price: %s, status: %s\n", timeh, roomID, roomName, storeName, storeID, storeAddress, price, status)

	// 插入数据
	statement, err := db.Prepare(
		`INSERT INTO rooms (time, room_id, room_name, store_name, store_id, store_address, price, status) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`)
	if err != nil {
		fmt.Println("ERROR: Failed to prepare the SQL statement")
		fmt.Println(err)
		return
	}
	_, err = statement.Exec(timeh, roomID, roomName, storeName, storeID, storeAddress, price, status)
	if err != nil {
		fmt.Println(err)
	}
}

func formatTime() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()

	return fmt.Sprintf("%d%d%d%d", year, month, day, hour)
}

func main() {
	db := initDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	timeHour := formatTime()
	stores := getStores()
	for _, store := range stores.([]interface{}) {
		store := store.(map[string]interface{})
		rooms := getRooms(store)
		for _, room := range rooms.([]interface{}) {
			insertRoom(db, timeHour, store, room)
		}
	}
}
