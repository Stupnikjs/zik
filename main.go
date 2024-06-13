package main

import (
	"fmt"
	"log"
	"net/http"

	_ "cloud.google.com/go/storage"
	repo "github.com/Stupnikjs/zik/pkg/db"
	"github.com/joho/godotenv"
	_ "google.golang.org/api/option"
)

type Application struct {
	repo *repo.Dbrepo
	Port int
}

var BucketName string = "firstappbucknamezikapp"

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app := Application{
		Port: 3322,
		repo: &repo.PostgresRepo,
	}

	conn, err := app.repo.connectToDB()
	if err != nil {
		log.Fatalf("Error loading db conn: %v", err)
	}

	app.repo.InitTable()

	http.ListenAndServe(fmt.Sprintf(":%d", app.Port), app.routes())

}
