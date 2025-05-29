package client

import _ "embed"
import "fmt"
import "log"

import "gopkg.in/yaml.v3"

//go:embed presets/mainnet/phase0.yaml
var mainnetPhase0 []byte

//go:embed presets/mainnet/altair.yaml
var mainnetAltair []byte

//go:embed presets/mainnet/bellatrix.yaml
var mainnetBellatrix []byte

//go:embed presets/mainnet/capella.yaml
var mainnetCapella []byte

//go:embed presets/mainnet/deneb.yaml
var mainnetDeneb []byte

//go:embed presets/mainnet/electra.yaml
var mainnetElectra []byte

//go:embed presets/mainnet/fulu.yaml
var mainnetFulu []byte

//go:embed presets/minimal/phase0.yaml
var minimalPhase0 []byte

//go:embed presets/minimal/altair.yaml
var minimalAltair []byte

//go:embed presets/minimal/bellatrix.yaml
var minimalBellatrix []byte

//go:embed presets/minimal/capella.yaml
var minimalCapella []byte

//go:embed presets/minimal/deneb.yaml
var minimalDeneb []byte

//go:embed presets/minimal/electra.yaml
var minimalElectra []byte

//go:embed presets/minimal/fulu.yaml
var minimalFulu []byte

var (
	MainnetPreset map[string]interface{}
	MinimalPreset map[string]interface{}
)

func init() {
	err := yaml.Unmarshal(minimalPhase0, &MinimalPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(minimalAltair, &MinimalPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(minimalBellatrix, &MinimalPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(minimalCapella, &MinimalPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(minimalDeneb, &MinimalPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(minimalElectra, &MinimalPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(minimalFulu, &MinimalPreset)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(mainnetPhase0, &MainnetPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(mainnetAltair, &MainnetPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(mainnetBellatrix, &MainnetPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(mainnetCapella, &MainnetPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(mainnetDeneb, &MainnetPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(mainnetElectra, &MainnetPreset)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(mainnetFulu, &MainnetPreset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("minimal preset: %v\n", MinimalPreset)
	fmt.Printf("mainnet preset: %v\n", MainnetPreset)
}
