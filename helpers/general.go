package helpers

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ExecuteCommand : Execute Shell Command
func ExecuteCommand(input string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", input)
	stdout, stderr := cmd.Output()
	return strings.TrimSpace(string(stdout)), stderr
}

// ExecuteWithLiveOutput : executing given command with realtime output
func ExecuteWithLiveOutput(command string) {
	cmd := exec.Command("bash", "-c", command)
	stdout, _ := cmd.StdoutPipe()

	_ = cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println(m)
	}
	_ = cmd.Wait()
}

// Abort : Abort system and print message
func Abort(message string) {
	fmt.Println(message)
	os.Exit(0)
}

// GetFreePort Get free port to listen on.
func GetFreePort() int {
	// Let's get port after 5454
	var port int
	var startPort = 5454
	for {
		addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:"+strconv.Itoa(startPort))
		_, err := net.ListenTCP("tcp", addr)
		if err == nil {
			port = startPort
			break
		} else {
			startPort++
		}
	}

	return port
}
