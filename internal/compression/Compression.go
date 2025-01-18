package compression

import (
	"FTPArchive/internal/config"
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// CompressToZip compresses the given directory path to a zip file
func CompressToZip(source string, config *config.Config) error {

	archiveFolder := config.ArchiveDirectory
	if _, err := os.Stat(archiveFolder); os.IsNotExist(err) {
		err = os.MkdirAll(archiveFolder, 0755)
		if err != nil {
			log.Fatal("Error creating log folder: ", err)
		}
	}

	basename := filepath.Base(source)
	archive, e := os.Create(archiveFolder + "/" + basename + ".zip")
	if e != nil {
		return fmt.Errorf("failed to create zip file: %w", e)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("could not create zip header: %w", err)
		}

		relativePath, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("could not determine relative path: %w", err)
		}
		if info.IsDir() {
			relativePath += "/"
		}
		header.Name = relativePath

		if !info.IsDir() {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("could not create writer for %s: %w", path, err)
		}

		if !info.IsDir() {
			sourceFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("could not open file %s: %w", path, err)
			}
			defer sourceFile.Close()

			if _, err := io.Copy(writer, sourceFile); err != nil {
				return fmt.Errorf("error writing file %s to zip: %w", path, err)
			}
		}

		return nil
	})
}
