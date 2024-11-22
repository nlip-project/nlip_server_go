package main

import (
	"nlip/auth"
	"nlip/handlers"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func init() {
	// initializers.InitDb()
	// initializers.SyncDB()
}

func main() {

	e := echo.New()

	e.Pre(middleware.AddTrailingSlash())

	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetOutput(os.Stdout)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/nlip/", handlers.HandleIncomingMessage)
	e.POST("/upload/", handlers.UploadHandler)

	e.GET("/login/", auth.HandleLogin)
	e.GET("/auth/", auth.HandleCallback)
	e.GET("/protected/", auth.ProtectedHandler)

	certFile := "/Users/hbzengin/src/go-server-example/nlip.crt"
	keyFile := "/Users/hbzengin/src/go-server-example/nlip.key"
	e.Logger.Fatal(e.StartTLS(":80", certFile, keyFile))
}
