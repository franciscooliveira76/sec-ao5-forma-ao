package main

import (
	"log"

	"github.com/sikozonpc/social/internal/env"
	"github.com/sikozonpc/social/internal/env/db"
	"github.com/sikozonpc/social/internal/env/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStorage(conn)

	defer conn.Close()

	if err := db.Seed(store); err != nil {
		log.Fatal(err)
	}
}
