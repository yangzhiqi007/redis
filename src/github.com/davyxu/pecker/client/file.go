package client

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/davyxu/pecker/model"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// 将服务器的文件拉到客户端
func RecvFile(remoteAddr, remoteFileName, localDir string) error {

	req, err := http.NewRequest("post", fmt.Sprintf("http://%s/download", model.GConfig.GetAddress(remoteAddr)), nil)
	if err != nil {
		return errors.Wrap(err, "post error")
	}

	req.Header.Set("Content-Type", "application/binary")
	req.Header.Set("FileName", remoteFileName)

	model.EncodeRequest(req)

	respond, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}

	defer respond.Body.Close()

	errStr := respond.Header.Get("Error")
	if errStr != "" {
		return errors.Errorf("RecvFile error: %s", errStr)
	}

	nameOnly := filepath.Base(remoteFileName)

	localFullName := filepath.ToSlash(filepath.Join(localDir, nameOnly))

	os.MkdirAll(localDir, model.DefaultFilePerm)

	f, err := os.Create(localFullName)
	if err != nil {
		return errors.Wrap(err, "recvfile os.OpenFile error")
	}

	defer f.Close()

	zipreader, err := zlib.NewReader(respond.Body)
	if err != nil {
		return errors.Wrap(err, "recvfile zlib.NewReader(respond.Body) error")
	}

	defer zipreader.Close()

	io.Copy(f, zipreader)

	log.Debugf("[RecvFile] '%s' -> '%s'", remoteFileName, localDir)

	return nil
}

func getRemoteName() string {

	addr := model.GConfig.GetAddress(*model.FlagAddr)
	if addr == *model.FlagAddr {
		return *model.FlagAddr
	}

	return fmt.Sprintf("%s(%s)", addr, *model.FlagAddr)
}

// 将客户端本地的文件发到服务器
func SendFile(remoteAddr, localFileName, remoteDir string) error {
	log.Debugf("[SendFile] %s ...", localFileName)
	f, err := os.Open(localFileName)

	if err != nil {
		return errors.Wrap(err, "send file os.Open error")
	}

	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		return errors.Wrap(err, "file info error")
	}

	var compresedBytes bytes.Buffer
	zipWriter := zlib.NewWriter(&compresedBytes)

	// 拷贝到压缩器
	io.Copy(zipWriter, f)

	zipWriter.Flush()
	zipWriter.Close()

	req, err := http.NewRequest("post", fmt.Sprintf("http://%s/upload", model.GConfig.GetAddress(remoteAddr)), &compresedBytes)
	if err != nil {
		return errors.Wrap(err, "post error")
	}

	req.Header.Set("Content-Type", "application/binary")
	model.EncodeRequest(req)

	nameOnly := filepath.Base(localFileName)

	remoteFileName := filepath.ToSlash(filepath.Join(remoteDir, nameOnly))

	req.Header.Set("FileName", remoteFileName)
	req.Header.Set("FileSize", strconv.Itoa(int(finfo.Size())))

	respond, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}

	defer respond.Body.Close()

	errStr := respond.Header.Get("Error")
	if errStr != "" {
		return errors.Errorf("SendFile error: %s", errStr)
	} else {
		repondText := addPrefixText(respond.Body, model.LogIndent)
		if repondText != "" {
			log.Debugf("%s", repondText)
		}
	}

	log.Debugf("[SendFile] '%s' -> %s@%s", localFileName, remoteFileName, getRemoteName())

	return nil
}
