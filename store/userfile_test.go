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

package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddRemoveUser(t *testing.T) {
	username := "test-addremove-user"

	u := NewUserFile(testStoreUserFile, username)

	if err := u.Add(); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if _, err := os.Stat(filepath.Join(testBaseDirUserFile, username)); err != nil {
		t.Fatal("cannot read test user file after add:", err)
	}

	if err := u.Add(); err == nil {
		t.Fatal("adding user a second time returned no error!")
	}

	u.Remove()
	if _, err := os.Stat(filepath.Join(testBaseDirUserFile, username)); err == nil {
		t.Fatal("test user does still exist after remove")
	} else if !os.IsNotExist(err) {
		t.Fatal("unexpected error:", err)
	}
}

func TestExistsUser(t *testing.T) {
	username := "test-exists-user"

	u := NewUserFile(testStoreUserFile, username)

	if exists, err := u.Exists(); err != nil {
		t.Fatal("unexpected error:", err)
	} else if exists {
		t.Fatal("file file for test user shouldn't exist")
	}

	if err := u.Add(); err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer u.Remove()

	if exists, err := u.Exists(); err != nil {
		t.Fatal("unexpected error:", err)
	} else if !exists {
		t.Fatal("file for test user should exist")
	}
}
