package testutils

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"net"
	"time"
)

func isPortAvailable(host, port string) bool {
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return true
	}
	defer conn.Close()
	return false
}

func FindAvailablePort(hostname string) string {
	minPort := 1024
	maxPort := 49151

	for i := 0; i < 20; i++ {
		port := fmt.Sprint(rand.Intn(maxPort-minPort) + minPort)

		if isPortAvailable(hostname, port) {
			return port
		}
	}
	return ""
}

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
