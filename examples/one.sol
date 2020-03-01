pragma solidity ^0.5.0;

contract BenchmarkTest {

    function BenchmarkOne() public returns(uint) {
        return 1;
    }

    function BenchmarkTwo() public returns(uint) {
        return 2;
    }

    function Lol() public returns(uint) { // Won't be executed
        return 3;
    }
}

