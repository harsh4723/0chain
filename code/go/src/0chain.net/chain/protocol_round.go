package chain

import (
	"context"
	"fmt"
	"time"

	"0chain.net/node"

	"0chain.net/block"
	"0chain.net/common"
	"0chain.net/datastore"
	. "0chain.net/logging"
	"0chain.net/round"
	"0chain.net/util"
	metrics "github.com/rcrowley/go-metrics"
	"go.uber.org/zap"
)

var DELTA = 200 * time.Millisecond
var FINALIZATION_TIME = 2 * DELTA

/*SetNetworkRelayTime - setup the network relay time */
func SetNetworkRelayTime(delta time.Duration) {
	DELTA = delta
	FINALIZATION_TIME = 2 * delta
}

//SteadyStateFinalizationTimer - a metric that tracks the steady state finality time (time between two successive finalized blocks in steady state)
var SteadyStateFinalizationTimer metrics.Timer
var ssFTs time.Time

//StartToFinalizeTimer - a metric that tracks the time a block is created to finalized
var StartToFinalizeTimer metrics.Timer

func init() {
	if SteadyStateFinalizationTimer != nil {
		metrics.Unregister("ss_finalization_time")
	}
	SteadyStateFinalizationTimer = metrics.GetOrRegisterTimer("ss_finalization_time", nil)
	if StartToFinalizeTimer != nil {
		metrics.Unregister("s2f_time")
	}
	StartToFinalizeTimer = metrics.GetOrRegisterTimer("s2f_time", nil)
}

/*FinalizeRound - starting from the given round work backwards and identify the round that can be
  assumed to be finalized as only one chain has survived.
  Note: It is that round and prior that actually get finalized.
*/
func (c *Chain) FinalizeRound(ctx context.Context, r *round.Round, bsh BlockStateHandler) {
	if !r.SetFinalizing() {
		return
	}
	time.Sleep(FINALIZATION_TIME)
	c.FinalizedRoundsChannel <- r
}

func (c *Chain) finalizeRound(ctx context.Context, r *round.Round, bsh BlockStateHandler) {
	lfb := c.ComputeFinalizedBlock(ctx, r)
	if lfb == nil {
		Logger.Debug("finalization - no decisive block to finalize yet or don't have all the necessary blocks", zap.Any("round", r.Number))
		return
	}
	if lfb.Hash == c.LatestFinalizedBlock.Hash {
		return
	}
	if lfb.Round < c.LatestFinalizedBlock.Round {
		Logger.Info("finalize round - TODO: need to repair", zap.Any("lf_round", c.LatestFinalizedBlock.Round), zap.Int64("new_lf_round", lfb.Round))
		return
	}
	plfb := c.LatestFinalizedBlock
	lfbHash := plfb.Hash
	c.LatestFinalizedBlock = lfb
	frchain := make([]*block.Block, 0, 1)
	for b := lfb; b != nil && b.Hash != lfbHash; b = b.PrevBlock {
		frchain = append(frchain, b)
	}
	if len(frchain) == 0 {
		return
	}
	fb := frchain[len(frchain)-1]
	if fb.PrevBlock == nil {
		Logger.Info("finalize round (missed blocks)", zap.Int64("from", plfb.Round+1), zap.Int64("to", fb.Round-1))
		c.MissedBlocks += fb.Round - 1 - plfb.Round
	}
	deadBlocks := make([]*block.Block, 0, 1)
	for idx := range frchain {
		fb := frchain[len(frchain)-1-idx]
		Logger.Info("finalize round", zap.Int64("round", r.Number), zap.Int64("finalized_round", fb.Round), zap.String("hash", fb.Hash))
		if time.Since(ssFTs) < 10*time.Second {
			SteadyStateFinalizationTimer.UpdateSince(ssFTs)
		}
		StartToFinalizeTimer.UpdateSince(fb.ToTime())
		ssFTs = time.Now()
		c.UpdateChainInfo(fb)
		if fb.ClientState != nil {
			ts := time.Now()
			err := fb.ClientState.SaveChanges(c.StateDB, util.Origin(fb.Round), false)
			if err != nil {
				Logger.Error("finalize round - save state", zap.Int64("round", fb.Round), zap.String("block", fb.Hash), zap.Duration("time", time.Since(ts)), zap.String("client_state", util.ToHex(fb.ClientStateHash)), zap.Int("changes", len(fb.ClientState.GetChangeCollector().GetChanges())), zap.Error(err))
			} else {
				Logger.Info("finalize round - save state", zap.Int64("round", fb.Round), zap.String("block", fb.Hash), zap.Duration("time", time.Since(ts)), zap.String("client_state", util.ToHex(fb.ClientStateHash)), zap.Int("changes", len(fb.ClientState.GetChangeCollector().GetChanges())))
			}
			c.rebaseState(fb)
			/*
				if config.DevConfiguration.State && stateOut != nil {
					fmt.Fprintf(stateOut, "round: %v block: %v state: %v prev_block: %v prev_state: %v\n", fb.Round, fb.Hash, util.ToHex(fb.ClientStateHash), fb.PrevHash, util.ToHex(fb.PrevBlock.ClientStateHash))
					fb.ClientState.PrettyPrint(stateOut)
					stateOut.Sync()
					if err = util.IsMPTValid(fb.ClientState); err != nil {
						fmt.Fprintf(stateOut, "round save: %v %v\n", util.ToHex(fb.ClientState.GetRoot()), fb.Round)
						panic(err)
					}
				}*/
		}
		bsh.UpdateFinalizedBlock(ctx, fb)
		c.BlockChain.Value = fb.GetSummary()
		c.BlockChain = c.BlockChain.Next()
		frb := c.GetRoundBlocks(fb.Round)
		for _, b := range frb {
			if b.Hash != fb.Hash {
				deadBlocks = append(deadBlocks, b)
			}
		}
	}
	// Prune all the dead blocks
	c.DeleteBlocks(deadBlocks)
	// Prune the chain from the oldest finalized block
	c.PruneChain(ctx, frchain[len(frchain)-1])
}

