package client

import "testing"

func TestPresetInitialization(t *testing.T) {
	t.Run("consensus-specs presets init", func(t *testing.T) {
		if MainnetPreset == nil {
			t.Error("MainnetConfig should be initialized")
		}
		if MinimalPreset == nil {
			t.Error("MinimalConfig should be initialized")
		}
	})
}
