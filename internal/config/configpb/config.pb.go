// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: config.proto

package configpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogLevel int32

const (
	LogLevel_LOG_LEVEL_UNSPECIFIED LogLevel = 0
	LogLevel_LOG_LEVEL_DEBUG       LogLevel = 1
	LogLevel_LOG_LEVEL_INFO        LogLevel = 2
	LogLevel_LOG_LEVEL_WARN        LogLevel = 3
	LogLevel_LOG_LEVEL_ERROR       LogLevel = 4
	LogLevel_LOG_LEVEL_FATAL       LogLevel = 5
	LogLevel_LOG_LEVEL_PANIC       LogLevel = 6
	LogLevel_LOG_LEVEL_DISABLED    LogLevel = 7
)

// Enum value maps for LogLevel.
var (
	LogLevel_name = map[int32]string{
		0: "LOG_LEVEL_UNSPECIFIED",
		1: "LOG_LEVEL_DEBUG",
		2: "LOG_LEVEL_INFO",
		3: "LOG_LEVEL_WARN",
		4: "LOG_LEVEL_ERROR",
		5: "LOG_LEVEL_FATAL",
		6: "LOG_LEVEL_PANIC",
		7: "LOG_LEVEL_DISABLED",
	}
	LogLevel_value = map[string]int32{
		"LOG_LEVEL_UNSPECIFIED": 0,
		"LOG_LEVEL_DEBUG":       1,
		"LOG_LEVEL_INFO":        2,
		"LOG_LEVEL_WARN":        3,
		"LOG_LEVEL_ERROR":       4,
		"LOG_LEVEL_FATAL":       5,
		"LOG_LEVEL_PANIC":       6,
		"LOG_LEVEL_DISABLED":    7,
	}
)

func (x LogLevel) Enum() *LogLevel {
	p := new(LogLevel)
	*p = x
	return p
}

func (x LogLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_config_proto_enumTypes[0].Descriptor()
}

func (LogLevel) Type() protoreflect.EnumType {
	return &file_config_proto_enumTypes[0]
}

func (x LogLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogLevel.Descriptor instead.
func (LogLevel) EnumDescriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{0}
}

type Net int32

const (
	Net_NET_UNSPECIFIED Net = 0
	Net_NET_UDP         Net = 1
	Net_NET_TCP         Net = 2
)

// Enum value maps for Net.
var (
	Net_name = map[int32]string{
		0: "NET_UNSPECIFIED",
		1: "NET_UDP",
		2: "NET_TCP",
	}
	Net_value = map[string]int32{
		"NET_UNSPECIFIED": 0,
		"NET_UDP":         1,
		"NET_TCP":         2,
	}
)

func (x Net) Enum() *Net {
	p := new(Net)
	*p = x
	return p
}

func (x Net) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Net) Descriptor() protoreflect.EnumDescriptor {
	return file_config_proto_enumTypes[1].Descriptor()
}

func (Net) Type() protoreflect.EnumType {
	return &file_config_proto_enumTypes[1]
}

func (x Net) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Net.Descriptor instead.
func (Net) EnumDescriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{1}
}

