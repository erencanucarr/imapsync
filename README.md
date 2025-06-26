# 📧 IMAPSYNC CLI Tool

Professional command-line interface for seamless email migration using IMAPSYNC. Built with Go for maximum performance, cross-platform compatibility, and zero dependencies.

![GitHub release](https://img.shields.io/github/v/release/erencanucarr/imapsync)
![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)


![Termius_OJbvTQ5Po1](https://github.com/user-attachments/assets/537291b3-1467-4e5f-8947-c3b4350669bc)
![Termius_gGlvCQmLz1](https://github.com/user-attachments/assets/8f0c368c-be9e-4fad-ad1a-9a0d2178b010)




## ✨ Features

- 🎨 **Beautiful ASCII Interface** - Eye-catching terminal UI with colored output
- 🐍 **Python Detection & Auto-Installation** - Automatic Python environment setup
- 📧 **IMAPSYNC Management** - Complete IMAPSYNC installation and configuration
- 🔒 **Secure Password Input** - Hidden password entry for maximum security
- 🌍 **Multi-Platform Support** - Works on all major Linux distributions, macOS, and Windows
- ⚡ **Single Binary** - No dependencies required after compilation
- 🎯 **Developer Info Section** - Easy access to contact and project information
- 🚀 **Cross-Platform Build** - Compile for any target platform

## 🖥️ System Requirements

### Supported Operating Systems
- **Linux**: Ubuntu, Debian, CentOS, RHEL, Rocky Linux, AlmaLinux, Fedora, openSUSE, Arch Linux, Alpine Linux
- **macOS**: All versions with Homebrew support
- **Windows**: Windows 10/11 (64-bit)

### Prerequisites
- **For building from source**: Go 1.19 or higher
- **For running**: No additional requirements (self-contained binary)
- **Internet connection**: Required for IMAPSYNC installation

## 🚀 Quick Start

### Option 1: Download Pre-compiled Binary (Recommended)

```bash
# For Linux (64-bit)
wget https://github.com/erencanucarr/imapsync/releases/latest/download/imapsync-tool-linux
chmod +x imapsync-tool-linux
./imapsync-tool-linux

# For macOS (64-bit)
wget https://github.com/erencanucarr/imapsync/releases/latest/download/imapsync-tool-mac
chmod +x imapsync-tool-mac
./imapsync-tool-mac

# For Windows (64-bit)
# Download imapsync-tool.exe from releases page and run it
```

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/erencanucarr/imapsync.git
cd imapsync

# Build the binary
go build -o imapsync-tool main.go

# Make it executable and run
chmod +x imapsync-tool
./imapsync-tool
```

### Option 3: One-Line Installation

```bash
# Automatic installation script (Linux/macOS)
curl -sSL https://raw.githubusercontent.com/erencanucarr/imapsync/main/install.sh | bash
```

## 📋 Detailed Installation Guide

### Step 1: Install Go (if building from source)

#### Linux (Ubuntu/Debian):
```bash
sudo apt update
sudo apt install golang-go
```

#### Linux (CentOS/RHEL/Rocky):
```bash
sudo yum install golang
# or for newer versions
sudo dnf install golang
```

#### macOS:
```bash
brew install go
```

#### Windows:
Download from [https://golang.org/dl/](https://golang.org/dl/)

### Step 2: Clone and Build

```bash
# Clone the repository
git clone https://github.com/erencanucarr/imapsync.git
cd imapsync

# Download dependencies
go mod tidy

# Build for your platform
go build -o imapsync-tool main.go

# Or build for specific platforms
# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o imapsync-tool-linux main.go

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o imapsync-tool.exe main.go

# macOS 64-bit
GOOS=darwin GOARCH=amd64 go build -o imapsync-tool-mac main.go
```

### Step 3: Run the Tool

```bash
./imapsync-tool
```

## 🎯 Usage Guide

### 1. System Installation

When you first run the tool, select **"1 - Install System"**:

- **System Detection**: Automatically detects your OS and distribution
- **Python Check**: Verifies Python 3 installation, installs if missing
- **IMAPSYNC Check**: Verifies IMAPSYNC installation, installs if missing
- **Package Manager Support**: Works with apt, yum, dnf, zypper, pacman, apk, brew

#### Supported Package Managers:
| Distribution | Package Manager | Auto-Install |
|-------------|----------------|--------------|
| Ubuntu/Debian | `apt` | ✅ |
| CentOS/RHEL/Rocky/Alma | `yum` + EPEL | ✅ |
| Fedora | `dnf` | ✅ |
| openSUSE/SLES | `zypper` | ✅ |
| Arch/Manjaro | `pacman` | ✅ |
| Alpine | `apk` | ✅ |
| macOS | `brew` | ✅ |

### 2. Email Migration

Select **"2 - Migrate Emails"** and follow the wizard:

1. **Source mail server IP address** (e.g., `mail.oldserver.com`)
2. **Destination mail server IP address** (e.g., `mail.newserver.com`)
3. **Source email address** (e.g., `user@oldserver.com`)
4. **Source email password** (hidden input)
5. **Destination email address** (e.g., `user@newserver.com`)
6. **Destination email password** (hidden input)
7. **Confirmation** of migration details

#### Migration Features:
- **SSL/TLS Support**: Automatic SSL encryption for both servers
- **Folder Exclusions**: Automatically excludes Junk, Deleted, Trash folders
- **Folder Mapping**: Maps common folders (Sent → Sent Items, Spam → Junk E-Mail)
- **UID Preservation**: Maintains unique message identifiers
- **Caching**: Uses intelligent caching for faster subsequent runs
- **Progress Monitoring**: Real-time migration progress display

### 3. Developer Information

Select **"3 - Developer Info"** to access:
- Developer email contact
- LinkedIn profile
- GitHub repository

## ⚙️ IMAPSYNC Configuration

The tool uses optimized IMAPSYNC settings for reliable migration:

```bash
imapsync \
  --host1 'source-server' --ssl1 \
  --user1 'source-email' --password1 'source-password' \
  --host2 'destination-server' --ssl2 \
  --user2 'destination-email' --password2 'destination-password' \
  --exclude '^Junk\ E-Mail' \
  --exclude '^Deleted\ Items' \
  --exclude '^Deleted' \
  --exclude '^Trash' \
  --regextrans2 's#^Sent$#Sent Items#' \
  --regextrans2 's#^Spam$#Junk E-Mail#' \
  --useuid --usecache \
  --tmpdir /tmp/imapsync
```

## 🛠️ Advanced Usage

### Cross-Platform Compilation

```bash
# Compile for all major platforms
make build-all

# Or manually:
# Linux ARM64 (Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o imapsync-tool-arm64 main.go

# Linux 32-bit
GOOS=linux GOARCH=386 go build -o imapsync-tool-386 main.go

# Windows 32-bit
GOOS=windows GOARCH=386 go build -o imapsync-tool-386.exe main.go
```

### Installation to System PATH

```bash
# Install to system (requires sudo)
sudo cp imapsync-tool /usr/local/bin/
sudo chmod +x /usr/local/bin/imapsync-tool

# Now you can run from anywhere
imapsync-tool
```

### Uninstall

```bash
# Remove from system
sudo rm /usr/local/bin/imapsync-tool
```

## 🔧 Troubleshooting

### Common Issues and Solutions

#### Permission Denied Error
```bash
chmod +x imapsync-tool
```

#### Go Module Errors
```bash
go mod tidy
go mod download
```

#### EPEL Installation Failed (RHEL/CentOS/Rocky)
The tool automatically tries alternative installation methods including source compilation.

#### Python Not Found
The tool will automatically install Python 3 for your distribution.

#### IMAPSYNC Not Working
1. Verify installation: `which imapsync`
2. Check permissions: `ls -la /usr/local/bin/imapsync`
3. Try manual installation: `sudo yum install imapsync`

### Debug Mode

```bash
# Enable verbose output
./imapsync-tool --verbose
```

### Log Files

Check system logs for detailed error information:
```bash
# Linux
journalctl -u imapsync
tail -f /var/log/messages

# Check tmp directory
ls -la /tmp/imapsync/
```

## 📊 Performance Information

- **Binary Size**: ~8MB (statically linked, no dependencies)
- **Memory Usage**: <50MB during operation
- **Startup Time**: <100ms
- **Migration Speed**: Depends on network and server performance
- **Concurrent Connections**: Optimized for email server limits

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork** the repository
2. **Create** your feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Setup

```bash
git clone https://github.com/erencanucarr/imapsync.git
cd imapsync
go mod tidy
go run main.go
```

### Building and Testing

```bash
# Run tests
go test -v ./...

# Format code
go fmt ./...

# Build and test
go build -o imapsync-tool main.go
./imapsync-tool
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or need help:

1. **Check** the [Issues](https://github.com/erencanucarr/imapsync/issues) page
2. **Search** for existing solutions
3. **Create** a new issue with detailed information
4. **Contact** the developer directly

## 👨‍💻 Developer

**Eren Can Ucar**
- 📧 **Email**: [dev@eren.gg](mailto:dev@eren.gg)
- 💼 **LinkedIn**: [linkedin.com/in/erencanucarr](https://www.linkedin.com/in/erencanucarr/)
- 🐙 **GitHub**: [github.com/erencanucarr](https://github.com/erencanucarr)

## ⭐ Show Your Support

If this tool helped you migrate your emails successfully, please give it a ⭐ on GitHub!

## 📈 Roadmap

- [ ] **GUI Version** - Cross-platform desktop application
- [ ] **Batch Migration** - Multiple account migration support
- [ ] **Configuration Files** - Save and reuse migration settings
- [ ] **Progress Bar** - Enhanced visual progress tracking
- [ ] **Email Filtering** - Advanced filtering options
- [ ] **Cloud Storage** - Direct migration to cloud email providers
- [ ] **Docker Support** - Containerized deployment
- [ ] **API Integration** - REST API for automation

## 🙏 Acknowledgments

- **IMAPSYNC** - The powerful email migration tool that powers this CLI
- **Go Community** - For the excellent programming language and ecosystem
- **Contributors** - Everyone who helps improve this project

---

<div align="center">

**Made with ❤️ by [Eren Can Ucar](https://github.com/erencanucarr)**

[⭐ Star this repo](https://github.com/erencanucarr/imapsync) • [🐛 Report Bug](https://github.com/erencanucarr/imapsync/issues) • [💡 Request Feature](https://github.com/erencanucarr/imapsync/issues)

</div>
