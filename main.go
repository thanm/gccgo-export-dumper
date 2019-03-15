//
// Reads gccgo-generated objectfiles/archives and dumps contents of
// export data.
//

package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var verbflag = flag.Int("v", 0, "Verbose trace output level")

func verb(vlevel int, s string, a ...interface{}) {
	if *verbflag >= vlevel {
		fmt.Printf(s, a...)
		fmt.Printf("\n")
	}
}

func warn(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, s, a...)
	fmt.Fprintf(os.Stderr, "\n")
}

func usage(msg string) {
	if len(msg) > 0 {
		fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	}
	fmt.Fprintf(os.Stderr, "usage: gccgo-export-dumper [flags] { files }\n")
	fmt.Fprintf(os.Stderr, `
Works on either objects (.o files) or archives (.a files); dumps out
gccgo export data for a package in text form, to standard output. Example:

  cd $GOPATH/src/mumble
  go build -o mumble.a -compiler gccgo .
  gccgo-export-dumper mumble.a
`)
	flag.PrintDefaults()
	os.Exit(2)
}

func examineElf(infile, objfile string) error {
	verb(1, "examining ELF file %s", objfile)
	var f *elf.File
	var err error
	if f, err = elf.Open(objfile); err != nil {
		return fmt.Errorf("elf.Open(%s) failed: %v", infile, err)
	}

	// Locate the .go_export section and dump it.
	var exportSection *elf.Section
	if exportSection = f.Section(".go_export"); exportSection == nil {
		warn("warning: %s contains no .go_export", infile)
		return nil
	}
	var edata []byte
	if edata, err = exportSection.Data(); err != nil {
		return fmt.Errorf("unable to extract export data from %s: %v",
			infile, err)
	}
	fmt.Printf(string(edata))
	fmt.Printf("\n")
	return nil
}

func examine(filename string) error {
	verb(1, "examining file %s", filename)

	objectfile := filename

	// Run "ar t" on the file and examine the exit status.
	cmd := exec.Command("ar", "t", filename)
	if _, err := cmd.CombinedOutput(); err != nil {
		// Some sort of error took place. Assume that the
		// input file is an on object file.
		verb(1, "assuming %s is object file", filename)
	} else {
		// Success. It was indeed an archive file. Extract
		// _go_.o to a byte slice.
		verb(1, "assuming %s is archive file", filename)
		cmd := exec.Command("ar", "p", filename, "_go_.o")
		var ocontents []byte
		if ocontents, err = cmd.Output(); err != nil {
			return fmt.Errorf("problems extracting _go_.o from %s: %v", filename, err)
		}
		// Emit the byte slice to a temporary file.
		tmpfile, err := ioutil.TempFile("", "objfile.o")
		if err != nil {
			return fmt.Errorf("can't open tempfile: %v", err)
		}
		verb(1, "emitting object contents into tempfile %s", tmpfile.Name())
		defer os.Remove(tmpfile.Name()) // clean up
		if _, err := tmpfile.Write(ocontents); err != nil {
			tmpfile.Close()
			return fmt.Errorf("can't write to %s: %v", tmpfile.Name(), err)
		}
		if err := tmpfile.Close(); err != nil {
			return fmt.Errorf("can't close tempfile %s: %v", tmpfile.Name(), err)
		}
		objectfile = tmpfile.Name()
	}

	// OK, now process file in question as ELF.
	if err := examineElf(filename, objectfile); err != nil {
		return err
	}
	verb(1, "done with %s", filename)
	return nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("gccgo-export-dumper: ")
	flag.Parse()
	verb(1, "in main")
	if flag.NArg() == 0 {
		usage("supply one or more objects/archives as arguments.")
	}
	for i := 0; i < flag.NArg(); i++ {
		arg := flag.Arg(i)
		infile, err := os.Open(arg)
		if err != nil {
			log.Fatal(err)
		} else {
			infile.Close()
		}
		err = examine(arg)
		if err != nil {
			log.Fatal(err)
		}
	}
	verb(1, "leaving main")
}
