package env

import (
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"shorterUrl/src/storage"
	"strconv"
)

const (
	keyRedisAddr string = "APP_REDIS_ADDRESS"
	keyRedisPwd  string = "APP_REDIS_PASSWORD"
	keyRedisDB   string = "APP_REDIS_DB"
)

type Env struct {
	S storage.Storage
}

func GetEnv() *Env {
	addr := os.Getenv(keyRedisAddr)
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	pwd := os.Getenv(keyRedisPwd)
	if pwd == "" {
		pwd = ""
	}

	dbs := os.Getenv(keyRedisDB)
	if dbs == "" {
		dbs = "0"
	}

	db, err := strconv.Atoi(dbs)
	if err != nil {
		log.Fatal(err)
	}

	r := storage.NewRedisClient(addr, pwd, db)
	log.Printf("connect to redis (addr : %s ; pwd : %s ; db : %d\n", addr, pwd, db)

	return &Env{S: r}
}
