# SMS Gateway for Android™ Go API Client

[![Go Report Card](https://goreportcard.com/badge/github.com/android-sms-gateway/client-go)](https://goreportcard.com/report/github.com/android-sms-gateway/client-go)
[![GoDoc](https://godoc.org/github.com/android-sms-gateway/client-go?status.svg)](https://godoc.org/github.com/android-sms-gateway/client-go)
[![codecov](https://codecov.io/gh/android-sms-gateway/client-go/branch/master/graph/badge.svg)](https://codecov.io/gh/android-sms-gateway/client-go)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/android-sms-gateway/client-go)
[![License](https://img.shields.io/github/license/android-sms-gateway/client-go)](https://github.com/android-sms-gateway/client-go/blob/master/LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/android-sms-gateway/client-go)](https://github.com/android-sms-gateway/client-go/releases)
[![GitHub stars](https://img.shields.io/github/stars/android-sms-gateway/client-go)](https://github.com/android-sms-gateway/client-go/stargazers)
![GitHub All Releases](https://img.shields.io/github/downloads/android-sms-gateway/client-go/total)
[![GitHub issues](https://img.shields.io/github/issues/android-sms-gateway/client-go)](https://github.com/android-sms-gateway/client-go/issues)

This is a Go client library for interfacing with the [SMS Gateway for Android API](https://sms-gate.app).

## Features

- Send SMS messages with a simple method call.
- Check the state of sent messages.
- Customizable base URL for use with local, cloud or private servers.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have a basic understanding of Go.
- You have Go installed on your local machine.

## Installation

To install the SMS Gateway API Client in the existing project, run this command in your terminal:

```bash
go get github.com/android-sms-gateway/client-go
```

## Usage

Here's how to get started with the SMS Gateway API Client:

1. Import the `github.com/android-sms-gateway/client-go/smsgateway` package.
2. Create a new client with configuration with `smsgateway.NewClient` method.
3. Use the `Send` method to send an SMS message.
4. Use the `GetState` method to check the status of a sent message.

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func main() {
	// Create a client with configuration from environment variables.
	client := smsgateway.NewClient(smsgateway.Config{
		User:     os.Getenv("GATEWAY_USER"),
		Password: os.Getenv("GATEWAY_PASSWORD"),
	})

	// Create a message to send.
	message := smsgateway.Message{
		Message: "Your SMS message text here",
		PhoneNumbers: []string{
			"+1234567890",
			"+9876543210",
		},
	}

	// Send the message and get the response.
	status, err := client.Send(context.Background(), message)
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	log.Printf("Send message response: %+v", status)

	// Get the state of the message and print the response.
	status, err = client.GetState(context.Background(), status.ID)
	if err != nil {
		log.Fatalf("failed to get state: %v", err)
	}

	log.Printf("Get state response: %+v", status)
}
```

## API Reference

For more information on the API endpoints and data structures, please consult the [SMS Gateway for Android API documentation](https://sms-gate.app/api).

# Contributing

Contributions are welcome! Please submit a pull request or create an issue for anything you'd like to add or change.

# License

This library is open-sourced software licensed under the [Apache-2.0 license](LICENSE).