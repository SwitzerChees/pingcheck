# Pingcheck

The `pingcheck` is a Go application designed to provide simple endpoints for ping checks. The application reads a YAML file containing a list of checks and pings and provides an endpoint to check the status of the checks.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

What things you need to install the software and how to install them:

- Go (version 1.20 or later)

### Installing

A step-by-step series of examples that tell you how to get a development environment running.

1. **Clone the repository**

   ```bash
   git clone https://github.com/switzerchees/stock-publisher.git
   cd stock-publisher
   ```

### Running

1. **Build the application**

   ```bash
   go build .
   ```

2. **Run the application**

   ```bash
   ./pingcheck
   ```

### Environment Variables

- `CHECKS_FILE`: The Path to the checks file (default: `data/checks.yml`)
- `PINGS_FILE`: The Path to the pings file (default: `data/pings.yml`)

## Docker

### Building your own Docker Image

```bash
docker build -t pingcheck .
```

### Running the Docker Image

```bash
docker run -p 1234:1234 -v ./checks.yml:/app/data/checks.yml pingcheck
```

### Official Docker Image

```bash
docker run -p 1234:1234 switzerchees/pingcheck
```
