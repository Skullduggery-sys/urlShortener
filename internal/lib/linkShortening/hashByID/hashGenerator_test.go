package hashByID

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math"
	"testing"
)

const getID = "getID"

type mockIDGenerator struct {
	mock.Mock
}

func (m *mockIDGenerator) getID() uint64 {
	args := m.Called()
	id := args.Get(0)
	return id.(uint64)
}

func TestHashLenLessThanHashLen(t *testing.T) {
	idGen := &mockIDGenerator{}
	hasher := HashGenerator{idGen}

	var id uint64 = 10
	expectedHash := "qqqqqqqqqa"
	idGen.On(getID).Return(id)

	resultHash, err := hasher.Hash()
	assert.NoError(t, err)
	assert.Equal(t, resultHash, expectedHash)

	assert.True(t, idGen.AssertExpectations(t))
}

func TestHashZeroID(t *testing.T) {
	idGen := &mockIDGenerator{}
	hasher := HashGenerator{idGen}

	var id uint64 = 0
	expectedHash := "qqqqqqqqqq"
	idGen.On(getID).Return(id)

	resultHash, err := hasher.Hash()
	assert.NoError(t, err)
	assert.Equal(t, resultHash, expectedHash)

	assert.True(t, idGen.AssertExpectations(t))
}

func TestHashRandomID(t *testing.T) {
	idGen := &mockIDGenerator{}
	hasher := HashGenerator{idGen}

	var id uint64 = 987
	expectedHash := "qqqqqqqqJh"
	idGen.On(getID).Return(id)

	resultHash, err := hasher.Hash()
	assert.NoError(t, err)
	assert.Equal(t, resultHash, expectedHash)

	assert.True(t, idGen.AssertExpectations(t))
}

func TestHashOverflow(t *testing.T) {
	idGen := &mockIDGenerator{}
	hasher := HashGenerator{idGen}

	id := uint64(math.Pow(float64(len(alphabet)), hashLen)) + 1
	idGen.On(getID).Return(id)

	_, err := hasher.Hash()
	assert.True(t, errors.Is(err, ErrOverFlow))

	assert.True(t, idGen.AssertExpectations(t))
}
