package repository

import (
	"encoding/json"
	"os"
	"test-mnc/model"
)

type LogActivityRepository struct {
	dbFile string
	log    []model.LogActivity
}

func NewLogActivityRepository(dbFile string) *LogActivityRepository {
	repo := &LogActivityRepository{dbFile: dbFile}
	repo.loadLog()
	return repo
}

func (r *LogActivityRepository) loadLog() []model.LogActivity {
	if _, err := os.Stat(r.dbFile); os.IsNotExist(err) {
		if err := os.WriteFile(r.dbFile, []byte("[]"), 0644); err != nil {
			return nil
		}
	}

	data, err := os.ReadFile(r.dbFile)
	if err != nil {
		return nil
	}

	var log []model.LogActivity
	if err := json.Unmarshal(data, &log); err != nil {
		return nil
	}

	return log
}

func (r *LogActivityRepository) AddLog(log model.LogActivity) error {
	r.log = append(r.log, log)
	return r.saveLog()
}

func (r *LogActivityRepository) saveLog() error {
	data, err := json.MarshalIndent(r.log, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.dbFile, data, 0644)
}

func (r *LogActivityRepository) GetLog() []model.LogActivity {
	return r.log
}
