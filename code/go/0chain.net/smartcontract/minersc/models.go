package minersc

import (
	"encoding/json"
	"errors"
	"net/url"

	"0chain.net/chaincore/smartcontractstate"
)

var allMinersKey = smartcontractstate.Key("all_miners")

//MinerNode struct that holds information about the registering miner
type MinerNode struct {
	ID        string `json:"id"`
	BaseURL   string `json:"url"`
	PublicKey string `json:"-"`
}

func (mn *MinerNode) getKey() smartcontractstate.Key {
	return smartcontractstate.Key("miner:" + mn.ID)
}

func (mn *MinerNode) encode() []byte {
	buff, _ := json.Marshal(mn)
	return buff
}

func (mn *MinerNode) decodeFromValues(params url.Values) error {
	mn.BaseURL = params.Get("baseurl")
	mn.ID = params.Get("id")

	if mn.BaseURL == "" || mn.ID == "" {
		return errors.New("BaseURL or ID is not specified")
	}
	return nil

}

func (mn *MinerNode) decode(input []byte) error {
	err := json.Unmarshal(input, mn)
	if err != nil {
		return err
	}
	return nil
}
