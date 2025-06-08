# 🚀 SSHCracker - Advanced SSH Brute Force Tool

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey?style=for-the-badge)](https://github.com/)
[![Release](https://img.shields.io/github/v/release/Matrix-Community-ORG/SSHCracker?style=for-the-badge)](https://github.com/Matrix-Community-ORG/SSHCracker/releases)

A powerful, high-performance SSH brute force tool written in Go with advanced honeypot detection, real-time statistics, and comprehensive system reconnaissance capabilities.

## 🌟 Key Features

### 🔥 Core Capabilities
- **⚡ High-Speed SSH Brute Force** - Multi-threaded concurrent attacks
- **🍯 Advanced Honeypot Detection (BETA)** - 9 intelligent detection algorithms
- **📊 Real-Time Dashboard** - Live progress tracking and statistics
- **🎯 Smart Target Management** - Efficient wordlist and target handling
- **🔍 Deep System Reconnaissance** - Comprehensive server information gathering
- **📁 Multiple Output Formats** - Detailed logging and reporting

### 🛡️ Security & Performance
- **🚀 Cross-Platform Support** - Linux, Windows, macOS compatibility
- **⚙️ Configurable Threading** - Adjustable concurrent connections
- **🔒 No License Requirements** - Completely free and open-source
- **🎛️ Advanced Configuration** - Timeout, stealth mode, performance tuning

## 🚀 Quick Start

### Option 1: Download Pre-built Binary (Recommended)
```bash
# Visit releases page and download for your platform:
# https://github.com/Matrix-Community-ORG/SSHCracker/releases/latest

# Make executable (Linux/macOS):
chmod +x ssh-cracker-*

# Run:
./ssh-cracker-*
```

### Option 2: Build from Source
```bash
# Clone repository
git clone https://github.com/Matrix-Community-ORG/SSHCracker.git
cd SSHCracker

# Build
go build ssh.go

# Run
./ssh
```

## 📋 Usage Guide

### Basic Usage
1. **Launch the tool**: `./ssh-cracker-*`
2. **Configure inputs**:
   - Username wordlist file (e.g., `users.txt`)
   - Password wordlist file (e.g., `passwords.txt`)
   - Target list file (e.g., `targets.txt`)
   - Connection timeout (recommended: 5-10 seconds)
   - Max concurrent connections (recommended: 10-50)

### File Format Examples

**Usernames (`users.txt`)**:
```
root
admin
administrator
user
ubuntu
```

**Passwords (`passwords.txt`)**:
```
123456
password
admin
root
12345678
```

**Targets (`targets.txt`)**:
```
192.168.1.1:22
10.0.0.1:22
example.com:2222
```

## 🍯 Honeypot Detection System (BETA)

Our advanced honeypot detection uses 9 sophisticated algorithms:

| Algorithm | Detection Method |
|-----------|------------------|
| **Pattern Recognition** | Known honeypot signatures and artifacts |
| **Response Time Analysis** | Unusual timing patterns |
| **Command Behavior** | Abnormal system command responses |
| **File System Analysis** | Fake or simulated file structures |
| **Network Configuration** | Suspicious port and service configurations |
| **Performance Testing** | System performance characteristics |
| **Anomaly Detection** | Unusual system behaviors |
| **Service Analysis** | Running processes and services |
| **Environment Analysis** | System environment variables |

> **⚠️ Note**: Currently in BETA - manual verification recommended for detected honeypots.

## 📊 Output Files

| File | Description |
|------|-------------|
| `su-goods.txt` | Successfully cracked SSH credentials |
| `detailed-results.txt` | Comprehensive scan results with system info |
| `honeypots.txt` | Detected honeypots with confidence scores |
| `combo.txt` | Generated credential combinations (temporary) |

## ⚙️ Advanced Configuration

### Performance Modes

**🏃 High-Speed Mode**:
- Timeout: 3 seconds
- Max Connections: 50
- Use for fast networks

**🥷 Stealth Mode**:
- Timeout: 10 seconds
- Max Connections: 5
- Use for careful reconnaissance

## 🔧 Installation Requirements

### Prerequisites
- **Go**: Version 1.18 or higher
- **Git**: For cloning (if building from source)
- **Network Access**: To target systems

### Supported Platforms
- ✅ Linux (x64, ARM64)
- ✅ Windows (x64)
- ✅ macOS (Intel, Apple Silicon)

## 🛠️ Troubleshooting

### Common Issues
```bash
# Permission denied
chmod +x ssh-cracker-*

# Module errors
go mod download && go mod tidy

# Too many open files
ulimit -n 65536
```

### Performance Tips
- Adjust timeout based on network latency
- Start with lower connection counts
- Monitor system resources during scanning

## 📱 Community & Support

### 🌐 Join Our Communities
- **English Community**: [@MatrixORG](https://t.me/MatrixORG)
- **Persian Community**: [@MatrixFa](https://t.me/MatrixFa)
- **Chat Group**: [@DD0SChat](https://t.me/DD0SChat)

### 💬 Get Help
1. Check [Issues](https://github.com/Matrix-Community-ORG/SSHCracker/issues)
2. Join our Telegram communities
3. Create detailed bug reports

## ⚠️ Legal & Ethical Use

**🚨 Important Notice**: This tool is designed for:
- ✅ Authorized penetration testing
- ✅ Educational purposes
- ✅ Security research
- ✅ Your own systems

**❌ DO NOT USE FOR**:
- Unauthorized access attempts
- Illegal activities
- Systems you don't own without permission

**Users are fully responsible for compliance with applicable laws and regulations.**

## 🤝 Contributing

We welcome contributions! Here's how:

1. **Fork** the repository
2. **Create** your feature branch: `git checkout -b feature/AmazingFeature`
3. **Commit** your changes: `git commit -m 'Add AmazingFeature'`
4. **Push** to branch: `git push origin feature/AmazingFeature`
5. **Open** a Pull Request

## 📄 License

This project is licensed under the **MIT License** - see [LICENSE](LICENSE) for details.

## 🏆 Acknowledgments

- **Matrix Community** - Development and maintenance
- **Go Community** - Excellent SSH libraries
- **Security Researchers** - Honeypot detection algorithms
- **Contributors** - Bug reports and feature requests

---

<div align="center">

**⭐ Star this project if you find it useful! ⭐**

Made with ❤️ by [Matrix Community](https://github.com/Matrix-Community-ORG)

</div>
