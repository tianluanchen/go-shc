# Go-SHC

Obfuscate Shell scripts or package them into binary programs. Note that obfuscation and packaging only provide a certain level of code concealment.

## Build

```bash
go build
```

## Usage

Obfuscation

```bash
# Print obfuscated result to standard output
go-shc --script 'echo hello;read a;echo $a'

# Obfuscate a.sh and b.sh, write the result to obfuscated.sh
go-shc -o obfuscated.sh a.sh b.sh
```

Packaging

```bash
# Scripts are encapsulated into the binary program 'app'
go-shc pack --script 'echo hello;read a;echo $a' -o app
# a.sh and b.sh are encapsulated into the binary program 'app'
go-shc pack -o app a.sh b.sh
```

Obfuscation and Packaging

```bash
# Obfuscate
go-shc -o obfuscated.sh a.sh b.sh
# Package
go-shc pack -o app obfuscated.sh
```

Enable the built-in web service

```bash
go-shc serve -a "127.0.0.1:8080"
```

## License

[GPL-3.0](./LICENSE)