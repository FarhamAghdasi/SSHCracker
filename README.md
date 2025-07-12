# 🚀 SSHCracker v2.5 - Advanced SSH Brute Force Tool

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-lightgrey?style=for-the-badge)](https://github.com/)
[![Release](https://img.shields.io/github/v/release/Matrix-Community-ORG/SSHCracker?style=for-the-badge)](https://github.com/Matrix-Community-ORG/SSHCracker/releases)
[![Version](https://img.shields.io/badge/Version-2.5-blue?style=for-the-badge)](https://github.com/Matrix-Community-ORG/SSHCracker)

A powerful, high-performance SSH brute force tool written in Go with **enhanced multi-layer worker architecture**, advanced honeypot detection, real-time statistics, and comprehensive system reconnaissance capabilities.

## 🌟 Key Features

### 🔥 Core Capabilities
- **⚡ Enhanced Multi-Layer Workers** - Revolutionary concurrent processing architecture
- **🚀 10x Performance Boost** - Up to 1000+ concurrent connections per worker
- **🍯 Advanced Honeypot Detection** - 9 intelligent detection algorithms with dedicated workers
- **📊 Real-Time Dashboard** - Live progress tracking with enhanced statistics
- **🎯 Smart Target Management** - Efficient wordlist and target handling
- **🔍 Deep System Reconnaissance** - Comprehensive server information gathering
- **📁 Beautiful Output Formats** - Enhanced logging with emojis and structured data

### 🛡️ Security & Performance
- **🚀 Cross-Platform Support** - Linux, Windows, macOS compatibility
- **⚙️ Intelligent Worker Management** - Separate pools for SSH and honeypot detection
- **🔒 Thread-Safe Operations** - Atomic operations for high-concurrency environments
- **🎛️ Advanced Configuration** - Timeout, stealth mode, performance tuning

## 🆕 What's New in v2.5

### 🎯 Enhanced Worker Architecture
- **Multi-Layer Concurrency**: Each worker handles 25+ concurrent connections
- **Dedicated Honeypot Workers**: 3 specialized workers for deep analysis
- **Thread-Safe Statistics**: Atomic operations for reliable counting
- **Structured Task Management**: Efficient memory usage and performance

### 🎨 Improved User Experience
- **Beautiful Console Output**: Enhanced real-time dashboard
- **Detailed Result Logging**: Emoji-rich structured output
- **Better Error Handling**: Graceful failure recovery
- **Performance Metrics**: Advanced speed and success rate tracking

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

## 🚀 Enhanced Multi-Layer Worker Architecture

### 🎯 Revolutionary Performance System
SSHCracker v2.5 introduces a groundbreaking multi-layer worker architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    Main Worker Pool                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Worker 1  │  │   Worker 2  │  │   Worker N  │        │
│  │             │  │             │  │             │        │
│  │ 25 Concurrent│  │ 25 Concurrent│  │ 25 Concurrent│        │
│  │ Connections │  │ Connections │  │ Connections │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                Honeypot Detection Pool                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ HoneyPot    │  │ HoneyPot    │  │ HoneyPot    │        │
│  │ Worker 1    │  │ Worker 2    │  │ Worker 3    │        │
│  │             │  │             │  │             │        │
│  │ Deep        │  │ Deep        │  │ Deep        │        │
│  │ Analysis    │  │ Analysis    │  │ Analysis    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

### 📊 Performance Metrics
- **Total Concurrent Capacity**: Workers × 25 connections
- **Example**: 50 workers × 25 = **1,250 concurrent connections**
- **Honeypot Processing**: 3 dedicated workers for deep analysis
- **Speed Improvement**: **5-10x faster** than previous versions

## 🍯 Advanced Honeypot Detection System

Our enhanced honeypot detection uses 9 sophisticated algorithms with dedicated workers:

| Algorithm | Detection Method | Performance Impact |
|-----------|------------------|-------------------|
| **Pattern Recognition** | Known honeypot signatures | ⚡ Fast |
| **Response Time Analysis** | Unusual timing patterns | ⚡ Fast |
| **Command Behavior** | Abnormal system responses | 🔄 Medium |
| **File System Analysis** | Fake file structures | 🔄 Medium |
| **Network Configuration** | Suspicious port configs | 🔄 Medium |
| **Performance Testing** | System characteristics | 🐌 Slow |
| **Anomaly Detection** | Unusual behaviors | 🐌 Slow |
| **Service Analysis** | Running processes | 🔄 Medium |
| **Environment Analysis** | System variables | ⚡ Fast |

> **✅ Production Ready**: Enhanced accuracy with dedicated processing workers.

## 📊 Enhanced Output Files

| File | Description | Format |
|------|-------------|---------|
| `su-goods.txt` | Successfully cracked SSH credentials | Simple list |
| `detailed-results.txt` | 🎨 Beautiful structured results with emojis | Enhanced format |
| `honeypots.txt` | Detected honeypots with confidence scores | Detailed analysis |
| `combo.txt` | Generated credential combinations | Temporary file |

## ⚙️ Advanced Configuration

### Performance Modes

**🚀 Ultra-High Speed Mode (v2.5)**:
- Timeout: 2 seconds
- Max Connections: 100
- Concurrent per Worker: 25
- Total Capacity: 2,500 connections

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
