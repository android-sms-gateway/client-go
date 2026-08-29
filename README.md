# 📱 SMSGate Go Client

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]
[![Go Version][version-shield]][version-url]

A typed Go client for the [SMSGate](https://sms-gate.app) API: send and track SMS messages through your Android devices with Basic or JWT authentication. Built on the Go standard library only, with zero runtime dependencies. See the [client libraries overview](https://docs.sms-gate.app/integration/client-libraries/) for the full ecosystem.

## 📖 About

`client-go` provides typed clients for the SMSGate ecosystem: the `smsgateway` package covers the 3rd-party API (messages, inbox, devices, health, logs, settings, webhooks, and token lifecycle), the `ca` package submits Certificate Signing Requests, and the shared `rest` package handles low-level HTTP and error classification. A bearer token, when set, takes priority over Basic credentials ([config.go](https://github.com/android-sms-gateway/client-go/blob/master/smsgateway/config.go)).

## 📚 Table of Contents

- [📱 SMSGate Go Client](#-smsgate-go-client)
	- [📖 About](#-about)
	- [📚 Table of Contents](#-table-of-contents)
	- [⭐ Features](#-features)
	- [📦 Installation](#-installation)
	- [🔑 Authentication](#-authentication)
		- [Basic Authentication](#basic-authentication)
		- [JWT Authentication](#jwt-authentication)
	- [🚀 Quickstart](#-quickstart)
	- [💻 Usage](#-usage)
	- [📖 API Reference](#-api-reference)
	- [🤝 Contributing](#-contributing)
	- [📄 License](#-license)

## ⭐ Features

- Text and data messages with priority, TTL, delivery reports, and scheduling
- Per-message and per-recipient state tracking, listing, filtering, and cancellation
- Inbox listing with pagination, inbox refresh, and webhook-based export
- Device management, health checks, logs, and device settings (get, patch, replace)
- Webhook registration with typed event constants and payload types
- JWT token lifecycle: generate, refresh, and revoke with scopes and TTL
- `ca` package for Certificate Signing Requests (webhook and private-server types)
- Customizable HTTP client and base URL for testing or private deployments
- Pure Go standard library, zero runtime dependencies (Go 1.22+)

## 📦 Installation

```bash
go get github.com/android-sms-gateway/client-go
```

Requires Go 1.22 or newer (see [go.mod](https://github.com/android-sms-gateway/client-go/blob/master/go.mod)).

## 🔑 Authentication

Two methods are supported: Basic authentication with account credentials, and JWT bearer tokens with scoped permissions. JWT is recommended for production.

### Basic Authentication

```go
client := smsgateway.NewClient(smsgateway.Config{
	User:     os.Getenv("ASG_USERNAME"),
	Password: os.Getenv("ASG_PASSWORD"),
})
```

### JWT Authentication

```go
token, err := client.GenerateToken(context.Background(), smsgateway.TokenRequest{
	Scopes: []smsgateway.JWTScope{smsgateway.ScopeMessagesSend, smsgateway.ScopeMessagesRead},
	TTL:    3600,
})
if err != nil {
	panic(err)
}

jwtClient := smsgateway.NewClient(smsgateway.Config{Token: token.AccessToken})
```

## 🚀 Quickstart

```go
package main

import (
	"context"
	"os"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func main() {
	client := smsgateway.NewClient(smsgateway.Config{User: os.Getenv("ASG_USERNAME"), Password: os.Getenv("ASG_PASSWORD")})

	msg := smsgateway.Message{PhoneNumbers: []string{"+15555550100"}, TextMessage: &smsgateway.TextMessage{Text: "Hello from Go"}}

	state, err := client.Send(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	println("message queued:", state.ID)
}
```

## 💻 Usage

Beyond sending, the client covers message listing and cancellation, inbox listing and refresh, device management, health checks, logs, settings (get, patch, replace), webhooks, and the full token lifecycle. See [smsgateway/client.go](https://github.com/android-sms-gateway/client-go/blob/master/smsgateway/client.go) for the complete method list with signatures, and the CA client in [ca/client.go](https://github.com/android-sms-gateway/client-go/blob/master/ca/client.go). API failures are wrapped in sentinel errors from the `rest` package; classify them with `errors.Is`.

## 📖 API Reference

- [Official API Reference](https://docs.sms-gate.app/integration/api/) - endpoints, payloads, and error codes
- [Authentication Guide](https://docs.sms-gate.app/integration/authentication/) - scopes and token management
- [Client libraries overview](https://docs.sms-gate.app/integration/client-libraries/)
- [Client source](https://github.com/android-sms-gateway/client-go/blob/master/smsgateway/client.go) - full method reference and examples

## 🤝 Contributing

Contributions are welcome. Open an issue to discuss major changes before submitting a pull request; PRs target the `master` branch.

## 📄 License

Distributed under the Apache License 2.0. See [LICENSE](https://github.com/android-sms-gateway/client-go/blob/master/LICENSE).

<!-- Badge references: Shields.io style=for-the-badge is mandatory -->
[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/client-go?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/client-go/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/client-go?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/client-go/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/client-go?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/client-go/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/client-go?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/client-go/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/client-go?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/client-go/blob/master/LICENSE
[version-shield]: https://img.shields.io/pkg.go.dev/v/github.com/android-sms-gateway/client-go?style=for-the-badge
[version-url]: https://pkg.go.dev/github.com/android-sms-gateway/client-go
