package host

import (
	"github.com/qiniu/qshell/v2/iqshell/common/data"
	"github.com/qiniu/qshell/v2/iqshell/common/provider"
)

type Provider interface {
	Available() (available bool, err *data.CodeError)
	Provide() (host *Host, err *data.CodeError)
	Freeze(host *Host)
}

func NewListProviderWithHostStrings(hostStrings []string) Provider {
	hosts := make([]*Host, 0, len(hostStrings))
	for _, h := range hostStrings {
		hosts = append(hosts, &Host{
			Host:   h,
			Domain: "",
		})
	}
	return NewListProvider(hosts)
}

func NewListProvider(hosts []*Host) Provider {
	items := make([]provider.Item, 0, len(hosts))
	for _, h := range hosts {
		items = append(items, h)
	}
	return &listProvider{
		p: provider.NewListProvider(items),
	}
}

type listProvider struct {
	p provider.Provider
}

func (l *listProvider) Available() (available bool, err *data.CodeError) {
	return l.p.Available()
}

func (l *listProvider) Provide() (host *Host, err *data.CodeError) {
	i, e := l.p.Provide()
	host, _ = i.(*Host)
	return host, e
}

func (l *listProvider) Freeze(host *Host) {
	l.p.Freeze(host)
}
