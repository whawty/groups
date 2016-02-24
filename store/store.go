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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	wl     = log.New(ioutil.Discard, "[whawty.groups]\t", log.LstdFlags)
	nameRe = regexp.MustCompile("^[-_.@A-Za-z0-9]+$")
)

const (
	userExt  string = ".user"
	groupExt string = ".group"
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
