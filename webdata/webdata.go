package webdata

import (
	"crypto/tls"
	//"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dyatlov/go-htmlinfo/htmlinfo"

	//"io"
	//"net/http/httputil"
	//"fmt"
	//"crypto/x509"
	//"os"
	//"net/http/cookieJar"
)

func FixUrl(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}
	return "https://" + url
}

func tlsRequestDefault() *http.Client {
	var (
		conn *tls.Conn
		err  error
	)

	tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig

	clientConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
		},
	}

	// jar, cookieErr := cookiejar.New(nil)
	// if cookieErr != nil {
	// 	log.Fatal(err)
	// }

	client := &http.Client{
		//Jar: jar,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
			DisableKeepAlives:   false,
			TLSClientConfig:     clientConfig,
			ForceAttemptHTTP2:   false,
			DialTLS: func(network, addr string) (net.Conn, error) {
				conn, err = tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}

	return client
}

func addHeaders(req *http.Request) {
	targetURL := req.URL.String()

	u, err := url.Parse(targetURL)
	if err != nil {
		log.Println("Failed to parse URL:", err)
		return
	}

	origin := u.Scheme + "://" + u.Host
	req.Header.Set("Referer", origin+"/")
	req.Header.Set("Origin", origin)

	// Optional: Add realistic browser headers
	//req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}

func GetWebData(url string) (WebInfo, error) {

	client := http.DefaultClient//tlsRequestDefault() //clientWithProxy()//tlsRequestTest()//tls_request()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WebInfo{}, err
	}

	addHeaders(req)

	// //dump, _ := httputil.DumpRequestOut(req, true)
	// //log.Println("=== REQUEST ===\n", string(dump))

	resp, err := client.Do(req)
	if err != nil {
		return WebInfo{}, err
	}

	defer resp.Body.Close()

	info := htmlinfo.NewHTMLInfo()
	err = info.Parse(resp.Body, &url, nil)
	if err != nil {
		return WebInfo{}, err
	}

	//log.Println("=== info ===\n", info)
	oembed := info.GenerateOembedFor(url)
	//log.Println("=== OEMBED ===\n", oembed)

	webInfo := WebInfo{url, oembed.Title, oembed.Description, oembed.ProviderName, oembed.ThumbnailURL, info.FaviconURL}

	return webInfo, nil
}
