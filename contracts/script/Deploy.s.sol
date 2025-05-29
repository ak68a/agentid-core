// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../src/AgentRegistry.sol";
import "../src/AgentDelegation.sol";

contract DeployScript is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        // Deploy AgentRegistry
        AgentRegistry registry = new AgentRegistry();
        console.log("AgentRegistry deployed to:", address(registry));

        // Deploy AgentDelegation with registry address
        AgentDelegation delegation = new AgentDelegation(address(registry));
        console.log("AgentDelegation deployed to:", address(delegation));

        // Add initial verifiers if needed
        address[] memory initialVerifiers = vm.envAddress("INITIAL_VERIFIERS", ",");
        for (uint i = 0; i < initialVerifiers.length; i++) {
            registry.addVerifier(initialVerifiers[i]);
            console.log("Added verifier:", initialVerifiers[i]);
        }

        vm.stopBroadcast();
    }
} 