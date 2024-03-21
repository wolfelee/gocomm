package base

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func EtHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	home := path.Join(dir, ".Et")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func EtHomeWithDir(dir string) string {
	home := path.Join(EtHome(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func copyFile(src, dst string, replaces []string) error {
	var err error
	// if strings.HasSuffix(dst, "go.sum") {
	// 	return err
	// }
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	var old string
	for i, next := range replaces {
		if i%2 == 0 {
			old = next
			continue
		}

		if strings.HasSuffix(dst, "appconst.go") {
			buf = bytes.ReplaceAll(buf, []byte("goinit"), []byte(next))
		} else if strings.HasSuffix(dst, "go.mod") {
			buf = bytes.ReplaceAll(buf, []byte(old), []byte(next))
		} else {
			buf = bytes.ReplaceAll(buf, []byte(old+"/"), []byte(next+"/"))
		}

	}
	return ioutil.WriteFile(dst, buf, srcinfo.Mode())
}

func copyDir(src, dst string, replaces, ignores []string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}

		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp, replaces, ignores); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcfp, dstfp, replaces); err != nil {
				return err
			}
		}
	}
	return nil
}

func hasSets(name string, sets []string) bool {
	for _, ig := range sets {
		if ig == name {
			return true
		}
	}
	return false
}
