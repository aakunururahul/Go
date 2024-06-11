package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func messageChannelServer() {
	e := echo.New()
	// simpleMessagesService(e)
	groupMessagesService(e)

	fmt.Println("Starting Server at :8001")
	e.Logger.Fatal(e.Start(":8001"))
}
