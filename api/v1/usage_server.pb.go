// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: api/v1/usage_server.proto

package v1

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

type GetAggregatedSummaryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TenantId string `protobuf:"bytes,1,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	// star_time is the UNIX timestamp for the summary start time (inclusive).
	// If start_time is not provided, the default is the 24 hours before end_time.
	StartTime int64 `protobuf:"varint,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// end_time is the UNIX timestamp for the summary end time (exclusive).
	// If end_time is not provided, the default is the current time.
	EndTime int64 `protobuf:"varint,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *GetAggregatedSummaryRequest) Reset() {
	*x = GetAggregatedSummaryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_usage_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetAggregatedSummaryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAggregatedSummaryRequest) ProtoMessage() {}

func (x *GetAggregatedSummaryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_usage_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAggregatedSummaryRequest.ProtoReflect.Descriptor instead.
func (*GetAggregatedSummaryRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_usage_server_proto_rawDescGZIP(), []int{0}
}

func (x *GetAggregatedSummaryRequest) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

func (x *GetAggregatedSummaryRequest) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *GetAggregatedSummaryRequest) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

type Summary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method          string  `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	TotalRequests   int64   `protobuf:"varint,2,opt,name=total_requests,json=totalRequests,proto3" json:"total_requests,omitempty"`
	SuccessRequests int64   `protobuf:"varint,3,opt,name=success_requests,json=successRequests,proto3" json:"success_requests,omitempty"`
	FailureRequests int64   `protobuf:"varint,4,opt,name=failure_requests,json=failureRequests,proto3" json:"failure_requests,omitempty"`
	AverageLatency  float64 `protobuf:"fixed64,5,opt,name=average_latency,json=averageLatency,proto3" json:"average_latency,omitempty"`
}

func (x *Summary) Reset() {
	*x = Summary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_usage_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Summary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Summary) ProtoMessage() {}

func (x *Summary) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_usage_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Summary.ProtoReflect.Descriptor instead.
func (*Summary) Descriptor() ([]byte, []int) {
	return file_api_v1_usage_server_proto_rawDescGZIP(), []int{1}
}

func (x *Summary) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Summary) GetTotalRequests() int64 {
	if x != nil {
		return x.TotalRequests
	}
	return 0
}

func (x *Summary) GetSuccessRequests() int64 {
	if x != nil {
		return x.SuccessRequests
	}
	return 0
}

func (x *Summary) GetFailureRequests() int64 {
	if x != nil {
		return x.FailureRequests
	}
	return 0
}

func (x *Summary) GetAverageLatency() float64 {
	if x != nil {
		return x.AverageLatency
	}
	return 0
}

type AggregatedSummary struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Summary         *Summary   `protobuf:"bytes,1,opt,name=summary,proto3" json:"summary,omitempty"`
	MethodSummaries []*Summary `protobuf:"bytes,2,rep,name=method_summaries,json=methodSummaries,proto3" json:"method_summaries,omitempty"`
}

func (x *AggregatedSummary) Reset() {
	*x = AggregatedSummary{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_usage_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregatedSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregatedSummary) ProtoMessage() {}

func (x *AggregatedSummary) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_usage_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregatedSummary.ProtoReflect.Descriptor instead.
func (*AggregatedSummary) Descriptor() ([]byte, []int) {
	return file_api_v1_usage_server_proto_rawDescGZIP(), []int{2}
}

func (x *AggregatedSummary) GetSummary() *Summary {
	if x != nil {
		return x.Summary
	}
	return nil
}

func (x *AggregatedSummary) GetMethodSummaries() []*Summary {
	if x != nil {
		return x.MethodSummaries
	}
	return nil
}

var File_api_v1_usage_server_proto protoreflect.FileDescriptor

var file_api_v1_usage_server_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x6c, 0x6c, 0x6d,
	0x61, 0x72, 0x69, 0x6e, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x22, 0x74, 0x0a, 0x1b, 0x47, 0x65, 0x74,
	0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e,
	0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22,
	0xc7, 0x01, 0x0a, 0x07, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74,
	0x68, 0x6f, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x73, 0x12, 0x29, 0x0a, 0x10, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73,
	0x12, 0x27, 0x0a, 0x0f, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x6c, 0x61, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0e, 0x61, 0x76, 0x65, 0x72, 0x61,
	0x67, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x22, 0xa6, 0x01, 0x0a, 0x11, 0x41, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12,
	0x3f, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x25, 0x2e, 0x6c, 0x6c, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69,
	0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x12, 0x50, 0x0a, 0x10, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x73, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6c, 0x6c, 0x6d,
	0x61, 0x72, 0x69, 0x6e, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72,
	0x79, 0x52, 0x0f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x69,
	0x65, 0x73, 0x32, 0x96, 0x01, 0x0a, 0x0f, 0x41, 0x50, 0x49, 0x55, 0x73, 0x61, 0x67, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x82, 0x01, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x41, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12,
	0x39, 0x2e, 0x6c, 0x6c, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x64, 0x53, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x6c, 0x6c, 0x6d,
	0x61, 0x72, 0x69, 0x6e, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x65, 0x64, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x42, 0x27, 0x5a, 0x25, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x6c, 0x6d, 0x61, 0x72, 0x69,
	0x6e, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2d, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_usage_server_proto_rawDescOnce sync.Once
	file_api_v1_usage_server_proto_rawDescData = file_api_v1_usage_server_proto_rawDesc
)

func file_api_v1_usage_server_proto_rawDescGZIP() []byte {
	file_api_v1_usage_server_proto_rawDescOnce.Do(func() {
		file_api_v1_usage_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_usage_server_proto_rawDescData)
	})
	return file_api_v1_usage_server_proto_rawDescData
}

var file_api_v1_usage_server_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_v1_usage_server_proto_goTypes = []interface{}{
	(*GetAggregatedSummaryRequest)(nil), // 0: llmariner.apiusage.server.v1.GetAggregatedSummaryRequest
	(*Summary)(nil),                     // 1: llmariner.apiusage.server.v1.Summary
	(*AggregatedSummary)(nil),           // 2: llmariner.apiusage.server.v1.AggregatedSummary
}
var file_api_v1_usage_server_proto_depIdxs = []int32{
	1, // 0: llmariner.apiusage.server.v1.AggregatedSummary.summary:type_name -> llmariner.apiusage.server.v1.Summary
	1, // 1: llmariner.apiusage.server.v1.AggregatedSummary.method_summaries:type_name -> llmariner.apiusage.server.v1.Summary
	0, // 2: llmariner.apiusage.server.v1.APIUsageService.GetAggregatedSummary:input_type -> llmariner.apiusage.server.v1.GetAggregatedSummaryRequest
	2, // 3: llmariner.apiusage.server.v1.APIUsageService.GetAggregatedSummary:output_type -> llmariner.apiusage.server.v1.AggregatedSummary
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_v1_usage_server_proto_init() }
func file_api_v1_usage_server_proto_init() {
	if File_api_v1_usage_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_usage_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetAggregatedSummaryRequest); i {
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
		file_api_v1_usage_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Summary); i {
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
		file_api_v1_usage_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregatedSummary); i {
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
			RawDescriptor: file_api_v1_usage_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_usage_server_proto_goTypes,
		DependencyIndexes: file_api_v1_usage_server_proto_depIdxs,
		MessageInfos:      file_api_v1_usage_server_proto_msgTypes,
	}.Build()
	File_api_v1_usage_server_proto = out.File
	file_api_v1_usage_server_proto_rawDesc = nil
	file_api_v1_usage_server_proto_goTypes = nil
	file_api_v1_usage_server_proto_depIdxs = nil
}
