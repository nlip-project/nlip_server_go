package main

import (
	"nlip/handlers"

	"github.com/labstack/echo/v4"
)

func init() {
	// initializers.InitEnv()
	// initializers.InitDb()
	// initializers.SyncDB()
}

func main() {
	e := echo.New()

	e.POST("/nlip/", handlers.StartConversationHandler)
	e.POST("/text/", handlers.TextHandler)
	e.POST("/image/", handlers.ImageHandler)
	e.POST("/register/", handlers.Register)
	e.POST("/login/", handlers.Login)

	certFile := "/Users/hbzengin/src/go-server-example/druid.eecs.umich.edu.pem"
	keyFile := "/Users/hbzengin/src/go-server-example/druid.eecs.umich.edu-key.pem"
	// log.Fatal(http.ListenAndServeTLS(addr, certFile, keyFile, mux))
	e.Logger.Fatal(e.StartTLS(":80", certFile, keyFile))
}
