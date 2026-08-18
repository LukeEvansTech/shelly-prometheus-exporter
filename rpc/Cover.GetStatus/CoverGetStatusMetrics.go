package CoverGetStatus

import (
	"fmt"

	"github.com/LukeEvansTech/shelly-prometheus-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type CoverGetStatusMetrics struct {
	State       *prometheus.GaugeVec
	APower      *prometheus.GaugeVec
	Voltage     *prometheus.GaugeVec
	Current     *prometheus.GaugeVec
	Pf          *prometheus.GaugeVec
	Freq        *prometheus.GaugeVec
	Energy      *prometheus.GaugeVec
	Temperature *prometheus.GaugeVec
	PosControl  *prometheus.GaugeVec
	Position    *prometheus.GaugeVec
}

var metrics *CoverGetStatusMetrics

func RegisterCoverGetStatusMetrics() {
	metrics = &CoverGetStatusMetrics{
		State: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "state",
			Help:      "Describes the current position aka state the cover is in. (1 = open, 0 = closed, 2 = in movement, 3 = stopped)",
		}, []string{"device_mac", "cover_id"}),
		APower: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "power",
			Help:      "Active power of the cover in Watts",
		}, []string{"device_mac", "cover_id"}),
		Voltage: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "voltage",
			Help:      "Present power in Volts",
		}, []string{"device_mac", "cover_id"}),
		Current: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "current",
			Help:      "Current draw by the cover in amps",
		}, []string{"device_mac", "cover_id"}),
		Pf: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "powerfactor",
			Help:      "Power factor of the cover",
		}, []string{"device_mac", "cover_id"}),
		Freq: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "frequency",
			Help:      "Current input frequency of the power source in Hz.",
		}, []string{"device_mac", "cover_id"}),
		Energy: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "energy",
			Help:      "Total consumption of the cover in Wh",
		}, []string{"device_mac", "cover_id"}),
		Temperature: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "temperature",
			Help:      "Temperature of the shelly device in C or F",
		}, []string{"device_mac", "cover_id", "temperature_unit"}),
		PosControl: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "pos_control",
			Help:      "Boolean indicating if position control is present",
		}, []string{"device_mac", "cover_id"}),
		Position: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "shelly",
			Subsystem: "cover",
			Name:      "position",
			Help:      "Current position of the cover",
		}, []string{"device_mac", "cover_id"}),
	}

	prometheus.MustRegister(
		metrics.State,
		metrics.APower,
		metrics.Voltage,
		metrics.Current,
		metrics.Pf,
		metrics.Freq,
		metrics.Energy,
		metrics.Temperature,
		metrics.PosControl,
		metrics.Position,
	)
}

func UpdateCoverGetStatusMetrics(apiClient *client.APIClient, coverID int, deviceMac string) error {
	var config client.CoverGetStatusResponse
	err := apiClient.FetchData(fmt.Sprintf("/rpc/Cover.GetStatus?id=%d", coverID), &config)
	if err != nil {
		return fmt.Errorf("error fetching config: %w", err)
	}

	metrics.UpdateMetrics(config, deviceMac)

	return nil
}

func (m *CoverGetStatusMetrics) UpdateMetrics(status client.CoverGetStatusResponse, deviceMac string) {
	coverID := fmt.Sprintf("%d", status.ID)

	switch state := status.State; state {
	case "open":
		m.State.WithLabelValues(deviceMac, coverID).Set(1)
	case "closed":
		m.State.WithLabelValues(deviceMac, coverID).Set(0)
	case "opening":
		m.State.WithLabelValues(deviceMac, coverID).Set(2)
	case "closing":
		m.State.WithLabelValues(deviceMac, coverID).Set(2)
	case "stopped":
		m.State.WithLabelValues(deviceMac, coverID).Set(3)
	case "calibrating":
		m.State.WithLabelValues(deviceMac, coverID).Set(2)
	default:
		m.State.WithLabelValues(deviceMac, coverID).Set(-1)
	}

	m.APower.WithLabelValues(deviceMac, coverID).Set(status.Apower)
	m.Voltage.WithLabelValues(deviceMac, coverID).Set(status.Voltage)
	m.Current.WithLabelValues(deviceMac, coverID).Set(status.Current)
	m.Pf.WithLabelValues(deviceMac, coverID).Set(status.Pf)
	m.Freq.WithLabelValues(deviceMac, coverID).Set(status.Freq)
	m.Energy.WithLabelValues(deviceMac, coverID).Set(status.Aenergy.Total)
	m.Temperature.WithLabelValues(deviceMac, coverID, "dC").Set(status.Temperature.TC)
	m.Temperature.WithLabelValues(deviceMac, coverID, "dF").Set(status.Temperature.TF)
	m.PosControl.WithLabelValues(deviceMac, coverID).Set(boolToFloat64(status.PosControl))
	m.Position.WithLabelValues(deviceMac, coverID).Set(float64(status.CurrentPos))
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
