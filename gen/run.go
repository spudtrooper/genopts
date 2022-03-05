package gen

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

var (
	verboseRun = flag.Bool("verbose_run", false, "print command lines when running")
)

func run(dir, command string, args ...string) error {
	if *verboseRun {
		log.Printf("Running: %s %v", command, args)
	}
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
