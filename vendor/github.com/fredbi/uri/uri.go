package uri

import (
	"bytes"
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// Validation errors
var (
	ErrNoSchemeFound    = errors.New("no scheme found in URI")
	ErrInvalidURI       = errors.New("not a valid URI")
	ErrInvalidCharacter = errors.New("invalid character in URI")
	ErrInvalidScheme    = errors.New("invalid scheme in URI")
	ErrInvalidQuery     = errors.New("invalid query string in URI")
	ErrInvalidFragment  = errors.New("invalid fragment in URI")
	ErrInvalidPath      = errors.New("invalid path in URI")
	ErrInvalidHost      = errors.New("invalid host in URI")
	ErrInvalidPort      = errors.New("invalid port in URI")
	ErrInvalidUserInfo  = errors.New("invalid userinfo in URI")
	ErrMissingHost      = errors.New("missing host in URI")
)

// SchemesWithDNSHost provides a list of schemes for which the host validation
// does not follow RFC3986 (which is quite generic), but assume a valid
// DNS hostname instead.
//
// See: https://www.iana.org/assignments/uri-schemes/uri-schemes.xhtml
//
var SchemesWithDNSHost map[string]bool

func init() {
	SchemesWithDNSHost = map[string]bool{
		"dns":      true,
		"dntp":     true,
		"finger":   true,
		"ftp":      true,
		"git":      true,
		"http":     true,
		"https":    true,
		"imap":     true,
		"irc":      true,
		"jms":      true,
		"mailto":   true,
		"nfs":      true,
		"nntp":     true,
		"ntp":      true,
		"postgres": true,
		"redis":    true,
		"rmi":      true,
		"rtsp":     true,
		"rsync":    true,
		"sftp":     true,
		"skype":    true,
		"smtp":     true,
		"snmp":     true,
		"soap":     true,
		"ssh":      true,
		"steam":    true,
		"svn":      true,
		"tcp":      true,
		"telnet":   true,
		"udp":      true,
		"vnc":      true,
		"wais":     true,
		"ws":       true,
		"wss":      true,
	}
}

// URI represents a general RFC3986 specified URI.
type URI interface {
	// Scheme is the scheme the URI conforms to.
	Scheme() string

	// Authority returns the authority information for the URI, including "//" prefix.
	Authority() Authority

	// Query returns a map of key/value pairs of all parameters
	// in the query string of the URI.
	Query() url.Values

	// Fragment returns the fragment (component proceeded by '#') in the
	// URI if there is one.
	Fragment() string

	// Builder returns a Builder that can be used to modify the URI.
	Builder() Builder

	// String return a string representation of the URI
	String() string

	// Validate the different components of the URI
	Validate() error
}

// Authority represents the authority information that a URI contains
// as specified by RFC3986. Username and password are given by UserInfo().
type Authority interface {
	UserInfo() string
	Host() string
	Port() string
	Path() string
	String() string
	Validate(...string) error
}

// Builder is a construct for building URIs.
type Builder interface {
	URI() URI
	SetScheme(scheme string) Builder
	SetUserInfo(userinfo string) Builder
	SetHost(host string) Builder
	SetPort(port string) Builder
	SetPath(path string) Builder
	SetQuery(query string) Builder
	SetFragment(fragment string) Builder

	// Returns the URI this Builder represents.
	String() string
}

const (
	// string literals
	colonMark       = ":"
	questionMark    = "?"
	fragmentMark    = "#"
	percentMark     = "%"
	atHost          = "@"
	authorityPrefix = "//"
)

var (
	// byte literals
	atBytes       = []byte(atHost)
	colonBytes    = []byte(colonMark)
	queryBytes    = []byte(questionMark)
	fragmentBytes = []byte(fragmentMark)
)

// IsURI tells if a URI is valid according to RFC3986/RFC397
func IsURI(raw string) bool {
	_, err := Parse(raw)
	return err == nil
}

// IsURIReference tells if a URI reference is valid according to RFC3986/RFC397
func IsURIReference(raw string) bool {
	_, err := ParseReference(raw)
	return err == nil
}

// Parse attempts to parse a URI and returns an error if the URI
// is not RFC3986 compliant.
func Parse(raw string) (URI, error) {
	return parse(raw, false)
}

// ParseReference attempts to parse a URI relative reference and returns an error if the URI
// is not RFC3986 compliant.
func ParseReference(raw string) (URI, error) {
	return parse(raw, true)
}

func parse(raw string, withURIReference bool) (URI, error) {
	var (
		schemeEnd   = strings.Index(raw, colonMark)
		hierPartEnd = strings.Index(raw, questionMark)
		queryEnd    = strings.Index(raw, fragmentMark)
		scheme      string

		curr int
	)

	// exclude pathological input
	if schemeEnd == 1 || hierPartEnd == 1 || queryEnd == 1 {
		return nil, ErrInvalidURI
	}

	if hierPartEnd > 0 && hierPartEnd < schemeEnd || queryEnd > 0 && queryEnd < schemeEnd {
		// e.g. htt?p: ; h#ttp: ..
		return nil, ErrInvalidURI
	}

	if queryEnd > 0 && queryEnd < hierPartEnd {
		// e.g.  https://abc#a?b
		hierPartEnd = queryEnd
	}

	isRelative := strings.HasPrefix(raw, authorityPrefix)
	switch {
	case schemeEnd > 0 && !isRelative:
		scheme = raw[curr:schemeEnd]
		if schemeEnd+1 == len(raw) {
			// trailing : (e.g. http:)
			u := &uri{
				scheme: scheme,
			}
			return u, u.Validate()
		}
	case !withURIReference:
		return nil, ErrNoSchemeFound
	case isRelative:
		// start with // and a ':' is following... e.g //example.com:8080/path
		schemeEnd = -1
	}

	curr = schemeEnd + 1

	if hierPartEnd == len(raw)-1 || (hierPartEnd < 0 && queryEnd < 0) {
		// trailing ? or (no query & no fragment)
		if hierPartEnd < 0 {
			hierPartEnd = len(raw)
		}
		authorityInfo, err := parseAuthority(raw[curr:hierPartEnd])
		if err != nil {
			return nil, ErrInvalidURI
		}
		u := &uri{
			scheme:    scheme,
			hierPart:  raw[curr:hierPartEnd],
			authority: authorityInfo,
		}
		return u, u.Validate()
	}

	var (
		hierPart, query, fragment string
		authorityInfo             *authorityInfo
	)
	var err error

	if hierPartEnd > 0 {
		hierPart = raw[curr:hierPartEnd]
		authorityInfo, err = parseAuthority(hierPart)
		if err != nil {
			return nil, ErrInvalidURI
		}
		if hierPartEnd+1 < len(raw) {
			if queryEnd < 0 {
				// query ?, no fragment
				query = raw[hierPartEnd+1:]
			} else if hierPartEnd < queryEnd-1 {
				// query ?, fragment
				query = raw[hierPartEnd+1 : queryEnd]
			}
		}
		curr = hierPartEnd + 1
	}

	if queryEnd == len(raw)-1 && hierPartEnd < 0 {
		// trailing #,  no query "?"
		hierPart = raw[curr:queryEnd]
		authorityInfo, err = parseAuthority(hierPart)
		if err != nil {
			return nil, ErrInvalidURI
		}
		u := &uri{
			scheme:    scheme,
			hierPart:  hierPart,
			authority: authorityInfo,
			query:     query,
		}
		return u, u.Validate()
	}

	if queryEnd > 0 {
		// there is a fragment
		if hierPartEnd < 0 {
			// no query
			hierPart = raw[curr:queryEnd]
			authorityInfo, err = parseAuthority(hierPart)
			if err != nil {
				return nil, ErrInvalidURI
			}
		}
		if queryEnd+1 < len(raw) {
			fragment = raw[queryEnd+1:]
		}
	}

	u := &uri{
		scheme:    scheme,
		hierPart:  hierPart,
		query:     query,
		fragment:  fragment,
		authority: authorityInfo,
	}

	return u, u.Validate()
}

type uri struct {
	// raw components
	scheme   string
	hierPart string
	query    string
	fragment string

	// parsed components
	authority *authorityInfo
}

func (u *uri) URI() URI {
	return u
}

func (u *uri) Scheme() string {
	return u.scheme
}

func (u *uri) Authority() Authority {
	u.ensureAuthorityExists()
	return u.authority
}

// Query returns parsed query parameters like standard lib URL.Query()
func (u *uri) Query() url.Values {
	v, _ := url.ParseQuery(u.query)
	return v
}

func (u *uri) Fragment() string {
	return u.fragment
}

var (
	rexScheme   = regexp.MustCompile(`^[\p{L}][\p{L}\d\+-\.]+$`)
	rexFragment = regexp.MustCompile(`^([\p{L}\d\-\._~\:@!\$\&'\(\)\*\+,;=\?/]|(%[[:xdigit:]]{2})+)+$`)
	rexQuery    = rexFragment
	rexSegment  = regexp.MustCompile(`^([\p{L}\d\-\._~\:@!\$\&'\(\)\*\+,;=]|(%[[:xdigit:]]{2})+)+$`)
	rexHostname = regexp.MustCompile(`^[a-zA-Z0-9\p{L}]((-?[a-zA-Z0-9\p{L}]+)?|(([a-zA-Z0-9-\p{L}]{0,63})(\.)){1,6}([a-zA-Z\p{L}]){2,})$`)

	// unreserved | pct-encoded | sub-delims
	rexRegname = regexp.MustCompile(`^([\p{L}\d\-\._~!\$\&'\(\)\*\+,;=]|(%[[:xdigit:]]{2})+)+$`)
	// unreserved | pct-encoded | sub-delims | ":"
	rexUserInfo = regexp.MustCompile(`^([\p{L}\d\-\._~\:!\$\&'\(\)\*\+,;=\?/]|(%[[:xdigit:]]{2})+)+$`)

	rexIPv6Zone = regexp.MustCompile(`:[^%:]+%25(([\p{L}\d\-\._~\:@!\$\&'\(\)\*\+,;=]|(%[[:xdigit:]]{2}))+)?$`)
	rexPort     = regexp.MustCompile(`^\d+$`)
)

// Validate checks that all parts of a URI abide by allowed characters
func (u *uri) Validate() error {
	if u.scheme != "" {
		if ok := rexScheme.MatchString(u.scheme); !ok {
			return ErrInvalidScheme
		}
	}
	if u.query != "" {
		if ok := rexQuery.MatchString(u.query); !ok {
			return ErrInvalidQuery
		}
	}
	if u.fragment != "" {
		if ok := rexFragment.MatchString(u.fragment); !ok {
			return ErrInvalidFragment
		}
	}
	if u.hierPart != "" {
		if u.authority != nil {
			a := u.Authority()
			if a != nil {
				return a.Validate(u.scheme)
			}
		}
	}
	// empty hierpart case
	return nil
}

type authorityInfo struct {
	prefix   string
	userinfo string
	host     string
	port     string
	path     string
}

func (a authorityInfo) UserInfo() string { return a.userinfo }
func (a authorityInfo) Host() string     { return a.host }
func (a authorityInfo) Port() string     { return a.port }
func (a authorityInfo) Path() string     { return a.path }
func (a authorityInfo) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(a.prefix)
	buf.WriteString(a.userinfo)
	if len(a.userinfo) > 0 {
		buf.Write(atBytes)
	}
	if strings.Index(a.host, colonMark) > 0 {
		// ipv6 address host
		buf.WriteString("[" + a.host + "]")
	} else {
		buf.WriteString(a.host)
	}
	if len(a.port) > 0 {
		buf.Write(colonBytes)
	}
	buf.WriteString(a.port)
	buf.WriteString(a.path)
	return buf.String()
}

