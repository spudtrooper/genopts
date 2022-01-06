package genopts

func postGenCleanup(goImportsBin, dir, f string) error {
	if err := run(dir, goImportsBin, "-w", f); err != nil {
		return err
	}
	if err := run(dir, "go", "fmt", f); err != nil {
		return err
	}
	return nil
}
