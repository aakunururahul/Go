package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func groupMessagesService(e *echo.Echo) {
	allGroups := make(map[string]*CappedQueue[Message])

	allGroupsAliveRequests := make(map[string]*PubSub)

	e.GET("groupUpdates", func(c echo.Context) error {
		lastUpdate := c.QueryParam("lastUpdate")
		groupName := c.QueryParam("groupName")

		if _, ok := allGroups[groupName]; !ok {
			allGroups[groupName] = NewCappedQueue[Message]()
		}
		if _, ok := allGroupsAliveRequests[groupName]; !ok {
			allGroupsAliveRequests[groupName] = NewPubSub()
		}

		lastUpdateUnix, _ := strconv.ParseInt(lastUpdate, 10, 64)

		// show it to user if we already have an update
		q := allGroups[groupName]
		updates := getNewUpdates(q.Copy(), lastUpdateUnix)
		if len(updates) > 0 {
			return c.JSON(200, updates)
		}

		ps := allGroupsAliveRequests[groupName]
		ch, close := ps.Subscribe()
		defer close()

		select {
		case <-ch:
			return c.JSON(200, getNewUpdates(q.Copy(), lastUpdateUnix))
		case <-c.Request().Context().Done():
			return c.String(http.StatusRequestTimeout, "timeout")
		}
	})

	e.POST("sendToGroup", func(c echo.Context) error {
		var request GroupMessageRequest
		if err := c.Bind(&request); err != nil {
			return c.String(400, fmt.Sprintf("Bad request: %v", err))
		}
		groupName := request.GroupName

		if _, ok := allGroups[groupName]; !ok {
			allGroups[groupName] = NewCappedQueue[Message]()
		}
		if _, ok := allGroupsAliveRequests[groupName]; !ok {
			allGroupsAliveRequests[groupName] = NewPubSub()
		}

		q := allGroups[request.GroupName]
		createdAt := time.Now().Unix()
		q.Append(Message{
			CreatedAt: createdAt,
			Content:   request.Message,
		})

		ps := allGroupsAliveRequests[request.GroupName]
		ps.Publish()
		return c.JSON(201, fmt.Sprintf("I've sent your request at %v", createdAt))
	})
}
