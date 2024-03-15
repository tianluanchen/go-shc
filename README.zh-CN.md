# Go-SHC

[![GoVersion](https://img.shields.io/badge/Go-v1.20.2-blue?logo=Go&style=flat-square)](https://go.dev/)

混淆 Shell 脚本或打包 Shell 脚本为二进制程序，注意，混淆和打包仅在一定程度上保障你的代码隐密。

[English](./README.md) | 中文

## 构建

```bash
go build
```

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

## License

[GPL-3.0](./LICENSE) © Ayouth
