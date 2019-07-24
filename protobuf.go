package structmapper

import (
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
)

// ProtobufModule is Transformer Module of between golang/protobuf/ptypes types and go types
func ProtobufModule(m Mapper) {
	registerTimestamp(m)
	registerWrappers(m)
}

// for ptypes/timestamp.Timestamp
func registerTimestamp(m Mapper) {
	// string <-> Timestamp
	m.RegisterTransformer(
		Target{
			From: reflect.TypeOf(""),
			To:   reflect.TypeOf(timestamp.Timestamp{}),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			t, err := time.Parse(time.RFC3339, from.String())
			if err != nil {
				return reflect.ValueOf(nil), errors.WithStack(err)
			}

			ts, err := ptypes.TimestampProto(t)
			if err != nil {
				return reflect.ValueOf(nil), errors.WithStack(err)
			}

			return reflect.ValueOf(*ts), nil
		},
	)
	m.RegisterTransformer(
		Target{
			From: reflect.TypeOf(timestamp.Timestamp{}),
			To:   reflect.TypeOf(""),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			ts, ok := from.Interface().(timestamp.Timestamp)
			if !ok {
				return reflect.ValueOf(nil), errors.Errorf("Invalid value was found, expected timestamp.Timestamp, but was %+v", from)
			}

			t, err := ptypes.Timestamp(&ts)
			if err != nil {
				return reflect.ValueOf(""), errors.WithStack(err)
			}

			return reflect.ValueOf(t.Format(time.RFC3339)), nil
		},
	)

	// time.Time <-> Timestamp
	m.RegisterTransformer(
		Target{
			From: reflect.TypeOf(time.Time{}),
			To:   reflect.TypeOf(timestamp.Timestamp{}),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			t, ok := from.Interface().(time.Time)
			if !ok {
				return reflect.ValueOf(nil), errors.Errorf("Invalid value was found, expected time.Time, but was %+v", from)
			}

			ts, err := ptypes.TimestampProto(t)
			if err != nil {
				return reflect.ValueOf(nil), errors.WithStack(err)
			}

			return reflect.ValueOf(*ts), nil
		},
	)
	m.RegisterTransformer(
		Target{
			From: reflect.TypeOf(timestamp.Timestamp{}),
			To:   reflect.TypeOf(time.Time{}),
		},
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			ts, ok := from.Interface().(timestamp.Timestamp)
			if !ok {
				return reflect.ValueOf(nil), errors.Errorf("Invalid value was found, expected timestamp.Timestamp, but was %+v", from)
			}

			t, err := ptypes.Timestamp(&ts)
			if err != nil {
				return reflect.ValueOf(""), errors.WithStack(err)
			}

			return reflect.ValueOf(t), nil
		},
	)

}

// for ptypes/wrappers.*
func registerWrappers(m Mapper) {
	for _, tm := range protoTypeMappings {
		registerWrapper(m, tm)
	}
}

func registerWrapper(m Mapper, tm protoTypeMapping) {
	m.RegisterTransformerFunc(
		// matcher
		func(target Target) bool {
			return target.To == tm.WrapperType && tm.ContainsInAcceptableTypes(target.From)
		},
		// mapper
		func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			return tm.AsProto(from)
		},
	)

	m.RegisterTransformerFunc(
		// matcher
		func(target Target) bool {
			return target.From == tm.WrapperType && tm.ContainsInAcceptableTypes(target.To)
		},
		// mapper
		func(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
			return tm.AsValue(from, toType)
		},
	)
}

var _ Module = ProtobufModule

type protoTypeMapping struct {
	AcceptableTypes []reflect.Type
	WrapperType     reflect.Type
	AsProto         func(reflect.Value) (reflect.Value, error)
	AsValue         func(reflect.Value, reflect.Type) (reflect.Value, error)
}

func (m *protoTypeMapping) ContainsInAcceptableTypes(t reflect.Type) bool {
	for _, at := range m.AcceptableTypes {
		if t.ConvertibleTo(at) {
			return true
		}
	}
	return false
}

