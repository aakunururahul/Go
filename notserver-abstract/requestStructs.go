package main

type SendMessageRequest struct {
	Message string `json:"message"`
}

type Message struct {
	CreatedAt int64
	Content   string
}

type GroupMessageRequest struct {
	Message   string `json:"message"`
	GroupName string `json:"groupName"`
}
