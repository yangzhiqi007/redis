package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/golog"
	"github.com/davyxu/pecker/client"
	"github.com/davyxu/pecker/model"
	"github.com/davyxu/pecker/server"
	"os"
	"time"
)

func main() {

	flag.Parse()

	if *model.FlagVersion {
		fmt.Println("version", model.Version)
		return
	}

	var err error

	switch *model.FlagMode {
	case "client", "cli":
		// pecker --mode=cli -cmd='ls'
		// pecker --mode=cli -cmdfile=x.sh
		// pecker --mode=cli -cmdpipe <<EOF
		//  your shell code
		// EOF

		if *model.FlagConfigFile != "" {
			err = model.GConfig.Load(*model.FlagConfigFile)
			if err != nil {
				goto ExitOnError
			}
		}

		err = client.Run()

	case "server", "sv":
		// pecker -mode=server -addr=localhost:10901

		if *model.FlagLogFile != "" {
			golog.SetOutputToFile(".", *model.FlagLogFile)
		}

		log.Infof("version: %s", model.Version)

		err = server.Run()
		// 发送文件到服务器目录
	case "sendfile": // pecker --mode=sendfile localfile remotedir

		var (
			localFile string
			remoteDir string
		)

		switch len(flag.Args()) {
		case 1:
			localFile = flag.Args()[0]
			remoteDir = "."
		case 2:
			localFile = flag.Args()[0]
			remoteDir = flag.Args()[1]
		default:
			fmt.Println("invalid arguments, pecker --mode=sendfile localfile remotedir")
			return
		}

		for i := 0; i < *model.FlagRetryTimes; i++ {
			if err = client.SendFile(*model.FlagAddr, localFile, remoteDir); err == nil {
				break
			} else {
				time.Sleep(1 * time.Second)
				log.Errorf("%s try times...%d", err.Error(), i+1)
			}
		}

		if err != nil {
			// 避免多报一次错误
			err = nil
			goto ExitOnError
		}

		// 从服务器目录拉取文件
	case "recvfile": // pecker --mode=recvfile remotefile localdir

		var (
			localDir   string
			remoteFile string
		)

		switch len(flag.Args()) {
		case 1:
			remoteFile = flag.Args()[0]
			localDir = "."
		case 2:
			remoteFile = flag.Args()[0]
			localDir = flag.Args()[1]
		default:
			log.Errorf("invalid arguments, pecker --mode=recvfile remotefile localdir")
			return
		}

		for i := 0; i < *model.FlagRetryTimes; i++ {
			if err = client.RecvFile(*model.FlagAddr, localDir, remoteFile); err == nil {
				break
			} else {
				log.Errorf("%s try times...%d", err.Error(), i+1)
			}
		}

		if err != nil {
			// 避免多报一次错误
			err = nil
			goto ExitOnError
		}
	default:
		log.Errorf("unknown mode %s", *model.FlagMode)

	}

	if err != nil {
		goto ExitOnError
	}

	return
ExitOnError:
	if err != nil {
		log.Errorf("%s", err.Error())
	}
	os.Exit(1)

}
