// Copyright 2015 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func generic(directory string, commands []Command) {

	// execute comand
	go execute(directory, commands)
}
