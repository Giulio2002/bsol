package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/params"
)

func bsol(runs uint64, time bool, source ...string) error {
	// Configure and generate a sample block chain
	var (
		memDb   = memorydb.New()
		db      = rawdb.NewDatabase(memDb)
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress(key.PublicKey)
		gspec   = &core.Genesis{
			GasLimit: 800000000,
			Config: &params.ChainConfig{
				HomesteadBlock:      big.NewInt(0),
				EIP150Block:         big.NewInt(0),
				EIP155Block:         big.NewInt(0),
				EIP158Block:         big.NewInt(0),
				ByzantiumBlock:      big.NewInt(0),
				ConstantinopleBlock: big.NewInt(0),
				PetersburgBlock:     big.NewInt(0),
				IstanbulBlock:       big.NewInt(0),
			},
			Alloc: core.GenesisAlloc{
				address: {Balance: big.NewInt(9000000000000000000)},
			},
		}
	)
	engine := ethash.NewFaker()
	chainConfig, _, err := core.SetupGenesisBlock(db, gspec)
	if err != nil {
		return fmt.Errorf("Error: %s\n", err.Error())
	}
	blockchain, err := core.NewBlockChain(db, &core.CacheConfig{
		TrieDirtyDisabled: true,
	}, chainConfig, engine, vm.Config{}, nil)
	if err != nil {
		return fmt.Errorf("Error: %s\n", err.Error())

	}
	_ = blockchain.StateCache().TrieDB()
	// construct the first diff

	contractBackend := backends.NewSimulatedBackend(gspec.Alloc, gspec.GasLimit)
	transactOpts := bind.NewKeyedTransactor(key)
	transactOpts.GasPrice = big.NewInt(1)
	_ = transactOpts
	_ = contractBackend
	data, err := deployBenchmarks(contractBackend, transactOpts, source)
	if err != nil {
		return fmt.Errorf("Error: %s\n", err.Error())

	}
	err = executeBenchmarks(contractBackend, transactOpts, data, runs, time)
	if err != nil {
		return fmt.Errorf("Error: %s\n", err.Error())
	}
	return nil
}
