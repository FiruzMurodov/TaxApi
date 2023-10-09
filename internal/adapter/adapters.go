package adapters

import (
	"encoding/json"
	"io"
	"net/http"
	"taxApi/internal/logs"
	"taxApi/internal/models"
	"time"
)

type Adapter struct {
	Client *http.Client
	Lg     *logs.Logger
	Cfg    models.Adapter
}

func NewAdapter(cfg *models.Config, lg *logs.Logger) *Adapter {
	client := http.Client{
		Timeout: time.Duration(cfg.Adapter.Timeout) * time.Second,
	}
	return &Adapter{Client: &client, Lg: lg, Cfg: cfg.Adapter}
}

func (a *Adapter) GetRoute(reqBody *models.RequestBody) (models.ResponseBody, error) {
	resp, err := a.Client.Get(a.Cfg.Url + reqBody.LongitudeSource + "," + reqBody.LatitudeSource + ";" + reqBody.LongitudeDestination + "," + reqBody.LatitudeDestination + "?sources=1&destinations=0&annotations=duration,distance")
	if err != nil {
		a.Lg.Error(err)
		return models.ResponseBody{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.Lg.Error(err)
		return models.ResponseBody{}, err
	}
	resp.Body.Close()

	var respBody models.ResponseBody
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		a.Lg.Error(err)
		return models.ResponseBody{}, err
	}

	return respBody, nil
}
