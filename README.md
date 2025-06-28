# IMAPSYNC CLI

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platforms-Linux%20%7C%20macOS%20%7C%20Windows-blue)

> A user-friendly, multi-language wrapper around [imapsync](https://imapsync.lamiral.info/) that makes mailbox migrations **safe, observable and cross-platform**.

---

## ✨ Key Highlights

| Category | Details |
|----------|---------|
| Interactivity | Guided CLI menus in English, Turkish, Spanish, German (extendable) |
| Security | Passwords collected with hidden input (`golang.org/x/term`) |
| Portability | Pre-built install scripts for Ubuntu, Debian, CentOS/RHEL, Arch, macOS |
| Observability | Realtime progress bar (`schollz/progressbar`) parses imapsync output |
| Extensibility | Modular Go code-base, JSON i18n, clean folder layout |

---

## 📂 Repository Layout

```
├── cmd/               # CLI entry point
│   └── imapsync/      # main.go
├── internal/
│   ├── app/           # business logic (setup, transfer, developer info)
│   ├── i18n/          # locale loader + JSON files (en, tr, es, de)
│   └── ui/            # colour helpers
├── install/           # OS-specific imapsync install scripts
├── README.md          # (EN) you are here
└── readmetr.md        # (TR) Türkçe sürüm
```

---

## 🚀 Quick Start

### 1. Build & Run

```bash
git clone https://github.com/yourname/imapsync-cli.git
cd imapsync-cli

go build -o imapsync ./cmd/imapsync
./imapsync -lang=en   # use -lang=<code> to switch language
```

### 2. Initial Setup

Choose **1) System Setup** in the menu. If Python / imapsync are missing you will see:

```
Python ✗
imapsync ✗
Local install scripts available in ./install directory:
 - ubuntu  - debian  - centos  - arch  - darwin
Enter distribution key to run installer or press Enter to skip:
```

Type your distribution keyword to run the automated installer **with sudo/root**. Logs are written to `/var/log/imapsync-install.log`.

### 3. Transfer E-mails

1. Pick **2) Transfer Mail**.
2. Enter source & destination hosts / e-mails / passwords (hidden).
3. Tool validates credentials with `imapsync --justlogin` then starts migration.
4. Watch the progress bar – you can safely cancel with `Ctrl+C` and rerun (thanks to `--useuid`).

---

## ⚙️ Recommended imapsync Flags

The wrapper launches imapsync with defaults proven in production:

```
--ssl1 --ssl2 \
--exclude "^Junk E-Mail" --exclude "^Trash" --exclude "^Deleted( Items)?$" \
--regextrans2 's#^Sent$#Sent Items#' --regextrans2 's#^Spam$#Junk E-Mail#' \
--useuid --usecache --tmpdir ./tmp --syncinternaldates --progress
```
You can modify flags inside `internal/app/transfer.go`.

---

## 🖥️ Installing imapsync Manually

If you prefer your own package manager, refer to `install/` scripts or the official docs:
<https://imapsync.lamiral.info/INSTALL.d/>.

---

## 🛠️ Development

```bash
go vet ./...
go test ./...
```

Run with `go run ./cmd/imapsync -lang=en` for rapid iterations.

### Adding a New Locale

1. Copy `internal/i18n/locales/en.json` ⇒ `internal/i18n/locales/fr.json` (example).
2. Translate the values.
3. Build & run with `-lang=fr`.

### Extending Install Scripts

Add `<distro>.txt` inside `install/` and reference the key in `internal/app/setup.go → scripts` map.

---

## 🙋‍♀️ FAQ

* **Does this store my passwords?**  No. They are passed directly to `imapsync` as process args.
* **Can I resume interrupted sync?**  Yes – `--useuid` makes runs idempotent.
* **GUI?**  Planned in roadmap – PRs welcome!

---

## 🤝 Contributing

1. Fork & branch.
2. Follow Go conventions (`go vet`, `golint`).
3. Submit a PR; CI will run `go test`.

All contributions, translations and bug reports are appreciated.

---

## 📅 Roadmap

- [ ] Config file / profile saving
- [ ] OAuth2 support (Gmail, Outlook 365)
- [ ] TUI with BubbleTea
- [ ] Build pipeline & release binaries

---

## 📝 License

MIT © 2025 Your Name
