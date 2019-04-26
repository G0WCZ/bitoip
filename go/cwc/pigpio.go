package cwc

import (
	"github.com/golang/glog"
	"github.com/stianeikeland/go-rpio"
	"log"
	"strconv"
)

// PWM settings
const OnDutyCycle = uint32(1)
const PWMCycleLength = uint32(32)

type PiGPIO struct {
	config ConfigMap
	output rpio.Pin
	input rpio.Pin
	pwm bool
	pwmOut rpio.Pin
}

func NewPiGPIO() *PiGPIO {
	pigpio := PiGPIO{
		config: make(ConfigMap),
		pwm: false,
	}
	return &pigpio
}

func (g *PiGPIO) Open() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	sFreq, err := strconv.Atoi(g.config["sidetoneFreq"])

	if err != nil {
		log.Fatalf("Bad sidetone frequency")
	}

	glog.Info("setting sidetone to %d", sFreq)

	// PCM output
	if (sFreq > 0) {
		g.pwm = true
		g.pwmOut = rpio.Pin(13)
		g.pwmOut.Mode(rpio.Pwm)
		g.pwmOut.Freq(sFreq * 32)
		g.pwmOut.DutyCycle(0, 32)
	}

	// Pin output
	g.output = rpio.Pin(17) // header pin 11 BCM17
	g.output.Output()
	g.output.Low()

    g.input = rpio.Pin(27) // header pin 13 BCM27
    g.input.Input()
    g.input.PullUp()

    return nil
}

func (g *PiGPIO) SetConfig(key string, value string) {
	g.config[key] = value
}

func (g *PiGPIO) ConfigMap() ConfigMap {
	return g.config
}

func (g *PiGPIO) Bit() bool {
	if g.input.Read() == rpio.High {
		return false
	} else {
		return true
	}
}

func (g *PiGPIO) SetBit(bit0 bool) {
	if bit0 {
		g.output.High()
		g.SetPwm(true)
	} else {
		g.output.Low()
		g.SetPwm(false)
	}
}

func (g *  PiGPIO) SetPwm(v bool) {
	if g.pwm {
		var dutyLen uint32

		if v {
			dutyLen = OnDutyCycle
		} else {
			dutyLen = 0
		}
		g.pwmOut.DutyCycle(dutyLen, PWMCycleLength)
	}
}

func (g *PiGPIO) Close() {
	// pass
}

