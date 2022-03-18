package server

import (
	"compress/zlib"
	"github.com/davyxu/pecker/model"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func serverRecvFile(request *http.Request) error {
	fileName := request.Header.Get("FileName")

	if fileName == "" {
		return errors.New("recvfile invalid filename")
	}

	log.Debugf("[RecvFile] %s -> %s", request.RemoteAddr, fileName)

	dirToCreate := filepath.Dir(fileName)
	os.MkdirAll(dirToCreate, model.DefaultFilePerm)

	f, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "recvfile os.OpenFile error")
	}

	defer f.Close()

	zipreader, err := zlib.NewReader(request.Body)
	if err != nil {
		return errors.Wrap(err, "recvfile copy to zlib.NewReader error")
	}

	defer zipreader.Close()

	io.Copy(f, zipreader)

	return nil
}

func serverSendFile(writer http.ResponseWriter, request *http.Request) error {
	fileName := request.Header.Get("FileName")

	if fileName == "" {
		return errors.New("sendfile invalid filename")
	}

	log.Debugf("[SendFile] %s --> %s", fileName, request.RemoteAddr)

	f, err := os.Open(fileName)
	if err != nil {
		return errors.Wrap(err, "sendfile os.OpenFile error")
	}

	defer f.Close()

	zipWriter := zlib.NewWriter(writer)

	io.Copy(zipWriter, f)

	zipWriter.Flush()

	zipWriter.Close()

	return nil
}

func onUpload(writer http.ResponseWriter, request *http.Request) {
	if !model.VerifyRequest(request) {
		return
	}

	err := serverRecvFile(request)
	if err != nil {
		log.Errorf("RecvFile failed: %s", err)
		writer.Header().Set("Error", err.Error())
	}
}

func onDownload(writer http.ResponseWriter, request *http.Request) {

	if !model.VerifyRequest(request) {
		return
	}

	err := serverSendFile(writer, request)
	if err != nil {
		log.Errorf("SendFile failed: %s", err)
		writer.Header().Set("Error", err.Error())
	}
}
