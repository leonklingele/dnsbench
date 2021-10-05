package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/leonklingele/dnsbench"
)

const (
	defaultProtocol   = "udp"
	defaultNumQueries = 64
	defaultNumWorkers = 1
)

var (
	//go:embed domains.txt
	defaultDomains string

	//nolint: gochecknoglobals
	defaultServers = []string{
		"1.1.1.1",
		"8.8.4.4",
		"8.8.8.8",
		"9.9.9.9",
	}
)

func run() error {
	defaultDomainsFlag := strings.Join(strings.Split(defaultDomains, "\n"), ", ")
	if maxLen, l := 52, len(defaultDomainsFlag); l > maxLen {
		defaultDomainsFlag = defaultDomainsFlag[:maxLen] + "…"
	}
	defaultServersFlag := strings.Join(defaultServers, ", ")

	domains := flag.String("domains", defaultDomainsFlag, "comma-separated list of domain names to benchmark")
	servers := flag.String("servers", defaultServersFlag, "comma-separated list of DNS servers to run benchmark against")
	proto := flag.String("protocol", defaultProtocol, "protocol to use (either tcp or udp)")
	numQueries := flag.Int("queries", defaultNumQueries, "number of queries for each domain name")
	numWorkers := flag.Int("workers", defaultNumWorkers, "number of concurrent workers")
	flag.Parse()

	// Domains
	var ds []dnsbench.Domain
	if *domains == defaultDomainsFlag {
		// Use default domains
		rawDomains := strings.Split(defaultDomains, "\n")

		for _, rd := range rawDomains {
			rd = strings.TrimSpace(rd)
			rd = strings.ToLower(rd)

			if len(rd) == 0 || strings.HasPrefix(rd, "#") {
				continue
			}

			ds = append(ds, dnsbench.NewDomain(rd))
		}
	} else {
		// Use user-provided domains
		rawDomains := strings.Split(*domains, ",")

		for _, rd := range rawDomains {
			rd = strings.TrimSpace(rd)
			rd = strings.ToLower(rd)

			if len(rd) == 0 {
				continue
			}

			ds = append(ds, dnsbench.NewDomain(rd))
		}
	}

	// Servers
	var ss []dnsbench.Server
	if *servers == defaultServersFlag {
		// Use default servers
		for _, ds := range defaultServers {
			s, err := dnsbench.NewServer(ds)
			if err != nil {
				return err //nolint: wrapcheck
			}

			ss = append(ss, s)
		}
	} else {
		// Use user-provided servers
		rawServers := strings.Split(*servers, ",")

		for _, rs := range rawServers {
			rs = strings.TrimSpace(rs)
			rs = strings.ToLower(rs)

			if len(rs) == 0 {
				continue
			}

			s, err := dnsbench.NewServer(rs)
			if err != nil {
				return err //nolint: wrapcheck
			}

			ss = append(ss, s)
		}
	}

	// Print config
	{
		domains := make([]string, 0, len(ds))
		for _, d := range ds {
			domains = append(domains, d.String())
		}

		servers := make([]string, 0, len(ss))
		for _, s := range ss {
			servers = append(servers, s.String())
		}

		//nolint: forbidigo
		fmt.Printf(
			"Domains: %s\nServers: %s\nProto:   %s\nQueries: %d\nWorkers: %d\n",
			strings.Join(domains, ", "),
			strings.Join(servers, ", "),
			*proto,
			*numQueries,
			*numWorkers,
		)
	}

	var maxDomainLen int
	for _, d := range ds {
		if l := len(d.String()); l > maxDomainLen {
			maxDomainLen = l
		}
	}

	type result struct {
		duration time.Duration
		server   dnsbench.Server
	}

	results := make([]result, 0, len(ss))
	for _, s := range ss {
		c, err := dnsbench.Bench(ds, s, *proto, *numQueries, *numWorkers)
		if err != nil {
			return err //nolint: wrapcheck
		}

		var acc time.Duration
		var acci int
		for res := range c {
			avg := res.Measurements.AverageDuration()
			acc += avg
			acci++

			domainFormatter := fmt.Sprintf("%%-%ds", maxDomainLen+1)
			//nolint: forbidigo
			fmt.Printf(
				"[%s]: avg query time for "+domainFormatter+": %s\n",
				s,
				res.WorkItem.Domain,
				avg,
			)
		}

		if acci == 0 {
			continue
		}

		d := time.Duration(acc.Nanoseconds() / int64(acci))
		results = append(results, result{
			duration: d,
			server:   s,
		})
	}

	// Print result
	{
		sort.Slice(results, func(i, j int) bool {
			return results[i].duration < results[j].duration
		})

		for _, r := range results {
			//nolint: forbidigo
			fmt.Printf(
				"Summary [%s]: avg query time: %s\n",
				r.server,
				r.duration.String(),
			)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
