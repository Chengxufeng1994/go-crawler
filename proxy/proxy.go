package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"sync/atomic"
)

var (
	ErrProxyUrlLength = fmt.Errorf("proxy url list empty")
)

type ProxyFunc func(*http.Request) (*url.URL, error)

type roundRobinSwitcher struct {
	proxyUrls []*url.URL
	index     uint32
}

func (r roundRobinSwitcher) GetProxy(pr *http.Request) (*url.URL, error) {
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyUrls[index%uint32(len(r.proxyUrls))]
	return u, nil
}

func NewRoundRobinProxySwitcher(proxyUrls ...string) (ProxyFunc, error) {
	if len(proxyUrls) < 1 {
		return nil, ErrProxyUrlLength
	}

	urls := make([]*url.URL, len(proxyUrls))
	for i, u := range proxyUrls {
		parseUrl, err := url.Parse(u)
		if err != nil {
			return nil, err
		}

		urls[i] = parseUrl
	}

	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}
