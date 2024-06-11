package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func simpleMessagesService(e *echo.Echo) {
	q := NewCappedQueue[Message]()
	ps := NewPubSub()

	e.GET("updates", func(c echo.Context) error {
		lastUpdate := c.QueryParam("lastUpdate")
		lastUpdateUnix, _ := strconv.ParseInt(lastUpdate, 10, 64)

		// show it to user if we already have an update
		updates := getNewUpdates(q.Copy(), lastUpdateUnix)
		if len(updates) > 0 {
			return c.JSON(200, updates)
		}
		ch, close := ps.Subscribe()
		defer close()

		select {
		case <-ch:
			return c.JSON(200, getNewUpdates(q.Copy(), lastUpdateUnix))
		case <-c.Request().Context().Done():
			return c.String(http.StatusRequestTimeout, "timeout")
		}
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
		ps.Publish()
		return c.JSON(201, "I've sent your request.")
	})
}
