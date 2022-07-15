package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"push-api-cron/core/models"
	"time"
)

const baseUrl = "https://push.api.appmetrica.yandex.net"

type GroupsService struct {
}

func NewGroupsService() *GroupsService {
	return &GroupsService{}
}

func (s *GroupsService) CreateGroup(body []byte) error {
	_, err := http.Post(baseUrl+"/push/v1/management/groups", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	//TODO:добавить сохранение группы в базу
	return nil
}

func (s *GroupsService) GetGroups(appId int) (models.Group, error) {
	//TODO:добавить возможность получать группы из БД, а также условную актуализацию данных раз в N время
	response, err := http.Get(baseUrl + "/push/v1/management/groups?app_id=" + fmt.Sprint(appId))
	if err != nil {
		log.Println(err.Error())
		return models.Group{}, err
	}
	defer response.Body.Close()
	var result models.Group
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Println(err.Error())
		return models.Group{}, err
	}

	return result, err
}

//Пуш токены будут лежать а пуш контенте(также после согласовки некоторых моментов
func (s *GroupsService) startGroup(stopChan chan struct{}, data models.PushContent, interval int) error {
	go func() {
		bl := []byte{1, 2, 3}
		//TODO:заменить на парсинг контента пуша
		for {
			select {
			case <-stopChan:
				return
			default:
				_, err := http.Post(baseUrl+"/push/v1/send-batch", "application/json", bytes.NewBuffer(bl))
				if err != nil {
					stopChan <- struct{}{}
				}
			}
			//TODO:возможно лучше принимать в минутах и конвертить в секунды
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()
	return nil
}

func (s *GroupsService) StopGroup(ch chan struct{}) {
	ch <- struct{}{}
}
