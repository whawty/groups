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

// Package store implements a simple storage backend for whawty.groups user data
// files. The schema of the whawty.groups store can be found in the doc directory.
// If the environment contains the variable WHAWTY_GROUPS_DEBUG logging will be enabled.
// By default whawty.groups doesn't log anything.
package store

import (
	"errors"
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

	ErrNotImplemented = errors.New("not implemented yet")
)

const (
	tmpDir        string = ".tmp"
	usersDir      string = "users"
	groupsDir     string = "groups"
	groupMetaFile string = "_meta.yaml"
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

// getTempFile provides a new, empty file in the base's .tmp directory,
//  suitable for atomic file updates (by create/write/rename)
func (d *Dir) getTempFile() (*os.File, error) {
	tmpDir := filepath.Join(d.basedir, tmpDir)
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return nil, err
	}

	return ioutil.TempFile(tmpDir, "")
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

// AddUser adds user to the store. It is an error if the user already exists.
func (d *Dir) AddUser(user string) (err error) {
	if !nameRe.MatchString(user) {
		return fmt.Errorf("user name '%s' is invalid", user)
	}
	return NewUserFile(d, user).Add()
}

// RemoveUser removes user from the store.
func (d *Dir) RemoveUser(user string) {
	NewUserFile(d, user).Remove()
	// TODO: remove user from all groups it is a member of
}

// AddGroup adds group to the store. It is an error if the group already exists.
func (d *Dir) AddGroup(group string) (err error) {
	if !nameRe.MatchString(group) {
		return fmt.Errorf("group name '%s' is invalid", group)
	}
	return NewGroupDir(d, group).Add()
}

// RemoveGroup removes group from the store.
func (d *Dir) RemoveGroup(group string) {
	NewGroupDir(d, group).Remove()
	// TODO: remove group from all groups it is a member of
}

// AddUserMember adds user to group. It is *not* an error if user is already
// a member.
func (d *Dir) AddUserMember(group, user string) error {
	u := NewUserFile(d, user)
	if exists, err := u.Exists(); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("whawty.groups.store: user '%s' does not exist", user)
	}

	return NewGroupDir(d, group).AddUserMember(user)
}

// RemoveUserMember removes user from group. It is *not* an error if user
// is not a member.
func (d *Dir) RemoveUserMember(group, user string) error {
	return NewGroupDir(d, group).RemoveUserMember(user)
}

// AddGroupMember adds groupToAdd to group. It is *not* an error if groupToAdd
// is already a member.
func (d *Dir) AddGroupMember(group, groupToAdd string) error {
	g := NewGroupDir(d, groupToAdd)
	if exists, err := g.Exists(); err != nil {
		return err
	} else if !exists {
		return fmt.Errorf("whawty.groups.store: group '%s' does not exist", groupToAdd)
	}
	// TODO: check for loops
	return NewGroupDir(d, group).AddGroupMember(groupToAdd)
}

// RemoveGroupMember removes groupToRemove from group. It is *not* an error
// if groupToRemove is not a member.
func (d *Dir) RemoveGroupMember(group, groupToRemove string) error {
	return NewGroupDir(d, group).RemoveGroupMember(groupToRemove)
}
