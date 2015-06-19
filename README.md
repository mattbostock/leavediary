[![Stories in Ready](https://badge.waffle.io/mattbostock/leavediary.png?label=ready&title=Ready)](https://waffle.io/mattbostock/leavediary)
[![Build Status](https://travis-ci.org/mattbostock/leavediary.svg?branch=master)](https://travis-ci.org/mattbostock/leavediary)
[![Go Report Card](http://goreportcard.com/badge/mattbostock/leavediary)](http://goreportcard.com/report/mattbostock/leavediary)

# LeaveDiary

LeaveDiary is a web application for tracking a annual leave (vacation or holiday).

## Project status

LeaveDiary is currently in alpha.

## Demo

You can try the application here:

https://leavediary.herokuapp.com/

## Roadmap

See the [Milestones](https://github.com/mattbostock/leavediary/milestones) and
[Issues](https://github.com/mattbostock/leavediary/issues) pages for planned
features.

## Requirements

### TLS

The application can be run with TLS support turned off but, due to the way
LeaveDiary configures cookies with the `Secure` flag, the application must be
frontend by a server providing TLS termination.

## HTTP/2

LeaveDiary supports the HTTP/2 protocol over TLS for browsers that understand it.
Other browsers will continue to use HTTP/1.1.

## Getting started

For a development environment:

    # Ensure you have Go 1.4 installed
    make

    # Run the application with debug logging enabled
    # Replace the GitHub client ID and secret with your own:
    # https://github.com/settings/applications
    DEBUG=1 ADDR=localhost:3000 TLS_CERT=test_fixtures/cert.pem TLS_KEY=test_fixtures/key.pem \
      GITHUB_CLIENT_ID=abc GITHUB_CLIENT_SECRET=xyz ./leavediary
