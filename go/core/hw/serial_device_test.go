package hw

import (
	"testing"
)

func TestClosePort(t *testing.T) {
	ClosePort("/dev/null")
}
