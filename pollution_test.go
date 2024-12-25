package openweathermap

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"
)

// TestNewPollution
func TestNewPollution(t *testing.T) {

	p, err := NewPollution(os.Getenv("OWM_API_KEY"))
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(p).String() != "*openweathermap.Pollution" {
		t.Error("incorrect data type returned")
	}
}

// TestNewPollution with custom http client
func TestNewPollutionWithCustomHttpClient(t *testing.T) {

	hc := http.DefaultClient
	hc.Timeout = time.Duration(1) * time.Second
	p, err := NewPollution(os.Getenv("OWM_API_KEY"), WithHttpClient(hc))
	if err != nil {
		t.Error(err)
	}

	if reflect.TypeOf(p).String() != "*openweathermap.Pollution" {
		t.Error("incorrect data type returned")
	}

	expected := time.Duration(1) * time.Second
	if p.client.Timeout != expected {
		t.Errorf("Expected Duration %v, but got %v", expected, p.client.Timeout)
	}
}

// TestNewPollutionWithInvalidOptions will verify that returns an error with
// invalid option
func TestNewPollutionWithInvalidOptions(t *testing.T) {

	optionsPattern := [][]Option{
		{nil},
		{nil, nil},
		{WithHttpClient(&http.Client{}), nil},
		{nil, WithHttpClient(&http.Client{})},
	}

	for _, options := range optionsPattern {
		c, err := NewPollution(os.Getenv("OWM_API_KEY"), options...)
		if err == errInvalidOption {
			t.Logf("Received expected invalid option error. message: %s", err.Error())
		} else if err != nil {
			t.Errorf("Expected %v, but got %v", errInvalidOption, err)
		}
		if c != nil {
			t.Errorf("Expected nil, but got %v", c)
		}
	}
}

// TestNewPollutionWithInvalidHttpClient will verify that returns an error with
// invalid http client
func TestNewPollutionWithInvalidHttpClient(t *testing.T) {

	p, err := NewPollution(os.Getenv("OWM_API_KEY"), WithHttpClient(nil))
	if err == errInvalidHttpClient {
		t.Logf("Received expected bad client error. message: %s", err.Error())
	} else if err != nil {
		t.Errorf("Expected %v, but got %v", errInvalidHttpClient, err)
	}
	if p != nil {
		t.Errorf("Expected nil, but got %v", p)
	}
}

// TestPollutionByParams tests the call to the current pollution API
func TestPollutionByParams(t *testing.T) {
	t.Parallel()
	p, err := NewPollution(os.Getenv("OWM_API_KEY"))
	if err != nil {
		t.Error(err)
	}
	params := &PollutionParameters{
		Location: Coordinates{
			Latitude:  0.0,
			Longitude: 10.0,
		},
	}
	if err := p.PollutionByParams(params); err != nil {
		t.Error(err)
	}
}

// TestForecastPollutionByParams tests the call to the forecast pollution API
func TestForecastPollutionByParams(t *testing.T) {
	t.Parallel()
	p, err := NewPollution(os.Getenv("OWM_API_KEY"))
	if err != nil {
		t.Error(err)
	}
	params := &PollutionParameters{
		Location: Coordinates{
			Latitude:  0.0,
			Longitude: 10.0,
		},
	}
	if err := p.ForecastPollutionByParams(params); err != nil {
		t.Error(err)
	}
}

// TestHistoricalPollutionByParams tests the call to the historical pollution API
func TestHistoricalPollutionByParams(t *testing.T) {
	t.Parallel()
	p, err := NewPollution(os.Getenv("OWM_API_KEY"))
	if err != nil {
		t.Error(err)
	}
	params := &HistoricalPollutionParameters{
		Location: Coordinates{
			Latitude:  0.0,
			Longitude: 10.0,
		},
		Start: 1606223802,
		End:   1606482999,
	}
	if err := p.HistoricalPollutionByParams(params); err != nil {
		t.Error(err)
	}
	t.Logf("p list: %v\n", p.List)
	for _, pollutionData := range p.List {
		t.Logf("Pollution data Aqi: %f\n", pollutionData.Main.Aqi)
	}
}
