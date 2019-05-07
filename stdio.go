// Copyright 2017 Timothy E. Peoples
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package stdio // import "toolman.org/io/stdio"

import (
	"fmt"
	"io"
	"os"
)

var (
	quietValue   = false
	verboseValue = false
	syncValue    = false
	defaultStdio = &Stdio{&verboseValue, &quietValue, &syncValue, true, os.Stdout, os.Stderr}
)

type Stdio struct {
	verbose *bool
	quiet   *bool
	sync    *bool
	auto    bool
	stdout  io.Writer
	stderr  io.Writer
}

func New(settings ...Setting) *Stdio {
	var vv, qv, sv bool
	s := &Stdio{&vv, &qv, &sv, true, os.Stdout, os.Stderr}
	for _, sf := range settings {
		sf(s)
	}
	return s
}

type Setting func(*Stdio)

func VerboseVar(v *bool) Setting {
	return func(s *Stdio) {
		s.verbose = v
	}
}

func QuietVar(v *bool) Setting {
	return func(s *Stdio) {
		s.quiet = v
	}
}

func SyncVar(v *bool) Setting {
	return func(s *Stdio) {
		s.sync = v
	}
}

func Stdout(w io.Writer) Setting {
	return func(s *Stdio) {
		s.stdout = w
	}
}

func Stderr(w io.Writer) Setting {
	return func(s *Stdio) {
		s.stderr = w
	}
}

func AutoNL(a bool) Setting {
	return func(s *Stdio) {
		s.auto = a
	}
}

func Clone(settings ...Setting) *Stdio {
	return defaultStdio.Clone(settings...)
}

func (s *Stdio) Clone(settings ...Setting) *Stdio {
	sc := &Stdio{s.verbose, s.quiet, s.sync, s.auto, s.stdout, s.stderr}
	for _, sf := range settings {
		sf(sc)
	}
	return sc
}

func (s *Stdio) VerboseVar(v *bool) *bool {
	old := s.verbose
	if v != nil {
		s.verbose = v
	}
	return old
}

func (s *Stdio) QuietVar(v *bool) *bool {
	old := s.quiet
	if v != nil {
		s.quiet = v
	}
	return old
}

func (s *Stdio) SyncVar(v *bool) *bool {
	old := s.sync
	if v != nil {
		s.sync = v
	}
	return old
}

func (s *Stdio) Stdout(w io.Writer) io.Writer {
	old := s.stdout
	if w != nil {
		s.stdout = w
	}
	return old
}

func (s *Stdio) Stderr(w io.Writer) io.Writer {
	old := s.stderr
	if w != nil {
		s.stderr = w
	}
	return old
}

func (s *Stdio) Reset() {
	*s.verbose = false
	*s.quiet = false
	*s.sync = false
	s.stdout = os.Stdout
	s.stderr = os.Stderr
}

type syncer interface {
	Sync()
}

func (s *Stdio) emit(v bool, w io.Writer, args ...interface{}) {
	if *s.quiet || v && !*s.verbose {
		return
	}

	out := fmt.Sprintln(args...)
	if !s.auto {
		out = out[:len(out)-1]
	}

	fmt.Fprint(w, out)

	if sw, ok := w.(syncer); ok && *s.sync {
		sw.Sync()
	}
}

func (s *Stdio) Warn(args ...interface{})    { s.emit(false, s.stderr, args...) }
func (s *Stdio) Babble(args ...interface{})  { s.emit(true, s.stdout, args...) }
func (s *Stdio) Caution(args ...interface{}) { s.emit(true, s.stderr, args...) }
func (s *Stdio) Mention(args ...interface{}) { s.emit(false, s.stdout, args...) }

func (s *Stdio) Warnf(msg string, args ...interface{})    { s.Warn(fmt.Sprintf(msg, args...)) }
func (s *Stdio) Babblef(msg string, args ...interface{})  { s.Babble(fmt.Sprintf(msg, args...)) }
func (s *Stdio) Cautionf(msg string, args ...interface{}) { s.Caution(fmt.Sprintf(msg, args...)) }
func (s *Stdio) Mentionf(msg string, args ...interface{}) { s.Mention(fmt.Sprintf(msg, args...)) }

func Defaults(settings ...Setting) {
	for _, sf := range settings {
		sf(defaultStdio)
	}
}

func Warn(args ...interface{})    { defaultStdio.Warn(args...) }
func Babble(args ...interface{})  { defaultStdio.Babble(args...) }
func Caution(args ...interface{}) { defaultStdio.Caution(args...) }
func Mention(args ...interface{}) { defaultStdio.Mention(args...) }

func Warnf(msg string, args ...interface{})    { defaultStdio.Warnf(msg, args...) }
func Babblef(msg string, args ...interface{})  { defaultStdio.Babblef(msg, args...) }
func Cautionf(msg string, args ...interface{}) { defaultStdio.Cautionf(msg, args...) }
func Mentionf(msg string, args ...interface{}) { defaultStdio.Mentionf(msg, args...) }
