package importer

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/vmrocha/bible-cli/internal/bible"
	"github.com/vmrocha/bible-cli/internal/translation"
)

// LoadArchive verifies a source ZIP and parses the manifest-selected VPL member.
func LoadArchive(path string, manifest translation.Manifest) (bible.Dataset, error) {
	file, err := os.Open(path)
	if err != nil {
		return bible.Dataset{}, fmt.Errorf("open source archive: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return bible.Dataset{}, fmt.Errorf("stat source archive: %w", err)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return bible.Dataset{}, fmt.Errorf("checksum source archive: %w", err)
	}
	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	if !strings.EqualFold(actualChecksum, manifest.Source.ArchiveSHA256) {
		return bible.Dataset{}, fmt.Errorf(
			"verify source archive: SHA-256 is %s, want %s",
			actualChecksum,
			manifest.Source.ArchiveSHA256,
		)
	}

	archive, err := zip.NewReader(file, info.Size())
	if err != nil {
		return bible.Dataset{}, fmt.Errorf("open source ZIP: %w", err)
	}

	var source *zip.File
	for _, member := range archive.File {
		if member.Name != manifest.Source.ArchiveMember {
			continue
		}
		if source != nil {
			return bible.Dataset{}, fmt.Errorf("open source ZIP: duplicate member %q", member.Name)
		}
		source = member
	}
	if source == nil {
		return bible.Dataset{}, fmt.Errorf("open source ZIP: member %q not found", manifest.Source.ArchiveMember)
	}

	reader, err := source.Open()
	if err != nil {
		return bible.Dataset{}, fmt.Errorf("open VPL source: %w", err)
	}
	defer reader.Close()

	dataset, err := ParseVPL(reader, manifest.Expected)
	if err != nil {
		return bible.Dataset{}, err
	}
	if err := reader.Close(); err != nil {
		return bible.Dataset{}, fmt.Errorf("verify VPL source: %w", err)
	}

	return dataset, nil
}
