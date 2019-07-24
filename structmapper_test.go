package structmapper

import (
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

func Int64(i int64) *int64 {
	return &i
}

func CustonInt64(i dto.CustomInt64) *dto.CustomInt64 {
	return &i
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
