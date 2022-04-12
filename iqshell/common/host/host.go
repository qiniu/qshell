package host

import "github.com/qiniu/qshell/v2/iqshell/common/provider"

type Host struct {
	Host   string
	Domain string // 可为 host，也可为 IP + 端口
}

var _ provider.Item = (*Host)(nil)

func (h *Host) Equal(item provider.Item) bool {
	host, _ := item.(*Host)
	if h == nil || host == nil {
		return false
	}
	return h.Host == host.Host && h.Domain == host.Domain
}

func (h *Host) GetServer() string {
	if len(h.Domain) > 0 {
		return h.Domain
	}
	return h.Host
}

func (h *Host) GetHost() string {
	return h.Host
}
