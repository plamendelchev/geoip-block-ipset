package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/biter777/countries"
	"github.com/coreos/go-iptables/iptables"
	"github.com/janeczku/go-ipset/ipset"
	"golang.org/x/exp/maps"
	"gopkg.in/ini.v1"
)

const (
	CONFIG_FILE = "/etc/geoip-block.conf"
	RIPE_URL    = "https://stat.ripe.net/data/country-resource-list/data.json?v4_format=prefix"
	USER_AGENT  = "geoip-block-ipset"
)

type config struct {
	AllowedCountries []string `ini:"allowed_countries"`
}

type allowedCountries map[string][]string

func Setup(configFile string, debug bool) error {
	// Ensure superuser
	isRoot, err := isRoot()
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

	// Read config file
	log.WithFields(log.Fields{"file": configFile}).Info("Reading configuration file")
	config, err := readConfig(configFile)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"file": configFile}).Info("Successfully read configuration file")

	// Obtain IP ranges from RIPE
	log.WithFields(log.Fields{"allowed_countries": config.AllowedCountries}).Info("Getting IP Ranges from RIPE")
	ranges, err := getIpRanges(config.AllowedCountries)
	if err != nil {
		return err
	}
	// Log the number of IP ranges per country
	fields := make(map[string]interface{})
	for country, ranges := range *ranges {
		fields[country] = len(ranges)
	}
	log.WithFields(log.Fields(fields)).Info("Successfully got IP Ranges from RIPE")

	// Create and populate IPSet sets
	log.WithFields(log.Fields{"sets": maps.Keys(*ranges)}).Info("Creating IPSet sets")
	err = createIpSets(*ranges)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"sets": maps.Keys(*ranges)}).Info("Successfully created IPSet sets")

	// Create IPTables rules
	rules := maps.Keys(*ranges)
	log.WithFields(log.Fields{"rules": rules}).Info("Creating IPTables rules")
	err = createIpTablesRules(rules)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{"rules": rules}).Info("Successfully created IPTables rules")

	log.Info("Done")
	return nil
}

// Determine if user is superuser
func isRoot() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, fmt.Errorf("Failed to determine user: %q", err)
	}
	return currentUser.Username == "root", nil
}

// Read config file
func readConfig(path string) (*config, error) {
	inidata, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("Configuration Error: %q", err)
	}

	var config config
	err = inidata.MapTo(&config)
	if err != nil {
		return nil, fmt.Errorf("Configuration Error: %q", err)
	}
	if len(config.AllowedCountries) == 0 {
		return nil, fmt.Errorf("Configuration Error: %q is empty", "allowed_countries")
	}

	// Ensure that all country codes are valid
	for _, c := range config.AllowedCountries {
		cc := countries.ByName(c)
		if !countries.CountryCode.IsValid(cc) {
			return nil, fmt.Errorf("Configuration Error: %q is not a valid Country Code", c)
		}
	}

	return &config, nil
}

// Download IP ranges from RIPE
func getIpRanges(configCountries []string) (*allowedCountries, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	ranges := make(allowedCountries)

	for _, cc := range configCountries {
		url := fmt.Sprintf("%s&resource=%s", RIPE_URL, cc)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("Request Error: %q", err)
		}

		req.Header.Set("User-Agent", USER_AGENT)

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Request Error: %q", err)
		}

		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("Request Error: %q", err)
		}
		if resp.StatusCode > 299 {
			return nil, fmt.Errorf("Request Error: %q %q", resp.StatusCode, string(body))
		}

		// Deserialize JSON response
		var res map[string]json.RawMessage
		var r []string
		err = json.Unmarshal(body, &res)
		err = json.Unmarshal(res["data"], &res)
		err = json.Unmarshal(res["resources"], &res)
		err = json.Unmarshal(res["ipv4"], &r)
		if err != nil {
			return nil, fmt.Errorf("Request Error: %q", err)
		}

		name := fmt.Sprintf("geoip_allow_%s", strings.ToLower(cc))
		ranges[name] = append(ranges[name], r...)
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("Request Error: %q", "RIPE returned 0 IP ranges")
	}

	return &ranges, nil
}

// Create ipset set and add ranges
func createIpSets(countries allowedCountries) error {
	for name, ranges := range countries {
		ips_type := "hash:net"
		ips_params := ipset.Params{}

		log.WithFields(log.Fields{"name": name, "type": ips_type, "params": fmt.Sprintf("%+v", ips_params)}).Debug("Creating set")

		set, err := ipset.New(name, ips_type, ips_params)
		if err != nil {
			return fmt.Errorf("IPSet Error: %q", err)
		}

		log.WithFields(log.Fields{"name": name, "num_ranges": len(ranges)}).Debug("Adding ranges to set")

		err = set.Refresh(ranges)
		if err != nil {
			return fmt.Errorf("IPSet Error: %q", err)
		}
	}

	return nil
}

// Block set in iptables
func createIpTablesRules(chains []string) error {
	ipt, err := iptables.New()
	if err != nil {
		return fmt.Errorf("IPTables Error: %q", err)
	}

	for _, chain := range chains {
		t := "filter"
		c := "INPUT"
		rs := []string{"-m", "set", "--match-set", chain, "src", "-j", "ACCEPT"}

		log.WithFields(log.Fields{"table": t, "chain": c, "rulespec": rs}).Debug("Creating rule")

		err = ipt.AppendUnique(t, c, rs...)
		if err != nil {
			return fmt.Errorf("IPTables Error: %q", err)
		}
	}

	return nil
}
