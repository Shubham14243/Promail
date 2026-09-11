package services

import "testing"

func TestAddClickTrackingBodyRejectsUnsafeURLs(t *testing.T) {
	input := `<a href="javascript:alert(1)">unsafe</a> <a href="https://example.com/path">safe</a>`

	got, trackings := AddClickTrackingBody(input)

	if len(trackings) != 1 {
		t.Fatalf("expected one tracked URL, got %d", len(trackings))
	}
	if trackings[0].OriginalURL != "https://example.com/path" {
		t.Fatalf("unexpected tracked URL: %q", trackings[0].OriginalURL)
	}
	if got == "" || got == input {
		t.Fatal("expected the safe URL to be rewritten")
	}
}