type Config struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	LogLevel      LogLevel               `protobuf:"varint,1,opt,name=log_level,json=logLevel,proto3,enum=configpb.LogLevel" json:"log_level,omitempty"`
	MaxRequests   uint32                 `protobuf:"varint,2,opt,name=max_requests,json=maxRequests,proto3" json:"max_requests,omitempty"`
	Listener      []*Listener            `protobuf:"bytes,10,rep,name=listener,proto3" json:"listener,omitempty"`
	Route         []*Route               `protobuf:"bytes,20,rep,name=route,proto3" json:"route,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Config) Reset() {
	*x = Config{}
	mi := &file_config_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetLogLevel() LogLevel {
	if x != nil {
		return x.LogLevel
	}
	return LogLevel_LOG_LEVEL_UNSPECIFIED
}

func (x *Config) GetMaxRequests() uint32 {
	if x != nil {
		return x.MaxRequests
	}
	return 0
}

func (x *Config) GetListener() []*Listener {
	if x != nil {
		return x.Listener
	}
	return nil
}

func (x *Config) GetRoute() []*Route {
	if x != nil {
		return x.Route
	}
	return nil
}

type Listener struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Listening address, e.g., "127.0.0.1:5553"
	Addr string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	// Network protocol: TCP or UDP
	Net Net `protobuf:"varint,2,opt,name=net,proto3,enum=configpb.Net" json:"net,omitempty"`
	// Request read timeout
	ReadTimeout *durationpb.Duration `protobuf:"bytes,3,opt,name=read_timeout,json=readTimeout,proto3" json:"read_timeout,omitempty"`
	// Response write timeout
	WriteTimeout  *durationpb.Duration `protobuf:"bytes,4,opt,name=write_timeout,json=writeTimeout,proto3" json:"write_timeout,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Listener) Reset() {
	*x = Listener{}
	mi := &file_config_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Listener) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Listener) ProtoMessage() {}

func (x *Listener) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Listener.ProtoReflect.Descriptor instead.
func (*Listener) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{1}
}

func (x *Listener) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Listener) GetNet() Net {
	if x != nil {
		return x.Net
	}
	return Net_NET_UNSPECIFIED
}

func (x *Listener) GetReadTimeout() *durationpb.Duration {
	if x != nil {
		return x.ReadTimeout
	}
	return nil
}

func (x *Listener) GetWriteTimeout() *durationpb.Duration {
	if x != nil {
		return x.WriteTimeout
	}
	return nil
}

type Route struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Route name, e.g., "semi-freedom"
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Finalize A/AAAA records (resolves CNAMEs, etc.)
	Finalize bool `protobuf:"varint,2,opt,name=finalize,proto3" json:"finalize,omitempty"`
	// DNS64 config to synthesize IPv6 addresses
	Dns64 *Dns64 `protobuf:"bytes,3,opt,name=dns64,proto3" json:"dns64,omitempty"`
	// Cache configuration
	Cache *Cache `protobuf:"bytes,4,opt,name=cache,proto3" json:"cache,omitempty"`
	// List of upstream resolvers
	Upstream []*Upstream `protobuf:"bytes,5,rep,name=upstream,proto3" json:"upstream,omitempty"`
	// Routing sources used for this route
	Source        []*Source `protobuf:"bytes,6,rep,name=source,proto3" json:"source,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Route) Reset() {
	*x = Route{}
	mi := &file_config_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Route) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Route) ProtoMessage() {}

func (x *Route) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Route.ProtoReflect.Descriptor instead.
func (*Route) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{2}
}

func (x *Route) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Route) GetFinalize() bool {
	if x != nil {
		return x.Finalize
	}
	return false
}

func (x *Route) GetDns64() *Dns64 {
	if x != nil {
		return x.Dns64
	}
	return nil
}

func (x *Route) GetCache() *Cache {
	if x != nil {
		return x.Cache
	}
	return nil
}

func (x *Route) GetUpstream() []*Upstream {
	if x != nil {
		return x.Upstream
	}
	return nil
}

func (x *Route) GetSource() []*Source {
	if x != nil {
		return x.Source
	}
	return nil
}

type Dns64 struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Prefix, e.g., "64:ff9b::/96"
	Prefix        string `protobuf:"bytes,1,opt,name=prefix,proto3" json:"prefix,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Dns64) Reset() {
	*x = Dns64{}
	mi := &file_config_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Dns64) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Dns64) ProtoMessage() {}

func (x *Dns64) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Dns64.ProtoReflect.Descriptor instead.
func (*Dns64) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{3}
}

func (x *Dns64) GetPrefix() string {
	if x != nil {
		return x.Prefix
	}
	return ""
}

