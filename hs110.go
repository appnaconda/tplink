package tplink

import (
	"encoding/json"
	"fmt"
)

type HS110 struct {
	HS100
}

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
