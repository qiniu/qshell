package config

import "github.com/qiniu/qshell/v2/iqshell/common/utils"

type Hosts struct {
	UC       []string `json:"uc,omitempty"`
	Api      []string `json:"api,omitempty"`
	Rs       []string `json:"rs,omitempty"`
	Rsf      []string `json:"rsf,omitempty"`
	Io       []string `json:"io,omitempty"`
	Up       []string `json:"up,omitempty"`
	Endpoint []string `json:"endpoint,omitempty"`
}

func (h *Hosts) GetOneUc() string {
	return getOneHostFromStringArray(h.UC)
}

func (h *Hosts) GetOneApi() string {
	return getOneHostFromStringArray(h.Api)
}

func (h *Hosts) GetOneRs() string {
	return getOneHostFromStringArray(h.Rs)
}

func (h *Hosts) GetOneRsf() string {
	return getOneHostFromStringArray(h.Rsf)
}

func (h *Hosts) GetOneIo() string {
	return getOneHostFromStringArray(h.Io)
}

func (h *Hosts) GetOneUp() string {
	return getOneHostFromStringArray(h.Up)
}

func (h *Hosts) GetOneEndpoint() string {
	return getOneHostFromStringArray(h.Endpoint)
}

func getOneHostFromStringArray(hosts []string) string {
	hosts = getRealHosts(hosts)
	if len(hosts) > 0 {
		return hosts[0]
	} else {
		return ""
	}
}

func getRealHosts(hosts []string) []string {
	if hosts == nil {
		return nil
	}

	newHosts := make([]string, 0, len(hosts))
	for _, host := range hosts {
		newHosts = append(newHosts, utils.RemoveUrlScheme(host))
	}
	return newHosts
}

func (h *Hosts) merge(from *Hosts) {
	if from == nil {
		return
	}

	if len(h.UC) == 0 {
		h.UC = getRealHosts(from.UC)
	}

	if len(h.Api) == 0 {
		h.Api = getRealHosts(from.Api)
	}

	if len(h.Rsf) == 0 {
		h.Rsf = getRealHosts(from.Rsf)
	}

	if len(h.Rs) == 0 {
		h.Rs = getRealHosts(from.Rs)
	}

	if len(h.Io) == 0 {
		h.Io = getRealHosts(from.Io)
	}

	if len(h.Up) == 0 {
		h.Up = getRealHosts(from.Up)
	}

	if len(h.Endpoint) == 0 {
		h.Endpoint = getRealHosts(from.Endpoint)
	}
}
