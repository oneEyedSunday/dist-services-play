Install  protocompiler

```sh
 wget https://github.com/protocolbuffers/protobuf/releases/download/v3.9.0/protoc-3.9.0-osx-x86_64.zip 
 unzip protoc-3.9.0-osx-x86_64.zip -d /usr/local/opt/protobuf
 ```

Install proto gen tools
```sh
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```
Ensure path is set

```sh
export GOPATH=$HOME/go
                                                        
export PATH=$PATH:$GOPATH/bin
```

Generate proto from defs
```sh
protoc api/v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.
```
