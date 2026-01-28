package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/crayon/wrap-bot/pkgs/feature/ai/tool"
)

func RegisterWeatherTools(registry tool.ToolRegistry, client *WeatherAPIClient) error {
	tools := []tool.Tool{
		{
			Name:        "get_weather",
			Description: "获取当前天气",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"city": map[string]interface{}{
						"type":        "string",
						"description": "城市名称",
					},
				},
				"required": []string{"city"},
			},
			Handler: GetWeather(client),
			Enabled: true,
		},
		{
			Name:        "get_weather_forecast",
			Description: "获取天气预报",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"city": map[string]interface{}{
						"type":        "string",
						"description": "城市名称",
					},
					"days": map[string]interface{}{
						"type":        "integer",
						"description": "天数（3-7）",
					},
				},
				"required": []string{"city"},
			},
			Handler: GetWeatherForecast(client),
			Enabled: true,
		},
	}

	for _, t := range tools {
		if err := registry.Register(t); err != nil {
			return err
		}
	}
	return nil
}

func GetWeather(client *WeatherAPIClient) tool.ToolHandler {
	return func(ctx context.Context, args string) (string, error) {
		var params struct {
			City string `json:"city"`
		}

		json.Unmarshal([]byte(args), &params)

		weather, err := client.GetCurrentWeather(ctx, params.City)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("城市：%s\n温度：%.1f°C\n天气：%s\n湿度：%d%%\n风速：%.1f km/h",
			weather.City,
			weather.Temperature,
			weather.Condition,
			weather.Humidity,
			weather.WindSpeed,
		), nil
	}
}

func GetWeatherForecast(client *WeatherAPIClient) tool.ToolHandler {
	return func(ctx context.Context, args string) (string, error) {
		var params struct {
			City string `json:"city"`
			Days int    `json:"days"`
		}

		json.Unmarshal([]byte(args), &params)

		forecast, err := client.GetForecast(ctx, params.City, params.Days)
		if err != nil {
			return "", err
		}

		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("%s 未来%d天天气预报：\n\n", forecast.City, len(forecast.Days)))

		for i, day := range forecast.Days {
			if i > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(fmt.Sprintf("%s\n温度：%.1f°C ~ %.1f°C\n天气：%s",
				day.Date,
				day.Temperature.Min,
				day.Temperature.Max,
				day.Condition,
			))
		}

		return builder.String(), nil
	}
}
