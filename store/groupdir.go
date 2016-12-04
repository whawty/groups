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

// GroupDir is the representation of a single group directory inside the store.
// Use NewGroupDir to create it.
type GroupDir struct {
	store *Dir
	group string
}

// NewGroupDir creates a new whawty.groups GroupDir for group inside basedir.
func NewGroupDir(store *Dir, group string) (g *GroupDir) {
	g = &GroupDir{}
	g.store = store
	g.group = group
	return
}

func (g *GroupDir) getDirname() string {
	return filepath.Join(g.store.basedir, groupsDir, g.group)
}

func (g *GroupDir) getMetafilename() string {
	return filepath.Join(g.store.basedir, groupsDir, g.group, groupMetaFile)
}

// Add creates the group directory. It is an error if the group already exists.
func (g *GroupDir) Add() (err error) {
	var exists bool
	if exists, err = g.Exists(); err != nil {
		return
	} else if exists {
		return fmt.Errorf("whawty.groups.store: group '%s' already exists", g.group)
	}

	if err = os.Mkdir(g.getDirname(), 0755); err != nil {
		return
	}

	var file *os.File
	file, err = os.Create(g.getMetafilename())
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

// Remove deletes the group directory.
func (g *GroupDir) Remove() {
	os.RemoveAll(g.getDirname())
	return
}

// Exists checks if group exists.
func (g *GroupDir) Exists() (exists bool, err error) {
	return fileExists(g.getDirname())
}

// AddUserMember adds link to user file
func (g *GroupDir) AddUserMember(user string) error {
	return ErrNotImplemented
}

// RemoveUserMember removes the link to user file
func (g *GroupDir) RemoveUserMember(user string) error {
	return ErrNotImplemented
}

// AddGroupMember adds link to group dir
func (g *GroupDir) AddGroupMember(group string) error {
	return ErrNotImplemented
}

// RemoveGroupMember removes the link to group dir
func (g *GroupDir) RemoveGroupMember(group string) error {
	return ErrNotImplemented
}
