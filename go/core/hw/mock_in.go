package hw

import "github.com/G0WCZ/cwc/config"

type MockIn struct {
	Config      *config.Config
	adapterName string
	keyer       bool
	bitValue    bool
	ditValue    bool
	dahValue    bool
	name        string
}

func NewMockIn(config *config.Config, name string, adapterName string) MorseIn {
	return &MockIn{
		Config:      config,
		adapterName: adapterName,
		keyer:       KEYER == adapterName,
		bitValue:    false,
		ditValue:    false,
		dahValue:    false,
		name:        name,
	}
}

func (m *MockIn) Open() error {
	m.bitValue = false
	m.ditValue = false
	m.dahValue = false

	return nil
}

func (m *MockIn) ConfigChanged() error {
	return nil
}

func (m *MockIn) Bit() bool {
	return m.bitValue
}

func (m *MockIn) Dit() bool {
	return m.ditValue
}

func (m *MockIn) Dah() bool {
	return m.dahValue
}

func (m *MockIn) Close() error {
	m.ditValue = false
	m.dahValue = false
	m.bitValue = false
	return nil
}

func (m *MockIn) UseKeyer() bool {
	return m.keyer
}

func (m *MockIn) SetBit(bit bool) {
	m.bitValue = bit
}

func (m *MockIn) SetDit(bit bool) {
	m.ditValue = bit
}

func (m *MockIn) SetDah(bit bool) {
	m.dahValue = bit
}

func (m *MockIn) Name() string {
	return m.name
}
