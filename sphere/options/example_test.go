package options_test

import (
	"fmt"

	"github.com/go-sphere/options/sphere/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// A code generator receives these options from a method descriptor. Here they
// are built by hand to show how a consumer selects its entry by key.
func Example() {
	opts := &descriptorpb.MethodOptions{}
	proto.SetExtension(opts, options.E_Options, []*options.KeyValuePair{
		{Key: "route", Value: &options.KeyValuePair_Text{Text: "/start"}},
		{Key: "bot", Extra: map[string]string{"command": "start"}},
	})

	pairs, _ := proto.GetExtension(opts, options.E_Options).([]*options.KeyValuePair)
	for _, kv := range pairs {
		if kv.GetKey() == "bot" {
			fmt.Println("command:", kv.GetExtra()["command"])
		}
	}
	// Output: command: start
}

func ExampleKeyValuePair_GetValue() {
	pairs := []*options.KeyValuePair{
		{Key: "auth", Value: &options.KeyValuePair_Flag{Flag: false}},
		{Key: "limit", Value: &options.KeyValuePair_Number{Number: 10}},
		{Key: "route", Value: &options.KeyValuePair_Text{Text: "/users"}},
		{Key: "empty"},
	}
	for _, kv := range pairs {
		// GetFlag would return false both for an explicit false and for an
		// unset value; the type switch distinguishes them.
		switch v := kv.GetValue().(type) {
		case *options.KeyValuePair_Flag:
			fmt.Println(kv.GetKey(), "flag", v.Flag)
		case *options.KeyValuePair_Number:
			fmt.Println(kv.GetKey(), "number", v.Number)
		case *options.KeyValuePair_Text:
			fmt.Println(kv.GetKey(), "text", v.Text)
		default:
			fmt.Println(kv.GetKey(), "unset")
		}
	}
	// Output:
	// auth flag false
	// limit number 10
	// route text /users
	// empty unset
}
