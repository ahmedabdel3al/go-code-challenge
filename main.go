package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/speps/go-hashids/v2"
	"net/http"
	"strconv"
)

var ctx = context.Background()

func redisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func routers() *gin.Engine {
	r := gin.Default()

	r.GET("/:key", func(c *gin.Context) {
		redisClient := redisConnection()
		val, err := redisClient.Get(ctx, c.Param("key")).Result()

		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusMovedPermanently, val)
	})

	r.GET("/", func(c *gin.Context) {

		url := c.Query("url")
		lastValue := incrementUrl()
		lastId, _ := strconv.Atoi(lastValue)
		hashId := generateHashFromUrl(lastId)
		hashValue := setHash(hashId, url)

		c.JSON(200, gin.H{
			"shorterUrl": hashValue,
			"url":        url,
		})
	})

	return r
}

func incrementUrl() string {
	redisClient := redisConnection()
	redisClient.Incr(context.Background(), "last:id")
	val, err := redisClient.Get(context.Background(), "last:id").Result()

	if err != nil {
		panic(err)
	}

	return val
}

func generateHashFromUrl(lastId int) string {

	hd := hashids.NewData()
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	generatedId, _ := h.Encode([]int{lastId})

	return generatedId
}

func setHash(key string, value string) string {
	redisClient := redisConnection()
	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
	return key
}

func main() {
	router := routers()
	router.Run()
}
