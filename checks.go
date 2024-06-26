package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

type Check struct {
	Name   string `yaml:"name"`
	Slug   string `yaml:"slug"`
	Token  string `yaml:"token"`
	Period int    `yaml:"period"`
}

type CheckHandler struct {
	checks    []Check
	pings     []Ping
	pingsFile string
}

func (h *CheckHandler) HandleCheck(c echo.Context) error {
	slug := c.Param("slug")
	token := c.Param("token")
	fmt.Printf("Check: %s \n", slug)
	if slug == "" || token == "" {
		return c.String(http.StatusBadRequest, "Slug and token are required")
	}

	check := FindCheckBySlug(slug, h.checks)
	if check == nil {
		return c.String(http.StatusNotFound, "Check not found")
	}
	if check.Token != token {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	ping := FindPingsById(slug, h.pings)
	if ping == nil {
		return c.String(http.StatusNotFound, "Ping not found")
	}
	if ping.LastPing == nil {
		return c.String(http.StatusInternalServerError, "No pings yet")
	}
	if time.Since(*ping.LastPing) > time.Duration(check.Period)*time.Minute {
		ping.Status = PingStatusDown
		SavePings(h.pingsFile, h.pings)
		return c.String(http.StatusInternalServerError, string(ping.Status))
	} else {
		ping.Status = PingStatusUp
		SavePings(h.pingsFile, h.pings)
	}
	return c.String(http.StatusOK, string(ping.Status))
}

func FindCheckBySlug(slug string, checks []Check) *Check {
	for i := range checks {
		if checks[i].Slug == slug {
			return &checks[i]
		}
	}
	return nil
}

func LoadChecks(path string) ([]Check, error) {
	var checks []Check
	data, err := os.ReadFile(path)
	if err != nil {
		return checks, err
	}
	err = yaml.Unmarshal(data, &checks)
	if err != nil {
		return checks, err
	}
	return checks, nil
}
