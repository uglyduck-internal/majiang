package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type GetRoomsFunc func(store map[string]interface{}) (interface{}, error)

func WithRetryGetRooms(action GetRoomsFunc, maxRetries int) GetRoomsFunc {
	return func(store map[string]interface{}) (interface{}, error) {
		for i := 0; i < maxRetries; i++ {
			result, err := action(store)
			if err != nil {
				log.Printf("[Retry %d] Failed to [GetRooms]: %s, param: %v", i, err, store)
				time.Sleep(30 * time.Second)
				continue
			}
			return result, nil
		}
		panic(fmt.Errorf("Failed to get rooms after %d retries", maxRetries))
	}
}

type GetResultFunc func(int, string) (interface{}, error)

func WithRetryGetResult(action GetResultFunc, maxRetries int) GetResultFunc {
	return func(a int, b string) (interface{}, error) {
		for i := 0; i < maxRetries; i++ {
			result, err := action(a, b)
			if err != nil {
				log.Printf("[Retry %d] Failed to [GetResult]: %s", i, err)
				time.Sleep(30 * time.Second)
				continue
			}
			return result, nil
		}
		panic(fmt.Errorf("Failed to get result after %d retries", maxRetries))
	}
}

type GetProxyFunc func() (*http.Client, error)

func WithRetryGetProxy(action GetProxyFunc, maxRetries int) GetProxyFunc {
	return func() (*http.Client, error) {
		for i := 0; i < maxRetries; i++ {
			result, err := action()
			if err != nil {
				log.Printf("[Retry %d] Failed to [GetProxy]: %s", i, err)
				time.Sleep(30 * time.Second)
				continue
			}
			return result, nil
		}
		panic(fmt.Errorf("Failed to get proxy after %d retries", maxRetries))
	}
}
