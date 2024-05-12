package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

type PingStatus string

const (
	PingStatusUnknown PingStatus = "UNKNOWN"
	PingStatusUp      PingStatus = "UP"
	PingStatusDown    PingStatus = "DOWN"
)

type Ping struct {
	Id       string     `yaml:"id"`
	Status   PingStatus `yaml:"status"`
	LastPing *time.Time `yaml:"last_ping"`
}

type PingHandler struct {
	checks    []Check
	pings     []Ping
	pingsFile string
}

func (h *PingHandler) HandlePing(c echo.Context) error {
	slug := c.Param("slug")
	token := c.Param("token")
	fmt.Printf("Ping: %s \n", slug)
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
	now := time.Now()
	ping.LastPing = &now
	ping.Status = PingStatusUp
	SavePings(h.pingsFile, h.pings)
	return c.String(http.StatusOK, "Pong")
}

func FindPingsById(id string, pings []Ping) *Ping {
	for i := range pings {
		if pings[i].Id == id {
			return &pings[i]
		}
	}
	return nil
}

func LoadPings(path string, checks []Check) ([]Ping, error) {
	var pings []Ping
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return pings, err
		}
		err = yaml.Unmarshal(data, &pings)
		if err != nil {
			return pings, err
		}
	}
	var currentPings []Ping
	for _, p := range pings {
		for _, c := range checks {
			if p.Id == c.Slug {
				currentPings = append(currentPings, p)
				break
			}
		}
	}
	for _, c := range checks {
		found := false
		for _, p := range currentPings {
			if p.Id == c.Slug {
				found = true
				break
			}
		}
		if !found {
			currentPings = append(currentPings, Ping{
				Id:       c.Slug,
				Status:   PingStatusUnknown,
				LastPing: nil,
			})
		}
	}
	return currentPings, nil
}

func SavePings(path string, pings []Ping) error {
	data, err := yaml.Marshal(pings)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
