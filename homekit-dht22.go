package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"

	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

var sensorHost string
var sensorPort int
var secondsBetweenReadings time.Duration
var developmentMode bool

func init() {
	flag.StringVar(&sensorHost, "host", "http://0.0.0.0", "sensor host, a string")
	flag.IntVar(&sensorPort, "port", 1006, "sensor port number, an int")
	flag.DurationVar(&secondsBetweenReadings, "sleep", 5*time.Second, "how many seconds between sensor readings, an int followed by the duration")
	flag.BoolVar(&developmentMode, "dev", false, "turn on development mode to return a random temperature reading, boolean")
	flag.Parse()

	if developmentMode == true {
		log.Println("Development mode on, ignoring sensor and returning random values...")
	}
}

func main() {
	info := accessory.Info{
		Name:             "DHT22",
		// Note: If running more than one sensor per home, update the serial number
		SerialNumber:     "DHT22",
		Manufacturer:     "Guangzhou Aosong Electronic Co., Ltd.",
		Model:            "DHT22",
		FirmwareRevision: "1.0.0",
	}

	acc := accessory.NewTemperatureSensor(
		info,
		0.0,   // Initial value
		-40.0, // Min sensor value
		80.0,  // Max sensor value
		0.1,   // Step value
	)

	humidity := service.NewHumiditySensor()
	acc.AddService(humidity.Service)
	acc.TempSensor.AddLinkedService(humidity.Service)

	config := hc.Config{
		// Change the default Apple Accessory Pin if you wish
		Pin: "00102003",
		// Port: "12345",
		// StoragePath: "./db",
	}

	t, err := hc.NewIPTransport(config, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	// Get the sensor readings every secondsBetweenReadings
	go func() {
		type Reading struct {
			Name  string
			Value float64
		}

		type Readings struct {
			Temperature Reading
			Humidity    Reading
		}

		readings := Readings{
			Temperature: Reading{
				Name:  "ambient_temperature",
				Value: 0,
			},
			Humidity: Reading{
				Name:  "ambient_humidity",
				Value: 0,
			},
		}
		values := reflect.ValueOf(readings)

		for {
			// Get readings from the Prometheus exporter
			resp, err := http.Get(fmt.Sprintf("%s:%d", sensorHost, sensorPort))
			if err == nil {
				defer resp.Body.Close()
				scanner := bufio.NewScanner(resp.Body)
				for scanner.Scan() {
					line := scanner.Text()
					// Parse the readings
					for i := 0; i < values.NumField(); i++ {
						fieldname := values.Field(i).Interface().(Reading).Name
						regexString := fmt.Sprintf("^%s", fieldname) + ` ([-+]?\d*\.\d+|\d+)`
						re := regexp.MustCompile(regexString)
						rs := re.FindStringSubmatch(line)
						if rs != nil {
							parsedValue, err := strconv.ParseFloat(rs[1], 64)
							if err == nil {
								if developmentMode {
									println(fmt.Sprintf("%s %f", fieldname, parsedValue))
								}
								switch fieldname {
								case "ambient_temperature":
									readings.Temperature.Value = parsedValue
								case "ambient_humidity":
									readings.Humidity.Value = parsedValue
								}
							}
						}
					}
				}
				scanner = nil
			} else {
				log.Println(err)
			}

			if developmentMode {
				// Return a random float between 15 and 30
				readings.Temperature.Value = 15 + rand.Float64()*(30-15)
			}

			// Set the sensor readings
			acc.TempSensor.CurrentTemperature.SetValue(readings.Temperature.Value)
			acc.TempSensor.CurrentTemperature.SetStepValue(0.1)
			humidity.CurrentRelativeHumidity.SetValue(readings.Humidity.Value)
			humidity.CurrentRelativeHumidity.SetStepValue(0.1)

			log.Println(fmt.Sprintf("Temperature: %f°C", readings.Temperature.Value))
			log.Println(fmt.Sprintf("Humidity: %f RH", readings.Humidity.Value))

			// Time between readings
			time.Sleep(secondsBetweenReadings)
		}
	}()

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
