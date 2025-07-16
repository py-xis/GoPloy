package builder

import (
	"bufio"
	"fmt"
	"os/exec"
)

func RunShellCommand(dir, command string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && %s", dir, command))

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("builder - RunShellCommand : failed to start command: %w", err)
	}

	go streamOutput(stdout, "RunShellCommand")
	go streamOutput(stderr, "RunShellCommand")

	return cmd.Wait()
}

func streamOutput(pipe interface{ Read([]byte) (int, error) }, functionName string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		formattedLine := fmt.Sprintf("builder - %s : %s", functionName, line)
		fmt.Println(formattedLine)
		PublishLog(formattedLine)
	}
}