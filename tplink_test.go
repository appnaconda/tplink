package tplink

import "testing"

func TestEncription(t *testing.T) {
	tt := []string{GET_INFO, GET_METER, GET_DAILY_STATS, GET_SCHEDULE_RULES_LIST}

	for _, v := range tt {
		e := encrypt(v)
		d := decrypt(e)
		if v != d {
			t.Errorf("expecting %s; got %s", v, d)
		}
	}
}

func TestDaysToString(t *testing.T) {
	tt := []struct {
		days      Days
		expecting string
	}{
		{Days{}, "[0,0,0,0,0,0,0]"},
		{Days{true, true, true, true, true, true, true}, "[1,1,1,1,1,1,1]"},
		{Days{true, false, true, false, true, false, true}, "[1,0,1,0,1,0,1]"},
	}

	for _, v := range tt {
		s := v.days.String()

		if s != v.expecting {
			t.Errorf("expecting %s; got %s", v.expecting, s)
		}
	}

}
