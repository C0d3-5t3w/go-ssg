package editor

import (
	"os"
	"os/exec"
	"runtime"
)

// OpenFileInEditor opens the specified file path in a preferred text editor.
// It prioritizes $EDITOR, then vim, then nano.
func OpenFileInEditor(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		if runtime.GOOS == "windows" {
			editor = "notepad" // A very basic default for Windows
		} else {
			// Check for common editors
			if _, err := exec.LookPath("vim"); err == nil {
				editor = "vim"
			} else if _, err := exec.LookPath("nano"); err == nil {
				editor = "nano"
			} else {
				// As a last resort, try 'edit' on Unix-like or 'open' on macOS
				// 'edit' is less common, 'open' might open with a GUI app.
				// For a CLI tool, sticking to CLI editors is better.
				// If no common CLI editor is found, this might be an issue.
				// Consider adding 'vi' or other common ones.
				// For now, if vim/nano aren't found, this will likely fail to find an editor.
				// A more robust solution might involve more checks or allowing user config.
				editor = "vim" // Default to vim if nothing else is found
			}
		}
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