/*PruneChain - prunes the chain */
func (c *Chain) PruneChain(ctx context.Context, b *block.Block) {
	c.DeleteBlocksBelowRound(b.Round - 50)
}

/*GetNotarizedBlockForRound - get a notarized block for a round */
func (c *Chain) GetNotarizedBlockForRound(r *round.Round, nbrequestor node.EntityReceiveHandler) *block.Block {
	params := map[string]string{"round": fmt.Sprintf("%v", r.Number)}
	ctx, cancelf := context.WithCancel(context.TODO())
	handler := func(ctx context.Context, entity datastore.Entity) (interface{}, error) {
		Logger.Info("get notarized block for round", zap.Int64("round", r.Number), zap.String("block", entity.GetKey()))
		if r.Number+1 != c.CurrentRound {
			cancelf()
			return nil, nil
		}
		if r.GetBestNotarizedBlock() != nil {
			cancelf()
			return nil, nil
		}
		b, ok := entity.(*block.Block)
		if !ok {
			return nil, common.NewError("invalid_entity", "Invalid entity")
		}
		if b.Round != r.Number {
			return nil, common.NewError("invalid_block", "Block not from the requested round")
		}
		if err := b.Validate(ctx); err != nil {
			return nil, err
		}
		if err := c.VerifyNotarization(ctx, b, b.VerificationTickets); err != nil {
			return nil, err
		}
		//TODO: this may not be the best round block or the best chain weight block. Do we do that extra work?
		c.AddBlock(b)
		r.AddNotarizedBlock(b)
		Logger.Info("get notarized block", zap.Int64("round", r.Number), zap.String("block", b.Hash), zap.String("state", util.ToHex(b.ClientStateHash)))
		return nil, nil
	}
	c.Miners.RequestEntity(ctx, nbrequestor(params, handler))
	return r.GetBestNotarizedBlock()
}

/*GetNotarizedBlock - get a notarized block for a round */
func (c *Chain) GetNotarizedBlock(blockHash string, nbrequestor node.EntityReceiveHandler) *block.Block {
	cround := c.CurrentRound
	params := map[string]string{"block": blockHash}
	ctx, cancelf := context.WithCancel(context.TODO())
	var b *block.Block
	handler := func(ctx context.Context, entity datastore.Entity) (interface{}, error) {
		Logger.Info("get notarized block", zap.String("block", blockHash), zap.Int64("cround", cround), zap.Int64("current_round", c.CurrentRound))
		if cround != c.CurrentRound {
			cancelf()
			return nil, nil
		}
		nb, ok := entity.(*block.Block)
		if !ok {
			return nil, common.NewError("invalid_entity", "Invalid entity")
		}
		if err := nb.Validate(ctx); err != nil {
			Logger.Error("get notarized block - validate", zap.String("block", blockHash), zap.Error(err), zap.Any("block_obj", nb))
			return nil, err
		}
		if err := c.VerifyNotarization(ctx, nb, nb.VerificationTickets); err != nil {
			Logger.Error("get notarized block - validate notarization", zap.String("block", blockHash), zap.Error(err))
			return nil, err
		}
		b = nb
		c.AddBlock(b)
		Logger.Info("get notarized block", zap.Int64("round", b.Round), zap.String("block", b.Hash))
		return b, nil
	}
	c.Miners.RequestEntity(ctx, nbrequestor(params, handler))
	return b
}
