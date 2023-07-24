// package main implements compile-mo command line to compile localised po file to mo via gettext.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s domain /path/to/po/dir /path/to/mo/dir\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) != 3 {
		fmt.Fprintln(os.Stderr, "ERROR: Incorrect number of arguments")
		flag.Usage()
		os.Exit(2)
	}

	if err := generateMos(flag.Arg(0), flag.Arg(1), flag.Arg(2)); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func generateMos(domain, poDir, moDir string) error {
	err := os.MkdirAll(moDir, 0750)
	if err != nil {
		return fmt.Errorf("can’t create destination directory: %v", err)
	}

	files, err := os.ReadDir(poDir)
	if err != nil {
		return fmt.Errorf("can't read %q: %v", poDir, err)
	}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".po") {
			continue
		}

		poFile := filepath.Join(poDir, f.Name())
		outDir := filepath.Join(moDir, strings.TrimSuffix(f.Name(), ".po"), "LC_MESSAGES")
		if err := os.MkdirAll(outDir, 0750); err != nil {
			return fmt.Errorf("couldn't create %q mo file: %v", f.Name(), err)
		}
		//nolint: gosec // this is only use for file generation, not in production code.
		if out, err := exec.Command("msgfmt", "--output-file="+filepath.Join(outDir, domain+".mo"),
			poFile).CombinedOutput(); err != nil {
			return fmt.Errorf("couldn't compile mo file from %q: %v.\nCommand output: %s", poFile, err, out)
		}
	}

	return nil
}
