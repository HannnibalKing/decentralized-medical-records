// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title MerkleAnchor
/// @notice Anchors Merkle roots and associated CIDs for integrity verification.
contract MerkleAnchor {
    event Anchored(
        bytes32 indexed root,
        bytes32 indexed cid,
        string resourceType,
        bytes32 indexed patient,
        address actor,
        uint256 ts
    );

    struct Anchor {
        bytes32 root;
        bytes32 cid;
        string resourceType;
        bytes32 patient;
        uint64 ts;
        address actor;
    }

    mapping(bytes32 => Anchor) private anchors; // key: root

    /// @notice Anchor a Merkle root for a CID and metadata.
    /// @dev Caller should ensure uniqueness or accept overwrite semantics.
    function anchor(
        bytes32 root,
        bytes32 cid,
        string calldata resourceType,
        bytes32 patient
    ) external {
        anchors[root] = Anchor({
            root: root,
            cid: cid,
            resourceType: resourceType,
            patient: patient,
            ts: uint64(block.timestamp),
            actor: msg.sender
        });
        emit Anchored(root, cid, resourceType, patient, msg.sender, block.timestamp);
    }

    function get(bytes32 root)
        external
        view
        returns (bytes32 cid, string memory resourceType, bytes32 patient, uint64 ts, address actor)
    {
        Anchor storage a = anchors[root];
        return (a.cid, a.resourceType, a.patient, a.ts, a.actor);
    }
}
