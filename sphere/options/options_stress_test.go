package options

import (
	"fmt"
	"sync"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// TestKeyValuePair_NilReceiverMethods tests all method calls on a nil *KeyValuePair.
func TestKeyValuePair_NilReceiverMethods(t *testing.T) {
	var nilKV *KeyValuePair

	if nilKV.GetKey() != "" {
		t.Errorf("expected empty string, got %q", nilKV.GetKey())
	}
	if nilKV.GetValue() != nil {
		t.Errorf("expected nil value, got %v", nilKV.GetValue())
	}
	if nilKV.GetFlag() != false {
		t.Errorf("expected false, got %v", nilKV.GetFlag())
	}
	if nilKV.GetText() != "" {
		t.Errorf("expected empty string, got %q", nilKV.GetText())
	}
	if nilKV.GetNumber() != 0 {
		t.Errorf("expected 0, got %d", nilKV.GetNumber())
	}
	if nilKV.GetExtra() != nil {
		t.Errorf("expected nil extra map, got %v", nilKV.GetExtra())
	}

	if nilKV.String() == "" {
		t.Errorf("expected non-empty string for nil receiver")
	}
	raw, indices := nilKV.Descriptor()
	if len(raw) == 0 || len(indices) == 0 {
		t.Errorf("Descriptor() on nil receiver should return valid descriptor info")
	}

	ref := nilKV.ProtoReflect()
	if ref.IsValid() {
		t.Errorf("ProtoReflect on nil receiver should not be valid")
	}

	// Proto standard operations on nil receiver
	if proto.Size(nilKV) != 0 {
		t.Errorf("proto.Size(nil) should be 0, got %d", proto.Size(nilKV))
	}
	cloned := proto.Clone(nilKV)
	if cloned != nil && !proto.Equal(cloned, nilKV) {
		t.Errorf("proto.Clone(nil) mismatch: %v", cloned)
	}
	if !proto.Equal(nilKV, nilKV) {
		t.Error("proto.Equal(nil, nil) should be true")
	}
	if proto.Equal(nilKV, &KeyValuePair{}) {
		t.Error("proto.Equal(nil, &KeyValuePair{}) should be false")
	}

	data, err := proto.Marshal(nilKV)
	if err != nil {
		t.Errorf("proto.Marshal(nil) error: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("proto.Marshal(nil) should produce empty bytes, got %d bytes", len(data))
	}

	// Reset on initialized instance
	resetTarget := &KeyValuePair{
		Key:   "k",
		Value: &KeyValuePair_Text{Text: "v"},
		Extra: map[string]string{"a": "b"},
	}
	resetTarget.Reset()
	if resetTarget.GetKey() != "" || resetTarget.GetValue() != nil || resetTarget.GetExtra() != nil {
		t.Errorf("Reset did not clear fields: %+v", resetTarget)
	}
}

// TestKeyValuePair_EmptyAndZeroOneofValues tests behavior and wire serialization
// when oneof values are unset or explicitly set to their respective scalar zero values.
func TestKeyValuePair_EmptyAndZeroOneofValues(t *testing.T) {
	// 1. Completely unset oneof
	t.Run("unset oneof", func(t *testing.T) {
		kv := &KeyValuePair{Key: "unset"}
		if kv.GetValue() != nil {
			t.Errorf("expected nil GetValue(), got %v", kv.GetValue())
		}
		if kv.GetFlag() != false || kv.GetText() != "" || kv.GetNumber() != 0 {
			t.Errorf("expected zero values for all getters on unset oneof")
		}

		data, err := proto.Marshal(kv)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}
		var parsed KeyValuePair
		if err := proto.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if parsed.GetValue() != nil {
			t.Errorf("expected parsed.GetValue() to be nil, got %v", parsed.GetValue())
		}
	})

	// 2. Explicit zero boolean flag (false)
	t.Run("zero flag false", func(t *testing.T) {
		kv := &KeyValuePair{
			Key:   "flag_zero",
			Value: &KeyValuePair_Flag{Flag: false},
		}
		if kv.GetValue() == nil {
			t.Fatal("expected non-nil GetValue()")
		}
		if kv.GetFlag() != false {
			t.Errorf("expected false, got %v", kv.GetFlag())
		}

		data, err := proto.Marshal(kv)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var parsed KeyValuePair
		if err := proto.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		// In proto3 oneof, setting a zero value must preserve the oneof discriminator!
		if parsed.GetValue() == nil {
			t.Fatal("expected oneof Value to be preserved after wire roundtrip")
		}
		flagVal, ok := parsed.GetValue().(*KeyValuePair_Flag)
		if !ok || flagVal.Flag != false {
			t.Errorf("expected KeyValuePair_Flag{false}, got %v", parsed.GetValue())
		}
	})

	// 3. Explicit zero text string ("")
	t.Run("zero text empty", func(t *testing.T) {
		kv := &KeyValuePair{
			Key:   "text_empty",
			Value: &KeyValuePair_Text{Text: ""},
		}
		if kv.GetValue() == nil {
			t.Fatal("expected non-nil GetValue()")
		}
		if kv.GetText() != "" {
			t.Errorf("expected empty string, got %q", kv.GetText())
		}

		data, err := proto.Marshal(kv)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var parsed KeyValuePair
		if err := proto.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if parsed.GetValue() == nil {
			t.Fatal("expected oneof Value to be preserved after wire roundtrip")
		}
		textVal, ok := parsed.GetValue().(*KeyValuePair_Text)
		if !ok || textVal.Text != "" {
			t.Errorf("expected KeyValuePair_Text{\x22\x22}, got %v", parsed.GetValue())
		}
	})

	// 4. Explicit zero number (0)
	t.Run("zero number 0", func(t *testing.T) {
		kv := &KeyValuePair{
			Key:   "num_zero",
			Value: &KeyValuePair_Number{Number: 0},
		}
		if kv.GetValue() == nil {
			t.Fatal("expected non-nil GetValue()")
		}
		if kv.GetNumber() != 0 {
			t.Errorf("expected 0, got %d", kv.GetNumber())
		}

		data, err := proto.Marshal(kv)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var parsed KeyValuePair
		if err := proto.Unmarshal(data, &parsed); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if parsed.GetValue() == nil {
			t.Fatal("expected oneof Value to be preserved after wire roundtrip")
		}
		numVal, ok := parsed.GetValue().(*KeyValuePair_Number)
		if !ok || numVal.Number != 0 {
			t.Errorf("expected KeyValuePair_Number{0}, got %v", parsed.GetValue())
		}
	})
}

// TestKeyValuePair_ExtraMapConcurrentReadStress tests concurrent reading of
// the extra map across 50 goroutines under race detection.
func TestKeyValuePair_ExtraMapConcurrentReadStress(t *testing.T) {
	kv := &KeyValuePair{
		Key:   "shared_kv",
		Value: &KeyValuePair_Text{Text: "read_only_test"},
		Extra: map[string]string{
			"env":     "production",
			"tier":    "2",
			"region":  "us-west",
			"cluster": "alpha",
		},
	}

	const goroutines = 50
	const iterations = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				// Concurrent map reads via GetExtra()
				extra := kv.GetExtra()
				if extra["env"] != "production" {
					t.Errorf("goroutine %d: env mismatch: %s", id, extra["env"])
					return
				}
				if extra["tier"] != "2" {
					t.Errorf("goroutine %d: tier mismatch: %s", id, extra["tier"])
					return
				}

				// Concurrent iteration of the map
				count := 0
				for k, v := range extra {
					if k == "" || v == "" {
						t.Errorf("goroutine %d: empty entry", id)
					}
					count++
				}
				if count != 4 {
					t.Errorf("goroutine %d: expected 4 keys, got %d", id, count)
				}

				// Concurrent marshaling under read access
				data, err := proto.Marshal(kv)
				if err != nil {
					t.Errorf("goroutine %d: proto.Marshal failed: %v", id, err)
					return
				}
				if len(data) == 0 {
					t.Errorf("goroutine %d: empty serialized data", id)
					return
				}

				// Concurrent proto.Clone
				cloned := proto.Clone(kv).(*KeyValuePair)
				if !proto.Equal(cloned, kv) {
					t.Errorf("goroutine %d: cloned not equal", id)
					return
				}
			}
		}(g)
	}

	wg.Wait()
}

