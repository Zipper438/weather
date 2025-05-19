package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type WeatherClient struct {
	Client  HTTPClient
	BaseURL string
	APIKey  string
}

func NewWeatherClient(apiKey string) *WeatherClient {
	return &WeatherClient{
		Client:  &http.Client{Timeout: 8 * time.Second},
		BaseURL: "http://api.weatherapi.com/v1/current.json",
		APIKey:  apiKey,
	}
}

func (c *WeatherClient) GetWeather(city string) (float64, error) {
	escapedCity := url.QueryEscape(city)

	url := fmt.Sprintf("%s?key=%s&q=%s", c.BaseURL, c.APIKey, escapedCity)

	resp, err := c.Client.Get(url)
	if err != nil {
		return float64(-274), fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return float64(-264), fmt.Errorf("сервер вернул код ошибки: %d", resp.StatusCode)
	}

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil && !os.IsNotExist(err) {
		return float64(-274), fmt.Errorf("ошибка парсинга из JSON: %v", err)
	}

	return data.Current.TempC, nil
}
