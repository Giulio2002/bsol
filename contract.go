package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
)

type ContractData struct {
	address      common.Address
	a            abi.ABI
	contract     *bind.BoundContract
	names        []string
	contractName string
}

func deployBenchmarks(path string, contractBackend *backends.SimulatedBackend, opts *bind.TransactOpts) ([]ContractData, error) {
	contracts, err := compiler.CompileSolidity("", path)
	if err != nil {
		return nil, err
	}
	var data []ContractData
	for contractName, contract := range contracts {
		var names []string
		if err != nil {
			return nil, err
		}
		abiData, err := json.Marshal(contract.Info.AbiDefinition.([]interface{}))
		if err != nil {
			return nil, err
		}
		abi, err := abi.JSON(bytes.NewBuffer(abiData))
		if err != nil {
			return nil, err
		}
		address, _, c, err := bind.DeployContract(opts, abi, common.Hex2Bytes(contract.Code[2:]), contractBackend)
		for _, method := range contract.Info.AbiDefinition.([]interface{}) {
			mapped := method.(map[string]interface{})
			name := mapped["name"].(string)
			if strings.Index(name, "Benchmark") != 0 {
				continue
			}
			if len(mapped["inputs"].([]interface{})) != 0 {
				return nil, fmt.Errorf("Invalid Benchmark: %s: function should require 0 arguments, but it requires %d", name, len(mapped["inputs"].([]interface{})))
			}
			names = append(names, name)
		}
		if err != nil {
			return nil, err
		}
		data = append(data, ContractData{address, abi, c, names, strings.Split(contractName, ":")[1]})
	}
	contractBackend.Commit()
	return data, nil
}

func executeBenchmarks(contractBackend *backends.SimulatedBackend, opts *bind.TransactOpts, data []ContractData, runs uint64) error {
	for _, contractData := range data {
		fmt.Printf("Contract: %s", contractData.contractName)
		if len(contractData.names) == 0 {
			fmt.Printf(" (No Benchmarks)\n")
			continue
		} else {
			fmt.Println()
		}
		for _, method := range contractData.names {
			var totalGas uint64
			var tx *types.Transaction
			var err error
			var i uint64
			for ; i < runs; i++ {
				tx, err = contractData.contract.Transact(opts, method)
				if err != nil {
					return err
				}
				totalGas += tx.Gas()
			}
			start := time.Now()
			contractBackend.Commit()
			totalTime := convertElapsedToNano(time.Since(start).String())
			fmt.Printf("Method: %s.%s()\n", contractData.contractName, method)
			fmt.Printf("Average Computation time: %fµs\n", totalTime/float64(runs))
			fmt.Printf("Average Gas Usage: %d Gas\n", totalGas/runs)
			fmt.Printf("Average Gas Usage per execution: %d Gas\n\n", (totalGas/runs)-21000)
			if err != nil {
				return err
			}

		}
	}
	return nil
}
