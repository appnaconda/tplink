package tplink

import (
	"encoding/json"
	"fmt"
	"time"
)

type HS100 struct {
	ip string
}

// Get System Info (Software & Hardware Versions, MAC, deviceID, hwID etc.)
func (p *HS100) Info() (*Info, error) {
	data, err := exec(p.ip, GET_INFO)
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
	return exec(p.ip, REBOOT)
}

// Reset
func (p *HS100) Reset() (string, error) {
	return exec(p.ip, RESET)
}

// Set alias/name
func (p *HS100) SetAlias(alias string) error {
	data, err := exec(p.ip, fmt.Sprintf(SET_ALIAS, alias))
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
	data, err := exec(p.ip, TURN_ON)
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
	data, err := exec(p.ip, TURN_OFF)
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
	data, err := exec(p.ip, TURN_LED_ON)
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
	data, err := exec(p.ip, TURN_LED_OFF)
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

// TODO: return a timezone instead of index
func (p *HS100) TimeZone() (int, error) {
	data, err := exec(p.ip, GET_TIMEZONE)
	if err != nil {
		return 0, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return 0, err
	}

	tz := r.Time.GetTimeZone

	if tz.ErrorCode != 0 {
		return 0, fmt.Errorf("failed to get time zone. Error code=%d", r.System.SetState.ErrorCode)
	}

	return tz.Index, nil
}

func (p *HS100) Time() (time.Time, error) {
	data, err := exec(p.ip, GET_TIME)
	if err != nil {
		return time.Time{}, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return time.Time{}, err
	}

	t := r.Time.GetTime

	if t.ErrorCode != 0 {
		return time.Time{}, fmt.Errorf("failed to get time. Error code=%d", r.System.SetState.ErrorCode)
	}

	// TODO: get timezone
	//timezone, err := p.TimeZone()
	// if err != nill {
	//	return "", fmt.Errorf("failed to get device timezone: %s", err)
	//}

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get device timezone: %s", err)
	}
	d := time.Date(t.Year, time.Month(t.Month), t.Day, t.Hour, t.Minutes, t.Seconds, 0, loc)

	return d, nil
}

func (p *HS100) SetTimeZone(t time.Time) error {
	// TODO: timezone
	cmd := fmt.Sprintf(SET_TIMEZONE, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 18)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.Time.SetTimeZone.ErrorCode != 0 {
		return fmt.Errorf("failed to set timezone. Error code=%d, msg: %s", r.Time.SetTimeZone.ErrorCode, r.Time.SetTimeZone.ErrorMessage)
	}

	return nil
}

func (p *HS100) ScanWifi() ([]AP, error) {
	data, err := exec(p.ip, SCAN_WIFI)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	if r.NetIf.GetScanInfo.ErrorCode != 0 {
		return nil, fmt.Errorf("failed to scan for wifi networks. Error code=%d, msg: %s", r.NetIf.GetScanInfo.ErrorCode, r.NetIf.GetScanInfo.ErrorMessage)
	}

	return r.NetIf.GetScanInfo.List, nil
}

func (p *HS100) SetWifi(ssid string, password string, keyType int) error {
	cmd := fmt.Sprintf(SET_WIFI, ssid, password, keyType)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.NetIf.SetWifi.ErrorCode != 0 {
		return fmt.Errorf("failed to set wifi. Error code=%d, msg: %s", r.NetIf.SetWifi.ErrorCode, r.NetIf.SetWifi.ErrorMessage)
	}
	fmt.Println(data)

	return nil
}
