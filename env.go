package glue

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var (
	EnvironmentPath = path.Base("")
	EnvironmentFile = ".env"
)

func EnvironmentFilePath() string {
	return path.Join(EnvironmentPath, EnvironmentFile)
}

func writeNewEnvironmentFileWith(key string) error {
	file, err := os.OpenFile(EnvironmentFilePath(), os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(bytes.NewReader(b))
	var out bytes.Buffer
	for s.Scan() {
		line := s.Text()
		switch {
		case len(line) == 0 || line[0] == '#': // too long line, empty, comment
			continue
		}
		if line == "APP_KEY=" {
			line += key
		}
		_, _ = fmt.Fprintln(&out, line)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if _, err := out.WriteTo(file); err != nil {
		return err
	}
	return nil
}
