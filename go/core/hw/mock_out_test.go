package hw

import (
	"github.com/G0WCZ/cwc/config"
	"gotest.tools/assert"
	"testing"
)

func makeMockOut() MorseOut {
	c := config.DefaultConfig()
	mo := NewMockOut(c, "mock", "")
	mo.Open()

	return mo
}

func TestNewMockOut(t *testing.T) {
	mo := makeMockOut()
	assert.Equal(t, mo.(*MockOut).Bit(), false)
}

func TestMockOut_Bit(t *testing.T) {
	mo := makeMockOut()

	assert.Equal(t, mo.(*MockOut).Bit(), false)

	mo.SetBit(true)

	assert.Equal(t, mo.(*MockOut).Bit(), true)
}

func TestMockOut_Tone(t *testing.T) {
	mo := makeMockOut()

	assert.Equal(t, mo.(*MockOut).Tone(), false)

	mo.SetToneOut(true)

	assert.Equal(t, mo.(*MockOut).Tone(), true)
}

func TestMockOut_Status(t *testing.T) {
	mo := makeMockOut()

	assert.Equal(t, mo.(*MockOut).Status(), false)

	mo.SetStatusLED(true)

	assert.Equal(t, mo.(*MockOut).Status(), true)
}
