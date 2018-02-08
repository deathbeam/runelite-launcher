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
	"github.com/andlabs/ui"
	"log"
  "strings"
)

const windowTitle = "RuneLite Launcher /*$mvn.project.version$*/"
const windowWidth  = 150
const windowHeight  = 80

var (
	window *ui.Window
	label  *ui.Label
	bar    *ui.ProgressBar
)

func wordWrap(text string, lineWidth int) string {
	words := strings.Fields(strings.TrimSpace(text))

	if len(words) == 0 {
		return text
	}

	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)

	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}

	return wrapped

}

func HideWindow() {
  window.Hide()
}

func CloseWindow() {
  ui.Quit()
}

func UpdateProgress(value float64) {
	ui.QueueMain(func() {
		bar.SetValue(int(value))
	})
}

func AppLog(format string, a ...interface{}) {
	formatted := fmt.Sprintf(format, a...)
	log.Print(formatted)

	ui.QueueMain(func() {
	  formattedWrapped := wordWrap(formatted, windowWidth - 20)
		label.SetText(label.Text() + formattedWrapped + "\n")
	})
}

func BuildGui(boot func()) func() {
	return func() {
		bar = ui.NewProgressBar()
		label = ui.NewLabel("")

		box := ui.NewVerticalBox()
		box.Append(label, false)
		box.Append(bar, true)
		window = ui.NewWindow(windowTitle, windowWidth, windowHeight, false)
		window.SetMargined(true)
		window.SetChild(box)
		window.Show()

		window.OnClosing(func(window *ui.Window) bool {
			ui.Quit()
			return false
		})

		go boot()
	}
}
