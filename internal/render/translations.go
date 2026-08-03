package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/vmrocha/bible-cli/internal/bible"
)

// Translations writes bundled translation metadata and attribution.
func Translations(writer io.Writer, translations []bible.Translation, options Options) error {
	if options.Plain {
		for _, translation := range translations {
			if _, err := fmt.Fprintf(
				writer,
				"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				plainMetadataField(translation.ID),
				plainMetadataField(translation.Abbreviation),
				plainMetadataField(translation.Name),
				plainMetadataField(translation.LanguageTag),
				plainMetadataField(translation.Edition),
				plainMetadataField(translation.TextEdition),
				plainMetadataField(translation.Canon),
				plainMetadataField(translation.RightsStatus),
				plainMetadataField(translation.SourceHomepage),
				plainMetadataField(translation.RightsNoticeURL),
				plainMetadataField(translation.TrademarkNotice),
				plainMetadataField(translation.TextPolicy),
			); err != nil {
				return err
			}
		}
		return nil
	}

	if len(translations) == 0 {
		_, err := fmt.Fprintln(writer, "No translations available.")
		return err
	}

	heading := styled("Available translations", ansiBold+ansiCyan, options.Color)
	if _, err := fmt.Fprintln(writer, heading); err != nil {
		return err
	}
	for _, translation := range translations {
		abbreviation := styled(translation.Abbreviation, ansiBold+ansiCyan, options.Color)
		if _, err := fmt.Fprintf(writer, "\n%s · %s\n", abbreviation, translation.Name); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(
			writer,
			"%s (%s) · %s\n",
			translation.LanguageName,
			translation.LanguageTag,
			translation.Edition,
		); err != nil {
			return err
		}
		lines := []struct {
			label string
			value string
		}{
			{"Text edition", translation.TextEdition},
			{"Canon", translation.Canon},
			{"Rights", rightsLabel(translation.RightsStatus) + " · " + translation.RightsNoticeURL},
			{"Source", translation.SourcePublisher + " · " + translation.SourceHomepage},
			{"Trademark", translation.TrademarkNotice},
			{"Text policy", translation.TextPolicy},
		}
		for _, line := range lines {
			label := styled(line.label+":", ansiDim, options.Color)
			if _, err := fmt.Fprintf(writer, "%s %s\n", label, line.value); err != nil {
				return err
			}
		}
	}
	return nil
}

func rightsLabel(status string) string {
	if status == "public-domain" {
		return "Public domain"
	}
	return status
}

func plainMetadataField(value string) string {
	return strings.NewReplacer("\t", " ", "\r", " ", "\n", " ").Replace(value)
}
