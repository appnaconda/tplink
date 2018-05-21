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
	INFO = `{"system":{"get_sysinfo":{}}}`
	ON   = `{"system":{"set_relay_state":{"state":1}}}`
	OFF  = `{"system":{"set_relay_state":{"state":0}}}`
	TIME = `{"time":{"get_time":{}}}`
)

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

func send(ip string, payload []byte) (data []byte, err error) {
	// 10 second timeout
	conn, err := net.DialTimeout("tcp", ip+":9999", time.Duration(10)*time.Second)
	if err != nil {
		fmt.Println("Cannot connnect to plug:", err)
		data = nil
		return
	}
	_, err = conn.Write(payload)
	data, err = ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println("Cannot read data from plug:", err)
	}
	return

}

func NewHS110(ip string) *HS110 {
	return &HS110{ip: ip}
}
