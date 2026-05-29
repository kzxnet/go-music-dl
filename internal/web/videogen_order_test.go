package web

import (
	"strings"
	"testing"
)

func TestVideogenUsesRomajiBeforeTranslationOrder(t *testing.T) {
	content, err := templateFS.ReadFile("templates/static/js/videogen.js")
	if err != nil {
		t.Fatalf("ReadFile(videogen.js): %v", err)
	}

	js := string(content)
	if strings.Contains(js, "const [orig, trans, roma] = group.lines") {
		t.Fatal("videogen.js still assumes the old original/translation/romaji order")
	}
	for _, want := range []string{
		"function splitLyricGroupLinesWorker(lines)",
		"function splitLyricGroupLines(lines)",
		"const { orig, roma, trans } = splitLyricGroupLinesWorker(group.lines)",
		"const { orig, roma, trans } = splitLyricGroupLines(group.lines)",
		"renderKaraokeLineHTML(roma, 'vg-line-roma'",
		"renderKaraokeLineHTML(trans, 'vg-line-trans'",
	} {
		if !strings.Contains(js, want) {
			t.Fatalf("videogen.js missing %q", want)
		}
	}
}

func TestVideogenRenderUploadsBinaryFrameBatches(t *testing.T) {
	content, err := templateFS.ReadFile("templates/static/js/videogen.js")
	if err != nil {
		t.Fatalf("ReadFile(videogen.js): %v", err)
	}

	js := string(content)
	for _, want := range []string{
		"targetCanvas.toBlob",
		"const form = new FormData()",
		`form.append("frames", blob,`,
		"framesBuffer.push(await canvasToJpegBlob(canvas, 0.92))",
	} {
		if !strings.Contains(js, want) {
			t.Fatalf("videogen.js missing %q", want)
		}
	}
	if strings.Contains(js, `framesBuffer.push(canvas.toDataURL("image/jpeg", 0.92))`) {
		t.Fatal("videogen.js still uploads render frames as base64 data URLs")
	}
}
