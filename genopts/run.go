package genopts

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func run(dir, command string, args ...string) error {
	log.Printf("running from %s: %s %s", dir, command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
