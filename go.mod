module github.com/atomix/go-framework

go 1.12

require (
	cloud.google.com/go v0.43.0 // indirect
	github.com/atomix/api v0.0.0-20200202100958-13b24edbe32d
	github.com/atomix/atomix-go-node v0.0.0-20200114212450-178a2dc70336
	github.com/atomix/go-client v0.0.0-20200203180003-61799b5ca7c2
	github.com/atomix/go-local v0.0.0-20200202105028-743d224c66eb
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/google/pprof v0.0.0-20190723021845-34ac40c74b70 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/mobile v0.0.0-20190806162312-597adff16ade // indirect
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/sys v0.0.0-20190804053845-51ab0e2deafa // indirect
	golang.org/x/tools v0.0.0-20190806215303-88ddfcebc769 // indirect
	google.golang.org/grpc v1.23.1
	honnef.co/go/tools v0.0.1-2019.2.2 // indirect
)

replace github.com/atomix/api => ../atomix-api

replace github.com/atomix/go-client => ../atomix-go-client

replace github.com/atomix/go-local => ../atomix-go-local
