package server

import (
	"github.com/davyxu/pecker/model"
	"net/http"
)

func Run() error {
	http.HandleFunc("/version", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(model.Version))
	})

	http.HandleFunc("/upload", onUpload)
	http.HandleFunc("/download", onDownload)

	runOnServer_cmd()

	log.Infof("listen at: '%s' ...\n", *model.FlagAddr)

	return http.ListenAndServe(*model.FlagAddr, nil) //设置监听的端口
}
