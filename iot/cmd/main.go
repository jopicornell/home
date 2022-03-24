package main

import (
	"bufio"
	formatter "github.com/antonfisher/nested-logrus-formatter"
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
	Logger.SetFormatter(&formatter.Formatter{})
	ChangesLogger = log.New()
	ChangesLogger.SetOutput(fileChanges)

	sensors, err := ds18b20.Sensors()
	if err != nil {
		panic(err)
	}
	Logger.Infof("sensor IDs: %v", sensors)
	ReadCommandsRoutine()
	for {
		baskingTemperature, err := getBaskingTemp()
		if err != nil {
			Logger.Errorf("Error getting basking temperature: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		coldTemperature, err := getColdTemp()
		if err != nil {
			Logger.Errorf("Error getting basking temperature: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		ChangeStatusTemperatures(baskingTemperature, coldTemperature)
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

		currentStatus.ColdTemp = coldTemperature
		time.Sleep(time.Minute)
	}
	//	response, err := http.Post("https://umczz0pvpc.execute-api.eu-central-1.amazonaws.com/prod/temperatures", "application/json", bytes.NewBuffer(body))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("Created temperatures. Response: %+v", response)
}

func ChangeStatusTemperatures(baskingTemperature float64, coldTemperature float64) {
	if baskingTemperature != currentStatus.BaskingTemp {
		Logger.Infof("Changed baskingTemperature: %.2f", baskingTemperature)
		currentStatus.BaskingTemp = baskingTemperature
	}
	if coldTemperature != currentStatus.ColdTemp {
		Logger.Infof("Changed cold temperature: %.2f", coldTemperature)
		currentStatus.ColdTemp = coldTemperature
	}
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
	ChangesLogger.Infof("Deactivating heater temperature reached %.2f more than %.2f", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(0)
	if err != nil {
		Logger.Errorf("Error deactivating heater: %s", err.Error())
	}
	currentStatus.HeaterOn = false
}

func ActivateHeater(baskingTemperature float64) {
	ChangesLogger.Infof("Activating heater temperature reached %.2f less than %.2f", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getHeaterLine().SetValue(1)
	if err != nil {
		Logger.Errorf("Error activating heater: %+v", err)
	}
	currentStatus.HeaterOn = true
}

func DeactivateLight(baskingTemperature float64) {
	ChangesLogger.Printf("Deactivating heater temperature reached %.2f more than %.2f", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getLightLine().SetValue(0)
	if err != nil {
		Logger.Errorf("Error deactivating heater: %s", err.Error())
	}
	currentStatus.LightOn = false
}

func ActivateLight(baskingTemperature float64) {
	ChangesLogger.Printf("Activating heater temperature reached %.2f less than %.2f", baskingTemperature, getCurrentCondition().IdealTemperature(time.Now()))
	err := getLightLine().SetValue(1)
	if err != nil {
		Logger.Errorf("Error activating heater: %+v", err)
	}
	currentStatus.LightOn = true
}

func getBaskingTemp() (float64, error) {
	return getTemperature(sensorNameMap["basking"])
}

func getColdTemp() (float64, error) {
	return getTemperature(sensorNameMap["cold"])
}

func getTemperature(sensorName string) (float64, error) {
	return ds18b20.Temperature(sensorName)
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
	baskingTemperature, err := getBaskingTemp()
	if err != nil {
		Logger.Errorf("Error getting basking temperature: %s", err.Error())
		return 0
	}
	if condition.IsHeaterOn(time.Now(), baskingTemperature) {
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
