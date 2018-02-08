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
	"github.com/zserge/webview"
	"log"
)

type Controller struct {
	Label string `json:"label"`
	Progress int `json:"progress"`
}

const windowTitle = "RuneLite Launcher /*$mvn.project.version$*/"
const windowWidth  = 600
const windowHeight  = 600

var (
	window webview.WebView
	controller Controller
	refresh func()
)

func CloseWindow() {
	window.Terminate()
	window = nil
}

func UpdateProgress(value float64) {
	window.Dispatch(func() {
		controller.Progress = int(value)
		refresh()
	})
}

func AppLog(format string, a ...interface{}) {
	formatted := fmt.Sprintf(format, a...)
	log.Print(formatted)

	if window != nil {
		window.Dispatch(func() {
			if controller.Label == "" {
				controller.Label = formatted
			} else {
				controller.Label += "\n" + formatted
			}

			refresh()
		})
	}
}

func CreateUI(boot func()) {
	controller = Controller{}
	window = webview.New(webview.Settings{
		Title: windowTitle,
		Width: windowWidth,
		Height: windowHeight,
		Debug: true,
	})

	window.Dispatch(func() {
		// Bind controller
		refresh, _ = window.Bind("controller", &controller)

		// Load bootstrap css
		//#local bootstrap=str2java(evalfile("./../../resources/vendor/bootstrap.min.css"),false)
		window.InjectCSS( "/*$bootstrap$*/")

		// Load picodom library
		//#local picodom=str2java(evalfile("./../../resources/vendor/picodom.min.js"),false)
		window.Eval( "/*$picodom$*/")

		// Load main application
		//#local app=str2java(evalfile("./../../resources/app.js"),false)
		window.Eval( "/*$app$*/")

		go boot()
	})

	window.Run()
}
