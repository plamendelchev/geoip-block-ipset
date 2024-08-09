package geoip

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"

	"github.com/plamendelchev/geoip-block-ipset/internal/config"
	"github.com/plamendelchev/geoip-block-ipset/internal/ipset"
	"github.com/plamendelchev/geoip-block-ipset/internal/iptables"
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
		return fmt.Errorf("You need superuser privileges to run this program.")
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
	fields := make(map[string]interface{})
	for country, ranges := range *ranges {
		fields[country] = len(ranges)
	}
	log.WithFields(log.Fields(fields)).Info("Successfully got IP Ranges from RIPE")

	// Convert country names from cc to geoip_block_cc
	ipSetRanges := make(ripe.AllowedCountries)
	for k, v := range *ranges {
		ipSetRanges[fmt.Sprintf("geoip_allow_%s", k)] = v
	}
	rules := maps.Keys(ipSetRanges)

	// Create and populate IPSet sets
	log.WithFields(log.Fields{"sets": rules}).Info("Creating IPSet sets")
	if err := ipset.Create(ipSetRanges); err != nil {
		return err
	}
	log.WithFields(log.Fields{"sets": rules}).Info("Successfully created IPSet sets")

	// Create IPTables rules
	log.WithFields(log.Fields{"rules": rules}).Info("Creating IPTables rules")
	if err := iptables.Create(rules); err != nil {
		return err
	}
	log.WithFields(log.Fields{"rules": rules}).Info("Successfully created IPTables rules")

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

	// Remove IPTables rules
	if err := iptables.Remove(rules); err != nil {
		return err
	}

	// Remove IPSet sets
	if err := ipset.Remove(rules); err != nil {
		return err
	}

	return nil
}
