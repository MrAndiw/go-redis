package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"youtube/redis/db"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// initialize redis
	db.RedisInit()

	// initialize echo webservice
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/insert", Insert)
	e.GET("/get", Get)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

type RespJson struct {
	Data   interface{}
	Status string
}

type RequestRedis struct {
	Name string
	Age  string
}

var key = "app_andi"

// Handler
func Insert(c echo.Context) error {
	id := c.QueryParam("id")
	name := c.QueryParam("name")
	age := c.QueryParam("age")

	// connection
	rdb := db.RedisConnect()

	reqRedis := RequestRedis{
		Name: name,
		Age:  age,
	}
	req, _ := json.Marshal(reqRedis)

	err := rdb.HSet(key, id, req).Err()
	if err != nil {
		return fmt.Errorf("error set redis %s", err)
	}

	// response
	resp := RespJson{
		Data:   id,
		Status: "Success",
	}

	return c.JSON(http.StatusOK, resp)
}

func Get(c echo.Context) error {
	id := c.QueryParam("id")

	// connection
	rdb := db.RedisConnect()

	val, err := rdb.HGet(key, id).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, "data tidak ditemukan")
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("error get redis %s", err.Error()))
	}

	var requestRedis RequestRedis
	err = json.Unmarshal([]byte(val), &requestRedis)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("error unmarshal redis %s", err.Error()))
	}

	// resp
	resp := RespJson{
		Data:   requestRedis,
		Status: "Success",
	}

	return c.JSON(http.StatusOK, resp)
}
