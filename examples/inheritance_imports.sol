pragma solidity ^0.5.0;
// Imports are supported. Execute also other benchmarks within the import.
import "./one.sol";

contract A {
    bool public activated;
    constructor() public {
        activated = true;
    }

    function active() public {
        activated = !activated;
    }
}

contract B is A {

    function BenchmarkActive() public {
        super.active();
    }
}

contract P is TestOne {
    A active;
    constructor() public {
        active = new A();
    }
}

