package stateutil_test

import (
	"testing"

	"github.com/prysmaticlabs/go-ssz"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	"github.com/prysmaticlabs/prysm/shared/featureconfig"
	"github.com/prysmaticlabs/prysm/shared/testutil"
)

func TestBlockRoot(t *testing.T) {
	f := featureconfig.Get()
	f.EnableBlockHTR = true
	featureconfig.Init(f)
	genState, keys := testutil.DeterministicGenesisState(t, 100)
	blk, err := testutil.GenerateFullBlock(genState, keys, testutil.DefaultBlockGenConfig(), 10)
	if err != nil {
		t.Fatal(err)
	}
	expectedRoot, err := ssz.HashTreeRoot(blk.Block)
	if err != nil {
		t.Fatal(err)
	}
	receivedRoot, err := stateutil.BlockRoot(blk.Block)
	if err != nil {
		t.Fatal(err)
	}
	if receivedRoot != expectedRoot {
		t.Fatalf("Wanted %#x but got %#x", expectedRoot, receivedRoot)
	}
	blk, err = testutil.GenerateFullBlock(genState, keys, testutil.DefaultBlockGenConfig(), 100)
	if err != nil {
		t.Fatal(err)
	}
	expectedRoot, err = ssz.HashTreeRoot(blk.Block)
	if err != nil {
		t.Fatal(err)
	}
	receivedRoot, err = stateutil.BlockRoot(blk.Block)
	if err != nil {
		t.Fatal(err)
	}
	if receivedRoot != expectedRoot {
		t.Fatalf("Wanted %#x but got %#x", expectedRoot, receivedRoot)
	}
}
