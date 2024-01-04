package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

const (
	defaultClientTimeoutSec = 10

	quicGoDisableEcn = "true"
	quicGoDisableGso = "true"
)

var (
	urlFlag           string
	ipv4flag          bool
	ipv6flag          bool
	clientTimeoutFlag int
	curlFlag          bool
	helpFlag          bool
)

func main() {
	flag.StringVar(&urlFlag, "url", "", "Specify the URL for the QUIC client")
	flag.StringVar(&urlFlag, "u", "", "Specify the URL for the QUIC client (shorthand)")
	flag.BoolVar(&ipv4flag, "ipv4", false, "Use IPv4 for QUIC client")
	flag.BoolVar(&ipv4flag, "4", false, "Use IPv4 for QUIC client (shorthand)")
	flag.BoolVar(&ipv6flag, "ipv6", false, "Use IPv6 for QUIC client")
	flag.BoolVar(&ipv6flag, "6", false, "Use IPv6 for QUIC client (shorthand)")
	flag.IntVar(&clientTimeoutFlag, "timeout", defaultClientTimeoutSec, "Timeout for request in seconds")
	flag.IntVar(&clientTimeoutFlag, "t", defaultClientTimeoutSec, "Timeout for request in seconds (shorthand)")
	flag.BoolVar(&curlFlag, "curl", false, "Use 'curl' user agent")
	flag.BoolVar(&curlFlag, "c", false, "Use 'curl' user agent (shorthand)")
	flag.BoolVar(&helpFlag, "help", false, "Print usage")
	flag.BoolVar(&helpFlag, "h", false, "Print usage (shorthand)")
	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if urlFlag == "" {
		fmt.Println("URL flag is not set")
		flag.Usage()
		os.Exit(1)
	}

	// I had some issues if these environment variables were not set
	// See https://github.com/quic-go/quic-go/issues/3911
	_, ecnSet := os.LookupEnv("QUIC_GO_DISABLE_ECN")
	if !ecnSet {
		err := os.Setenv("QUIC_GO_DISABLE_ECN", quicGoDisableEcn)
		if err != nil {
			fmt.Println("[ERROR]:", err)
			os.Exit(1)
		}
	}
	_, gsoSet := os.LookupEnv("QUIC_GO_DISABLE_GSO")
	if !gsoSet {
		err := os.Setenv("QUIC_GO_DISABLE_GSO", quicGoDisableGso)
		if err != nil {
			fmt.Println("[ERROR]:", err)
			os.Exit(1)
		}
	}

	var err error
	if ipv4flag && ipv6flag {
		fmt.Println("conflicting IP versions in flags")
		os.Exit(1)
	} else if ipv4flag {
		err = os.Setenv("QUIC_GO_CLIENT_NETWORK_TYPE", "udp4")
	} else if ipv6flag {
		err = os.Setenv("QUIC_GO_CLIENT_NETWORK_TYPE", "udp6")
	}
	if err != nil {
		fmt.Println("[ERROR]:", err)
		os.Exit(1)
	}

	client := &http.Client{
		Transport: &http3.RoundTripper{
			Dial: quic.DialAddrEarly,
		},
		Timeout: time.Duration(clientTimeoutFlag) * time.Second,
	}

	req, err := http.NewRequest("GET", urlFlag, nil)
	if err != nil {
		fmt.Println("[ERROR]:", err)
		os.Exit(1)
	}

	if curlFlag {
		req.Header.Set("User-Agent", "curl")
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR]:", err)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[ERROR]:", err)
		os.Exit(1)
	}

	err = resp.Body.Close()
	if err != nil {
		fmt.Println("[ERROR]:", err)
		os.Exit(1)
	}

	fmt.Print(string(body))
}
