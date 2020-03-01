package gifview

import (
	"path/filepath"
	"testing"
)

func TestFromImagePath(t *testing.T) {
	// Exercise failure paths
	errorTests := []struct {
		path string
		err  string
	}{
		{"missing.gif", "Unable to open file: open testdata/missing.gif: no such file or directory"},
		{"bad.gif", "Unable to decode GIF: gif: reading header: unexpected EOF"},
		{"odd.gif", "Unable to convert frame 0: pixelview: Can't process image with uneven height"},
	}
	for _, s := range errorTests {
		t.Run(s.path, func(t *testing.T) {
			_, err := FromImagePath(filepath.Join("testdata", s.path))
			if err.Error() != s.err {
				t.Errorf("Expected: %v, Got: %v", s.err, err)
			}
		})
	}

	// Load a good gif
	t.Run("good.gif", func(t *testing.T) {
		_, err := FromImagePath(filepath.Join("testdata", "good.gif"))
		if err != nil {
			t.Errorf("Expected: nil, Got: %v", err)
		}
	})
}
