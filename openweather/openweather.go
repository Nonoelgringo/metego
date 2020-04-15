// Package openweather provides a simple wrapper around the openweather 5 day / 3 hour forecast API for metego
package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// APIEndpoint is the base endpoint for the OpenWeatfer 5 day / 3 hour forecast API
var APIEndpoint = "http://api.openweathermap.org/data/2.5/forecast"

// Custom http client
var myClient = &http.Client{Timeout: 10 * time.Second}

// OpenWeather errors
var (
	ErrEmptyToken = errors.New("openweather: empty token")
	ErrEmptyCity  = errors.New("openweather: empty City")
	ErrNumberDays = errors.New("openweather: days value should be between 1 and 5 (included)")
)

// OpenWeather is the representation of an app using the openweather API.
type OpenWeather struct {
	token string
}

// CityForecast represents a City and a number of Days for us to make a forecast on
type CityForecast struct {
	City string
	Days int
}

// JSONWeather is the top level JSON response from the API
type JSONWeather struct {
	List []Entry `json:"list"`
}

// Entry is unit of List from the JSONWeather struct
type Entry struct {
	DtTxt   string    `json:"dt_txt"`
	Main    Main      `json:"main"`
	Weather []Weather `json:"weather"`
}

// Main is part of the Entry struct
type Main struct {
	Humidity float32 `json:"humidity"`
	Temp     float32 `json:"temp"`
}

// Weather is part of the Entry struct
type Weather struct {
	Description string `json:"description"`
}

// New returns a new OpenWeather
func New(token string) (*OpenWeather, error) {
	if token == "" {
		return nil, ErrEmptyToken
	}
	return &OpenWeather{token}, nil
}

// NewCityForecast returns a new CityForecast
func NewCityForecast(city string, days int) (*CityForecast, error) {
	if city == "" {
		return nil, ErrEmptyCity
	}
	if days < 1 || days > 5 {
		return nil, ErrNumberDays
	}
	return &CityForecast{city, days}, nil
}

// Forecast returns a forecast message from a CityForecast
func (p *OpenWeather) Forecast(cf *CityForecast) (*[]string, error) {

	// API url (units=metric returns JSON response)
	url := fmt.Sprintf("%s?q=%s&units=metric&appid=%s", APIEndpoint, cf.City, p.token)

	// Decode the JSON response
	weatherJSON := new(JSONWeather)

	// Returns JSON from the API
	err := getJSON(url, weatherJSON)
	if err != nil {
		return nil, err
	}

	forecast, err := formatForecast(weatherJSON, cf.Days)
	if err != nil {
		return nil, err
	}

	return forecast, nil
}

// formatForecast return a formated forecast in a slice, each string in the slice represents a day forecast
func formatForecast(p *JSONWeather, days int) (*[]string, error) {
	// 24(h) / 3(h) * 5(d) = 40 entries maximum in the slice
	forecastLen := days * 8
	forecast := make([]string, 0, forecastLen)
	// Intermediate variables
	var daysCounter int
	var previousDate string
	var concatEntry string
	var dayForecast string
	// Looping over the JSON entries
	for _, entry := range p.List {
		// JSON variables
		dtTxt := entry.DtTxt
		date := dtTxt[:10]
		hour := dtTxt[11:]
		humidity := entry.Main.Humidity
		temperature := entry.Main.Temp
		for _, desc := range entry.Weather {
			description := desc.Description
			// when we reach a new date, append the day forecast to the slice + check the number of days passed (returns if so)
			// if its the same date, simply add the entry to the current day forecast
			if date != previousDate {
				daysCounter++
				concatEntry = fmt.Sprintf("%v \n%v | %v - h:%v%% - %.2f°C\n", date, hour, description, humidity, temperature)
				previousDate = date
				if dayForecast != "" {
					forecast = append(forecast, dayForecast)
				}
				if daysCounter > days {
					return &forecast, nil
				}
				dayForecast = concatEntry
			} else {
				concatEntry = fmt.Sprintf("%v | %v - h:%v%% - %.2f°C\n", hour, description, humidity, temperature)
				dayForecast += concatEntry
			}
		}
	}

	return nil, fmt.Errorf("formatForecast didnt end properly, shouldve returned earlier")
}

// getJSON gets and url and decode the response body into a target
func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
