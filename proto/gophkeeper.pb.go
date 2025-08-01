// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: gophkeeper.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

// Регистрация и аутентификация
type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Login         string                 `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_gophkeeper_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type RegisterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_gophkeeper_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Login         string                 `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_gophkeeper_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{2}
}

func (x *LoginRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_gophkeeper_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{3}
}

func (x *LoginResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *LoginResponse) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

// Синхронизация и управление секретами
type Secret struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type          string                 `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"` // "login", "text", "binary", "card"
	Metadata      map[string]string      `protobuf:"bytes,3,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Data          []byte                 `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Version       int32                  `protobuf:"varint,5,opt,name=version,proto3" json:"version,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Secret) Reset() {
	*x = Secret{}
	mi := &file_gophkeeper_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Secret) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Secret) ProtoMessage() {}

func (x *Secret) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Secret.ProtoReflect.Descriptor instead.
func (*Secret) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{4}
}

func (x *Secret) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Secret) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Secret) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Secret) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Secret) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Secret) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type SyncRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	LocalSecrets  []*Secret              `protobuf:"bytes,2,rep,name=local_secrets,json=localSecrets,proto3" json:"local_secrets,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncRequest) Reset() {
	*x = SyncRequest{}
	mi := &file_gophkeeper_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncRequest) ProtoMessage() {}

func (x *SyncRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncRequest.ProtoReflect.Descriptor instead.
func (*SyncRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{5}
}

func (x *SyncRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *SyncRequest) GetLocalSecrets() []*Secret {
	if x != nil {
		return x.LocalSecrets
	}
	return nil
}

type SyncResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ServerSecrets []*Secret              `protobuf:"bytes,1,rep,name=server_secrets,json=serverSecrets,proto3" json:"server_secrets,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncResponse) Reset() {
	*x = SyncResponse{}
	mi := &file_gophkeeper_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncResponse) ProtoMessage() {}

func (x *SyncResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncResponse.ProtoReflect.Descriptor instead.
func (*SyncResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{6}
}

func (x *SyncResponse) GetServerSecrets() []*Secret {
	if x != nil {
		return x.ServerSecrets
	}
	return nil
}

type GetSecretRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	SecretId      string                 `protobuf:"bytes,2,opt,name=secret_id,json=secretId,proto3" json:"secret_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSecretRequest) Reset() {
	*x = GetSecretRequest{}
	mi := &file_gophkeeper_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSecretRequest) ProtoMessage() {}

func (x *GetSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSecretRequest.ProtoReflect.Descriptor instead.
func (*GetSecretRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{7}
}

func (x *GetSecretRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *GetSecretRequest) GetSecretId() string {
	if x != nil {
		return x.SecretId
	}
	return ""
}

type GetSecretResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Secret        *Secret                `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetSecretResponse) Reset() {
	*x = GetSecretResponse{}
	mi := &file_gophkeeper_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetSecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSecretResponse) ProtoMessage() {}

func (x *GetSecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSecretResponse.ProtoReflect.Descriptor instead.
func (*GetSecretResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{8}
}

func (x *GetSecretResponse) GetSecret() *Secret {
	if x != nil {
		return x.Secret
	}
	return nil
}

type UpdateSecretRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Secret        *Secret                `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateSecretRequest) Reset() {
	*x = UpdateSecretRequest{}
	mi := &file_gophkeeper_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSecretRequest) ProtoMessage() {}

func (x *UpdateSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSecretRequest.ProtoReflect.Descriptor instead.
func (*UpdateSecretRequest) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{9}
}

func (x *UpdateSecretRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *UpdateSecretRequest) GetSecret() *Secret {
	if x != nil {
		return x.Secret
	}
	return nil
}

type UpdateSecretResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateSecretResponse) Reset() {
	*x = UpdateSecretResponse{}
	mi := &file_gophkeeper_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateSecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateSecretResponse) ProtoMessage() {}

