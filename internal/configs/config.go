package configs

import (
	"encoding/json"
	"os"
	"taxApi/internal/models"
)

func InitConfigs() (*models.Config, error) {
	bytes, err := os.ReadFile("internal/configs/config.json")
	if err != nil {
		return nil, err
	}
	var config models.Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
