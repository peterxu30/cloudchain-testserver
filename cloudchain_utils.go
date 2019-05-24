package main

import (
	"context"
	"sync"

	"github.com/peterxu30/cloudchain"
)

const (
	testProjectId    = "cloudchaintestserver"
	genesisBlockData = "genesis"
)

// type TestCloudChainManager struct {
// 	cc *cloudchain.CloudChain
// }

// func NewTestCloudChainManager(ctx context.Context) (*TestCloudChainManager, error) {
// 	cc, err := cloudchain.NewCloudChain(ctx, testProjectId, 10, []byte(genesisBlockData))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &TestCloudChainManager{
// 		cc: cc,
// 	}, nil
// }

var _testCloudChain *cloudchain.CloudChain
var _testCloudChainLock sync.RWMutex

func GetTestCloudChain(ctx context.Context) *cloudchain.CloudChain {
	_testCloudChainLock.RLock()
	defer _testCloudChainLock.RUnlock()

	if _testCloudChain == nil {
		cc, err := cloudchain.NewCloudChain(ctx, testProjectId, 10, []byte(genesisBlockData))
		if err != nil {
			panic(err)
		}
		_testCloudChain = cc
	}

	return _testCloudChain
}

func DeleteTestCloudChain(ctx context.Context) error {
	_testCloudChainLock.Lock()
	defer _testCloudChainLock.Unlock()

	if _testCloudChain == nil {
		return nil
	}

	err := cloudchain.DeleteCloudChain(ctx, _testCloudChain)
	if err != nil {
		return err
	}

	_testCloudChain = nil
	return nil
}
