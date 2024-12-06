package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics
var (
	CustomRegistry = prometheus.NewRegistry() // Create a custom registry

	metricPrefix = "smart_maic_"

	voltage     = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "voltage", Help: "Voltage per line (V)"}, []string{"line"})
	current     = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "current", Help: "Current per line (A)"}, []string{"line"})
	power       = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "power", Help: "Active Power per line (W)"}, []string{"line"})
	energy      = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "energy", Help: "Energy per line (Wh)"}, []string{"line"})
	powerFactor = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "power_factor", Help: "Power Factor per line"}, []string{"line"})
	frequency   = prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricPrefix + "frequency", Help: "Frequency per line (Hz)"}, []string{"line"})

	totalCurrent    = prometheus.NewGauge(prometheus.GaugeOpts{Name: metricPrefix + "total_current", Help: "Total Current (A)"})
	totalPower      = prometheus.NewGauge(prometheus.GaugeOpts{Name: metricPrefix + "total_power", Help: "Total Active Power (W)"})
	totalEnergy     = prometheus.NewGauge(prometheus.GaugeOpts{Name: metricPrefix + "total_energy", Help: "Total Energy (Wh)"})
	temperature     = prometheus.NewGauge(prometheus.GaugeOpts{Name: metricPrefix + "temperature", Help: "Device Temperature (°C)"})
	deviceAPIStatus = prometheus.NewGauge(prometheus.GaugeOpts{Name: metricPrefix + "device_api_status", Help: "Device API Status (0 = Offline, 1 = OK, 2 = Too Many Requests)"})
)

func init() {
	// Register Prometheus metrics
	CustomRegistry.MustRegister(voltage, current, power, energy, powerFactor, frequency, totalCurrent, totalPower, totalEnergy, temperature, deviceAPIStatus)
}

func SetMetrics(v T) {
	data := v.Data

	SetDeviceAPIStatus(DeviceAPIStatusOK)

	voltage.WithLabelValues("1").Set(data.V1.MustGetFloat64Value())
	voltage.WithLabelValues("2").Set(data.V2.MustGetFloat64Value())
	voltage.WithLabelValues("3").Set(data.V3.MustGetFloat64Value())

	current.WithLabelValues("1").Set(data.A1.MustGetFloat64Value())
	current.WithLabelValues("2").Set(data.A2.MustGetFloat64Value())
	current.WithLabelValues("3").Set(data.A3.MustGetFloat64Value())

	power.WithLabelValues("1").Set(data.W1.MustGetFloat64Value())
	power.WithLabelValues("2").Set(data.W2.MustGetFloat64Value())
	power.WithLabelValues("3").Set(data.W3.MustGetFloat64Value())

	energy.WithLabelValues("1").Set(data.Wh1.MustGetFloat64Value())
	energy.WithLabelValues("2").Set(data.Wh2.MustGetFloat64Value())
	energy.WithLabelValues("3").Set(data.Wh3.MustGetFloat64Value())

	powerFactor.WithLabelValues("1").Set(data.PF1.MustGetFloat64Value())
	powerFactor.WithLabelValues("2").Set(data.PF2.MustGetFloat64Value())
	powerFactor.WithLabelValues("3").Set(data.PF3.MustGetFloat64Value())

	frequency.WithLabelValues("1").Set(data.Fr1.MustGetFloat64Value())
	frequency.WithLabelValues("2").Set(data.Fr2.MustGetFloat64Value())
	frequency.WithLabelValues("3").Set(data.Fr3.MustGetFloat64Value())

	totalCurrent.Set(data.A.MustGetFloat64Value())
	totalPower.Set(data.W.MustGetFloat64Value())
	totalEnergy.Set(data.TWh.MustGetFloat64Value())
	temperature.Set(data.T.MustGetFloat64Value())
}

func SetDeviceAPIStatus(v DeviceAPIStatus) {
	deviceAPIStatus.Set(float64(v))
}
