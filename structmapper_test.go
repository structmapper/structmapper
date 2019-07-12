package structmapper

import (
	"testing"
	"time"

	"git.dmm.com/cto-tech/graphql-opencrud/lib/gql"
	"git.dmm.com/cto-tech/graphql-opencrud/lib/structmapper/test/graphql"
	"git.dmm.com/cto-tech/graphql-opencrud/lib/structmapper/test/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
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
			From: &graphql.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           graphql.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int(123),
				OptionalNum64: Int64(123),
				CreatedAt:     mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt:    mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
			},
			EmptyTo: new(graphql.User),
			ExpectedTo: &graphql.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           graphql.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int(123),
				OptionalNum64: Int64(123),
				CreatedAt:     mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt:    mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
			},
		},
		{
			Name: "graphql struct to protobuf struct",
			From: &graphql.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           graphql.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int(123),
				OptionalNum64: Int64(123),
				CreatedAt:     mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt:    mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
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
				CreatedAt:     mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
				ModifiedAt:    mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
			},
		},
		{
			Name: "protobuf struct to graphql struct",
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
				CreatedAt:     mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
				ModifiedAt:    mustTimestampProto(ptypes.TimestampProto(mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")))),
			},
			EmptyTo: new(graphql.User),
			ExpectedTo: &graphql.User{
				ID:            "12345",
				Name:          "Satoshi Nakamoto",
				Age:           47,
				Weight:        12.3,
				Sex:           graphql.SexFemale,
				Alive:         true,
				BirthDate:     String("1999-11-17"),
				Num64:         123,
				OptionalNum:   Int(123),
				OptionalNum64: Int64(123),
				CreatedAt:     mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
				ModifiedAt:    mustTime(time.Parse(time.RFC3339, "2019-07-07T12:34:56Z")),
			},
		},
	}

	mapper := New().Install(ProtobufModule)
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

func Int(i int) *int {
	return &i
}

func Int64(i gql.Int64) *gql.Int64 {
	return &i
}

func String(s string) *string {
	return &s
}
