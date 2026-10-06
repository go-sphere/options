// Package options holds the Go bindings for sphere/options/options.proto, a
// generic protobuf extension that attaches keyed metadata to RPC methods.
//
// The (sphere.options.options) method option is a repeated [KeyValuePair].
// Each entry has a key naming its consumer, an optional typed value (flag,
// text, or number), and an extra map of string parameters. Code generators
// pick the entries they understand by key: protoc-gen-route, for example,
// generates routing code for non-streaming methods that carry an entry whose
// key matches its options_key flag ("route" by default). Nothing in this
// package runs at request time.
//
// # Declaring options in .proto
//
// Depend on buf.build/go-sphere/options, then:
//
//	import "sphere/options/options.proto";
//
//	service BotService {
//	  rpc Start(StartRequest) returns (StartResponse) {
//	    option (sphere.options.options) = {
//	      key: "bot"
//	      extra: { key: "command" value: "start" }
//	    };
//	  }
//	}
//
// # Reading options in Go
//
// A generator reads the entries from a method descriptor with [E_Options]
// and selects one by key:
//
//	import (
//		"github.com/go-sphere/options/sphere/options"
//		"google.golang.org/protobuf/proto"
//	)
//
//	opts := methodDescriptor.Options() // protoreflect.MethodDescriptor
//	pairs, _ := proto.GetExtension(opts, options.E_Options).([]*options.KeyValuePair)
//	for _, kv := range pairs {
//		if kv.GetKey() == "bot" {
//			command := kv.GetExtra()["command"]
//			_ = command
//		}
//	}
//
// GetFlag, GetText, and GetNumber return zero values when the value is unset
// or holds another type; switch on GetValue to tell them apart. Messages and
// slices returned by proto.GetExtension belong to the descriptor and must not
// be modified.
package options