type Cache struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MaxItems      uint32                 `protobuf:"varint,1,opt,name=max_items,json=maxItems,proto3" json:"max_items,omitempty"`
	MinTtl        uint32                 `protobuf:"varint,2,opt,name=min_ttl,json=minTtl,proto3" json:"min_ttl,omitempty"`
	MaxTtl        uint32                 `protobuf:"varint,3,opt,name=max_ttl,json=maxTtl,proto3" json:"max_ttl,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Cache) Reset() {
	*x = Cache{}
	mi := &file_config_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Cache) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cache) ProtoMessage() {}

func (x *Cache) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cache.ProtoReflect.Descriptor instead.
func (*Cache) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{4}
}

func (x *Cache) GetMaxItems() uint32 {
	if x != nil {
		return x.MaxItems
	}
	return 0
}

func (x *Cache) GetMinTtl() uint32 {
	if x != nil {
		return x.MinTtl
	}
	return 0
}

func (x *Cache) GetMaxTtl() uint32 {
	if x != nil {
		return x.MaxTtl
	}
	return 0
}

type Upstream struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Kind:
	//
	//	*Upstream_Udp
	//	*Upstream_Tcp
	//	*Upstream_Dot
	Kind          isUpstream_Kind `protobuf_oneof:"kind"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Upstream) Reset() {
	*x = Upstream{}
	mi := &file_config_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Upstream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Upstream) ProtoMessage() {}

func (x *Upstream) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Upstream.ProtoReflect.Descriptor instead.
func (*Upstream) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{5}
}

func (x *Upstream) GetKind() isUpstream_Kind {
	if x != nil {
		return x.Kind
	}
	return nil
}

func (x *Upstream) GetUdp() *UdpUpstream {
	if x != nil {
		if x, ok := x.Kind.(*Upstream_Udp); ok {
			return x.Udp
		}
	}
	return nil
}

func (x *Upstream) GetTcp() *TcpUpstream {
	if x != nil {
		if x, ok := x.Kind.(*Upstream_Tcp); ok {
			return x.Tcp
		}
	}
	return nil
}

func (x *Upstream) GetDot() *DotUpstream {
	if x != nil {
		if x, ok := x.Kind.(*Upstream_Dot); ok {
			return x.Dot
		}
	}
	return nil
}

type isUpstream_Kind interface {
	isUpstream_Kind()
}

type Upstream_Udp struct {
	// Plain UDP upstream
	Udp *UdpUpstream `protobuf:"bytes,10,opt,name=udp,proto3,oneof"`
}

type Upstream_Tcp struct {
	// Plain TCP upstream
	Tcp *TcpUpstream `protobuf:"bytes,20,opt,name=tcp,proto3,oneof"`
}

type Upstream_Dot struct {
	// DNS-over-TLS upstream
	Dot *DotUpstream `protobuf:"bytes,30,opt,name=dot,proto3,oneof"`
}

func (*Upstream_Udp) isUpstream_Kind() {}

func (*Upstream_Tcp) isUpstream_Kind() {}

func (*Upstream_Dot) isUpstream_Kind() {}

type UdpUpstream struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Address, e.g., "1.1.1.1:53" or "1.1.1.1"
	Addr string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	// Dial timeout
	DialTimeout *durationpb.Duration `protobuf:"bytes,2,opt,name=dial_timeout,json=dialTimeout,proto3" json:"dial_timeout,omitempty"`
	// Query timeout
	Timeout *durationpb.Duration `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// Disable UDP → TCP fallback
	NoTcpFallback bool `protobuf:"varint,4,opt,name=no_tcp_fallback,json=noTcpFallback,proto3" json:"no_tcp_fallback,omitempty"`
	// Network interface to use
	Iface         string `protobuf:"bytes,5,opt,name=iface,proto3" json:"iface,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UdpUpstream) Reset() {
	*x = UdpUpstream{}
	mi := &file_config_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UdpUpstream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UdpUpstream) ProtoMessage() {}

func (x *UdpUpstream) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UdpUpstream.ProtoReflect.Descriptor instead.
func (*UdpUpstream) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{6}
}

func (x *UdpUpstream) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *UdpUpstream) GetDialTimeout() *durationpb.Duration {
	if x != nil {
		return x.DialTimeout
	}
	return nil
}

