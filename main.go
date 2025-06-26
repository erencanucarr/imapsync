package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

// Colors for terminal output
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBold    = "\033[1m"
)

// SystemInfo holds system information
type SystemInfo struct {
	OS             string
	Distribution   string
	PackageManager string
}

// MigrationInfo holds email migration details
type MigrationInfo struct {
	SourceHost string
	DestHost   string
	SourceUser string
	SourcePass string
	DestUser   string
	DestPass   string
}

// ImapSyncCLI main structure
type ImapSyncCLI struct {
	SystemInfo SystemInfo
}

// NewImapSyncCLI creates a new CLI instance
func NewImapSyncCLI() *ImapSyncCLI {
	cli := &ImapSyncCLI{}
	cli.detectSystem()
	return cli
}

// detectSystem detects the operating system and package manager
func (cli *ImapSyncCLI) detectSystem() {
	cli.SystemInfo.OS = runtime.GOOS

	if cli.SystemInfo.OS == "linux" {
		cli.detectLinuxDistribution()
	} else if cli.SystemInfo.OS == "darwin" {
		cli.SystemInfo.Distribution = "macOS"
		cli.SystemInfo.PackageManager = "brew"
	}
}

// detectLinuxDistribution detects Linux distribution and package manager
func (cli *ImapSyncCLI) detectLinuxDistribution() {
	// Try to read /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				cli.SystemInfo.Distribution = strings.Trim(strings.Split(line, "=")[1], "\"")
				break
			}
		}
	}

	// Determine package manager based on distribution
	switch cli.SystemInfo.Distribution {
	case "ubuntu", "debian", "linuxmint":
		cli.SystemInfo.PackageManager = "apt"
	case "centos", "rhel", "rocky", "almalinux", "ol":
		cli.SystemInfo.PackageManager = "yum"
	case "fedora":
		cli.SystemInfo.PackageManager = "dnf"
	case "opensuse", "sles":
		cli.SystemInfo.PackageManager = "zypper"
	case "arch", "manjaro":
		cli.SystemInfo.PackageManager = "pacman"
	case "alpine":
		cli.SystemInfo.PackageManager = "apk"
	default:
		// Try to detect by checking which package manager exists
		if cli.commandExists("apt") {
			cli.SystemInfo.PackageManager = "apt"
		} else if cli.commandExists("yum") {
			cli.SystemInfo.PackageManager = "yum"
		} else if cli.commandExists("dnf") {
			cli.SystemInfo.PackageManager = "dnf"
		} else if cli.commandExists("zypper") {
			cli.SystemInfo.PackageManager = "zypper"
		} else if cli.commandExists("pacman") {
			cli.SystemInfo.PackageManager = "pacman"
		} else if cli.commandExists("apk") {
			cli.SystemInfo.PackageManager = "apk"
		} else {
			cli.SystemInfo.PackageManager = "unknown"
		}
	}
}

// commandExists checks if a command exists
func (cli *ImapSyncCLI) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// printWelcome prints ASCII welcome message
func (cli *ImapSyncCLI) printWelcome() {
	welcomeASCII := fmt.Sprintf(`
%s%s
██╗███╗   ███╗ █████╗ ██████╗ ███████╗██╗   ██╗███╗   ██╗ ██████╗
██║████╗ ████║██╔══██╗██╔══██╗██╔════╝╚██╗ ██╔╝████╗  ██║██╔════╝
██║██╔████╔██║███████║██████╔╝███████╗ ╚████╔╝ ██╔██╗ ██║██║     
██║██║╚██╔╝██║██╔══██║██╔═══╝ ╚════██║  ╚██╔╝  ██║╚██╗██║██║     
██║██║ ╚═╝ ██║██║  ██║██║     ███████║   ██║   ██║ ╚████║╚██████╗
╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝     ╚══════╝   ╚═╝   ╚═╝  ╚═══╝ ╚═════╝
%s
%s%s                    Mail Migration Tool%s
%s%s                   Powered by IMAPSYNC%s
`, ColorBlue, ColorBold, ColorReset, ColorMagenta, ColorBold, ColorReset, ColorCyan, ColorBold, ColorReset)

	fmt.Println(welcomeASCII)
}

