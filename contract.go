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
	a            abi.ABI
	code         string
	names        []string
	contractName string
}

func deployBenchmarks(contractBackend *backends.SimulatedBackend, opts *bind.TransactOpts, path []string) ([]ContractData, error) {
	contracts, err := compiler.CompileSolidity("", path...)
	if err != nil {
		return nil, err
	}
	var data []ContractData
	for contractName, contract := range contracts {
		if strings.Index(strings.Split(contractName, ":")[1], "Benchmark") != 0 {
			continue
		}
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
		for _, method := range contract.Info.AbiDefinition.([]interface{}) {
			mapped := method.(map[string]interface{})
			if mapped["name"] == nil {
				if len(mapped["inputs"].([]interface{})) != 0 {
					return nil, fmt.Errorf("Invalid Benchmark: %s: constructor should require 0 arguments")
				}
				continue
			}
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
		data = append(data, ContractData{abi, contract.Code, names, strings.Split(contractName, ":")[1]})
	}
	contractBackend.Commit()
	return data, nil
}

func executeBenchmarks(contractBackend *backends.SimulatedBackend, opts *bind.TransactOpts, data []ContractData, runs uint64, isTime bool) error {
	for _, contractData := range data {
		fmt.Printf("\nContract: %s", contractData.contractName)
		if len(contractData.names) == 0 {
			fmt.Printf(" (No Benchmarks)\n")
			continue
		} else {
			fmt.Println()
		}
		if isTime == false {
			runs = 1
		}
		for _, method := range contractData.names {
			var tx *types.Transaction
			var err error
			var i uint64
			var addresses []common.Address
			var totalTime float64
			for ; i < runs; i++ {
				addr, _, _, err := bind.DeployContract(opts, contractData.a, common.Hex2Bytes(contractData.code[2:]), contractBackend)
				if err != nil {
					return err
				}
				addresses = append(addresses, addr)
			}
			contractBackend.Commit()
			i = 0
			for ; i < runs; i++ {
				c := bind.NewBoundContract(addresses[i], contractData.a, contractBackend, contractBackend, contractBackend)
				if err != nil {
					return err
				}
				tx, err = c.Transact(opts, method)
				if err != nil {
					return err
				}
			}
			if isTime {
				start := time.Now()
				contractBackend.Commit()
				totalTime = convertElapsedToNano(time.Since(start).String())
			} else {
				contractBackend.Commit()
			}

			fmt.Printf("Method: %s.%s()\n", contractData.contractName, method)
			if isTime {
				fmt.Printf("Average Computation time: %fµs\n", totalTime/float64(runs))
			}
			fmt.Printf("Gas Usage: %d Gas\n", tx.Gas())
			fmt.Printf("Gas Usage per execution: %d Gas\n", tx.Gas()-21000)
			if err != nil {
				return err
			}

		}
	}
	return nil
}
