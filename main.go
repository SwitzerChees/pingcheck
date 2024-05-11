package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	checksFile := "checks.yml"
	if os.Getenv("CHECKS_FILE") != "" {
		checksFile = os.Getenv("CHECKS_FILE")
	}
	checks, err := LoadChecks(checksFile)
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}

	pingsFile := "pings.yml"
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
		checks: checks,
		pings:  pings,
	}
	app.GET("/ping/:slug/:token", pingHandler.HandlePing)

	checkHandler := &CheckHandler{
		checks: checks,
		pings:  pings,
	}
	app.GET("/check/:slug/:token", checkHandler.HandleCheck)

	app.Logger.Fatal(app.Start(":1234"))
}
