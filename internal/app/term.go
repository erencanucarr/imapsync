package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ReadPassword reads a password from stdin without echoing
func ReadPassword() (string, error) {
	fmt.Print("Password: ")
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}
