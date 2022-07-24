package service

import (
	"bytes"
	"github.com/goccy/go-json"
	"net/http"
	"push-api-cron/core/models"
	"time"
)

const baseUrl = "https://push.api.appmetrica.yandex.net"

type GroupsService struct {
	client *http.Client
}

func NewGroupsService() *GroupsService {
	return &GroupsService{}
}

func (s *GroupsService) CreateGroup(input models.InputGroup) error {

	p, _ := json.Marshal(input)
	_, err := http.Post(baseUrl+"/push/v1/management/groups", "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	//TODO:добавить сохранение группы в базу
	return nil
}

//func (s *GroupsService) GetGroups(appId int) (models.InputGroup, error) {
//	//TODO:добавить возможность получать группы из БД, а также условную актуализацию данных раз в N время
//	response, err := http.Get(baseUrl + "/push/v1/management/groups?app_id=" + fmt.Sprint(appId))
//	if err != nil {
//		log.Println(err.Error())
//		return models.Group{}, err
//	}
//	defer response.Body.Close()
//	var result models.Group
//	err = json.NewDecoder(response.Body).Decode(&result)
//	if err != nil {
//		log.Println(err.Error())
//		return models.Group{}, err
//	}
//
//	return result, err
//}

//Пуш токены будут лежать а пуш контенте(также после согласовки некоторых моментов
func (s *GroupsService) Start(stopChan chan struct{}, data models.PushContent, interval int) error {
	bl := []byte{1, 2, 3}
	req, _ := http.NewRequest("POST", baseUrl+"/push/v1/send-batch", bytes.NewBuffer(bl))
	req.Header.Set("Authorization", "")
	go func() {

		//TODO:заменить на парсинг контента пуша
		for {
			select {
			case <-stopChan:
				return
			default:
				_, err := s.client.Do(req)
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

func (s *GroupsService) Stop(ch chan struct{}) {
	ch <- struct{}{}
}
