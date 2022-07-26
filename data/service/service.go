package service

import (
	"bytes"
	"github.com/goccy/go-json"
	"net/http"
	"push-api-cron/core/models"
	"push-api-cron/core/models/device"
	"push-api-cron/data/repository"
	"time"
)

const baseUrl = "https://push.api.appmetrica.yandex.net"

type Service struct {
	repo   repository.Repository
	client http.Client
}

func NewService(r repository.Repository) *Service {
	return &Service{
		repo:   r,
		client: http.Client{},
	}
}

func (s *Service) CreateGroup(input models.InputGroup) error {

	p, _ := json.Marshal(input)
	resp, err := http.Post(baseUrl+"/push/v1/management/groups", "application/json", bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//TODO:читать ответ
	s.repo.CreateGroup(input)
	return nil
}

func (s *Service) Start(stopChan chan struct{}, groupId int, data models.Batch, interval int) error {
	go func() {
		for {
			select {
			case <-stopChan:
				return
			default:
				req, _ := http.NewRequest("POST", baseUrl+"/push/v1/send-batch", bytes.NewBuffer(s.prepareData(groupId, data)))
				req.Header.Set("Authorization", "OAuth AQAAAAA16k9pAAhCeSFuHQpfykDPta5srg51zdw")
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

func (s *Service) AddDevice(device device.Device) error {
	return s.repo.AddDevice(device)
}

func (s *Service) prepareData(group int, data models.Batch) []byte {
	newBatch := models.Batch{
		Messages: data.Messages,
		Device:   s.FillDevices(),
	}
	var b []models.Batch
	b = append(b, newBatch)
	push := models.Push{
		models.NewPushBatchRequest(group, b),
	}

	res, _ := json.Marshal(push)

	return res
}

func (s *Service) FillDevices() []models.Device {
	result := make([]models.Device, 5)
	dev := s.repo.GetAllDevices()
	for _, v := range dev {
		result = append(result, models.Device{
			IdType:   "android_push_token",
			IdValues: v.PushToken,
		})
	}

	return result
}

func (s *Service) Stop(ch chan struct{}) {
	ch <- struct{}{}
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
