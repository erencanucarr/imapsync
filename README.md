# IMAPSYNC CLI

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platforms-Linux%20%7C%20macOS%20%7C%20Windows-blue)
![Zero Dependency](https://img.shields.io/badge/Zero%20Dependency-100%25%20Go%20StdLib-brightgreen)

> A **zero-dependency** IMAP mailbox synchronization tool with modern Terminal User Interface (TUI). Built entirely with Go standard library - no external dependencies!

---

## ✨ Key Features

| Feature | Description |
|---------|-------------|
| 🎯 **Zero Dependencies** | Built entirely with Go standard library - no external packages |
| 🖥️ **Modern TUI** | Beautiful terminal interface with colors, progress bars, and real-time stats |
| 🔄 **Parallel Transfers** | Sync multiple mailboxes simultaneously with configurable concurrency |
| 📊 **Real-time Monitoring** | Live progress tracking, performance metrics, and transfer statistics |
| 🛡️ **Safe & Reliable** | Uses `--useuid` for idempotent transfers, resume interrupted syncs |
| 🚀 **Auto Setup** | Automatic imapsync installation for multiple Linux distributions |
| 📝 **Comprehensive Logging** | Detailed logs with history and performance tracking |

---

## 🚀 Quick Start

### 1. Build & Run

```bash
git clone https://github.com/erencanucarr/imapsync-cli.git
cd imapsync-cli

# Build the application
go build -o imapsync ./cmd/imapsync

# Run with modern TUI (default)
./imapsync

# Or run in CLI mode
./imapsync -cli
```

### 2. First Time Setup

The application will automatically start in TUI mode. Choose **1) System Setup** to install required dependencies:

```
┌─ IMAPSYNC CLI ──────────────────────────────────────────────┐
│                                                             │
│  🔧 1) System Setup                                         │
│  📧 2) Transfer Mail                                        │
│  📊 3) View Statistics                                      │
│  📝 4) View Logs/History                                    │
│  ℹ️  5) Developer Info                                      │
│  ❌ 6) Exit                                                  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

If Python or imapsync are missing, the setup will offer automatic installation:

```
Python ✗
imapsync ✗
Local install scripts available:
 - ubuntu  - debian  - centos  - arch  - darwin
Enter distribution key to run installer or press Enter to skip:
```

### 3. Transfer Emails

1. Select **2) Transfer Mail** from the main menu
2. Enter source and destination server details:
   - Host addresses
   - Email accounts
   - Passwords (hidden input)
3. The tool validates credentials with `imapsync --justlogin`
4. Watch real-time progress with beautiful progress bars
5. Cancel safely with `Ctrl+C` - transfers are resumable

---

## 🎯 Zero Dependency Architecture

This project replaces all external libraries with custom implementations:

| External Library | Our Implementation | Features |
|------------------|-------------------|----------|
| `github.com/schollz/progressbar` | `internal/app/progressbar.go` | Real-time progress bars with ETA |
| `github.com/patrickmn/go-cache` | `internal/app/cache.go` | Thread-safe cache with expiration |
| `github.com/sirupsen/logrus` | `internal/app/logger.go` | Structured logging with levels |
| `golang.org/x/sync/semaphore` | `internal/app/semaphore.go` | Concurrency control |
| `golang.org/x/term` | `internal/app/term.go` | Cross-platform terminal handling |

### Benefits:
- **No external dependencies** - 100% Go standard library
- **Full control** over all code and behavior
- **Better performance** with minimal overhead
- **Enhanced security** - no third-party code execution
- **Easy customization** - modify any component as needed

---

## 📂 Project Structure

```
imapsync/
├── cmd/
│   └── imapsync/
│       └── main.go              # Application entry point
├── internal/
│   ├── app/
│   │   ├── cache.go             # Custom cache implementation
│   │   ├── developer.go         # Developer information
│   │   ├── logger.go            # Custom logging system
│   │   ├── parallel.go          # Parallel transfer management
│   │   ├── performance.go       # Performance metrics
│   │   ├── progressbar.go       # Custom progress bars
│   │   ├── semaphore.go         # Concurrency control
│   │   ├── setup.go             # System setup logic
│   │   ├── simple_interface.go  # TUI application logic
│   │   ├── term.go              # Terminal input handling
│   │   └── transfer.go          # Mail transfer logic
│   └── ui/
│       └── console.go           # Color and UI helpers
├── install/                     # OS-specific install scripts
│   ├── ubuntu.txt
│   ├── debian.txt
│   ├── centos.txt
│   ├── archlinux.txt
│   └── darwin.txt
└── README.md                    # This file
```

---

## 🖥️ Terminal User Interface

The modern TUI provides an intuitive experience:

### Main Menu
- Clean, colorized interface
- Easy navigation with number keys
- Real-time status indicators

### Transfer Interface
- Live progress bars with ETA
- Transfer statistics (speed, success rate)
- Memory usage monitoring
- Cache performance metrics

### Logs & History
- View detailed transfer logs
- Performance history
- Error tracking and debugging

---

## ⚙️ Recommended imapsync Configuration

The wrapper uses production-tested defaults:

```bash
--ssl1 --ssl2 \
--exclude "^Junk E-Mail" --exclude "^Trash" --exclude "^Deleted( Items)?$" \
--regextrans2 's#^Sent$#Sent Items#' --regextrans2 's#^Spam$#Junk E-Mail#' \
--useuid --usecache --tmpdir ./tmp --syncinternaldates --progress
```

You can modify these flags in `internal/app/transfer.go`.

---

## 🔧 Development

### Building

```bash
# Build for current platform
go build -o imapsync ./cmd/imapsync

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o imapsync-linux ./cmd/imapsync
GOOS=darwin GOARCH=amd64 go build -o imapsync-macos ./cmd/imapsync
GOOS=windows GOARCH=amd64 go build -o imapsync-windows.exe ./cmd/imapsync
```

### Testing

```bash
go vet ./...
go test ./...
```

### Running in Development

```bash
# Run with TUI (default)
go run ./cmd/imapsync

# Run in CLI mode
go run ./cmd/imapsync -cli
```

---

## 🚀 Performance Features

### Parallel Processing
- **Connection Pooling**: Semaphore-based concurrency control
- **Cache System**: Successful transfers are cached
- **Memory Management**: Automatic memory optimization
- **Progress Tracking**: Real-time performance metrics

### Statistics & Monitoring
- Transfer success rates
- Average speed calculations
- Memory usage tracking
- Cache performance metrics
- Real-time progress updates

---

## 🛠️ System Requirements

- **Go**: 1.21 or higher
- **Python**: 3.6+ (for imapsync)
- **imapsync**: Will be installed automatically
- **Platforms**: Linux, macOS, Windows

---

## 📝 License

MIT © 2025 - Zero Dependency IMAPSYNC CLI

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow Go conventions (`go vet`, `golint`)
4. Submit a pull request

All contributions are welcome!

---

## 🎯 Roadmap

- [x] Zero dependency architecture
- [x] Modern TUI interface
- [x] Parallel transfer support
- [x] Real-time statistics
- [x] Comprehensive logging
- [ ] OAuth2 support (Gmail, Outlook 365)
- [ ] Configuration file support
- [ ] Web dashboard
- [ ] Advanced filtering options
- [ ] Backup and restore features 