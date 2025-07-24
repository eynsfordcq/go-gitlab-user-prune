package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	APIBaseUrl     string
	APIToken       string
	InactivityDays int
	DryRun         bool
	Whitelist      map[string]struct{}
}

func Load() (*Config, error) {
	apiURL := os.Getenv("API_BASE_URL")
	if apiURL == "" {
		return nil, fmt.Errorf("environment variable API_BASE_URL not set")
	}

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("environment variable API_TOKEN not set")
	}

	inactivityDays := 90
	inactivityDaysStr := os.Getenv("INACTIVITY_DAYS")
	if inactivityDaysStr != "" {
		var err error
		inactivityDays, err = strconv.Atoi(inactivityDaysStr)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid value for INACTIVITY_DAYS: %w", err,
			)
		}
	}

	whitelist := make(map[string]struct{})
	whitelistStr := os.Getenv("WHITELIST")
	if whitelistStr != "" {
		users := strings.SplitSeq(whitelistStr, ",")
		for user := range users {
			trimmedUser := strings.TrimSpace(user)
			if trimmedUser != "" {
				whitelist[trimmedUser] = struct{}{}
			}
		}
	}

	dryRun := true
	dryRunStr := os.Getenv("DRY_RUN")
	if strings.ToLower(dryRunStr) == "false" {
		dryRun = false
	}

	return &Config{
		APIBaseUrl:     strings.TrimRight(apiURL, "/"),
		APIToken:       apiToken,
		Whitelist:      whitelist,
		DryRun:         dryRun,
		InactivityDays: inactivityDays,
	}, nil
}
