package rpc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpPostJson(proxyUrl, httpUrl, bodyJson string) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		netTransport := &http.Transport{
			Proxy:                 http.ProxyURL(proxy),
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second * time.Duration(10),
		}
		httpClient.Transport = netTransport
	}
	b := strings.NewReader(bodyJson)
	res, err := httpClient.Post(httpUrl, "application/json", b)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("StatusCode: %d", res.StatusCode))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

}
