package share_code

import (
	"math/big"
	"regexp"

	"github.com/Cludch/csgo-tools/pkg/util"

	"strings"
)

type errInvalidShareCode struct{}

func (e errInvalidShareCode) Error() string {
	return "share_code: invalid share code"
}

func IsInvalidShareCodeError(err error) bool {
	_, ok := err.(errInvalidShareCode)
	return ok
}

// ShareCodeData holds the decoded match data and encoded share code.
type ShareCodeData struct {
	Encoded   string `json:"encoded" bson:"encoded"`
	OutcomeID uint64 `json:"outcomeId" bson:"outcomeId"`
	MatchID   uint64 `json:"matchId" bson:"matchId"`
	Token     uint32 `json:"token" bson:"token"`
}

// dictionary is used for the share code decoding.
const dictionary = "ABCDEFGHJKLMNOPQRSTUVWXYZabcdefhijkmnopqrstuvwxyz23456789"

// Used for share code decoding.
var bitmask64 uint64 = 18446744073709551615

// Validate validates an string whether the format matches a valid share code.
func Validate(code string) bool {
	var validateRe = regexp.MustCompile(`^CSGO(-?[\w]{5}){5}$`)
	return validateRe.MatchString(code)
}

// DecodeShareCode decodes the share code. Taken from ValvePython/csgo.
func Decode(code string) (*ShareCodeData, error) {
	if !Validate(code) {
		return nil, &errInvalidShareCode{}
	}

	var re = regexp.MustCompile(`^CSGO|\-`)
	s := re.ReplaceAllString(code, "")
	s = util.Reverse(s)

	bigNumber := big.NewInt(0)

	for _, c := range s {
		bigNumber = bigNumber.Mul(bigNumber, big.NewInt(int64(len(dictionary))))
		bigNumber = bigNumber.Add(bigNumber, big.NewInt(int64(strings.Index(dictionary, string(c)))))
	}

	a := SwapEndianness(bigNumber)

	matchid := big.NewInt(0)
	outcomeid := big.NewInt(0)
	token := big.NewInt(0)

	matchid = matchid.And(a, big.NewInt(0).SetUint64(bitmask64))
	outcomeid = outcomeid.Rsh(a, 64)
	outcomeid = outcomeid.And(outcomeid, big.NewInt(0).SetUint64(bitmask64))
	token = token.Rsh(a, 128)
	token = token.And(token, big.NewInt(0xFFFF))

	shareCode := &ShareCodeData{MatchID: matchid.Uint64(), OutcomeID: outcomeid.Uint64(), Token: uint32(token.Uint64()), Encoded: code}
	return shareCode, nil
}

// swapEndianness changes the byte order.
func SwapEndianness(number *big.Int) *big.Int {
	result := big.NewInt(0)

	left := big.NewInt(0)
	rightTemp := big.NewInt(0)
	rightResult := big.NewInt(0)

	for i := 0; i < 144; i += 8 {
		left = left.Lsh(result, 8)
		rightTemp = rightTemp.Rsh(number, uint(i))
		rightResult = rightResult.And(rightTemp, big.NewInt(0xFF))
		result = left.Add(left, rightResult)
	}

	return result
}
