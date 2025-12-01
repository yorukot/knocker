package handler

import (
	"github.com/hibiken/asynq"
	"github.com/yorukot/knocker/repository"
)

type Handler struct {
	repo       repository.Repository
	notifier   *asynq.Client
	pingBuffer *PingRecorder
}

func NewHandler(repo repository.Repository, notifier *asynq.Client) *Handler {
	return &Handler{
		repo:       repo,
		notifier:   notifier,
		pingBuffer: NewPingRecorder(repo),
	}
}
