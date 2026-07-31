<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Go Report Card][reportcard-shield]][reportcard-url]
[![Codecov][codecov-shield]][codecov-url]
[![Go Version][goversion-shield]][goversion-url]
[![License][license-shield]][license-url]
[![Release][release-shield]][release-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/android-sms-gateway/client-go">
    <img src="https://github.com/capcom6/android-sms-gateway/raw/master/assets/logo.png" alt="Logo" width="100" height="100">
  </a>

<h3 align="center">client-go</h3>

  <p align="center">
    Go client library for the SMSGate APIs.
    <br />
    <a href="https://api.sms-gate.app/"><strong>Explore API docs &gt;</strong></a>
    <br />
    <br />
    <a href="https://github.com/android-sms-gateway/client-go/issues">Report Bug</a>
    |
    <a href="https://github.com/android-sms-gateway/client-go/issues">Request Feature</a>
  </p>
</div>

## Table of Contents

- [Table of Contents](#table-of-contents)
- [About The Project](#about-the-project)
- [Features](#features)
- [Getting Started](#getting-started)
	- [Prerequisites](#prerequisites)
	- [Installation](#installation)
- [Usage](#usage)
	- [Sending a text message](#sending-a-text-message)
	- [Sending a data message](#sending-a-data-message)
	- [Listing incoming messages](#listing-incoming-messages)
	- [Certificate Authority client](#certificate-authority-client)
	- [Error handling](#error-handling)
- [Configuration](#configuration)
	- [`smsgateway.Config`](#smsgatewayconfig)
	- [`ca` package](#ca-package)
	- [Custom base URLs](#custom-base-urls)
- [API Coverage](#api-coverage)
	- [`smsgateway.Client`](#smsgatewayclient)
	- [`ca.Client`](#caclient)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About The Project

`client-go` provides typed Go clients for the SMSGate ecosystem:

- `smsgateway` - client for the SMSGate 3rd-party API (messages, inbox, devices, health, logs, settings, webhooks, and token lifecycle).
- `ca` - client for the SMSGate Certificate Authority service (submit CSRs and poll their status).
- `rest` - shared low-level HTTP handling used by both clients.

The library supports Basic authentication (`user` + `password`) and Bearer token authentication. A token, when set, takes priority over Basic credentials.

The library is built on the Go standard library only and has zero runtime dependencies.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Features

- Send text and data SMS messages, cancel pending messages, and track per-message and per-recipient state.
- List messages with filtering, pagination, content inclusion, and sorting.
- Read incoming messages from the device inbox and trigger inbox refreshes.
- List and delete registered devices.
- Check API health and retrieve device log entries.
- Read, partially update, and fully replace device settings.
- Register, list, and delete webhooks, with typed event constants and webhook payload types.
- Generate, refresh, and revoke API tokens with scopes and TTL.
- Submit Certificate Signing Requests and poll CSR status (`webhook` and `private_server` types).
- Inject a custom `http.Client` and override the base URL for testing or private deployments.
- Classify API errors with `errors.Is` and the sentinel errors from the `rest` package.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Getting Started

### Prerequisites

- Go 1.22 or newer (see [`go.mod`](go.mod)).
- An SMSGate account with device credentials (username/password) or an API token.

### Installation

```bash
go get github.com/android-sms-gateway/client-go
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

### Sending a text message

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

func main() {
	ctx := context.Background()

	client := smsgateway.NewClient(smsgateway.Config{
		User:     os.Getenv("ASG_USERNAME"),
		Password: os.Getenv("ASG_PASSWORD"),
		// or use Token: os.Getenv("ASG_TOKEN"),
	})

	state, err := client.Send(ctx, smsgateway.Message{
		TextMessage: &smsgateway.TextMessage{Text: "Hello from Go"},
		PhoneNumbers: []string{
			"+15555550100",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("message queued: %s", state.ID)
}
```

`Send` accepts optional `SendOption` values, for example `smsgateway.WithSkipPhoneValidation(true)` and `smsgateway.WithDeviceActiveWithin(24)`.

### Sending a data message

```go
state, err := client.Send(ctx, smsgateway.Message{
	DataMessage: &smsgateway.DataMessage{
		Data: "SGVsbG8gV29ybGQh", // base64-encoded payload
		Port: 53739,
	},
	PhoneNumbers: []string{
		"+15555550100",
	},
})
if err != nil {
	log.Fatal(err)
}
log.Printf("message queued: %s", state.ID)
```

### Listing incoming messages

```go
limit := 50

inbox, total, err := client.ListInboxMessages(ctx, smsgateway.ListInboxOptions{
	Limit: &limit,
})
if err != nil {
	log.Fatal(err)
}
log.Printf("%d messages of %d total", len(inbox), total)
```

`ListMessages` follows the same pattern for outgoing messages and returns the total count from the `X-Total-Count` response header.

### Certificate Authority client

```go
package main

import (
	"context"
	"log"

	"github.com/android-sms-gateway/client-go/ca"
)

func main() {
	ctx := context.Background()
	client := ca.NewClient()

	resp, err := client.PostCSR(ctx, ca.PostCSRRequest{
		Type:    ca.CSRTypeWebhook,
		Content: "-----BEGIN CERTIFICATE REQUEST-----...",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("request id: %s, status: %s", resp.RequestID, resp.Status)

	status, err := client.GetCSRStatus(ctx, resp.RequestID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("CSR status: %s", status.Status.Description())
}
```

### Error handling

The `rest` package exposes sentinel errors that you can match with `errors.Is` or the provided helpers:

- `rest.ErrAPIError` — base API error; `rest.IsAPIError` for any API failure
- `rest.ErrClient` / `rest.IsClientError` — 4xx client errors
- `rest.ErrServer` / `rest.IsServerError` — 5xx server errors
- `rest.ErrBadRequest` / `rest.IsBadRequest` — 400 validation failures
- `rest.ErrConflict` / `rest.IsConflict` — 409 conflicts

Continuing the `smsgateway` example above:

```go
import (
	"errors"

	"github.com/android-sms-gateway/client-go/rest"
)

_, err := client.Send(ctx, msg)
switch {
case errors.Is(err, rest.ErrBadRequest):
	// 400: message payload rejected
case errors.Is(err, rest.ErrConflict):
	// 409: conflicts with current state
case errors.Is(err, rest.ErrServer):
	// 5xx: service unavailable
default:
	// other client or transport errors
}
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Configuration

The library reads no environment variables; credentials and options are set on the config structs in code.

### `smsgateway.Config`

| Field      | Type           | Default                                | Description                                  |
| ---------- | -------------- | -------------------------------------- | -------------------------------------------- |
| `Client`   | `*http.Client` | `http.DefaultClient`                   | HTTP client used for requests                |
| `BaseURL`  | `string`       | `https://api.sms-gate.app/3rdparty/v1` | API base URL (constant `BaseURL`)            |
| `User`     | `string`       | empty                                  | Basic auth username                          |
| `Password` | `string`       | empty                                  | Basic auth password                          |
| `Token`    | `string`       | empty                                  | Bearer token, takes priority over Basic auth |

Chained helpers are available for the same fields: `Config.WithClient`, `Config.WithBaseURL`, `Config.WithBasicAuth(user, password)`, and `Config.WithJWTAuth(token)`. `Config.Validate()` reports an error wrapping `ErrInvalidConfig` when no credentials are set.

### `ca` package

The CA client is configured with functional options:

| Option        | Description                                                          |
| ------------- | -------------------------------------------------------------------- |
| `WithClient`  | Sets the HTTP client (defaults to `http.DefaultClient`)              |
| `WithBaseURL` | Sets the API base URL (defaults to `https://ca.sms-gate.app/api/v1`) |

### Custom base URLs

Override `BaseURL` (or use `WithBaseURL`) to point the clients at a private deployment or a mock server:

```go
client := smsgateway.NewClient(smsgateway.Config{
	Token:   os.Getenv("ASG_TOKEN"),
	BaseURL: "https://example.com/3rdparty/v1",
})
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## API Coverage

Endpoint semantics and payload details: <https://api.sms-gate.app/>

### `smsgateway.Client`

| Area     | Methods                                                                                  |
| -------- | ---------------------------------------------------------------------------------------- |
| Messages | `Send`, `CancelMessage`, `GetState`, `ListMessages`                                      |
| Inbox    | `ListInboxMessages`, `RefreshInbox` (`ExportInbox` is deprecated)                        |
| Devices  | `ListDevices`, `DeleteDevice`                                                            |
| Health   | `CheckHealth`                                                                            |
| Logs     | `GetLogs`                                                                                |
| Settings | `GetSettings`, `UpdateSettings`, `ReplaceSettings`                                       |
| Webhooks | `ListWebhooks`, `RegisterWebhook`, `DeleteWebhook`, plus event constants in `smsgateway` |
| Tokens   | `GenerateToken`, `RefreshToken`, `RevokeToken`                                           |

Typed webhook payloads (for example `PushNotification`, `MmsReceivedPayload`, `SmsBatchReceivedPayload`) live in the `smsgateway` package. The `smsgateway/webhooks` subpackage exists only as a deprecated compatibility alias.

### `ca.Client`

| Area | Methods                   |
| ---- | ------------------------- |
| CSR  | `PostCSR`, `GetCSRStatus` |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

Contributions are welcome. Please open an issue to discuss major changes before submitting a pull request.

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/my-change`).
3. Commit your changes (`git commit -m 'Describe change'`).
4. Push to your branch and open a pull request against `master`.

Pull requests are checked by CI (see [`.github/workflows/go.yml`](.github/workflows/go.yml)): `golangci-lint` and `go test -race` with coverage, reported to Codecov.

Run the same checks locally:

```bash
make lint
make test
make coverage
make benchmark
```

`make help` lists all available targets.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## License

Distributed under the Apache License 2.0. See [`LICENSE`](LICENSE) for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contact

- Repository: <https://github.com/android-sms-gateway/client-go>
- Issues: <https://github.com/android-sms-gateway/client-go/issues>

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[reportcard-shield]: https://goreportcard.com/badge/github.com/android-sms-gateway/client-go?style=for-the-badge
[reportcard-url]: https://goreportcard.com/report/github.com/android-sms-gateway/client-go
[codecov-shield]: https://img.shields.io/codecov/c/gh/android-sms-gateway/client-go?style=for-the-badge
[codecov-url]: https://codecov.io/gh/android-sms-gateway/client-go
[goversion-shield]: https://img.shields.io/github/go-mod/go-version/android-sms-gateway/client-go?style=for-the-badge
[goversion-url]: https://github.com/android-sms-gateway/client-go/blob/HEAD/go.mod
[license-shield]: https://img.shields.io/badge/License-Apache_2.0-blue.svg?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/client-go/blob/master/LICENSE
[release-shield]: https://img.shields.io/github/v/release/android-sms-gateway/client-go?style=for-the-badge
[release-url]: https://github.com/android-sms-gateway/client-go/releases
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/client-go?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/client-go/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/client-go?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/client-go/issues