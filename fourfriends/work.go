package fourfriends

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"majiang/utils"
	"strings"
	"sync"
	"time"
)

func initTable(db *sql.DB) {
	statement, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS fourfriends_data ("year" TEXT, 
"month" TEXT,
"day" TEXT,
"hour" TEXT,
"city_code" TEXT,
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

func calculateMD5(input string) string {
	// 计算 MD5 哈希
	hash := md5.Sum([]byte(input))

	// 将哈希值转换为十六进制字符串
	md5Hash := fmt.Sprintf("%x", hash)

	return md5Hash
}

func GetResult(pageNum int, cityCode string) (interface{}, error) {
	// body:
	//oem_id	300ab330835844d58a8bccfc1c8b0800
	//lat	114.37408539862041
	//lng	30.452313752565868
	//store_name
	//city_code	420100
	//page	1
	//limit	5
	//is_included_qipai	1
	//is_included_billiards	0
	//is_experience	0
	//timestamp	1718350117211
	//api_version_interceptor	1
	//sign	4366673cbe1f1999905383c116e4ff11

	oem_id := "300ab330835844d58a8bccfc1c8b0800"
	lat := "114.37409153788359"
	lng := "30.45228049075855"
	store_name := ""
	city_code := cityCode
	page := fmt.Sprintf("%d", pageNum)
	limit := ""
	is_included_qipai := "1"
	is_included_billiards := "0"
	is_experience := "0"
	api_version_interceptor := "1"
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	combineStr := fmt.Sprintf("api_version_interceptor=%scity_code=%sis_experience=%sis_included_billiards=%sis_included_qipai=%slat=%slimit=%slng=%soem_id=%spage=%sstore_name=%stimestamp=%ssgpy@2023!hsjt05", api_version_interceptor, city_code, is_experience, is_included_billiards, is_included_qipai, lat, limit, lng, oem_id, page, store_name, timestamp)
	// calculate md5 of combineStr
	sign := calculateMD5(combineStr)

	bodyStr := fmt.Sprintf("oem_id=%s&lat=%s&lng=%s&store_name=%s&city_code=%s&page=%s&limit=%s&is_included_qipai=%s&is_included_billiards=%s&is_experience=%s&timestamp=%s&api_version_interceptor=%s&sign=%s", oem_id, lat, lng, store_name, city_code, page, limit, is_included_qipai, is_included_billiards, is_experience, timestamp, api_version_interceptor, sign)
	body := strings.NewReader(bodyStr)
	getProxy := utils.WithRetryGetProxy(utils.GetProxy, 10)
	client, _ := getProxy()
	resp, err := client.Post("https://iot.hs499.com/applet/user/selectStore", "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, fmt.Errorf("Failed to send the HTTP request: %s", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response body")
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the response body as JSON: %s, client: %s", respBody, client)
	}

	// check the response code
	returnCode := fmt.Sprintf("%f", result["code"].(float64))
	if strings.HasPrefix(returnCode, "5") {
		return nil, fmt.Errorf("Response code is: %s", returnCode)
	}
	return result, nil
}

func GetStores(pageNum int, cityCode string) interface{} {
	getResult := utils.WithRetryGetResult(GetResult, 10)
	result, _ := getResult(pageNum, cityCode)
	storeList := result.(map[string]interface{})["result"].(map[string]interface{})["store_list"]
	return storeList
}

func GetCities() interface{} {
	getResult := utils.WithRetryGetResult(GetResult, 10)
	result, _ := getResult(1, "")
	cityList := result.(map[string]interface{})["result"].(map[string]interface{})["open_city_list"]
	return cityList
}

func GetStoreStatus(storeDetail map[string]interface{}) float64 {
	return storeDetail["store_status"].(float64)
}

func GetRooms(store map[string]interface{}) (interface{}, error) {
	//api_version=1api_version_interceptor=1is_experience=0lat=lng=oem_id=300ab330835844d58a8bccfc1c8b0800store_id=56373a8f3f60448e8b11c4e3aef8e8c6timestamp=1718409036033user_id=
	apiVersion := "1"
	apiVersionInterceptor := "1"
	isExperience := "0"
	lat := ""
	lng := ""
	storeID := fmt.Sprintf("%s", store["store_id"])
	oemID := "300ab330835844d58a8bccfc1c8b0800"

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	userID := ""

	combineStr := fmt.Sprintf("api_version=%sapi_version_interceptor=%sis_experience=%slat=%slng=%soem_id=%sstore_id=%stimestamp=%suser_id=%ssgpy@2023!hsjt05", apiVersion, apiVersionInterceptor, isExperience, lat, lng, oemID, storeID, timestamp, userID)
	// calculate md5 of combineStr
	sign := calculateMD5(combineStr)
	bodyStr := fmt.Sprintf("oem_id=%s&lat=%s&lng=%s&store_id=%s&is_experience=%s&user_id=%s&api_version=%s&timestamp=%s&api_version_interceptor=%s&sign=%s", oemID, lat, lng, storeID, isExperience, userID, apiVersion, timestamp, apiVersionInterceptor, sign)
	body := strings.NewReader(bodyStr)
	getProxy := utils.WithRetryGetProxy(utils.GetProxy, 10)
	client, _ := getProxy()
	resp, err := client.Post("https://iot.hs499.com/applet/user/home", "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Println("ERROR: [GetRooms] Failed to send the HTTP request", err)
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read the response body")
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the response body as JSON: %s, client: %v", respBody, client)
	}
	if result["result"].(map[string]interface{})["room_list"] == nil {
		return nil, fmt.Errorf("No room list in the response body: %s", respBody)
	}
	return result["result"].(map[string]interface{})["room_list"], nil
}

func insertRoom(db *sql.DB, datetime map[string]string, cityCode string, store interface{}, room interface{}) {
	stored := store.(map[string]interface{})
	roomd := room.(map[string]interface{})

	roomID := fmt.Sprintf("%s", roomd["room_id"])
	roomName := fmt.Sprintf("%s", roomd["room_name"].(string))
	storeName := fmt.Sprintf("%s", stored["store_name"].(string))
	storeID := fmt.Sprintf("%s", stored["store_id"])
	storeAddress := fmt.Sprintf("%s", stored["store_address"].(string))
	price := fmt.Sprintf("%s", strings.Trim(roomd["room_price"].(string), "元"))
	statusCode := fmt.Sprintf("%d", int(roomd["room_status"].(float64))) // 2: 使用中
	var status string
	if statusCode == "2" {
		status = "使用中"
	}

	statement, err := db.Prepare(
		`INSERT INTO fourfriends_data (year, month, day, hour, city_code, room_id, room_name, store_name, store_id, store_address, price, status)
    				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`)
	if err != nil {
		fmt.Println("ERROR: Failed to prepare the SQL statement")
		fmt.Println(err)
		return
	}
	year := datetime["year"]
	month := datetime["month"]
	day := datetime["day"]
	hour := datetime["hour"]
	statement.Exec(year, month, day, hour, cityCode, roomID, roomName, storeName, storeID, storeAddress, price, status)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func StartWorkFourFriends(db *sql.DB, datetime map[string]string) {
	initTable(db)
	cities := GetCities()
	var wg sync.WaitGroup
	for _, city := range cities.([]interface{}) {
		cityCode := city.(map[string]interface{})["city_code"].(string)
		if cityCode != "430100" {
			continue
		}
		pageNum := 0
		for {
			pageNum++
			storeList := GetStores(pageNum, cityCode)
			if len(storeList.([]interface{})) == 0 {
				log.Printf("INFO: No more stores in city %s", cityCode)
				break
			}
			for _, store := range storeList.([]interface{}) {
				// status 1: 营业中, 2: 筹备中
				if GetStoreStatus(store.(map[string]interface{})) != float64(1) {
					continue
				}
				curCityCode := cityCode
				curStore := store
				wg.Add(1)
				go func() {
					handleStore(db, datetime, curCityCode, curStore)
					wg.Done()
				}()
			}
		}
	}
	wg.Wait()
}

func handleStore(db *sql.DB, datetime map[string]string, cityCode string, store interface{}) {
	getRooms := utils.WithRetryGetRooms(GetRooms, 10)
	rooms, err := getRooms(store.(map[string]interface{}))
	if err != nil {
		log.Printf("ERROR: [handleStore] Failed to get rooms for store %s", store.(map[string]interface{})["store_name"])
		panic(err)
	}

	for _, room := range rooms.([]interface{}) {
		insertRoom(db, datetime, cityCode, store, room)
	}
}
