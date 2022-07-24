package service

import "push-api-cron/core/models"

type GroupService interface {
	Start(stopChan chan struct{}, data models.PushContent, interval int) error
	Stop(stopChan chan struct{})
}

type Service struct {
	GroupService
}

func NewService(groupService GroupService) *Service {
	return &Service{GroupService: NewGroupsService()}
}
