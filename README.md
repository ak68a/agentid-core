# agentid-core

`agentid-core` is a Go + Solidity toolkit for building secure, verifiable agent identities based on the [ACK-ID specification](https://www.agentcommercekit.com/ack-id/introduction). It implements the key building blocks needed to support permissioned, autonomous agent activity across both offchain and onchain environments.

> For detailed documentation about the ACK-ID specification, trust framework, and non-technical overview, please see our [documentation](./docs/README.md).

## Features

- 🔐 Agent keypair generation (ACK-compatible, secp256k1)
- ✍️ Structured identity claim signing (EIP-712-ready)
- 🧾 Scoped delegation of authority to other agents
- 🧠 Solidity verifier for onchain identity and delegation checks

This library supports the full ACK-ID lifecycle: identity creation, verification, delegation, and enforcement. It is designed to be composable with other Agent Commerce Kit protocols such as [ACK-Pay](https://www.agentcommercekit.com/ack-pay/introduction).

## Installation

```bash
go get github.com/ak68a/agentid-core
```

## Project Structure

The project is organized into several key packages:

- `key/` - Core functionality for agent keypair generation and management
- `models/` - Data structures and types for identity claims and delegations
- `signer/` - EIP-712 compatible signing utilities for identity claims

## Usage

### Generating Agent Keys

```go
import "github.com/ak68a/agentid-core/key"

// Generate a new agent keypair
keypair, err := key.GenerateKeyPair()
if err != nil {
    // Handle error
}

// Access the public and private keys
publicKey := keypair.PublicKey
privateKey := keypair.PrivateKey
```

### Signing Identity Claims

```go
import (
    "github.com/ak68a/agentid-core/models"
    "github.com/ak68a/agentid-core/signer"
)

// Create and sign an identity claim
claim := models.NewIdentityClaim(...)
signature, err := signer.SignClaim(claim, privateKey)
if err != nil {
    // Handle error
}
```

## Development

### Requirements

- Go 1.24.3 or later
- Ethereum development tools (for Solidity contract development)

### Testing

Run the test suite:

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Security

This library is designed for security-critical applications. Please report any security issues to our security team.
