package hw

import (
	"github.com/G0WCZ/cwc/config"
)

type MockOut struct {
	Config      *config.Config
	bitValue    bool
	toneValue   bool
	statusValue bool
	adapterName string
	name        string
}

func NewMockOut(config *config.Config, name string, adapterName string) MorseOut {
	return &MockOut{
		Config:      config,
		bitValue:    false,
		toneValue:   false,
		statusValue: false,
		adapterName: adapterName,
		name:        name,
	}
}

func (m *MockOut) Open() error {
	m.bitValue = false
	m.toneValue = false
	m.statusValue = false

	return nil
}

func (m *MockOut) SetBit(bit bool) {
	m.bitValue = bit
}

func (m *MockOut) SetToneOut(bit bool) {
	m.toneValue = bit
}

func (m *MockOut) SetStatusLED(bit bool) {
	m.statusValue = bit
}

func (m *MockOut) Close() error {
	m.bitValue = false
	m.toneValue = false
	m.statusValue = false

	return nil
}

func (m *MockOut) Bit() bool {
	return m.bitValue
}

func (m *MockOut) SetBitValue(bit bool) {
	m.bitValue = bit
}

func (m *MockOut) Tone() bool {
	return m.toneValue
}

func (m *MockOut) Status() bool {
	return m.statusValue
}

func (m *MockOut) ConfigChanged() error {
	return nil
}

func (m *MockOut) Name() string {
	return m.name
}
