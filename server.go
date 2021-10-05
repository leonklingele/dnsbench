package dnsbench

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultPort = 53
)

type Server struct {
	Server string
	Port   int
}

func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.Server, s.Port)
}

func NewServer(s string) (Server, error) {
	split := strings.Split(s, ":")

	server := split[0]
	port := defaultPort
	if nel := 2; len(split) == nel {
		p := split[1]

		var err error
		port, err = strconv.Atoi(p)
		if err != nil {
			return Server{}, fmt.Errorf("failed to parse port %s: %w", p, err)
		}
	}

	return Server{
		Server: server,
		Port:   port,
	}, nil
}
