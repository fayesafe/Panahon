# Panahon

Weather Station in Go.

## Install

#### Server
* `$ cd Panahan`
* Install dependencies with `go get`.
* Compile by `go install` or `go build -o main.o main.go`.

#### Front End Application
* `$ cd Panahan/main/app`
* Run `bower install` to install all used dependencies.
* Place `/app` in same directory as the executable.

#### Database

* Start InfluxDB CLI `$ influx`
* `> CREATE DATABASE weather_data`
* `> USE weather_data`
* Create continues queries which group sensor data into the tables `hours`, `days` and `weeks`
* `> CREATE CONTINUOUS QUERY cq_hour ON weather_data BEGIN SELECT mean(temperature) AS temperature, max(temperature) AS max_temperature, min(temperature) AS min_temperature, mean(humidity) as humidity, mean(pressure) as pressure, sum(rain)/count(rain) as rain, sum(sun)/count(sun) as sun INTO hours FROM sensors where time > '2015-12-01' GROUP BY time(1h) END`
* `> CREATE CONTINUOUS QUERY cq_hour ON weather_data BEGIN SELECT mean(temperature) AS temperature, max(max_temperature) AS max_temperature, min(min_temperature) AS min_temperature, mean(humidity) as humidity, mean(pressure) as pressure, sum(rain)/count(rain) as rain, sum(sun)/count(sun) as sun INTO days FROM hours where time > '2015-12-01' GROUP BY time(1d) END`
* `> CREATE CONTINUOUS QUERY cq_hour ON weather_data BEGIN SELECT mean(temperature) AS temperature, max(max_temperature) AS max_temperature, min(min_temperature) AS min_temperature, mean(humidity) as humidity, mean(pressure) as pressure, sum(rain)/count(rain) as rain, sum(sun)/count(sun) as sun INTO weeks FROM days where time > '2015-12-01' GROUP BY time(1w) END`

## Run

* make sure `/app` folder is in the same directory as binary
* make sure `config.toml` is in the same directory as binary
* set server and sensor options in `config.toml`
* `sudo $GO_PATH/bin/main`

## About

* Weather station using sensors for temperature, humidity,
 air pressure, air quality, brightness and rain
* Written entirely in Go
* Front end making use of modern technologies such as Angular & jQuery
* influxDB as database for working with time-series data

## Dependencies

* Mux
* Influxdb client
