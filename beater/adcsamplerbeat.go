package beater

import (
	"fmt"
	"time"
	"strings"
	"strconv"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/dimtass/adcsamplerbeat/config"
	"github.com/tarm/serial"
)

// adcsamplerbeat configuration.
type adcsamplerbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of adcsamplerbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &adcsamplerbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts adcsamplerbeat.
func (bt *adcsamplerbeat) Run(b *beat.Beat) error {
	logp.Info("adcsamplerbeat is running! Hit CTRL-C to stop it.")

	var err error

	serial_config := serial.Config {
		Name: bt.config.SerialPort,
		Baud: bt.config.SerialBaud,
		Size: 8,
		StopBits: 1,
		Parity: 'N',
		ReadTimeout: bt.config.SerialTimeout,
	}

	fmt.Println("Opening serial: %s,%d",
				bt.config.SerialPort,
				bt.config.SerialBaud )

	// Open the TTY port.
	port, err := serial.OpenPort(&serial_config)
	if err != nil {
		fmt.Errorf("serial.Open: %v", err)
		return err
	}

	// On the STM32MP1 we need to send a char first
	n, err := port.Write([]byte("start"))
	if err != nil {
		fmt.Errorf("serial.Write: %v", err)
		return err
	}

	bt.client, err = b.Publisher.Connect()
	if err != nil {
		fmt.Errorf("Publisher.Connect: %v", err)
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		buf := make([]byte, 512)
		
		n, _ = port.Read(buf)
		s := string(buf[:n])
		s1 := strings.Split(s,"\n")	// split new lines
		if len(s1) > 2 && len(s1[1]) > 16 {
			fmt.Println("s1: ", s1[1])
			s2 := strings.SplitAfterN(s1[1], ":", 2)
			fmt.Println("s2: ", s2[1])
			s3 := strings.Split(s2[1], ",")
			fmt.Println("adc1_val: ", s3[0])
			fmt.Println("adc2_val: ", s3[1])
			fmt.Println("adc3_val: ", s3[2])
			fmt.Println("adc4_val: ", s3[3])
			adc1_val, _ := strconv.ParseFloat(s3[0], 32)
			adc2_val, _ := strconv.ParseFloat(s3[1], 32)
			adc3_val, _ := strconv.ParseFloat(s3[2], 32)
			adc4_val, _ := strconv.ParseFloat(s3[3], 32)

			event := beat.Event {
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"type":    b.Info.Name,
					"counter": counter,
					"adc1_val": adc1_val,
					"adc2_val": adc2_val,
					"adc3_val": adc3_val,
					"adc4_val": adc4_val,
				},
			}
			bt.client.Publish(event)
			logp.Info("Event sent")
			counter++
		}
	}
}

// Stop stops adcsamplerbeat.
func (bt *adcsamplerbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
