// Package i18n is responsible for internationalization/translation handling and generation.
package i18n

import (
	"io"
	"io/fs"
	"os"
	"sync"

	"github.com/leonelquinteros/gotext"

	_ "unsafe"
)

// gotextConfig is the config from gotext.
type gotextConfig struct {
	sync.RWMutex

	// Default domain to look at when no domain is specified. Used by package level functions.
	domain string

	// Language set.
	language string

	// Path to library directory where all locale directories and Translation files are.
	library string

	// Storage for package level methods
	storage *gotext.Locale
}

//go:linkname globalGotextConfig github.com/leonelquinteros/gotext.globalConfig
var globalGotextConfig *gotextConfig

// InitI18nDomain loads domain for the user current locale.
// If a poDir is passed as fs.FS, then, it will override for the domain any translations
// potentially present on disk.
func InitI18nDomain(domain string, poDir fs.FS) {
	lang := getCurrentLanguage()
	if lang == "C" {
		return
	}

	// System configuration.
	loadFromSystem(lang, domain)

	// Override with embedded po.
	loadFromEmbeddedPos(poDir, lang, domain)
}

// getCurrentLanguage returns the language name from the system.
func getCurrentLanguage() string {
	for _, k := range []string{"LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"} {
		if lang := os.Getenv(k); lang != "" {
			return gotext.SimplifiedLocale(lang)
		}
	}

	return "C"
}

// loadFromSystem uses known system path to load l10n translations.
func loadFromSystem(lang, domain string) {
	for _, p := range []string{"/usr/local/share/locales", "/usr/share/locale", "/usr/share/locale-langpack"} {
		gotext.Configure(p, lang, domain)
		// Stop as soon as we found something to load for this domain.
		if len(globalGotextConfig.storage.Domains) > 0 {
			break
		}
	}
}

// loadFromEmbeddedPos loads any po files if embedded in the FS directory.
func loadFromEmbeddedPos(poDir fs.FS, lang, domain string) {
	if poDir == nil {
		return
	}

	for _, language := range []string{lang, lang[:2]} {
		f, err := poDir.Open(language + ".po")
		if err != nil {
			continue
		}

		buf, err := io.ReadAll(f)
		if err != nil {
			// TODO: log here with 1.21 (as file opened but can't read)
			continue
		}

		translator := gotext.NewPo()
		translator.Parse(buf)
		globalGotextConfig.storage.AddTranslator(domain, translator)

		return
	}
}
