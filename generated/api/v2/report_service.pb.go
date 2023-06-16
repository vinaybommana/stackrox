// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/v2/report_service.proto

package v2

import (
	context "context"
	fmt "fmt"
	types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ReportStatus_RunState int32

const (
	ReportStatus_WAITING   ReportStatus_RunState = 0
	ReportStatus_PREPARING ReportStatus_RunState = 1
	ReportStatus_SUCCESS   ReportStatus_RunState = 2
	ReportStatus_FAILURE   ReportStatus_RunState = 3
)

var ReportStatus_RunState_name = map[int32]string{
	0: "WAITING",
	1: "PREPARING",
	2: "SUCCESS",
	3: "FAILURE",
}

var ReportStatus_RunState_value = map[string]int32{
	"WAITING":   0,
	"PREPARING": 1,
	"SUCCESS":   2,
	"FAILURE":   3,
}

func (x ReportStatus_RunState) String() string {
	return proto.EnumName(ReportStatus_RunState_name, int32(x))
}

func (ReportStatus_RunState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c1e2917f181293be, []int{0, 0}
}

type ReportStatus_ReportMethod int32

const (
	ReportStatus_ON_DEMAND ReportStatus_ReportMethod = 0
	ReportStatus_SCHEDULED ReportStatus_ReportMethod = 1
)

var ReportStatus_ReportMethod_name = map[int32]string{
	0: "ON_DEMAND",
	1: "SCHEDULED",
}

var ReportStatus_ReportMethod_value = map[string]int32{
	"ON_DEMAND": 0,
	"SCHEDULED": 1,
}

func (x ReportStatus_ReportMethod) String() string {
	return proto.EnumName(ReportStatus_ReportMethod_name, int32(x))
}

func (ReportStatus_ReportMethod) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c1e2917f181293be, []int{0, 1}
}

type ReportStatus_NotificationMethod int32

const (
	ReportStatus_UNSET    ReportStatus_NotificationMethod = 0
	ReportStatus_EMAIL    ReportStatus_NotificationMethod = 1
	ReportStatus_DOWNLOAD ReportStatus_NotificationMethod = 2
)

var ReportStatus_NotificationMethod_name = map[int32]string{
	0: "UNSET",
	1: "EMAIL",
	2: "DOWNLOAD",
}

var ReportStatus_NotificationMethod_value = map[string]int32{
	"UNSET":    0,
	"EMAIL":    1,
	"DOWNLOAD": 2,
}

func (x ReportStatus_NotificationMethod) String() string {
	return proto.EnumName(ReportStatus_NotificationMethod_name, int32(x))
}

func (ReportStatus_NotificationMethod) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c1e2917f181293be, []int{0, 2}
}

type ReportStatus struct {
	RunState                 ReportStatus_RunState           `protobuf:"varint,1,opt,name=run_state,json=runState,proto3,enum=v2.ReportStatus_RunState" json:"run_state,omitempty"`
	RunTime                  *types.Timestamp                `protobuf:"bytes,2,opt,name=run_time,json=runTime,proto3" json:"run_time,omitempty"`
	ErrorMsg                 string                          `protobuf:"bytes,3,opt,name=error_msg,json=errorMsg,proto3" json:"error_msg,omitempty"`
	ReportMethod             ReportStatus_ReportMethod       `protobuf:"varint,4,opt,name=report_method,json=reportMethod,proto3,enum=v2.ReportStatus_ReportMethod" json:"report_method,omitempty"`
	ReportNotificationMethod ReportStatus_NotificationMethod `protobuf:"varint,5,opt,name=report_notification_method,json=reportNotificationMethod,proto3,enum=v2.ReportStatus_NotificationMethod" json:"report_notification_method,omitempty"`
	XXX_NoUnkeyedLiteral     struct{}                        `json:"-"`
	XXX_unrecognized         []byte                          `json:"-"`
	XXX_sizecache            int32                           `json:"-"`
}

