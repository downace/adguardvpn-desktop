package adguard

import (
	"fmt"
	"github.com/acarl005/stripansi"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	Connected bool `json:"connected"`
}

type Cli struct {
	CliBin string
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

func (a *Cli) Status() (*Status, error) {
	statusOutput, err := a.exec("status")
	if err != nil {
		return nil, err
	}

	status := Status{
		Connected: strings.Contains(statusOutput, "is connected"),
	}

	return &status, nil
}

var (
	licenseUsernameRe   = regexp.MustCompile("Logged in as (.+)")
	licenseSubTypeRe    = regexp.MustCompile("You are using the (FREE|PREMIUM) version")
	maxDevicesOnRe      = regexp.MustCompile("Up to (.+) devices simultaneously")
	licenseValidUntilRe = regexp.MustCompile("Your subscription will be renewed on (.+)")
)

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
