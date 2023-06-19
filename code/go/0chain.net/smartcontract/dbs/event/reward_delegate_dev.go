package event

import (
	"0chain.net/smartcontract/stakepool/spenum"
	"github.com/pkg/errors"
)

func (edb *EventDb) GetRewardsToDelegates(blockNumber, startBlockNumber, endBlockNumber string, rewardType int) ([]RewardDelegate, error) {

	if blockNumber != "" {
		var rds []RewardDelegate
		err := edb.Get().Where("block_number = ? AND reward_type = ?", blockNumber, rewardType).Find(&rds).Error
		return rds, err
	}

	if startBlockNumber != "" && endBlockNumber != "" {
		var rds []RewardDelegate
		err := edb.Get().Where("block_number >= ? AND block_number <= ? AND reward_type = ?", startBlockNumber, endBlockNumber, rewardType).Find(&rds).Error

		return rds, err
	}

	return nil, errors.Errorf("start or end block number can't be empty")

}

func (edb *EventDb) GetAllocationCancellationRewardsToDelegates(startBlock, endBlock string) ([]RewardDelegate, error) {

	var rps []RewardDelegate
	err := edb.Get().Where("block_number >= ? AND block_number <= ? AND reward_type = ?", startBlock, endBlock, spenum.CancellationChargeReward).Find(&rps).Error

	return rps, err
}
