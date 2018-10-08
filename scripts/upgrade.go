package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/sipt/shuttle/extension/config"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	var (
		fileName     string
		err          error
		cmd          *exec.Cmd
		configBackup []byte
	)
	flag.StringVar(&fileName, "f", "shuttle.zip", "shuttle upgrade zip")
	flag.Parse()
	configBackup, _ = ioutil.ReadFile(filepath.Join(config.ShuttleHomeDir, "shuttle.yaml"))
	time.Sleep(3 * time.Second) // wait for shuttle shutdown.
	err = ClearDir(getCurrentDirectory())
	if err != nil {
		goto Failed
	}
	err = Unzip(filepath.Join(config.HomeDir, "Downloads", fileName), ".."+string(os.PathSeparator))
	if err != nil {
		goto Failed
	}
	ioutil.WriteFile(filepath.Join(config.ShuttleHomeDir, "shuttle.yaml"), configBackup, 0644)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("startup")
	} else {
		cmd = exec.Command("./start.sh")
	}
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
	os.MkdirAll(filepath.Join(config.ShuttleHomeDir, "logs"), 0755)
	ioutil.WriteFile(filepath.Join(config.ShuttleHomeDir, "logs", "upgrade.log"), []byte(err.Error()), 0664)
}

func ClearDir(dir string) error {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, v := range infos {
		if v.IsDir() {
			err = os.RemoveAll(v.Name())
		} else {
			err = os.Remove(v.Name())
		}
		if err != nil {
			fmt.Println(err)
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
