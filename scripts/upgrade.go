package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"io"
	"fmt"
	"io/ioutil"
	"flag"
	"github.com/sipt/shuttle/extension/config"
	"os/exec"
	"time"
)

func main() {
	var (
		fileName string
		err      error
		cmd      *exec.Cmd
	)
	time.Sleep(time.Second) // wait for shuttle shutdown.
	flag.StringVar(&fileName, "f", "shuttle.zip", "shuttle upgrade zip")
	err = ClearDir(getCurrentDirectory())
	if err != nil {
		goto Failed
	}
	err = Unzip(filepath.Join(config.HomeDir, "Downloads", fileName), "../")
	if err != nil {
		goto Failed
	}
	cmd = exec.Command("./start.sh")
	err = cmd.Start()
	if err != nil {
		goto Failed
	}
	err = cmd.Wait()
	if err != nil {
		goto Failed
	}
	return
Failed:
	ioutil.WriteFile(filepath.Join(config.HomeDir, "Documents", "shuttle", "logs", "upgrade.log"), []byte(err.Error()), 0664)
}

func ClearDir(dir string) error {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range infos {
		if v.IsDir() {
			os.RemoveAll(v.Name())
		} else {
			os.Remove(v.Name())
		}
	}
	return nil
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func Unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Printf("[Upgrade] [Unzip] failed: %s", err.Error())
		}
	}()

	if len(dst) > 0 {
		os.MkdirAll(dst, 0755)
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				fmt.Printf("[Upgrade] [Unzip] failed: %s", err.Error())
			}
		}()
		path := f.Name
		if len(dst) > 0 {
			path = filepath.Join(dst, f.Name)
		}

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
					if err := r.Close(); err != nil {
						fmt.Printf("[Upgrade] [Unzip] failed: %s", err.Error())
					}
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
