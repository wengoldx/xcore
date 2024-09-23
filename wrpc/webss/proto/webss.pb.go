// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.4
// source: proto/webss.proto

package proto

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

// Empty response for Webss
type WEmpty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *WEmpty) Reset() {
	*x = WEmpty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WEmpty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WEmpty) ProtoMessage() {}

func (x *WEmpty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WEmpty.ProtoReflect.Descriptor instead.
func (*WEmpty) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{0}
}

// Unique id to indicate lifecycle
type ID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ID) Reset() {
	*x = ID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ID) ProtoMessage() {}

func (x *ID) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ID.ProtoReflect.Descriptor instead.
func (*ID) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{1}
}

func (x *ID) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Object file path releative bucket
type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{2}
}

func (x *File) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

// Multiple object files of indicate bucket
type Files struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bucket string   `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Files  []string `protobuf:"bytes,2,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *Files) Reset() {
	*x = Files{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Files) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Files) ProtoMessage() {}

func (x *Files) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Files.ProtoReflect.Descriptor instead.
func (*Files) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{3}
}

func (x *Files) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *Files) GetFiles() []string {
	if x != nil {
		return x.Files
	}
	return nil
}

// Lifecycle status
type Life struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Bucket string `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Life) Reset() {
	*x = Life{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Life) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Life) ProtoMessage() {}

func (x *Life) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Life.ProtoReflect.Descriptor instead.
func (*Life) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{4}
}

func (x *Life) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Life) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *Life) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

// Lifesycle ids for delete
type Lifes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids    []string `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	Bucket string   `protobuf:"bytes,2,opt,name=bucket,proto3" json:"bucket,omitempty"`
}

func (x *Lifes) Reset() {
	*x = Lifes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Lifes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Lifes) ProtoMessage() {}

func (x *Lifes) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Lifes.ProtoReflect.Descriptor instead.
func (*Lifes) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{5}
}

func (x *Lifes) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *Lifes) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

// Object file tag infos
type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bucket string   `protobuf:"bytes,1,opt,name=bucket,proto3" json:"bucket,omitempty"`
	Paths  []string `protobuf:"bytes,2,rep,name=paths,proto3" json:"paths,omitempty"`
	Status string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{6}
}

func (x *Tag) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

func (x *Tag) GetPaths() []string {
	if x != nil {
		return x.Paths
	}
	return nil
}

func (x *Tag) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

// Sign datas to get upload sign url
type Sign struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Res    string `protobuf:"bytes,1,opt,name=res,proto3" json:"res,omitempty"`
	Add    string `protobuf:"bytes,2,opt,name=add,proto3" json:"add,omitempty"`
	Suffix string `protobuf:"bytes,3,opt,name=suffix,proto3" json:"suffix,omitempty"`
}

func (x *Sign) Reset() {
	*x = Sign{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Sign) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Sign) ProtoMessage() {}

func (x *Sign) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Sign.ProtoReflect.Descriptor instead.
func (*Sign) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{7}
}

func (x *Sign) GetRes() string {
	if x != nil {
		return x.Res
	}
	return ""
}

func (x *Sign) GetAdd() string {
	if x != nil {
		return x.Add
	}
	return ""
}

func (x *Sign) GetSuffix() string {
	if x != nil {
		return x.Suffix
	}
	return ""
}

// Multiple sign urls to upload
type Signs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Res      string   `protobuf:"bytes,1,opt,name=res,proto3" json:"res,omitempty"`
	Add      string   `protobuf:"bytes,2,opt,name=add,proto3" json:"add,omitempty"`
	Suffixes []string `protobuf:"bytes,3,rep,name=suffixes,proto3" json:"suffixes,omitempty"`
}

func (x *Signs) Reset() {
	*x = Signs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Signs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Signs) ProtoMessage() {}

func (x *Signs) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Signs.ProtoReflect.Descriptor instead.
func (*Signs) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{8}
}

func (x *Signs) GetRes() string {
	if x != nil {
		return x.Res
	}
	return ""
}

func (x *Signs) GetAdd() string {
	if x != nil {
		return x.Add
	}
	return ""
}

func (x *Signs) GetSuffixes() []string {
	if x != nil {
		return x.Suffixes
	}
	return nil
}

// Sign url to upload
type SignUrl struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url  string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Path string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *SignUrl) Reset() {
	*x = SignUrl{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignUrl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignUrl) ProtoMessage() {}

func (x *SignUrl) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignUrl.ProtoReflect.Descriptor instead.
func (*SignUrl) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{9}
}

func (x *SignUrl) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *SignUrl) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

