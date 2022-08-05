package repository

import (
	"github.com/jmoiron/sqlx"
	"push-api-cron/core"
	"push-api-cron/core/models"
	"push-api-cron/core/models/device"
)

type Repository struct {
	db    *sqlx.DB
	cache *core.Cache
}

func NewRepository(db *sqlx.DB) Repository {
	rep := Repository{db: db, cache: core.NewCache(0, 0)}
	rep.UpdateCache()
	return rep
}

func (s *Repository) GetGroup(id int) (models.OutputGroup, error) {
	var res models.OutputGroup
	err := s.db.Get(&res, "SELECT * FROM groups WHERE id=$1", id)
	if err != nil {
		return models.OutputGroup{}, nil
	}
	return res, nil
}

func (s *Repository) CreateGroup(input models.OutputGroup) (models.OutputGroup, error) {
	var id int
	row := s.db.QueryRow("INSERT INTO groups (id, app_id, name, send_rate) VALUES( $1, $2,$3,$4) RETURNING id", input.Id, input.AppId, input.Name, input.SendRate)
	if err := row.Scan(&id); err != nil {
		return models.OutputGroup{}, err
	}
	res, err := s.GetGroup(id)
	if err != nil {
		return models.OutputGroup{}, err
	}
	return res, nil
}

func (s *Repository) AddDevice(device device.Device) error {
	elem, _ := s.cache.Get(string(rune(device.TimeZone)))
	elem = append(elem, device)
	s.cache.Delete(string(rune(device.TimeZone)))
	s.cache.Set(string(rune(device.TimeZone)), elem, 0)
	var res string
	row := s.db.QueryRow("INSERT INTO devices (id, push_token) values ($1, $2) RETURNING id", device.Id, device.PushToken)
	if err := row.Scan(&res); err != nil {
		return err
	}
	return nil
}

func (s *Repository) UpdateDevice(pushToken string, id string) error {
	s.db.QueryRow("update devices set push_token=$1 where id=$2", pushToken, id)
	return nil
}

func (s *Repository) DeleteDevice(id string) error {
	s.db.QueryRow("delete from devices where id=$1", id)
	return nil
}

func (s *Repository) UpdateCache() {
	c := core.NewCache(0, 0)
	for _, v := range s.GetAllFromDb() {
		c.Set(string(rune(v[0].TimeZone)), v, 0)
	}
	s.cache = c
}

func (s *Repository) GetAllDevices() [][]device.Device {
	var res [][]device.Device
	for _, v := range s.cache.GetAll() {
		res = append(res, v.Value)
	}
	return res
}

func (s *Repository) GetAllFromDb() [][]device.Device {
	var times []int
	var result [][]device.Device
	err := s.db.Select(&times, "select time_zone from devices group by time_zone")
	if err != nil {
		return nil
	}
	for _, v := range times {
		var r []device.Device
		err := s.db.Select(&r, "SELECT * FROM devices where time_zone=$1", v)
		if err != nil {
			return nil
		}
		result = append(result, r)
		s.cache.Set(string(rune(v)), r, 0)

	}
	return result
}
