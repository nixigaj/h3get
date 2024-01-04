# `h3get`
A dead simple curl-like HTTP/3 client tool I use for debugging HTTP/3 servers. It can be set to explicitly use IPv4 or IPv6.

> [!NOTE]  
> This program uses a [modified version](https://github.com/nixigaj/quic-go) of the [quic-go library](https://github.com/quic-go/quic-go) to allow for explicitly setting the network string for the [`ListenUDP`](https://pkg.go.dev/net#ListenUDP) function through the `QUIC_GO_CLIENT_NETWORK_TYPE` environment variable.

## Build
### Dependencies
- Git
- Go 1.21 or later

### Instructions
Clone the repository and enter it:
```shell
git clone https://github.com/nixigaj/h3get
cd h3get
```

Build the binary:
```shell
go build --ldflags="-w -s"
```

## Usage
To specify URL, use the `--url` or `-u` flag.

To specify explicit usage of IPv4 or IPv6, use the `--ipv4` or `-4`, and `--ipv6` or `-6` flags.

To specify timeout for request in seconds, use the `--timeout` or `-t` flag.

To use `curl` as user agent, use the `--curl` or `-c` flag.

To print usage use the `--help` or `-h` flag.

### Example
To query the URL [`https://h3.erix.dev`](https://h3.erix.dev) using IPv6 explicitly and a request timeout of 5 seconds:
```shell
./h3get -u https://h3.erix.dev -6 -t 5
```

## License
MIT. See [LICENSE](LICENSE).
