package model

import "flag"

var (
	FlagAddr       = flag.String("addr", ":7701", "server address to listen")
	FlagMode       = flag.String("mode", "client", "client/server to run")
	FlagCmdFile    = flag.String("cmdfile", "", "command from file")
	FlagCmd        = flag.String("cmd", "", "command to execute on remote")
	FlagCmdPipe    = flag.Bool("cmdpipe", false, "use unix pipe(|) to input command")
	FlagConfigFile = flag.String("cfgfile", "", "config file")
	FlagLogFile    = flag.String("log", "", "log to file")
	FlagVersion    = flag.Bool("version", false, "show version")
	FlagSkipError  = flag.Bool("skiperr", false, "skip execute error")
	FlagRetryTimes = flag.Int("retrytimes", 1, "retry times when err, affect recv/send file")
	FlagPrivateKey = flag.String("secrets", "", "private key to communicate")
)
