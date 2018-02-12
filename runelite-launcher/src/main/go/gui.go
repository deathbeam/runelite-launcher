/*
 * Copyright (c) 2018, Tomas Slusny <slusnucky@gmail.com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
package main

import (
	"fmt"
	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"log"
	"time"
)

func CreateUI(boot func()) {
	const title  = "/*$mvn.project.name$*/"
	const theme = nstyle.DarkTheme
	const scaling = 1

	var open bool
	var lines []string
	var curProgress int

	window := nucular.NewMasterWindow(0, title, func(window *nucular.Window) {
		window.Row(24).Dynamic(1)
		window.Progress(&curProgress, 100, false)

		for _, line := range lines  {
			window.Row(24).Dynamic(1)
			window.Label(line, "LT")
		}
	})

	// Set global logger
	logger = func (format string, a ...interface{}) {
		formatted := fmt.Sprintf(format, a...)

		if !open {
			log.Printf(formatted)
			return
		}

		lines = append(lines, formatted)
		window.Changed()
	}


	// Set global progress indicator
	progress = func (value int) {
		if !open {
			log.Printf("Downloaded %s", value)
			return
		}

		curProgress = value
		window.Changed()
	}

	// Set global window closing function
	closeWindow = func() {
		time.Sleep(time.Second * 2)
		window.Close()
	}

	window.SetStyle(nstyle.FromTheme(theme, scaling))
	go boot()
	open = true
	window.Main()
	open = false
}
