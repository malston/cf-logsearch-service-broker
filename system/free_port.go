package system

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func FindFreePort() (int, error) {
	ln, _ := net.Listen("tcp", ":0")
	defer ln.Close()

	fmt.Fprintf(os.Stdout, "Listening for tcp traffic on %q", ln.Addr().String())

	parsedPort, parseErr := strconv.ParseInt(ln.Addr().String()[5:], 10, 32)
	if parseErr != nil {
		return -1, parseErr
	}

	return int(parsedPort), nil
}
