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
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const (
	testBaseDir         string = "test-store"
	testBaseDirUserFile string = "test-store-user"
)

var (
	testStoreUserFile *Dir
)

func TestInitDir(t *testing.T) {
	store := NewDir(testBaseDir)

	if err := store.Init(); err == nil {
		t.Fatalf("Initializing a not existing dir should give an error")
	}

	if file, err := os.Create(testBaseDir); err != nil {
		t.Fatal("unexpected error:", err)
	} else {
		file.Close()
	}

	if err := store.Init(); err == nil {
		t.Fatalf("Initializing where path is a not a dir should give an error")
	}

	if err := os.Remove(testBaseDir); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := os.Mkdir(testBaseDir, 0000); err != nil {
		t.Fatal("unexpected error:", err)
	}
	defer os.RemoveAll(testBaseDir)

	if err := store.Init(); err == nil {
		t.Fatalf("Initializing of a directory with wrong permissions shouldn't work")
	}

	if err := os.Chmod(testBaseDir, 0755); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if file, err := os.Create(filepath.Join(testBaseDir, "testfile")); err != nil {
		t.Fatal("unexpected error:", err)
	} else {
		file.Close()
	}

	if err := store.Init(); err == nil {
		t.Fatalf("Initializing a non-empty directory should give an error")
	}

	if err := os.Remove(filepath.Join(testBaseDir, "testfile")); err != nil {
		t.Fatal("unexpected error:", err)
	}

	if err := store.Init(); err != nil {
		t.Fatalf("unexpected error")
	}
}

func TestMain(m *testing.M) {
	if err := os.Mkdir(testBaseDirUserFile, 0755); err != nil {
		fmt.Println("Error creating store base directory:", err)
		os.Exit(-1)
	}

	testStoreUserFile = NewDir(testBaseDirUserFile)

	ret := m.Run()

	if err := os.RemoveAll(testBaseDirUserFile); err != nil {
		fmt.Println("Error removing store base directory:", err)
	}
	os.Exit(ret)
}
