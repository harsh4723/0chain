package event

import (
	"fmt"

	"gorm.io/gorm"

	"0chain.net/chaincore/state"
)

type Validator struct {
	gorm.Model
	ValidatorID string `json:"validator_id" gorm:"index:validator_id"`
	BaseUrl     string `json:"url" gorm:"index:url"`
	Stake       int64  `json:"stake" gorm:"index:stake"`

	// StakePoolSettings
	DelegateWallet string        `json:"delegate_wallet"`
	MinStake       state.Balance `json:"min_stake"`
	MaxStake       state.Balance `json:"max_stake"`
	NumDelegates   int           `json:"num_delegates"`
	ServiceCharge  float64       `json:"service_charge"`
}

func (vn *Validator) exists(edb *EventDb) (bool, error) {
	var count int64
	result := edb.Get().
		Model(&Validator{}).
		Where(&Validator{ValidatorID: vn.ValidatorID}).
		Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error searching for Validator %v, error %v",
			vn.ValidatorID, result.Error)
	}
	return count > 0, nil
}

func (edb *EventDb) GetValidatorByValidatorID(validatorID string) (Validator, error) {
	var vn Validator

	result := edb.Store.Get().Model(&Validator{}).Where(&Validator{ValidatorID: validatorID}).First(vn)

	if result.Error != nil {
		return vn, fmt.Errorf("error retriving Validation node with ID %v; error: %v", validatorID, result.Error)
	}

	return vn, nil
}

func (edb *EventDb) overwriteValidator(vn Validator) error {

	result := edb.Store.Get().Model(&Validator{}).Where(&Validator{ValidatorID: vn.ValidatorID}).Updates(&vn)
	return result.Error
}

func (edb *EventDb) addOrOverwriteValidator(vn Validator) error {
	exists, err := vn.exists(edb)
	if err != nil {
		return err
	}
	if exists {
		return edb.overwriteValidator(vn)
	}

	result := edb.Store.Get().Create(&vn)

	return result.Error
}