// TestKeyValuePair_ExtraMapSynchronizedWriteStress verifies that concurrent writes
// to extra map with mutual exclusion operate safely under the race detector.
func TestKeyValuePair_ExtraMapSynchronizedWriteStress(t *testing.T) {
	kv := &KeyValuePair{
		Key:   "synced_kv",
		Value: &KeyValuePair_Number{Number: 42},
		Extra: make(map[string]string),
	}

	var mu sync.RWMutex
	const writers = 20
	const readers = 20
	const iterations = 50

	var wg sync.WaitGroup
	wg.Add(writers + readers)

	// Writers
	for w := 0; w < writers; w++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				mu.Lock()
				extra := kv.GetExtra()
				extra[fmt.Sprintf("key_%d", id)] = fmt.Sprintf("val_%d_%d", id, i)
				mu.Unlock()
			}
		}(w)
	}

	// Readers
	for r := 0; r < readers; r++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				mu.RLock()
				extra := kv.GetExtra()
				_ = len(extra)
				mu.RUnlock()
			}
		}(r)
	}

	wg.Wait()

	mu.RLock()
	defer mu.RUnlock()
	if len(kv.GetExtra()) != writers {
		t.Errorf("expected %d entries, got %d", writers, len(kv.GetExtra()))
	}
}

