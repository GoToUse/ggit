# ggit

## 有go环境

> Linux和MacOS下使用

克隆该代码，并 `bash` 执行 `install.sh`，他会自动编译为可执行文件并注册到全局使用。

安装:

```bash
bash install.sh
```

然后就可以使用 `ggit clone https://github.com/xxx/xxx.git`


> windows下使用

克隆该代码，并编译为windows可执行文件，编译命令: 

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ggit.exe main.go
```

## 没有go环境

那么这里有现成编译好的文件，[下载地址](https://github.com/Abeautifulsnow/ggit/releases)

- windows/amd64
- linux/amd64
- darwin/amd64

将对应的可执行文件放入全局变量。

使用方式:

```nashorn js
[可执行文件路径] clone https://github.com/xxx/xxx.git
```
