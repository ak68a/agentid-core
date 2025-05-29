# agentid-core

`agentid-core` is a Go + Solidity toolkit for building secure, verifiable agent identities based on the [ACK-ID specification](https://www.agentcommercekit.com/ack-id/introduction). It implements the key building blocks needed to support permissioned, autonomous agent activity across both offchain and onchain environments.

### It includes:
- 🔐 Agent keypair generation (ACK-compatible, secp256k1)
- ✍️ Structured identity claim signing (EIP-712-ready)
- 🧾 Scoped delegation of authority to other agents
- 🧠 Solidity verifier for onchain identity and delegation checks

This library supports the full ACK-ID lifecycle: identity creation, verification, delegation, and enforcement. It is designed to be composable with other Agent Commerce Kit protocols such as [ACK-Pay](https://www.agentcommercekit.com/ack-pay/introduction).
