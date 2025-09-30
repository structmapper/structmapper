package structmapper

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/structmapper/structmapper/test/dto"
	"github.com/structmapper/structmapper/test/proto"
)

func TestCopy(t *testing.T) {
	cases := []struct {
		Name       string
		From       interface{}
		EmptyTo    interface{}
		ExpectedTo interface{}
	}{
		{
			Name:       "struct{} to struct{}",
			From:       &struct{}{},
			EmptyTo:    &struct{}{},
			ExpectedTo: &struct{}{},
		},
		{
			Name: "same struct copy",
			From: &dto.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           dto.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int32(123),
				OptionalNum64: CustonInt64(123),
				Numbers:       []int64{1, 2, 3},
				Times: []time.Time{
					mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-08T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-09T12:34:56Z")),
				},
				CreatedAt:  mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt: mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
			},
			EmptyTo: new(dto.User),
			ExpectedTo: &dto.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           dto.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int32(123),
				OptionalNum64: CustonInt64(123),
				Numbers:       []int64{1, 2, 3},
				Times: []time.Time{
					mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-08T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-09T12:34:56Z")),
				},
				CreatedAt:  mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt: mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
			},
		},
		{
			Name: "dto struct to protobuf struct",
			From: &dto.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           dto.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int32(123),
				OptionalNum64: CustonInt64(123),
				Numbers:       []int64{1, 2, 3},
				Times: []time.Time{
					mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-08T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-09T12:34:56Z")),
				},
				CreatedAt:  mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt: mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
			},
			EmptyTo: new(proto.User),
			ExpectedTo: &proto.User{
				Id:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           "Female",
				Alive:         true,
				BirthDate:     "1999-11-17",
				Num64:         123,
				OptionalNum:   &wrappers.Int64Value{Value: 123},
				OptionalNum64: &wrappers.Int64Value{Value: 123},
				Numbers:       []int64{1, 2, 3},
				Times: []*timestamp.Timestamp{
					mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
					mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-08T12:34:56Z")))),
					mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-09T12:34:56Z")))),
				},
				CreatedAt:  mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
				ModifiedAt: mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
			},
		},
		{
			Name: "protobuf struct to dto struct",
			From: &proto.User{
				Id:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           "Female",
				Alive:         true,
				BirthDate:     "1999-11-17",
				Num64:         123,
				OptionalNum:   &wrappers.Int64Value{Value: 123},
				OptionalNum64: &wrappers.Int64Value{Value: 123},
				Numbers:       []int64{1, 2, 3},
				Times: []*timestamp.Timestamp{
					mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
					mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-08T12:34:56Z")))),
					mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-09T12:34:56Z")))),
				},
				CreatedAt:  mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
				ModifiedAt: mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
			},
			EmptyTo: new(dto.User),
			ExpectedTo: &dto.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           dto.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int32(123),
				OptionalNum64: CustonInt64(123),
				Numbers:       []int64{1, 2, 3},
				Times: []time.Time{
					mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-08T12:34:56Z")),
					mustTime(time.Parse(time.RFC3339, "2019-07-09T12:34:56Z")),
				},
				CreatedAt:  mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt: mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
			},
		},
		{
			Name: "enum pointer to string",
			From: &EnumPtrStruct{
				SexPtr: func() *dto.Sex { sex := dto.SexFemale; return &sex }(),
			},
			EmptyTo: new(EnumStringStruct),
			ExpectedTo: &EnumStringStruct{
				Sex: "Female",
			},
		},
		{
			Name: "nil enum pointer to string",
			From: &EnumPtrStruct{
				SexPtr: func() *dto.Sex { var sex *dto.Sex; return sex }(),
			},
			EmptyTo: new(EnumStringStruct),
			ExpectedTo: &EnumStringStruct{
				Sex: "",
			},
		},
		{
			Name: "empty slice handling",
			From: &struct {
				Numbers []int64 `structmapper:"numbers"`
			}{
				Numbers: []int64{},
			},
			EmptyTo: new(struct {
				Numbers []int64 `structmapper:"numbers"`
			}),
			ExpectedTo: &struct {
				Numbers []int64 `structmapper:"numbers"`
			}{
				Numbers: []int64{},
			},
		},
		{
			Name: "nil slice handling",
			From: &struct {
				Numbers []int64 `structmapper:"numbers"`
			}{
				Numbers: nil,
			},
			EmptyTo: new(struct {
				Numbers []int64 `structmapper:"numbers"`
			}),
			ExpectedTo: &struct {
				Numbers []int64 `structmapper:"numbers"`
			}{
				Numbers: nil,
			},
		},
		{
			Name: "ignore tags skip sensitive data",
			From: &struct {
				ID       string
				Secret   string `structmapper:"-"`
				Internal string `json:"-"`
			}{
				ID:       "user-123",
				Secret:   "top-secret",
				Internal: "internal-value",
			},
			EmptyTo: new(struct {
				ID       string
				Secret   string
				Internal string
			}),
			ExpectedTo: &struct {
				ID       string
				Secret   string
				Internal string
			}{
				ID:       "user-123",
				Secret:   "",
				Internal: "",
			},
		},
		{
			Name: "protobuf wrapper types",
			From: &struct {
				OptionalBool   *wrappers.BoolValue   `structmapper:"optional_bool"`
				OptionalInt32  *wrappers.Int32Value  `structmapper:"optional_int32"`
				OptionalString *wrappers.StringValue `structmapper:"optional_string"`
			}{
				OptionalBool:   &wrappers.BoolValue{Value: true},
				OptionalInt32:  &wrappers.Int32Value{Value: 42},
				OptionalString: &wrappers.StringValue{Value: "test"},
			},
			EmptyTo: new(struct {
				OptionalBool   *bool   `structmapper:"optional_bool"`
				OptionalInt32  *int32  `structmapper:"optional_int32"`
				OptionalString *string `structmapper:"optional_string"`
			}),
			ExpectedTo: &struct {
				OptionalBool   *bool   `structmapper:"optional_bool"`
				OptionalInt32  *int32  `structmapper:"optional_int32"`
				OptionalString *string `structmapper:"optional_string"`
			}{
				OptionalBool:   func() *bool { b := true; return &b }(),
				OptionalInt32:  func() *int32 { i := int32(42); return &i }(),
				OptionalString: func() *string { s := "test"; return &s }(),
			},
		},
	}

	mapper := New().
		EnableLogging().
		Install(ProtobufModule).
		Install(StringerModule)
	for _, _c := range cases {
		c := _c
		t.Run(c.Name, func(t *testing.T) {
			to := c.EmptyTo
			if assert.NoError(t, mapper.From(c.From).CopyTo(to)) {
				assert.EqualValues(t, c.ExpectedTo, to)
			}
		})
	}
}

