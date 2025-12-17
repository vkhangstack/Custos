package main

import (
	"log"

	"github.com/miekg/dns"
)

const upstream = "8.8.8.8:5353"

func handleDNS(w dns.ResponseWriter, r *dns.Msg) {
	c := new(dns.Client)
	resp, _, err := c.Exchange(r, upstream)
	if err != nil {
		log.Println("upstream error:", err)
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
		_ = w.WriteMsg(m)
		return
	}
	resp.Id = r.Id
	_ = w.WriteMsg(resp)
}

func StartDNS() {
	dns.HandleFunc(".", handleDNS)

	srv := &dns.Server{Addr: ":5352", Net: "udp"}
	log.Println("DNS proxy on :5352 ->", upstream)
	log.Fatal(srv.ListenAndServe())
}