// showMenu displays main menu options
func (cli *ImapSyncCLI) showMenu() {
	fmt.Printf("\n%s%sPlease select an option:%s\n", ColorBold, ColorWhite, ColorReset)
	fmt.Printf("%s%s1 - Install System%s\n", ColorGreen, ColorBold, ColorReset)
	fmt.Printf("%s%s2 - Migrate Emails%s\n", ColorGreen, ColorBold, ColorReset)
	fmt.Printf("%s%s3 - Developer Info%s\n", ColorBlue, ColorBold, ColorReset)
	fmt.Printf("%s%sq - Exit%s\n", ColorYellow, ColorBold, ColorReset)
}

// checkPythonInstalled checks if Python is installed
func (cli *ImapSyncCLI) checkPythonInstalled() (bool, string, string) {
	// Check Python 3
	if cmd := exec.Command("python3", "--version"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			version := strings.TrimSpace(string(output))
			return true, version, "python3"
		}
	}

	// Check Python
	if cmd := exec.Command("python", "--version"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			version := strings.TrimSpace(string(output))
			if strings.Contains(version, "Python 3.") {
				return true, version, "python"
			}
			return false, version, "python2"
		}
	}

	return false, "", ""
}

// checkImapSyncInstalled checks if imapsync is installed
func (cli *ImapSyncCLI) checkImapSyncInstalled() bool {
	_, err := exec.LookPath("imapsync")
	return err == nil
}

// installPython installs Python based on the operating system
func (cli *ImapSyncCLI) installPython() bool {
	fmt.Printf("\n%s%sInstalling Python...%s\n", ColorYellow, ColorBold, ColorReset)

	var commands [][]string

	switch cli.SystemInfo.PackageManager {
	case "apt":
		commands = [][]string{
			{"sudo", "apt", "update"},
			{"sudo", "apt", "install", "-y", "python3", "python3-pip"},
		}
	case "yum":
		commands = [][]string{
			{"sudo", "yum", "install", "-y", "python3", "python3-pip"},
		}
	case "dnf":
		commands = [][]string{
			{"sudo", "dnf", "install", "-y", "python3", "python3-pip"},
		}
	case "zypper":
		commands = [][]string{
			{"sudo", "zypper", "install", "-y", "python3", "python3-pip"},
		}
	case "pacman":
		commands = [][]string{
			{"sudo", "pacman", "-Sy", "--noconfirm", "python", "python-pip"},
		}
	case "apk":
		commands = [][]string{
			{"sudo", "apk", "add", "python3", "py3-pip"},
		}
	case "brew":
		commands = [][]string{
			{"brew", "install", "python3"},
		}
	default:
		fmt.Printf("%s%sUnsupported package manager. Please install Python manually.%s\n", ColorRed, ColorBold, ColorReset)
		return false
	}

	for _, cmd := range commands {
		fmt.Printf("Running: %s\n", strings.Join(cmd, " "))
		execCmd := exec.Command(cmd[0], cmd[1:]...)
		if err := execCmd.Run(); err != nil {
			fmt.Printf("%s%sError: %v%s\n", ColorRed, ColorBold, err, ColorReset)
			return false
		}
	}

	fmt.Printf("%s%sPython installed successfully!%s\n", ColorGreen, ColorBold, ColorReset)
	return true
}

