package share_code_test

import (
	"testing"

	"github.com/Cludch/csgo-tools/pkg/share_code"
	"github.com/stretchr/testify/assert"
)

const validShareCode = "CSGO-2cLcm-AiUKj-abhb4-kDVWK-ixnkP"
const invalidShareCode = "test"

func TestValidate_ValidCode(t *testing.T) {
	valid := share_code.Validate(validShareCode)
	assert.True(t, valid)
}

func TestValidate_InalidCode(t *testing.T) {
	valid := share_code.Validate(invalidShareCode)
	assert.False(t, valid)
}

func TestDecode_ValidCode(t *testing.T) {
	sc, err := share_code.Decode(validShareCode)
	assert.NotNil(t, sc)
	assert.Nil(t, err)
	assert.Equal(t, sc.Encoded, validShareCode)
	assert.NotEmpty(t, sc.MatchID)
	assert.NotEmpty(t, sc.OutcomeID)
	assert.NotEmpty(t, sc.Token)
}

func TestDecode_InvalidCode(t *testing.T) {
	sc, err := share_code.Decode(invalidShareCode)
	assert.Nil(t, sc)
	assert.NotNil(t, err)
}
