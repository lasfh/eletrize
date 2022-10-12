//go:build !linux

package notify

func Send(title, message string, ignore bool) error {
	return nil
}
