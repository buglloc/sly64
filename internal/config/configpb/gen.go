package configpb

//go:generate protoc --go_out=. config.proto
//go:generate mv github.com/buglloc/sly64/v2/internal/config/configpb/config.pb.go .
//go:generate rmdir -p github.com/buglloc/sly64/v2/internal/config/configpb
