package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/crayon/wrap-bot/pkgs/logger"
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
	if c.apiKey == "" {
		logger.Error("[WeatherAPI] API Key is empty")
		return nil, fmt.Errorf("WeatherAPI error: API Key is not configured")
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", c.apiKey, city)
	logger.Info(fmt.Sprintf("[WeatherAPI] Requesting weather for city: %s", city))
	logger.Debug(fmt.Sprintf("[WeatherAPI] Request URL: %s", url))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to create request: %v", err))
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to execute request: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to read response body: %v", err))
		return nil, err
	}

	logger.Info(fmt.Sprintf("[WeatherAPI] Response status: %d", resp.StatusCode))
	logger.Debug(fmt.Sprintf("[WeatherAPI] Response body: %s", string(body)))

	if resp.StatusCode != 200 {
		logger.Error(fmt.Sprintf("[WeatherAPI] API returned error %d: %s", resp.StatusCode, string(body)))
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
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to parse JSON response: %v", err))
		logger.Error(fmt.Sprintf("[WeatherAPI] Response body: %s", string(body)))
		return nil, err
	}

	logger.Info(fmt.Sprintf("[WeatherAPI] Successfully retrieved weather for %s: %.1f°C, %s",
		weatherResp.Location.Name,
		weatherResp.Current.TempC,
		weatherResp.Current.Condition.Text))

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

	if c.apiKey == "" {
		logger.Error("[WeatherAPI] API Key is empty")
		return nil, fmt.Errorf("WeatherAPI error: API Key is not configured")
	}

	url := fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=%d&aqi=no", c.apiKey, city, days)
	logger.Info(fmt.Sprintf("[WeatherAPI] Requesting forecast for city: %s, days: %d", city, days))
	logger.Debug(fmt.Sprintf("[WeatherAPI] Request URL: %s", url))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to create request: %v", err))
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to execute request: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to read response body: %v", err))
		return nil, err
	}

	logger.Info(fmt.Sprintf("[WeatherAPI] Response status: %d", resp.StatusCode))
	logger.Debug(fmt.Sprintf("[WeatherAPI] Response body: %s", string(body)))

	if resp.StatusCode != 200 {
		logger.Error(fmt.Sprintf("[WeatherAPI] API returned error %d: %s", resp.StatusCode, string(body)))
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
		logger.Error(fmt.Sprintf("[WeatherAPI] Failed to parse JSON response: %v", err))
		logger.Error(fmt.Sprintf("[WeatherAPI] Response body: %s", string(body)))
		return nil, err
	}

	logger.Info(fmt.Sprintf("[WeatherAPI] Successfully retrieved forecast for %s: %d days",
		forecastResp.Location.Name,
		len(forecastResp.Forecast.Forecastday)))

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