// installImapSync installs imapsync based on the operating system
func (cli *ImapSyncCLI) installImapSync() bool {
	fmt.Printf("\n%s%sInstalling IMAPSYNC...%s\n", ColorYellow, ColorBold, ColorReset)

	var commands [][]string

	switch cli.SystemInfo.PackageManager {
	case "apt":
		commands = [][]string{
			{"sudo", "apt", "update"},
			{"sudo", "apt", "install", "-y", "imapsync"},
		}
	case "yum":
		// For RHEL/CentOS/Rocky/Alma - need EPEL
		commands = [][]string{
			{"sudo", "yum", "install", "-y", "epel-release"},
			{"sudo", "yum", "install", "-y", "imapsync"},
		}
	case "dnf":
		// For Fedora - try EPEL first, then RPM Fusion
		commands = [][]string{
			{"sudo", "dnf", "install", "-y", "epel-release"},
			{"sudo", "dnf", "install", "-y", "imapsync"},
		}
	case "zypper":
		commands = [][]string{
			{"sudo", "zypper", "install", "-y", "imapsync"},
		}
	case "pacman":
		// Arch AUR - try with yay or manual compilation
		commands = [][]string{
			{"sudo", "pacman", "-Sy", "--noconfirm", "base-devel"},
			{"yay", "-S", "--noconfirm", "imapsync"},
		}
	case "apk":
		commands = [][]string{
			{"sudo", "apk", "add", "imapsync"},
		}
	case "brew":
		commands = [][]string{
			{"brew", "install", "imapsync"},
		}
	default:
		fmt.Printf("%s%sUnsupported package manager. Please install IMAPSYNC manually.%s\n", ColorRed, ColorBold, ColorReset)
		return false
	}

	for _, cmd := range commands {
		fmt.Printf("Running: %s\n", strings.Join(cmd, " "))
		execCmd := exec.Command(cmd[0], cmd[1:]...)
		if err := execCmd.Run(); err != nil {
			fmt.Printf("%s%sError: %v%s\n", ColorRed, ColorBold, err, ColorReset)

			// If EPEL installation fails, try alternative methods
			if strings.Contains(strings.Join(cmd, " "), "epel-release") {
				fmt.Printf("%s%sEPEL installation failed. Trying alternative method...%s\n", ColorYellow, ColorBold, ColorReset)
				return cli.installImapSyncAlternative()
			}
			return false
		}
	}

	fmt.Printf("%s%sIMAPSYNC installed successfully!%s\n", ColorGreen, ColorBold, ColorReset)
	return true
}

// installImapSyncAlternative tries alternative installation methods
func (cli *ImapSyncCLI) installImapSyncAlternative() bool {
	fmt.Printf("%s%sTrying alternative IMAPSYNC installation method...%s\n", ColorCyan, ColorBold, ColorReset)

	// Try to install from source
	commands := [][]string{
		{"sudo", "yum", "groupinstall", "-y", "Development Tools"},
		{"sudo", "yum", "install", "-y", "perl-CPAN", "perl-App-cpanminus", "openssl-devel"},
		{"sudo", "cpanm", "Mail::IMAPClient", "IO::Socket::SSL", "Digest::MD5", "Term::ReadKey", "File::Spec", "IO::Socket::INET6"},
	}

	for _, cmd := range commands {
		fmt.Printf("Running: %s\n", strings.Join(cmd, " "))
		execCmd := exec.Command(cmd[0], cmd[1:]...)
		if err := execCmd.Run(); err != nil {
			fmt.Printf("%s%sWarning: %v%s\n", ColorYellow, ColorBold, err, ColorReset)
		}
	}

	// Download and install imapsync from source
	sourceCommands := [][]string{
		{"wget", "-O", "/tmp/imapsync", "https://raw.githubusercontent.com/imapsync/imapsync/master/imapsync"},
		{"sudo", "chmod", "+x", "/tmp/imapsync"},
		{"sudo", "mv", "/tmp/imapsync", "/usr/local/bin/imapsync"},
	}

	for _, cmd := range sourceCommands {
		fmt.Printf("Running: %s\n", strings.Join(cmd, " "))
		execCmd := exec.Command(cmd[0], cmd[1:]...)
		if err := execCmd.Run(); err != nil {
			fmt.Printf("%s%sError: %v%s\n", ColorRed, ColorBold, err, ColorReset)
			return false
		}
	}

	fmt.Printf("%s%sIMAPSYNC installed successfully from source!%s\n", ColorGreen, ColorBold, ColorReset)
	return true
}

