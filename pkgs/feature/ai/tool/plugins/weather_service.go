package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WeatherAPIClient struct {
	apiKey string
	client *http.Client
}

func NewWeatherAPIClient(apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *WeatherAPIClient) GetCurrentWeather(ctx context.Context, city string) (*Weather, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", c.apiKey, city)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("WeatherAPI error %d: %s", resp.StatusCode, string(body))
	}

	var weatherResp struct {
		Current struct {
			TempC     float64 `json:"temp_c"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
			WindKmph float64 `json:"wind_kph"`
			Humidity int     `json:"humidity"`
		} `json:"current"`
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
	}

	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, err
	}

	return &Weather{
		City:        weatherResp.Location.Name,
		Temperature: weatherResp.Current.TempC,
		Condition:   weatherResp.Current.Condition.Text,
		WindSpeed:   weatherResp.Current.WindKmph,
		Humidity:    weatherResp.Current.Humidity,
	}, nil
}

func (c *WeatherAPIClient) GetForecast(ctx context.Context, city string, days int) (*Forecast, error) {
	if days == 0 {
		days = 3
	}
	if days > 7 {
		days = 7
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=%d&aqi=no", c.apiKey, city, days)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("WeatherAPI error %d: %s", resp.StatusCode, string(body))
	}

	var forecastResp struct {
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
		Forecast struct {
			Forecastday []struct {
				Date string `json:"date"`
				Day  struct {
					MaxtempC  float64 `json:"maxtemp_c"`
					MintempC  float64 `json:"mintemp_c"`
					Condition struct {
						Text string `json:"text"`
					} `json:"condition"`
				} `json:"day"`
			} `json:"forecastday"`
		} `json:"forecast"`
	}

	if err := json.Unmarshal(body, &forecastResp); err != nil {
		return nil, err
	}

	result := &Forecast{
		City: forecastResp.Location.Name,
		Days: make([]DayWeather, 0, len(forecastResp.Forecast.Forecastday)),
	}

	for _, day := range forecastResp.Forecast.Forecastday {
		result.Days = append(result.Days, DayWeather{
			Date: day.Date,
			Temperature: struct {
				Min float64
				Max float64
			}{
				Min: day.Day.MintempC,
				Max: day.Day.MaxtempC,
			},
			Condition: day.Day.Condition.Text,
		})
	}

	return result, nil
}

type Weather struct {
	City        string
	Temperature float64
	Condition   string
	Humidity    int
	WindSpeed   float64
}

type Forecast struct {
	City string
	Days []DayWeather
}

type DayWeather struct {
	Date        string
	Temperature struct {
		Min float64
		Max float64
	}
	Condition string
}
