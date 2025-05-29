// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";
import "./AgentRegistry.sol";

contract AgentDelegation is Ownable, Pausable {
    AgentRegistry public registry;

    struct Delegation {
        address delegator;
        address delegate;
        bytes32[] capabilities;
        uint256 validFrom;
        uint256 validUntil;
        bool isActive;
    }

    // Events
    event DelegationCreated(
        address indexed delegator,
        address indexed delegate,
        bytes32[] capabilities,
        uint256 validUntil
    );
    event DelegationRevoked(
        address indexed delegator,
        address indexed delegate,
        string reason
    );
    event DelegationExpired(
        address indexed delegator,
        address indexed delegate
    );

    // State
    mapping(address => mapping(address => Delegation)) public delegations;
    mapping(address => address[]) public delegateList; // List of delegates for each delegator

    constructor(address _registry) Ownable(msg.sender) {
        registry = AgentRegistry(_registry);
    }

    // Modifiers
    modifier onlyRegisteredAgent() {
        require(registry.getAgentTrustLevel(msg.sender) != AgentRegistry.TrustLevel.Unverified, "Not a registered agent");
        _;
    }

    modifier onlyActiveDelegation(address delegator, address delegate) {
        Delegation storage d = delegations[delegator][delegate];
        require(d.isActive, "Delegation not active");
        require(block.timestamp >= d.validFrom, "Delegation not started");
        require(block.timestamp <= d.validUntil, "Delegation expired");
        _;
    }

    // Functions
    function createDelegation(
        address delegate,
        bytes32[] calldata capabilities,
        uint256 validUntil
    ) external onlyRegisteredAgent whenNotPaused {
        require(delegate != address(0), "Invalid delegate");
        require(validUntil > block.timestamp, "Invalid validity period");
        require(registry.getAgentTrustLevel(delegate) != AgentRegistry.TrustLevel.Unverified, "Delegate not registered");

        // Check if delegation already exists
        require(!delegations[msg.sender][delegate].isActive, "Delegation already exists");

        // Create new delegation
        delegations[msg.sender][delegate] = Delegation({
            delegator: msg.sender,
            delegate: delegate,
            capabilities: capabilities,
            validFrom: block.timestamp,
            validUntil: validUntil,
            isActive: true
        });

        // Add to delegate list if not already present
        bool exists = false;
        for (uint i = 0; i < delegateList[msg.sender].length; i++) {
            if (delegateList[msg.sender][i] == delegate) {
                exists = true;
                break;
            }
        }
        if (!exists) {
            delegateList[msg.sender].push(delegate);
        }

        emit DelegationCreated(msg.sender, delegate, capabilities, validUntil);
    }

    function revokeDelegation(address delegate, string memory reason) external {
        require(delegations[msg.sender][delegate].isActive, "No active delegation");
        
        delegations[msg.sender][delegate].isActive = false;
        
        // Remove from delegate list
        for (uint i = 0; i < delegateList[msg.sender].length; i++) {
            if (delegateList[msg.sender][i] == delegate) {
                delegateList[msg.sender][i] = delegateList[msg.sender][delegateList[msg.sender].length - 1];
                delegateList[msg.sender].pop();
                break;
            }
        }

        emit DelegationRevoked(msg.sender, delegate, reason);
    }

    function hasDelegatedCapability(
        address delegator,
        address delegate,
        bytes32 capability
    ) external view returns (bool) {
        Delegation storage d = delegations[delegator][delegate];
        if (!d.isActive || block.timestamp < d.validFrom || block.timestamp > d.validUntil) {
            return false;
        }

        for (uint i = 0; i < d.capabilities.length; i++) {
            if (d.capabilities[i] == capability) {
                return true;
            }
        }
        return false;
    }

    function getDelegationDetails(
        address delegator,
        address delegate
    ) external view returns (
        bytes32[] memory capabilities,
        uint256 validFrom,
        uint256 validUntil,
        bool isActive
    ) {
        Delegation storage d = delegations[delegator][delegate];
        return (d.capabilities, d.validFrom, d.validUntil, d.isActive);
    }

    function getDelegateList(address delegator) external view returns (address[] memory) {
        return delegateList[delegator];
    }

    // Admin functions
    function setRegistry(address _registry) external onlyOwner {
        registry = AgentRegistry(_registry);
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }
} 