func (x *UpdateSecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gophkeeper_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateSecretResponse.ProtoReflect.Descriptor instead.
func (*UpdateSecretResponse) Descriptor() ([]byte, []int) {
	return file_gophkeeper_proto_rawDescGZIP(), []int{10}
}

func (x *UpdateSecretResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_gophkeeper_proto protoreflect.FileDescriptor

const file_gophkeeper_proto_rawDesc = "" +
	"\n" +
	"\x10gophkeeper.proto\x12\n" +
	"gophkeeper\x1a\x1fgoogle/protobuf/timestamp.proto\"C\n" +
	"\x0fRegisterRequest\x12\x14\n" +
	"\x05login\x18\x01 \x01(\tR\x05login\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"+\n" +
	"\x10RegisterResponse\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\"@\n" +
	"\fLoginRequest\x12\x14\n" +
	"\x05login\x18\x01 \x01(\tR\x05login\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\">\n" +
	"\rLoginResponse\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\tR\x06userId\"\x90\x02\n" +
	"\x06Secret\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04type\x18\x02 \x01(\tR\x04type\x12<\n" +
	"\bmetadata\x18\x03 \x03(\v2 .gophkeeper.Secret.MetadataEntryR\bmetadata\x12\x12\n" +
	"\x04data\x18\x04 \x01(\fR\x04data\x12\x18\n" +
	"\aversion\x18\x05 \x01(\x05R\aversion\x129\n" +
	"\n" +
	"updated_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\x1a;\n" +
	"\rMetadataEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\"\\\n" +
	"\vSyncRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x127\n" +
	"\rlocal_secrets\x18\x02 \x03(\v2\x12.gophkeeper.SecretR\flocalSecrets\"I\n" +
	"\fSyncResponse\x129\n" +
	"\x0eserver_secrets\x18\x01 \x03(\v2\x12.gophkeeper.SecretR\rserverSecrets\"E\n" +
	"\x10GetSecretRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12\x1b\n" +
	"\tsecret_id\x18\x02 \x01(\tR\bsecretId\"?\n" +
	"\x11GetSecretResponse\x12*\n" +
	"\x06secret\x18\x01 \x01(\v2\x12.gophkeeper.SecretR\x06secret\"W\n" +
	"\x13UpdateSecretRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12*\n" +
	"\x06secret\x18\x02 \x01(\v2\x12.gophkeeper.SecretR\x06secret\"0\n" +
	"\x14UpdateSecretResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess2\xe9\x02\n" +
	"\n" +
	"Gophkeeper\x12E\n" +
	"\bRegister\x12\x1b.gophkeeper.RegisterRequest\x1a\x1c.gophkeeper.RegisterResponse\x12<\n" +
	"\x05Login\x12\x18.gophkeeper.LoginRequest\x1a\x19.gophkeeper.LoginResponse\x129\n" +
	"\x04Sync\x12\x17.gophkeeper.SyncRequest\x1a\x18.gophkeeper.SyncResponse\x12H\n" +
	"\tGetSecret\x12\x1c.gophkeeper.GetSecretRequest\x1a\x1d.gophkeeper.GetSecretResponse\x12Q\n" +
	"\fUpdateSecret\x12\x1f.gophkeeper.UpdateSecretRequest\x1a .gophkeeper.UpdateSecretResponseB;Z9github.com/darkseear/gophkeeper/server/internal/api/protob\x06proto3"

var (
	file_gophkeeper_proto_rawDescOnce sync.Once
	file_gophkeeper_proto_rawDescData []byte
)

func file_gophkeeper_proto_rawDescGZIP() []byte {
	file_gophkeeper_proto_rawDescOnce.Do(func() {
		file_gophkeeper_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_gophkeeper_proto_rawDesc), len(file_gophkeeper_proto_rawDesc)))
	})
	return file_gophkeeper_proto_rawDescData
}

var file_gophkeeper_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_gophkeeper_proto_goTypes = []any{
	(*RegisterRequest)(nil),       // 0: gophkeeper.RegisterRequest
	(*RegisterResponse)(nil),      // 1: gophkeeper.RegisterResponse
	(*LoginRequest)(nil),          // 2: gophkeeper.LoginRequest
	(*LoginResponse)(nil),         // 3: gophkeeper.LoginResponse
	(*Secret)(nil),                // 4: gophkeeper.Secret
	(*SyncRequest)(nil),           // 5: gophkeeper.SyncRequest
	(*SyncResponse)(nil),          // 6: gophkeeper.SyncResponse
	(*GetSecretRequest)(nil),      // 7: gophkeeper.GetSecretRequest
	(*GetSecretResponse)(nil),     // 8: gophkeeper.GetSecretResponse
	(*UpdateSecretRequest)(nil),   // 9: gophkeeper.UpdateSecretRequest
	(*UpdateSecretResponse)(nil),  // 10: gophkeeper.UpdateSecretResponse
	nil,                           // 11: gophkeeper.Secret.MetadataEntry
	(*timestamppb.Timestamp)(nil), // 12: google.protobuf.Timestamp
}
var file_gophkeeper_proto_depIdxs = []int32{
	11, // 0: gophkeeper.Secret.metadata:type_name -> gophkeeper.Secret.MetadataEntry
	12, // 1: gophkeeper.Secret.updated_at:type_name -> google.protobuf.Timestamp
	4,  // 2: gophkeeper.SyncRequest.local_secrets:type_name -> gophkeeper.Secret
	4,  // 3: gophkeeper.SyncResponse.server_secrets:type_name -> gophkeeper.Secret
	4,  // 4: gophkeeper.GetSecretResponse.secret:type_name -> gophkeeper.Secret
	4,  // 5: gophkeeper.UpdateSecretRequest.secret:type_name -> gophkeeper.Secret
	0,  // 6: gophkeeper.Gophkeeper.Register:input_type -> gophkeeper.RegisterRequest
	2,  // 7: gophkeeper.Gophkeeper.Login:input_type -> gophkeeper.LoginRequest
	5,  // 8: gophkeeper.Gophkeeper.Sync:input_type -> gophkeeper.SyncRequest
	7,  // 9: gophkeeper.Gophkeeper.GetSecret:input_type -> gophkeeper.GetSecretRequest
	9,  // 10: gophkeeper.Gophkeeper.UpdateSecret:input_type -> gophkeeper.UpdateSecretRequest
	1,  // 11: gophkeeper.Gophkeeper.Register:output_type -> gophkeeper.RegisterResponse
	3,  // 12: gophkeeper.Gophkeeper.Login:output_type -> gophkeeper.LoginResponse
	6,  // 13: gophkeeper.Gophkeeper.Sync:output_type -> gophkeeper.SyncResponse
	8,  // 14: gophkeeper.Gophkeeper.GetSecret:output_type -> gophkeeper.GetSecretResponse
	10, // 15: gophkeeper.Gophkeeper.UpdateSecret:output_type -> gophkeeper.UpdateSecretResponse
	11, // [11:16] is the sub-list for method output_type
	6,  // [6:11] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_gophkeeper_proto_init() }
func file_gophkeeper_proto_init() {
	if File_gophkeeper_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_gophkeeper_proto_rawDesc), len(file_gophkeeper_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gophkeeper_proto_goTypes,
		DependencyIndexes: file_gophkeeper_proto_depIdxs,
		MessageInfos:      file_gophkeeper_proto_msgTypes,
	}.Build()
	File_gophkeeper_proto = out.File
	file_gophkeeper_proto_goTypes = nil
	file_gophkeeper_proto_depIdxs = nil
}
