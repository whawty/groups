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
	"os"
	"path/filepath"
	"testing"
)

func TestAddRemoveGroup(t *testing.T) {
	groupname := "test-addremove-group"

	u := NewGroupDir(testStoreGroupDir, groupname)

	if err := u.Add(); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if fi, err := os.Stat(filepath.Join(testBaseDirGroupDir, groupsDir, groupname)); err != nil {
		t.Fatal("cannot find test group dir after add:", err)
	} else if !fi.IsDir() {
		t.Fatal("group dir is not a directory")
	}
	if _, err := os.Stat(filepath.Join(testBaseDirGroupDir, groupsDir, groupname, groupMetaFile)); err != nil {
		t.Fatal("cannot read test group's meta file after add:", err)
	}

	if err := u.Add(); err == nil {
		t.Fatal("adding group a second time returned no error!")
	}

	u.Remove()
	if _, err := os.Stat(filepath.Join(testBaseDirGroupDir, groupsDir, groupname)); err == nil {
		t.Fatal("test group does still exist after remove")
	} else if !os.IsNotExist(err) {
		t.Fatal("unexpected error:", err)
	}
}

func TestExistsGroup(t *testing.T) {
	groupname := "test-exists-group"

	u := NewGroupDir(testStoreGroupDir, groupname)

	if exists, err := u.Exists(); err != nil {
		t.Fatal("unexpected error:", err)
	} else if exists {
		t.Fatal("file file for test group shouldn't exist")
	}

	if err := u.Add(); err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer u.Remove()

	if exists, err := u.Exists(); err != nil {
		t.Fatal("unexpected error:", err)
	} else if !exists {
		t.Fatal("file for test group should exist")
	}
}
