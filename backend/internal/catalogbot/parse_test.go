package catalogbot

import "testing"

func TestParseScores(t *testing.T) {
	scores, err := parseScores("20,20,40,20")
	if err != nil {
		t.Fatalf("parseScores returned error: %v", err)
	}
	if scores.Drive != 20 || scores.Focus != 20 || scores.Aesthetic != 40 || scores.Power != 20 {
		t.Fatalf("unexpected scores: %#v", scores)
	}
}

func TestParseScoresSupportsLegacyKeyValueFormat(t *testing.T) {
	scores, err := parseScores("drive:20, focus:20, aesthetic:40, power:20")
	if err != nil {
		t.Fatalf("parseScores returned error: %v", err)
	}
	if scores.Drive != 20 || scores.Focus != 20 || scores.Aesthetic != 40 || scores.Power != 20 {
		t.Fatalf("unexpected scores: %#v", scores)
	}
}

func TestParseScoresRejectsInvalidValue(t *testing.T) {
	_, err := parseScores("120,0,0,-20")
	if err == nil {
		t.Fatal("expected error for invalid score")
	}
}

func TestParseScoresRejectsInvalidSum(t *testing.T) {
	_, err := parseScores("30,30,30,30")
	if err == nil {
		t.Fatal("expected error for score sum different from 100")
	}
}

func TestParseScoresRejectsSpacesInOrderedFormat(t *testing.T) {
	_, err := parseScores("20, 20,40,20")
	if err == nil {
		t.Fatal("expected error for spaces in ordered score format")
	}
}

func TestParseVolumes(t *testing.T) {
	volumes, err := parseVolumes("3:8393, 5:12990, 10:18990")
	if err != nil {
		t.Fatalf("parseVolumes returned error: %v", err)
	}
	if len(volumes) != 3 || volumes[0].VolumeML != 3 || volumes[1].Price != 12990 {
		t.Fatalf("unexpected volumes: %#v", volumes)
	}
}

func TestParseVolumesRejectsUnsupportedVolume(t *testing.T) {
	_, err := parseVolumes("50:8393")
	if err == nil {
		t.Fatal("expected unsupported volume error")
	}
}

func TestParseGender(t *testing.T) {
	gender, err := parseGender("женский")
	if err != nil {
		t.Fatalf("parseGender returned error: %v", err)
	}
	if gender != "female" {
		t.Fatalf("expected female, got %q", gender)
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
