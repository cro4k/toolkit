// Package clients
// The code of this file is copied from github.com/gin-gonic/gin
package clients

import (
	"net"
	"net/http"
	"strings"
)

// RemoteIP parses the IP from Request.RemoteAddr, normalizes and returns the IP (without the port).
func RemoteIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

// ClientIP implements one best effort algorithm to return the real client IP.
// It calls c.RemoteIP() under the hood, to check if the remote IP is a trusted proxy or not.
// If it is it will then try to parse the headers defined in Engine.RemoteIPHeaders (defaulting to [X-Forwarded-For, X-Real-Ip]).
// If the headers are not syntactically valid OR the remote IP does not correspond to a trusted proxy,
// the remote IP (coming from Request.RemoteAddr) is returned.
func (c *IPResolver) ClientIP(request *http.Request) string {
	// Check if we're running on a trusted platform, continue running backwards if error
	if c.TrustedPlatform != "" {
		// Developers can define their own header of Trusted Platform or use predefined constants
		if addr := request.Header.Get(c.TrustedPlatform); addr != "" {
			return addr
		}
	}

	// It also checks if the remoteIP is a trusted proxy or not.
	// In order to perform this validation, it will see if the IP is contained within at least one of the CIDR blocks
	// defined by Engine.SetTrustedProxies()
	remoteIP := net.ParseIP(RemoteIP(request))
	if remoteIP == nil {
		return ""
	}
	trusted := c.isTrustedProxy(remoteIP)

	if trusted && c.ForwardedByClientIP && c.RemoteIPHeaders != nil {
		for _, headerName := range c.RemoteIPHeaders {
			ip, valid := c.validateHeader(request.Header.Get(headerName))
			if valid {
				return ip
			}
		}
	}
	return remoteIP.String()
}

// isTrustedProxy will check whether the IP address is included in the trusted list according to Engine.trustedCIDRs
func (c *IPResolver) isTrustedProxy(ip net.IP) bool {
	if c.trustedCIDRs == nil {
		return false
	}
	for _, cidr := range c.trustedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

// validateHeader will parse X-Forwarded-For header and return the trusted client IP address
func (c *IPResolver) validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, ",")
	for i := len(items) - 1; i >= 0; i-- {
		ipStr := strings.TrimSpace(items[i])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		// X-Forwarded-For is appended by proxy
		// Check IPs in reverse order and stop when find untrusted proxy
		if (i == 0) || (!c.isTrustedProxy(ip)) {
			return ipStr, true
		}
	}
	return "", false
}

// SetTrustedProxies set a list of network origins (IPv4 addresses,
// IPv4 CIDRs, IPv6 addresses or IPv6 CIDRs) from which to trust
// request's headers that contain alternative client IP when
// `(*gin.Engine).ForwardedByClientIP` is `true`. `TrustedProxies`
// feature is enabled by default, and it also trusts all proxies
// by default. If you want to disable this feature, use
// Engine.SetTrustedProxies(nil), then Context.ClientIP() will
// return the remote address directly.
func (c *IPResolver) SetTrustedProxies(trustedProxies []string) error {
	c.trustedProxies = trustedProxies
	return c.parseTrustedProxies()
}

// parseTrustedProxies parse Engine.trustedProxies to Engine.trustedCIDRs
func (c *IPResolver) parseTrustedProxies() error {
	trustedCIDRs, err := c.prepareTrustedCIDRs()
	c.trustedCIDRs = trustedCIDRs
	return err
}

func (c *IPResolver) prepareTrustedCIDRs() ([]*net.IPNet, error) {
	if c.trustedProxies == nil {
		return nil, nil
	}

	cidr := make([]*net.IPNet, 0, len(c.trustedProxies))
	for _, trustedProxy := range c.trustedProxies {
		if !strings.Contains(trustedProxy, "/") {
			ip := parseIP(trustedProxy)
			if ip == nil {
				return cidr, &net.ParseError{Type: "IP address", Text: trustedProxy}
			}

			switch len(ip) {
			case net.IPv4len:
				trustedProxy += "/32"
			case net.IPv6len:
				trustedProxy += "/128"
			}
		}
		_, cidrNet, err := net.ParseCIDR(trustedProxy)
		if err != nil {
			return cidr, err
		}
		cidr = append(cidr, cidrNet)
	}
	return cidr, nil
}

// parseIP parse a string representation of an IP and returns a net.IP with the
// minimum byte representation or nil if input is invalid.
func parseIP(ip string) net.IP {
	parsedIP := net.ParseIP(ip)

	if ipv4 := parsedIP.To4(); ipv4 != nil {
		// return ip in a 4-byte representation
		return ipv4
	}

	// return ip in a 16-byte representation or nil
	return parsedIP
}

type IPResolver struct {
	// TrustedPlatform if set to a constant of value gin.Platform*, trusts the headers set by
	// that platform, for example to determine the client IP
	TrustedPlatform string
	// ForwardedByClientIP if enabled, client IP will be parsed from the request's headers that
	// match those stored at `(*gin.Engine).RemoteIPHeaders`. If no IP was
	// fetched, it falls back to the IP obtained from
	// `(*gin.Context).Request.RemoteAddr`.
	ForwardedByClientIP bool

	// AppEngine was deprecated.
	// Deprecated: USE `TrustedPlatform` WITH VALUE `gin.PlatformGoogleAppEngine` INSTEAD
	// #726 #755 If enabled, it will trust some headers starting with
	// 'X-AppEngine...' for better integration with that PaaS.
	// AppEngine bool

	// RemoteIPHeaders list of headers used to obtain the client IP when
	// `(*gin.Engine).ForwardedByClientIP` is `true` and
	// `(*gin.Context).Request.RemoteAddr` is matched by at least one of the
	// network origins of list defined by `(*gin.Engine).SetTrustedProxies()`.
	RemoteIPHeaders []string

	trustedCIDRs   []*net.IPNet
	trustedProxies []string
}

var DefaultIPResolver = &IPResolver{}

func SetTrustedProxies(trustedProxies []string) error {
	return DefaultIPResolver.SetTrustedProxies(trustedProxies)
}

func SetTrustedPlatform(trustedPlatform string) {
	DefaultIPResolver.TrustedPlatform = trustedPlatform
}

func SetForwardedByClientIP(forwardedByClientIP bool) {
	DefaultIPResolver.ForwardedByClientIP = forwardedByClientIP
}

func SetRemoteIPHeaders(remoteIPHeaders []string) {
	DefaultIPResolver.RemoteIPHeaders = remoteIPHeaders
}

func ClientIP(request *http.Request) string {
	return DefaultIPResolver.ClientIP(request)
}
