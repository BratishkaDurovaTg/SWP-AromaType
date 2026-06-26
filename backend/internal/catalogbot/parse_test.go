package catalogbot

import "testing"

func TestParseScores(t *testing.T) {
	scores, err := parseScores("drive:20, focus:35, aesthetic:90, power:25")
	if err != nil {
		t.Fatalf("parseScores returned error: %v", err)
	}
	if scores.Drive != 20 || scores.Focus != 35 || scores.Aesthetic != 90 || scores.Power != 25 {
		t.Fatalf("unexpected scores: %#v", scores)
	}
}

func TestParseScoresRejectsInvalidValue(t *testing.T) {
	_, err := parseScores("drive:120")
	if err == nil {
		t.Fatal("expected error for score over 100")
	}
}

func TestParseVolumes(t *testing.T) {
	volumes, err := parseVolumes("50:8393, 100:12990")
	if err != nil {
		t.Fatalf("parseVolumes returned error: %v", err)
	}
	if len(volumes) != 2 || volumes[0].VolumeML != 50 || volumes[1].Price != 12990 {
		t.Fatalf("unexpected volumes: %#v", volumes)
	}
}

func TestValidateID(t *testing.T) {
	id, err := validateID("Miami-Shake")
	if err != nil {
		t.Fatalf("validateID returned error: %v", err)
	}
	if id != "miami-shake" {
		t.Fatalf("expected normalized id, got %q", id)
	}
}
