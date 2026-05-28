package zip

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failReader struct{}

func (f failReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read failed")
}

type failWriter struct{}

func (f failWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write failed")
}

type failReaderAt struct {
	io.ReaderAt
}

func (f failReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 100 { // Fails reads occurring at beginning like local file headers
		return 0, errors.New("read failed")
	}
	return f.ReaderAt.ReadAt(p, off)
}

func TestZipper_CompressStream(t *testing.T) {
	z := NewZipper()

	t.Run("success", func(t *testing.T) {
		out := new(bytes.Buffer)
		files := []ArchiveFile{
			{Name: "test.txt", Content: strings.NewReader("hello world")},
			{Name: "empty.txt", Content: nil},
		}

		err := z.CompressStream(files, out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify zip format
		r, err := zip.NewReader(bytes.NewReader(out.Bytes()), int64(out.Len()))
		if err != nil {
			t.Fatalf("invalid zip created: %v", err)
		}
		if len(r.File) != 2 {
			t.Errorf("expected 2 files, found %d", len(r.File))
		}
	})

	t.Run("io error during copy", func(t *testing.T) {
		out := new(bytes.Buffer)
		files := []ArchiveFile{
			{Name: "test.txt", Content: failReader{}},
		}

		err := z.CompressStream(files, out)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "failed to write content") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("create entry error", func(t *testing.T) {
		files := []ArchiveFile{
			{Name: strings.Repeat("a", 70000), Content: nil},
		}
		out := new(bytes.Buffer)
		err := z.CompressStream(files, out)
		if err == nil {
			t.Fatal("expected error for header write failure")
		}
	})
}

func generateVulnerableZip(t *testing.T) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// Malicious file name
	f, err := w.Create("../../../malicious.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write([]byte("malicious content"))
	w.Close()
	return buf.Bytes()
}

func generateValidZip(t *testing.T) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	f, _ := w.Create("safe.txt")
	f.Write([]byte("safe content"))

	dirHeader := &zip.FileHeader{
		Name: "subdir/",
	}
	dirHeader.SetMode(os.ModeDir)
	w.CreateHeader(dirHeader)

	w.Close()
	return buf.Bytes()
}

func TestZipper_UnzipSecure(t *testing.T) {
	z := NewZipper()

	t.Run("invalid zip", func(t *testing.T) {
		err := z.UnzipSecure(bytes.NewReader(nil), 0, os.TempDir())
		if err == nil {
			t.Fatal("expected error for empty zip")
		}
	})

	t.Run("zip slip vulnerability", func(t *testing.T) {
		vulnZip := generateVulnerableZip(t)
		err := z.UnzipSecure(bytes.NewReader(vulnZip), int64(len(vulnZip)), os.TempDir())
		if err == nil {
			t.Fatal("expected error for zip slip attack")
		}
		if !strings.Contains(err.Error(), "zip slip") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		safeZip := generateValidZip(t)

		destDir := filepath.Join(os.TempDir(), "zipper_test_extract")
		os.RemoveAll(destDir) // Ensure clean state
		defer os.RemoveAll(destDir)

		err := z.UnzipSecure(bytes.NewReader(safeZip), int64(len(safeZip)), destDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify extraction
		if _, err := os.Stat(filepath.Join(destDir, "safe.txt")); os.IsNotExist(err) {
			t.Error("safe.txt not extracted")
		}
		info, err := os.Stat(filepath.Join(destDir, "subdir"))
		if err != nil || !info.IsDir() {
			t.Error("subdir directory not created")
		}
	})

	t.Run("filepath abs error", func(t *testing.T) {
		origAbs := filepathAbs
		defer func() { filepathAbs = origAbs }()
		filepathAbs = func(path string) (string, error) {
			return "", errors.New("mock abs error")
		}

		safeZip := generateValidZip(t)
		err := z.UnzipSecure(bytes.NewReader(safeZip), int64(len(safeZip)), os.TempDir())
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("mkdir info error", func(t *testing.T) {
		origMkdir := osMkdirAll
		defer func() { osMkdirAll = origMkdir }()
		osMkdirAll = func(path string, perm os.FileMode) error {
			if strings.HasSuffix(path, "subdir") || strings.HasSuffix(path, "subdir/") {
				return errors.New("mock mkdir dir error")
			}
			return nil
		}

		safeZip := generateValidZip(t)
		err := z.UnzipSecure(bytes.NewReader(safeZip), int64(len(safeZip)), os.TempDir())
		if err == nil {
			t.Fatal("expected err")
		}
	})

	t.Run("mkdir parent error", func(t *testing.T) {
		origMkdir := osMkdirAll
		defer func() { osMkdirAll = origMkdir }()
		osMkdirAll = func(path string, perm os.FileMode) error {
			if !strings.HasSuffix(path, "subdir") && !strings.HasSuffix(path, "subdir/") {
				return errors.New("mock mkdir file error")
			}
			return nil // allow the directory creation itself
		}

		safeZip := generateValidZip(t)
		err := z.UnzipSecure(bytes.NewReader(safeZip), int64(len(safeZip)), os.TempDir())
		if err == nil {
			t.Fatal("expected err")
		}
	})

	t.Run("open file error", func(t *testing.T) {
		origOpen := osOpenFile
		defer func() { osOpenFile = origOpen }()
		osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
			return nil, errors.New("mock open file error")
		}

		safeZip := generateValidZip(t)
		err := z.UnzipSecure(bytes.NewReader(safeZip), int64(len(safeZip)), os.TempDir())
		if err == nil {
			t.Fatal("expected err")
		}
	})

	t.Run("f open error", func(t *testing.T) {
		safeZip := generateValidZip(t)
		// Break the local file header signature to make f.Open() fail.
		// Local file header signature is 0x50, 0x4b, 0x03, 0x04.
		idx := bytes.Index(safeZip, []byte{0x50, 0x4b, 0x03, 0x04})
		if idx >= 0 {
			safeZip[idx] = 0x99
		}

		err := z.UnzipSecure(bytes.NewReader(safeZip), int64(len(safeZip)), os.TempDir())
		if err == nil {
			t.Fatal("expected f open error due to corrupt local file header")
		}
	})

	t.Run("copy error", func(t *testing.T) {
		origOpen := osOpenFile
		defer func() { osOpenFile = origOpen }()

		destDir, _ := os.MkdirTemp("", "zipper_copy_test")
		defer os.RemoveAll(destDir)

		osOpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
			// To make io.Copy fail, we return a closed file
			f, err := origOpen(name, flag, perm)
			if f != nil {
				f.Close()
			}
			return f, err
		}

		safeZip := generateValidZip(t)
		err := z.UnzipSecure(bytes.NewReader(safeZip), int64(len(safeZip)), destDir)
		if err == nil {
			t.Fatal("expected copy error due to closed output file")
		}
	})
}
