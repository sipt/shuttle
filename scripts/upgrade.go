package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"io"
	"os/exec"
)

func main() {
	//CoverApp()
	// start app
	cmd := exec.Command("cd shuttle; ./shuttle")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func CoverApp() error {
	// cover app
	os.RemoveAll("shuttle")
	// unzip
	err := Unzip("shuttle_macos_amd64_beta_v0.4.1.zip", "")
	if err != nil {
		return err
	}
	return nil
}

func Unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dst, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
