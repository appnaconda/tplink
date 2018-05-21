package tplink

type HS110 struct {
	HS100
	ip string
}

func (p *HS110) MeterInfo() (string, error) {
	json := `{"system":{"get_sysinfo":{}}, "emeter":{"get_realtime":{},"get_vgain_igain":{}}}`
	data := encrypt(json)
	reading, err := send(p.ip, data)
	if err != nil {
		return "", err
	}

	return decrypt(reading[4:]), nil
}
