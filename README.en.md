# go监管程序

#### 介绍
go监管程序
- 增加java类服务监控
- 增加docker类服务监控
- 增加top命令、或者htop命令
- 增加json配置
- 及时kill掉相关进程 [待开发]

#### 软件架构
1. 使用 go mod init main （main 模块名称）
2. 增加缺失的包，移除没用的包  go mod tidy

#### 安装教程
1.  打包
```
Mac 下编译 Linux 和 Windows 64位可执行程序
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

Linux 下编译 Mac 和 Windows 64位可执行程序
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

Windows 下编译 Mac 和 Linux 64位可执行程序
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go

GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
GOARCH：目标平台的体系架构（386、amd64、arm）
交叉编译不支持 CGO 所以要禁用它
```
1. 部署
``` 
liunx下，赋予脚本操作权限
chmod 777 main
后台执行,会在同级生成nohup.out
nohup ./main &

将输出结果重定向到/tmp/myprogram.log文件
nohup ./myprogram >/tmp/myprogram.log 2>&1 &
线上使用版本
nohup ./main >/home/goCode/go-regulatory-procedures/main_run.log 2>&1 &

查看进程占用或者htop
top -p 12811

htop 直接输入进程数字
```

- 设置代理路径

- 方法1 ：cmd -> 打开终端
  输入： go env -w GOPROXY=https://goproxy.cn,direct

- 方法2： 如果安装golang
  打开setting -> 选择Go -> Go Modules 回到下图这个界面
  勾选 Enable Go modules integration选项
  输入 -> GOPROXY=https://goproxy.cn  即可
	
