// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Test.sol";
import "../src/AgentRegistry.sol";
import "../src/AgentDelegation.sol";

contract AgentDelegationTest is Test {
    AgentRegistry public registry;
    AgentDelegation public delegation;
    address public owner;
    address public verifier;
    address public delegator;
    address public delegate;
    string public testDID;

    function setUp() public {
        owner = address(this);
        verifier = address(0x1);
        delegator = address(0x2);
        delegate = address(0x3);
        testDID = "did:ackid:0x1234";

        // Deploy and setup registry
        registry = new AgentRegistry();
        registry.addVerifier(verifier);

        // Register and verify both agents
        vm.prank(delegator);
        registry.registerAgent("did:ackid:0x2");
        vm.prank(delegate);
        registry.registerAgent("did:ackid:0x3");

        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        vm.prank(verifier);
        registry.verifyAgent(delegator, AgentRegistry.TrustLevel.Verified, capabilities);
        vm.prank(verifier);
        registry.verifyAgent(delegate, AgentRegistry.TrustLevel.Verified, capabilities);

        // Deploy delegation contract
        delegation = new AgentDelegation(address(registry));
    }

    function test_CreateDelegation() public {
        bytes32[] memory capabilities = new bytes32[](2);
        capabilities[0] = keccak256("TRADE");
        capabilities[1] = keccak256("TRANSFER");

        vm.prank(delegator);
        delegation.createDelegation(delegate, capabilities, block.timestamp + 1 days);

        (bytes32[] memory delegatedCapabilities, uint256 validFrom, uint256 validUntil, bool isActive) = 
            delegation.getDelegationDetails(delegator, delegate);

        assertEq(delegatedCapabilities.length, 2);
        assertEq(delegatedCapabilities[0], capabilities[0]);
        assertEq(delegatedCapabilities[1], capabilities[1]);
        assertEq(validFrom, block.timestamp);
        assertEq(validUntil, block.timestamp + 1 days);
        assertTrue(isActive);
    }

    function test_RevokeDelegation() public {
        // First create a delegation
        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        vm.prank(delegator);
        delegation.createDelegation(delegate, capabilities, block.timestamp + 1 days);

        // Then revoke it
        vm.prank(delegator);
        delegation.revokeDelegation(delegate, "Test revocation");

        (,,, bool isActive) = delegation.getDelegationDetails(delegator, delegate);
        assertFalse(isActive);
    }

    function test_HasDelegatedCapability() public {
        bytes32[] memory capabilities = new bytes32[](2);
        capabilities[0] = keccak256("TRADE");
        capabilities[1] = keccak256("TRANSFER");

        vm.prank(delegator);
        delegation.createDelegation(delegate, capabilities, block.timestamp + 1 days);

        assertTrue(delegation.hasDelegatedCapability(delegator, delegate, capabilities[0]));
        assertTrue(delegation.hasDelegatedCapability(delegator, delegate, capabilities[1]));
        assertFalse(delegation.hasDelegatedCapability(delegator, delegate, keccak256("OTHER")));
    }

    function test_DelegationExpiration() public {
        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        vm.prank(delegator);
        delegation.createDelegation(delegate, capabilities, block.timestamp + 1 days);

        // Fast forward past expiration
        vm.warp(block.timestamp + 2 days);

        assertFalse(delegation.hasDelegatedCapability(delegator, delegate, capabilities[0]));
    }

    function test_CannotDelegateToUnregisteredAgent() public {
        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        vm.prank(delegator);
        vm.expectRevert("Delegate not registered");
        delegation.createDelegation(address(0x4), capabilities, block.timestamp + 1 days);
    }

    function test_CannotCreateDuplicateDelegation() public {
        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        vm.prank(delegator);
        delegation.createDelegation(delegate, capabilities, block.timestamp + 1 days);

        vm.prank(delegator);
        vm.expectRevert("Delegation already exists");
        delegation.createDelegation(delegate, capabilities, block.timestamp + 1 days);
    }

    function test_GetDelegateList() public {
        bytes32[] memory capabilities = new bytes32[](1);
        capabilities[0] = keccak256("TRADE");

        // Create multiple delegations
        vm.prank(delegator);
        delegation.createDelegation(delegate, capabilities, block.timestamp + 1 days);

        address delegate2 = address(0x4);
        vm.prank(delegate2);
        registry.registerAgent("did:ackid:0x4");
        vm.prank(verifier);
        registry.verifyAgent(delegate2, AgentRegistry.TrustLevel.Verified, capabilities);

        vm.prank(delegator);
        delegation.createDelegation(delegate2, capabilities, block.timestamp + 1 days);

        address[] memory delegates = delegation.getDelegateList(delegator);
        assertEq(delegates.length, 2);
        assertTrue(delegates[0] == delegate || delegates[1] == delegate);
        assertTrue(delegates[0] == delegate2 || delegates[1] == delegate2);
    }
} 