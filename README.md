# openrvs-stats

Two components:

1. A stats server which outputs JSON for all active servers at :8081/servers
1. A web app (HTML+CSS+JS) which acts on the JSON

Run the stats server as a Go app under systemd.
Run the web app behind Nginx.