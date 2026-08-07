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
    Go client library for SMSGate APIs.
    <br />
    <a href="https://api.sms-gate.app/"><strong>Explore API docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/android-sms-gateway/client-go/issues">Report Bug</a>
    ·
    <a href="https://github.com/android-sms-gateway/client-go/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
- [About The Project](#about-the-project)
- [Built With](#built-with)
- [Getting Started](#getting-started)
	- [Prerequisites](#prerequisites)
	- [Installation](#installation)
- [Usage](#usage)
	- [SMSGate client (`smsgateway`)](#smsgate-client-smsgateway)
	- [Certificate Authority client (`ca`)](#certificate-authority-client-ca)
- [API Coverage](#api-coverage)
	- [`smsgateway.Client`](#smsgatewayclient)
	- [`ca.Client`](#caclient)
- [Contributing](#contributing)
- [License](#license)


<!-- ABOUT THE PROJECT -->
## About The Project

`client-go` provides typed clients for the SMSGate ecosystem:

- `smsgateway` package for 3rd-party API operations (messages, devices, health, logs, settings, webhooks, and token lifecycle).
- `ca` package for Certificate Authority workflows (submit CSR and check CSR status).
- Shared low-level HTTP handling in the `rest` package.

The library supports both Basic authentication (`user` + `password`) and Bearer token authentication for the SMSGate client.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Built With

- [Go](https://go.dev/) 1.22+
- Standard `net/http` client with configurable transport

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

- Go 1.22 or newer
- SMSGate account/device credentials and/or API token

### Installation

```bash
go get github.com/android-sms-gateway/client-go
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

### SMSGate client (`smsgateway`)

#### Basic send

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

#### E2E encryption

End-to-end encryption is optional and injected into the client. Without an
encryptor, messages are sent as plaintext exactly as before and no device
listing lookup happens. To enable E2E, install an encryptor via
`Config.WithEncryptor`:

```go
client := smsgateway.NewClient(
    smsgateway.Config{
        User:     os.Getenv("ASG_USERNAME"),
        Password: os.Getenv("ASG_PASSWORD"),
    }.WithEncryptor(smsgateway.NewEncryptor()),
)
```

With an encryptor installed, E2E is applied **by default**: set `DeviceID` to
the target device so the client can resolve its `publicKey` + `keyVersion`
from the device listing, and the encryptor encrypts the body and every phone
number with hybrid RSA-OAEP (SHA-256) + AES-256-GCM and marks the message with
`IsEncrypted`.

```go
state, err := client.Send(ctx, smsgateway.Message{
    DeviceID:     "PyDmBQZZXYmyxMwED8Fzy",
    TextMessage:  &smsgateway.TextMessage{Text: "Secret payload"},
    PhoneNumbers: []string{"+15555550100"},
})
if err != nil {
    log.Fatal(err)
}
```

Wire format: `$rsa-oaep-aes-256-gcm$v=1$k={keyVersion}${base64(encrypted_aes_key)}${base64(iv)}${base64(ciphertext || tag)}`.

- `deviceId` is required for encryption; a missing `deviceId` with
  `WithE2EEncryption(true)` returns `ErrDeviceIDRequired`.
- Pass `WithE2EEncryption(true)` to require encryption: a device without a
  usable `publicKey`/`keyVersion` (or an unknown device) then returns the typed
  `ErrE2ENotConfigured` / `ErrDeviceNotFound`, and nothing is sent.
- In the default auto mode, devices without a key (or unknown devices) fall
  back to plaintext - `ErrE2ENotConfigured` is only ever returned when
  encryption is explicitly required.
- Pass `WithE2EEncryption(false)` to disable E2E for a message even when a
  device has a key.
- Passphrase encryption (`$aes-256-cbc/pbkdf2-sha1$`) is not implemented yet;
  the wire prefix is reserved and messages already carrying an encrypted-format
  prefix are always sent verbatim, never re-encrypted.

### Certificate Authority client (`ca`)

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
}
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## API Coverage

### `smsgateway.Client`

- Messages: `Send`, `GetState`
- Devices: `ListDevices`, `DeleteDevice`
- Health: `CheckHealth`
- Inbox export: `ExportInbox`
- Logs: `GetLogs`
- Settings: `GetSettings`, `UpdateSettings`, `ReplaceSettings`
- Webhooks: `ListWebhooks`, `RegisterWebhook`, `DeleteWebhook`
- Token management: `GenerateToken`, `RefreshToken`, `RevokeToken`

For endpoint semantics and payload details, see https://api.sms-gate.app/

### `ca.Client`

- CSR workflows: `PostCSR`, `GetCSRStatus`

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are welcome. Please open an issue to discuss major changes before submitting a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/my-change`)
3. Commit your changes (`git commit -m 'Describe change'`)
4. Push to your branch
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under the Apache-2.0 License. See [`LICENSE`](LICENSE) for details.

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
