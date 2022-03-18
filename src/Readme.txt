1.获取msgpack包依赖
	go get github.com/tinylib/msgp
	
	运行 go run main.go 检查下载依赖
	
	
2.model内加入注释
	//go:generate msgp

3.model下运行go.exe generate，自动生成代码
	
4.获取radix.v2

	go get github.com/mediocregopher/radix.v2
	
5.将davyxu包转入github.com文件里

6.运行redis
	
	
	