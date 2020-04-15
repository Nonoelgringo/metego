# MeteGo

MeteGo uses the openweather package to produce an up to 5 days weather 
forecast for a city that can be sent to a pushover application (more information on the Go pushover package [here](https://github.com/gregdel/pushover)).  

## openweather package

Wrapper around the 5 day / 3 hour forecast [OpenWeather API](https://openweathermap.org/api) for MeteGo.

## Usage

You need to create a *tokens.txt* file with your tokens where :
  * **openweather** is your openweather *API key*
  * **pushover** is the *API token* (Application)
  * **recipient** is your *User Key*
  
openweather is the only **mandatory** token. Pushover tokens are only required if you want to send to Pushover.
  
```
openweather youropenweathertoken
pushover yourpushovertoken
recipient yourrecipienttoken
```

### Options

```
-city <string>    City for the weather forecast
-days <int>       Number of forecasted days - between 1 and 5
-debug            Show debug lines
-pushover         Send the forecast to Pushover
```

## Ouput

Current day is included as part of the forecast.
Each line represents a 3 hours slice and is made of:
```
beginning hour | weather description - humidity - temperature(°C)
```
example : 
```
INFO    Weather forecast for Paris (2 day(s))
2020-04-15 
15:00:00 | few clouds - h:40% - 20.01°C
18:00:00 | clear sky - h:56% - 17.00°C
21:00:00 | clear sky - h:72% - 12.89°C

2020-04-16 
00:00:00 | clear sky - h:78% - 11.23°C
03:00:00 | scattered clouds - h:82% - 10.08°C
06:00:00 | broken clouds - h:77% - 10.52°C
09:00:00 | overcast clouds - h:54% - 15.89°C
12:00:00 | broken clouds - h:45% - 20.02°C
15:00:00 | scattered clouds - h:46% - 21.68°C
18:00:00 | scattered clouds - h:59% - 18.98°C
21:00:00 | light rain - h:72% - 15.58°C
```
