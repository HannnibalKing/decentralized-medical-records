// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../AttestationRegistry.sol";

contract AttestationRegistryTest is Test {
    AttestationRegistry registry;

    function setUp() public {
        registry = new AttestationRegistry();
    }

    function testRecordAndGet() public {
        bytes32 cidMr = keccak256("mr");
        bytes32 cidMd = keccak256("md");
        bytes32 aggSig = keccak256("agg");
        registry.record(cidMr, cidMd, aggSig, AttestationRegistry.Status.Issued);

        (bytes32 gotCidMd, bytes32 gotAgg, AttestationRegistry.Status status, uint64 ts) = registry.get(cidMr);
        assertEq(gotCidMd, cidMd, "cidMd");
        assertEq(gotAgg, aggSig, "aggSig");
        assertEq(uint256(status), uint256(AttestationRegistry.Status.Issued), "status");
        assertGt(ts, 0, "timestamp");
    }

    function testUpdateStatus() public {
        bytes32 cidMr = keccak256("mr2");
        registry.record(cidMr, bytes32(0), bytes32(0), AttestationRegistry.Status.Issued);
        registry.updateStatus(cidMr, AttestationRegistry.Status.Void);
        (, , AttestationRegistry.Status status,) = registry.get(cidMr);
        assertEq(uint256(status), uint256(AttestationRegistry.Status.Void), "status updated");
    }
}
