package tplink

import (
	"encoding/json"
	"fmt"
)

// TP-Link HS110 smart plug
type HS110 struct {
	HS100
}

// Gets Realtime Current and Voltage Reading
func (p *HS110) Meter() (*Meter, error) {
	data, err := exec(p.ip, GET_METER)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return r.EMeter.Meter, nil
}

// Gets Daily Statistic for given Month
func (p *HS110) DailyStats(month int, year int) ([]*DailyUsage, error) {
	data, err := exec(p.ip, fmt.Sprintf(GET_DAILY_STATS, month, year))
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return r.EMeter.DailyStats.DailyUsageList, nil
}

// Get Montly Statistic for given Year
func (p *HS110) MonthlyStats(year int) ([]*MonthlyUsage, error) {
	data, err := exec(p.ip, fmt.Sprintf(GET_MONTHLY_STATS, year))
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return r.EMeter.MonthlyStats.MonthlyUsageList, nil
}

// Erase All EMeter Statistics
func (p *HS110) EraseAllStats() error {
	data, err := exec(p.ip, ERASE_ALL_STATS)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.EMeter.EraseMeterStat.ErrorCode != 0 {
		return fmt.Errorf("failed to erase meter stats. Error code=%d, msg: %s", r.EMeter.EraseMeterStat.ErrorCode, r.EMeter.EraseMeterStat.ErrorMessage)
	}

	return nil
}
