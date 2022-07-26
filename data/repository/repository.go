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

func NewRepository(db *sqlx.DB, cache *core.Cache) Repository {
	return Repository{db: db, cache: cache}
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
	s.cache.Set(device.Id, device, 0)
	var res string
	row := s.db.QueryRow("INSERT INTO devices (id, push_token) values ($1, $2) RETURNING id", device.Id, device.PushToken)
	if err := row.Scan(&res); err != nil {
		return err
	}
	return nil
}

func (s *Repository) GetAllDevices() []device.Device {
	if len(s.cache.GetAll()) != 0 {
		var res []device.Device
		for _, v := range s.cache.GetAll() {
			res = append(res, v.Value)
		}
		return res
	}
	var r []device.Device

	err := s.db.Select(&r, "SELECT * FROM devices")
	if err != nil {
		return nil
	}
	for _, v := range r {
		s.cache.Set(v.Id, v, 0)
	}
	return r
}
