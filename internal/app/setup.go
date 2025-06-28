package app

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strings"

    "imapsync/internal/i18n"
    "imapsync/internal/ui"
)

// checkBinary returns true if command exists in PATH.
func checkBinary(name string) bool {
    _, err := exec.LookPath(name)
    return err == nil
}

// SetupSystem verifies Python and imapsync; prints install guidance if missing.
func SetupSystem(lang string) {
    // Check Python
    pythonOK := checkBinary("python") || checkBinary("python3")
    if pythonOK {
        fmt.Println(ui.Green("Python ✔"))
    } else {
        fmt.Println(ui.Red("Python ✗"))
        switch runtime.GOOS {
        case "windows":
            fmt.Println("Download Python: https://www.python.org/downloads/windows/")
        case "darwin":
            fmt.Println("brew install python")
        default:
            fmt.Println("sudo apt install python3")
            if checkBinary("apt") {
                exec.Command("sudo", "apt", "update").Run()
                exec.Command("sudo", "apt", "install", "-y", "python3").Run()
            } else if checkBinary("yum") {
                exec.Command("sudo", "yum", "install", "-y", "python3").Run()
            }
        }
    }

    // Check imapsync
    imapOK := checkBinary("imapsync")
    if imapOK {
        fmt.Println(ui.Green("imapsync ✔"))
    } else {
        fmt.Println(ui.Red("imapsync ✗"))
        promptInstall()
    }

    if pythonOK && imapOK {
        fmt.Println(ui.Green(i18n.T(lang, "menu_setup") + " ✅"))
    }
}

// promptInstall lets user choose a local install script and executes it via bash.
func promptInstall() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Local install scripts available in ./install directory:")
    scripts := map[string]string{
        "ubuntu": "install/ubuntu.txt",
        "debian": "install/debian.txt",
        "centos": "install/centos.txt",
        "arch": "install/archlinux.txt",
        "darwin": "install/darwin.txt",
    }
    for k := range scripts {
        fmt.Println(" -", k)
    }
    fmt.Print("Enter distribution key to run installer or press Enter to skip: ")
    choice, _ := reader.ReadString('\n')
    choice = strings.TrimSpace(strings.ToLower(choice))
    path, ok := scripts[choice]
    if !ok || choice == "" {
        fmt.Println("Skipping automatic install. Follow manual instructions in INSTALL.d directory.")
        return
    }
    fmt.Println(ui.Cyan("Running installation script:"), path)
    if err := exec.Command("bash", path).Run(); err != nil {
        fmt.Println(ui.Red("Installer failed:"), err)
    }
}

// installImapsync tries to install imapsync according to the OS.
func installImapsync() {
    switch runtime.GOOS {
    case "windows":
        // Use chocolatey or scoop if available
        if checkBinary("choco") {
            exec.Command("choco", "install", "imapsync", "-y").Run()
        } else if checkBinary("scoop") {
            exec.Command("scoop", "install", "imapsync").Run()
        } else {
            fmt.Println("Please install Chocolatey or Scoop, or download the binary from https://imapsync.lamiral.info/")
        }
    case "darwin":
        exec.Command("brew", "install", "imapsync").Run()
    default: // assume linux
        if checkBinary("apt") {
            exec.Command("sudo", "apt", "update").Run()
            exec.Command("sudo", "apt", "install", "-y", "imapsync").Run()
        } else if checkBinary("yum") {
            exec.Command("sudo", "yum", "install", "-y", "imapsync").Run()
        } else {
            fmt.Println("Please install imapsync manually: https://imapsync.lamiral.info/#install")
        }
    }
}
