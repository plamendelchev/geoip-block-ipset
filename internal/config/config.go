package config

import (
	"fmt"

	"github.com/biter777/countries"
	"gopkg.in/ini.v1"
)

type Config struct {
	AllowedCountries []string `ini:"allowed_countries"`
}

// Read config file
func Read(path string) (*Config, error) {
	inidata, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Configuration Error: %q", err)
	}

	var c Config

	err = inidata.MapTo(&c)
	if err != nil {
		return nil, fmt.Errorf("Configuration Error: %q", err)
	}
	if len(c.AllowedCountries) == 0 {
		return nil, fmt.Errorf("Configuration Error: %q is empty", "allowed_countries")
	}

	// Ensure that all country codes are valid
	for _, country := range c.AllowedCountries {
		cc := countries.ByName(country)
		if !countries.CountryCode.IsValid(cc) {
			return nil, fmt.Errorf("Configuration Error: %q is not a valid Country Code", c)
		}
	}

	return &c, nil
}
