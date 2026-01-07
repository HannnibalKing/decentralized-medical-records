// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../RevocationRegistry.sol";

contract RevocationRegistryTest is Test {
    RevocationRegistry registry;

    function setUp() public {
        registry = new RevocationRegistry();
    }

    function testRevokeOnce() public {
        bytes32 handle = keccak256("cap");
        registry.revoke(handle, "test");
        (bool revoked,, address revoker,) = registry.get(handle);
        assertTrue(revoked, "handle should be revoked");
        assertEq(revoker, address(this), "revoker should be caller");
    }

    function testIdempotent() public {
        bytes32 handle = keccak256("cap2");
        registry.revoke(handle, "first");
        registry.revoke(handle, "second");
        (bool revoked, uint64 ts, address revoker, string memory reason) = registry.get(handle);
        assertTrue(revoked, "revoked");
        assertGt(ts, 0, "timestamp set");
        assertEq(revoker, address(this), "revoker stable");
        assertEq(reason, "first", "reason should not change on repeated revoke");
    }
}
