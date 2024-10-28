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

	// This is unused for now.
	// e.POST("/nlip/", handlers.StartConversationHandler)
	e.POST("/", handlers.HandleIncomingMessage)
	e.POST("/register/", handlers.Register)
	e.POST("/login/", handlers.Login)

	certFile := "/Users/hbzengin/src/go-server-example/nlip.crt"
	keyFile := "/Users/hbzengin/src/go-server-example/nlip.key"
	e.Logger.Fatal(e.StartTLS(":80", certFile, keyFile))
}
