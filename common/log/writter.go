package log

import (
	"bufio"
	"fmt"
	"os"
)

type LogWriter struct {
	FileName string
	FileDir  string
	Buffer   []byte
}

var endOfLine = []byte("\n")

func (l *LogWriter) Write(str string) {
	l.Buffer = append(l.Buffer, []byte(str)...)
	// break the line
	l.Buffer = append(l.Buffer, endOfLine...)
}

func (l *LogWriter) Writef(format string, args ...interface{}) {
	l.Buffer = append(l.Buffer, []byte(fmt.Sprintf(format, args...))...)
	// break the line
	l.Buffer = append(l.Buffer, endOfLine...)
}

func (l *LogWriter) Save() (err error) {
	f, err := os.Create(l.FileDir + "/" + l.FileName)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	_, err = w.Write(l.Buffer)
	if err != nil {
		return
	}
	return
}
