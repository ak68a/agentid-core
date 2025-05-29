// contracts/src/AgentRegistry.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

contract AgentRegistry is Ownable, Pausable {
    constructor() Ownable(msg.sender) {}

    // Trust levels as defined in our framework
    enum TrustLevel {
        Unverified,    // Level 0
        Verified,      // Level 1
        Certified,     // Level 2
        CrossChain     // Level 3
    }

    struct Agent {
        address agentAddress;
        string did;
        TrustLevel trustLevel;
        uint256 registeredAt;
        uint256 lastVerifiedAt;
        bool isActive;
        mapping(bytes32 => bool) capabilities;
    }

    // Events
    event AgentRegistered(address indexed agent, string did, TrustLevel trustLevel);
    event AgentVerified(address indexed agent, TrustLevel newTrustLevel);
    event AgentRevoked(address indexed agent, string reason);
    event CapabilityAdded(address indexed agent, bytes32 capability);
    event CapabilityRemoved(address indexed agent, bytes32 capability);

    // State
    mapping(address => Agent) public agents;
    mapping(string => address) public didToAddress;
    mapping(address => bool) public verifiers;

    // Modifiers
    modifier onlyVerifier() {
        require(verifiers[msg.sender], "Not a verifier");
        _;
    }

    modifier onlyRegisteredAgent() {
        require(agents[msg.sender].isActive, "Not a registered agent");
        _;
    }

    // Functions
    function registerAgent(string memory did) external whenNotPaused {
        require(bytes(did).length > 0, "Invalid DID");
        require(didToAddress[did] == address(0), "DID already registered");
        require(!agents[msg.sender].isActive, "Agent already registered");

        Agent storage agent = agents[msg.sender];
        agent.agentAddress = msg.sender;
        agent.did = did;
        agent.trustLevel = TrustLevel.Unverified;
        agent.registeredAt = block.timestamp;
        agent.lastVerifiedAt = block.timestamp;
        agent.isActive = true;

        didToAddress[did] = msg.sender;
        emit AgentRegistered(msg.sender, did, TrustLevel.Unverified);
    }

    function verifyAgent(
        address agent,
        TrustLevel newTrustLevel,
        bytes32[] calldata capabilities
    ) external onlyVerifier whenNotPaused {
        require(agents[agent].isActive, "Agent not registered");
        require(uint256(newTrustLevel) > uint256(agents[agent].trustLevel), "Invalid trust level upgrade");

        agents[agent].trustLevel = newTrustLevel;
        agents[agent].lastVerifiedAt = block.timestamp;

        for (uint i = 0; i < capabilities.length; i++) {
            agents[agent].capabilities[capabilities[i]] = true;
            emit CapabilityAdded(agent, capabilities[i]);
        }

        emit AgentVerified(agent, newTrustLevel);
    }

    function revokeAgent(address agent, string memory reason) external onlyVerifier {
        require(agents[agent].isActive, "Agent not registered");
        
        agents[agent].isActive = false;
        delete didToAddress[agents[agent].did];
        
        emit AgentRevoked(agent, reason);
    }

    function hasCapability(address agent, bytes32 capability) external view returns (bool) {
        return agents[agent].isActive && agents[agent].capabilities[capability];
    }

    function getAgentTrustLevel(address agent) external view returns (TrustLevel) {
        return agents[agent].trustLevel;
    }

    // Admin functions
    function addVerifier(address verifier) external onlyOwner {
        verifiers[verifier] = true;
    }

    function removeVerifier(address verifier) external onlyOwner {
        verifiers[verifier] = false;
    }

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }
}