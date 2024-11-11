package cluster

import (
	"errors"
	"net"
	"os"
)

const (
	envPodIP      = "POD_IP"
	envExternalIP = "EXTERNAL_IP"
)

func IP() (string, error) {
	if val := os.Getenv(envPodIP); val != "" {
		return val, nil
	}
	if val := os.Getenv(envExternalIP); val != "" {
		return val, nil
	}
	return ExternalIP()
}

func ExternalIP(names ...string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	for _, iface := range interfaces {
		if name != "" && iface.Name != name {
			continue
		}
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsPrivate() || ip.IsLoopback() {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", errors.Join(err, errors.New("external ip not found"))
}
