//
// Copyright (c) 2016 Christian Pointner <equinox@spreadspace.org>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of whawty.groups nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

// Package store implements a simple storage backend for whawty.groups user data
// files. The schema of the whawty.groups store can be found in the doc directory.
// If the environment contains the variable WHAWTY_GROUPS_DEBUG logging will be enabled.
// By default whawty.groups doesn't log anything.
package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	wl     = log.New(ioutil.Discard, "[whawty.groups]\t", log.LstdFlags)
	nameRe = regexp.MustCompile("^[A-Za-z0-9][-_.@A-Za-z0-9]*$")
)

const (
	tmpDir    string = ".tmp"
	usersDir  string = "users"
	groupsDir string = "groups"
)

func init() {
	if _, exists := os.LookupEnv("WHAWTY_GROUPS_DEBUG"); exists {
		wl.SetOutput(os.Stderr)
	}
}

// Dir represents a directoy containing a whawty.groups store. Use NewDir to create it.
type Dir struct {
	basedir string
}

// NewDir creates a new whawty.groups store using basedir as base directory.
func NewDir(basedir string) (d *Dir) {
	d = &Dir{}
	d.basedir = filepath.Clean(basedir)
	return
}

func openDir(path string) (*os.File, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	i, err := dir.Stat()
	if err != nil {
		dir.Close()
		return nil, err
	}
	if !i.IsDir() {
		dir.Close()
		return nil, fmt.Errorf("Error: '%s' is not a directory", path)
	}

	return dir, nil
}

func isDirEmpty(dir *os.File) bool {
	if _, err := dir.Readdir(1); err == nil {
		return false
	}
	return true
}

func isDir(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("Error: %s exists but is not a directory", path)
	}
	return nil
}

// getTempFile provides a new, empty file in the base's .tmp directory,
//  suitable for atomic file updates (by create/write/rename)
func (d *Dir) getTempFile() (tmp *os.File, err error) {
	tmpDir := filepath.Join(d.basedir, tmpDir)
	err = os.MkdirAll(tmpDir, 0700)
	if err != nil {
		return nil, err
	}

	tmp, err = ioutil.TempFile(tmpDir, "")
	if err != nil {
		return tmp, err
	}

	return tmp, nil
}

// Init initializes the store by creating directories for users and groups
func (d *Dir) Init() error {
	dir, err := openDir(d.basedir)
	if err != nil {
		return err
	}
	defer dir.Close()

	if empty := isDirEmpty(dir); !empty {
		return fmt.Errorf("Error: '%s' is not empty", d.basedir)
	}

	if err = os.Mkdir(filepath.Join(d.basedir, usersDir), 0700); err != nil {
		return err
	}
	if err = os.Mkdir(filepath.Join(d.basedir, groupsDir), 0700); err != nil {
		return err
	}
	return nil
}

// Check tests if the directory is a valid whawty.group base directory.
func (d *Dir) Check() (err error) {
	dir, err := openDir(d.basedir)
	if err != nil {
		return err
	}
	defer dir.Close()

	hasUsersDir := false
	hasGroupsDir := false
	names, err := dir.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, name := range names {
		switch name {
		case tmpDir:
		case usersDir:
			hasUsersDir = true
			if err = isDir(filepath.Join(d.basedir, name)); err != nil {
				return err
			}
		case groupsDir:
			hasGroupsDir = true
			if err = isDir(filepath.Join(d.basedir, name)); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Error: found invalid file or directory: %s", name)
		}
	}

	if !hasGroupsDir {
		return fmt.Errorf("Error: groups directory not found!")
	}
	if !hasUsersDir {
		return fmt.Errorf("Error: users directory not found!")
	}

	// TODO: check usersdir and groups dir
	return nil
}
