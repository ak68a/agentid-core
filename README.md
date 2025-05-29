# agentid-core

<div align="center">

![AgentID Core](https://img.shields.io/badge/AgentID-Core-blue)
![Status](https://img.shields.io/badge/Status-Active%20Development-orange)
![License](https://img.shields.io/badge/License-Apache%202.0-green)

> 🚧 **Active Development Notice**  
> This project is currently under active development.  
> Core features are being implemented and tested.  
> The Solidity components are still in progress.  
> Not ready for production use.

</div>

`agentid-core` is a Go + Solidity toolkit for building secure, verifiable agent identities based on the [ACK-ID specification](https://www.agentcommercekit.com/ack-id/introduction). It implements the key building blocks needed to support permissioned, autonomous agent activity across both offchain and onchain environments.

> For detailed documentation about the ACK-ID specification, trust framework, and non-technical overview, please see our [documentation](./docs/README.md).

## Features

- 🔐 Agent keypair generation (ACK-compatible, secp256k1)
- ✍️ Structured identity claim signing (EIP-712-ready)
- 🧾 Scoped delegation of authority to other agents
- 🧠 Solidity verifier for onchain identity and delegation checks

This library supports the full ACK-ID lifecycle: identity creation, verification, delegation, and enforcement. It is designed to be composable with other Agent Commerce Kit protocols such as [ACK-Pay](https://www.agentcommercekit.com/ack-pay/introduction).

## Enterprise Version

We also offer an enterprise version of AgentID designed for traditional infrastructure environments. This version includes:

- 🔐 Ed25519-based key management
- 📜 X.509 certificate integration
- 🏢 Enterprise PKI support
- 🔒 Centralized verification service
- 💾 Traditional database storage
- 📊 Enterprise-grade monitoring and logging
- 🛡️ Advanced security features

The enterprise version is ideal for organizations that need:
- Traditional infrastructure integration
- Enterprise security compliance
- PKI-based identity management
- Centralized control and monitoring
- Enterprise support and SLAs

For inquiries about the enterprise version, please contact us at [hey@ak68a.co](mailto:hey@ak68a.co).

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

## Security

AgentID Core is designed for security-critical applications. We take security seriously and have implemented several measures to ensure the safety of our users:

### Security Features

- 🔒 Cryptographically secure key generation and management
- 🛡️ EIP-712 structured data signing for human-readable signatures
- 🔐 Hardware Security Module (HSM) support
- 🎯 Capability-based access control
- 🔄 Automatic key rotation support
- 🚨 Comprehensive security monitoring and logging

### Reporting Security Issues

We take the security of AgentID Core seriously. If you believe you have found a security vulnerability, please report it to us as described in our [Security Policy](./SECURITY.md).

**Please do not report security vulnerabilities through public GitHub issues.**

### Security Best Practices

When using AgentID Core, we recommend following these security best practices:

1. **Key Management**
   - Use HSMs for key storage
   - Implement regular key rotation
   - Use secure key backup procedures
   - Implement proper access controls

2. **Deployment**
   - Keep all dependencies up to date
   - Use secure communication channels
   - Implement proper access controls
   - Monitor for suspicious activity

For more detailed security information, please see our [Security & Privacy documentation](./docs/security-privacy.md).

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](./LICENSE) file for details.

The Apache License 2.0 provides:
- A permissive license that allows for commercial use
- Patent protection for contributors and users
- Clear terms for modification and distribution
- A strong community-oriented license

## Contributing

We welcome contributions to AgentID Core! Please see our [Contributing Guidelines](./CONTRIBUTING.md) for more details.

### Development Process

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code of Conduct

Please note that this project is released with a [Contributor Code of Conduct](./CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Support

For support, please:
1. Check our [documentation](./docs/README.md)
2. Search [existing issues](https://github.com/ak68a/agentid-core/issues)
3. Create a new issue if needed

## Roadmap

See our [Roadmap](./docs/ROADMAP.md) for planned features and improvements.

## Acknowledgments

- [Agent Commerce Kit](https://www.agentcommercekit.com) for the ACK-ID specification
- [Ethereum Foundation](https://ethereum.org) for EIP-712
- All our contributors and users
