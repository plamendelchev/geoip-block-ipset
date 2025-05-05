package geoip

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/plamendelchev/geoip-block-ipset/internal/config"
	"github.com/plamendelchev/geoip-block-ipset/internal/ipset"
	"github.com/plamendelchev/geoip-block-ipset/internal/ripe"
	"github.com/plamendelchev/geoip-block-ipset/internal/utils"
)

func setup(debug bool) error {
	// Ensure superuser
	isRoot, err := utils.IsRoot()
	if err != nil {
		return err
	}
	if !isRoot {
		return fmt.Errorf("you need superuser privileges to run this program")
	}

	// Set Up logger
	log.SetOutput(os.Stdout)
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}

// Create sets up the geoip whitelist
func Create(configFile string, debug bool) error {
	// Initial setup
	if err := setup(debug); err != nil {
		return err
	}

	// Read config file
	log.WithFields(log.Fields{"file": configFile}).Info("Reading configuration file")
	config, err := config.Read(configFile)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"file": configFile}).Info("Successfully read configuration file")

	// Obtain IP ranges from RIPE
	log.WithFields(log.Fields{"allowed_countries": config.AllowedCountries}).Info("Getting IP Ranges from RIPE")
	ranges, err := ripe.Ranges(config.AllowedCountries)
	if err != nil {
		return err
	}
	// Log the number of IP ranges per country
	ranges_per_cc := utils.TotalRangesPerCountry(*ranges)
	log.WithFields(log.Fields(*ranges_per_cc)).Info("Successfully got IP Ranges from RIPE")

	// Convert country names from cc to geoip_allow_cc
	ipSetRanges := make(ripe.AllowedCountries)
	for k, v := range *ranges {
		ipSetRanges[utils.ToIpSetName(k)] = v
	}

	// Create set for allowed_ranges
	if len(config.AllowedRanges) > 0 {
		log.WithFields(log.Fields{"allowed_ranges": len(config.AllowedRanges)}).Info("Found allowed_ranges in config")
		ipSetRanges["geoip_allowed_ranges"] = config.AllowedRanges
	}

	// Create slice with the names of the sets
	set_names := maps.Keys(ipSetRanges)

	// Create and populate IPSet sets
	log.WithFields(log.Fields{"sets": set_names}).Info("Creating IPSet sets")
	if err := ipset.Create(ipSetRanges); err != nil {
		return err
	}
	log.WithFields(log.Fields{"sets": set_names}).Info("Successfully created IPSet sets")

	log.Info("Done")
	return nil
}

// Delete removes the geoip whitelist
func Delete(configFile string, debug bool) error {
	// Initial setup
	if err := setup(debug); err != nil {
		return err
	}

	// Read config file
	log.WithFields(log.Fields{"file": configFile}).Info("Reading configuration file")
	config, err := config.Read(configFile)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"file": configFile}).Info("Successfully read configuration file")

	// Convert country names from cc to geoip_block_cc
	var rules []string
	for _, country := range config.AllowedCountries {
		rules = append(rules, utils.ToIpSetName(country))
	}

	// Remove geoip_allowed_ranges set if configured
	if len(config.AllowedRanges) > 0 {
		rules = append(rules, "geoip_allowed_ranges")
	}

	// Remove IPSet sets
	log.WithFields(log.Fields{"sets": rules}).Info("Deleting IPSet sets")
	if err := ipset.Remove(rules); err != nil {
		return err
	}
	log.WithFields(log.Fields{"sets": rules}).Info("Successfully deleted IPSet sets")

	return nil
}
