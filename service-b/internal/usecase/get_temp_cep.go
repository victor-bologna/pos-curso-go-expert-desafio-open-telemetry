package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/victor-bologna/pos-curso-go-expert-desafio-open-telemetry-b/internal/infra/dto"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

const name = "weather-usecase-b"

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
	logger = otelslog.NewLogger(name)
	cepCnt metric.Int64Counter
)

func init() {
	var err error
	cepCnt, err = meter.Int64Counter("weather.cep",
		metric.WithDescription("Get Weather by temperature"),
		metric.WithUnit("{cep}"))
	if err != nil {
		panic(err)
	}
}

func (ws WeatherService) Execute(ctx context.Context, cep string) (OutputTempDTO, error) {
	localidade, err := getLocalidade(ctx, cep)
	if err != nil {
		return OutputTempDTO{}, err
	}

	tempC, err := getTemperaturesInC(ctx, localidade)
	if err != nil {
		return OutputTempDTO{}, err
	}
	return createOutputDTO(localidade, tempC)
}

func getLocalidade(ctx context.Context, cep string) (string, error) {
	ctx, span := tracer.Start(ctx, "viaCep")
	defer span.End()

	logger.Info("Sending CEP " + cep + " to via cep service.")

	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Error while retrieving CEP from via cep: " + err.Error())
		return "", err
	}
	defer resp.Body.Close()

	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var viaCepResponse dto.ViaCepResponse
	err = json.Unmarshal(reader, &viaCepResponse)
	if err != nil {
		return "", err
	}
	if viaCepResponse.Erro {
		logger.Error("can not find zipcode")
		return "", errors.New("can not find zipcode")
	}
	logger.Info("ViaCepResponse:" + viaCepResponse.Cep + " - " + viaCepResponse.Localidade)

	viaCepAttr := attribute.String("cep", viaCepResponse.Localidade)
	span.SetAttributes(viaCepAttr)
	cepCnt.Add(ctx, 1, metric.WithAttributes(viaCepAttr))

	return viaCepResponse.Localidade, nil
}

func getTemperaturesInC(ctx context.Context, localidade string) (float64, error) {
	ctx, span := tracer.Start(ctx, "weatherAPI")
	defer span.End()

	encodedCity := url.QueryEscape(localidade)
	apiKey := "06140e1756914c55b5a213915242605"
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedCity)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	logger.Info("Getting weather from City: " + localidade)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Error while weather from weather api:" + err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var weatherApiResponse dto.WeatherApiResponse
	err = json.Unmarshal(reader, &weatherApiResponse)
	if err != nil {
		return 0, err
	}

	logger.Info("Weather api Response: " + fmt.Sprintf("%f", weatherApiResponse.Current.TempC))

	weatherApiAttr := attribute.Float64("Temperature in C", weatherApiResponse.Current.TempC)
	span.SetAttributes(weatherApiAttr)
	cepCnt.Add(ctx, 1, metric.WithAttributes(weatherApiAttr))

	return weatherApiResponse.Current.TempC, nil
}

func createOutputDTO(city string, tempC float64) (OutputTempDTO, error) {
	return OutputTempDTO{
		City:   city,
		Temp_C: tempC,
		Temp_F: (tempC*1.8 + 32),
		Temp_K: (tempC + 273),
	}, nil
}
