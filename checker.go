package dnsbench

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
)

type checker struct {
	client *dns.Client
	config *dns.ClientConfig
	conn   *dns.Conn
	msg    *dns.Msg
}

func (c *checker) dnsQuery(domain string) (*dns.Msg, error) {
	c.msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	r, _, err := c.client.ExchangeWithConn(c.msg, c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s: %w", domain, err)
	}
	return r, nil //nolint: nlreturn
}

func (c *checker) Check(wi workItem) (*Result, error) {
	numQueries := c.config.Attempts

	ms := make(Measurements, 0, numQueries)
	for i := 0; i < numQueries; i++ {
		start := time.Now()

		_, err := c.dnsQuery(wi.Domain.String())
		if err != nil {
			return nil, err
		}

		end := time.Now()
		duration := end.Sub(start)

		ms = append(ms, Measurement{
			Duration: duration,
		})
	}

	return &Result{
		WorkItem:     wi,
		Measurements: ms,
	}, nil
}

func newChecker(c *dns.ClientConfig, proto string) (*checker, error) {
	//nolint: exhaustivestruct
	client := &dns.Client{
		Timeout: DNSTimeout,
		Net:     proto,
	}
	server := c.Servers[0]
	port := c.Port
	socket := fmt.Sprintf("%s:%s", server, port)

	conn, err := client.Dial(socket)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", socket, err)
	}

	//nolint: exhaustivestruct
	msg := &dns.Msg{}

	return &checker{
		client: client,
		config: c,
		conn:   conn,
		msg:    msg,
	}, nil
}
