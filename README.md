# TimeOff

TimeOff is a web application for tracking a worker's annual leave (vacation and public holidays),
sickness and time in lieu.

## Requirements

### TLS

The application can be run with TLS support turned off but, due to the way
TimeOff configures cookies with the `Secure` flag, the application must be
frontend by a server providing TLS termination.

## HTTP/2

Timeoff supports the HTTP/2 protocol over TLS for browsers that understand it.
Other browsers will continue to use HTTP/1.1.

## Getting started

The application supports Heroku and Docker for deployment to Production using the
provided `Procfile` and `Dockerfile` respectively. For deployment to Heroku, use
the [Go build pack](https://github.com/kr/heroku-buildpack-go).

For a development environment:

    # Ensure you have Go 1.4 installed
    godep go get -t -v ./...
    godep go build

    # Generate TLS certificates for development use only
    go run $GOROOT/src/crypto/tls/generate_cert.go -host localhost

    # Run the application with debug logging enabled
    # Replace the GitHub client ID and secret with your own:
    # https://github.com/settings/applications
    DEBUG=1 ADDR=localhost:3000 TLS_CERT=cert.pem TLS_KEY=key.pem GITHUB_CLIENT_ID=abc GITHUB_CLIENT_SECRET=xyz ./timeoff

## Project status

Currently in alpha. Many features are not yet implemented.

### Roadmap

#### Version 0.1 - MVP

- Track annual leave (no approvals yet)
- Annotations (e.g. "holiday to Spain")
- Adjustments (rolled-over allowance)

Limitations/assumptions:

- one, full time, job

### Future TODO

- Sharing (read-only view)
- Approvals for holiday requests
- Ability to change leave year
- ICAL export (including leave for those who have shared it)
- CSV export
- JSON export
- Support for multiple jobs
- Manager request approvals
- Welcome email
- Google Oauth login
- Public holidays/bank holidays
- Sickness
- Payment in lieu (above statutory minimum)
- Custom leave day types (e.g. privilege days)
- Notifications to employee and manager to notify of time not taken, nearing end of year, etc.
- Currently assumes full time, need to allow for part-time