// installSystemRequirements installs Python and IMAPSYNC with comprehensive system checks
func (cli *ImapSyncCLI) installSystemRequirements() {
	fmt.Printf("\n%s%sSystem Analysis and Installation%s\n", ColorMagenta, ColorBold, ColorReset)
	fmt.Println(strings.Repeat("=", 50))

	// System Info
	fmt.Printf("\n%s%s📋 System Information:%s\n", ColorBold, ColorWhite, ColorReset)
	fmt.Printf("Operating System: %s%s%s\n", ColorCyan, strings.Title(cli.SystemInfo.OS), ColorReset)
	if cli.SystemInfo.Distribution != "" {
		fmt.Printf("Distribution: %s%s%s\n", ColorCyan, strings.Title(cli.SystemInfo.Distribution), ColorReset)
	}
	fmt.Printf("Package Manager: %s%s%s\n", ColorCyan, cli.SystemInfo.PackageManager, ColorReset)

	// Python Check
	fmt.Printf("\n%s%s🐍 Python Check:%s\n", ColorBold, ColorWhite, ColorReset)
	pythonInstalled, pythonVersion, pythonCmd := cli.checkPythonInstalled()

	if pythonInstalled {
		fmt.Printf("%s%s✓ Python installed: %s%s\n", ColorGreen, ColorBold, pythonVersion, ColorReset)
		if pythonCmd == "python2" {
			fmt.Printf("%s%s⚠️  Python 2 detected. Python 3 installation is recommended.%s\n", ColorYellow, ColorBold, ColorReset)
			if cli.getUserConfirmation("Do you want to install Python 3? (y/n): ") {
				if !cli.installPython() {
					return
				}
			}
		}
	} else {
		fmt.Printf("%s%s✗ Python not installed%s\n", ColorRed, ColorBold, ColorReset)
		fmt.Printf("%s%sPython 3 is required for IMAPSYNC CLI Tool to work.%s\n", ColorYellow, ColorBold, ColorReset)
		if cli.getUserConfirmation("Do you want to install Python 3? (y/n): ") {
			if !cli.installPython() {
				return
			}
		} else {
			fmt.Printf("%s%sPython installation cancelled.%s\n", ColorYellow, ColorBold, ColorReset)
			return
		}
	}

	// IMAPSYNC Check
	fmt.Printf("\n%s%s📧 IMAPSYNC Check:%s\n", ColorBold, ColorWhite, ColorReset)
	if cli.checkImapSyncInstalled() {
		fmt.Printf("%s%s✓ IMAPSYNC installed - No installation required%s\n", ColorGreen, ColorBold, ColorReset)
	} else {
		fmt.Printf("%s%s✗ IMAPSYNC not installed%s\n", ColorRed, ColorBold, ColorReset)
		if cli.getUserConfirmation("Do you want to install IMAPSYNC? (y/n): ") {
			if !cli.installImapSync() {
				return
			}
		} else {
			fmt.Printf("%s%sIMAPSYNC installation cancelled.%s\n", ColorYellow, ColorBold, ColorReset)
			return
		}
	}

	// Final Status
	fmt.Printf("\n%s%s📊 Installation Summary:%s\n", ColorMagenta, ColorBold, ColorReset)
	fmt.Println(strings.Repeat("=", 50))

	// Re-check everything
	pythonInstalled, pythonVersion, _ = cli.checkPythonInstalled()
	imapSyncInstalled := cli.checkImapSyncInstalled()

	pythonStatus := "✗ Not installed"
	if pythonInstalled {
		pythonStatus = fmt.Sprintf("✓ %s", pythonVersion)
	}

	imapSyncStatus := "✗ Not installed"
	if imapSyncInstalled {
		imapSyncStatus = "✓ Installed"
	}

	fmt.Printf("Python: %s%s%s%s\n",
		func() string {
			if pythonInstalled {
				return ColorGreen
			} else {
				return ColorRed
			}
		}(),
		ColorBold, pythonStatus, ColorReset)
	fmt.Printf("IMAPSYNC: %s%s%s%s\n",
		func() string {
			if imapSyncInstalled {
				return ColorGreen
			} else {
				return ColorRed
			}
		}(),
		ColorBold, imapSyncStatus, ColorReset)

	if pythonInstalled && imapSyncInstalled {
		fmt.Printf("\n%s%s🎉 All system requirements met! Ready for email migration.%s\n", ColorGreen, ColorBold, ColorReset)
	} else {
		fmt.Printf("\n%s%s⚠️  Some requirements are missing. Please install missing components.%s\n", ColorYellow, ColorBold, ColorReset)
	}
}

// getUserInput gets user input with prompt
func (cli *ImapSyncCLI) getUserInput(prompt string) string {
	fmt.Printf("%s%s%s%s", ColorCyan, ColorBold, prompt, ColorReset)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			input := strings.TrimSpace(scanner.Text())
			if input != "" {
				return input
			}
		}
		fmt.Printf("%s%sThis field cannot be empty. Please try again.%s\n", ColorYellow, ColorBold, ColorReset)
		fmt.Printf("%s%s%s%s", ColorCyan, ColorBold, prompt, ColorReset)
	}
}