func (a authorityInfo) Validate(schemes ...string) error {
	for _, segment := range strings.Split(a.path, "/") {
		if segment == "" {
			continue
		}
		if ok := rexSegment.MatchString(segment); !ok {
			return ErrInvalidPath
		}
	}

	if a.host != "" {
		var isIP bool
		if ok := rexIPv6Zone.MatchString(a.host); ok {
			z := strings.Index(a.host, percentMark)
			isIP = net.ParseIP(a.host[0:z]) != nil
		} else {
			isIP = net.ParseIP(a.host) != nil
		}
		if !isIP {
			var isHost bool
			unescapedHost, err := url.PathUnescape(a.host)
			if err != nil {
				return ErrInvalidHost
			}
			for _, scheme := range schemes {
				if SchemesWithDNSHost[scheme] {
					// DNS name
					isHost = rexHostname.MatchString(unescapedHost)
				} else {
					// standard RFC 3986
					isHost = rexRegname.MatchString(unescapedHost)
				}
				if !isHost {
					return ErrInvalidHost
				}
			}
		}
	}

	if a.port != "" {
		if ok := rexPort.MatchString(a.port); !ok {
			return ErrInvalidPort
		}
		if a.host == "" {
			return ErrMissingHost
		}
	}

	if a.userinfo != "" {
		if ok := rexUserInfo.MatchString(a.userinfo); !ok {
			return ErrInvalidUserInfo
		}
	}

	return nil
}

