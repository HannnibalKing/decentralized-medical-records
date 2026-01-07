// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title AttestationRegistry
/// @notice Tracks prescription attestation state (issue, dispense, approve) with optional aggregate signature hash.
contract AttestationRegistry {
    enum Status {
        Issued,
        Dispensed,
        Approved,
        Void
    }

    event Recorded(bytes32 indexed cidMr, bytes32 cidMd, bytes32 aggSig, Status status, address indexed actor);
    event StatusUpdated(bytes32 indexed cidMr, Status status, address indexed actor);

    struct Attestation {
        bytes32 cidMr; // MedicationRequest CID
        bytes32 cidMd; // MedicationDispense CID
        bytes32 aggSig; // aggregate sig or hash of sig set
        Status status;
        uint64 ts;
    }

    mapping(bytes32 => Attestation) private attestations; // key: cidMr

    /// @notice Record a new attestation set.
    /// @dev Overwrites existing entry for cidMr; callers should ensure correctness.
    function record(
        bytes32 cidMr,
        bytes32 cidMd,
        bytes32 aggSig,
        Status status
    ) external {
        Attestation storage a = attestations[cidMr];
        a.cidMr = cidMr;
        a.cidMd = cidMd;
        a.aggSig = aggSig;
        a.status = status;
        a.ts = uint64(block.timestamp);
        emit Recorded(cidMr, cidMd, aggSig, status, msg.sender);
    }

    /// @notice Update status (e.g., void).
    function updateStatus(bytes32 cidMr, Status status) external {
        Attestation storage a = attestations[cidMr];
        a.status = status;
        a.ts = uint64(block.timestamp);
        emit StatusUpdated(cidMr, status, msg.sender);
    }

    function get(bytes32 cidMr)
        external
        view
        returns (bytes32 cidMd, bytes32 aggSig, Status status, uint64 ts)
    {
        Attestation storage a = attestations[cidMr];
        return (a.cidMd, a.aggSig, a.status, a.ts);
    }
}
