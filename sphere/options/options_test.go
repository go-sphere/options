package options_test

import (
	"testing"

	"github.com/go-sphere/options/sphere/options"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestKeyValuePair_GettersAndNilSafety(t *testing.T) {
	var nilKV *options.KeyValuePair
	if nilKV.GetKey() != "" {
		t.Errorf("expected empty string for nil Key")
	}
	if nilKV.GetValue() != nil {
		t.Errorf("expected nil for nil Value")
	}
	if nilKV.GetFlag() != false {
		t.Errorf("expected false for nil Flag")
	}
	if nilKV.GetText() != "" {
		t.Errorf("expected empty string for nil Text")
	}
	if nilKV.GetNumber() != 0 {
		t.Errorf("expected 0 for nil Number")
	}
	if nilKV.GetExtra() != nil {
		t.Errorf("expected nil for nil Extra")
	}
}

func TestKeyValuePair_OneofValueDiscrimination(t *testing.T) {
	// Flag
	kvFlag := &options.KeyValuePair{
		Key:   "auth_required",
		Value: &options.KeyValuePair_Flag{Flag: true},
	}
	if !kvFlag.GetFlag() {
		t.Errorf("expected flag to be true")
	}
	if kvFlag.GetText() != "" {
		t.Errorf("expected empty text when flag is set")
	}
	if kvFlag.GetNumber() != 0 {
		t.Errorf("expected 0 number when flag is set")
	}

	// Text
	kvText := &options.KeyValuePair{
		Key:   "role",
		Value: &options.KeyValuePair_Text{Text: "admin"},
	}
	if kvText.GetText() != "admin" {
		t.Errorf("expected text 'admin'")
	}
	if kvText.GetFlag() != false {
		t.Errorf("expected false flag when text is set")
	}

	// Number
	kvNum := &options.KeyValuePair{
		Key:   "rate_limit",
		Value: &options.KeyValuePair_Number{Number: 100},
	}
	if kvNum.GetNumber() != 100 {
		t.Errorf("expected number 100")
	}
	if kvNum.GetFlag() != false {
		t.Errorf("expected false flag when number is set")
	}
}

func TestKeyValuePair_ExtraAndMutation(t *testing.T) {
	kv := &options.KeyValuePair{
		Key:   "bot",
		Extra: map[string]string{"command": "start", "desc": "Start bot"},
	}
	extra := kv.GetExtra()
	if extra["command"] != "start" {
		t.Errorf("expected start, got %s", extra["command"])
	}

	extra["new_key"] = "new_val"
	if kv.GetExtra()["new_key"] != "new_val" {
		t.Errorf("expected mutation to reflect in underlying map")
	}
}

func TestKeyValuePair_ProtoRoundTrip(t *testing.T) {
	orig := &options.KeyValuePair{
		Key:   "timeout",
		Value: &options.KeyValuePair_Number{Number: 30},
		Extra: map[string]string{"unit": "seconds"},
	}

	data, err := proto.Marshal(orig)
	if err != nil {
		t.Fatalf("proto.Marshal failed: %v", err)
	}

	var parsed options.KeyValuePair
	if err := proto.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("proto.Unmarshal failed: %v", err)
	}

	if parsed.GetKey() != orig.GetKey() ||
		parsed.GetNumber() != orig.GetNumber() ||
		parsed.GetExtra()["unit"] != "seconds" {
		t.Errorf("roundtrip mismatch: got %+v, want %+v", &parsed, orig)
	}
}

func TestOptionsExtension(t *testing.T) {
	methodOpts := &descriptorpb.MethodOptions{}
	optsList := []*options.KeyValuePair{
		{Key: "route", Value: &options.KeyValuePair_Text{Text: "/v1/test"}},
		{Key: "public", Value: &options.KeyValuePair_Flag{Flag: true}},
	}
	proto.SetExtension(methodOpts, options.E_Options, optsList)

	if !proto.HasExtension(methodOpts, options.E_Options) {
		t.Fatal("expected Options extension to be set")
	}
	retrieved := proto.GetExtension(methodOpts, options.E_Options).([]*options.KeyValuePair)
	if len(retrieved) != 2 {
		t.Fatalf("expected 2 items, got %d", len(retrieved))
	}
	if retrieved[0].GetKey() != "route" || retrieved[0].GetText() != "/v1/test" {
		t.Errorf("unexpected first item: %+v", retrieved[0])
	}
}
