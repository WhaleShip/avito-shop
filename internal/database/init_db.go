package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Host                 string
	Port                 string
	Username             string
	Password             string
	DBName               string
	SSLMode              string
	PreferSimpleProtocol bool
}

func getEnvVariable(name string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		log.Fatalf("enviroment error: %s variable not exist", name)
	}
	return value
}

func GetInitializedDb() (*pgx.Conn, error) {
	conn, err := ConnectPostgres(Config{
		Host:                 getEnvVariable("POSTGRES_HOST"),
		Port:                 getEnvVariable("POSTGRES_PORT"),
		Username:             getEnvVariable("POSTGRES_USER"),
		Password:             getEnvVariable("POSTGRES_PASSWORD"),
		DBName:               getEnvVariable("POSTGRES_DB"),
		SSLMode:              getEnvVariable("SSL_MODE"),
		PreferSimpleProtocol: true,
	})
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("DB ping success")
	return conn, nil
}