// TestKeyValuePair_ConcurrentSerializationRoundtripStress tests end-to-end
// serialization roundtrips under 50 concurrent goroutines with -race.
func TestKeyValuePair_ConcurrentSerializationRoundtripStress(t *testing.T) {
	const goroutines = 50
	const iterations = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				key := fmt.Sprintf("opt_%d_%d", id, i)
				var val isKeyValuePair_Value
				switch (id + i) % 3 {
				case 0:
					val = &KeyValuePair_Flag{Flag: (id+i)%2 == 0}
				case 1:
					val = &KeyValuePair_Text{Text: fmt.Sprintf("text_%d_%d", id, i)}
				case 2:
					val = &KeyValuePair_Number{Number: int64(id*1000 + i)}
				}

				kv := &KeyValuePair{
					Key:   key,
					Value: val,
					Extra: map[string]string{
						"worker": fmt.Sprintf("%d", id),
						"iter":   fmt.Sprintf("%d", i),
					},
				}

				// Roundtrip KeyValuePair directly
				data, err := proto.Marshal(kv)
				if err != nil {
					t.Errorf("goroutine %d: Marshal error: %v", id, err)
					return
				}

				var parsed KeyValuePair
				if err := proto.Unmarshal(data, &parsed); err != nil {
					t.Errorf("goroutine %d: Unmarshal error: %v", id, err)
					return
				}

				if !proto.Equal(kv, &parsed) {
					t.Errorf("goroutine %d: equality failed after roundtrip", id)
					return
				}

				// Roundtrip inside MethodOptions extension
				mOpts := &descriptorpb.MethodOptions{}
				proto.SetExtension(mOpts, E_Options, []*KeyValuePair{kv})

				mBytes, err := proto.Marshal(mOpts)
				if err != nil {
					t.Errorf("goroutine %d: Marshal MethodOptions error: %v", id, err)
					return
				}

				parsedMOpts := &descriptorpb.MethodOptions{}
				if err := proto.Unmarshal(mBytes, parsedMOpts); err != nil {
					t.Errorf("goroutine %d: Unmarshal MethodOptions error: %v", id, err)
					return
				}

				if !proto.HasExtension(parsedMOpts, E_Options) {
					t.Errorf("goroutine %d: missing E_Options extension", id)
					return
				}
				retrieved := proto.GetExtension(parsedMOpts, E_Options).([]*KeyValuePair)
				if len(retrieved) != 1 || !proto.Equal(retrieved[0], kv) {
					t.Errorf("goroutine %d: extension KeyValuePair mismatch", id)
					return
				}
			}
		}(g)
	}

	wg.Wait()
}