func (m *ReportStatus) Reset()         { *m = ReportStatus{} }
func (m *ReportStatus) String() string { return proto.CompactTextString(m) }
func (*ReportStatus) ProtoMessage()    {}
func (*ReportStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1e2917f181293be, []int{0}
}
func (m *ReportStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ReportStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ReportStatus.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ReportStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportStatus.Merge(m, src)
}
func (m *ReportStatus) XXX_Size() int {
	return m.Size()
}
func (m *ReportStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ReportStatus proto.InternalMessageInfo

func (m *ReportStatus) GetRunState() ReportStatus_RunState {
	if m != nil {
		return m.RunState
	}
	return ReportStatus_WAITING
}

func (m *ReportStatus) GetRunTime() *types.Timestamp {
	if m != nil {
		return m.RunTime
	}
	return nil
}

func (m *ReportStatus) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}

func (m *ReportStatus) GetReportMethod() ReportStatus_ReportMethod {
	if m != nil {
		return m.ReportMethod
	}
	return ReportStatus_ON_DEMAND
}

func (m *ReportStatus) GetReportNotificationMethod() ReportStatus_NotificationMethod {
	if m != nil {
		return m.ReportNotificationMethod
	}
	return ReportStatus_UNSET
}

func (m *ReportStatus) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ReportStatus) Clone() *ReportStatus {
	if m == nil {
		return nil
	}
	cloned := new(ReportStatus)
	*cloned = *m

	cloned.RunTime = m.RunTime.Clone()
	return cloned
}

func init() {
	proto.RegisterEnum("v2.ReportStatus_RunState", ReportStatus_RunState_name, ReportStatus_RunState_value)
	proto.RegisterEnum("v2.ReportStatus_ReportMethod", ReportStatus_ReportMethod_name, ReportStatus_ReportMethod_value)
	proto.RegisterEnum("v2.ReportStatus_NotificationMethod", ReportStatus_NotificationMethod_name, ReportStatus_NotificationMethod_value)
	proto.RegisterType((*ReportStatus)(nil), "v2.ReportStatus")
}

func init() { proto.RegisterFile("api/v2/report_service.proto", fileDescriptor_c1e2917f181293be) }

