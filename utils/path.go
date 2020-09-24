package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathCreate(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func FileCreate(content bytes.Buffer, name string) {
	file, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	file.WriteString(content.String())
	//for i := 0; i < len(content); i++ {
	//	file.Write(content)
	//}
	file.Close()
}

func PathRemove(name string) {
	err := os.RemoveAll(name)
	if err != nil {
		log.Println(err)
	}
}

func FileRemove(name string) {
	err := os.Remove(name)
	if err != nil {
		log.Println(err)
	}
}

func FileZip(dst, src string, notContPath string) (err error) {
	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}
	zw := zip.NewWriter(fw)
	defer func() {
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		if fi.IsDir() {
			fh.Name += "/"
		}
		fh.Name = strings.Replace(fh.Name, notContPath, "", -1)

		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		if !fh.Mode().IsRegular() {
			return nil
		}

		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		n, err := io.Copy(w, fr)
		if err != nil {
			return
		}

		log.Printf("Zip file： %s, write %d char\n", path, n)

		return nil
	})
}

type ReplaceHelper struct {
	Root    string
	OldText string
	NewText string
}

func (h *ReplaceHelper) DoWrok() error {
	return filepath.Walk(h.Root, h.walkCallback)
}

func (h ReplaceHelper) walkCallback(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if f == nil {
		return nil
	}
	if f.IsDir() {
		log.Println("DIR:", path)
		return nil
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		//err
		return err
	}
	content := string(buf)
	//log.Printf("h.OldText: %s \n", h.OldText)
	//log.Printf("h.NewText: %s \n", h.NewText)
	newContent := strings.Replace(content, h.OldText, h.NewText, -1)
	ioutil.WriteFile(path, []byte(newContent), 0)
	return err
}
