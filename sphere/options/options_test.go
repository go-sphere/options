package options_test

import (
	"testing"

	"github.com/go-sphere/options/sphere/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestKeyValuePairDescriptor(t *testing.T) {
	descriptor := (&options.KeyValuePair{}).ProtoReflect().Descriptor()
	if got, want := descriptor.FullName(), protoreflect.FullName("sphere.options.KeyValuePair"); got != want {
		t.Fatalf("full name = %q, want %q", got, want)
	}

	tests := []struct {
		name   protoreflect.Name
		number protoreflect.FieldNumber
		kind   protoreflect.Kind
		oneof  bool
		isMap  bool
	}{
		{name: "key", number: 1, kind: protoreflect.StringKind},
		{name: "flag", number: 2, kind: protoreflect.BoolKind, oneof: true},
		{name: "text", number: 3, kind: protoreflect.StringKind, oneof: true},
		{name: "number", number: 4, kind: protoreflect.Int64Kind, oneof: true},
		{name: "extra", number: 5, kind: protoreflect.MessageKind, isMap: true},
	}
	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			field := descriptor.Fields().ByName(tt.name)
			if field == nil {
				t.Fatalf("field %q not found", tt.name)
			}
			if got := field.Number(); got != tt.number {
				t.Errorf("number = %d, want %d", got, tt.number)
			}
			if got := field.Kind(); got != tt.kind {
				t.Errorf("kind = %s, want %s", got, tt.kind)
			}
			if got := field.ContainingOneof() != nil; got != tt.oneof {
				t.Errorf("belongs to oneof = %v, want %v", got, tt.oneof)
			}
			if got := field.IsMap(); got != tt.isMap {
				t.Errorf("is map = %v, want %v", got, tt.isMap)
			}
		})
	}
	if got := descriptor.Oneofs().ByName("value"); got == nil || got.IsSynthetic() {
		t.Errorf("value oneof = %v, want non-synthetic oneof", got)
	}
}

func TestOptionsExtensionDescriptor(t *testing.T) {
	descriptor := options.E_Options.TypeDescriptor()
	if got, want := descriptor.FullName(), protoreflect.FullName("sphere.options.options"); got != want {
		t.Errorf("full name = %q, want %q", got, want)
	}
	if got, want := descriptor.Number(), protoreflect.FieldNumber(501319300); got != want {
		t.Errorf("number = %d, want %d", got, want)
	}
	if got, want := descriptor.Kind(), protoreflect.MessageKind; got != want {
		t.Errorf("kind = %s, want %s", got, want)
	}
	if got, want := descriptor.Cardinality(), protoreflect.Repeated; got != want {
		t.Errorf("cardinality = %s, want %s", got, want)
	}
	if got, want := descriptor.ContainingMessage().FullName(), protoreflect.FullName("google.protobuf.MethodOptions"); got != want {
		t.Errorf("extendee = %q, want %q", got, want)
	}
}

func TestOptionsExtensionWireRoundTrip(t *testing.T) {
	input := &descriptorpb.MethodOptions{}
	proto.SetExtension(input, options.E_Options, []*options.KeyValuePair{
		{Key: "public", Value: &options.KeyValuePair_Flag{Flag: false}},
		{Key: "route", Value: &options.KeyValuePair_Text{Text: "/v1/users"}},
		{Key: "timeout", Value: &options.KeyValuePair_Number{Number: 30}, Extra: map[string]string{"unit": "seconds"}},
	})

	data, err := proto.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	output := &descriptorpb.MethodOptions{}
	if err := proto.Unmarshal(data, output); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !proto.Equal(output, input) {
		t.Errorf("round trip mismatch:\n got: %v\nwant: %v", output, input)
	}

	values := proto.GetExtension(output, options.E_Options).([]*options.KeyValuePair)
	if _, ok := values[0].GetValue().(*options.KeyValuePair_Flag); !ok {
		t.Errorf("explicit false flag lost its oneof discriminator: %T", values[0].GetValue())
	}
}
