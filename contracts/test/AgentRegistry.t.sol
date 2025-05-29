// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "../src/AgentRegistry.sol";

contract AgentRegistryTest is Test {
    AgentRegistry public registry;
    address public owner;
    address public verifier;
    address public agent;
    string public testDID;

    function setUp() public {
        owner = address(this);
        verifier = address(0x1);
        agent = address(0x2);
        testDID = "did:ackid:0x1234";

        registry = new AgentRegistry();
        registry.addVerifier(verifier);
    }

    function test_RegisterAgent() public {
        vm.prank(agent);
        registry.registerAgent(testDID);

        (address agentAddress, string memory did, AgentRegistry.TrustLevel trustLevel, uint256 registeredAt, uint256 lastVerifiedAt, bool isActive) = registry.agents(agent);
        
        assertEq(agentAddress, agent);
        assertEq(did, testDID);
        assertEq(uint256(trustLevel), uint256(AgentRegistry.TrustLevel.Unverified));
        assertTrue(isActive);
        assertEq(registry.didToAddress(testDID), agent);
    }

    function test_VerifyAgent() public {
        // First register the agent
        vm.prank(agent);
        registry.registerAgent(testDID);

        // Then verify the agent
        bytes32[] memory capabilities = new bytes32[](2);
        capabilities[0] = keccak256("TRADE");
        capabilities[1] = keccak256("TRANSFER");

        vm.prank(verifier);
        registry.verifyAgent(agent, AgentRegistry.TrustLevel.Verified, capabilities);

        (,, AgentRegistry.TrustLevel trustLevel,,,) = registry.agents(agent);
        assertEq(uint256(trustLevel), uint256(AgentRegistry.TrustLevel.Verified));
        assertTrue(registry.hasCapability(agent, capabilities[0]));
        assertTrue(registry.hasCapability(agent, capabilities[1]));
    }

    function test_RevokeAgent() public {
        // First register and verify the agent
        vm.prank(agent);
        registry.registerAgent(testDID);

        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        vm.prank(verifier);
        registry.verifyAgent(agent, AgentRegistry.TrustLevel.Verified, capabilities);

        // Then revoke the agent
        vm.prank(verifier);
        registry.revokeAgent(agent, "Test revocation");

        (,,,,, bool isActive) = registry.agents(agent);
        assertFalse(isActive);
        assertEq(registry.didToAddress(testDID), address(0));
    }

    function test_OnlyVerifierCanVerify() public {
        vm.prank(agent);
        registry.registerAgent(testDID);

        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        vm.prank(agent);
        vm.expectRevert("Not a verifier");
        registry.verifyAgent(agent, AgentRegistry.TrustLevel.Verified, capabilities);
    }

    function test_CannotRegisterDuplicateDID() public {
        vm.prank(agent);
        registry.registerAgent(testDID);

        vm.prank(address(0x3));
        vm.expectRevert("DID already registered");
        registry.registerAgent(testDID);
    }

    function test_CannotRegisterDuplicateAgent() public {
        vm.prank(agent);
        registry.registerAgent(testDID);

        vm.prank(agent);
        vm.expectRevert("Agent already registered");
        registry.registerAgent("did:ackid:0x5678");
    }
} 