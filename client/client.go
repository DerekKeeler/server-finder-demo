package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/DerekKeeler/server-finder-demo/server"
)

type ServerMapping struct {
	IP   string
	Resp server.AnnounceResponse
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip, nil
		}
	}
	return nil, errors.New("no network connection found")
}

func makeRequest(addr string) (server.AnnounceResponse, error) {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("http://%v:8080/announce", addr), nil)
	if err != nil {
		return server.AnnounceResponse{}, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return server.AnnounceResponse{}, err
	}

	var result server.AnnounceResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

func Scan() ([]ServerMapping, error) {
	responses := []ServerMapping{}

	currIP, err := externalIP()
	if err != nil {
		return nil, err
	}

	ip := currIP.To4()
	if ip == nil {
		return nil, errors.New("non ipv4 address")
	}

	ip = ip.Mask(ip.DefaultMask())
	responseValues := make(chan ServerMapping, 20)
	emptyValues := make(chan struct{}, 20)
	reqSemaphore := make(chan struct{}, 150)
	for i := uint8(255); i >= 2; i-- {
		dupIP := ip.To4()
		dupIP[3] = byte(i)
		reqAddr := dupIP.String()

		go func(c chan ServerMapping, emptyValues chan struct{}, reqSemaphore chan struct{}, reqAddr string) {
			reqSemaphore <- struct{}{}
			resp, err := makeRequest(reqAddr)

			<-reqSemaphore

			if err != nil {
				emptyValues <- struct{}{}
				return
			}

			c <- ServerMapping{
				IP:   reqAddr,
				Resp: resp,
			}
		}(responseValues, emptyValues, reqSemaphore, reqAddr)
	}

	recvCount := 0
	progress := progressOutput{
		total: 254,
	}
	for recvCount < 254 {
		select {
		case val := <-responseValues:
			responses = append(responses, val)
		case <-emptyValues:
		}

		recvCount++
		progress.writeProgress(recvCount)
	}

	return responses, nil
}
