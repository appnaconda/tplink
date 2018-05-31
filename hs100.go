package tplink

import (
	"encoding/json"
	"fmt"
)

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

// Set alias/name
func (p *HS100) SetAlias(alias string) error {
	data, err := p.exec(fmt.Sprintf(SET_ALIAS, alias))
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetAlias.ErrorCode != 0 {
		return fmt.Errorf("failed to set alias. Error code=%d", r.System.SetAlias.ErrorCode)
	}
	return nil
}

// Turn On
func (p *HS100) TurnOn() error {
	data, err := p.exec(TURN_ON)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
}

// Turn Off
func (p *HS100) TurnOff() error {
	data, err := p.exec(TURN_OFF)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
}

// Turn Led Light On
func (p *HS100) TurnLedOn() error {
	data, err := p.exec(TURN_LED_ON)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
}

// Turn Led Light Off
func (p *HS100) TurnLedOff() error {
	data, err := p.exec(TURN_LED_OFF)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
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
