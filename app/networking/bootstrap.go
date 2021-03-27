package networking

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/itzmeanjan/harmony/app/data"
)

var memPool *data.MemPool
var redisClient *redis.Client
var reconnectionManager *ReconnectionManager

// InitMemPool - Initializing mempool handle, in this module
// so that it can be used updating local mempool state, when new
// deserialisable tx chunk is received from any peer, over p2p network
func InitMemPool(pool *data.MemPool) error {

	if pool != nil {
		memPool = pool
		return nil
	}

	return errors.New("bad mempool received in p2p networking handler")

}

// InitRedisClient - Initializing redis client handle, so that all
// subscriptions can be done using this client
func InitRedisClient(client *redis.Client) error {

	if client != nil {
		redisClient = client
		return nil
	}

	return errors.New("bad redis client received in p2p networking handler")

}

// InitReconnectionManager - Initialising it so that it's available
// in global scope with in this package & it can be used by multiple
// workers, when they encounter problem i.e. stream gets reset
func InitReconnectionManager(reconnMgr *ReconnectionManager) error {

	if reconnMgr != nil {
		reconnectionManager = reconnMgr
		return nil
	}

	return errors.New("bad reconnection manager received in p2p networking handler")

}

// Setup - Bootstraps `harmony`'s p2p networking stack
func Setup(ctx context.Context, comm chan struct{}) error {

	if !(memPool != nil && redisClient != nil) {
		return errors.New("mempool/ redisClient instance not initialised")
	}

	// Attempt to create a new `harmony` node
	// with p2p networking capabilities
	host, err := CreateHost(ctx)
	if err != nil {
		return err
	}

	// Display info regarding this node
	ShowHost(host)
	// Start listening for incoming streams, for supported protocol
	Listen(host)

	go SetUpPeerDiscovery(ctx, host, comm)

	reconnMgr := NewReconnectionManager(host)
	if err := InitReconnectionManager(reconnMgr); err != nil {
		return err
	}

	// Starting this worker as a seperate go routine
	go reconnMgr.Start(ctx)

	return nil

}
