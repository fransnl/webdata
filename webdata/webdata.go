package webdata

import (
	"crypto/tls"
	"encoding/json"
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
	"crypto/x509"
	"os"
	//"net/http/cookieJar"
)

func FixUrl(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}
	return "https://" + url
}

//not used
func tls_request() *http.Client {
	var (
		conn *tls.Conn
		err  error
	)

	tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig

	c := &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
			DisableKeepAlives:   false,

			TLSClientConfig: &tls.Config{
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_AES_128_GCM_SHA256,
					tls.VersionTLS13,
					tls.VersionTLS10,
				},
			},
			DialTLS: func(network, addr string) (net.Conn, error) {
				conn, err = tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}

	// returns client with the necessary setup bypass basic cloudflare checks
	return c
}

//not used
func tlsRequest() *http.Client {
	var (
		conn *tls.Conn
		err  error
	)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
			DisableKeepAlives:   false,
			TLSClientConfig:     tlsConfig,
			DialTLS: func(network, addr string) (net.Conn, error) {
				conn, err = tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}

	return client
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

//not used
func clientWithProxy() *http.Client {
	var (
		conn *tls.Conn
		err  error
	)

	proxyURL, err := url.Parse("http://brd-customer-hl_d424f34a-zone-linksrepo_proxy:iz5pg7j95j2j@brd.superproxy.io:33335")
	if err != nil {
		panic(err)
	}

	caCert, certErr := os.ReadFile("./brightdata.crt")
	if certErr != nil {
		panic(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig

	clientConfig := &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
		},
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyURL(proxyURL),
			TLSHandshakeTimeout: 30 * time.Second,
			DisableKeepAlives:   false,
			TLSClientConfig:     clientConfig,
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

func checkIfYoutube(req *http.Request) bool {
	targetURL := req.URL.String()

	u, err := url.Parse(targetURL)
	if err != nil {
		log.Println("Failed to parse URL:", err)
		return false
	}

	origin := u.Scheme + "://" + u.Host

	log.Println(origin)

	if origin == "https://www.youtube.com" || origin == "https://youtube.com" {
		return true
	}

	return false
}

func fallbackOEmbed(url string) (string, string) {

	u := "https://www.youtube.com/oembed?url=" + url + "&format=json"

	client := tlsRequestDefault() //tls_request()

	req, err := http.NewRequest("GET", u, nil)

	addHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Fallback oEmbed request failed:", err)
		return "", ""
	}
	defer resp.Body.Close()

	var result struct {
		ThumbnailURL string `json:"thumbnail_url"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		Type         string `json:"type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding fallback oEmbed:", err)
		return "", ""
	}

	return result.ThumbnailURL, result.Title
}

func GetWebData(url string) WebInfo {

	client := tlsRequestDefault() //clientWithProxy()//tlsRequestTest()//tls_request()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error while retrieving site 1", err)
		return WebInfo{}
	}

	if checkIfYoutube(req) {
		p, t := fallbackOEmbed(url)
		return WebInfo{url, t, "", "Youtube", p, "https://www.youtube.com/s/desktop/9fda8632/img/logos/favicon.ico"}
	}

	addHeaders(req)

	// //dump, _ := httputil.DumpRequestOut(req, true)
	// //log.Println("=== REQUEST ===\n", string(dump))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error while retrieving site 2", err)
		return WebInfo{}
	}

	defer resp.Body.Close()

	info := htmlinfo.NewHTMLInfo()
	err = info.Parse(resp.Body, &url, nil)
	if err != nil {
		log.Println("Info Parse error:", err)
		return WebInfo{}
	}

	//log.Println("=== info ===\n", info)
	oembed := info.GenerateOembedFor(url)
	//log.Println("=== OEMBED ===\n", oembed)

	webInfo := WebInfo{url, oembed.Title, oembed.Description, oembed.ProviderName, oembed.ThumbnailURL, info.TouchIcons[0].URL}

	return webInfo
}
