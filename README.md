<!-- Anchor for back to top links -->
<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Go Report Card][reportcard-shield]][reportcard-url]
[![Codecov][codecov-shield]][codecov-url]
[![Go Version][goversion-shield]][goversion-url]
[![License][license-shield]][license-url]
[![Release][release-shield]][release-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]

<!-- TABLE OF CONTENTS -->
- [📱 About The Project](#-about-the-project)
- [🌟 Features](#-features)
- [⚙️ Prerequisites](#️-prerequisites)
- [📦 Installation](#-installation)
- [🛠️ Usage Examples](#️-usage-examples)
- [📚 API Reference](#-api-reference)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)


<!-- ABOUT THE PROJECT -->
## 📱 About The Project

This is a Go client library for interfacing with the [SMS Gateway for Android™ API](https://sms-gate.app). It provides a simple and efficient way to integrate SMS functionality into your Go applications, with features like message sending, status checking, and webhook management.

Key value propositions:
- Lightweight and easy to integrate
- Comprehensive API coverage
- Built with Go best practices

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- FEATURES -->
## 🌟 Features

- Send SMS messages with a simple method call.
- Check the state of sent messages.
- Webhooks management.
- Customizable base URL for use with local, cloud or private servers.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- PREREQUISITES -->
## ⚙️ Prerequisites

Before you begin, ensure you have met the following requirements:
- You have a basic understanding of Go.
- You have Go installed on your local machine.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- INSTALLATION -->
## 📦 Installation

To install the SMS Gateway API Client in the existing project, run this command in your terminal:

```bash
go get github.com/android-sms-gateway/client-go
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## 🛠️ Usage Examples

Here's how to get started with the SMS Gateway API Client:

1. Import the [`github.com/android-sms-gateway/client-go/smsgateway`](smsgateway/) package.
2. Create a new client with configuration using the `smsgateway.NewClient()` method.
3. Use the `Send()` method to send an SMS message.
4. Use the `GetState()` method to check the status of a sent message.

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
		User:     os.Getenv("ASG_USERNAME"),
		Password: os.Getenv("ASG_PASSWORD"),
	})

	// Create a message to send.
	message := smsgateway.Message{
		TextMessage: &smsgateway.TextMessage{
			Text: "Hello, doctors!",
		},
		PhoneNumbers: []string{
			"+15555550100",
			"+15555550101",
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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- API REFERENCE -->
## 📚 API Reference

For more information on the API endpoints and data structures, please consult the [SMS Gateway for Android API documentation](https://api.sms-gate.app/).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## 🤝 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## 📄 License

Distributed under the Apache-2.0 License. See [`LICENSE`](LICENSE) for more information.

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
