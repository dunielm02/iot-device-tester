# Iot Device Tester

## Why
Sometimes the IoT devices should be placed on a site where the internet connection is very bad. So before connecting them to the main platform, you should test that those devices are well connected to the network and also send correct values. This Project is a lightweight solution to facilitate the network testing process and also offers a functional interface to check the correctness of the data sent by the device.

## Features

* Complete compiled to secure its use in every Operative System without troubles.

* It is capable of receive data throw HTTP and MQTT protocols.

* Counts with a simple UI that can show the received data in real time using the WebSockets protocol.

* An sqLite database to store all logs and access them every time that is needed.

## Getting Started

### Set environment variables

Before running you have to set this variables in your environment.

```env
#DATABASE

#Name of the database file
DB_NAME = "test.db"

#HTTP SERVER
HTTP_PORT = "8000"

#MQTT

#This variable sets if you want to use Mqtt
USE_MQTT = TRUE

#Mqtt Options
BROKER = "broker.dunielm02.com" 
MQTT_PORT = "8883"
CLIENT_ID = "iotDeviceTester"
USERNAME = ""
PASSWORD = ""
MQTT_USE_SSL = TRUE

#Topic of the first subscription
TOPIC = "/api/"
```

### Database structure

A simple database to store received messages containing the path to where they were sent and the body of the request in JSON format.

```go
type Log struct {
    ID       int `gorm:"primaryKey"`
    DeviceID string
    Path     string
    Data     datatypes.JSON
}
```

## Compiling

```shell
git clone https://github.com/dunielm02/iot-device-tester.git
go mod tidy
go build -o build/iot-device-tester
```

***Important: because [sqlite Driver](https://github.com/mattn/go-sqlite3) is a `CGO` enabled package, you are required to set the environment variable `CGO_ENABLED=1` and have a `gcc` compiler present within your path.***
