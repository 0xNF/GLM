package fsops

import (
	"fmt"
	"os"
)

func Fsops() {
	fmt.Printf("fsops lol")
}

// CheckExists checks whether the given `path` exists
func CheckExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil

	} else {
		return false, err
	}
}
