package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var regexps []*regexp.Regexp

func init() {
	regexps = make([]*regexp.Regexp, 0)
	file, err := os.Open(".hashignore")

	if err != nil {
		fmt.Println("error:", err, os.IsExist(err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.Index(text, "#") == 0 {
			continue
		}
		if strings.Index(text, "*") == 0 || strings.Index(text, "?") == 0 {
			text = "." + text
		}
		r, err := regexp.Compile(text)
		if err != nil {
			fmt.Println("error:", err, text)

		}
		//fmt.Println(r)
		regexps = append(regexps, r)

	}
}
func hashCode(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return "error", err
	}
	defer file.Close()
	h := sha1.New()
	_, erro := io.Copy(h, file)
	if erro != nil {
		fmt.Println(erro)
		return "error", erro
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

//获取指定目录及所有子目录下的所有文件。
func WalkDir(dirPth string) {
	var out *os.File
	var err1 error
	out, err1 = os.OpenFile("out.dat", os.O_CREATE, 0666) //打开文件,如文件不存在则创建

	if err1 != nil {
		fmt.Println("error:", err1)
		return
	}
	defer out.Close()

	filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() {
			return nil
		}
		for _, r := range regexps {
			if r.MatchString(fi.Name()) {
				return nil
			}
		}
		code, err := hashCode(filename)
		if err != nil {
			fmt.Println(err)
			return err
		}
		io.WriteString(out, filename+","+code+","+strconv.FormatInt(fi.Size(), 10)+"\n")

		return nil
	})

	//return files, err
}

func main() {

	WalkDir("D:\\temp")
}
