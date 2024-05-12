package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	checksFile := "data/checks.yml"
	if os.Getenv("CHECKS_FILE") != "" {
		checksFile = os.Getenv("CHECKS_FILE")
	}
	checks, err := LoadChecks(checksFile)
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}

	pingsFile := "data/pings.yml"
	if os.Getenv("PINGS_FILE") != "" {
		pingsFile = os.Getenv("PINGS_FILE")
	}
	pings, err := LoadPings(pingsFile, checks)
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
	SavePings(pingsFile, pings)

	pingHandler := &PingHandler{
		checks:    checks,
		pings:     pings,
		pingsFile: pingsFile,
	}
	app.GET("/ping/:slug/:token", pingHandler.HandlePing)

	checkHandler := &CheckHandler{
		checks:    checks,
		pings:     pings,
		pingsFile: pingsFile,
	}
	app.GET("/check/:slug/:token", checkHandler.HandleCheck)

	app.GET("", func(c echo.Context) error {
		return c.String(200, fmt.Sprintf("PingCheck Running! %s Checks Active.", strconv.Itoa(len(checks))))
	})

	app.Logger.Fatal(app.Start(":1234"))
}
