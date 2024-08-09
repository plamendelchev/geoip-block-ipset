package iptables

import (
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	log "github.com/sirupsen/logrus"
)

// Block set in iptables
func Create(chains []string) error {
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

// Remove IPtables rules
func Remove(chains []string) error {
	ipt, err := iptables.New()
	if err != nil {
		return fmt.Errorf("IPTables Error: %q", err)
	}

	for _, chain := range chains {
		t := "filter"
		c := "INPUT"
		rs := []string{"-m", "set", "--match-set", chain, "src", "-j", "ACCEPT"}

		log.WithFields(log.Fields{"table": t, "chain": c, "rulespec": rs}).Debug("Deleting rule")

		ipt.Delete(t, c, rs...)
	}

	return nil
}
