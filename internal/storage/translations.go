package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/vmrocha/bible-cli/internal/bible"
)

// Translations lists every registered embedded translation.
func (reader *Reader) Translations(ctx context.Context) ([]bible.Translation, error) {
	var translations []bible.Translation
	for _, id := range embeddedTranslationIDs() {
		if id == reader.translationID {
			entries, err := translationsFromConnection(ctx, reader.connection)
			if err != nil {
				return nil, err
			}
			translations = append(translations, entries...)
			continue
		}

		other, err := OpenEmbedded(ctx, id)
		if err != nil {
			return nil, err
		}
		entries, readErr := translationsFromConnection(ctx, other.connection)
		closeErr := other.Close()
		if readErr != nil || closeErr != nil {
			return nil, errors.Join(readErr, closeErr)
		}
		translations = append(translations, entries...)
	}
	sort.Slice(translations, func(i, j int) bool {
		if translations[i].LanguageTag != translations[j].LanguageTag {
			return translations[i].LanguageTag < translations[j].LanguageTag
		}
		if translations[i].Name != translations[j].Name {
			return translations[i].Name < translations[j].Name
		}
		return translations[i].ID < translations[j].ID
	})
	return translations, nil
}

func translationsFromConnection(ctx context.Context, connection *sql.Conn) ([]bible.Translation, error) {
	rows, err := connection.QueryContext(ctx, `
        SELECT
            id,
            name,
            abbreviation,
            language_tag,
            language_name,
            edition,
            canon,
            text_edition,
            source_publisher,
            source_homepage,
            rights_status,
            rights_notice_url,
            trademark_notice,
            text_policy
        FROM translations
        ORDER BY language_tag, name, id
    `)
	if err != nil {
		return nil, fmt.Errorf("list translations: %w", err)
	}
	defer rows.Close()

	var translations []bible.Translation
	for rows.Next() {
		var translation bible.Translation
		if err := rows.Scan(
			&translation.ID,
			&translation.Name,
			&translation.Abbreviation,
			&translation.LanguageTag,
			&translation.LanguageName,
			&translation.Edition,
			&translation.Canon,
			&translation.TextEdition,
			&translation.SourcePublisher,
			&translation.SourceHomepage,
			&translation.RightsStatus,
			&translation.RightsNoticeURL,
			&translation.TrademarkNotice,
			&translation.TextPolicy,
		); err != nil {
			return nil, fmt.Errorf("scan translation: %w", err)
		}
		translations = append(translations, translation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read translation rows: %w", err)
	}
	return translations, nil
}
