package generate

import (
	"bytes"
	"os"
	"path"

	"github.com/cortesi/modd/shell"
)

var (
	headerPound     = []byte("#k8s:generate ")
	headerSlash     = []byte("//k8s:generate ")
	headerKeepPound = []byte("#k8s:generate(keep) ")
	headerKeepSlash = []byte("//k8s:generate(keep) ")
)

func Generate(filename string, data []byte) ([]byte, error) {
	cmdLine, header, input := extractCommand(data)
	if cmdLine == "" {
		return data, nil
	}

	cmd, err := shell.Command("bash", cmdLine)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if len(header) > 0 {
		buf.Write(header)
		buf.WriteByte('\n')
	}

	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	cmd.Dir = path.Dir(filename)
	cmd.Env = append(os.Environ(),
		"K8S_GEN_FILE="+path.Base(filename),
	)

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func extractCommand(data []byte) (string, []byte, []byte) {
	var header []byte
	if bytes.HasPrefix(data, headerPound) {
		data = bytes.TrimPrefix(data, headerPound)
	} else if bytes.HasPrefix(data, headerSlash) {
		data = bytes.TrimPrefix(data, headerSlash)
	} else if bytes.HasPrefix(data, headerKeepPound) {
		header = data
		data = bytes.TrimPrefix(data, headerKeepPound)
	} else if bytes.HasPrefix(data, headerKeepSlash) {
		header = data
		data = bytes.TrimPrefix(data, headerKeepSlash)
	} else {
		return "", nil, data
	}

	cmdLine := data
	if idx := bytes.IndexRune(header, '\n'); idx >= 0 && len(header) > 0 {
		header = header[:idx]
	}
	if idx := bytes.IndexRune(data, '\n'); idx >= 0 {
		cmdLine = data[:idx]
		data = data[idx+1:]
	}

	return string(cmdLine), header, data
}
