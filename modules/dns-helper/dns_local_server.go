package dns_helper

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var testDomain = "gruntwork.io"

// dnsDatabase holds DNSAnswers to a collection of DNSQuery, to be used by a local dnsTestServer
type dnsDatabase map[DNSQuery]DNSAnswers

// dnsTestServer helper for testing this package using local DNS nameservers with test records
type dnsTestServer struct {
	Server           *dns.Server
	DNSDatabase      dnsDatabase
	DNSDatabaseRetry dnsDatabase
}

// newDNSTestServer returns a new instance of dnsTestServer
func newDNSTestServer(server *dns.Server) *dnsTestServer {
	return &dnsTestServer{Server: server, DNSDatabase: make(dnsDatabase), DNSDatabaseRetry: make(dnsDatabase)}
}

// Address returns the host:port string of the server listener
func (s *dnsTestServer) Address() string {
	return s.Server.PacketConn.LocalAddr().String()
}

// AddEntryToDNSDatabase adds DNSAnswers to the DNSQuery in the database of the server
func (s *dnsTestServer) AddEntryToDNSDatabase(q DNSQuery, a DNSAnswers) {
	s.DNSDatabase[q] = append(s.DNSDatabase[q], a...)
}

// AddEntryToDNSDatabaseRetry adds DNSAnswers to the DNSQuery in the database used when retrying
func (s *dnsTestServer) AddEntryToDNSDatabaseRetry(q DNSQuery, a DNSAnswers) {
	s.DNSDatabaseRetry[q] = append(s.DNSDatabaseRetry[q], a...)
}

func setupTestDNSServers() (s1, s2 *dnsTestServer) {
	s1 = runTestDNSServer("0")
	s2 = runTestDNSServer("0")

	q := DNSQuery{"NS", testDomain}
	a := DNSAnswers{{"NS", s1.Address() + "."}, {"NS", s2.Address() + "."}}
	s1.AddEntryToDNSDatabase(q, a)
	s2.AddEntryToDNSDatabase(q, a)

	s1.Server.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) {
		stdDNSHandler(w, r, s1, false)
	})
	s2.Server.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) {
		stdDNSHandler(w, r, s2, true)
	})

	return s1, s2
}

func setupTestDNSServersRetry() (s1, s2 *dnsTestServer) {
	s1 = runTestDNSServer("0")
	s2 = runTestDNSServer("0")

	q := DNSQuery{"NS", testDomain}
	a := DNSAnswers{{"NS", s1.Address() + "."}, {"NS", s2.Address() + "."}}
	s1.AddEntryToDNSDatabase(q, a)
	s2.AddEntryToDNSDatabase(q, a)
	s1.AddEntryToDNSDatabaseRetry(q, a)
	s2.AddEntryToDNSDatabaseRetry(q, a)

	s1.Server.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) {
		retryDNSHandler(w, r, s1, false)
	})
	s2.Server.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) {
		retryDNSHandler(w, r, s2, true)
	})

	return s1, s2
}

func runTestDNSServer(port string) *dnsTestServer {
	listener, err := net.ListenPacket("udp", "127.0.0.1:"+port)

	if err != nil {
		log.Fatal(err)
	}

	mux := dns.NewServeMux()
	server := &dns.Server{PacketConn: listener, Net: "udp", Handler: mux}

	go func() {
		if err := server.ActivateAndServe(); err != nil {
			log.Printf("Error in local DNS server: %s", err)
		}
	}()

	return newDNSTestServer(server)
}

func doDNSAnswer(w dns.ResponseWriter, r *dns.Msg, d dnsDatabase, invertAnswers bool) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	q := m.Question[0]
	qtype := dns.TypeToString[q.Qtype]
	answers := d[DNSQuery{qtype, strings.TrimSuffix(q.Name, ".")}]

	var seen = make(map[DNSAnswer]bool)

	for _, r := range answers {
		if seen[r] {
			continue
		}
		seen[r] = true

		rr, err := dns.NewRR(fmt.Sprintf("%s %s", q.Name, r.String()))

		if err != nil {
			log.Fatalf("err: %s", err)
		}

		m.Answer = append(m.Answer, rr)
	}

	if invertAnswers {
		for i, j := 0, len(m.Answer)-1; i < j; i, j = i+1, j-1 {
			m.Answer[i], m.Answer[j] = m.Answer[j], m.Answer[i]
		}
	}

	w.WriteMsg(m)
}

func stdDNSHandler(w dns.ResponseWriter, r *dns.Msg, s *dnsTestServer, invertAnswers bool) {
	doDNSAnswer(w, r, s.DNSDatabase, invertAnswers)
}

var startTime = time.Now()

func retryDNSHandler(w dns.ResponseWriter, r *dns.Msg, s *dnsTestServer, invertAnswers bool) {
	if time.Now().Sub(startTime).Seconds() > 3 {
		doDNSAnswer(w, r, s.DNSDatabaseRetry, invertAnswers)
	} else {
		doDNSAnswer(w, r, s.DNSDatabase, invertAnswers)
	}
}
