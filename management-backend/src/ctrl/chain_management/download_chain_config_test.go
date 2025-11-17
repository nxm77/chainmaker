package chain_management

import (
	"testing"

	"gotest.tools/assert"
)

const confYml = "../../../dependence"

// TestCreateBinAndLib
//
//	@Description:
//	@param t
func TestCreateBinAndLib(t *testing.T) {
	err := createBinAndLib("chain1", "org1", "node1", confYml, 1)
	assert.Equal(t, err, nil)
}

// TestCreateLib
//
//	@Description:
//	@param t
func TestCreateLib(t *testing.T) {
	err := createLib("org1", "node1", confYml)
	assert.Equal(t, err, nil)
}
