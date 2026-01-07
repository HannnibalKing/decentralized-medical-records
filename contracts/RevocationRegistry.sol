// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title RevocationRegistry
/// @notice On-chain handle revocation for capabilities; simple, gas-cheap registry.
contract RevocationRegistry {
    event Revoked(bytes32 indexed handle, address indexed revoker, uint256 ts, string reason);

    struct Revocation {
        bool revoked;
        uint64 ts;
        address revoker;
        string reason;
    }

    mapping(bytes32 => Revocation) private revocations;

    /// @notice Revoke a capability handle.
    /// @dev Idempotent: subsequent calls on an already revoked handle are no-ops except for emitting nothing.
    function revoke(bytes32 handle, string calldata reason) external {
        Revocation storage r = revocations[handle];
        if (r.revoked) return;
        r.revoked = true;
        r.ts = uint64(block.timestamp);
        r.revoker = msg.sender;
        r.reason = reason;
        emit Revoked(handle, msg.sender, block.timestamp, reason);
    }

    /// @return revoked True if handle is revoked.
    /// @return ts Unix timestamp when revocation happened.
    /// @return revoker Address that revoked.
    /// @return reason Human-readable reason string.
    function get(bytes32 handle)
        external
        view
        returns (bool revoked, uint64 ts, address revoker, string memory reason)
    {
        Revocation storage r = revocations[handle];
        return (r.revoked, r.ts, r.revoker, r.reason);
    }

    function isRevoked(bytes32 handle) external view returns (bool) {
        return revocations[handle].revoked;
    }
}
