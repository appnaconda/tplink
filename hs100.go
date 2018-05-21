package tplink

import (
	"fmt"
)

type HS100 struct {
	ip string
}

func (p *HS100) TurnOn() error {
	data := encrypt(ON)
	_, err := send(p.ip, data)
	return err
}

func (p *HS100) TurnOff() error {
	data := encrypt(OFF)
	_, err := send(p.ip, data)
	return err
}

func (p *HS100) Info() (string, error) {
	data := encrypt(INFO)
	reading, err := send(p.ip, data)
	if err != nil {
		return "", err
	}

	return decrypt(reading[4:]), nil
}

func (p *HS100) Time() (string, error) {
	data := encrypt(TIME)
	reading, err := send(p.ip, data)
	if err != nil {
		return "", err
	}

	return decrypt(reading[4:]), nil
}

func (p *HS100) DailyStats(month int, year int) (string, error) {
	json := fmt.Sprintf(`{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`, month, year)
	data := encrypt(json)
	reading, err := send(p.ip, data)
	if err != nil {
		return "", err
	}

	return decrypt(reading[4:]), nil
}
