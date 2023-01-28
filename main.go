package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func main() {
	start := net.ParseIP("104.16.0.0")
	end := net.ParseIP("104.31.255.255")

	for ip := start; ip.To4() != nil && !ip.Equal(end); inc(ip) {
		testIP(ip)
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func testIP(ip net.IP) {
	fmt.Print(ip.String() + " ")

	url := "https://archlinux.cloudflaremirrors.com/archlinux/iso/latest/archlinux-x86_64.iso"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, ip.String()+":443")
		},
	}
	client := &http.Client{Transport: transport}

	startTime := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	size := 1 << 20
	data := make([]byte, size)
	nr, err := io.ReadAtLeast(resp.Body, data, size)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(nr, time.Since(startTime))
}
