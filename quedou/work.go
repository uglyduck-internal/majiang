package quedou

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"majiang/utils"
	"net/http"
	"sync"
)

func initTable(db *sql.DB) {
	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS quedou_data ("year" TEXT,
"month" TEXT,
"day" TEXT, 
"hour" TEXT,
"room_id" TEXT,
"room_name" TEXT,
"store_name" TEXT, 
"store_id" TEXT,
"store_address" TEXT,
"price" TEXT,
"status" TEXT,
UNIQUE(year, month, day, hour, room_id)
);`)
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func GetRooms(store map[string]interface{}, token string) interface{} {
	storeID := int(store["id"].(float64))

	url := fmt.Sprintf("https://portal.quedou.cn/api/dms/device/list?storeId=%d&memberId=76965", storeID)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 MicroMessenger/6.8.0(0x16080000) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080710) XWEB/1191")
	req.Header.Set("Referer", "https://servicewechat.com/wx697f0b89354ff12e/26/page-frame.html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	getProxy := utils.WithRetryGetProxy(utils.GetProxy, 10)
	client, _ := getProxy()
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

func getStores(token string) interface{} {
	req, err := http.NewRequest("GET", "https://portal.quedou.cn/api/qsg/store/list", nil)
	if err != nil {
		fmt.Println("ERROR: Failed to create the HTTP request")
		return nil
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 MicroMessenger/6.8.0(0x16080000) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.7(0x13080710) XWEB/1191")
	req.Header.Set("Referer", "https://servicewechat.com/wx697f0b89354ff12e/26/page-frame.html")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	getProxy := utils.WithRetryGetProxy(utils.GetProxy, 10)
	client, _ := getProxy()
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

func insertRoom(db *sql.DB, datetime map[string]string, store interface{}, room interface{}) {
	stored := store.(map[string]interface{})
	roomd := room.(map[string]interface{})

	roomID := fmt.Sprintf("%d", int(roomd["id"].(float64)))
	roomName := fmt.Sprintf("%s", roomd["machineName"].(string))
	storeName := fmt.Sprintf("%s", stored["name"].(string))
	storeID := fmt.Sprintf("%d", int(stored["id"].(float64)))
	storeAddress := fmt.Sprintf("%s", stored["addressLat"].(string))
	price := fmt.Sprintf("%f", roomd["price"].(float64))
	statusCode := fmt.Sprintf("%d", int(roomd["room"].(map[string]interface{})["status"].(float64)))
	var status string
	if statusCode == "1" {
		status = "使用中"
	}
	// 插入数据
	statement, err := db.Prepare(
		`INSERT INTO quedou_data (year, month, day, hour, room_id, room_name, store_name, store_id, store_address, price, status) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)
	if err != nil {
		fmt.Println("ERROR: Failed to prepare the SQL statement")
		fmt.Println(err)
		return
	}

	year := datetime["year"]
	month := datetime["month"]
	day := datetime["day"]
	hour := datetime["hour"]
	statement.Exec(year, month, day, hour, roomID, roomName, storeName, storeID, storeAddress, price, status)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func StartWorkOnQuedou(db *sql.DB, datetime map[string]string, token string) {
	initTable(db)
	stores := getStores(token)
	var wg sync.WaitGroup
	for _, store := range stores.([]interface{}) {
		curStore := store
		wg.Add(1)
		go func() {
			handleStore(db, datetime, curStore, token)
			wg.Done()
		}()
	}
	wg.Wait()
}

func handleStore(db *sql.DB, datetime map[string]string, store interface{}, token string) {
	rooms := GetRooms(store.(map[string]interface{}), token)
	for _, room := range rooms.([]interface{}) {
		insertRoom(db, datetime, store, room)
	}
}
