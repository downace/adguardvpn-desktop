package adguard

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"github.com/downace/adguardvpn-desktop/internal/common"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os/exec"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

type VpnMode = string

const (
	VpnModeTun   = "tun"
	VpnModeSocks = "socks"
)

type ExclusionMode string

const (
	ExclusionModeGeneral   ExclusionMode = "general"
	ExclusionModeSelective ExclusionMode = "selective"
)

type SubscriptionType = string

const (
	SubscriptionFree    = "free"
	SubscriptionPremium = "premium"
)

type Subscription struct {
	Type       SubscriptionType `json:"type"`
	ValidUntil time.Time        `json:"validUntil" ts_type:"string"`
	MaxDevices uint8            `json:"maxDevices"`
}

type Account struct {
	Username     string       `json:"username"`
	Subscription Subscription `json:"subscription"`
}

type Status struct {
	Connecting bool      `json:"connecting"`
	Connected  bool      `json:"connected"`
	Location   *Location `json:"location"`
	VpnMode    VpnMode   `json:"mode"`
}

type Location struct {
	Iso     string `json:"iso"`
	Country string `json:"country"`
	City    string `json:"city"`
	Ping    int    `json:"ping"`
	// Estimate string // don't know what is it
}

type Cli struct {
	CliBin         string
	OnStatusChange func(status *Status)

	locations []Location
	status    *Status
}

func (a *Cli) exec(args ...string) (string, error) {
	cmd := exec.Command(a.CliBin, args...)

	outputBytes, err := cmd.CombinedOutput()

	output := string(outputBytes)

	if err != nil && output != "" {
		err = fmt.Errorf("%s%s", output, err)
	}

	return stripansi.Strip(output), err
}

func (a *Cli) Version() (string, error) {
	return a.exec("--version")
}

var statusRe = regexp.MustCompile("Connected to (.+) in (.+) mode")

func (a *Cli) GetStatus() (*Status, error) {
	if a.status == nil {
		err := a.RefreshStatus()
		if err != nil {
			return nil, err
		}
	}
	return a.status, nil
}

func (a *Cli) RefreshStatus() (err error) {
	a.status, err = a.fetchStatus()
	if err == nil {
		a.OnStatusChange(a.status)
	}
	return
}

func (a *Cli) fetchStatus() (*Status, error) {
	statusOutput, err := a.exec("status")
	if err != nil {
		return nil, err
	}

	matches := statusRe.FindStringSubmatch(statusOutput)

	if matches == nil {
		return &Status{
			Connecting: false,
			Connected:  false,
			Location:   nil,
		}, nil
	} else {
		location, _ := a.getLocationByCityName(matches[1])
		if location == nil {
			location = &Location{
				City: cases.Title(language.English).String(matches[1]),
			}
		}
		return &Status{
			Connecting: false,
			Connected:  true,
			Location:   location,
			VpnMode:    matches[2],
		}, nil
	}
}

func (a *Cli) getLocationByCityName(city string) (*Location, error) {
	locations, err := a.GetLocations()
	if err != nil {
		return nil, err
	}
	for _, location := range locations {
		if strings.ToLower(location.City) == strings.ToLower(city) {
			return &location, nil
		}
	}

	return nil, nil
}

func (a *Cli) GetLocations() ([]Location, error) {
	if a.locations == nil {
		newLocations, err := a.RefreshLocations()
		if err != nil {
			return nil, err
		}
		a.locations = newLocations
	}

	return a.locations, nil
}

var (
	licenseUsernameRe   = regexp.MustCompile("Logged in as (.+)")
	licenseSubTypeRe    = regexp.MustCompile("You are using the (FREE|PREMIUM) version")
	maxDevicesOnRe      = regexp.MustCompile("Up to (.+) devices simultaneously")
	licenseValidUntilRe = regexp.MustCompile("Your subscription will be renewed on (.+)")
)

func (a *Cli) Connect(location string) error {
	status, err := a.GetStatus()
	if err != nil {
		return err
	}
	status.Connecting = true
	a.OnStatusChange(status)
	args := []string{"connect", "--yes"}
	if location != "" {
		args = append(args, "--location", location)
	}
	_, err = a.exec(args...)

	status.Connecting = false
	if err == nil {
		_ = a.RefreshStatus()
	}

	return err
}