// getUserPassword gets password input (hidden)
func (cli *ImapSyncCLI) getUserPassword(prompt string) string {
	fmt.Printf("%s%s%s%s", ColorCyan, ColorBold, prompt, ColorReset)
	for {
		password, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // New line after password input
		if err == nil && len(strings.TrimSpace(string(password))) > 0 {
			return strings.TrimSpace(string(password))
		}
		fmt.Printf("%s%sThis field cannot be empty. Please try again.%s\n", ColorYellow, ColorBold, ColorReset)
		fmt.Printf("%s%s%s%s", ColorCyan, ColorBold, prompt, ColorReset)
	}
}

// getUserConfirmation gets user confirmation
func (cli *ImapSyncCLI) getUserConfirmation(prompt string) bool {
	fmt.Printf("%s%s%s%s", ColorCyan, ColorBold, prompt, ColorReset)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			input := strings.ToLower(strings.TrimSpace(scanner.Text()))
			switch input {
			case "y", "yes":
				return true
			case "n", "no":
				return false
			default:
				fmt.Printf("%s%sPlease enter 'y' (yes) or 'n' (no).%s\n", ColorYellow, ColorBold, ColorReset)
				fmt.Printf("%s%s%s%s", ColorCyan, ColorBold, prompt, ColorReset)
			}
		}
	}
}

// collectMigrationInfo collects all necessary information for email migration
func (cli *ImapSyncCLI) collectMigrationInfo() *MigrationInfo {
	fmt.Printf("\n%s%sEmail Migration Wizard%s\n", ColorMagenta, ColorBold, ColorReset)
	fmt.Printf("%s%sPlease enter the following information step by step:%s\n\n", ColorBold, ColorWhite, ColorReset)

	info := &MigrationInfo{}

	// Collect information step by step
	info.SourceHost = cli.getUserInput("1. Enter current mail server IP address: ")
	info.DestHost = cli.getUserInput("2. Enter destination mail server IP address: ")
	info.SourceUser = cli.getUserInput("3. Enter source email address: ")
	info.SourcePass = cli.getUserPassword("4. Enter source email password: ")
	info.DestUser = cli.getUserInput("5. Enter destination email address: ")
	info.DestPass = cli.getUserPassword("6. Enter destination email password: ")

	// Confirmation
	fmt.Printf("\n%s%s7. All email content will be migrated. Do you confirm?%s\n", ColorYellow, ColorBold, ColorReset)
	fmt.Printf("%s%sSource:%s %s (%s)\n", ColorBold, ColorWhite, ColorReset, info.SourceUser, info.SourceHost)
	fmt.Printf("%s%sDestination:%s %s (%s)\n", ColorBold, ColorWhite, ColorReset, info.DestUser, info.DestHost)

	if !cli.getUserConfirmation("Do you confirm? (y/n): ") {
		fmt.Printf("%s%sOperation cancelled.%s\n", ColorYellow, ColorBold, ColorReset)
		return nil
	}

	return info
}

// runImapSync executes imapsync with the provided information
func (cli *ImapSyncCLI) runImapSync(info *MigrationInfo) bool {
	if !cli.checkImapSyncInstalled() {
		fmt.Printf("%s%sIMAPSYNC not installed. Please use '1 - Install System' option first.%s\n", ColorRed, ColorBold, ColorReset)
		return false
	}

	// Create tmp directory
	tmpDir := "/tmp/imapsync"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		fmt.Printf("%s%sCould not create temporary directory: %v%s\n", ColorRed, ColorBold, err, ColorReset)
		return false
	}

	// Build imapsync command
	args := []string{
		"--host1", info.SourceHost,
		"--ssl1",
		"--user1", info.SourceUser,
		"--password1", info.SourcePass,
		"--host2", info.DestHost,
		"--ssl2",
		"--user2", info.DestUser,
		"--password2", info.DestPass,
		"--exclude", "^Junk\\ E-Mail",
		"--exclude", "^Deleted\\ Items",
		"--exclude", "^Deleted",
		"--exclude", "^Trash",
		"--regextrans2", "s#^Sent$#Sent Items#",
		"--regextrans2", "s#^Spam$#Junk E-Mail#",
		"--useuid",
		"--usecache",
		"--tmpdir", tmpDir,
	}

	fmt.Printf("\n%s%sStarting email migration process...%s\n", ColorMagenta, ColorBold, ColorReset)
	fmt.Printf("%s%sSource:%s %s@%s\n", ColorBold, ColorWhite, ColorReset, info.SourceUser, info.SourceHost)
	fmt.Printf("%s%sDestination:%s %s@%s\n", ColorBold, ColorWhite, ColorReset, info.DestUser, info.DestHost)
	fmt.Printf("%s%sTemporary Directory:%s %s\n", ColorBold, ColorWhite, ColorReset, tmpDir)
	fmt.Printf("\n%s%sThis process may take a long time. Please wait...%s\n\n", ColorYellow, ColorBold, ColorReset)

	// Execute imapsync command
	cmd := exec.Command("imapsync", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n%s%s✗ Email migration process failed: %v%s\n", ColorRed, ColorBold, err, ColorReset)
		return false
	}

	fmt.Printf("\n%s%s✓ Email migration completed successfully!%s\n", ColorGreen, ColorBold, ColorReset)
	return true
}

