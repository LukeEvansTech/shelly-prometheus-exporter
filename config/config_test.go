package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestApplyCredentialDefaultsEnvFallback(t *testing.T) {
	t.Setenv("SHELLY_USERNAME", "")
	t.Setenv("SHELLY_PASSWORD", "test-envpass")

	cfg := &YamlConfig{Devices: []DeviceYamlConfig{
		{Host: "a"}, // empty -> env pass, admin user
		{Host: "b", Username: "u", Password: "p"}, // explicit -> preserved
		{Host: "c", Password: "explicit"},         // pass kept, user -> admin
		{Host: "d", Username: "custom"},           // explicit user, no pass -> env pass, user kept
	}}
	applyCredentialDefaults(cfg)

	if cfg.Devices[0].Password != "test-envpass" || cfg.Devices[0].Username != "admin" {
		t.Errorf("device a: got %+v, want password=envpass username=admin", cfg.Devices[0])
	}
	if cfg.Devices[1].Username != "u" || cfg.Devices[1].Password != "p" {
		t.Errorf("device b: explicit creds not preserved: %+v", cfg.Devices[1])
	}
	if cfg.Devices[2].Password != "explicit" || cfg.Devices[2].Username != "admin" {
		t.Errorf("device c: got %+v, want password=explicit username=admin", cfg.Devices[2])
	}
	if cfg.Devices[3].Username != "custom" || cfg.Devices[3].Password != "test-envpass" {
		t.Errorf("device d: got %+v, want username=custom password=envpass", cfg.Devices[3])
	}
}

func TestApplyCredentialDefaultsExplicitEnvUser(t *testing.T) {
	t.Setenv("SHELLY_USERNAME", "operator")
	t.Setenv("SHELLY_PASSWORD", "test-envpass")

	cfg := &YamlConfig{Devices: []DeviceYamlConfig{{Host: "a"}}}
	applyCredentialDefaults(cfg)

	if cfg.Devices[0].Username != "operator" || cfg.Devices[0].Password != "test-envpass" {
		t.Errorf("got %+v, want username=operator password=envpass", cfg.Devices[0])
	}
}

func TestApplyCredentialDefaultsNoEnvLeavesEmpty(t *testing.T) {
	t.Setenv("SHELLY_USERNAME", "")
	t.Setenv("SHELLY_PASSWORD", "")

	cfg := &YamlConfig{Devices: []DeviceYamlConfig{{Host: "a"}}}
	applyCredentialDefaults(cfg)

	if cfg.Devices[0].Username != "" || cfg.Devices[0].Password != "" {
		t.Errorf("expected empty creds with no env, got %+v", cfg.Devices[0])
	}
}

// writeConfig writes body to a temp file and returns its path.
func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return path
}

// deviceUpdateInterval is a count of seconds. It used to be typed
// time.Duration, which happened to work only because yaml.v2 decodes a bare
// integer into a Duration as a raw nanosecond count and the ticker then
// multiplied by time.Second to compensate. This pins the units so that
// compensating multiply cannot come back.
func TestNewConfigDeviceUpdateIntervalIsSeconds(t *testing.T) {
	cfg, err := NewConfig(writeConfig(t, "deviceUpdateInterval: 30\ndevices:\n- host: a\n"))
	if err != nil {
		t.Fatalf("NewConfig: %v", err)
	}
	if cfg.DeviceUpdateInterval != 30 {
		t.Errorf("DeviceUpdateInterval = %d, want 30", cfg.DeviceUpdateInterval)
	}
	if got := time.Duration(cfg.DeviceUpdateInterval) * time.Second; got != 30*time.Second {
		t.Errorf("poll interval = %v, want 30s", got)
	}
}

// The old time.Duration field also accepted a duration string, and yaml.v2
// decoded "30s" to 30000000000 -- which the ticker then multiplied by
// time.Second, overflowing int64 into a negative interval and panicking
// time.NewTicker at startup. An int rejects it at decode time instead.
func TestNewConfigDeviceUpdateIntervalRejectsDurationString(t *testing.T) {
	if _, err := NewConfig(writeConfig(t, "deviceUpdateInterval: 30s\ndevices:\n- host: a\n")); err == nil {
		t.Fatal("expected a decode error for a duration string, got nil")
	}
}
