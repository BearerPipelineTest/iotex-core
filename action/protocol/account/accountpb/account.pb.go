// Copyright (c) 2022 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

// To compile the proto, run:
//      protoc --go_out=plugins=grpc:. *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.20.1
// source: account.proto

package accountpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AccountType int32

const (
	AccountType_DEFAULT    AccountType = 0
	AccountType_ZERO_NONCE AccountType = 1
)

// Enum value maps for AccountType.
var (
	AccountType_name = map[int32]string{
		0: "DEFAULT",
		1: "ZERO_NONCE",
	}
	AccountType_value = map[string]int32{
		"DEFAULT":    0,
		"ZERO_NONCE": 1,
	}
)

func (x AccountType) Enum() *AccountType {
	p := new(AccountType)
	*p = x
	return p
}

func (x AccountType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccountType) Descriptor() protoreflect.EnumDescriptor {
	return file_account_proto_enumTypes[0].Descriptor()
}

func (AccountType) Type() protoreflect.EnumType {
	return &file_account_proto_enumTypes[0]
}

func (x AccountType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccountType.Descriptor instead.
func (AccountType) EnumDescriptor() ([]byte, []int) {
	return file_account_proto_rawDescGZIP(), []int{0}
}

type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// used by state-based model
	Nonce        uint64      `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Balance      string      `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	Root         []byte      `protobuf:"bytes,3,opt,name=root,proto3" json:"root,omitempty"`
	CodeHash     []byte      `protobuf:"bytes,4,opt,name=codeHash,proto3" json:"codeHash,omitempty"`
	IsCandidate  bool        `protobuf:"varint,5,opt,name=isCandidate,proto3" json:"isCandidate,omitempty"`
	VotingWeight []byte      `protobuf:"bytes,6,opt,name=votingWeight,proto3" json:"votingWeight,omitempty"`
	Type         AccountType `protobuf:"varint,7,opt,name=type,proto3,enum=accountpb.AccountType" json:"type,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_account_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetNonce() uint64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *Account) GetBalance() string {
	if x != nil {
		return x.Balance
	}
	return ""
}

func (x *Account) GetRoot() []byte {
	if x != nil {
		return x.Root
	}
	return nil
}

func (x *Account) GetCodeHash() []byte {
	if x != nil {
		return x.CodeHash
	}
	return nil
}

func (x *Account) GetIsCandidate() bool {
	if x != nil {
		return x.IsCandidate
	}
	return false
}

func (x *Account) GetVotingWeight() []byte {
	if x != nil {
		return x.VotingWeight
	}
	return nil
}

func (x *Account) GetType() AccountType {
	if x != nil {
		return x.Type
	}
	return AccountType_DEFAULT
}

var File_account_proto protoreflect.FileDescriptor

var file_account_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x22, 0xdb, 0x01, 0x0a, 0x07, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6f, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x72, 0x6f, 0x6f, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f,
	0x64, 0x65, 0x48, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x6f,
	0x64, 0x65, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x73, 0x43, 0x61, 0x6e, 0x64,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x43,
	0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x76, 0x6f, 0x74, 0x69,
	0x6e, 0x67, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c,
	0x76, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x2a, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x2a, 0x2a, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x46, 0x41, 0x55,
	0x4c, 0x54, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x5a, 0x45, 0x52, 0x4f, 0x5f, 0x4e, 0x4f, 0x4e,
	0x43, 0x45, 0x10, 0x01, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x69, 0x6f, 0x74, 0x65, 0x78, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x69, 0x6f, 0x74, 0x65, 0x78, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_account_proto_rawDescOnce sync.Once
	file_account_proto_rawDescData = file_account_proto_rawDesc
)

func file_account_proto_rawDescGZIP() []byte {
	file_account_proto_rawDescOnce.Do(func() {
		file_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_account_proto_rawDescData)
	})
	return file_account_proto_rawDescData
}

var file_account_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_account_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_account_proto_goTypes = []interface{}{
	(AccountType)(0), // 0: accountpb.AccountType
	(*Account)(nil),  // 1: accountpb.Account
}
var file_account_proto_depIdxs = []int32{
	0, // 0: accountpb.Account.type:type_name -> accountpb.AccountType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_account_proto_init() }
func file_account_proto_init() {
	if File_account_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_account_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_account_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_account_proto_goTypes,
		DependencyIndexes: file_account_proto_depIdxs,
		EnumInfos:         file_account_proto_enumTypes,
		MessageInfos:      file_account_proto_msgTypes,
	}.Build()
	File_account_proto = out.File
	file_account_proto_rawDesc = nil
	file_account_proto_goTypes = nil
	file_account_proto_depIdxs = nil
}
