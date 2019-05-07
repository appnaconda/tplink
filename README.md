# tplink

A Go library for controlling tp-link devices, based on this information: [https://www.softscheck.com/en/reverse-engineering-tp-link-hs110](https://www.softscheck.com/en/reverse-engineering-tp-link-hs110)

# Supported Devices

* HS100
* HS110 
* HS105

# Supported Features

Most of these features were implemented: [https://github.com/softScheck/tplink-smartplug/blob/master/tplink-smarthome-commands.txt](https://github.com/softScheck/tplink-smartplug/blob/master/tplink-smarthome-commands.txt)


# Usage

### Scan

Scan your network:
```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/appnaconda/tplink"
)

func main() {
	timeout := 2 * time.Second
	devices, err := tplink.Scan(timeout)
	if err != nil {
		log.Fatalf("failed: %s\n", err)
	}

	for _, v := range devices {
		fmt.Printf("%#v\n", v)
	}
}
```

Response:
```
tplink.Device{IPAddress:"10.0.1.XX1", Info:tplink.Info{SoftwareVersion:"1.2.5 Build 171206 Rel.085954", HardwareVersion:"1.0", HardwareID:"60FF6B258734EA6880E186F8C96DDC61", Type:"IOT.SMARTPLUGSWITCH", Model:"HS110(US)", MacAddr:"XX:XX:XX:XX:XX:XX", DeviceID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", FirmwareID:"00000000000000000000000000000000", OEMID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", Alias:"Plug1", IconHash:"", State:1, ActiveMode:"none", Feature:"TIM:ENE", Updating:0, RSSI:-38, LedOff:0, Latitude:0, Longitude:0}}
tplink.Device{IPAddress:"10.0.1.XX2", Info:tplink.Info{SoftwareVersion:"1.2.5 Build 171206 Rel.085954", HardwareVersion:"1.0", HardwareID:"60FF6B258734EA6880E186F8C96DDC61", Type:"IOT.SMARTPLUGSWITCH", Model:"HS110(US)", MacAddr:"XX:XX:XX:XX:XX:XX", DeviceID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", FirmwareID:"00000000000000000000000000000000", OEMID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", Alias:"Plug2", IconHash:"", State:1, ActiveMode:"schedule", Feature:"TIM:ENE", Updating:0, RSSI:-55, LedOff:0, Latitude:0, Longitude:0}}
tplink.Device{IPAddress:"10.0.1.XX3", Info:tplink.Info{SoftwareVersion:"1.5.1 Build 171109 Rel.165709", HardwareVersion:"2.0", HardwareID:"0DC28CDD0B7E6C55F52AD35B8B68277E", Type:"IOT.SMARTPLUGSWITCH", Model:"HS100(US)", MacAddr:"XX:XX:XX:XX:XX:XX", DeviceID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", FirmwareID:"00000000000000000000000000000000", OEMID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", Alias:"Plug3", IconHash:"", State:1, ActiveMode:"none", Feature:"TIM", Updating:0, RSSI:-58, LedOff:0, Latitude:0, Longitude:0}}

```

### Info

Get device info:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/appnaconda/tplink"
)

func main() {
	ip := "10.0.1.XXX" // Your device IP
	plug := tplink.NewHS100(ip, 2 * time.Second)

	info, err := plug.Info()
	if err != nil {
		log.Fatalf("failed: %s\n", err)
	}

	fmt.Printf("Result: %+v\n", info)
}
```

Response:
```
&tplink.Info{SoftwareVersion:"1.5.1 Build 171109 Rel.165709", HardwareVersion:"2.0", HardwareID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", Type:"IOT.SMARTPLUGSWITCH", Model:"HS100(US)", MacAddr:"XX:XX:XX:XX:XX:XX", DeviceID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", FirmwareID:"00000000000000000000000000000000", OEMID:"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", Alias:"Plug3", IconHash:"", State:1, ActiveMode:"none", Feature:"TIM", Updating:0, RSSI:-59, LedOff:0, Latitude:0, Longitude:0}
```