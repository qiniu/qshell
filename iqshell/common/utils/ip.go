package utils

import (
	"net"
	"strings"
)

func IsIPString(host string) bool {
	if len(host) == 0 {
		return false
	}
	return net.ParseIP(host) != nil
}

func IsIPUrlString(host string) bool {
	host = RemoveUrlScheme(host)
	for i := 0; i < len(host); i++ {
		switch host[i] {
		case '.':
			return isIPV4UrlString(host)
		case ':':
			return isIPV6UrlString(host)
		}
	}
	return false
}

func isIPV4UrlString(s string) bool {
	for i := 0; i < net.IPv4len; i++ {
		if len(s) == 0 {
			return false
		}
		if i > 0 {
			if i < net.IPv4len && s[0] != '.' {
				return false
			}
			s = s[1:]
		}
		n, c, ok := dtoi(s)
		if !ok || (i < net.IPv4len && n > 0xFF) {
			return false
		}
		if c > 1 && s[0] == '0' {
			// Reject non-zero components with leading zeroes.
			return false
		}
		s = s[c:]
	}
	return true
}

// Bigger than we need, not too big to worry about overflow
const big = 0xFFFFFF

// Decimal to integer.
// Returns number, characters consumed, success.
func dtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= big {
			return big, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

func isIPV6UrlString(host string) bool {
	host = strings.ReplaceAll(host, "[", "")
	host = strings.ReplaceAll(host, "]", "")

	ellipsis := -1 // position of ellipsis in ip

	// Might have leading ellipsis
	if len(host) >= 2 && host[0] == ':' && host[1] == ':' {
		ellipsis = 0
		host = host[2:]
		// Might be only ellipsis
		if len(host) == 0 {
			return true
		}
	}

	// Loop, parsing hex numbers followed by colon.
	i := 0
	for i < net.IPv6len {
		// Hex number.
		n, c, ok := xtoi(host)
		if !ok || n > 0xFFFF {
			return false
		}

		// If followed by dot, might be in trailing IPv4.
		if c < len(host) && host[c] == '.' {
			if ellipsis < 0 && i != net.IPv6len-net.IPv4len {
				// Not the right place.
				return false
			}
			if i+net.IPv4len > net.IPv6len {
				// Not enough room.
				return false
			}
			return isIPV4UrlString(host)
		}
		i += 2
		// Stop at end of string.
		host = host[c:]
		if len(host) == 0 {
			break
		}

		// Otherwise must be followed by colon and more.
		if host[0] != ':' || len(host) == 1 {
			return false
		}
		host = host[1:]

		// Look for ellipsis.
		if host[0] == ':' {
			if ellipsis >= 0 { // already have one
				return false
			}
			ellipsis = i
			host = host[1:]
			if len(host) == 0 { // can be at end
				break
			}
		}
	}

	// If didn't parse enough, expand ellipsis.
	if i < net.IPv6len {
		if ellipsis < 0 {
			return false
		}
		return true
	} else if ellipsis >= 0 {
		// Ellipsis must represent at least one 0 group.
		return false
	}
	return true
}

// Hexadecimal to integer.
// Returns number, characters consumed, success.
func xtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i]-'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i]-'A') + 10
		} else {
			break
		}
		if n >= big {
			return 0, i, false
		}
	}
	if i == 0 {
		return 0, i, false
	}
	return n, i, true
}
