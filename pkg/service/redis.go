package database

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Employee struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

func SetToRedis(json interface{}) {

	err := client.Set("aaa", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
func GetToRedis() string {
	val, err := client.Get("aaa").Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}
