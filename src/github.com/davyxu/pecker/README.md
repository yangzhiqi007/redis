# pecker
服务器工具集

绕过ssh远程操作远程服务器

# 功能
## 开启服务器

```bash
    ./pecker -mode=server -addr=localhost:10901
```

## 发送文件
将本地的文件发送到远程服务器

```bash
    ./pecker -mode=sendfile FILE1 FILE2
```
注意: 发送文件放在pecker目录以 FILE1的文件名命名,例如: /path/to/a.txt 发过去将保存为a.txt

## 远程执行
将指令在远程机执行:
```bash
    ./pecker -addr=localhost:10901 CMD
```
CMD: 希望在远程执行的指令

注意: 执行权限为远程pecker服务器所在的用户的权限

如需批量执行,可使用:
```bash
    ./pecker -addr=localhost:10901 -cmdfile=CMDFILE
```

# 提示
本工具仅供调试及测试用. 使用本工具作为