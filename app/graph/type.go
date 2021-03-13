package graph

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/itzmeanjan/harmony/app/data"
)

// PublishingCriteria - Message publishing criteria is expected
// in this form, which is to be invoked, everytime new message
// is received in any topic client is subscribed to
type PublishingCriteria func(*data.MemPoolTx, ...interface{}) bool

// NoCriteria - When you want to listen to
// any tx being published on your topic of interest
// simply pass this function to `ListenToMessages`
// so that all criteria check always returns `true`
// & graphQL client receives all tx(s)
func NoCriteria(*data.MemPoolTx, ...interface{}) bool {
	return true
}

// CheckFromAddress - Just checks `from` address of tx, so that client
// is only notified when tx from that address is detected to be entering/ leaving
// mempool
func CheckFromAddress(m *data.MemPoolTx, params ...interface{}) bool {

	if len(params) != 1 {
		return false
	}

	// Attempting to assert type
	addr, ok := params[0].(common.Address)
	if !ok {
		return false
	}

	return addr == m.From

}

// CheckToAddress - Just checks `to` address of tx, so that client
// is only notified when tx `to` that address is detected to be entering/ leaving
// mempool
func CheckToAddress(m *data.MemPoolTx, params ...interface{}) bool {

	if len(params) != 1 {
		return false
	}

	// Attempting to assert type
	addr, ok := params[0].(common.Address)
	if !ok {
		return false
	}

	// Checking with `to` address of tx
	//
	// @note For tx(s) trying to deploy
	// contract, there'll be no `to` address
	//
	// That's why 👇 check
	if m.To == nil {
		return false
	}

	return *m.To == addr

}