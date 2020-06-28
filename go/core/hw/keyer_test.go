package hw

import (
	"github.com/G0WCZ/cwc/config"
	"gotest.tools/assert"
	"testing"
)

func TestInitKeyer(t *testing.T) {
	c := config.DefaultConfig()
	InitKeyer(c)

	assert.Equal(t, len(keyerStates), 1)
	assert.Equal(t, keyerStates[0].state, CHECK)
	assert.Equal(t, keyerStates[0].dotMemory, false)
	assert.Equal(t, keyerStates[0].dashMemory, false)
	assert.Equal(t, keyerStates[0].kDelay, 0)
	assert.Equal(t, keyerStates[0].dotDelay, 1200/c.Keyer.Speed)
	assert.Equal(t, keyerStates[0].dashDelay, keyerStates[0].dotDelay*3*c.Keyer.Weight/50)
	assert.Equal(t, keyerStates[0].out, false)
	assert.Equal(t, keyerStates[0].weight, c.Keyer.Weight)
	assert.Equal(t, keyerStates[0].mode, c.Keyer.Mode)
	assert.Equal(t, keyerStates[0].speed, c.Keyer.Speed)
	assert.Equal(t, keyerStates[0].reverse, c.Keyer.Reverse)
	assert.Equal(t, keyerStates[0].spacing, c.Keyer.LetterSpace)
}

func TestResetKeyers(t *testing.T) {
	c := config.DefaultConfig()
	ResetKeyers()
	InitKeyer(c)
	assert.Equal(t, len(keyerStates), 1)

	ResetKeyers()
	assert.Equal(t, len(keyerStates), 0)
}

func TestSampleDitAndDahPaddle(t *testing.T) {
	c := config.DefaultConfig()
	i := NewMockIn(c, "mock", KEYER)
	i.(*MockIn).SetDit(true)
	i.(*MockIn).SetDah(false)

	assert.Equal(t, SampleDitPaddle(false, i), true)
	assert.Equal(t, SampleDitPaddle(true, i), false)
	assert.Equal(t, SampleDahPaddle(false, i), false)
	assert.Equal(t, SampleDahPaddle(true, i), true)
}
