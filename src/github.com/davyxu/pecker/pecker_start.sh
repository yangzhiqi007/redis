#!/bin/bash


StartDaemon()
{
	if [ ! -x "./pecker" ]; then
		chmod +x ./pecker
	fi
	nohup ./pecker -mode=server -addr=:17701 >> pecker.log 2>&1 &
}

# 上传pecker二进制
#  ./pecker -mode=sendfile /path/to/new/pecker .

# 更新本shell
#  ./pecker -mode=sendfile /path/to/new/pecker_start .

# 更新
#  ./pecker -mode=client -cmd="sh ./pecker_start.sh upgrade" -skiperr

# 查看更新后的版本
#  ./pecker -mode=client -cmd="./pecker --version"


case "${1}" in  
daemon)  
	pkill pecker
	StartDaemon
    ;;
upgrade)
	if [ ! -f "upgrade/pecker" ]; then
		echo "upgrade not exists"
		exit 1
	fi
	
	echo "$(date '+%Y-%m-%d %H:%M:%S') Upgrade begin..."
	
	pkill pecker
	
	mv -f upgrade/pecker pecker	
	StartDaemon
	
	echo "$(date '+%Y-%m-%d %H:%M:%S') Upgrade done!"
	
esac  