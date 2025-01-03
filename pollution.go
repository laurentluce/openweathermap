package openweathermap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// PollutionParameters holds the parameters needed to make
// a call to the current and forecast pollution API
type PollutionParameters struct {
	Location Coordinates
}

// HistoricalPollutionParameters holds the parameters needed to make
// a call to the historical pollution API
type HistoricalPollutionParameters struct {
	Location Coordinates
	Start    int64 // Data start (unix time, UTC time zone)
	End      int64 // Data end (unix time, UTC time zone)
}

// Pollution holds the data returnd from the pollution API
type Pollution struct {
	Dt       string          `json:"dt"`
	Location Coordinates     `json:"coord"`
	List     []PollutionData `json:"list"`
	Key      string
	*Settings
}

// PollutionData holds the pollution specific data from the call
type PollutionData struct {
	// Coord []float64 `json:"coord"`
	// List  []struct {
	Dt   int `json:"dt"`
	Main struct {
		Aqi float64 `json:"aqi"`
	} `json:"main"`
	Components struct {
		Co   float64 `json:"co"`
		No   float64 `json:"no"`
		No2  float64 `json:"no2"`
		O3   float64 `json:"o3"`
		So2  float64 `json:"so2"`
		Pm25 float64 `json:"pm2_5"`
		Pm10 float64 `json:"pm10"`
		Nh3  float64 `json:"nh3"`
	} `json:"components"`
	// } `json:"list"`
}

// NewPollution creates a new reference to Pollution
func NewPollution(key string, options ...Option) (*Pollution, error) {
	k, err := setKey(key)
	if err != nil {
		return nil, err
	}
	p := &Pollution{
		Key:      k,
		Settings: NewSettings(),
	}

	if err := setOptions(p.Settings, options); err != nil {
		return nil, err
	}
	return p, nil
}

// PollutionByParams gets the current pollution data based on the given parameters
func (p *Pollution) PollutionByParams(params *PollutionParameters) error {
	url := fmt.Sprintf(pollutionURL,
		p.Key,
		strconv.FormatFloat(params.Location.Latitude, 'f', -1, 64),
		strconv.FormatFloat(params.Location.Longitude, 'f', -1, 64),
	)
	response, err := p.client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return errInvalidKey
	}

	if err = json.NewDecoder(response.Body).Decode(&p); err != nil {
		return err
	}

	return nil
}

// ForecastPollutionByParams gets the forecast forecast pollution data based on the given parameters
func (p *Pollution) ForecastPollutionByParams(params *PollutionParameters) error {
	url := fmt.Sprintf(forecastPollutionURL,
		p.Key,
		strconv.FormatFloat(params.Location.Latitude, 'f', -1, 64),
		strconv.FormatFloat(params.Location.Longitude, 'f', -1, 64),
	)
	response, err := p.client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return errInvalidKey
	}

	if err = json.NewDecoder(response.Body).Decode(&p); err != nil {
		return err
	}

	return nil
}

// HistoricalPollutionByParams gets the historical pollution data based on the given parameters
func (p *Pollution) HistoricalPollutionByParams(params *HistoricalPollutionParameters) error {
	url := fmt.Sprintf(historicalPollutionURL,
		p.Key,
		strconv.FormatFloat(params.Location.Latitude, 'f', -1, 64),
		strconv.FormatFloat(params.Location.Longitude, 'f', -1, 64),
		params.Start,
		params.End,
	)
	response, err := p.client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return errInvalidKey
	}

	if err = json.NewDecoder(response.Body).Decode(&p); err != nil {
		return err
	}

	return nil
}
