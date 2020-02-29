package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	runs := flag.Uint64("runs", 1500, "Count of runs per execution to calculate average. Default: 1500")
	sol := flag.String("sol", "", "Solidity file to benchmark")
	sol_dir := flag.String("sol-dir", "", "Directory of benchmark smart contracts")
	flag.Parse()
	if *sol != "" && *sol_dir != "" {
		fmt.Println("--sol and --sol-dir cannot be set at the same time.")
	}

	if *sol == "" && *sol_dir == "" {
		fmt.Println("No Solidity files to benchmark specified, please set --sol(single source) or --sol-dir (iterate over directory).")
	}

	if *sol != "" {
		err := beth(*runs, *sol)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if *sol_dir != "" {
		err := filepath.Walk(*sol_dir, func(path string, info os.FileInfo, err error) error {
			if !strings.Contains(path, ".sol") {
				return nil
			}
			return beth(*runs, path)
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
