package hashByID

import (
	"errors"
	"math"
	"strings"
	"urlShortener/utils/e"
)

const (
	alphabet = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM_0123456789"
	hashLen  = 10
)

var ErrOverFlow = errors.New("shortenUrl len overflow happend")
var maxURlCount = uint64(math.Pow(float64(len(alphabet)), hashLen))

type SeedGenerator interface {
	getID() uint64
}

type HashGenerator struct {
	SeedGenerator
}

func New(id uint64) *HashGenerator {
	return &HashGenerator{newIDGenerator(id)}
}

func (h *HashGenerator) Hash() (string, error) {
	const fn = "lib.linkShortening.Hash"
	hashBuilder := strings.Builder{}

	seed := h.getID()
	if seed > maxURlCount {
		return "", e.WrapError(fn, ErrOverFlow)
	}

	for ; seed > 0; seed /= uint64(len(alphabet)) {
		letter := alphabet[seed%uint64(len(alphabet))]
		hashBuilder.WriteByte(letter)
	}

	result := ""
	if hashBuilder.Len() < hashLen {
		result = strings.Repeat(string(alphabet[0]), hashLen-hashBuilder.Len())
	}
	return result + hashBuilder.String(), nil
}
