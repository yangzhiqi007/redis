package server

import (
	"bytes"
	"github.com/davyxu/pecker/model"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func convGBKToUTF8(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ExecFile(textReader io.Reader, src string, output *string) (err error) {

	var dir string

	rand.Seed(time.Now().Unix())

	dir, err = ioutil.TempDir("", "peckerexec_")
	defer os.Remove(dir)

	if err != nil {
		return
	}

	shellFile := filepath.Join(dir, "script.json")

	defer os.Remove(shellFile)

	data, err := ioutil.ReadAll(textReader)
	if err != nil {
		return err
	}
	ioutil.WriteFile(shellFile, data, 0666)

	log.Debugf("[ExecFile] %s:\n%s", src, string(data))

	cmd := exec.Command("bash", shellFile)

	var outBytes []byte
	outBytes, err = cmd.CombinedOutput()

	if output != nil {
		*output = string(outBytes)
	}

	log.Infof("[Output] %s:\n%s%s", src, model.LogIndent, string(outBytes))

	if err != nil {
		return errors.Wrap(err, "server.execFile")
	}

	return
}

func ExecCommand(text, src string, output *string) (err error) {

	log.Debugf("[Exec] %s:\n%s%s", src, model.LogIndent, text)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", text)
	} else {
		cmd = exec.Command("bash", "-c", text)
	}

	var outBytes []byte
	outBytes, err = cmd.CombinedOutput()

	if runtime.GOOS == "windows" {
		utfBytes, _ := convGBKToUTF8(outBytes)
		outBytes = utfBytes
	}

	if output != nil {
		*output = string(outBytes)
	} else {
		log.Infof("[Output] %s:\n%s%s", src, model.LogIndent, string(outBytes))
	}

	if err != nil {
		return errors.Wrap(err, "server.execCommand")
	}

	return
}

func executeServerCommand(commandReader io.Reader, src, shellMode string, output *string) (err error) {

	switch shellMode {
	case "", "cmd":
		data, err := ioutil.ReadAll(commandReader)
		if err != nil {
			return err
		}
		return ExecCommand(string(data), src, output)

	case "file":
		return ExecFile(commandReader, src, output)
	}

	return errors.New("unknown shell mode")

}

func runOnServer_cmd() {
	http.HandleFunc("/cmd", func(writer http.ResponseWriter, request *http.Request) {

		if !model.VerifyRequest(request) {
			return
		}

		shellMode := request.Header.Get("Shell-Mode")

		var output string
		err := executeServerCommand(request.Body, request.RemoteAddr, shellMode, &output)

		if err != nil {
			writer.Header().Set("Error", err.Error())
		}

		writer.Write([]byte(output))
	})
}
