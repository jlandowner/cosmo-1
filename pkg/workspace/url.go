package workspace

import (
	"net/url"
	"strings"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
)

const (
	DefaultHostBase     = "{{NETRULE}}-{{WORKSPACE}}-{{USER}}"
	URLVarNetRule       = "{{NETRULE}}"
	URLVarWorkspaceName = "{{WORKSPACE}}"
	URLVarUserName      = "{{USER}}"
)

type URLBase struct {
	url.URL
}

func NewURLBase(protocol, host string) URLBase {
	return URLBase{
		url.URL{
			Scheme: protocol,
			Host:   host,
		},
	}
}

func (u URLBase) GenHost(netRule, wsName, userName string) string {
	host := u.Host

	host = strings.ReplaceAll(host, URLVarNetRule, netRule)
	host = strings.ReplaceAll(host, URLVarWorkspaceName, wsName)
	host = strings.ReplaceAll(host, URLVarUserName, userName)

	return host
}

func (u URLBase) GenURL(r cosmov1alpha1.NetworkRule) string {
	u.Scheme = r.Protocol
	u.Host = r.Host
	u.Path = r.HTTPPath
	return u.String()
}
