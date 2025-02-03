package proxy

import (
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"resty.dev/v3"
)

type Proxy struct {
	id      string
	rest    *resty.Client
	proxies []ProxyDetail
	current int
	mu      sync.Mutex
}

func (p *Proxy) ID() string {
	return p.id
}

func (p *Proxy) init() {
	token := os.Getenv("PROXY_TOKEN")
	p.rest = resty.New()
	resp, err := p.rest.R().
		SetHeader("Authorization", fmt.Sprintf("Token %s", token)).
		SetQueryParam("mode", "direct").
		SetResult(&ProxyResponse{}).
		Get("https://proxy.webshare.io/api/v2/proxy/list/")
	if err != nil {
		panic(err)
	}

	result := resp.Result().(*ProxyResponse)
	p.proxies = result.Results

	log.WithFields(log.Fields{
		"id":      p.id,
		"proxies": len(p.proxies),
	}).Info("proxy initialized")
}

func (p *Proxy) GetProxy() string {
	pass := os.Getenv("PROXY_PASS")
	p.mu.Lock()
	proxy := p.proxies[p.current]
	p.current = (p.current + 1) % len(p.proxies)
	p.mu.Unlock()
	return fmt.Sprintf("http://%s:%s@%s:%d", pass, pass, proxy.ProxyAddress, proxy.Port)
}

func Create(id string) *Proxy {
	p := &Proxy{id: id}
	p.init()
	return p
}
