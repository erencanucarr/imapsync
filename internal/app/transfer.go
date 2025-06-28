package app

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "os/exec"
    "regexp"
    "strconv"
    "golang.org/x/term"
    "github.com/schollz/progressbar/v3"
    "imapsync/internal/i18n"
    "imapsync/internal/ui"
)

// TransferMail runs imapsync and shows a progress bar.
// It parses stdout looking for "Transferred:" lines to update progress.
func TransferMail(lang string) {
    fmt.Println(ui.Cyan(i18n.T(lang, "transfer_start")))

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Source IMAP host: ")
    srcHost, _ := reader.ReadString('\n')
    srcHost = strings.TrimSpace(srcHost)

    fmt.Print("Source email: ")
    srcEmail, _ := reader.ReadString('\n')
    srcEmail = strings.TrimSpace(srcEmail)

    fmt.Print("Source password: ")
    srcPass, _ := term.ReadPassword(int(os.Stdin.Fd()))
    fmt.Println()

    fmt.Print("Destination IMAP host: ")
    dstHost, _ := reader.ReadString('\n')
    dstHost = strings.TrimSpace(dstHost)

    fmt.Print("Destination email: ")
    dstEmail, _ := reader.ReadString('\n')
    dstEmail = strings.TrimSpace(dstEmail)

    fmt.Print("Destination password: ")
    dstPass, _ := term.ReadPassword(int(os.Stdin.Fd()))
    fmt.Println()

    fmt.Println(ui.Cyan("Testing credentials..."))
    testCmd := exec.Command("imapsync", "--justlogin", "--host1", srcHost, "--ssl1", "--user1", srcEmail, "--password1", string(srcPass), "--host2", dstHost, "--ssl2", "--user2", dstEmail, "--password2", string(dstPass))
    if err := testCmd.Run(); err != nil {
        fmt.Println(ui.Red(i18n.T(lang, "error")), err)
        return
    }

    args := []string{
        "--host1", srcHost, "--ssl1",
        "--user1", srcEmail, "--password1", string(srcPass),
        "--host2", dstHost, "--ssl2",
        "--user2", dstEmail, "--password2", string(dstPass),
        "--exclude", "^Junk\\ E-Mail",
        "--exclude", "^Deleted\\ Items",
        "--exclude", "^Deleted",
        "--exclude", "^Trash",
        "--regextrans2", "s#^Sent$#Sent Items#",
        "--regextrans2", "s#^Spam$#Junk E-Mail#",
        "--useuid",
        "--usecache",
        "--tmpdir", "./tmp",
        "--syncinternaldates",
        "--progress",
    }
    cmd := exec.Command("imapsync", args...)

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        fmt.Println(ui.Red(i18n.T(lang, "error")), err)
        return
    }

    if err := cmd.Start(); err != nil {
        fmt.Println(ui.Red(i18n.T(lang, "error")), err)
        return
    }

    bar := progressbar.NewOptions(100,
        progressbar.OptionSetDescription("IMAPSYNC"),
        progressbar.OptionSetTheme(progressbar.Theme{Saucer: "=", SaucerHead: ">", SaucerPadding: " ", BarStart: "[", BarEnd: "]"}),
    )

    scanner := bufio.NewScanner(stdout)
    percentRe := regexp.MustCompile(`([0-9]{1,3}(?:\.[0-9]+)?)%`)
    ratioRe := regexp.MustCompile(`(?i)(\d+)/(\d+)`)
    for scanner.Scan() {
        line := scanner.Text()

        if m := percentRe.FindStringSubmatch(line); len(m) == 2 {
            p, _ := strconv.ParseFloat(m[1], 64)
            bar.Set(int(p))
            continue
        }
        if m := ratioRe.FindStringSubmatch(line); len(m) == 3 {
            current, _ := strconv.Atoi(m[1])
            total, _ := strconv.Atoi(m[2])
            if total > 0 {
                bar.Set(int(float64(current) / float64(total) * 100))
            }
        }
    }

    if err := cmd.Wait(); err != nil {
        fmt.Println() // newline after bar
        fmt.Println(ui.Red(i18n.T(lang, "transfer_fail")))
        fmt.Println(ui.Red(i18n.T(lang, "error")), err)
        return
    }

    bar.Finish()
    fmt.Println() // newline after bar
    fmt.Println(ui.Green(i18n.T(lang, "transfer_success")))
}
