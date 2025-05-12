package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func OpenFileInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad"
		} else {
			if _, err := exec.LookPath("vim"); err == nil {
				editor = "vim"
			} else if _, err := exec.LookPath("nano"); err == nil {
				editor = "nano"
			} else if _, err := exec.LookPath("nvim"); err == nil {
				editor = "nvim"
			} else {
				return fmt.Errorf("no suitable editor found. Please install vim, nano, or nvim.")
			}
		}
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
