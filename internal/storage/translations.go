package storage

import (
	"context"
	"fmt"

	"github.com/vmrocha/bible-cli/internal/bible"
)

// Translations lists every translation bundled in the embedded database.
func (reader *Reader) Translations(ctx context.Context) ([]bible.Translation, error) {
	rows, err := reader.connection.QueryContext(ctx, `
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
