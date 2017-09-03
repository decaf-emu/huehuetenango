package importer

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/phayes/permbits"
)

type zipSource struct {
	path      string
	processed bool
	targetDir string
	dirSource Source
}

func NewZIPSource(path string) Source {
	return &zipSource{
		path: path,
	}
}

func (s *zipSource) Open() error {
	// skip extracting the ZIP if it's already been processed
	if s.processed {
		return nil
	}

	if err := s.extract(); err != nil {
		return err
	}

	// use the directory the files were extracted to as the source
	s.dirSource = NewSource(s.targetDir)
	if err := s.dirSource.Open(); err != nil {
		return err
	}

	s.processed = true
	return nil
}

func (s *zipSource) Titles() ([]*Title, error) {
	if s.dirSource != nil {
		return s.dirSource.Titles()
	}
	return nil, errors.New("source isn't open")
}

func (s *zipSource) Close() error {
	// dirSource will be nil if the source hasn't been opened
	if s.dirSource != nil {
		if err := s.dirSource.Close(); err != nil {
			return err
		}
	}

	return s.removeTargetDir()
}

func (s *zipSource) extract() error {
	// clean-up if the ZIP has already been extracted
	if s.targetDir != "" {
		if err := s.removeTargetDir(); err != nil {
			return err
		}
	}
	// create the target directory for extracting the files
	if err := s.createTargetDir(); err != nil {
		return err
	}

	reader, err := zip.OpenReader(s.path)
	if err != nil {
		return err
	}
	defer reader.Close()

	// extract each file to the target directory
	for _, file := range reader.File {
		if err := s.extractFile(file); err != nil {
			return err
		}
	}

	return nil
}

func (s *zipSource) extractFile(file *zip.File) error {
	// create the output path for the file or directory
	path := filepath.Join(s.targetDir, file.Name)

	directory := path
	if !file.FileInfo().IsDir() {
		directory = filepath.Dir(directory)
	}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		permissions := os.FileMode(0755)
		if file.FileInfo().IsDir() {
			permissions = file.Mode()
		}

		// create the directory if it doesn't already exist
		if err := os.MkdirAll(directory, permissions); err != nil {
			return err
		}
	} else if file.FileInfo().IsDir() {
		// set the directory's file mode to the zip entry's
		// the directory could have been created when extracting a file
		// in which case the directory was created as 0755
		permissions := permbits.FileMode(file.Mode())
		if err := permbits.Chmod(directory, permissions); err != nil {
			return err
		}
	}

	if !file.FileInfo().IsDir() {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer output.Close()

		// extract the file
		_, err = io.Copy(output, fileReader)
		return err
	}

	return nil
}

func (s *zipSource) createTargetDir() error {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	s.targetDir = dir
	return nil
}

func (s *zipSource) removeTargetDir() error {
	if s.targetDir != "" {
		if err := os.RemoveAll(s.targetDir); err != nil {
			return err
		}
		s.targetDir = ""
	}
	return nil
}
