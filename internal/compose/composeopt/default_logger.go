/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package composeopt

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/isolateminds/go-conduit-cli/internal/compose/types"
)

// LogConsumer consume logs from services and format them
type logConsumer struct {
	ctx        context.Context
	presenters sync.Map // map[string]*presenter
	width      int
	stdout     io.Writer
	stderr     io.Writer
	color      bool
	prefix     bool
	timestamp  bool
}

// NewLogConsumer creates a new LogConsumer
func DefaultComposeLogConsumer(ctx context.Context) SetComposerOptions {
	return func(opt *types.ComposerOptions) error {
		opt.LogConsumer = &logConsumer{
			ctx:        ctx,
			presenters: sync.Map{},
			width:      0,
			stdout:     os.Stdout,
			stderr:     os.Stderr,
			color:      true,
			prefix:     true,
			timestamp:  false,
		}
		return nil
	}
}

func (l *logConsumer) Register(name string) {
	l.register(name)
}

func (l *logConsumer) register(name string) *presenter {
	cf := monochrome
	if l.color {
		cf = nextColor()
	}
	p := &presenter{
		colors: cf,
		name:   name,
	}
	l.presenters.Store(name, p)
	if l.prefix {
		l.computeWidth()
		l.presenters.Range(func(key, value interface{}) bool {
			p := value.(*presenter)
			p.setPrefix(l.width)
			return true
		})
	}
	return p
}

func (l *logConsumer) getPresenter(container string) *presenter {
	p, ok := l.presenters.Load(container)
	if !ok { // should have been registered, but ¯\_(ツ)_/¯
		return l.register(container)
	}
	return p.(*presenter)
}

// Log formats a log message as received from name/container
func (l *logConsumer) Log(container, message string) {
	l.write(l.stdout, container, message)
}

// Log formats a log message as received from name/container
func (l *logConsumer) Err(container, message string) {
	l.write(l.stderr, container, message)
}

func (l *logConsumer) write(w io.Writer, container, message string) {
	if l.ctx.Err() != nil {
		return
	}
	p := l.getPresenter(container)
	timestamp := time.Now().Format(jsonmessage.RFC3339NanoFixed)
	for _, line := range strings.Split(message, "\n") {
		if l.timestamp {
			fmt.Fprintf(w, "%s%s%s\n", p.prefix, timestamp, line)
		} else {
			fmt.Fprintf(w, "%s%s\n", p.prefix, line)
		}
	}
}

func (l *logConsumer) Status(container, msg string) {
	p := l.getPresenter(container)
	s := p.colors(fmt.Sprintf("%s %s\n", container, msg))
	l.stdout.Write([]byte(s)) //nolint:errcheck
}

func (l *logConsumer) computeWidth() {
	width := 0
	l.presenters.Range(func(key, value interface{}) bool {
		p := value.(*presenter)
		if len(p.name) > width {
			width = len(p.name)
		}
		return true
	})
	l.width = width + 1
}

type presenter struct {
	colors colorFunc
	name   string
	prefix string
}

func (p *presenter) setPrefix(width int) {
	p.prefix = p.colors(fmt.Sprintf("%-"+strconv.Itoa(width)+"s | ", p.name))
}

var names = []string{
	"grey",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",
}

const (
	// Never use ANSI codes
	Never = "never"

	// Always use ANSI codes
	Always = "always"

	// Auto detect terminal is a tty and can use ANSI codes
	Auto = "auto"
)

// SetANSIMode configure formatter for colored output on ANSI-compliant console
func SetANSIMode(streams api.Streams, ansi string) {
	if !useAnsi(streams, ansi) {
		nextColor = func() colorFunc {
			return monochrome
		}
	}
}

func useAnsi(streams api.Streams, ansi string) bool {
	switch ansi {
	case Always:
		return true
	case Auto:
		return streams.Out().IsTerminal()
	}
	return false
}

// colorFunc use ANSI codes to render colored text on console
type colorFunc func(s string) string

var monochrome = func(s string) string {
	return s
}

func ansiColor(code, s string) string {
	return fmt.Sprintf("%s%s%s", ansi(code), s, ansi("0"))
}

func ansi(code string) string {
	return fmt.Sprintf("\033[%sm", code)
}

func makeColorFunc(code string) colorFunc {
	return func(s string) string {
		return ansiColor(code, s)
	}
}

var nextColor = rainbowColor
var rainbow []colorFunc
var currentIndex = 0
var mutex sync.Mutex

func rainbowColor() colorFunc {
	mutex.Lock()
	defer mutex.Unlock()
	result := rainbow[currentIndex]
	currentIndex = (currentIndex + 1) % len(rainbow)
	return result
}