// printDevASCII prints ASCII DEV message
func (cli *ImapSyncCLI) printDevASCII() {
	devASCII := fmt.Sprintf(`
%s%s
██████╗ ███████╗██╗   ██╗
██╔══██╗██╔════╝██║   ██║
██║  ██║█████╗  ██║   ██║
██║  ██║██╔══╝  ╚██╗ ██╔╝
██████╔╝███████╗ ╚████╔╝ 
╚═════╝ ╚══════╝  ╚═══╝  
%s
%s%s        Developer Information%s
%s%s         Created by Eren Can Ucar%s
`, ColorMagenta, ColorBold, ColorReset, ColorCyan, ColorBold, ColorReset, ColorBold, ColorWhite, ColorReset)

	fmt.Println(devASCII)
}

// showDevMenu displays developer information menu
func (cli *ImapSyncCLI) showDevMenu() {
	fmt.Printf("\n%s%sDeveloper Information:%s\n", ColorBold, ColorWhite, ColorReset)
	fmt.Printf("%s%s1 - Developer Email: %s%sdev@eren.gg%s\n", ColorGreen, ColorBold, ColorCyan, ColorBold, ColorReset)
	fmt.Printf("%s%s2 - Developer LinkedIn: %s%shttps://www.linkedin.com/in/erencanucarr/%s\n", ColorGreen, ColorBold, ColorCyan, ColorBold, ColorReset)
	fmt.Printf("%s%s3 - Developer Github: %s%shttps://github.com/erencanucarr%s\n", ColorGreen, ColorBold, ColorCyan, ColorBold, ColorReset)
	fmt.Printf("%s%sb - Main Menu%s\n", ColorYellow, ColorBold, ColorReset)
}

// handleDevSection handles developer information section
func (cli *ImapSyncCLI) handleDevSection() {
	cli.printDevASCII()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		cli.showDevMenu()

		fmt.Printf("\n%s%sYour choice: %s", ColorBold, ColorWhite, ColorReset)
		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Printf("\n%s%s📧 Developer Email:%s\n", ColorMagenta, ColorBold, ColorReset)
			fmt.Printf("%s%s✉️  dev@eren.gg%s\n", ColorGreen, ColorBold, ColorReset)
			fmt.Printf("%s%sYou can contact me by copying this email address.%s\n", ColorCyan, ColorBold, ColorReset)
			cli.waitForEnter()

		case "2":
			fmt.Printf("\n%s%s💼 Developer LinkedIn:%s\n", ColorMagenta, ColorBold, ColorReset)
			fmt.Printf("%s%s🔗 https://www.linkedin.com/in/erencanucarr/%s\n", ColorGreen, ColorBold, ColorReset)
			fmt.Printf("%s%sYou can visit my LinkedIn profile for my professional network and work experience.%s\n", ColorCyan, ColorBold, ColorReset)
			cli.waitForEnter()

		case "3":
			fmt.Printf("\n%s%s💻 Developer Github:%s\n", ColorMagenta, ColorBold, ColorReset)
			fmt.Printf("%s%s🐙 https://github.com/erencanucarr%s\n", ColorGreen, ColorBold, ColorReset)
			fmt.Printf("%s%sYou can check my GitHub profile for open source projects and code examples.%s\n", ColorCyan, ColorBold, ColorReset)
			cli.waitForEnter()

		case "b", "back", "main":
			return

		default:
			fmt.Printf("%s%sInvalid choice. Please enter 1, 2, 3 or b.%s\n", ColorYellow, ColorBold, ColorReset)
		}
	}
}