var protoTypeMappings = []protoTypeMapping{
	{
		AcceptableTypes: []reflect.Type{reflect.TypeOf(int(0)), reflect.TypeOf(int32(0)), reflect.TypeOf(int64(0))},
		WrapperType:     reflect.TypeOf(wrappers.Int64Value{}),
		AsProto: func(from reflect.Value) (reflect.Value, error) {
			return reflect.ValueOf(wrappers.Int64Value{Value: from.Int()}), nil
		},
		AsValue: func(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
			v, ok := from.Interface().(wrappers.Int64Value)
			if !ok {
				return reflect.Zero(toType), errors.Errorf("Invalid value type: %+v", from)
			}
			return reflect.ValueOf(v.Value).Convert(toType), nil
		},
	},
	{
		AcceptableTypes: []reflect.Type{reflect.TypeOf(int(0)), reflect.TypeOf(int32(0)), reflect.TypeOf(int64(0))},
		WrapperType:     reflect.TypeOf(wrappers.Int32Value{}),
		AsProto: func(from reflect.Value) (reflect.Value, error) {
			return reflect.ValueOf(wrappers.Int32Value{Value: int32(from.Int())}), nil
		},
		AsValue: func(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
			v, ok := from.Interface().(wrappers.Int32Value)
			if !ok {
				return reflect.Zero(toType), errors.Errorf("Invalid value type: %+v", from)
			}
			return reflect.ValueOf(v.Value).Convert(toType), nil
		},
	},
	{
		AcceptableTypes: []reflect.Type{reflect.TypeOf(float32(0)), reflect.TypeOf(float64(0))},
		WrapperType:     reflect.TypeOf(wrappers.DoubleValue{}),
		AsProto: func(from reflect.Value) (reflect.Value, error) {
			return reflect.ValueOf(wrappers.DoubleValue{Value: from.Float()}), nil
		},
		AsValue: func(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
			v, ok := from.Interface().(wrappers.DoubleValue)
			if !ok {
				return reflect.Zero(toType), errors.Errorf("Invalid value type: %+v", from)
			}
			return reflect.ValueOf(v.Value).Convert(toType), nil
		},
	},
	{
		AcceptableTypes: []reflect.Type{reflect.TypeOf(float32(0)), reflect.TypeOf(float64(0))},
		WrapperType:     reflect.TypeOf(wrappers.FloatValue{}),
		AsProto: func(from reflect.Value) (reflect.Value, error) {
			return reflect.ValueOf(wrappers.FloatValue{Value: float32(from.Float())}), nil
		},
		AsValue: func(from reflect.Value, toType reflect.Type) (reflect.Value, error) {
			v, ok := from.Interface().(wrappers.FloatValue)
			if !ok {
				return reflect.Zero(toType), errors.Errorf("Invalid value type: %+v", from)
			}
			return reflect.ValueOf(v.Value).Convert(toType), nil
		},
	},
	{
		AcceptableTypes: []reflect.Type{reflect.TypeOf(true)},
		WrapperType:     reflect.TypeOf(wrappers.BoolValue{}),
		AsProto: func(from reflect.Value) (reflect.Value, error) {
			return reflect.ValueOf(wrappers.BoolValue{Value: from.Bool()}), nil
		},
		AsValue: func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			v, ok := from.Interface().(wrappers.BoolValue)
			if !ok {
				return reflect.ValueOf(false), errors.Errorf("Invalid value type: %+v", from)
			}
			return reflect.ValueOf(v.Value), nil
		},
	},
	{
		AcceptableTypes: []reflect.Type{reflect.TypeOf("")},
		WrapperType:     reflect.TypeOf(wrappers.StringValue{}),
		AsProto: func(from reflect.Value) (reflect.Value, error) {
			return reflect.ValueOf(wrappers.StringValue{Value: from.String()}), nil
		},
		AsValue: func(from reflect.Value, _ reflect.Type) (reflect.Value, error) {
			v, ok := from.Interface().(wrappers.StringValue)
			if !ok {
				return reflect.ValueOf(""), errors.Errorf("Invalid value type: %+v", from)
			}
			return reflect.ValueOf(v.Value), nil
		},
	},
}
