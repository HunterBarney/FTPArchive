package cleanup

import (
	"FTPArchive/internal/config"
	"fmt"
	"os"
)

func Cleanup(profile *config.Profile) error {
	if profile.CleanupDownloads {
		err := os.RemoveAll(profile.OutputName)
		if err != nil {
			return fmt.Errorf("could not remove downloaded files. Error: %s", err)
		}
	}

	if profile.CleanupArchives {
		err := os.Remove(profile.ArchivePath)
		if err != nil {
			return fmt.Errorf("could not remove archives directory. Error: %s", err)
		}
	}

	return nil
}