func parseAuthority(hier string) (*authorityInfo, error) {
	// as per RFC 3986 Section 3.6
	var prefix, userinfo, host, port, path string

	// authority sections MUST begin with a '//'
	if strings.HasPrefix(hier, authorityPrefix) {
		prefix = authorityPrefix
		hier = strings.TrimPrefix(hier, authorityPrefix)
	}

	if prefix == "" {
		path = hier
	} else {
		// authority   = [ userinfo "@" ] host [ ":" port ]
		slashEnd := strings.Index(hier, "/")
		if slashEnd > 0 {
			if slashEnd < len(hier) {
				path = hier[slashEnd:]
			}
			hier = hier[:slashEnd]
		}

		host = hier
		if at := strings.Index(host, atHost); at > 0 {
			userinfo = host[:at]
			if at+1 < len(host) {
				host = host[at+1:]
			}
		}

		if bracket := strings.Index(host, "["); bracket >= 0 {
			// ipv6 addresses: "[" xx:yy:zz "]":port
			rawHost := host
			closingbracket := strings.Index(host, "]")
			if closingbracket > 0 {
				host = host[bracket+1 : closingbracket-bracket]
				rawHost = rawHost[closingbracket+1:]
			} else {
				return nil, ErrInvalidURI
			}
			if colon := strings.Index(rawHost, colonMark); colon >= 0 {
				if colon+1 < len(rawHost) {
					port = rawHost[colon+1:]
				}
			}
		} else {
			if colon := strings.Index(host, colonMark); colon >= 0 {
				if colon+1 < len(host) {
					port = host[colon+1:]
				}
				host = host[:colon]
			}
		}
	}

	return &authorityInfo{
		prefix:   prefix,
		userinfo: userinfo,
		host:     host,
		port:     port,
		path:     path,
	}, nil
}

