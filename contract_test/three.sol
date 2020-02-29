pragma solidity ^0.5.0;

contract TestThree {

    function BenchmarkLoop() public {
        for (uint index = 0; index < 1000; index++) {}
    }
}

