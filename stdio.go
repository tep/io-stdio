// Copyright © 2017 Tim Peoples <coders@toolman.org>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
