package service

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"math"
	"net/http"
	"push-api-cron/core/models"
	"push-api-cron/core/models/device"
	"push-api-cron/data/repository"
	"sync"
	"time"
)

const baseUrl = "https://push.api.appmetrica.yandex.net/push"

type Service struct {
	repo   repository.Repository
	client *http.Client
}

func NewService(r repository.Repository) Service {
	return Service{
		repo:   r,
		client: &http.Client{},
	}
}

func (s *Service) UpdateToken(old, new string) {
	s.repo.UpdateToken(old, new)
}

func (s *Service) DeleteDevice(token string) {
	s.repo.DeleteDevice(token)
}

func (s *Service) CreateGroup(input models.InputGroup) (models.OutputGroup, error) {
	tmp := map[string]interface{}{
		"group": input,
	}
	p, _ := json.Marshal(tmp)
	buf := bytes.NewReader(p)
	req, _ := http.NewRequest("POST", "https://push.api.appmetrica.yandex.net/push/v1/management/groups", buf)
	tm := string(p)
	fmt.Println(tm)
	req.Header.Set("Authorization", "OAuth AQAAAAA16k9pAAhCeSFuHQpfykDPta5srg51zdw")
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return models.OutputGroup{}, err
	}
	defer resp.Body.Close()
	var b models.ResponseGroup
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&b)
	if err != nil {
		return models.OutputGroup{}, err
	}
	result, err := s.repo.CreateGroup(b.Group)
	if err != nil {
		return models.OutputGroup{}, err
	}
	return result, nil
}

func (s *Service) Start(stopChan chan struct{}, groupId int, data models.Messages, sendHour int) error {
	go func() {
		wg := sync.WaitGroup{}
		for {
			select {
			case <-stopChan:
				fmt.Println("all routines stopped")
				return
			default:
				dev := s.repo.GetAllDevices()
				for _, v := range dev {
					prepared := s.prepareData(groupId, data, v)
					go func(data []byte, timezone int) {
					again:
						wg.Add(1)
						var push models.Push
						json.Unmarshal(data, &push)
						utcHour := time.Now().UTC().Hour()
						localTime := utcHour + timezone
						if localTime >= 24 {
							localTime -= 24
						}
						temp := localTime - sendHour
						fmt.Println(temp, time.Now())
						if localTime != sendHour {
							if math.Abs(float64(localTime-sendHour)) > 2 {
								time.Sleep(1 * time.Hour)
								goto again
							} else {
								time.Sleep(30 * time.Minute)
								goto again
							}
						}
						tmp := bytes.NewReader(data)
						t := string(data)
						fmt.Println(t)
						req, _ := http.NewRequest("POST", baseUrl+"/v1/send-batch", tmp)
						req.Header.Set("Authorization", "OAuth AQAAAAA16k9pAAhCeSFuHQpfykDPta5srg51zdw")
						resp, err := s.client.Do(req)
						fmt.Println(resp)
						if err != nil {
							stopChan <- struct{}{}
						}
						body, _ := io.ReadAll(resp.Body)
						fmt.Println(string(body))
						wg.Done()

					}(prepared, v[0].TimeZone)
				}
				if len(dev) == 1 {
					time.Sleep(1 * time.Hour)
				}
			}
			wg.Wait()
		}
	}()
	return nil
}

func (s *Service) AddDevice(device device.Device) error {
	return s.repo.AddDevice(device)
}

func (s *Service) prepareData(group int, data models.Messages, devices []device.Device) []byte {
	var batches []models.Batch
	dev := s.FillDevices(devices)
	batches = append(batches, models.Batch{
		Messages: data,
		Device:   dev,
	})
	var push models.Push
	push = models.Push{models.NewPushBatchRequest(group, batches)}
	res, _ := json.Marshal(push)
	return res
}

func (s *Service) FillDevices(devices []device.Device) []models.Device {
	var result []models.Device
	var values []string
	for _, v := range devices {
		values = append(values, v.PushToken)
		result = append(result, models.Device{
			IdType:   "android_push_token",
			IdValues: values,
		})

	}

	return result
}

func (s *Service) Stop(ch chan struct{}) {
	ch <- struct{}{}
	return
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
