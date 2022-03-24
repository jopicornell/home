package main

import (
	"bufio"
	"github.com/jopicornell/thermonick/models"
	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
	"github.com/yryz/ds18b20"
	"io"
	"os"
	"time"
	//	"net/http"
)

type Status struct {
	HeaterOn    bool
	LightOn     bool
	BaskingTemp float64
	ColdTemp    float64
}

var sensorNameMap = map[string]string{
	"basking": "28-012032eb3734",
	"cold":    "28-012032cd28ce",
}

var currentStatus Status = Status{
	HeaterOn:    false,
	LightOn:     false,
	BaskingTemp: 0,
	ColdTemp:    0,
}

var HeaterLine *gpiod.Line
var HeaterLineNumber int = 14
var LightLine *gpiod.Line
var LightLineNumber int = 15
var Logger *log.Logger
var ChangesLogger *log.Logger

func main() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	fileChanges, err := os.OpenFile("changes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	stdAndFile := io.MultiWriter(file, os.Stdout)

	Logger = log.New()
	Logger.SetOutput(stdAndFile)
	ChangesLogger = log.New()
	ChangesLogger.SetOutput(fileChanges)

	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}
	Logger.Printf("sensor IDs: %v\n", sensors)
	currentBaskingTemp := getBaskingTemp()
	Logger.Printf("currentBaskingTemp: %f\n", currentBaskingTemp)
	ReadCommandsRoutine()
	for {
		baskingTemperature := getBaskingTemp()
		currentStatus.BaskingTemp = baskingTemperature
		if ShouldActivateHeater(baskingTemperature) {
			ActivateHeater(baskingTemperature)
		} else if ShouldDeactivateHeater(baskingTemperature) {
			DeactivateHeater(baskingTemperature)
		}
		if ShouldActivateLight() {
			ActivateLight(baskingTemperature)
		} else if ShouldDeactivateLight() {
			DeactivateLight(baskingTemperature)
		}
		currentStatus.ColdTemp = getColdTemp()
		time.Sleep(5 * time.Second)
	}
	//	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("Created temperatures. Response: %+v", response)
}

func ReadCommandsRoutine() {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			cmd, err := reader.ReadString('\n')
			if err != nil {
				Logger.Fatal(err)
			}
			if cmd == "status\n" {
				Logger.Infof("Current status \n Heater: %v \n Light: %v \n BaskingTemp: %f \n ColdTemp: %f \n", currentStatus.HeaterOn, currentStatus.LightOn, currentStatus.BaskingTemp, currentStatus.ColdTemp)
			}
		}

	}()
}

func DeactivateHeater(baskingTemperature float64) {
	ChangesLogger.Info("Deactivating heater temperature reached %.2f more than %.2f\n", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(0)
	if err != nil {
		Logger.Errorf("Error deactivating heater: %s", err.Error())
	}
	currentStatus.HeaterOn = false
}

func ActivateHeater(baskingTemperature float64) {
	ChangesLogger.Info("Activating heater temperature reached %.2f less than %.2f\n", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(1)
	if err != nil {
		Logger.Errorf("Error activating heater: %+v", err)
	}
	currentStatus.HeaterOn = true
}

func DeactivateLight(baskingTemperature float64) {
	ChangesLogger.Printf("Deactivating heater temperature reached %.2f more than %.2f\n", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getLightLine().SetValue(0)
	if err != nil {
		Logger.Errorf("Error deactivating heater: %s", err.Error())
	}
	currentStatus.LightOn = false
}

func ActivateLight(baskingTemperature float64) {
	ChangesLogger.Printf("Activating heater temperature reached %.2f less than %.2f\n", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getLightLine().SetValue(1)
	if err != nil {
		Logger.Errorf("Error activating heater: %+v", err)
	}
	currentStatus.LightOn = true
}

func getBaskingTemp() float64 {
	return getTemperature(sensorNameMap["basking"])
}

func getColdTemp() float64 {
	return getTemperature(sensorNameMap["cold"])
}

func getTemperature(sensorName string) float64 {
	temperature, err := ds18b20.Temperature(sensorName)
	if err != nil {
		Logger.Errorf("Error getting basking temperature: %s", err.Error())
		return 21
	}
	return temperature
}

func ShouldActivateHeater(temperature float64) bool {
	condition := getCurrentCondition()
	return !currentStatus.HeaterOn && condition.IsHeaterOn(time.Now(), temperature)
}

func ShouldDeactivateHeater(temperature float64) bool {
	condition := getCurrentCondition()
	return currentStatus.HeaterOn && !condition.IsHeaterOn(time.Now(), temperature)
}

func ShouldActivateLight() bool {
	condition := getCurrentCondition()
	return !currentStatus.HeaterOn && condition.IsLightOn(time.Now())
}

func ShouldDeactivateLight() bool {
	condition := getCurrentCondition()
	return currentStatus.HeaterOn && !condition.IsLightOn(time.Now())
}

func getCurrentCondition() models.Condition {
	return &models.SunsetSunriseCondition{
		MinimumHours:     9,
		MaximumHours:     12,
		TemperatureNight: 21,
		TemperatureDay:   31,
	}
}

func getDefaultGpioValue() int {
	condition := getCurrentCondition()
	if condition.IsHeaterOn(time.Now(), getBaskingTemp()) {
		return 1
	}
	return 0
}

func getHeaterLine() *gpiod.Line {
	if HeaterLine == nil {
		line, err := gpiod.RequestLine("gpiochip0", HeaterLineNumber, gpiod.AsOutput(getDefaultGpioValue()))
		if err != nil {
			Logger.Fatal(err)
		}
		HeaterLine = line
		return line
	}
	return HeaterLine
}

func getLightLine() *gpiod.Line {
	if LightLine == nil {
		line, err := gpiod.RequestLine("gpiochip0", LightLineNumber, gpiod.AsOutput(getDefaultGpioValue()))
		if err != nil {
			Logger.Fatal(err)
		}
		LightLine = line
		return line
	}
	return HeaterLine
}
