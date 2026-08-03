package canon

import (
	"fmt"
	"strings"

	"github.com/vmrocha/bible-cli/internal/bible"
)

// Entry combines canonical book metadata with accepted CLI aliases.
type Entry struct {
	Book    bible.Book
	Aliases []string
}

// ProtestantBooks returns the ordered 66-book catalog.
func ProtestantBooks() []Entry {
	books := make([]Entry, len(protestantBooks))
	for index, entry := range protestantBooks {
		books[index] = entry
		books[index].Aliases = append([]string(nil), entry.Aliases...)
	}
	return books
}

// Resolve returns the canonical book ID for a name, source code, ID, or alias.
func Resolve(value string) (string, bool) {
	id, ok := booksByAlias[normalize(value)]
	return id, ok
}

func normalize(value string) string {
	value = strings.NewReplacer(".", "", "_", " ", "-", " ").Replace(value)
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}

var booksByAlias = func() map[string]string {
	aliases := make(map[string]string, len(protestantBooks)*5)
	for _, entry := range protestantBooks {
		values := append([]string{entry.Book.ID, entry.Book.Name, entry.Book.SourceCode}, entry.Aliases...)
		for _, value := range values {
			key := normalize(value)
			if existing, ok := aliases[key]; ok && existing != entry.Book.ID {
				panic(fmt.Sprintf("book alias %q resolves to both %s and %s", value, existing, entry.Book.ID))
			}
			aliases[key] = entry.Book.ID
		}
	}
	return aliases
}()

