// https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/
package tplink

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

const (
	defaultPort = 9999
	connTimeout = 10 * time.Second
)

// TODO: check for panic when a bad or misspell command is passed

// https://github.com/softScheck/tplink-smartplug/blob/master/tplink-smarthome-commands.txt
const (
	// Plug HS100 and HS110
	GET_INFO     = `{"system":{"get_sysinfo":{}}}`
	REBOOT       = `{"system":{"reboot":{"delay":1}}}`
	RESET        = `{"system":{"reset":{"delay":1}}}`
	SET_ALIAS    = `{"system":{"set_dev_alias":{"alias":"%s"}}}`
	TURN_LED_ON  = `{"system":{"set_led_off":{"off":1}}}`
	TURN_LED_OFF = `{"system":{"set_led_off":{"off":0}}}` // need to be tested
	SET_LOCATION = `{"system":{"set_dev_location":{"longitude":6.9582814,"latitude":50.9412784}}}`
	GET_ICON     = `{"system":{"get_dev_icon":null}}`
	SET_ICON     = `{"system":{"set_dev_icon":{"icon":"xxxx","hash":"ABCD"}}}`
	TURN_ON      = `{"system":{"set_relay_state":{"state":1}}}`
	TURN_OFF     = `{"system":{"set_relay_state":{"state":0}}}`
	GET_TIME     = `{"time":{"get_time":{}}}`
	GET_TIMEZONE = `{"time":{"get_timezone":null}}`
	SET_TIMEZONE = `{"time":{"set_timezone":{"year":2016,"month":1,"mday":1,"hour":10,"min":10,"sec":10,"index":42}}}`
	// HS110
	GET_METER         = `{"system":{"get_sysinfo":{}}, "emeter":{"get_realtime":{},"get_vgain_igain":{}}}`
	GET_DAILY_STATS   = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`
	GET_MONTHLY_STATS = `{"emeter":{""get_monthstat":{"year":2016}}}`
	ERASE_ALL_STATS   = `{"emeter":{"erase_emeter_stat":null}}`

	// Schedule
	GET_SCHEDULE_RULES_LIST = `{"schedule":{"get_rules":null}}`
)

type Response struct {
	System struct {
		*Info    `json:"get_sysinfo"`
		SetAlias struct {
			ErrorCode int `json:"err_code"`
		} `json:"set_dev_alias"`
		SetState struct {
			ErrorCode int `json:"err_code"`
		} `json:"set_relay_state"`
	}

	EMeter struct {
		*Meter     `json:"get_realtime"`
		DailyStats *DailyStats `json:"get_daystat"`
	}
}

type Info struct {
	SoftwareVersion string  `json:"sw_ver"`      // Software version
	HardwareVersion string  `json:"hw_ver"`      // Hardware version
	HardwareID      string  `json:"hwId"`        // Hardware ID
	Type            string  `json:"type"`        // Type
	Model           string  `json:"model"`       // Model
	MacAddr         string  `json:"mac"`         // Mac Address
	DeviceID        string  `json:"deviceId"`    // Device ID
	FirmwareID      string  `json:"fwId"`        // Firmware ID
	OEMID           string  `json:"oemId"`       // OEM ID
	Alias           string  `json:"alias"`       // Description. e.g "Basement light"
	IconHash        string  `json:"icon_hash"`   // hash for custom picture
	State           int     `json:"relay_state"` // State:  0 = OFF; 1 = ON
	ActiveMode      string  `json:"active_mode"` // "schedule" for schedule mode
	Feature         string  `json:"feature"`     // "TIM:ENE" (Timer, Energy Monitor)
	Updating        int     `json:"updating"`    // 0 = not updating
	RSSI            int     `json:"rssi"`        // Signal Strength Indicator in dBm (e.g. -35)
	LedOff          int     `json:"led_off"`     // 0 = Led ON (default); 1 = Led OFF
	Latitude        float64 `json:"latitude"`    // Optional Geolocation information
	Longitude       float64 `json:"longitude"`   // Optional Geolocation information
}

func (i Info) IsOn() bool {
	return i.State == 1
}

func (i Info) IsLedOn() bool {
	return i.LedOff == 0
}

type Meter struct {
	Current float64 `json:"current"`
	Voltage float64 `json:"voltage"`
	Power   float64 `json:"power"`
	Total   float64 `json:"total"`
}

type DailyStats struct {
	DailyUsageList []*DailyUsage `json:"day_list"`
}

type DailyUsage struct {
	Year   int
	Month  int
	Day    int
	Energy float64
}

// encript message
func encrypt(s string) []byte {
	n := len(s)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint32(n))
	ciphertext := []byte(buf.Bytes())

	key := byte(0xAB)
	payload := make([]byte, n)
	for i := 0; i < n; i++ {
		payload[i] = s[i] ^ key
		key = payload[i]
	}

	for i := 0; i < len(payload); i++ {
		ciphertext = append(ciphertext, payload[i])
	}

	return ciphertext
}

func decrypt(ciphertext []byte) string {
	n := len(ciphertext)
	key := byte(0xAB)
	var nextKey byte
	for i := 0; i < n; i++ {
		nextKey = ciphertext[i]
		ciphertext[i] = ciphertext[i] ^ key
		key = nextKey
	}
	return string(ciphertext)
}

func exec(ip string, payload []byte) ([]byte, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, defaultPort), connTimeout)
	if err != nil {
		return nil, fmt.Errorf("cannot connnect to plug: %s", err)
	}
	defer conn.Close()

	_, err = conn.Write(payload)
	data, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("cannot read data from plug: %s", err)
	}
	return data, nil

}

func NewHS110(ip string) *HS110 {
	return &HS110{HS100{ip: ip}}
}

func NewHS100(ip string) *HS100 {
	return &HS100{ip: ip}
}
