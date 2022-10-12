package notify

import (
	"os/exec"
)

func Send(title, message string, ignore bool) error {
	if ignore {
		return nil
	}

	cmd := exec.Command("notify-send", title, message)

	return cmd.Run()
}
