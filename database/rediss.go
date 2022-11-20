package database

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Employee struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func ConnectToRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	return client
	//json, err := json.Marshal(model.CredentialsJsonLess{Username: "aa", Password: "bb"})
	//if err != nil {
	//	fmt.Println(err)
	//}
}
func SetToRedis(client *redis.Client, json interface{}) {

	err := client.Set("aaa", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
func GetToRedis(client *redis.Client) string {
	val, err := client.Get("aaa").Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}
