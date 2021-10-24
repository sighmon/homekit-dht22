# HomeKit DHT22

An [Apple HomeKit](https://developer.apple.com/homekit/) accessory for the [DHT22](https://www.adafruit.com/product/385) temperature and humidity sensor running on a Raspberry Pi.

![The accessory added to iOS](_images/homekit-dht22.jpg)

## Dependencies

* [**Go**](http://golang.org/doc/install) - this accessory is written in Go
* [**HomeControl**](https://github.com/brutella/hc) - to expose climate readings from the Enviro+ as an Apple HomeKit accessory
* [**climate-sensor**](https://github.com/ACMILabs/climate-sensor) - to read the DHT22 sensor and export it for scraping by [Prometheus](https://prometheus.io)

## Installation

Install this on a Raspberry Pi, or test it on macOS.

### Setup

1. Install [Go](http://golang.org/doc/install) >= 1.14 ([useful Gist](https://gist.github.com/pcgeek86/0206d688e6760fe4504ba405024e887c) for Raspberry Pi)
1. Clone this project: `git clone https://github.com/sighmon/homekit-dht22` and then `cd homekit-dht22`
1. Install the Go dependencies: `go get`
1. Install and run the Climate Sensor exporter [climate-sensor](https://github.com/ACMILabs/climate-sensor)

### Build

1. To build this accessory: `go build homekit-dht22.go`
1. To cross-compile for Raspberry Pi on macOS: `env GOOS=linux GOARCH=arm GOARM=7 go build homekit-dht22.go`

### Run

1. Execute the executable: `./homekit-dht22`
1. Or run with the command: `go run homekit-dht22.go`

### Start automatically at boot

1. sudo cp homekit-dht22.service /lib/systemd/system/homekit-dht22.service
2. sudo systemctl daemon-reload
3. sudo systemctl enable homekit-dht22.service
4. sudo systemctl start homekit-dht22.service

### Optional flags

The flag defaults can be overridden by handing them in at runtime:

* `-host=http://0.0.0.0` The host of your DHT22 sensor
* `-port=1006` The port of your DHT22 sensor
* `-sleep=5s` The [time](https://golang.org/pkg/time/#ParseDuration) between updating the accessory with sensor readings (`5s` equals five seconds)
* `-dev` This turns on development mode to return a random temperature reading without needing to have a DHT22

e.g. to override the port run: `go run homekit-dht22.go -port=8000` or `./homekit-dht22 -port=8000`

## Reset this accessory

If you uninstall this accessory from your Apple Home, you'll also need to delete the stored data for it to be able to be re-added.

### macOS

1. Delete the data in the folder created: `homekit-dht22/DHT22/`
1. Restart the executable

### Raspberry Pi

1. Remove the persistent data: `rm -rf /var/lib/homekit-dht22/data`
1. Restart the Raspberry Pi

## Thanks

This project uses the amazing work of [Matthias](https://github.com/brutella). Please consider donating if you found this useful.
