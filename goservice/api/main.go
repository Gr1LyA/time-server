package main

import (
	"api/internal/app/apiserver"
	"api/internal/app/redis"
	"fmt"
	"os"
)

func main() {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	cache, err := redis.NewRedisClient(addr, os.Getenv("REDIS_PASSWORD"), 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	apiserver.Run(cache)
}
