package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/victor-bologna/pos-curso-go-expert-desafio-open-telemetry-a/internal/usecase"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var WeatherService usecase.WeatherUsecaseInterface

const name = "weather-handler-a"

var (
	tracer = otel.Tracer(name)
	logger = otelslog.NewLogger(name)
)

func GetTempByCep(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)
	cep := r.URL.Query().Get("cep")

	ctx, span := tracer.Start(ctx, name+"-"+cep)
	defer span.End()

	if len(cep) != 8 {
		msg := "invalid zipcode"
		logger.Error(name + " -> " + msg)
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}

	logger.Info("Sending CEP " + cep + " to Service B")

	response, err := WeatherService.Execute(ctx, cep)

	logger.Info("City: " + response.City + ", Temperature in C: " + fmt.Sprintf("%f", response.Temp_C))

	if err != nil {
		if err.Error() == "can not find zipcode" {
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
}
