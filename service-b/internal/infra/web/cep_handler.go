package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/victor-bologna/pos-curso-go-expert-desafio-open-telemetry-b/internal/usecase"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
)

var WeatherService usecase.WeatherUsecaseInterface

const name = "weather-handler-b"

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

func GetTempByCep(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)

	cep := r.URL.Query().Get("cep")

	ctx, span := tracer.Start(ctx, name+"-"+cep)
	defer span.End()

	if len(cep) != 8 {
		msg := "invalid zipcode"
		logger.Error(name + " -> " + msg)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	response, err := WeatherService.Execute(ctx, cep)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "can not find zipcode") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attrSet := attribute.NewSet(
		attribute.Float64("Temperature in C", response.Temp_C),
		attribute.Float64("Temperature in F", response.Temp_F),
		attribute.Float64("Temperature in K", response.Temp_K),
		attribute.String("City", response.City))
	span.SetAttributes(attrSet.ToSlice()...)
	cepCnt.Add(ctx, 1, metric.WithAttributeSet(attrSet))
}
