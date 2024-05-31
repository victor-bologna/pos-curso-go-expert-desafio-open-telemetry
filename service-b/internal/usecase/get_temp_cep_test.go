package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockWeatherService struct{}

func (m MockWeatherService) Execute(cep string) (OutputTempDTO, error) {
	if cep == "12345678" {
		return OutputTempDTO{
			City:   "São Paulo",
			Temp_C: 20.0,
			Temp_F: 68.0,
			Temp_K: 293.15,
		}, nil
	}
	return OutputTempDTO{}, errors.New("invalid cep")
}

func TestWeatherService_Execute(t *testing.T) {
	mockService := MockWeatherService{}

	t.Run("valid cep", func(t *testing.T) {
		output, err := mockService.Execute("12345678")
		assert.NoError(t, err)
		assert.Equal(t, "São Paulo", output.City)
		assert.Equal(t, 20.0, output.Temp_C)
		assert.Equal(t, 68.0, output.Temp_F)
		assert.Equal(t, 293.15, output.Temp_K)
	})

	t.Run("invalid cep", func(t *testing.T) {
		_, err := mockService.Execute("invalid")
		assert.Error(t, err)
	})
}
