package protocol

import (
	"net/http"
	neturl "net/url"
)

func (c *Client) ImportCookies(cookies map[string]string, domain string) {
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
	// Make cookie slice
	cks := make([]*http.Cookie, 0, len(cookies))
	for name, value := range cookies {
		cookie := &http.Cookie{
			Name:     name,
			Value:    value,
			Path:     "/",
			Domain:   domain,
			HttpOnly: true,
		}
		cks = append(cks, cookie)
	}
	// Save cookie
	c.cj.SetCookies(url, cks)
}

func (c *Client) ExportCookies(url string) map[string]string {
	u, _ := neturl.Parse(url)
	cookies := make(map[string]string)
	for _, ck := range c.cj.Cookies(u) {
		cookies[ck.Name] = ck.Value
	}
	return cookies
}