// Multiple sign urls to upload
type SignUrls struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls []*SignUrl `protobuf:"bytes,1,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *SignUrls) Reset() {
	*x = SignUrls{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignUrls) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignUrls) ProtoMessage() {}

func (x *SignUrls) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignUrls.ProtoReflect.Descriptor instead.
func (*SignUrls) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{10}
}

func (x *SignUrls) GetUrls() []*SignUrl {
	if x != nil {
		return x.Urls
	}
	return nil
}

// Multiple original file names
type FNames struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Res   string   `protobuf:"bytes,1,opt,name=res,proto3" json:"res,omitempty"`
	Add   string   `protobuf:"bytes,2,opt,name=add,proto3" json:"add,omitempty"`
	Names []*FName `protobuf:"bytes,3,rep,name=names,proto3" json:"names,omitempty"`
}

func (x *FNames) Reset() {
	*x = FNames{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FNames) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FNames) ProtoMessage() {}

func (x *FNames) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FNames.ProtoReflect.Descriptor instead.
func (*FNames) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{11}
}

func (x *FNames) GetRes() string {
	if x != nil {
		return x.Res
	}
	return ""
}

func (x *FNames) GetAdd() string {
	if x != nil {
		return x.Add
	}
	return ""
}

func (x *FNames) GetNames() []*FName {
	if x != nil {
		return x.Names
	}
	return nil
}

// Original uploaded file name
type FName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Suffix string `protobuf:"bytes,2,opt,name=suffix,proto3" json:"suffix,omitempty"`
}

func (x *FName) Reset() {
	*x = FName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FName) ProtoMessage() {}

func (x *FName) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FName.ProtoReflect.Descriptor instead.
func (*FName) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{12}
}

func (x *FName) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FName) GetSuffix() string {
	if x != nil {
		return x.Suffix
	}
	return ""
}

// Object file infos
type Info struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Etag string `protobuf:"bytes,2,opt,name=etag,proto3" json:"etag,omitempty"`
	Last string `protobuf:"bytes,3,opt,name=last,proto3" json:"last,omitempty"`
	Size int64  `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
}

func (x *Info) Reset() {
	*x = Info{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_webss_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Info) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Info) ProtoMessage() {}

func (x *Info) ProtoReflect() protoreflect.Message {
	mi := &file_proto_webss_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Info.ProtoReflect.Descriptor instead.
func (*Info) Descriptor() ([]byte, []int) {
	return file_proto_webss_proto_rawDescGZIP(), []int{13}
}

func (x *Info) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Info) GetEtag() string {
	if x != nil {
		return x.Etag
	}
	return ""
}

func (x *Info) GetLast() string {
	if x != nil {
		return x.Last
	}
	return ""
}

func (x *Info) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

var File_proto_webss_proto protoreflect.FileDescriptor

var file_proto_webss_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x65, 0x62, 0x73, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x08, 0x0a, 0x06, 0x57, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x1a, 0x0a, 0x04, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x35, 0x0a, 0x05, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x46, 0x0a,
	0x04, 0x4c, 0x69, 0x66, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x31, 0x0a, 0x05, 0x4c, 0x69, 0x66, 0x65, 0x73, 0x12, 0x10,
	0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x4b, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12,
	0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x70, 0x61, 0x74, 0x68, 0x73, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x42, 0x0a, 0x04, 0x53, 0x69, 0x67, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x61, 0x64, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x64,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x22, 0x47, 0x0a, 0x05, 0x53, 0x69, 0x67,
	0x6e, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x72, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x64, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x61, 0x64, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78,
	0x65, 0x73, 0x22, 0x2f, 0x0a, 0x07, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x22, 0x2e, 0x0a, 0x08, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x72, 0x6c, 0x73, 0x12,
	0x22, 0x0a, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x72, 0x6c, 0x52, 0x04, 0x75,
	0x72, 0x6c, 0x73, 0x22, 0x50, 0x0a, 0x06, 0x46, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x10, 0x0a,
	0x03, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x72, 0x65, 0x73, 0x12,
	0x10, 0x0a, 0x03, 0x61, 0x64, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x61, 0x64,
	0x64, 0x12, 0x22, 0x0a, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x05,
	0x6e, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x33, 0x0a, 0x05, 0x46, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x75, 0x66, 0x66, 0x69, 0x78, 0x22, 0x56, 0x0a, 0x04, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x74, 0x61, 0x67, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x74, 0x61, 0x67, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x61,
	0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x61, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69,
	0x7a, 0x65, 0x32, 0xec, 0x02, 0x0a, 0x05, 0x57, 0x65, 0x62, 0x73, 0x73, 0x12, 0x2a, 0x0a, 0x0b,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x0c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x57, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2b, 0x0a, 0x0d, 0x53, 0x65, 0x74, 0x42,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x4c, 0x69, 0x66, 0x65, 0x12, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4c, 0x69, 0x66, 0x65, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2c, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x42, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x4c, 0x69, 0x66, 0x65, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4c,
	0x69, 0x66, 0x65, 0x73, 0x1a, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x12, 0x28, 0x0a, 0x0b, 0x53, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x69,
	0x66, 0x65, 0x12, 0x0a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x67, 0x1a, 0x0d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2a, 0x0a,
	0x0b, 0x53, 0x69, 0x67, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x0b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x1a, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x2d, 0x0a, 0x0c, 0x53, 0x69, 0x67,
	0x6e, 0x46, 0x69, 0x6c, 0x65, 0x55, 0x72, 0x6c, 0x73, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x73, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x69, 0x67, 0x6e, 0x55, 0x72, 0x6c, 0x73, 0x12, 0x2e, 0x0a, 0x0c, 0x4f, 0x72, 0x69, 0x67,
	0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c, 0x73, 0x12, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x46, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x69, 0x67, 0x6e, 0x55, 0x72, 0x6c, 0x73, 0x12, 0x27, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x46,
	0x69, 0x6c, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x46, 0x69, 0x6c, 0x65, 0x1a, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x49, 0x6e, 0x66,
	0x6f, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_webss_proto_rawDescOnce sync.Once
	file_proto_webss_proto_rawDescData = file_proto_webss_proto_rawDesc
)

