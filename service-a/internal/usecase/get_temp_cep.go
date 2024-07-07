package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type WeatherUsecaseInterface interface {
	Execute(ctx context.Context, cep string) (OutputTempDTO, error)
}

type WeatherService struct{}

type OutputTempDTO struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_C"`
	Temp_F float64 `json:"temp_F"`
	Temp_K float64 `json:"temp_K"`
}

func (ws WeatherService) Execute(ctx context.Context, cep string) (OutputTempDTO, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://service-b:8081/temperature?cep="+cep, nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	if err != nil {
		return OutputTempDTO{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return OutputTempDTO{}, err
	}
	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return OutputTempDTO{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return OutputTempDTO{}, errors.New("can not find zipcode")
	}

	var outputTempDTO OutputTempDTO
	json.Unmarshal(reader, &outputTempDTO)
	return outputTempDTO, nil
}
