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

// DNSData type
type DNSData map[DNSQuery][]DNSAnswer

func setupTestDNSServers() (s1, s2 *dns.Server, ns1, ns2 string, dnsData1, dnsData2 *DNSData) {
	s1 = runTestDNSServer("0")
	s2 = runTestDNSServer("0")

	ns1 = s1.PacketConn.LocalAddr().String()
	ns2 = s2.PacketConn.LocalAddr().String()

	q := DNSQuery{"NS", testDomain}
	dnsData1 = &DNSData{q: DNSAnswers{{"NS", ns1}, {"NS", ns2}}}
	dnsData2 = &DNSData{q: DNSAnswers{{"NS", ns1}, {"NS", ns2}}}

	s1.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) { stdDNSHandler(w, r, dnsData1, false) })
	s2.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) { stdDNSHandler(w, r, dnsData2, true) })

	return s1, s2, ns1, ns2, dnsData1, dnsData2
}

func setupTestDNSServersRetry() (s1, s2 *dns.Server, ns1, ns2 string, dnsData1, dnsData2, dnsDataRetry1, dnsDataRetry2 *DNSData) {
	s1 = runTestDNSServer("0")
	s2 = runTestDNSServer("0")

	ns1 = s1.PacketConn.LocalAddr().String()
	ns2 = s2.PacketConn.LocalAddr().String()

	q := DNSQuery{"NS", testDomain}
	dnsData1 = &DNSData{q: DNSAnswers{{"NS", ns1}, {"NS", ns2}}}
	dnsData2 = &DNSData{q: DNSAnswers{{"NS", ns1}, {"NS", ns2}}}
	dnsDataRetry1 = &DNSData{q: DNSAnswers{{"NS", ns1}, {"NS", ns2}}}
	dnsDataRetry2 = &DNSData{q: DNSAnswers{{"NS", ns1}, {"NS", ns2}}}

	s1.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) { retryDNSHandler(w, r, dnsData1, dnsDataRetry1, false) })
	s2.Handler.(*dns.ServeMux).HandleFunc(testDomain+".", func(w dns.ResponseWriter, r *dns.Msg) { retryDNSHandler(w, r, dnsData2, dnsDataRetry2, true) })

	return s1, s2, ns1, ns2, dnsData1, dnsData2, dnsDataRetry1, dnsDataRetry2
}

func runTestDNSServer(port string) *dns.Server {
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

	return server
}

func doDNSAnswer(w dns.ResponseWriter, r *dns.Msg, dnsData *DNSData, invertAnswers bool) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	q := m.Question[0]
	qtype := dns.TypeToString[q.Qtype]
	answers := (*dnsData)[DNSQuery{qtype, strings.TrimSuffix(q.Name, ".")}]

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

func stdDNSHandler(w dns.ResponseWriter, r *dns.Msg, dnsData *DNSData, invertAnswers bool) {
	doDNSAnswer(w, r, dnsData, invertAnswers)
}

var startTime = time.Now()

func retryDNSHandler(w dns.ResponseWriter, r *dns.Msg, dnsData, dnsDataRetry *DNSData, invertAnswers bool) {
	if time.Now().Sub(startTime).Seconds() > 3 {
		dnsData = dnsDataRetry
	}

	doDNSAnswer(w, r, dnsData, invertAnswers)
}
