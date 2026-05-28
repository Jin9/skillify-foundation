package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	filepathAbs    = filepath.Abs
	osMkdirAll     = os.MkdirAll
	osOpenFile     = os.OpenFile
)

// ArchiveFile represents a single file to be placed into a compressed archive.
type ArchiveFile struct {
	Name    string
	Content io.Reader
}

// Zipper defines secure archiving and extraction operations.
type Zipper interface {
	// CompressStream combines multiple files into a zip archive and streams it to the given writer.
	CompressStream(files []ArchiveFile, out io.Writer) error

	// UnzipSecure extracts a zip archive safely to the given directory, guarding against Zip Slip attacks.
	UnzipSecure(in io.ReaderAt, size int64, destDir string) error
}

type zipper struct{}

// NewZipper creates a new instance of Zipper.
func NewZipper() Zipper {
	return &zipper{}
}

// CompressStream combines multiple files into a zip archive and streams it to the given writer.
func (z *zipper) CompressStream(files []ArchiveFile, out io.Writer) error {
	w := zip.NewWriter(out)
	defer w.Close()

	for _, file := range files {
		fWriter, err := w.Create(file.Name)
		if err != nil {
			return fmt.Errorf("failed to create entry %s in zip: %w", file.Name, err)
		}

		if file.Content != nil {
			if _, err := io.Copy(fWriter, file.Content); err != nil {
				return fmt.Errorf("failed to write content to zip entry %s: %w", file.Name, err)
			}
		}
	}

	return nil
}

// UnzipSecure extracts a zip archive safely to the given directory, guarding against Zip Slip attacks.
func (z *zipper) UnzipSecure(in io.ReaderAt, size int64, destDir string) error {
	r, err := zip.NewReader(in, size)
	if err != nil {
		return fmt.Errorf("failed to open zip reader: %w", err)
	}

	// Ensure destination directory is absolute to prevent relative path escapes
	destDir, err = filepathAbs(destDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for destDir: %w", err)
	}

	for _, f := range r.File {
		// Calculate the ultimate target path of the file
		fpath := filepath.Join(destDir, f.Name)

		// Check for Zip Slip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path %s: vulnerability 'zip slip' detected", fpath)
		}

		if f.FileInfo().IsDir() {
			err := osMkdirAll(fpath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fpath, err)
			}
			continue
		}

		if err = osMkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directory for file %s: %w", fpath, err)
		}

		outFile, err := osOpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to open output file %s: %w", fpath, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open compressed file %s: %w", f.Name, err)
		}

		_, err = io.Copy(outFile, rc)

		// Clean up handles early inside loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file %s: %w", f.Name, err)
		}
	}
	return nil
}
