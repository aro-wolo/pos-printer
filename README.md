# pos-printer

![GitHub Repo stars](https://img.shields.io/github/stars/aro-wolo/pos-printer?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/aro-wolo/pos-printer)
![GitHub issues](https://img.shields.io/github/issues/aro-wolo/pos-printer)
![GitHub contributors](https://img.shields.io/github/contributors/aro-wolo/pos-printer)
[![Go Report Card](https://goreportcard.com/badge/github.com/aro-wolo/pos-printer)](https://goreportcard.com/report/github.com/aro-wolo/pos-printer)
[![codecov](https://codecov.io/gh/aro-wolo/pos-printer/branch/main/graph/badge.svg)](https://codecov.io/gh/aro-wolo/pos-printer)
![Go Module](https://img.shields.io/github/go-mod/go-version/aro-wolo/pos-printer)
[![License](https://img.shields.io/github/license/aro-wolo/pos-printer.svg)](https://github.com/aro-wolo/pos-printer/blob/main/LICENSE)

## Overview

`pos-printer` is a Go application designed to fetch order details from an API and print a formatted receipt on a thermal printer. It uses environment variables for configuration and supports token-based authentication.

## Setup

### Prerequisites

- Go 1.16 or later
- A `.env` file with the necessary environment variables

### Environment Variables

Create a `.env` file in the root directory of the project with the following variables:

```
API_USERNAME=your_api_username
API_KEY=your_api_key
TOKEN_URL=https://example.com/get-token
ORDER_URL=https://example.com/orders/%s
HEADER_KEY=your_header_key
HEADER_VALUE=your_header_value
STORE_NAME=Your Store Name
STORE_LINK=https://yourstore.com
STORE_ADDRESS=Your Store Address
```

### Install Dependencies

Run the following command to install the required dependencies:

```sh
go mod tidy
```

## Usage

1. Build the application:

```sh
go build -o pos-printer
```

2. Run the application:

```sh
./pos-printer
```

3. Enter the order ID when prompted to fetch and print the receipt.

## Customization

- **Footer Template**: Customize the footer of the receipt by editing the `footer_template.txt` file.
- **Receipt Formatting**: Modify the `printReceipt` function in `main.go` to change the receipt format.

## License

This project is licensed under the MIT License.
