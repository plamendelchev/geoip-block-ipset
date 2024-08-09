package ipset

import (
	"fmt"

	"github.com/plamendelchev/geoip-block-ipset/internal/ripe"
	log "github.com/sirupsen/logrus"
)

// Create ipset set and add ranges
func Create(countries ripe.AllowedCountries) error {
	for name, ranges := range countries {
		ips_type := "hash:net"
		ips_params := Params{}

		log.WithFields(log.Fields{"name": name, "type": ips_type, "params": fmt.Sprintf("%+v", ips_params)}).
			Debug("Creating set")

		set, err := New(name, ips_type, &ips_params)
		if err != nil {
			return fmt.Errorf("IPSet Error: %q", err)
		}

		log.WithFields(log.Fields{"name": name, "num_ranges": len(ranges)}).
			Debug("Adding ranges to set")

		if err := set.Refresh(ranges); err != nil {
			return fmt.Errorf("IPSet Error: %q", err)
		}
	}

	return nil
}

// Remove IPSet sets
func Remove(sets []string) error {
	for _, name := range sets {
		ips_type := "hash:net"
		ips_params := Params{}

		set, err := New(name, ips_type, &ips_params)
		if err != nil {
			return fmt.Errorf("IPSet Error: %q", err)
		}

		if err := set.Destroy(); err != nil {
			return fmt.Errorf("IPSet Error: %q", err)
		}
	}
	return nil
}
