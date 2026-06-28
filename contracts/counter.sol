// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Counter {
    uint256 public count;

    event Incremented(uint256 newCount, address indexed caller);

    function getCount() public view returns (uint256) {
        return count;
    }

    function increment() public {
        count += 1;
        emit Incremented(count, msg.sender);
    }
}
