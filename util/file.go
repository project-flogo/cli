package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func IsRemote(path string) bool {
	return strings.HasPrefix(path, "http")
}

func LoadRemoteFile(sourceURL string) (string, error) {

	resp, err := http.Get(sourceURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func LoadLocalFile(path string) (string, error) {

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func Rename(srcPath string) error {
	oldPath := fmt.Sprintf("%s.old", srcPath)
	_ = os.Remove(oldPath)

	err := os.Rename(srcPath, oldPath)
	if err != nil {
		return err
	}

	return nil
}

func CopyFile(srcFile, destFile string) error {
	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destFile, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Copy(src, dst string, copyMode bool) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDir(src, dst, info, copyMode)
	}

	return copyFile(src, dst, info, copyMode)
}

func copyDir(srcDir, dstDir string, info os.FileInfo, copyMode bool) error {

	if err := os.MkdirAll(dstDir, os.FileMode(0755)); err != nil {
		return err
	}

	if copyMode {
		defer os.Chmod(dstDir, info.Mode())
	}

	items, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, item := range items {

		srcPath := path.Join(srcDir, item.Name())
		dstPath := path.Join(dstDir, item.Name())

		if item.IsDir() {
			if err = copyDir(srcPath, dstPath, item, copyMode); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcPath, dstPath, item, copyMode); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string, srcInfo os.FileInfo, copyMode bool) error {
	var err error

	if err = os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	var sf *os.File
	if sf, err = os.Open(src); err != nil {
		return err
	}
	defer sf.Close()

	var df *os.File
	if df, err = os.Create(dst); err != nil {
		return err
	}
	defer df.Close()

	if err = os.Chmod(sf.Name(), srcInfo.Mode()); err != nil {
		return err
	}

	if _, err = io.Copy(df, sf); err != nil {
		return err
	}

	return nil
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func DeleteFile(path string) error {

	if _, err := os.Stat(path); err == nil {
		err = os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

// SwapFile is like a copy but use a temporary file and rename it
// to allow running executable replacement.
// Thanks to https://gist.github.com/fenollp/7e31e6462b10c96aef443351bce6aea7
func SwapFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	destDir := filepath.Dir(dst)
	tmpFile := filepath.Join(destDir, "exe_swap")

	defer DeleteFile(tmpFile)

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(tmpFile, data, srcInfo.Mode()); err != nil {
		return err
	}

	if err := os.Rename(tmpFile, dst); err != nil {
		return err
	}

	return nil

}
