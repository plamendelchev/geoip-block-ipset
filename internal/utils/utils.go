package utils

import (
	"fmt"
	"os/user"
	"strings"
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
