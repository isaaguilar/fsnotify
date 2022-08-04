// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9
// +build !plan9

// Package fsnotify provides a cross-platform interface for file system
// notifications.
package fsnotify

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// These are the generalized file operations that can trigger a notification.
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// Common errors that can be reported by a watcher
var (
	ErrNonExistentWatch     = errors.New("can't remove non-existent watcher")
	ErrEventOverflow        = errors.New("fsnotify queue overflow")
	ErrNotDirectory         = errors.New("not a directory")
	ErrRecursionUnsupported = errors.New("recursion not supported")
)

// Event represents a single file system notification.
type Event struct {
	// Path to the file or directory.
	//
	// Paths are relative to the input; for example with Add("dir") the Name
	// will be set to "dir/file" if you create that file, but if you use
	// Add("/path/to/dir") it will be "/path/to/dir/file".
	Name string

	// File operation that triggered the event.
	//
	// This is a bitmask as some systems may send multiple operations at once.
	// Use the Op.Has() or Event.Has() method instead of comparing with ==.
	Op Op
}

// Op describes a set of file operations.
type Op uint32

func (op Op) String() string {
	var b strings.Builder
	if op.Has(Create) {
		b.WriteString("|CREATE")
	}
	if op.Has(Remove) {
		b.WriteString("|REMOVE")
	}
	if op.Has(Write) {
		b.WriteString("|WRITE")
	}
	if op.Has(Rename) {
		b.WriteString("|RENAME")
	}
	if op.Has(Chmod) {
		b.WriteString("|CHMOD")
	}
	if b.Len() == 0 {
		return ""
	}
	return b.String()[1:]
}

// Has reports if this operation has the given operation.
func (o Op) Has(h Op) bool { return o&h == h }

// Has reports if this event has the given operation.
func (e Event) Has(op Op) bool { return e.Op.Has(op) }

// String returns a string representation of the event in the form
// "file: REMOVE|WRITE|..."
func (e Event) String() string {
	return fmt.Sprintf("%q: %s", e.Name, e.Op.String())
}

// findDirs finds all directories under path (return value *includes* path as
// the first entry).
func findDirs(path string) ([]string, error) {
	dirs := make([]string, 0, 8)
	err := filepath.WalkDir(path, func(root string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if root == path && !d.IsDir() {
			return fmt.Errorf("%q: %w", path, ErrNotDirectory)
		}
		if d.IsDir() {
			dirs = append(dirs, root)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

// Check if this path is recursive (ends with "/..."), and return the path with
// the /... stripped.
func recursivePath(path string) (string, bool) {
	if filepath.Base(path) == "..." {
		return filepath.Dir(path), true
	}
	return path, false
}