func mustTime(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

func mustTimestampProto(t *timestamp.Timestamp, err error) *timestamp.Timestamp {
	if err != nil {
		panic(err)
	}
	return t
}

func Int32(i int32) *int32 {
	return &i
}

func CustonInt64(i dto.CustomInt64) *dto.CustomInt64 {
	return &i
}

func TestNamesOfSkip(t *testing.T) {
	typ := reflect.TypeOf(struct {
		Secret   string `structmapper:"-"`
		Internal string `json:"-"`
	}{})

	secretField, ok := typ.FieldByName("Secret")
	assert.True(t, ok)
	assert.Empty(t, namesOf(secretField))

	internalField, ok := typ.FieldByName("Internal")
	assert.True(t, ok)
	assert.Empty(t, namesOf(internalField))
}

func TestDeepFieldsNames(t *testing.T) {
	typ := reflect.TypeOf(struct {
		Secret   string `structmapper:"-"`
		Internal string `json:"-"`
	}{})

	fields := deepFields(typ)
	assert.Len(t, fields, 2)
	assert.Empty(t, namesOf(fields[0]))
	assert.Empty(t, namesOf(fields[1]))
}

func TestCopyIgnoreTags(t *testing.T) {
	type sensitive struct {
		ID       string
		Secret   string `structmapper:"-"`
		Internal string `json:"-"`
	}

	src := &sensitive{ID: "user-123", Secret: "top-secret", Internal: "internal-value"}
	dst := &struct {
		ID       string
		Secret   string
		Internal string
	}{}

	require := assert.New(t)
	mapper := New()
	require.NoError(mapper.From(src).CopyTo(dst))
	require.Equal("user-123", dst.ID)
	require.Empty(dst.Secret)
	require.Empty(dst.Internal)
}

func String(s string) *string {
	return &s
}

type EnumPtrStruct struct {
	SexPtr *dto.Sex `structmapper:"sex"`
}

type EnumStringStruct struct {
	Sex string `structmapper:"sex"`
}
