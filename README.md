# 📱 SMSGate Go Client

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]
[![Module][module-shield]][module-url]

`client-go` provides a typed Go client for the SMSGate 3rd-party API and a separate client for the SMSGate Certificate Authority API. It is built on the Go standard library and has no runtime dependencies.

## 📖 About the Project

The module contains three packages:

- `smsgateway` provides the SMSGate 3rd-party API client for messages, inbox, devices, health, logs, settings, webhooks, and JWT tokens.
- `ca` submits Certificate Signing Requests (CSRs) and retrieves their status.
- `rest` provides the shared low-level HTTP client and error classification used by the other packages.

## 📚 Table of Contents

- [📱 SMSGate Go Client](#-smsgate-go-client)
	- [📖 About the Project](#-about-the-project)
	- [📚 Table of Contents](#-table-of-contents)
	- [⭐ Features](#-features)
	- [🚀 Getting Started](#-getting-started)
		- [Prerequisites](#prerequisites)
		- [Installation](#installation)
		- [Authentication](#authentication)
			- [Basic Authentication](#basic-authentication)
			- [JWT Authentication](#jwt-authentication)
	- [🚀 Quickstart](#-quickstart)
	- [💻 Usage](#-usage)
	- [Configuration](#configuration)
		- [`smsgateway.Config`](#smsgatewayconfig)
		- [CA Client Options](#ca-client-options)
		- [Versioned Public-Key DTOs](#versioned-public-key-dtos)
	- [📖 API Reference](#-api-reference)
	- [🤝 Contributing](#-contributing)
	- [📄 License](#-license)

## ⭐ Features

- Text, data, and outgoing MMS messages with base64 attachment support
- Message priority, device and SIM selection, delivery reports, TTL, validity windows, and scheduled delivery
- Message state retrieval, listing, filtering, sorting, content inclusion, and cancellation
- Inbox listing and refresh with pagination, filters, MMS attachment metadata, encryption-state flags, and disabled, individual, or batch webhook delivery
- Device listing and deletion, health checks, logs, and device settings with get, patch, and replace operations
- Webhook registration with typed event constants and payload DTOs
- JWT access-token generation, refresh, and revocation with scopes and TTL
- Optional versioned public-key fields on device and mobile request DTOs for E2E key metadata
- CSR submission and status retrieval through the `ca` package
- Custom HTTP clients and base URLs for testing or private deployments
- Go 1.22+ and zero runtime dependencies

## 🚀 Getting Started

### Prerequisites

- Go 1.22 or newer
- Credentials for the target SMSGate API

### Installation

```bash
go get github.com/android-sms-gateway/client-go
```

The module version is available on [pkg.go.dev](https://pkg.go.dev/github.com/android-sms-gateway/client-go). The Go version requirement is declared in [`go.mod`](go.mod).

### Authentication

The `smsgateway` client supports Basic authentication and JWT bearer tokens. A configured token takes priority over Basic credentials. The `ca` client does not use these SMSGate 3rd-party credentials.

#### Basic Authentication

```go
client := smsgateway.NewClient(smsgateway.Config{
	User:     os.Getenv("ASG_USERNAME"),
	Password: os.Getenv("ASG_PASSWORD"),
})
```

The `ASG_USERNAME` and `ASG_PASSWORD` names are examples for the consuming application. The library itself does not read environment variables.

#### JWT Authentication

Create a Basic-authenticated client to request a scoped token, then create a client with the returned access token:

```go
client := smsgateway.NewClient(smsgateway.Config{
	User:     os.Getenv("ASG_USERNAME"),
	Password: os.Getenv("ASG_PASSWORD"),
})

token, err := client.GenerateToken(context.Background(), smsgateway.TokenRequest{
	Scopes: []smsgateway.JWTScope{
		smsgateway.ScopeMessagesSend,
		smsgateway.ScopeMessagesRead,
	},
	TTL: 3600,
})
if err != nil {
	panic(err)
}

jwtClient := smsgateway.NewClient(smsgateway.Config{
	Token: token.AccessToken,
})
```

The available scopes are defined in [`smsgateway/domain_auth.go`](smsgateway/domain_auth.go).

## 🚀 Quickstart

The following example sends a text message using credentials supplied by the application:

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func main() {
	client := smsgateway.NewClient(smsgateway.Config{
		User:     os.Getenv("ASG_USERNAME"),
		Password: os.Getenv("ASG_PASSWORD"),
	})

	message := smsgateway.Message{
		PhoneNumbers: []string{"+15555550100"},
		TextMessage:  &smsgateway.TextMessage{Text: "Hello from Go"},
	}

	state, err := client.Send(context.Background(), message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("message queued: %s", state.ID)
}
```

## 💻 Usage

`Client.Send` validates `SendOptions` before making the request.

Use send options to control phone-number validation and active-device filtering:

```go
state, err := client.Send(ctx, message,
	smsgateway.WithSkipPhoneValidation(false),
	smsgateway.WithDeviceActiveWithin(24),
)
if err != nil {
	return err
}
log.Printf("message state: %s", state.State)
```

List messages with shared date, pagination, state, device, sort, and content options:

```go
limit := 50
includeContent := true
messages, total, err := client.ListMessages(ctx, smsgateway.ListMessagesOptions{
	PaginationOptions: smsgateway.PaginationOptions{Limit: &limit},
	IncludeContent:   &includeContent,
})
if err != nil {
	return err
}
log.Printf("loaded %d of %d messages", len(messages), total)
```

`RefreshInbox` supports disabled, individual, and batch webhook delivery. `ExportInbox`, `MessagesExportRequest`, and `TriggerWebhooks` are deprecated in favor of the newer inbox API.

The `ca` client exposes `PostCSR` and `GetCSRStatus`. The `smsgateway/webhooks` package is retained for compatibility; new code should use the webhook types in the `smsgateway` package.

## Configuration

### `smsgateway.Config`

| Field      | Type           | Default                                | Purpose                                                     |
| ---------- | -------------- | -------------------------------------- | ----------------------------------------------------------- |
| `Client`   | `*http.Client` | `http.DefaultClient`                   | Custom HTTP transport, proxy, TLS, or timeout.              |
| `BaseURL`  | `string`       | `https://api.sms-gate.app/3rdparty/v1` | Base URL for the 3rd-party API.                             |
| `User`     | `string`       | `""`                                   | Basic-authentication username.                              |
| `Password` | `string`       | `""`                                   | Basic-authentication password.                              |
| `Token`    | `string`       | `""`                                   | JWT bearer token; takes priority over Basic authentication. |

The exported `WithClient`, `WithBaseURL`, `WithBasicAuth`, and `WithJWTAuth` methods return updated `Config` values.

### CA Client Options

The `ca` package accepts functional options:

- `ca.WithClient` supplies a custom `*http.Client` and defaults to `http.DefaultClient`.
- `ca.WithBaseURL` overrides the CA API base URL, which defaults to `https://ca.sms-gate.app/api/v1`.

### Versioned Public-Key DTOs

`VersionedPublicKey` models optional E2E key metadata:

- `PublicKey` is a base64-encoded RSA public key.
- `KeyVersion` identifies the key version and is intended for rotation tracking.

The type is embedded in `Device` and `MobileRegisterRequest`; `MobileUpdateRequest` receives the fields through its embedded `MobileRegisterRequest`. These fields provide API schema and serialization support only. The client does not generate keys, encrypt or decrypt message payloads, or automatically attach keys to messages.

## 📖 API Reference

- [Official API Reference](https://docs.sms-gate.app/integration/api/) — endpoints, payloads, and error codes
- [Authentication Guide](https://docs.sms-gate.app/integration/authentication/) — scopes and token management
- [Client libraries overview](https://docs.sms-gate.app/integration/client-libraries/)
- [`smsgateway`](smsgateway/) — third-party API client and request/response types
- [`ca`](ca/) — Certificate Authority client
- [`rest`](rest/) — shared HTTP client and error types

## 🤝 Contributing

Contributions are welcome. Pull requests target the `master` branch. Before submitting a change, run:

```bash
make lint
make test
```

The CI workflow also runs linting and race-enabled tests with coverage.

## 📄 License

Distributed under the Apache License 2.0. See [LICENSE](LICENSE).

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
[module-shield]: https://img.shields.io/pkg.go.dev/v/github.com/android-sms-gateway/client-go?style=for-the-badge
[module-url]: https://pkg.go.dev/github.com/android-sms-gateway/client-go
