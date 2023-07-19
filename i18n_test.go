package i18n_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/leonelquinteros/gotext"
	"github.com/stretchr/testify/require"
	"github.com/ubuntu/go-i18n"
	"github.com/ubuntu/go-i18n/testdata/po"
)

func TestLanguageEnVariables(t *testing.T) {
	tests := map[string]struct {
		// default is singular/translated singular
		envs map[string]string

		want string
	}{
		"LANGUAGE has precedence over all": {
			envs: map[string]string{"LANGUAGE": "fr_FR", "LC_ALL": "de_DE", "LC_MESSAGES": "en_US", "LANG": "dk"},
			want: "fr_FR",
		},
		"LC_ALL has precedence over others": {
			envs: map[string]string{"LANGUAGE": "", "LC_ALL": "de_DE", "LC_MESSAGES": "en_US", "LANG": "dk"},
			want: "de_DE",
		},
		"LC_MESSAGES has precedence over others": {
			envs: map[string]string{"LANGUAGE": "", "LC_ALL": "", "LC_MESSAGES": "en_US", "LANG": "dk"},
			want: "en_US",
		},
		"LANG has precedence over others": {
			envs: map[string]string{"LANGUAGE": "", "LC_ALL": "", "LC_MESSAGES": "", "LANG": "dk"},
			want: "dk",
		},

		"Defaults to C": {
			envs: map[string]string{"LANGUAGE": "", "LC_ALL": "", "LC_MESSAGES": "", "LANG": ""},
			want: "C",
		},

		"Locale is simplified": {
			envs: map[string]string{"LANGUAGE": "fr_FR@UTF-8"},
			want: "fr_FR",
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			i18n.InitI18nDomain("domain", nil)
			require.Equal(t, tc.want, gotext.GetLanguage(), "Should select desired language")
		})
	}
}

func TestTranslations(t *testing.T) {
	tests := map[string]struct {
		msg     string
		lang    string
		localFS fs.FS

		want string
	}{
		"Load system translations":                   {lang: "fr", want: "inconnu"},
		"System translation fallbacks to simplified": {lang: "fr_FR", want: "inconnu"},

		// Local translations.
		"Local translation wins over system":          {localFS: po.Files, lang: "fr", want: "inconnu from local"},
		"Local translation is simplified too":         {lang: "fr_FR", localFS: po.Files, want: "inconnu from local"},
		"Local translation with string not in system": {msg: "translation not in system apt", localFS: po.Files, lang: "fr", want: "traduction pas dans system apt"},
		"Complex domain wins over simplified one":     {lang: "aa_BB", localFS: po.Files, want: "nekonata from aa_BB"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv("LANGUAGE", tc.lang)

			// Check that the system has the system translation available for running the tests.
			if tc.localFS != nil {
				if _, err := os.Lstat("/usr/share/locale/fr/LC_MESSAGES/apt.mo"); err != nil {
					t.Skipf("apt translation is not available on the system: %v", err)
				}
			}

			i18n.InitI18nDomain("apt", tc.localFS)

			msg := "unknown"
			if tc.msg != "" {
				msg = tc.msg
			}

			got := gotext.Get(msg)

			require.Equal(t, tc.want, got, "Should be translated")
		})
	}
}
