package utils

import (
	"fmt"
	"os/user"
	"strings"

	"github.com/plamendelchev/geoip-block-ipset/internal/ripe"
)

// Determine if user is superuser
func IsRoot() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, fmt.Errorf("Failed to determine user: %q", err)
	}
	return currentUser.Username == "root", nil
}

// Convert cc to geoip_allow_cc
func ToIpSetName(name string) string {
	return fmt.Sprintf("geoip_allow_%s", strings.ToLower(name))
}

//
func TotalRangesPerCountry(a ripe.AllowedCountries) *map[string]interface{} {
	r := make(map[string]interface{})
	for country, ranges := range a {
		r[country] = len(ranges)
	}
	return &r
}
