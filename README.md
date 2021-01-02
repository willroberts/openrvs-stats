# openrvs-stats

This repo contains:

- A Go service (`main.go`) which fetches healthy servers from [openrvs-registry](https://github.com/willroberts/openrvs-registry), fetches server info using [openrvs-beacon](https://github.com/willroberts/openrvs-beacon), and then serves the data as JSON over HTTP on port `8081`
- A build script (`build.bat`) which compiles the server for Windows (`stats.exe`) and Linux (`stats`)
- Rudimentary HTML, JS, and CSS (`web` directory) to display the JSON

## Running Stats App Locally

1. Start [openrvs-registry](https://github.com/willroberts/openrvs-registry) locally
2. Start the stats app (`go run main.go` or `build.bat; stats.exe`)
3. Visit http://localhost:8081 in your browser to get stats as JSON

## Testing Frontend Locally

1. Once the stats app is running, you can visit `file:///path/to/web/index.html` to test HTML/JS/CSS changes

## Deployment

The web frontend is currently deployed at http://64.225.54.237/live, where it runs alongside the stats app and registry
