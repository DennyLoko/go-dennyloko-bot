package bot

import (
	"math/rand"
	"net"
	"time"

	"github.com/miekg/dns"
)

type DNS struct {
	API     *dns.Client
	servers []string
	rand    *rand.Rand
}

func NewDNS() *DNS {
	return &DNS{
		API: &dns.Client{},
		servers: []string{
			"208.67.222.222",
			"208.67.220.220",
		},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (i *DNS) OwnAddress() ([]net.IP, error) {
	return i.lookupddress("myip.opendns.com")
}

func (i *DNS) lookupddress(host string) (ips []net.IP, err error) {
	m := dns.Msg{}
	m.SetQuestion(host+".", dns.TypeA)

	r, _, err := i.API.Exchange(&m, i.servers[i.rand.Intn(len(i.servers))]+":53")
	if err != nil {
		return
	}

	for _, rec := range r.Answer {
		if t, ok := rec.(*dns.A); ok {
			ips = append(ips, t.A)
		}
	}

	return
}
