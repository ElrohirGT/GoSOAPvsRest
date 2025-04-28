package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const JSON_RESPONSE = `
[
  {
    "city": "Bogotá",
    "temperature": 18.3,
    "condition": "Lluvia ligera"
  },
  {
    "city": "Medellín",
    "temperature": 24.1,
    "condition": "Soleado"
  },
  {
    "city": "Cali",
    "temperature": 29.0,
    "condition": "Nublado"
  },
  {
    "city": "Cartagena",
    "temperature": 31.7,
    "condition": "Caluroso y húmedo"
  },
  {
    "city": "Barranquilla",
    "temperature": 30.4,
    "condition": "Soleado con brisa"
  },
  {
    "city": "Pasto",
    "temperature": 14.2,
    "condition": "Frío y lluvioso"
  },
  {
    "city": "Manizales",
    "temperature": 17.0,
    "condition": "Niebla parcial"
  },
  {
    "city": "Armenia",
    "temperature": 22.5,
    "condition": "Clima templado"
  }
]`

const XML_RESPONSE = `
<WeatherReports>
  <WeatherReport>
    <city>Bogotá</city>
    <temperature>18.3</temperature>
    <condition>Lluvia ligera</condition>
  </WeatherReport>
  <WeatherReport>
    <city>Medellín</city>
    <temperature>24.1</temperature>
    <condition>Soleado</condition>
  </WeatherReport>
  <WeatherReport>
    <city>Cali</city>
    <temperature>29.0</temperature>
    <condition>Nublado</condition>
  </WeatherReport>
  <WeatherReport>
    <city>Cartagena</city>
    <temperature>31.7</temperature>
    <condition>Caluroso y húmedo</condition>
  </WeatherReport>
  <WeatherReport>
    <city>Barranquilla</city>
    <temperature>30.4</temperature>
    <condition>Soleado con brisa</condition>
  </WeatherReport>
  <WeatherReport>
    <city>Pasto</city>
    <temperature>14.2</temperature>
    <condition>Frío y lluvioso</condition>
  </WeatherReport>
  <WeatherReport>
    <city>Manizales</city>
    <temperature>17.0</temperature>
    <condition>Niebla parcial</condition>
  </WeatherReport>
  <WeatherReport>
    <city>Armenia</city>
    <temperature>22.5</temperature>
    <condition>Clima templado</condition>
  </WeatherReport>
</WeatherReports>`

type WeatherReport struct {
	City        string  `json:"city" xml:"city"`
	Temperature float32 `json:"temperature" xml:"temperature"`
	Condition   string  `json:"condition" xml:"condition"`
}

func marshallWeatherRequest(report WeatherReport) ([]byte, error) {
	bytes, err := xml.Marshal(report)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func encodeWeatherRequest(report WeatherReport) ([]byte, error) {

	buff := bytes.NewBuffer([]byte{})
	encoder := xml.NewEncoder(buff)

	err := encoder.Encode(report)
	if err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

func postWeatherReport(report WeatherReport) error {
	// Using marshal/unmarshall
	// bytes, err := marshallWeatherRequest(report)
	// if err != nil {
	// 	return err
	// }

	// Using encode/decode
	bytes, err := encodeWeatherRequest(report)
	if err != nil {
		return err
	}

	log.Default().Printf("Sending bytes to SOAP: %s\n", string(bytes))

	time.Sleep(3 * time.Second)
	shouldFail := rand.Float32() < 0.5

	if shouldFail {
		return errors.New("Failed to post weather report!")
	}

	return nil
}

func decodeWeatherResponse(response string) ([]WeatherReport, error) {
	reader := strings.NewReader(response)
	decoder := json.NewDecoder(reader)

	var reports []WeatherReport
	err := decoder.Decode(&reports)
	if err != nil {
		return []WeatherReport{}, errors.New("Failed to decode JSON response!")
	}

	return reports, nil
}

func unmarshallWeatherResponse(response string) ([]WeatherReport, error) {
	var reports []WeatherReport
	err := json.Unmarshal([]byte(response), &reports)
	if err != nil {
		return []WeatherReport{}, err
	}

	return reports, nil
}

func getWeatherReport(city string) (WeatherReport, error) {
	time.Sleep(2 * time.Second)
	shouldFail := rand.Float32() < 0.5

	if shouldFail {
		return WeatherReport{}, errors.New("Failed to get the Weather Report!")
	}

	// Using encode/decode
	// reports, err := decodeWeatherResponse(JSON_RESPONSE)
	// if err != nil {
	// 	return WeatherReport{}, err
	// }

	// Using marshal/Unmarshall
	reports, err := unmarshallWeatherResponse(JSON_RESPONSE)
	if err != nil {
		return WeatherReport{}, err
	}

	var report *WeatherReport = nil

	for _, val := range reports {
		if val.City == city {
			report = &val
			break
		}
	}

	if report == nil {
		return WeatherReport{}, errors.New("City not found!")
	}

	return *report, nil
}

func main() {
	group := sync.WaitGroup{}

	group.Add(1)
	go func() {
		defer group.Done()

		city := "Armenia"
		fmt.Printf("Getting all weather report from %s...\n", city)
		report, err := getWeatherReport(city)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Done getting report!\nCity:\t%s\nTemperature:\t%f\nCondition:\t%s\n", report.City, report.Temperature, report.Condition)
	}()

	group.Add(1)
	go func() {
		defer group.Done()

		city := "Guatemala"
		fmt.Printf("Updating weather report from %s...\n", city)
		err := postWeatherReport(WeatherReport{City: city, Temperature: 23.5, Condition: "Windy"})
		if err != nil {
			panic(err)
		}
		fmt.Println("Done updating report!")
	}()

	group.Wait()
}
