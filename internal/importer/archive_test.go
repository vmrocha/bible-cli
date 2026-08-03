package importer

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vmrocha/bible-cli/internal/translation"
)

func TestLoadArchive(t *testing.T) {
	manifest, archivePath := archiveFixture(t, "GEN 1:1 First.\nGEN 1:2 Second.")

	dataset, err := LoadArchive(archivePath, manifest)
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if len(dataset.Verses) != 2 {
		t.Fatalf("got %d verses, want 2", len(dataset.Verses))
	}
}

func TestLoadArchiveRejectsWrongChecksum(t *testing.T) {
	manifest, archivePath := archiveFixture(t, "GEN 1:1 First.\nGEN 1:2 Second.")
	manifest.Source.ArchiveSHA256 = strings.Repeat("0", 64)

	_, err := LoadArchive(archivePath, manifest)
	if err == nil || !strings.Contains(err.Error(), "SHA-256") {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func archiveFixture(t *testing.T, source string) (translation.Manifest, string) {
	t.Helper()

	var contents bytes.Buffer
	archive := zip.NewWriter(&contents)
	member, err := archive.Create("source.txt")
	if err != nil {
		t.Fatalf("create ZIP member: %v", err)
	}
	if _, err := member.Write([]byte(source)); err != nil {
		t.Fatalf("write ZIP member: %v", err)
	}
	if err := archive.Close(); err != nil {
		t.Fatalf("close ZIP: %v", err)
	}

	path := filepath.Join(t.TempDir(), "source.zip")
	if err := os.WriteFile(path, contents.Bytes(), 0o600); err != nil {
		t.Fatalf("write ZIP fixture: %v", err)
	}
	digest := sha256.Sum256(contents.Bytes())

	return translation.Manifest{
		Source: translation.Source{
			ArchiveMember: "source.txt",
			ArchiveSHA256: hex.EncodeToString(digest[:]),
		},
		Expected: translation.Expected{
			Books:          1,
			Verses:         2,
			FirstReference: "GEN 1:1",
			LastReference:  "GEN 1:2",
		},
	}, path
}
