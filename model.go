package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type T struct {
	Devid  string      `json:"devid"`
	Time   json.Number `json:"time"`
	Pout   string      `json:"pout"`
	Powset string      `json:"powset"`
	Data   Data        `json:"data"`
}

type Data struct {
	A   V `json:"A"`
	W   V `json:"W"`
	TWh V `json:"TWh"`

	V1 V `json:"V1"`
	V2 V `json:"V2"`
	V3 V `json:"V3"`

	A1 V `json:"A1"`
	A2 V `json:"A2"`
	A3 V `json:"A3"`

	W1 V `json:"W1"`
	W2 V `json:"W2"`
	W3 V `json:"W3"`

	Wh1 V `json:"Wh1"`
	Wh2 V `json:"Wh2"`
	Wh3 V `json:"Wh3"`

	PF1 V `json:"PF1"`
	PF2 V `json:"PF2"`
	PF3 V `json:"PF3"`

	Fr1 V `json:"Fr1"`
	Fr2 V `json:"Fr2"`
	Fr3 V `json:"Fr3"`

	Br0 BR `json:"br0"`
	Br1 BR `json:"br1"`
	Br2 BR `json:"br2"`
	Br3 BR `json:"br3"`

	T V `json:"T"`
}

type V struct {
	Name  string `json:"name"`
	Unit  string `json:"unit"`
	Value string `json:"value"`
}

type BR struct {
	Name string `json:"name"`
}

func (v V) MustGetFloat64Value() float64 {
	value := strings.ReplaceAll(v.Value, " ", "")

	f64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Panicf("value %s is not float64", v.Value)
	}

	return f64
}

type DeviceAPIStatus float64

const (
	DeviceAPIStatusOffline         DeviceAPIStatus = 0
	DeviceAPIStatusOK              DeviceAPIStatus = 1
	DeviceAPIStatusTooManuRequests DeviceAPIStatus = 2
)
