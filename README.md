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

将对应的可执行文件放入全局变量，比如我将其命名为ggit。

用法示例:

```text
╰─± ggit -h
Speed up the repo cloning from the github.com

Usage:
    ggit [flags]
ggit [command]

Available Commands:
clone       Clone the specified repo from github.com
completion  Generate the autocompletion script for the specified shell
help        Help about any command
version     Prints the version of ggit

Flags:
    -h, --help      help for ggit
                             -v, --version   Prints the version of ggit

Use "ggit [command] --help" for more information about a command.

-------------------------------------------------------
╰─± ggit clone -h  
Clone the specified repo from github.com

Usage:
  ggit clone <git repo url> [flags]

Flags:
  -h, --help                help for clone
  -o, --other stringArray   other sub-commands of clone-command in git. 
                            Wrap it in double quotation marks. 
                            eg. "--depth=1"
-------------------------------------------------------
```

下载示例：

```text
╰─± ggit clone https://github.com/mergestat/mergestat.git -o "--depth=1" -o "-v"
[folderAbsPath] /Users/dapeng/Desktop/code/Git/ggit/mergestat
######################## 🥳 Sort By Ping RTT Value 🥳 ########################
PING gitclone.com (47.96.130.35)
PING hub.fastgit.org (89.31.125.6)
PING github.com.cnpmjs.org (47.241.4.205)
PING github.wuyanzheshui.workers.dev (104.21.74.35)
gitclone.com done!
hub.fastgit.org done!
github.com.cnpmjs.org done!
github.wuyanzheshui.workers.dev done!
Sorted list: [{https://github.com.cnpmjs.org/ 0} {https://gitclone.com/ 40782499} {https://hub.fastgit.org/ 105558000} {https://github.wuyanzheshui.workers.dev/ 462600000}]
********************************************************************************
# Current mirror's url is:  https://github.com.cnpmjs.org/
Folder name: mergestat
Command: [/usr/local/bin/git clone https://github.com.cnpmjs.org/mergestat/mergestat.git --depth=1 -v]
----------------- CLONE -----------------
Cloning into 'mergestat'...
POST git-upload-pack (175 bytes)
POST git-upload-pack (229 bytes)
remote: Enumerating objects: 192, done.
remote: Counting objects: 100% (192/192), done.
remote: Compressing objects: 100% (165/165), done.
remote: Total 192 (delta 54), reused 67 (delta 14), pack-reused 0
Receiving objects: 100% (192/192), 2.68 MiB | 1.43 MiB/s, done.
Resolving deltas: 100% (54/54), done.
Clone done!!!
Command: [/usr/local/bin/git remote set-url origin https://github.com/mergestat/mergestat.git]
----------------- REMOTE -----------------
Set remote done!!!
All done!!!
```