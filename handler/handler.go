package handler

import "github.com/sashabaranov/go-openai"

type Handler struct {
	*openai.Client
}

func NewHandler(openaiKey string) *Handler {
	handler := &Handler{
		openai.NewClient(openaiKey),
	}
	return handler
}
