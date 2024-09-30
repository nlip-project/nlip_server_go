package main

import (
	"fmt"
	"log"
	"net/http"
	"nlip/initializers"
	"nlip/router"
)

func init() {
	initializers.InitEnv()
	initializers.InitDb()
	initializers.SyncDB()
}

func main() {
	mux := router.NewRouter()
	addr := "localhost:8080"
	fmt.Println("Starting the server on localhost:8080!")
	log.Fatal(http.ListenAndServe(addr, mux))
}