func (x *UdpUpstream) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *UdpUpstream) GetNoTcpFallback() bool {
	if x != nil {
		return x.NoTcpFallback
	}
	return false
}

func (x *UdpUpstream) GetIface() string {
	if x != nil {
		return x.Iface
	}
	return ""
}

type TcpUpstream struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Address, e.g., "1.1.1.1:53" or "1.1.1.1"
	Addr string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	// Dial timeout
	DialTimeout *durationpb.Duration `protobuf:"bytes,2,opt,name=dial_timeout,json=dialTimeout,proto3" json:"dial_timeout,omitempty"`
	// Query timeout
	Timeout *durationpb.Duration `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// Network interface to use
	Iface string `protobuf:"bytes,5,opt,name=iface,proto3" json:"iface,omitempty"`
	// Connection pool config
	ConnPool      *ConnPool `protobuf:"bytes,6,opt,name=conn_pool,json=connPool,proto3" json:"conn_pool,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TcpUpstream) Reset() {
	*x = TcpUpstream{}
	mi := &file_config_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TcpUpstream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TcpUpstream) ProtoMessage() {}

func (x *TcpUpstream) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TcpUpstream.ProtoReflect.Descriptor instead.
func (*TcpUpstream) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{7}
}

func (x *TcpUpstream) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *TcpUpstream) GetDialTimeout() *durationpb.Duration {
	if x != nil {
		return x.DialTimeout
	}
	return nil
}

func (x *TcpUpstream) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *TcpUpstream) GetIface() string {
	if x != nil {
		return x.Iface
	}
	return ""
}

func (x *TcpUpstream) GetConnPool() *ConnPool {
	if x != nil {
		return x.ConnPool
	}
	return nil
}

type DotUpstream struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Address, e.g., "1.1.1.1:853" or "1.1.1.1"
	Addr string `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	// Dial timeout
	DialTimeout *durationpb.Duration `protobuf:"bytes,2,opt,name=dial_timeout,json=dialTimeout,proto3" json:"dial_timeout,omitempty"`
	// Query timeout
	Timeout *durationpb.Duration `protobuf:"bytes,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// Network interface to use
	Iface string `protobuf:"bytes,4,opt,name=iface,proto3" json:"iface,omitempty"`
	// TLS config
	Tls *TLS `protobuf:"bytes,5,opt,name=tls,proto3" json:"tls,omitempty"`
	// Connection pool config
	ConnPool      *ConnPool `protobuf:"bytes,6,opt,name=conn_pool,json=connPool,proto3" json:"conn_pool,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DotUpstream) Reset() {
	*x = DotUpstream{}
	mi := &file_config_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DotUpstream) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DotUpstream) ProtoMessage() {}

func (x *DotUpstream) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DotUpstream.ProtoReflect.Descriptor instead.
func (*DotUpstream) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{8}
}

func (x *DotUpstream) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *DotUpstream) GetDialTimeout() *durationpb.Duration {
	if x != nil {
		return x.DialTimeout
	}
	return nil
}

func (x *DotUpstream) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *DotUpstream) GetIface() string {
	if x != nil {
		return x.Iface
	}
	return ""
}

func (x *DotUpstream) GetTls() *TLS {
	if x != nil {
		return x.Tls
	}
	return nil
}

func (x *DotUpstream) GetConnPool() *ConnPool {
	if x != nil {
		return x.ConnPool
	}
	return nil
}

type TLS struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Path to the CA certificate file
	CaCert string `protobuf:"bytes,1,opt,name=ca_cert,json=caCert,proto3" json:"ca_cert,omitempty"`
	// Server name for TLS verification
	ServerName string `protobuf:"bytes,2,opt,name=server_name,json=serverName,proto3" json:"server_name,omitempty"`
	// Controls whether a client verifies the server's certificate chain and host name.
	InsecureSkipVerify bool `protobuf:"varint,3,opt,name=insecure_skip_verify,json=insecureSkipVerify,proto3" json:"insecure_skip_verify,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *TLS) Reset() {
	*x = TLS{}
	mi := &file_config_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TLS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLS) ProtoMessage() {}

