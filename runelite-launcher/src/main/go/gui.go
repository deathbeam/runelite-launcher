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
	"image"
	"strings"
	"time"
)

func CreateUI(boot func()) {
	const title = "/*$mvn.project.name$*/ /*$mvn.project.version$*/"
	const lineSize = 16
	const theme = style.DefaultTheme
	const scaling = 1
	const windowWidth = 640
	const windowHeight = 300
	const windowFlags = nucular.WindowBorder |
		nucular.WindowMovable |
		nucular.WindowTitle |
		nucular.WindowClosable

	var lines []string
	var curProgress int

	wordWrap := func(text string, lineWidth int) []string {
		words := strings.Fields(strings.TrimSpace(text))

		if len(words) == 0 {
			return []string{text}
		}

		wrapped := words[0]
		spaceLeft := lineWidth - len(wrapped)

		for _, word := range words[1:] {
			if len(word)+1 > spaceLeft {
				wrapped += "\n  " + word
				spaceLeft = lineWidth - len(word)
			} else {
				wrapped += " " + word
				spaceLeft -= 1 + len(word)
			}
		}

		return strings.Split(wrapped, "\n")

	}

	// Create main window with layout
	window := nucular.NewMasterWindowSize(windowFlags, title, image.Point{X: windowWidth, Y: windowHeight},
		func(window *nucular.Window) {
			window.Row(lineSize).Dynamic(1)
			window.Progress(&curProgress, 100, false)
			linesLen := len(lines)

			for i := range lines {
				line := "> " + lines[linesLen-1-i]
				wrappedLines := wordWrap(line, window.Bounds.W/7-2)
				window.Row(2).Dynamic(1)

				for _, wrappedLine := range wrappedLines {
					window.Row(lineSize - 2).Dynamic(1)
					window.Label(wrappedLine, "LT")
				}
			}
		})

	// Set GUI style to dark theme
	window.SetStyle(style.FromTheme(theme, scaling))

	// Create custom GUI logger
	guiLogger := Logger{
		LogLine: func(format string, a ...interface{}) {
			defaultLogger.LogLine(format, a...)
			formatted := fmt.Sprintf(format, a...)
			lines = append(lines, formatted)
			window.Changed()
		},
		UpdateProgress: func(value int) {
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