func (u *uri) ensureAuthorityExists() {
	if u.authority == nil {
		u.authority = &authorityInfo{}
	} else {
		if u.authority.userinfo != "" ||
			u.authority.host != "" ||
			u.authority.port != "" {
			u.authority.prefix = "//"
		}
	}
}

func (u *uri) SetScheme(scheme string) Builder {
	u.scheme = scheme
	return u
}

func (u *uri) SetUserInfo(userinfo string) Builder {
	u.ensureAuthorityExists()
	u.authority.userinfo = userinfo
	return u
}

func (u *uri) SetHost(host string) Builder {
	u.ensureAuthorityExists()
	u.authority.host = host
	return u
}

func (u *uri) SetPort(port string) Builder {
	u.ensureAuthorityExists()
	u.authority.port = port
	return u
}

func (u *uri) SetPath(path string) Builder {
	u.ensureAuthorityExists()
	u.authority.path = path
	return u
}

func (u *uri) SetQuery(query string) Builder {
	u.query = query
	return u
}

func (u *uri) SetFragment(fragment string) Builder {
	u.fragment = fragment
	return u
}

func (u *uri) Builder() Builder {
	return u
}

func (u *uri) String() string {
	buf := bytes.NewBuffer(nil)
	if len(u.scheme) > 0 {
		buf.WriteString(u.scheme)
		buf.Write(colonBytes)
	}

	buf.WriteString(u.authority.String())

	if len(u.query) > 0 {
		buf.Write(queryBytes)
		buf.WriteString(u.query)
	}

	if len(u.fragment) > 0 {
		buf.Write(fragmentBytes)
		buf.WriteString(u.fragment)
	}

	return buf.String()
}
