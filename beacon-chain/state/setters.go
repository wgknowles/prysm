package state

import (
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/go-bitfield"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/stateutil"
)

type fieldIndex int

const (
	slot fieldIndex = iota
	genesisTime
	fork
	latestBlockHeader
	blockRoots
	stateRoots
	historicalRoots
	eth1Data
	eth1DataVotes
	eth1DepositIndex
	validators
	balances
	randaoMixes
	slashings
	previousEpochAttestations
	currentEpochAttestations
	justificationBits
	previousJustifiedCheckpoint
	currentJustifiedCheckpoint
	finalizedCheckpoint
)

func (b *BeaconState) SetGenesisTime(val uint64) error {
	b.state.GenesisTime = val
	root := stateutil.Uint64Root(val)
	b.lock.Lock()
	b.merkleLayers[0][genesisTime] = root[:]
	b.recomputeRoot(int(genesisTime))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetSlot(val uint64) error {
	b.state.Slot = val
	root := stateutil.Uint64Root(val)
	b.lock.Lock()
	b.merkleLayers[0][slot] = root[:]
	b.recomputeRoot(int(slot))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetFork(val *pbp2p.Fork) error {
	root, err := stateutil.ForkRoot(val)
	if err != nil {
		return err
	}
	b.state.Fork = val
	b.lock.Lock()
	b.merkleLayers[0][fork] = root[:]
	b.recomputeRoot(int(fork))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetLatestBlockHeader(val *ethpb.BeaconBlockHeader) error {
	root, err := stateutil.BlockHeaderRoot(val)
	if err != nil {
		return err
	}
	b.state.LatestBlockHeader = val
	b.lock.Lock()
	b.merkleLayers[0][latestBlockHeader] = root[:]
	b.recomputeRoot(int(latestBlockHeader))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetBlockRoots(val [][]byte) error {
	root, err := stateutil.ArraysRoot(val, params.BeaconConfig().SlotsPerHistoricalRoot, "BlockRoots")
	if err != nil {
		return err
	}
	b.state.BlockRoots = val
	b.lock.Lock()
	b.merkleLayers[0][blockRoots] = root[:]
	b.recomputeRoot(int(blockRoots))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetStateRoots(val [][]byte) error {
	root, err := stateutil.ArraysRoot(val, params.BeaconConfig().SlotsPerHistoricalRoot, "StateRoots")
	if err != nil {
		return err
	}
	b.state.StateRoots = val
	b.lock.Lock()
	b.merkleLayers[0][stateRoots] = root[:]
	b.recomputeRoot(int(stateRoots))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetHistoricalRoots(val [][]byte) error {
	root, err := stateutil.HistoricalRootsRoot(b.state.HistoricalRoots)
	if err != nil {
		return err
	}
	b.state.HistoricalRoots = val
	b.lock.Lock()
	b.merkleLayers[0][historicalRoots] = root[:]
	b.recomputeRoot(int(historicalRoots))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetEth1Data(val *ethpb.Eth1Data) error {
	root, err := stateutil.Eth1Root(b.state.Eth1Data)
	if err != nil {
		return err
	}
	b.state.Eth1Data = val
	b.lock.Lock()
	b.merkleLayers[0][eth1Data] = root[:]
	b.recomputeRoot(int(eth1Data))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetEth1DataVotes(val []*ethpb.Eth1Data) error {
	root, err := stateutil.Eth1DataVotesRoot(b.state.Eth1DataVotes)
	if err != nil {
		return err
	}
	b.state.Eth1DataVotes = val
	b.lock.Lock()
	b.merkleLayers[0][eth1DataVotes] = root[:]
	b.recomputeRoot(int(eth1DataVotes))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetEth1DepositIndex(val uint64) error {
	b.state.Eth1DepositIndex = val
	root := stateutil.Uint64Root(val)
	b.lock.Lock()
	b.merkleLayers[0][eth1DepositIndex] = root[:]
	b.recomputeRoot(int(eth1DepositIndex))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetValidators(val []*ethpb.Validator) {
	b.state.Validators = val
}

func (b *BeaconState) SetBalances(val []uint64) error {
	root, err := stateutil.ValidatorBalancesRoot(b.state.Balances)
	if err != nil {
		return err
	}
	b.state.Balances = val
	b.lock.Lock()
	b.merkleLayers[0][balances] = root[:]
	b.recomputeRoot(int(balances))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetRandaoMixes(val [][]byte) error {
	root, err := stateutil.ArraysRoot(val, params.BeaconConfig().EpochsPerHistoricalVector, "RandaoMixes")
	if err != nil {
		return err
	}
	b.state.RandaoMixes = val
	b.lock.Lock()
	b.merkleLayers[0][randaoMixes] = root[:]
	b.recomputeRoot(int(randaoMixes))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetSlashings(val []uint64) error {
	root, err := stateutil.SlashingsRoot(b.state.Slashings)
	if err != nil {
		return err
	}
	b.state.Slashings = val
	b.lock.Lock()
	b.merkleLayers[0][slashings] = root[:]
	b.recomputeRoot(int(slashings))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetPreviousEpochAttestations(val []*pbp2p.PendingAttestation) {
	b.state.PreviousEpochAttestations = val
}

func (b *BeaconState) SetCurrentEpochAttestations(val []*pbp2p.PendingAttestation) {
	b.state.CurrentEpochAttestations = val
}

func (b *BeaconState) SetJustificationBits(val bitfield.Bitvector4) error {
	root := bytesutil.ToBytes32(b.state.JustificationBits)
	b.state.JustificationBits = val
	b.lock.Lock()
	b.merkleLayers[0][justificationBits] = root[:]
	b.recomputeRoot(int(justificationBits))
	b.lock.Unlock()
}

func (b *BeaconState) SetPreviousJustifiedCheckpoint(val *ethpb.Checkpoint) error {
	root, err := stateutil.CheckpointRoot(b.state.PreviousJustifiedCheckpoint)
	if err != nil {
		return err
	}
	b.state.PreviousJustifiedCheckpoint = val
	b.lock.Lock()
	b.merkleLayers[0][previousJustifiedCheckpoint] = root[:]
	b.recomputeRoot(int(previousJustifiedCheckpoint))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetCurrentJustifiedCheckpoint(val *ethpb.Checkpoint) error {
	root, err := stateutil.CheckpointRoot(b.state.CurrentJustifiedCheckpoint)
	if err != nil {
		return err
	}
	b.state.CurrentJustifiedCheckpoint = val
	b.lock.Lock()
	b.merkleLayers[0][currentJustifiedCheckpoint] = root[:]
	b.recomputeRoot(int(currentJustifiedCheckpoint))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) SetFinalizedCheckpoint(val *ethpb.Checkpoint) error {
	root, err := stateutil.CheckpointRoot(b.state.FinalizedCheckpoint)
	if err != nil {
		return err
	}
	b.state.FinalizedCheckpoint = val
	b.lock.Lock()
	b.merkleLayers[0][finalizedCheckpoint] = root[:]
	b.recomputeRoot(int(finalizedCheckpoint))
	b.lock.Unlock()
	return nil
}

func (b *BeaconState) recomputeRoot(idx int) {
	layers := b.merkleLayers
	// The merkle tree structure looks as follows:
	// [[r1, r2, r3, r4], [parent1, parent2], [root]]
	// Using information about the index which changed, idx, we recompute
	// only its branch up the tree.
	currentIndex := idx
	root := b.merkleLayers[0][idx]
	for i := 0; i < len(layers)-1; i++ {
		isLeft := currentIndex%2 == 0
		neighborIdx := currentIndex ^ 1

		neighbor := make([]byte, 32)
		if layers[i] != nil && len(layers[i]) != 0 && neighborIdx < len(layers[i]) {
			neighbor = layers[i][neighborIdx]
		}
		if isLeft {
			parentHash := hashutil.Hash(append(root, neighbor...))
			root = parentHash[:]
		} else {
			parentHash := hashutil.Hash(append(neighbor, root...))
			root = parentHash[:]
		}
		parentIdx := currentIndex / 2
		// Update the cached layers at the parent index.
		layers[i+1][parentIdx] = root
		currentIndex = parentIdx
	}
	b.merkleLayers = layers
}