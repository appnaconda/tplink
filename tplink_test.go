package tplink

import "testing"

func TestEncription(t *testing.T) {
	tt := []string{GET_INFO, GET_METER, GET_DAILY_STATS, GET_SCHEDULE_RULES_LIST}

	for _, v := range tt {
		e := encrypt(v)
		d := decrypt(e[4:])
		if v != d {
			t.Errorf("expecting %s; got %s", v, d)
		}
	}
}
