// package main implements update-po command line to update one domain pot file and merge with any existing translations.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	"github.com/leonelquinteros/gotext"
	"github.com/leonelquinteros/gotext/cli/xgotext/parser"
	pkg_tree "github.com/leonelquinteros/gotext/cli/xgotext/parser/pkg-tree"
)

func main() {
	verbose := flag.Bool("v", false, "Print currently handled directory pkg")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s domain /path/to/po/dir /path/to/first/pkg [/path/to/other/pkg …]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) < 3 {
		fmt.Fprintln(os.Stderr, "ERROR: Missing input arguments")
		flag.Usage()
		os.Exit(2)
	}

	domain := flag.Arg(0)
	poDir := flag.Arg(1)
	potFile := filepath.Join(poDir, domain+".pot")

	// Create or update pot file.
	if err := i18nToPot(domain, flag.Args()[2:], potFile, *verbose); err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}

	// Update existing po files.
	files, err := os.ReadDir(poDir)
	if err != nil {
		fmt.Printf("ERROR: can't read %q: %v", poDir, err)
		os.Exit(1)
	}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".po") {
			continue
		}
		poFilePath := filepath.Join(poDir, f.Name())
		if err := updatePo(poFilePath, potFile); err != nil {
			fmt.Printf("ERROR: can't update %q: %v", poFilePath, err)
		}
	}
}

// i18nToPot extract i18n entries and save it into the pot file.
func i18nToPot(domain string, pkgs []string, potFile string, verbose bool) error {
	data := &parser.DomainMap{
		Default: domain,
	}

	for _, pkgPath := range pkgs {
		err := pkg_tree.ParsePkgTree(pkgPath, data, verbose)
		if err != nil {
			return fmt.Errorf("can't parse packages: %v", err)
		}
	}
	if _, ok := data.Domains[domain]; !ok {
		return fmt.Errorf("no strings marked up for i18n for %s", domain)
	}

	err := os.MkdirAll(filepath.Dir(potFile), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output dir: %v", err)
	}
	f, err := os.Create(potFile)
	if err != nil {
		return fmt.Errorf("failed to create domain file: %v", err)
	}
	defer f.Close()

	// write header
	if _, err := fmt.Fprintf(f, `msgid ""
msgstr ""
"Project-Id-Version: %s\n"
"Language: \n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
"Content-Transfer-Encoding: 8bit\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"

`, domain); err != nil {
		return err
	}

	// write domain content
	if _, err := f.WriteString(data.Domains[domain].Dump()); err != nil {
		return fmt.Errorf("can't unmarshall translations to pot file: %v", err)
	}

	return nil
}

func updatePo(poPath, potPath string) error {
	pot := gotext.NewPo()
	buf, err := os.ReadFile(potPath)
	if err != nil {
		return err
	}
	pot.Parse(buf)
	potTranslations := pot.GetDomain().GetTranslations()

	// We keep original po file headers, domain and references
	localizedPo := gotext.NewPo()
	buf, err = os.ReadFile(poPath)
	if err != nil {
		return err
	}
	localizedPo.Parse(buf)
	existingTranslations := localizedPo.GetDomain().GetTranslations()

	// Reuse existing translations from original file, or take the one from the .pot file.
	newTranslations := make(map[string]*gotext.Translation, len(potTranslations))
	for id, trans := range potTranslations {
		var newTrans gotext.Translation
		if existingsTrans, ok := existingTranslations[id]; ok {
			newTrans = *existingsTrans
		} else {
			newTrans = *trans
		}

		newTranslations[id] = &newTrans
	}

	// domain is a pointer, swap out the old translations with the new ones.
	// That also has the benefit to purge deprecated translation.
	domain := localizedPo.GetDomain()
	field := reflect.ValueOf(domain).Elem().FieldByName("translations")
	//nolint:gosec // We are changing a private pointer to avoid remapping against "n" which changes
	// between languages. This is a generator, limited to po generation, so any breakage will be noticed
	// before rolling out to production.
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(newTranslations))

	data, err := localizedPo.MarshalText()
	if err != nil {
		return fmt.Errorf("could not marshal po file %s to text: %v", poPath, err)
	}

	if err := os.WriteFile(poPath, data, 0600); err != nil {
		return fmt.Errorf("could not save to %q: %v", poPath, err)
	}

	return nil
}
