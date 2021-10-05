package dnsbench

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/miekg/dns"
)

const (
	DNSTimeout = 8 * time.Second
)

func Bench(domains []Domain, server Server, proto string, queries, workers int) (ResultChan, error) {
	servers := []string{server.Server}
	port := strconv.Itoa(server.Port)

	//nolint: exhaustivestruct
	config := &dns.ClientConfig{
		Servers:  servers,
		Port:     port,
		Timeout:  int(DNSTimeout.Seconds()),
		Attempts: queries,
	}

	var wg sync.WaitGroup
	wch := make(chan workItem)
	rch := make(chan Result)

	for i := 0; i < workers; i++ {
		c, err := newChecker(config, proto)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			for wi := range wch {
				res, err := c.Check(wi)
				if err != nil {
					log.Printf("failed to check: %v\n", err)
					continue //nolint: nlreturn
				}
				rch <- *res
			}
		}()
	}

	go func() {
		for _, d := range domains {
			wi := workItem{
				Domain: d,
			}
			wch <- wi
		}

		close(wch)
		wg.Wait()
		close(rch)
	}()

	return rch, nil
}
