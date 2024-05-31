package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/victor-bologna/pos-curso-go-expert-desafio-open-telemetry-a/internal/usecase"
)

type MockUsecase struct {
	mock.Mock
}

func (m *MockUsecase) Execute(ctx context.Context, cep string) (usecase.OutputTempDTO, error) {
	args := m.Called(cep)
	return args.Get(0).(usecase.OutputTempDTO), args.Error(1)
}

func TestGetTempByCep_InvalidZipcode(t *testing.T) {
	mockUsecase := new(MockUsecase)
	WeatherService = mockUsecase

	req, err := http.NewRequest("GET", "/temp?cep=123", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTempByCep)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, "invalid zipcode\n", rr.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestGetTempByCep_ZipcodeNotFound(t *testing.T) {
	mockUsecase := new(MockUsecase)
	WeatherService = mockUsecase

	mockUsecase.On("Execute", "12345678").Return(usecase.OutputTempDTO{}, errors.New("can not find zipcode"))

	req, err := http.NewRequest("GET", "/temp?cep=12345678", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTempByCep)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "can not find zipcode\n", rr.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestGetTempByCep_InternalServerError(t *testing.T) {
	mockUsecase := new(MockUsecase)
	WeatherService = mockUsecase

	mockUsecase.On("Execute", "12345678").Return(usecase.OutputTempDTO{}, errors.New("some internal error"))

	req, err := http.NewRequest("GET", "/temp?cep=12345678", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTempByCep)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "some internal error\n", rr.Body.String())

	mockUsecase.AssertExpectations(t)
}

func TestGetTempByCep_SuccessfulResponse(t *testing.T) {
	mockUsecase := new(MockUsecase)
	WeatherService = mockUsecase

	mockUsecase.On("Execute", "12345678").Return(usecase.OutputTempDTO{
		City:   "São Paulo",
		Temp_C: 25.0,
		Temp_F: 77.0,
		Temp_K: 298.0,
	}, nil)

	req, err := http.NewRequest("GET", "/temp?cep=12345678", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTempByCep)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := `{"city":"São Paulo","temp_C":25,"temp_F":77,"temp_K":298}` + "\n"
	assert.Equal(t, expectedBody, rr.Body.String())

	mockUsecase.AssertExpectations(t)
}