func (x *TLS) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLS.ProtoReflect.Descriptor instead.
func (*TLS) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{9}
}

func (x *TLS) GetCaCert() string {
	if x != nil {
		return x.CaCert
	}
	return ""
}

func (x *TLS) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

func (x *TLS) GetInsecureSkipVerify() bool {
	if x != nil {
		return x.InsecureSkipVerify
	}
	return false
}

type ConnPool struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Maximum items in pool, <= 1 mean no pool
	MaxItems      int32 `protobuf:"varint,1,opt,name=max_items,json=maxItems,proto3" json:"max_items,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ConnPool) Reset() {
	*x = ConnPool{}
	mi := &file_config_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ConnPool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnPool) ProtoMessage() {}

func (x *ConnPool) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnPool.ProtoReflect.Descriptor instead.
func (*ConnPool) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{10}
}

func (x *ConnPool) GetMaxItems() int32 {
	if x != nil {
		return x.MaxItems
	}
	return 0
}

type Source struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Kind:
	//
	//	*Source_Static
	//	*Source_File
	Kind          isSource_Kind `protobuf_oneof:"kind"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Source) Reset() {
	*x = Source{}
	mi := &file_config_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{11}
}

func (x *Source) GetKind() isSource_Kind {
	if x != nil {
		return x.Kind
	}
	return nil
}

func (x *Source) GetStatic() *StaticSource {
	if x != nil {
		if x, ok := x.Kind.(*Source_Static); ok {
			return x.Static
		}
	}
	return nil
}

func (x *Source) GetFile() *FileSource {
	if x != nil {
		if x, ok := x.Kind.(*Source_File); ok {
			return x.File
		}
	}
	return nil
}

type isSource_Kind interface {
	isSource_Kind()
}

type Source_Static struct {
	Static *StaticSource `protobuf:"bytes,10,opt,name=static,proto3,oneof"`
}

type Source_File struct {
	File *FileSource `protobuf:"bytes,20,opt,name=file,proto3,oneof"`
}

func (*Source_Static) isSource_Kind() {}

func (*Source_File) isSource_Kind() {}

type StaticSource struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// List of domains to match
	Domain        []string `protobuf:"bytes,1,rep,name=domain,proto3" json:"domain,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StaticSource) Reset() {
	*x = StaticSource{}
	mi := &file_config_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StaticSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StaticSource) ProtoMessage() {}

func (x *StaticSource) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StaticSource.ProtoReflect.Descriptor instead.
func (*StaticSource) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{12}
}

func (x *StaticSource) GetDomain() []string {
	if x != nil {
		return x.Domain
	}
	return nil
}

type FileSource struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Path to the file to load
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// Reload interval (based on mtime and content)
	ReloadInterval *durationpb.Duration `protobuf:"bytes,4,opt,name=reload_interval,json=reloadInterval,proto3" json:"reload_interval,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *FileSource) Reset() {
	*x = FileSource{}
	mi := &file_config_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FileSource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileSource) ProtoMessage() {}

func (x *FileSource) ProtoReflect() protoreflect.Message {
	mi := &file_config_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileSource.ProtoReflect.Descriptor instead.
func (*FileSource) Descriptor() ([]byte, []int) {
	return file_config_proto_rawDescGZIP(), []int{13}
}

func (x *FileSource) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *FileSource) GetReloadInterval() *durationpb.Duration {
	if x != nil {
		return x.ReloadInterval
	}
	return nil
}

var File_config_proto protoreflect.FileDescriptor

