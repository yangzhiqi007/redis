package client

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/davyxu/pecker/model"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Run() error {

	if *model.FlagCmdFile != "" {
		f, err := os.Open(*model.FlagCmdFile)
		if err != nil {
			return err
		}

		defer f.Close()

		_, err = ExecuteRemoteCommandText(*model.FlagAddr, *model.FlagCmdFile, f, "file")
		if err != nil {
			return err
		}

	} else if *model.FlagCmd != "" {

		_, err := ExecuteRemoteCommandText(*model.FlagAddr, *model.FlagCmd, strings.NewReader(*model.FlagCmd), "cmd")
		if err != nil {
			return err
		}
		// 标准管道输入
	} else if *model.FlagCmdPipe {

		if !isatty.IsTerminal(os.Stdin.Fd()) {

			data, err := ioutil.ReadAll(os.Stdin)

			if err != nil {
				return err
			}

			_, err = ExecuteRemoteCommandText(*model.FlagAddr, string(data), bytes.NewReader(data), "file")
			if err != nil {
				return err
			}

		} else {
			log.Errorf("tty is terminal")
		}
	}

	return nil
}

func addPrefixText(text io.Reader, prefix string) string {

	reader := bufio.NewScanner(text)

	reader.Split(bufio.ScanLines)

	var sb strings.Builder

	for reader.Scan() {
		sb.WriteString(prefix)
		sb.WriteString(reader.Text())
		sb.WriteString("\n")
	}

	return sb.String()
}

func ExecuteRemoteCommand(remoteAddr, text string) (string, error) {
	return ExecuteRemoteCommandText(remoteAddr, text, strings.NewReader(text), "cmd")
}

func ExecuteRemoteCommandText(remoteAddr, text string, reader io.Reader, shellMode string) (string, error) {
	log.Infof("[Exec] %s:\n%s%s", remoteAddr, model.LogIndent, text)

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/cmd", model.GConfig.GetAddress(remoteAddr)), reader)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Shell-Mode", shellMode)
	model.EncodeRequest(req)

	respond, err := http.DefaultClient.Do(req)

	if err != nil {

		if *model.FlagSkipError {
			return "", nil
		} else {
			return "", errors.Wrap(err, "ExecuteRemoteCommand.post")
		}

	}

	defer respond.Body.Close()

	repondText := addPrefixText(respond.Body, model.LogIndent)

	log.Infof("[Output] %s:\n%s", remoteAddr, repondText)

	errStr := respond.Header.Get("Error")

	if errStr != "" {
		if !*model.FlagSkipError {
			return repondText, errors.New(errStr)
		}

	}

	return repondText, nil
}
