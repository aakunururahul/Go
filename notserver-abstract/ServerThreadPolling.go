package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func threadPollingServer() {
	q := NewCappedQueue[Message]()

	e := echo.New()

	e.GET("updates", func(c echo.Context) error {
		lastUpdate := c.QueryParam("lastUpdate")
		lastUpdateUnix, _ := strconv.ParseInt(lastUpdate, 10, 64)
		var updates []Message
		for {
			updates = getNewUpdates(q.Copy(), lastUpdateUnix)
			if len(updates) != 0 {
				break
			}
			select {
			case <-c.Request().Context().Done():
			case <-time.After(time.Second):
			}
		}
		return c.JSON(200, updates)
	})

	e.POST("send", func(c echo.Context) error {
		var request SendMessageRequest
		if err := c.Bind(&request); err != nil {
			return c.String(400, fmt.Sprintf("Bad request: %v", err))
		}
		q.Append(Message{
			CreatedAt: time.Now().Unix(),
			Content:   request.Message,
		})
		return c.JSON(201, "I've sent your request.")
	})
	e.Logger.Fatal(e.Start(":8001"))
}