const file_config_proto_rawDesc = "" +
	"\n" +
	"\fconfig.proto\x12\bconfigpb\x1a\x1egoogle/protobuf/duration.proto\"\xb3\x01\n" +
	"\x06Config\x12/\n" +
	"\tlog_level\x18\x01 \x01(\x0e2\x12.configpb.LogLevelR\blogLevel\x12!\n" +
	"\fmax_requests\x18\x02 \x01(\rR\vmaxRequests\x12.\n" +
	"\blistener\x18\n" +
	" \x03(\v2\x12.configpb.ListenerR\blistener\x12%\n" +
	"\x05route\x18\x14 \x03(\v2\x0f.configpb.RouteR\x05route\"\xbd\x01\n" +
	"\bListener\x12\x12\n" +
	"\x04addr\x18\x01 \x01(\tR\x04addr\x12\x1f\n" +
	"\x03net\x18\x02 \x01(\x0e2\r.configpb.NetR\x03net\x12<\n" +
	"\fread_timeout\x18\x03 \x01(\v2\x19.google.protobuf.DurationR\vreadTimeout\x12>\n" +
	"\rwrite_timeout\x18\x04 \x01(\v2\x19.google.protobuf.DurationR\fwriteTimeout\"\xdf\x01\n" +
	"\x05Route\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x1a\n" +
	"\bfinalize\x18\x02 \x01(\bR\bfinalize\x12%\n" +
	"\x05dns64\x18\x03 \x01(\v2\x0f.configpb.Dns64R\x05dns64\x12%\n" +
	"\x05cache\x18\x04 \x01(\v2\x0f.configpb.CacheR\x05cache\x12.\n" +
	"\bupstream\x18\x05 \x03(\v2\x12.configpb.UpstreamR\bupstream\x12(\n" +
	"\x06source\x18\x06 \x03(\v2\x10.configpb.SourceR\x06source\"\x1f\n" +
	"\x05Dns64\x12\x16\n" +
	"\x06prefix\x18\x01 \x01(\tR\x06prefix\"V\n" +
	"\x05Cache\x12\x1b\n" +
	"\tmax_items\x18\x01 \x01(\rR\bmaxItems\x12\x17\n" +
	"\amin_ttl\x18\x02 \x01(\rR\x06minTtl\x12\x17\n" +
	"\amax_ttl\x18\x03 \x01(\rR\x06maxTtl\"\x93\x01\n" +
	"\bUpstream\x12)\n" +
	"\x03udp\x18\n" +
	" \x01(\v2\x15.configpb.UdpUpstreamH\x00R\x03udp\x12)\n" +
	"\x03tcp\x18\x14 \x01(\v2\x15.configpb.TcpUpstreamH\x00R\x03tcp\x12)\n" +
	"\x03dot\x18\x1e \x01(\v2\x15.configpb.DotUpstreamH\x00R\x03dotB\x06\n" +
	"\x04kind\"\xd2\x01\n" +
	"\vUdpUpstream\x12\x12\n" +
	"\x04addr\x18\x01 \x01(\tR\x04addr\x12<\n" +
	"\fdial_timeout\x18\x02 \x01(\v2\x19.google.protobuf.DurationR\vdialTimeout\x123\n" +
	"\atimeout\x18\x03 \x01(\v2\x19.google.protobuf.DurationR\atimeout\x12&\n" +
	"\x0fno_tcp_fallback\x18\x04 \x01(\bR\rnoTcpFallback\x12\x14\n" +
	"\x05iface\x18\x05 \x01(\tR\x05iface\"\xdb\x01\n" +
	"\vTcpUpstream\x12\x12\n" +
	"\x04addr\x18\x01 \x01(\tR\x04addr\x12<\n" +
	"\fdial_timeout\x18\x02 \x01(\v2\x19.google.protobuf.DurationR\vdialTimeout\x123\n" +
	"\atimeout\x18\x03 \x01(\v2\x19.google.protobuf.DurationR\atimeout\x12\x14\n" +
	"\x05iface\x18\x05 \x01(\tR\x05iface\x12/\n" +
	"\tconn_pool\x18\x06 \x01(\v2\x12.configpb.ConnPoolR\bconnPool\"\xfc\x01\n" +
	"\vDotUpstream\x12\x12\n" +
	"\x04addr\x18\x01 \x01(\tR\x04addr\x12<\n" +
	"\fdial_timeout\x18\x02 \x01(\v2\x19.google.protobuf.DurationR\vdialTimeout\x123\n" +
	"\atimeout\x18\x03 \x01(\v2\x19.google.protobuf.DurationR\atimeout\x12\x14\n" +
	"\x05iface\x18\x04 \x01(\tR\x05iface\x12\x1f\n" +
	"\x03tls\x18\x05 \x01(\v2\r.configpb.TLSR\x03tls\x12/\n" +
	"\tconn_pool\x18\x06 \x01(\v2\x12.configpb.ConnPoolR\bconnPool\"q\n" +
	"\x03TLS\x12\x17\n" +
	"\aca_cert\x18\x01 \x01(\tR\x06caCert\x12\x1f\n" +
	"\vserver_name\x18\x02 \x01(\tR\n" +
	"serverName\x120\n" +
	"\x14insecure_skip_verify\x18\x03 \x01(\bR\x12insecureSkipVerify\"'\n" +
	"\bConnPool\x12\x1b\n" +
	"\tmax_items\x18\x01 \x01(\x05R\bmaxItems\"n\n" +
	"\x06Source\x120\n" +
	"\x06static\x18\n" +
	" \x01(\v2\x16.configpb.StaticSourceH\x00R\x06static\x12*\n" +
	"\x04file\x18\x14 \x01(\v2\x14.configpb.FileSourceH\x00R\x04fileB\x06\n" +
	"\x04kind\"&\n" +
	"\fStaticSource\x12\x16\n" +
	"\x06domain\x18\x01 \x03(\tR\x06domain\"d\n" +
	"\n" +
	"FileSource\x12\x12\n" +
	"\x04path\x18\x01 \x01(\tR\x04path\x12B\n" +
	"\x0freload_interval\x18\x04 \x01(\v2\x19.google.protobuf.DurationR\x0ereloadInterval*\xb9\x01\n" +
	"\bLogLevel\x12\x19\n" +
	"\x15LOG_LEVEL_UNSPECIFIED\x10\x00\x12\x13\n" +
	"\x0fLOG_LEVEL_DEBUG\x10\x01\x12\x12\n" +
	"\x0eLOG_LEVEL_INFO\x10\x02\x12\x12\n" +
	"\x0eLOG_LEVEL_WARN\x10\x03\x12\x13\n" +
	"\x0fLOG_LEVEL_ERROR\x10\x04\x12\x13\n" +
	"\x0fLOG_LEVEL_FATAL\x10\x05\x12\x13\n" +
	"\x0fLOG_LEVEL_PANIC\x10\x06\x12\x16\n" +
	"\x12LOG_LEVEL_DISABLED\x10\a*4\n" +
	"\x03Net\x12\x13\n" +
	"\x0fNET_UNSPECIFIED\x10\x00\x12\v\n" +
	"\aNET_UDP\x10\x01\x12\v\n" +
	"\aNET_TCP\x10\x02B6Z4github.com/buglloc/sly64/v2/internal/config/configpbb\x06proto3"

