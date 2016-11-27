//
// Copyright (c) 2016 whawty contributors (see AUTHORS file)
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

package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// fileExists returns whether the given file or directory exists or not
// this is from: stackoverflow.com/questions/10510691
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// UserFile is the representation of a single user file inside the store.
// Use NewUserFile to create it.
type UserFile struct {
	store *Dir
	user  string
}

// NewUserFile creates a new whawty.groups UserFile for user inside basedir.
func NewUserFile(store *Dir, user string) (u *UserFile) {
	u = &UserFile{}
	u.store = store
	u.user = user
	return
}

func (u *UserFile) getFilename() string {
	return filepath.Join(u.store.basedir, u.user)
}

// Add creates the user file. It is an error if the user already exists.
func (u *UserFile) Add() (err error) {
	var exists bool
	if exists, err = u.Exists(); err != nil {
		return
	} else if exists {
		return fmt.Errorf("whawty.groups.store: user '%s' already exists", u.user)
	}
	var file *os.File
	file, err = os.Create(u.getFilename())
	defer file.Close()

	m := make(map[string]interface{})
	m["changed"] = time.Now()
	var data []byte
	if data, err = yaml.Marshal(m); err != nil {
		return
	}
	file.Write(data)
	return nil
}

// Remove deletes the user file.
func (u *UserFile) Remove() {
	os.Remove(u.getFilename())
	return
}

// Exists checks if user exists.
func (u *UserFile) Exists() (exists bool, err error) {
	return fileExists(u.getFilename())
}
