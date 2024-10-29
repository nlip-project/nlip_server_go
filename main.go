package main

import (
	"nlip/handlers"
	"os"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func init() {
	// initializers.InitEnv()
	// initializers.InitDb()
	// initializers.SyncDB()
}

func main() {
	e := echo.New()

	// Logging settings
	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetOutput(os.Stdout)
	e.Use(middleware.Logger())

	// These are unused for now.
	// e.POST("/nlip/", handlers.StartConversationHandler)
	// e.POST("/register/", handlers.Register)
	// e.POST("/login/", handlers.Login)

	e.POST("/", handlers.HandleIncomingMessage)

	certFile := "/Users/hbzengin/src/go-server-example/nlip.crt"
	keyFile := "/Users/hbzengin/src/go-server-example/nlip.key"
	e.Logger.Fatal(e.StartTLS(":80", certFile, keyFile))
}
