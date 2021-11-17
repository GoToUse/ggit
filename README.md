# ggit

> Linux和MacOS下使用

目前只支持Linux和MacOS系统。如果你有go环境，那么克隆该代码，并 `bash` 执行 `install.sh`.

安装:

```bash
bash install.sh
```

然后就可以使用 `ggit clone https://github.com/xxx/xxx.git`

> windows下使用

使用当前目录下的ggit.exe文件或者自己重新编译,编译命令: 

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ggit.exe main.go
```
