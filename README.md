# video-annotation-tool

# 视频数据标注工具

## 使用

1. 下载zip文件并解压
2. 修改配置文件`config.json`
3. 运行当前平台的可执行文件

## 编译可执行文件

```bash
# Linux 下编译 Windows 版本
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build server.go

# Linux 下编译 MacOS 版本
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build server.go
```