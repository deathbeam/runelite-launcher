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
	"github.com/aarzilli/nucular/style"
	"time"
)

func CreateUI(boot func()) {
	const title  = "/*$mvn.project.name$*/ /*$mvn.project.version$*/"
	const lineSize = 16
	const theme = style.DarkTheme
	const scaling = 1

	var lines []string
	var curProgress int

	// Create main window with layout
	window := nucular.NewMasterWindow(0, title, func(window *nucular.Window) {
		window.Row(lineSize).Dynamic(1)
		window.Progress(&curProgress, 100, false)

		for _, line := range lines  {
			window.Row(lineSize).Dynamic(1)
			window.Label(line, "LT")
		}
	})

	// Set GUI style to dark theme
	window.SetStyle(style.FromTheme(theme, scaling))

	// Create custom GUI logger
	guiLogger := Logger{
		LogLine: func (format string, a ...interface{}) {
			defaultLogger.LogLine(format, a...)
			formatted := fmt.Sprintf(format, a...)
			lines = append(lines, formatted)
			window.Changed()
		},
		UpdateProgress: func (value int) {
			defaultLogger.UpdateProgress(value)
			curProgress = value
			window.Changed()
		},
	}

	// Create main function
	main := func() {
		// Set logger to use GUI instead of console
		logger = guiLogger

		// Run provided boot function
		boot()

		// Close window after booting is done
		// Wait 2 seconds before closing the window to prevent closing before window has been fully initialized
		time.Sleep(time.Second * 2)
		window.Close()

		// Restore default logger after GUI window is closed
		logger = defaultLogger
	}

	// Run main function and start window main loop
	go main()
	window.Main()
}