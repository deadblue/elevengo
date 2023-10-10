package protocol

import (
	"net/http"
	neturl "net/url"
)

func (c *ClientImpl) ImportCookies(cookies map[string]string, domains ...string) {
	for _, domain := range domains {
		c.internalImportCookies(cookies, domain, "/")
	}
}

func (c *ClientImpl) internalImportCookies(cookies map[string]string, domain string, path string) {
	// Make a dummy URL for saving cookie
	url := &neturl.URL{
		Scheme: "https",
		Path:   "/",
	}
	if domain[0] == '.' {
		url.Host = "www" + domain
	} else {
		url.Host = domain
	}
	// Prepare cookies
	cks := make([]*http.Cookie, 0, len(cookies))
	for name, value := range cookies {
		cookie := &http.Cookie{
			Name:     name,
			Value:    value,
			Domain:   domain,
			Path:     path,
			HttpOnly: true,
		}
		cks = append(cks, cookie)
	}
	// Save cookies
	c.cj.SetCookies(url, cks)
}

func (c *ClientImpl) ExportCookies(url string) map[string]string {
	u, _ := neturl.Parse(url)
	cookies := make(map[string]string)
	for _, ck := range c.cj.Cookies(u) {
		cookies[ck.Name] = ck.Value
	}
	return cookies
}
