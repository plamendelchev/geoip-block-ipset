package ripe

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	RIPE_URL     = "https://stat.ripe.net/data/country-resource-list/data.json?v4_format=prefix"
	USER_AGENT   = "geoip-block-ipset"
	HTTP_TIMEOUT = 20 * time.Second
)

type AllowedCountries map[string][]string

// Download IP ranges from RIPE
func Ranges(configCountries []string) (*AllowedCountries, error) {
	client := &http.Client{
		Timeout: HTTP_TIMEOUT,
	}

	ranges := make(AllowedCountries)

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

		ranges[cc] = append(ranges[cc], r...)
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("Request Error: %q", "RIPE returned 0 IP ranges")
	}

	return &ranges, nil
}
