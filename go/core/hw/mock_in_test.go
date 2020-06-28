package hw

import (
	"github.com/G0WCZ/cwc/config"
	"gotest.tools/assert"
	"testing"
)

func makeMockIn(adapter string) MorseIn {
	c := config.DefaultConfig()
	mi := NewMockIn(c, "mock", adapter)
	mi.Open()

	return mi
}

func TestNewMockIn(t *testing.T) {
	mi := makeMockIn("")
	assert.Equal(t, mi.Bit(), false)
	assert.Equal(t, mi.Dit(), false)
	assert.Equal(t, mi.Dah(), false)
}

func TestMockIn_Dah(t *testing.T) {
	mi := makeMockIn("")

	assert.Equal(t, mi.Dah(), false)
	mi.(*MockIn).SetDah(true)
	assert.Equal(t, mi.Dah(), true)
}

func TestMockIn_Dit(t *testing.T) {
	mi := makeMockIn("")

	assert.Equal(t, mi.Dit(), false)
	mi.(*MockIn).SetDit(true)
	assert.Equal(t, mi.Dit(), true)
}

func TestMockIn_Bit(t *testing.T) {
	mi := makeMockIn("")

	assert.Equal(t, mi.Bit(), false)
	mi.(*MockIn).SetBit(true)
	assert.Equal(t, mi.Bit(), true)
}

func TestMockIn_UseKeyer(t *testing.T) {
	mi := makeMockIn(KEYER)
	assert.Equal(t, mi.UseKeyer(), true)

	mi = makeMockIn("")
	assert.Equal(t, mi.UseKeyer(), false)
}

func TestMockIn_Close(t *testing.T) {
	mi := makeMockIn("")

	mi.(*MockIn).SetDit(true)
	mi.(*MockIn).SetDah(true)
	mi.(*MockIn).SetBit(true)

	mi.Close()

	assert.Equal(t, mi.Bit(), false)
	assert.Equal(t, mi.Dit(), false)
	assert.Equal(t, mi.Dah(), false)
}
