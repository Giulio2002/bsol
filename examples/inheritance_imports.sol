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

contract BenchmarkB is A {

    function Active() public {
        super.active();
    }
}

contract BenchmarkInheritance is BenchmarkTest {
    A active;
    constructor() public {
        active = new A();
    }
}

