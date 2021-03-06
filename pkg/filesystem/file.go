package filesystem

import (
	"bufio"
	"fmt"
	util2 "github.com/ankorstore/ankorstore-cli-core/core/util"
	"github.com/go-errors/errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func SaveBinaryFile(target string, assetReader io.ReadCloser) error {
	f := fmt.Sprintf("%s/%s", target, util2.AppName)
	err := FolderExist(f)
	if err == nil {
		err = Delete(f)
		if err != nil {
			return errors.Wrap(err, 0)
		}
	}
	assetFile, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0700)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	defer assetFile.Close()

	_, err = io.Copy(assetFile, assetReader)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	return nil
}

// SaveFile Stores content into a file
func SaveFile(rc io.ReadCloser, executablePath string, mode fs.FileMode) error {

	dir := path.Dir(executablePath)
	ext := path.Ext(executablePath)
	filename := path.Base(executablePath)

	dirs := util2.NewDirs()
	binDir := dirs.GetBinDir()

	if len(ext) > 0 {
		filename = fmt.Sprintf("%s.%s", filename, ext)
	}

	tmpPath := fmt.Sprintf("%s/%s", dir, "tmp-file")
	finalPath := fmt.Sprintf("%s/%s", binDir, filename)

	outputFile, err := os.Create(tmpPath)

	defer outputFile.Close()

	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}

	_, err = io.Copy(outputFile, rc)
	if err != nil {
		return fmt.Errorf("failed to store file: %v", err)
	}

	err = os.Chmod(tmpPath, mode)
	if err != nil {
		return fmt.Errorf("failed to set permission: %v", err)
	}

	err = os.Rename(tmpPath, finalPath)
	if err != nil {
		return fmt.Errorf("failed move file: %v", err)
	}

	return nil
}

// FolderExist Checks if a folder exists in the filesystem.
func FolderExist(target string) error {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return err
	}
	return nil
}

// Delete Removes a folder or file from the filesystem.
func Delete(target string) error {
	err := os.RemoveAll(target)
	if err != nil {
		return err
	}
	return nil
}

// CreateFolder recursively creates required directories.
func CreateFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0770)
	}
	return nil
}

// GetFileContent Retrieves the content of a file.
func GetFileContent(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// WriteToFile Writes text to a file
func WriteToFile(target string, text string) error {
	f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer f.Close()

	if err != nil {
		return err
	}

	if _, err := f.WriteString(text); err != nil {
		return err
	}

	return nil
}

// FileHasText Checks if a file already have specific text
func FileHasText(target string, text string) (bool, error) {
	f, err := os.Open(target)
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			return true, nil
		}
	}

	return false, nil
}