var fileDescriptor_c1e2917f181293be = []byte{
	// 514 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xad, 0x1d, 0x4a, 0xe3, 0x6d, 0x02, 0xd6, 0x72, 0xc0, 0x49, 0x21, 0x8d, 0xc2, 0x25, 0x07,
	0xb4, 0x96, 0x8c, 0x40, 0x1c, 0x71, 0x62, 0x53, 0x22, 0x25, 0x4e, 0x65, 0x27, 0x14, 0x71, 0xb1,
	0x5c, 0x77, 0x9b, 0x5a, 0x60, 0x6f, 0xd8, 0x5d, 0x5b, 0x54, 0x88, 0x0b, 0xbf, 0xc0, 0x85, 0x4f,
	0xe2, 0x88, 0xc4, 0x0f, 0xa0, 0xc0, 0x8d, 0x9f, 0x40, 0x5e, 0xaf, 0x45, 0x20, 0xb7, 0x99, 0xf7,
	0x9e, 0xde, 0xbe, 0x99, 0x1d, 0x70, 0x14, 0xad, 0x13, 0xb3, 0xb0, 0x4c, 0x8a, 0xd7, 0x84, 0xf2,
	0x90, 0x61, 0x5a, 0x24, 0x31, 0x46, 0x6b, 0x4a, 0x38, 0x81, 0x6a, 0x61, 0x75, 0x8f, 0x57, 0x84,
	0xac, 0xde, 0x62, 0x53, 0x20, 0xe7, 0xf9, 0xa5, 0xc9, 0x93, 0x14, 0x33, 0x1e, 0xa5, 0xeb, 0x4a,
	0xd4, 0xbd, 0x23, 0x1d, 0x62, 0x92, 0xa6, 0x24, 0x93, 0x60, 0x47, 0x82, 0x0c, 0x47, 0x34, 0xbe,
	0x0a, 0xdf, 0xe5, 0x98, 0x5e, 0x4b, 0xea, 0x9e, 0x34, 0x2c, 0x15, 0x51, 0x96, 0x11, 0x1e, 0xf1,
	0x84, 0x64, 0xac, 0x62, 0x07, 0xbf, 0x1b, 0xa0, 0xe5, 0x8b, 0x2c, 0x01, 0x8f, 0x78, 0xce, 0xe0,
	0x13, 0xa0, 0xd1, 0x3c, 0x0b, 0x19, 0x8f, 0x38, 0x36, 0x94, 0xbe, 0x32, 0xbc, 0x65, 0x75, 0x50,
	0x61, 0xa1, 0x6d, 0x11, 0xf2, 0xf3, 0xac, 0xac, 0xb0, 0xdf, 0xa4, 0xb2, 0x82, 0x8f, 0x41, 0x59,
	0x87, 0x65, 0x5a, 0x43, 0xed, 0x2b, 0xc3, 0x43, 0xab, 0x8b, 0xaa, 0x97, 0x51, 0x3d, 0x0a, 0x5a,
	0xd4, 0xa3, 0xf8, 0x07, 0x34, 0xcf, 0xca, 0x0e, 0x1e, 0x01, 0x0d, 0x53, 0x4a, 0x68, 0x98, 0xb2,
	0x95, 0xd1, 0xe8, 0x2b, 0x43, 0xcd, 0x6f, 0x0a, 0x60, 0xc6, 0x56, 0x70, 0x04, 0xda, 0x72, 0x4f,
	0x29, 0xe6, 0x57, 0xe4, 0xc2, 0xb8, 0x21, 0xf2, 0xdc, 0xdf, 0xcd, 0x23, 0x9a, 0x99, 0x10, 0xf9,
	0x2d, 0xba, 0xd5, 0xc1, 0x08, 0x74, 0xa5, 0x47, 0x46, 0x78, 0x72, 0x99, 0xc4, 0x62, 0xfc, 0xda,
	0x70, 0x5f, 0x18, 0x3e, 0xd8, 0x31, 0xf4, 0xb6, 0xb4, 0xd2, 0xd6, 0xa8, 0x6c, 0x76, 0x99, 0xc1,
	0x33, 0xd0, 0xac, 0x17, 0x02, 0x0f, 0xc1, 0xc1, 0x99, 0x3d, 0x59, 0x4c, 0xbc, 0x13, 0x7d, 0x0f,
	0xb6, 0x81, 0x76, 0xea, 0xbb, 0xa7, 0xb6, 0x5f, 0xb6, 0x4a, 0xc9, 0x05, 0xcb, 0xf1, 0xd8, 0x0d,
	0x02, 0x5d, 0x2d, 0x9b, 0xe7, 0xf6, 0x64, 0xba, 0xf4, 0x5d, 0xbd, 0x31, 0x78, 0x58, 0x7f, 0x82,
	0x0c, 0xdd, 0x06, 0xda, 0xdc, 0x0b, 0x1d, 0x77, 0x66, 0x7b, 0x4e, 0xe5, 0x13, 0x8c, 0x5f, 0xb8,
	0xce, 0x72, 0xea, 0x3a, 0xba, 0x32, 0x78, 0x0a, 0xe0, 0x6e, 0x0a, 0xa8, 0x81, 0xfd, 0xa5, 0x17,
	0xb8, 0x0b, 0x7d, 0xaf, 0x2c, 0xdd, 0x99, 0x3d, 0x99, 0xea, 0x0a, 0x6c, 0x81, 0xa6, 0x33, 0x3f,
	0xf3, 0xa6, 0x73, 0xdb, 0xd1, 0x55, 0x6b, 0x05, 0xda, 0x72, 0xcc, 0xea, 0xee, 0xe0, 0x4b, 0x70,
	0xfb, 0x04, 0xf3, 0x7f, 0x0e, 0x40, 0xaf, 0x96, 0xc1, 0x48, 0x4e, 0x63, 0x3c, 0xba, 0x9e, 0x38,
	0x5d, 0xfd, 0xff, 0xf5, 0x0c, 0x8e, 0x3f, 0x7d, 0xff, 0xf5, 0x59, 0xed, 0xc0, 0xbb, 0x7f, 0x4f,
	0x99, 0x99, 0x4c, 0x70, 0xe6, 0x87, 0xe4, 0xe2, 0xe3, 0x08, 0x7d, 0xdd, 0xf4, 0x94, 0x6f, 0x9b,
	0x9e, 0xf2, 0x63, 0xd3, 0x53, 0xbe, 0xfc, 0xec, 0xed, 0x01, 0x23, 0x21, 0x88, 0xf1, 0x28, 0x7e,
	0x43, 0xc9, 0xfb, 0xea, 0x22, 0x50, 0xb4, 0x4e, 0x50, 0x61, 0xbd, 0x56, 0x0b, 0xeb, 0x55, 0xe3,
	0xfc, 0xa6, 0xc0, 0x1e, 0xfd, 0x09, 0x00, 0x00, 0xff, 0xff, 0x60, 0x38, 0xd6, 0xe2, 0x21, 0x03,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ReportServiceClient is the client API for ReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConnInterface.NewStream.
type ReportServiceClient interface {
	// ListReportStatus returns report status
	GetReportStatus(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*ReportStatus, error)
}

type reportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReportServiceClient(cc grpc.ClientConnInterface) ReportServiceClient {
	return &reportServiceClient{cc}
}

func (c *reportServiceClient) GetReportStatus(ctx context.Context, in *ResourceByID, opts ...grpc.CallOption) (*ReportStatus, error) {
	out := new(ReportStatus)
	err := c.cc.Invoke(ctx, "/v2.ReportService/GetReportStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReportServiceServer is the server API for ReportService service.
type ReportServiceServer interface {
	// ListReportStatus returns report status
	GetReportStatus(context.Context, *ResourceByID) (*ReportStatus, error)
}

// UnimplementedReportServiceServer can be embedded to have forward compatible implementations.
type UnimplementedReportServiceServer struct {
}

func (*UnimplementedReportServiceServer) GetReportStatus(ctx context.Context, req *ResourceByID) (*ReportStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReportStatus not implemented")
}

func RegisterReportServiceServer(s *grpc.Server, srv ReportServiceServer) {
	s.RegisterService(&_ReportService_serviceDesc, srv)
}

func _ReportService_GetReportStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResourceByID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetReportStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v2.ReportService/GetReportStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetReportStatus(ctx, req.(*ResourceByID))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReportService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v2.ReportService",
	HandlerType: (*ReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReportStatus",
			Handler:    _ReportService_GetReportStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v2/report_service.proto",
}

func (m *ReportStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ReportStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ReportStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ReportNotificationMethod != 0 {
		i = encodeVarintReportService(dAtA, i, uint64(m.ReportNotificationMethod))
		i--
		dAtA[i] = 0x28
	}
	if m.ReportMethod != 0 {
		i = encodeVarintReportService(dAtA, i, uint64(m.ReportMethod))
		i--
		dAtA[i] = 0x20
	}
	if len(m.ErrorMsg) > 0 {
		i -= len(m.ErrorMsg)
		copy(dAtA[i:], m.ErrorMsg)
		i = encodeVarintReportService(dAtA, i, uint64(len(m.ErrorMsg)))
		i--
		dAtA[i] = 0x1a
	}
	if m.RunTime != nil {
		{
			size, err := m.RunTime.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintReportService(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RunState != 0 {
		i = encodeVarintReportService(dAtA, i, uint64(m.RunState))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintReportService(dAtA []byte, offset int, v uint64) int {
	offset -= sovReportService(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ReportStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RunState != 0 {
		n += 1 + sovReportService(uint64(m.RunState))
	}
	if m.RunTime != nil {
		l = m.RunTime.Size()
		n += 1 + l + sovReportService(uint64(l))
	}
	l = len(m.ErrorMsg)
	if l > 0 {
		n += 1 + l + sovReportService(uint64(l))
	}
	if m.ReportMethod != 0 {
		n += 1 + sovReportService(uint64(m.ReportMethod))
	}
	if m.ReportNotificationMethod != 0 {
		n += 1 + sovReportService(uint64(m.ReportNotificationMethod))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovReportService(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozReportService(x uint64) (n int) {
	return sovReportService(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ReportStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowReportService
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ReportStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ReportStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RunState", wireType)
			}
			m.RunState = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReportService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RunState |= ReportStatus_RunState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RunTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReportService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthReportService
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthReportService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RunTime == nil {
				m.RunTime = &types.Timestamp{}
			}
			if err := m.RunTime.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ErrorMsg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReportService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthReportService
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthReportService
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ErrorMsg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReportMethod", wireType)
			}
			m.ReportMethod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReportService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReportMethod |= ReportStatus_ReportMethod(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReportNotificationMethod", wireType)
			}
			m.ReportNotificationMethod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowReportService
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReportNotificationMethod |= ReportStatus_NotificationMethod(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipReportService(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthReportService
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipReportService(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowReportService
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowReportService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowReportService
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthReportService
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupReportService
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthReportService
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthReportService        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowReportService          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupReportService = fmt.Errorf("proto: unexpected end of group")
)
