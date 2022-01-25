package config

type Hosts struct {
	UC  []string `json:"uc,omitempty"`
	Api []string `json:"api,omitempty"`
	Rs  []string `json:"rs,omitempty"`
	Rsf []string `json:"rsf,omitempty"`
	Io  []string `json:"io,omitempty"`
	Up  []string `json:"up,omitempty"`
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

func getOneHostFromStringArray(hosts []string) string {
	if len(hosts) > 0 {
		return hosts[0]
	} else {
		return ""
	}
}

func (h *Hosts) merge(from *Hosts) {
	if from == nil {
		return
	}

	if len(h.UC) == 0 {
		h.UC = from.UC
	}

	if len(h.Api) == 0 {
		h.Api = from.Api
	}

	if len(h.Rsf) == 0 {
		h.Rsf = from.Rsf
	}

	if len(h.Rs) == 0 {
		h.Rs = from.Rs
	}

	if len(h.Io) == 0 {
		h.Io = from.Io
	}

	if len(h.Up) == 0 {
		h.Up = from.Up
	}
}
