package util_test

import (
	"testing"

	"github.com/Cludch/csgo-tools/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestReverse(t *testing.T) {
	input := "abcdef"
	expectedOutput := "fedcba"

	output := util.Reverse(input)
	assert.NotNil(t, output)
	assert.Equal(t, output, expectedOutput)
}
