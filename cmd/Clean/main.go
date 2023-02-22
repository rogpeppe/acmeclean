package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	// Use vendored version that provides extra info when available.
	"github.com/rogpeppe/acmeclean/internal/acme"
)

var allFlag = flag.Bool("a", false, "clean windows with history too")

func main() {
	flag.Parse()
	wins, err := acme.Windows()
	if err != nil {
		log.Fatal(err)
	}
	for _, info := range wins {
		if err := clean(info, *allFlag); err != nil {
			fmt.Printf("cannot clean window %q (id %d): %v]\n", info.Name, info.ID, err)
		}
	}
}

func clean(info acme.WinInfo, all bool) error {
	w, err := acme.Open(info.ID, nil)
	if err != nil {
		return err
	}
	defer w.CloseFiles()

	winfo, err := w.Info()
	if err != nil {
		return fmt.Errorf("cannot get info for window: %v", err)
	}
	if winfo.History != nil {
		info.History = winfo.History
	}
	if isClean(info, all) {
		if err := w.Ctl("del"); err != nil {
			return fmt.Errorf("cannot delete: %v", err)
		}
	}
	return nil
}

func isClean(info acme.WinInfo, all bool) bool {
	if strings.HasSuffix(info.Name, "/+Errors") || info.Name == "+Errors" {
		// +Errors windows report as modified even when the
		// tag says they're not. We probably want to delete them always.
		return true
	}
	if info.IsModified {
		return false
	}
	if all {
		return true
	}
	// Maybe don't delete if the tag contains something long?
	//if !strings.HasSuffix(info.Tag, "| Look ") {
	//	return false
	//}
	if info.History == nil {
		return true
	}
	// TODO could return true on +Errors windows re
	if info.History.CanUndo || info.History.CanRedo {
		return false
	}
	return true
}