func file_proto_webss_proto_rawDescGZIP() []byte {
	file_proto_webss_proto_rawDescOnce.Do(func() {
		file_proto_webss_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_webss_proto_rawDescData)
	})
	return file_proto_webss_proto_rawDescData
}

var file_proto_webss_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_proto_webss_proto_goTypes = []interface{}{
	(*WEmpty)(nil),   // 0: proto.WEmpty
	(*ID)(nil),       // 1: proto.ID
	(*File)(nil),     // 2: proto.File
	(*Files)(nil),    // 3: proto.Files
	(*Life)(nil),     // 4: proto.Life
	(*Lifes)(nil),    // 5: proto.Lifes
	(*Tag)(nil),      // 6: proto.Tag
	(*Sign)(nil),     // 7: proto.Sign
	(*Signs)(nil),    // 8: proto.Signs
	(*SignUrl)(nil),  // 9: proto.SignUrl
	(*SignUrls)(nil), // 10: proto.SignUrls
	(*FNames)(nil),   // 11: proto.FNames
	(*FName)(nil),    // 12: proto.FName
	(*Info)(nil),     // 13: proto.Info
}
var file_proto_webss_proto_depIdxs = []int32{
	9,  // 0: proto.SignUrls.urls:type_name -> proto.SignUrl
	12, // 1: proto.FNames.names:type_name -> proto.FName
	3,  // 2: proto.Webss.DeleteFiles:input_type -> proto.Files
	4,  // 3: proto.Webss.SetBucketLife:input_type -> proto.Life
	5,  // 4: proto.Webss.DelBucketLife:input_type -> proto.Lifes
	6,  // 5: proto.Webss.SetFileLife:input_type -> proto.Tag
	7,  // 6: proto.Webss.SignFileUrl:input_type -> proto.Sign
	8,  // 7: proto.Webss.SignFileUrls:input_type -> proto.Signs
	11, // 8: proto.Webss.OriginalUrls:input_type -> proto.FNames
	2,  // 9: proto.Webss.GetFileInfo:input_type -> proto.File
	0,  // 10: proto.Webss.DeleteFiles:output_type -> proto.WEmpty
	0,  // 11: proto.Webss.SetBucketLife:output_type -> proto.WEmpty
	0,  // 12: proto.Webss.DelBucketLife:output_type -> proto.WEmpty
	0,  // 13: proto.Webss.SetFileLife:output_type -> proto.WEmpty
	9,  // 14: proto.Webss.SignFileUrl:output_type -> proto.SignUrl
	10, // 15: proto.Webss.SignFileUrls:output_type -> proto.SignUrls
	10, // 16: proto.Webss.OriginalUrls:output_type -> proto.SignUrls
	13, // 17: proto.Webss.GetFileInfo:output_type -> proto.Info
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_proto_webss_proto_init() }
func file_proto_webss_proto_init() {
	if File_proto_webss_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_webss_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WEmpty); i {
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
		file_proto_webss_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ID); i {
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
		file_proto_webss_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
		file_proto_webss_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Files); i {
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
		file_proto_webss_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Life); i {
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
		file_proto_webss_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Lifes); i {
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
		file_proto_webss_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
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
		file_proto_webss_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Sign); i {
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
		file_proto_webss_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Signs); i {
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
		file_proto_webss_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignUrl); i {
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
		file_proto_webss_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignUrls); i {
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
		file_proto_webss_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FNames); i {
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
		file_proto_webss_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FName); i {
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
		file_proto_webss_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Info); i {
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
			RawDescriptor: file_proto_webss_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_webss_proto_goTypes,
		DependencyIndexes: file_proto_webss_proto_depIdxs,
		MessageInfos:      file_proto_webss_proto_msgTypes,
	}.Build()
	File_proto_webss_proto = out.File
	file_proto_webss_proto_rawDesc = nil
	file_proto_webss_proto_goTypes = nil
	file_proto_webss_proto_depIdxs = nil
}
