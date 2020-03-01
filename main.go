package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	runs := flag.Uint64("runs", 500, "Count of runs per execution to calculate average. Default: 1500")
	time := flag.Bool("execution-time", false, "calculate average execution time")
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
		err := bsol(*runs, *time, *sol)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	_ = time
	if *sol_dir != "" {
		var src []string
		err := filepath.Walk(*sol_dir, func(path string, info os.FileInfo, err error) error {
			if !strings.Contains(path, ".sol") {
				return nil
			}
			src = append(src, path)
			return nil
		})
		err = bsol(*runs, *time, src...)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}
