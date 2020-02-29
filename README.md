# BETH (Benchmarking Ethereum)
Beth is a tool to write benchmark for Solidity snippets and Smart Contract.

Beth gives as an out put:
* The Average Gas Usage
* The Average Gas Usage for execution only
* Average Computation time for Ethereum to execute the code in µs(nanoseconds).
## Install
```
git clone https://github.com/Giulio2002/beth
cd beth
sudo sh install.sh
```
## Usage
Given this benchmark contract
```js
pragma solidity ^0.5.0;

contract N {

    function BenchmarkOne() public returns(uint) {
        return 1;
    }

    function BenchmarkTwo() public returns(uint) {
        return 2;
    }
}
```
to execute the benchmarks just uses `beth --sol N.sol`.

```
Contract: N
Method: N.BenchmarkOne()
Average Computation time: 32.742638µs
Average Gas Usage: 21262 Gas
Average Gas Usage per execution: 262 Gas

Method: N.BenchmarkTwo()
Average Computation time: 34.421126µs
Average Gas Usage: 21284 Gas
Average Gas Usage per execution: 284 Gas
```
BETH benchmarks every method of every smart contract in a given solidity file that:
* Benchmark at the beggining of its name (Are ignored)
* Requires 0 arguments (Throws error)

so the following bechmarks won't be executed:
```js
pragma solidity ^0.5.0;

contract N {

    function One() public returns(uint) { // Does no have Benchmark at the beggining of the name
        return 1;
    }

    function BenchmarkTwo(uint a) public returns(uint) { // Requires an argument
        return a;
    }
}
```
## Flag Options
```
--runs: Number of Runs to calculate averages
--sol: Solidity source file to benchmark
--sol-dir: Solidity source directory to benchmark
```