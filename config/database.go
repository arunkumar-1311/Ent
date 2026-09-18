package config

import (
	"log"

	"shipment-service/ent"

	_ "github.com/lib/pq"
)

func ConnectDatabase() *ent.Client {
	dsn := "host=localhost " +
		"port=5432 " +
		"user=postgres " +
		"password=password " +
		"dbname=rte " +
		"sslmode=disable " +
		"search_path=logistics"

	client, err := ent.Open("postgres", dsn)

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	return client
}
