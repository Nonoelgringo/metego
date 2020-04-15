package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Nonoelgringo/metego/openweather"
	"github.com/gregdel/pushover"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Args
	city := flag.String("city", "Paris", "Wanted city for the forecast")
	days := flag.Int("days", 5, "Wanted number of days to forecast")
	debug := flag.Bool("debug", false, "Set logging to debug level")
	sendPushover := flag.Bool("pushover", false, "Send the forecast to pushover")
	flag.Parse()

	// Logging
	atom := zap.NewAtomicLevel()
	if *debug {
		atom.SetLevel(zap.DebugLevel)
	}
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.TimeKey = ""
	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()
	sugar := logger.Sugar()

	// Get tokens
	tokens, err := getTokens("tokens.txt")
	if err != nil {
		sugar.Errorf("Error when getting tokens from file : %v\n", err)
		os.Exit(1)
	}

	// Initialize pushover, recipient and openweather
	pushoverApp := pushover.New((*tokens)["pushover"])
	recipient := pushover.NewRecipient((*tokens)["recipient"])
	openweatherApp, err := openweather.New((*tokens)["openweather"])
	if err != nil {
		sugar.Errorf("Error when creating an OpenWeather app : %v\n", err)
		os.Exit(1)
	}

	// Creates a new CityForecast instance
	cityForecast, err := openweather.NewCityForecast(*city, *days)
	if err != nil {
		sugar.Errorf("Error when creating a CityForecast : %v\n", err)
		os.Exit(1)
	}

	forecastSlice, err := openweatherApp.Forecast(cityForecast)
	if err != nil {
		sugar.Errorf("Error when getting Forecast : %v\n", err)
		os.Exit(1)
	}

	// Forecast title
	forecastTitle := fmt.Sprintf("Weather forecast for %s (%d day(s))", *city, *days)
	sugar.Infof(forecastTitle)

	// Debug prints
	sugar.Debugf("Forecast len='%v'", len(*forecastSlice))

	// Print forecast
	fmt.Println((*forecastSlice)[0])
	for _, v := range (*forecastSlice)[1:] {
		fmt.Println(v)
	}

	// Pushover
	if *sendPushover {
		sugar.Infof("Sending message to pushover")
		// Pushover message made of the first forecast entry + title
		message := pushover.NewMessageWithTitle((*forecastSlice)[0], forecastTitle)

		// Send message to recipient
		_, err = pushoverApp.SendMessage(message, recipient)
		if err != nil {
			sugar.Errorf("Error when sending message to Pushover : %v\n", err)
			os.Exit(1)
		}

		// Send the rest of the forecast (when days>1)
		for _, v := range (*forecastSlice)[1:] {
			message = pushover.NewMessage(v)
			_, err = pushoverApp.SendMessage(message, recipient)
			if err != nil {
				sugar.Errorf("Error when sending message to Pushover : %v\n", err)
				os.Exit(1)
			}
		}
	}
}

// getCredentials
func getTokens(filename string) (*map[string]string, error) {
	//
	tokensMap := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tokenName, token := sliceToStrings(strings.Split(scanner.Text(), " "))
		tokensMap[tokenName] = token
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &tokensMap, nil
}

func sliceToStrings(p []string) (string, string) {
	return p[0], p[1]
}