var (
	file_config_proto_rawDescOnce sync.Once
	file_config_proto_rawDescData []byte
)

func file_config_proto_rawDescGZIP() []byte {
	file_config_proto_rawDescOnce.Do(func() {
		file_config_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_config_proto_rawDesc), len(file_config_proto_rawDesc)))
	})
	return file_config_proto_rawDescData
}

var file_config_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_config_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_config_proto_goTypes = []any{
	(LogLevel)(0),               // 0: configpb.LogLevel
	(Net)(0),                    // 1: configpb.Net
	(*Config)(nil),              // 2: configpb.Config
	(*Listener)(nil),            // 3: configpb.Listener
	(*Route)(nil),               // 4: configpb.Route
	(*Dns64)(nil),               // 5: configpb.Dns64
	(*Cache)(nil),               // 6: configpb.Cache
	(*Upstream)(nil),            // 7: configpb.Upstream
	(*UdpUpstream)(nil),         // 8: configpb.UdpUpstream
	(*TcpUpstream)(nil),         // 9: configpb.TcpUpstream
	(*DotUpstream)(nil),         // 10: configpb.DotUpstream
	(*TLS)(nil),                 // 11: configpb.TLS
	(*ConnPool)(nil),            // 12: configpb.ConnPool
	(*Source)(nil),              // 13: configpb.Source
	(*StaticSource)(nil),        // 14: configpb.StaticSource
	(*FileSource)(nil),          // 15: configpb.FileSource
	(*durationpb.Duration)(nil), // 16: google.protobuf.Duration
}
var file_config_proto_depIdxs = []int32{
	0,  // 0: configpb.Config.log_level:type_name -> configpb.LogLevel
	3,  // 1: configpb.Config.listener:type_name -> configpb.Listener
	4,  // 2: configpb.Config.route:type_name -> configpb.Route
	1,  // 3: configpb.Listener.net:type_name -> configpb.Net
	16, // 4: configpb.Listener.read_timeout:type_name -> google.protobuf.Duration
	16, // 5: configpb.Listener.write_timeout:type_name -> google.protobuf.Duration
	5,  // 6: configpb.Route.dns64:type_name -> configpb.Dns64
	6,  // 7: configpb.Route.cache:type_name -> configpb.Cache
	7,  // 8: configpb.Route.upstream:type_name -> configpb.Upstream
	13, // 9: configpb.Route.source:type_name -> configpb.Source
	8,  // 10: configpb.Upstream.udp:type_name -> configpb.UdpUpstream
	9,  // 11: configpb.Upstream.tcp:type_name -> configpb.TcpUpstream
	10, // 12: configpb.Upstream.dot:type_name -> configpb.DotUpstream
	16, // 13: configpb.UdpUpstream.dial_timeout:type_name -> google.protobuf.Duration
	16, // 14: configpb.UdpUpstream.timeout:type_name -> google.protobuf.Duration
	16, // 15: configpb.TcpUpstream.dial_timeout:type_name -> google.protobuf.Duration
	16, // 16: configpb.TcpUpstream.timeout:type_name -> google.protobuf.Duration
	12, // 17: configpb.TcpUpstream.conn_pool:type_name -> configpb.ConnPool
	16, // 18: configpb.DotUpstream.dial_timeout:type_name -> google.protobuf.Duration
	16, // 19: configpb.DotUpstream.timeout:type_name -> google.protobuf.Duration
	11, // 20: configpb.DotUpstream.tls:type_name -> configpb.TLS
	12, // 21: configpb.DotUpstream.conn_pool:type_name -> configpb.ConnPool
	14, // 22: configpb.Source.static:type_name -> configpb.StaticSource
	15, // 23: configpb.Source.file:type_name -> configpb.FileSource
	16, // 24: configpb.FileSource.reload_interval:type_name -> google.protobuf.Duration
	25, // [25:25] is the sub-list for method output_type
	25, // [25:25] is the sub-list for method input_type
	25, // [25:25] is the sub-list for extension type_name
	25, // [25:25] is the sub-list for extension extendee
	0,  // [0:25] is the sub-list for field type_name
}

func init() { file_config_proto_init() }
func file_config_proto_init() {
	if File_config_proto != nil {
		return
	}
	file_config_proto_msgTypes[5].OneofWrappers = []any{
		(*Upstream_Udp)(nil),
		(*Upstream_Tcp)(nil),
		(*Upstream_Dot)(nil),
	}
	file_config_proto_msgTypes[11].OneofWrappers = []any{
		(*Source_Static)(nil),
		(*Source_File)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_config_proto_rawDesc), len(file_config_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_config_proto_goTypes,
		DependencyIndexes: file_config_proto_depIdxs,
		EnumInfos:         file_config_proto_enumTypes,
		MessageInfos:      file_config_proto_msgTypes,
	}.Build()
	File_config_proto = out.File
	file_config_proto_goTypes = nil
	file_config_proto_depIdxs = nil
}
