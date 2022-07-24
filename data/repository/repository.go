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

func NewRepository(db *sqlx.DB, cache *core.Cache) *Repository {
	return &Repository{db: db, cache: cache}
}

func (s *Repository) GetGroup(id int) (models.OutputGroup, error) {
	var res models.OutputGroup
	err := s.db.Get(&res, "SELECT * FROM groups WHERE id=$1", id)
	if err != nil {
		return models.OutputGroup{}, nil
	}
	return res, nil
}

func (s *Repository) CreateGroup(input models.InputGroup) (models.OutputGroup, error) {
	var id int
	row := s.db.QueryRow("INSERT INTO groups (app_id, name, send_rate) VALUES $1,$2,$3 RETURNING id", input.AppId, input.Name, input.SendRate)
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
	s.cache.Set("", device, 0)
	return nil
	//TODO:заполнить сущность девайсов и добавлять в бд
}
