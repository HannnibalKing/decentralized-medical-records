// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../MerkleAnchor.sol";

contract MerkleAnchorTest is Test {
    MerkleAnchor anchor;

    function setUp() public {
        anchor = new MerkleAnchor();
    }

    function testAnchorAndGet() public {
        bytes32 root = keccak256("root");
        bytes32 cid = keccak256("cid");
        string memory rtype = "Observation";
        bytes32 patient = keccak256("patient");
        anchor.anchor(root, cid, rtype, patient);
        (bytes32 gotCid, string memory gotType, bytes32 gotPatient, uint64 ts, address actor) = anchor.get(root);
        assertEq(gotCid, cid, "cid");
        assertEq(keccak256(bytes(gotType)), keccak256(bytes(rtype)), "type");
        assertEq(gotPatient, patient, "patient");
        assertGt(ts, 0, "timestamp");
        assertEq(actor, address(this), "actor");
    }
}
