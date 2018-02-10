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
	"github.com/gizak/termui"
	"log"
	"time"
)

func CreateTUI(boot func()) {
	var bar *termui.Gauge
	var open bool

	// Set global logger
	logger = func (format string, a ...interface{}) {
		formatted := fmt.Sprintf(format, a...)

		if !open {
			log.Printf(formatted)
			return
		}

		par := termui.NewPar(formatted)
		par.Border = false
		par.Height = 1

		termui.Body.AddRows(termui.NewRow(termui.NewCol(12, 0, par)))
		termui.Body.Align()
		termui.Render(termui.Body)
	}

	// Set global progress indicator
	progress = func (value int) {
		bar.Percent = value

		if !open {
			log.Printf("Downloaded %s", value)
			return
		}

		termui.Body.Align()
		termui.Render(termui.Body)
	}

	// Set global close function
	closeWindow = func() {
		// This does nothing here
	}

	err := termui.Init()
	open = true

	if err != nil {
		panic(err)
	}

	bar = termui.NewGauge()
	bar.BarColor = termui.ColorCyan
	termui.Body.AddRows(termui.NewRow(termui.NewCol(12, 0, bar)))
	termui.Body.Align()
	termui.Render(termui.Body)
	boot()
	open = false
	time.Sleep(time.Second * 2)
	termui.Close()
}
