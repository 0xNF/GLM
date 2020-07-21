package monitor

import (
	"os"
	"path/filepath"
	"regexp"

	conf "github.com/0xNF/glm/src/internal/conf"
	fsops "github.com/0xNF/glm/src/internal/fsops"
)

// ScanForTrigger scans the folder asccording to the Trigger for the trigger file,
// returning true if the trigger file exists, false otherwise,
// and errors if errors are encountered
func scanForTrigger(triggerFile string) (bool, error) {
	return fsops.CheckExists(triggerFile)
}

// Monitors runs the main monitor logic
func Monitor(trigger *conf.GLMTrigger) (bool, error) {
	ret := false
	files, err := findMatchingFiles(trigger)
	if err != nil {
		return ret, err
	}

	/* early quit */
	if len(files) == 0 {
		return ret, nil
	}

	/* check if our trigger is present */
	found, err := scanForTrigger(trigger.TriggerFile)
	if err != nil {
		return ret, err
	}
	if found {
		/* if so, preserve the files */
		// if err = moveFilesToSaveLocation(trigger.SaveTo, files); err != nil {
		// 	return ret, err
		// }
		ret = true
	} else {
		/* if not, delete the files */
		// if err = deleteExcess(files); err != nil {
		// 	return ret, err
		// }
	}
	return ret, nil
}

// FindMatchingFiles takes a trigger and returns any files that match its specifications
func findMatchingFiles(trigger *conf.GLMTrigger) ([]string, error) {
	var filePattern = regexp.MustCompile(trigger.SavePattern)

	var files []string

	root := trigger.SaveFromFolder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		/* check if the file matches the pattern */
		if !info.IsDir() {
			matched := filePattern.MatchString(info.Name())
			if matched {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// MoveFilesToSaveLocation moves any matching files found to the specified Save Location
// returns an error if there was a problem, nothing otherwise
func moveFilesToSaveLocation(saveTo string, moveThese []string) error {

	/* return early if nothing matches */
	if len(moveThese) == 0 {
		return nil
	}
	/* create the save directory, if necessary */
	err := os.MkdirAll(saveTo, os.FileMode(0777))
	if err != nil {
		return err
	}
	/* move files */
	for _, file := range moveThese {
		fname := filepath.Base(file)
		if err = os.Rename(file, filepath.Join(saveTo, fname)); err != nil {
			return err
		}
	}
	return nil
}

// DeleteExcess deletes all files for the pattern that weren't saved.
// returns an error if there was a problem, nothing otherwise
func deleteExcess(deleteThese []string) error {
	for _, file := range deleteThese {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}
