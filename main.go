package main

import (
	"fmt"
	"log"
	"net/http"
	"nlip/router"
)

func init() {
	// initializers.InitEnv()
	// initializers.InitDb()
	// initializers.SyncDB()
}

func main() {
	mux := router.NewRouter()
	addr := "0.0.0.0:80"
	fmt.Printf("Starting the HTTPS server on %s!", addr)
	// TODO: Paths here are absolute paths! They are here because
	// systemctl needs to have paths to this! Maybe use .env file later.
	certFile := "/Users/hbzengin/src/go-server-example/localhost+2.pem"
	keyFile := "/Users/hbzengin/src/go-server-example/localhost+2-key.pem"
	log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, mux))
}
