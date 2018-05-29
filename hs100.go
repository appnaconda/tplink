package tplink

import "encoding/json"

type HS100 struct {
	ip string
}

func (p *HS100) exec(cmd string) (string, error) {
	data := encrypt(cmd)
	reading, err := exec(p.ip, data)
	if err != nil {
		return "", err
	}

	return decrypt(reading[4:]), nil
}

// Get System Info (Software & Hardware Versions, MAC, deviceID, hwID etc.)
func (p *HS100) Info() (*Info, error) {
	data, err := p.exec(GET_INFO)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return r.System.Info, nil

}

// Reboot
func (p *HS100) Reboot() (string, error) {
	return p.exec(REBOOT)
}

// Reset
func (p *HS100) Reset() (string, error) {
	return p.exec(RESET)
}

// Turn On
func (p *HS100) TurnOn() error {
	_, err := p.exec(TURN_ON)
	return err
}

// Turn Off
func (p *HS100) TurnOff() error {
	_, err := p.exec(TURN_OFF)
	return err
}

func (p *HS100) Time() (string, error) {
	data := encrypt(GET_TIME)
	reading, err := exec(p.ip, data)
	if err != nil {
		return "", err
	}

	return decrypt(reading[4:]), nil
}

func (p *HS110) ScheduleRules() (string, error) {
	data := encrypt(GET_SCHEDULE_RULES_LIST)
	reading, err := exec(p.ip, data)
	if err != nil {
		return "", err
	}

	return decrypt(reading[4:]), nil
}