var protestantBooks = []Entry{
	{Book: bible.Book{ID: "genesis", SourceCode: "GEN", Position: 1, Name: "Genesis"}, Aliases: []string{"Gen", "Ge", "Gn"}},
	{Book: bible.Book{ID: "exodus", SourceCode: "EXO", Position: 2, Name: "Exodus"}, Aliases: []string{"Ex", "Exod"}},
	{Book: bible.Book{ID: "leviticus", SourceCode: "LEV", Position: 3, Name: "Leviticus"}, Aliases: []string{"Lev", "Le"}},
	{Book: bible.Book{ID: "numbers", SourceCode: "NUM", Position: 4, Name: "Numbers"}, Aliases: []string{"Num", "Nu", "Nm"}},
	{Book: bible.Book{ID: "deuteronomy", SourceCode: "DEU", Position: 5, Name: "Deuteronomy"}, Aliases: []string{"Deut", "Dt"}},
	{Book: bible.Book{ID: "joshua", SourceCode: "JOS", Position: 6, Name: "Joshua"}, Aliases: []string{"Josh"}},
	{Book: bible.Book{ID: "judges", SourceCode: "JDG", Position: 7, Name: "Judges"}, Aliases: []string{"Judg", "Jdg"}},
	{Book: bible.Book{ID: "ruth", SourceCode: "RUT", Position: 8, Name: "Ruth"}, Aliases: []string{"Ru"}},
	{Book: bible.Book{ID: "1-samuel", SourceCode: "1SA", Position: 9, Name: "1 Samuel"}, Aliases: []string{"1 Sam", "1Sa", "1 Sm"}},
	{Book: bible.Book{ID: "2-samuel", SourceCode: "2SA", Position: 10, Name: "2 Samuel"}, Aliases: []string{"2 Sam", "2Sa", "2 Sm"}},
	{Book: bible.Book{ID: "1-kings", SourceCode: "1KI", Position: 11, Name: "1 Kings"}, Aliases: []string{"1 Kgs", "1Ki"}},
	{Book: bible.Book{ID: "2-kings", SourceCode: "2KI", Position: 12, Name: "2 Kings"}, Aliases: []string{"2 Kgs", "2Ki"}},
	{Book: bible.Book{ID: "1-chronicles", SourceCode: "1CH", Position: 13, Name: "1 Chronicles"}, Aliases: []string{"1 Chr", "1Ch"}},
	{Book: bible.Book{ID: "2-chronicles", SourceCode: "2CH", Position: 14, Name: "2 Chronicles"}, Aliases: []string{"2 Chr", "2Ch"}},
	{Book: bible.Book{ID: "ezra", SourceCode: "EZR", Position: 15, Name: "Ezra"}, Aliases: []string{"Ezr"}},
	{Book: bible.Book{ID: "nehemiah", SourceCode: "NEH", Position: 16, Name: "Nehemiah"}, Aliases: []string{"Neh"}},
	{Book: bible.Book{ID: "esther", SourceCode: "EST", Position: 17, Name: "Esther"}, Aliases: []string{"Est", "Esth"}},
	{Book: bible.Book{ID: "job", SourceCode: "JOB", Position: 18, Name: "Job"}, Aliases: []string{"Jb"}},
	{Book: bible.Book{ID: "psalms", SourceCode: "PSA", Position: 19, Name: "Psalms"}, Aliases: []string{"Ps", "Psalm", "Psa"}},
	{Book: bible.Book{ID: "proverbs", SourceCode: "PRO", Position: 20, Name: "Proverbs"}, Aliases: []string{"Prov", "Pr"}},
	{Book: bible.Book{ID: "ecclesiastes", SourceCode: "ECC", Position: 21, Name: "Ecclesiastes"}, Aliases: []string{"Eccl", "Ecc"}},
	{Book: bible.Book{ID: "song-of-solomon", SourceCode: "SOL", Position: 22, Name: "Song of Solomon"}, Aliases: []string{"Song", "Song of Songs", "SOS", "Canticles"}},
	{Book: bible.Book{ID: "isaiah", SourceCode: "ISA", Position: 23, Name: "Isaiah"}, Aliases: []string{"Isa"}},
	{Book: bible.Book{ID: "jeremiah", SourceCode: "JER", Position: 24, Name: "Jeremiah"}, Aliases: []string{"Jer"}},
	{Book: bible.Book{ID: "lamentations", SourceCode: "LAM", Position: 25, Name: "Lamentations"}, Aliases: []string{"Lam"}},
	{Book: bible.Book{ID: "ezekiel", SourceCode: "EZE", Position: 26, Name: "Ezekiel"}, Aliases: []string{"Ezek", "Eze"}},
	{Book: bible.Book{ID: "daniel", SourceCode: "DAN", Position: 27, Name: "Daniel"}, Aliases: []string{"Dan", "Dn"}},
	{Book: bible.Book{ID: "hosea", SourceCode: "HOS", Position: 28, Name: "Hosea"}, Aliases: []string{"Hos"}},
	{Book: bible.Book{ID: "joel", SourceCode: "JOE", Position: 29, Name: "Joel"}, Aliases: []string{"Joe"}},
	{Book: bible.Book{ID: "amos", SourceCode: "AMO", Position: 30, Name: "Amos"}, Aliases: []string{"Am"}},
	{Book: bible.Book{ID: "obadiah", SourceCode: "OBA", Position: 31, Name: "Obadiah"}, Aliases: []string{"Obad", "Ob"}},
	{Book: bible.Book{ID: "jonah", SourceCode: "JON", Position: 32, Name: "Jonah"}, Aliases: []string{"Jon"}},
	{Book: bible.Book{ID: "micah", SourceCode: "MIC", Position: 33, Name: "Micah"}, Aliases: []string{"Mic"}},
	{Book: bible.Book{ID: "nahum", SourceCode: "NAH", Position: 34, Name: "Nahum"}, Aliases: []string{"Nah"}},
	{Book: bible.Book{ID: "habakkuk", SourceCode: "HAB", Position: 35, Name: "Habakkuk"}, Aliases: []string{"Hab"}},
	{Book: bible.Book{ID: "zephaniah", SourceCode: "ZEP", Position: 36, Name: "Zephaniah"}, Aliases: []string{"Zeph", "Zep"}},
	{Book: bible.Book{ID: "haggai", SourceCode: "HAG", Position: 37, Name: "Haggai"}, Aliases: []string{"Hag"}},
	{Book: bible.Book{ID: "zechariah", SourceCode: "ZEC", Position: 38, Name: "Zechariah"}, Aliases: []string{"Zech", "Zec"}},
	{Book: bible.Book{ID: "malachi", SourceCode: "MAL", Position: 39, Name: "Malachi"}, Aliases: []string{"Mal"}},
	{Book: bible.Book{ID: "matthew", SourceCode: "MAT", Position: 40, Name: "Matthew"}, Aliases: []string{"Matt", "Mt"}},
	{Book: bible.Book{ID: "mark", SourceCode: "MAR", Position: 41, Name: "Mark"}, Aliases: []string{"Mk", "Mrk"}},
	{Book: bible.Book{ID: "luke", SourceCode: "LUK", Position: 42, Name: "Luke"}, Aliases: []string{"Lk"}},
	{Book: bible.Book{ID: "john", SourceCode: "JOH", Position: 43, Name: "John"}, Aliases: []string{"Jn", "Jhn"}},
	{Book: bible.Book{ID: "acts", SourceCode: "ACT", Position: 44, Name: "Acts"}, Aliases: []string{"Ac"}},
	{Book: bible.Book{ID: "romans", SourceCode: "ROM", Position: 45, Name: "Romans"}, Aliases: []string{"Rom", "Ro"}},
	{Book: bible.Book{ID: "1-corinthians", SourceCode: "1CO", Position: 46, Name: "1 Corinthians"}, Aliases: []string{"1 Cor", "1Co"}},
	{Book: bible.Book{ID: "2-corinthians", SourceCode: "2CO", Position: 47, Name: "2 Corinthians"}, Aliases: []string{"2 Cor", "2Co"}},
	{Book: bible.Book{ID: "galatians", SourceCode: "GAL", Position: 48, Name: "Galatians"}, Aliases: []string{"Gal"}},
	{Book: bible.Book{ID: "ephesians", SourceCode: "EPH", Position: 49, Name: "Ephesians"}, Aliases: []string{"Eph"}},
	{Book: bible.Book{ID: "philippians", SourceCode: "PHI", Position: 50, Name: "Philippians"}, Aliases: []string{"Phil", "Php"}},
	{Book: bible.Book{ID: "colossians", SourceCode: "COL", Position: 51, Name: "Colossians"}, Aliases: []string{"Col"}},
	{Book: bible.Book{ID: "1-thessalonians", SourceCode: "1TH", Position: 52, Name: "1 Thessalonians"}, Aliases: []string{"1 Thess", "1Th"}},
	{Book: bible.Book{ID: "2-thessalonians", SourceCode: "2TH", Position: 53, Name: "2 Thessalonians"}, Aliases: []string{"2 Thess", "2Th"}},
	{Book: bible.Book{ID: "1-timothy", SourceCode: "1TI", Position: 54, Name: "1 Timothy"}, Aliases: []string{"1 Tim", "1Ti"}},
	{Book: bible.Book{ID: "2-timothy", SourceCode: "2TI", Position: 55, Name: "2 Timothy"}, Aliases: []string{"2 Tim", "2Ti"}},
	{Book: bible.Book{ID: "titus", SourceCode: "TIT", Position: 56, Name: "Titus"}, Aliases: []string{"Tit"}},
	{Book: bible.Book{ID: "philemon", SourceCode: "PHM", Position: 57, Name: "Philemon"}, Aliases: []string{"Phlm", "Phm"}},
	{Book: bible.Book{ID: "hebrews", SourceCode: "HEB", Position: 58, Name: "Hebrews"}, Aliases: []string{"Heb"}},
	{Book: bible.Book{ID: "james", SourceCode: "JAM", Position: 59, Name: "James"}, Aliases: []string{"Jas", "Jam"}},
	{Book: bible.Book{ID: "1-peter", SourceCode: "1PE", Position: 60, Name: "1 Peter"}, Aliases: []string{"1 Pet", "1Pe"}},
	{Book: bible.Book{ID: "2-peter", SourceCode: "2PE", Position: 61, Name: "2 Peter"}, Aliases: []string{"2 Pet", "2Pe"}},
	{Book: bible.Book{ID: "1-john", SourceCode: "1JO", Position: 62, Name: "1 John"}, Aliases: []string{"1 Jn", "1Jn", "1Jo"}},
	{Book: bible.Book{ID: "2-john", SourceCode: "2JO", Position: 63, Name: "2 John"}, Aliases: []string{"2 Jn", "2Jn", "2Jo"}},
	{Book: bible.Book{ID: "3-john", SourceCode: "3JO", Position: 64, Name: "3 John"}, Aliases: []string{"3 Jn", "3Jn", "3Jo"}},
	{Book: bible.Book{ID: "jude", SourceCode: "JUD", Position: 65, Name: "Jude"}, Aliases: []string{"Jud"}},
	{Book: bible.Book{ID: "revelation", SourceCode: "REV", Position: 66, Name: "Revelation"}, Aliases: []string{"Rev", "Re", "Apocalypse"}},
}