// waitForEnter waits for user to press Enter
func (cli *ImapSyncCLI) waitForEnter() {
	fmt.Printf("\n%s%sPress Enter to continue...%s", ColorBold, ColorWhite, ColorReset)
	bufio.NewScanner(os.Stdin).Scan()
} ColorCyan, ColorBold, ColorReset)
	fmt.Printf("%s%s3 - Developer Github: %s%shttps://github.com/erencanucarr%s\n", ColorGreen, ColorBold, ColorCyan, ColorBold, ColorReset)
	fmt.Printf("%s%sb - Ana Menü%s\n", ColorYellow, ColorBold, ColorReset)
}

// handleDevSection handles developer information section
func (cli *ImapSyncCLI) handleDevSection() {
	cli.printDevASCII()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		cli.showDevMenu()

		fmt.Printf("\n%s%sSeçiminiz: %s", ColorBold, ColorWhite, ColorReset)
		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Printf("\n%s%s📧 Developer Mail:%s\n", ColorMagenta, ColorBold, ColorReset)
			fmt.Printf("%s%s✉️  dev@eren.gg%s\n", ColorGreen, ColorBold, ColorReset)
			fmt.Printf("%s%sBu mail adresini kopyalayarak benimle iletişime geçebilirsiniz.%s\n", ColorCyan, ColorBold, ColorReset)
			cli.waitForEnter()

		case "2":
			fmt.Printf("\n%s%s💼 Developer LinkedIn:%s\n", ColorMagenta, ColorBold, ColorReset)
			fmt.Printf("%s%s🔗 https://www.linkedin.com/in/erencanucarr/%s\n", ColorGreen, ColorBold, ColorReset)
			fmt.Printf("%s%sProfesyonel ağım ve iş deneyimlerim için LinkedIn profilimi ziyaret edebilirsiniz.%s\n", ColorCyan, ColorBold, ColorReset)
			cli.waitForEnter()

		case "3":
			fmt.Printf("\n%s%s💻 Developer Github:%s\n", ColorMagenta, ColorBold, ColorReset)
			fmt.Printf("%s%s🐙 https://github.com/erencanucarr%s\n", ColorGreen, ColorBold, ColorReset)
			fmt.Printf("%s%sAçık kaynak projelerim ve kod örneklerim için GitHub profilimi inceleyebilirsiniz.%s\n", ColorCyan, ColorBold, ColorReset)
			cli.waitForEnter()

		case "b", "back", "geri", "ana":
			return

		default:
			fmt.Printf("%s%sGeçersiz seçim. Lütfen 1, 2, 3 veya b girin.%s\n", ColorYellow, ColorBold, ColorReset)
		}
	}
}

// waitForEnter waits for user to press Enter
func (cli *ImapSyncCLI) waitForEnter() {
	fmt.Printf("\n%s%sDevam etmek için Enter'a basın...%s", ColorBold, ColorWhite, ColorReset)
	bufio.NewScanner(os.Stdin).Scan()
}

// run starts the main application loop
func (cli *ImapSyncCLI) run() {
	cli.printWelcome()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		cli.showMenu()

		fmt.Printf("\n%s%sYour choice: %s", ColorBold, ColorWhite, ColorReset)
		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Printf("\n%s%sSystem Installation%s\n", ColorMagenta, ColorBold, ColorReset)
			cli.installSystemRequirements()

		case "2":
			fmt.Printf("\n%s%sEmail Migration%s\n", ColorMagenta, ColorBold, ColorReset)
			if migrationInfo := cli.collectMigrationInfo(); migrationInfo != nil {
				cli.runImapSync(migrationInfo)
			}

		case "3":
			cli.handleDevSection()

		case "q", "quit", "exit":
			fmt.Printf("\n%s%sSee you later!%s\n", ColorCyan, ColorBold, ColorReset)
			return

		default:
			fmt.Printf("%s%sInvalid choice. Please enter 1, 2, 3 or q.%s\n", ColorYellow, ColorBold, ColorReset)
		}

		// Add a small delay to make output more readable
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	// Handle Ctrl+C gracefully
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n\n%s%sSee you later!%s\n", ColorCyan, ColorBold, ColorReset)
		}
	}()

	cli := NewImapSyncCLI()
	cli.run()
}
