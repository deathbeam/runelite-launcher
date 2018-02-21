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
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

var consoleLog = logrus.New()
var fileLog = logrus.New()

// Initialize logger from file path
func initLogger(filePath string) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)

	if err == nil {
		fileLog.Out = file
	}
}

type Logger struct {
	LogLine        func(format string, a ...interface{})
	UpdateProgress func(progress int)
}

// Create default implementation of logger
var defaultLogger = Logger{
	LogLine: func(format string, a ...interface{}) {
		consoleLog.Printf(format, a...)
		fileLog.Printf(format, a...)
	},
	UpdateProgress: func(progress int) {
		consoleLog.Printf("Progress is %d%%", progress)
		fileLog.Printf("Progress is %d%%", progress)
	},
}

var logger = defaultLogger
