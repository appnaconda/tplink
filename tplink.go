// https://www.softscheck.com/en/reverse-engineering-tp-link-hs110/
package tplink

import (
	"bufio"
	"encoding/json"
	"net"
	"strconv"
	"time"
)

const (
	defaultPort = 9999
	connTimeout = 1 * time.Second
)

// https://github.com/softScheck/tplink-smartplug/blob/master/tplink-smarthome-commands.txt
const (
	// Plug HS100 and HS110
	GET_INFO     = `{"system":{"get_sysinfo":{}}}`
	REBOOT       = `{"system":{"reboot":{"delay":1}}}`
	RESET        = `{"system":{"reset":{"delay":1}}}`
	SET_ALIAS    = `{"system":{"set_dev_alias":{"alias":"%s"}}}`
	TURN_LED_ON  = `{"system":{"set_led_off":{"off":0}}}`
	TURN_LED_OFF = `{"system":{"set_led_off":{"off":1}}}`
	TURN_ON      = `{"system":{"set_relay_state":{"state":1}}}`
	TURN_OFF     = `{"system":{"set_relay_state":{"state":0}}}`
	GET_TIME     = `{"time":{"get_time":{}}}`
	GET_TIMEZONE = `{"time":{"get_timezone":null}}`
	SET_TIMEZONE = `{"time":{"set_timezone":{"year":%d,"month":%d,"mday":%d,"hour":%d,"min":%d,"sec":%d,"index":%d}}}`

	// HS110
	GET_METER         = `{"system":{"get_sysinfo":{}}, "emeter":{"get_realtime":{},"get_vgain_igain":{}}}`
	GET_DAILY_STATS   = `{"emeter":{"get_daystat":{"month":%d,"year":%d}}}`
	GET_MONTHLY_STATS = `{"emeter":{""get_monthstat":{"year":2016}}}`
	ERASE_ALL_STATS   = `{"emeter":{"erase_emeter_stat":null}}`

	// Schedule
	GET_SCHEDULE_RULES_LIST = `{"schedule":{"get_rules":null}}`

	// WLAN Commands
	SCAN_WIFI = `{"netif":{"get_scaninfo":{"refresh":1}}}`
	SET_WIFI  = `{"netif":{"set_stainfo":{"ssid":"%s","password":"%s","key_type":%d}}}`
)

type Device struct {
	IPAddress string
	Info      Info
}

type Response struct {
	System struct {
		*Info    `json:"get_sysinfo"`
		SetAlias struct {
			ErrorCode    int    `json:"err_code"`
			ErrorMessage string `json:"err_msg"`
		} `json:"set_dev_alias"`
		SetState struct {
			ErrorCode    int    `json:"err_code"`
			ErrorMessage string `json:"err_msg"`
		} `json:"set_relay_state"`
	}

	EMeter struct {
		*Meter     `json:"get_realtime"`
		DailyStats *DailyStats `json:"get_daystat"`
	}

	Time struct {
		GetTime struct {
			Year         int    `json:"year"`
			Month        int    `json:"month"`
			Day          int    `json:"mday"`
			Hour         int    `json:"hour"`
			Minutes      int    `json:"min"`
			Seconds      int    `json:"sec"`
			ErrorCode    int    `json:"err_code"`
			ErrorMessage string `json:"err_msg"`
		} `json:"get_time"`

		GetTimeZone struct {
			Index        int    `json:"index"`
			ErrorCode    int    `json:"err_code"`
			ErrorMessage string `json:"err_msg"`
		} `json:"get_timezone"`

		SetTimeZone struct {
			ErrorCode    int    `json:"err_code"`
			ErrorMessage string `json:"err_msg"`
		} `json:"set_timezone"`
	}

	NetIf struct {
		GetScanInfo struct {
			List         []AP   `json:"ap_list"`
			ErrorCode    int    `json:"err_code"`
			ErrorMessage string `json:"err_msg"`
		} `json:"get_scaninfo"`

		SetWifi struct {
			ErrorCode    int    `json:"err_code"`
			ErrorMessage string `json:"err_msg"`
		} `json:"set_stainfo"`
	} `json:"netif"`
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

type AP struct {
	SSID    string `json:"ssid"`
	KeyType int    `json:"key_type"`
}

func decrypt(request []byte) string {
	result := make([]byte, len(request))
	key := byte(0xAB)
	for i, c := range request {
		var a = key ^ uint8(c)
		key = uint8(c)
		result[i] = a
	}
	return string(result)
}

func encrypt(s string) []byte {
	request := []byte(s)
	key := byte(0xAB)
	result := make([]byte, 4+len(request))
	result[0] = 0x0
	result[1] = 0x0
	result[2] = 0x0
	result[3] = 0x0
	for i, c := range request {
		var a = key ^ uint8(c)
		key = uint8(a)
		result[i+4] = a
	}
	return result[4:]
}

func exec(ip string, cmd string) (string, error) {
	data := encrypt(cmd)
	port := 9999
	conn, err := net.Dial("udp4", ip+":"+strconv.Itoa(port))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		return "", err
	}
	rData := make([]byte, 1500)
	rLen, err := bufio.NewReader(conn).Read(rData)
	if err != nil {
		return "", err
	}

	return decrypt(rData[:rLen]), nil
}

func Scan(timeout time.Duration) ([]Device, error) {
	devices := []Device{}

	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:9999")
	if err != nil {
		return nil, err
	}

	fromAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8755")
	if err != nil {
		return nil, err
	}

	sock, err := net.ListenUDP("udp", fromAddr)
	defer sock.Close()
	if err != nil {
		return nil, err
	}
	sock.SetReadBuffer(2048)

	cmd := encrypt(GET_INFO)
	_, err = sock.WriteToUDP(cmd, broadcastAddr)
	if err != nil {
		return nil, err
	}

	err = sock.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}

	for {
		buff := make([]byte, 2048)
		rlen, addr, err := sock.ReadFromUDP(buff)
		if err != nil {
			break
		}

		data := decrypt(buff[:rlen])

		r := Response{}
		if err := json.Unmarshal([]byte(data), &r); err != nil {
			return nil, err
		}

		devices = append(devices, Device{
			IPAddress: addr.IP.String(),
			Info:      *r.System.Info,
		})
	}

	return devices, err
}

func NewHS110(ip string) *HS110 {
	return &HS110{HS100{ip: ip}}
}

func NewHS100(ip string) *HS100 {
	return &HS100{ip: ip}
}