func (a *Cli) Disconnect() error {
	_, err := a.exec("disconnect")

	if err == nil {
		_ = a.RefreshStatus()
	}

	return err
}

func (a *Cli) ToggleConnection() error {
	status, err := a.GetStatus()
	if err != nil {
		return err
	}
	if status.Connected {
		return a.Disconnect()
	} else {
		return a.Connect("")
	}
}

const logInMessage = "Please log in"

func (a *Cli) Account() (*Account, error) {
	output, err := a.exec("license")

	if err != nil {
		if strings.Contains(err.Error(), logInMessage) {
			return nil, nil
		}
		return nil, err
	}

	if strings.Contains(output, logInMessage) {
		return nil, nil
	}

	account := Account{}

	for line := range strings.Lines(output) {
		usernameMatches := licenseUsernameRe.FindStringSubmatch(line)
		if usernameMatches != nil {
			account.Username = strings.TrimSpace(usernameMatches[1])
		}
		subscriptionTypeMatches := licenseSubTypeRe.FindStringSubmatch(line)
		if subscriptionTypeMatches != nil {
			account.Subscription.Type = strings.TrimSpace(subscriptionTypeMatches[1])
		}
		maxDevicesMatches := maxDevicesOnRe.FindStringSubmatch(line)
		if maxDevicesMatches != nil {
			maxDevices, err := strconv.ParseUint(strings.TrimSpace(maxDevicesMatches[1]), 10, 8)
			if err == nil {
				account.Subscription.MaxDevices = uint8(maxDevices)
			}
		}
		validUntilMatches := licenseValidUntilRe.FindStringSubmatch(line)
		if validUntilMatches != nil {
			date, err := time.Parse(time.DateOnly, strings.TrimSpace(validUntilMatches[1]))
			if err == nil {
				account.Subscription.ValidUntil = date
			}
		}
	}

	return &account, nil
}

func (a *Cli) RefreshLocations() ([]Location, error) {
	output, err := a.exec("list-locations")

	if err != nil {
		return nil, err
	}

	if strings.Contains(output, "Please log in") {
		return make([]Location, 0), nil
	}

	return common.ParseTable(output, func(row map[string]string) Location {
		location := Location{}
		for key, value := range row {
			switch key {
			case "ISO":
				location.Iso = value
			case "COUNTRY":
				location.Country = value
			case "CITY":
				location.City = value
			case "PING":
				ping, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					ping = -1
				}
				location.Ping = int(ping)
			}
		}
		return location
	}), nil
}

var exclusionModeRe = regexp.MustCompile(`exclusion mode is (.+)`)

func (a *Cli) GetExclusionMode() (ExclusionMode, error) {
	output, err := a.exec("site-exclusions", "mode")

	if err != nil {
		return "", err
	}

	matches := exclusionModeRe.FindStringSubmatch(output)

	if matches == nil {
		return "", fmt.Errorf("unable to parse exclusion mode: %s", output)
	}

	return ExclusionMode(strings.ToLower(matches[1])), nil
}

func (a *Cli) SetExclusionMode(newMode ExclusionMode) error {
	_, err := a.exec("site-exclusions", "mode", string(newMode))

	return err
}

func (a *Cli) ExclusionsShow() ([]string, error) {
	output, err := a.exec("site-exclusions", "show")

	if err != nil {
		return nil, err
	}

	result := make([]string, 0)

	for line := range strings.Lines(output) {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Exclusions for") || line == "" {
			continue
		}

		result = append(result, line)
	}

	return result, nil
}

func (a *Cli) ExclusionsAdd(exclusions []string) error {
	args := slices.Concat([]string{"site-exclusions", "add"}, exclusions)
	_, err := a.exec(args...)

	return err
}

func (a *Cli) ExclusionsRemove(exclusion string) error {
	_, err := a.exec("site-exclusions", "remove", exclusion)

	return err
}
