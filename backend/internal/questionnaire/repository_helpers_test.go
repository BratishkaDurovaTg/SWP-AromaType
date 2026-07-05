package questionnaire

import (
	"reflect"
	"testing"
)

func TestStringArrayJSONHelpers(t *testing.T) {
	encoded, err := encodeStringArray([]string{"Клубника", "Мороженое"})
	if err != nil {
		t.Fatalf("encodeStringArray returned error: %v", err)
	}

	var decoded []string
	if err := decodeStringArray([]byte(encoded), &decoded); err != nil {
		t.Fatalf("decodeStringArray returned error: %v", err)
	}

	want := []string{"Клубника", "Мороженое"}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("unexpected decoded strings: got %#v want %#v", decoded, want)
	}
}

func TestStringArrayDecodeEmptyDefaultsToEmptySlice(t *testing.T) {
	var decoded []string
	if err := decodeStringArray(nil, &decoded); err != nil {
		t.Fatalf("decodeStringArray returned error: %v", err)
	}
	if decoded == nil || len(decoded) != 0 {
		t.Fatalf("expected empty slice, got %#v", decoded)
	}
}

func TestVolumeOptionsJSONHelpers(t *testing.T) {
	encoded, err := encodeVolumeOptions([]VolumeOption{{VolumeML: 3, Price: 8393}})
	if err != nil {
		t.Fatalf("encodeVolumeOptions returned error: %v", err)
	}

	var decoded []VolumeOption
	if err := decodeVolumeOptions([]byte(encoded), &decoded); err != nil {
		t.Fatalf("decodeVolumeOptions returned error: %v", err)
	}

	want := []VolumeOption{{VolumeML: 3, Price: 8393}}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("unexpected decoded volumes: got %#v want %#v", decoded, want)
	}
}

func TestPsychotypeScoresJSONHelpers(t *testing.T) {
	scores := PsychotypeScores{Drive: 10, Focus: 20, Aesthetic: 90, Power: 30}
	encoded, err := encodePsychotypeScores(scores)
	if err != nil {
		t.Fatalf("encodePsychotypeScores returned error: %v", err)
	}

	var decoded PsychotypeScores
	if err := decodePsychotypeScores([]byte(encoded), &decoded); err != nil {
		t.Fatalf("decodePsychotypeScores returned error: %v", err)
	}
	if decoded != scores {
		t.Fatalf("unexpected decoded scores: got %#v want %#v", decoded, scores)
	}
}
