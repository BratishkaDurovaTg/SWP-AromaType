package httpapi

import "testing"

func TestAllowedPhotoExtension(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		extension   string
		allowed     bool
	}{
		{name: "jpeg", contentType: "image/jpeg", extension: ".jpg", allowed: true},
		{name: "png", contentType: "image/png", extension: ".png", allowed: true},
		{name: "webp", contentType: "image/webp", extension: ".webp", allowed: true},
		{name: "svg", contentType: "image/svg+xml", allowed: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			extension, allowed := allowedPhotoExtension(test.contentType)
			if allowed != test.allowed {
				t.Fatalf("expected allowed=%v, got %v", test.allowed, allowed)
			}
			if extension != test.extension {
				t.Fatalf("expected extension %q, got %q", test.extension, extension)
			}
		})
	}
}
