package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func clientPollingServer() {
	q := NewCappedQueue[string]()

	e := echo.New()
	e.GET("updates", func(c echo.Context) error {
		return c.JSON(200, q.Copy())
	})
	e.POST("send", func(c echo.Context) error {
		var request SendMessageRequest
		if err := c.Bind(&request); err != nil {
			return c.String(400, fmt.Sprintf("Bad request: %v", err))
		}
		q.Append(request.Message)
		return c.JSON(201, "I've sent your request.")
	})
	e.Logger.Fatal(e.Start(":8000"))
}
