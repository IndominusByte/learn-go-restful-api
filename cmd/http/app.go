package main

import (
	"fmt"
	"log"

	"github.com/IndominusByte/learn-go-restful-api/internal/config"
)

func startApp(cfg *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// connect the db
	db, err := config.DBConnect(cfg)
	if err != nil {
		return err
	}
	log.Printf("DB Connected")

	// connect redis
	redisCli, err := config.RedisConnect(cfg)
	if err != nil {
		return err
	}
	log.Println("Redis connected")

	fmt.Println("debug", db, redisCli)
	return
}
