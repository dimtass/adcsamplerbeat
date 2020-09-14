// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period 			time.Duration `config:"period"`
	SerialPort		string	`config:"serial_port"`
	SerialBaud		int		`config:"serial_baud"`
	SerialTimeout	time.Duration	`config:"serial_timeout"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	SerialPort: "/dev/ttyUSB0",
	SerialBaud: 115200,
	SerialTimeout: 50,
}
