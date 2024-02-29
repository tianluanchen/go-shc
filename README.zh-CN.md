# Go-SHC

[![GoVersion](https://img.shields.io/badge/Go-v1.20.2-blue?logo=Go&style=flat-square)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/tianluanchen/go-shc.svg)](https://pkg.go.dev/github.com/tianluanchen/go-shc)
[![Beta](https://img.shields.io/badge/-Beta-orange?style=flat-square)](./)

混淆 Shell 脚本或打包 Shell 脚本为二进制程序，注意，混淆和打包仅在一定程度上保障你的代码隐密。

[English](./README.md) | 中文

## 安装

```bash
go install github.com/tianluanchen/go-shc@latest
```

你也可以进入仓库的 Actions 界面，从 Artifact 中获取已编译好的程序

注意，想要将 shell 脚本打包成二进制程序，需要你的系统已经配置好 Go 的开发环境

## 使用

混淆

```bash
# 混淆结果打印到标准输出
go-shc --script 'echo hello;read a;echo $a'

# a.sh和b.sh混淆结果写入到 obfuscated.sh
go-shc -o obfuscated.sh a.sh b.sh
```

打包

```bash
# 脚本被封装到二进制程序app中
go-shc pack --script 'echo hello;read a;echo $a' -o app
# a.sh和b.sh被封装到二进制程序app中
go-shc pack  -o app a.sh b.sh
```

混淆并打包

```bash
# 混淆
go-shc -o obfuscated.sh a.sh b.sh
# 打包
go-shc pack -o app obfuscated.sh
```

启用内置 Web 服务

```bash
go-shc serve -a "127.0.0.1:8080"
```

## 导入库使用

下载

```bash
go get -u github.com/tianluanchen/go-shc@latest
```

自定义封装 Shell 脚本的 Go 代码生成器

```go
shc.PackShellScript(script, shc.PackOption{
    CustomGenerator func(script string) string {
        // ....
    }
})
```

## License

[GPL-3.0](./LICENSE) © Ayouth
