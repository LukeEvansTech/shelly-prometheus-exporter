package metrics

import (
	"github.com/LukeEvansTech/shelly-prometheus-exporter/config"
	"github.com/LukeEvansTech/shelly-prometheus-exporter/rpc"
	CoverGetStatus "github.com/LukeEvansTech/shelly-prometheus-exporter/rpc/Cover.GetStatus"
	ShellyGetConfig "github.com/LukeEvansTech/shelly-prometheus-exporter/rpc/Shelly.GetConfig"
	ShellyGetDeviceInfo "github.com/LukeEvansTech/shelly-prometheus-exporter/rpc/Shelly.GetDeviceInfo"
	ShellyGetStatus "github.com/LukeEvansTech/shelly-prometheus-exporter/rpc/Shelly.GetStatus"
	SwitchGetConfig "github.com/LukeEvansTech/shelly-prometheus-exporter/rpc/Switch.GetConfig"
	SwitchGetStatus "github.com/LukeEvansTech/shelly-prometheus-exporter/rpc/Switch.GetStatus"
	WiFiGetStatus "github.com/LukeEvansTech/shelly-prometheus-exporter/rpc/WiFi.GetStatus"
)

// Register initializes Prometheus metrics and starts periodic API fetching.
func Register(cfg *config.YamlConfig, cfgPath *string) {
	ShellyGetConfig.RegisterShellyGetConfigMetrics()
	ShellyGetStatus.RegisterShellyGetStatusMetrics()
	ShellyGetDeviceInfo.RegisterShellyGetDeviceInfoMetrics()
	CoverGetStatus.RegisterCoverGetStatusMetrics()
	SwitchGetStatus.RegisterSwitchGetStatusMetrics()
	SwitchGetConfig.RegisterSwitchGetConfigMetrics()
	WiFiGetStatus.RegisterWiFiGetStatusMetrics()

	dm := rpc.NewDeviceManager()

	for _, device := range cfg.Devices {
		dm.RegisterDevice(&rpc.DeviceConfig{
			Host:     device.Host,
			Username: device.Username,
			Password: device.Password,
		}, cfg.DeviceUpdateInterval)
	}
}
