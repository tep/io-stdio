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

package stdio

import (
	"fmt"
	"io"
	"os"

	"toolman.org/base/toolman"
)

var (
	quietValue   = false
	verboseValue = false
	defaultStdio = &Stdio{&verboseValue, &quietValue, os.Stdout, os.Stderr}
)

type Stdio struct {
	verbose *bool
	quiet   *bool
	stdout  io.Writer
	stderr  io.Writer
}

func New() *Stdio {
	var v, q bool
	return &Stdio{&v, &q, os.Stdout, os.Stderr}
}

func (s *Stdio) Verbose(v *bool) *bool {
	old := s.verbose
	if v != nil {
		s.verbose = v
	}
	return old
}

func (s *Stdio) Quiet(q *bool) *bool {
	old := s.quiet
	if q != nil {
		s.quiet = q
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
	s.stdout = os.Stdout
	s.stderr = os.Stderr
}

func (s *Stdio) emit(v bool, w io.Writer, args ...interface{}) {
	if *s.quiet || v && !*s.verbose {
		return
	}

	fmt.Fprintln(w, args...)
}

func (s *Stdio) Abort(err error) {
	if *s.quiet {
		err = nil
	}
	toolman.Abort(err)
}

func (s *Stdio) Die(msg string, args ...interface{}) {
	s.Abort(fmt.Errorf(msg, args...))
}

func (s *Stdio) Warn(args ...interface{})    { s.emit(false, s.stderr, args...) }
func (s *Stdio) Babble(args ...interface{})  { s.emit(true, s.stdout, args...) }
func (s *Stdio) Caution(args ...interface{}) { s.emit(true, s.stderr, args...) }
func (s *Stdio) Mention(args ...interface{}) { s.emit(false, s.stdout, args...) }

func (s *Stdio) Warnf(msg string, args ...interface{})    { s.Warn(fmt.Sprintf(msg, args...)) }
func (s *Stdio) Babblef(msg string, args ...interface{})  { s.Babble(fmt.Sprintf(msg, args...)) }
func (s *Stdio) Cautionf(msg string, args ...interface{}) { s.Caution(fmt.Sprintf(msg, args...)) }
func (s *Stdio) Mentionf(msg string, args ...interface{}) { s.Mention(fmt.Sprintf(msg, args...)) }

func Abort(err error)                     { defaultStdio.Abort(err) }
func Die(msg string, args ...interface{}) { defaultStdio.Die(msg, args...) }

func Warn(args ...interface{})    { defaultStdio.Warn(args...) }
func Babble(args ...interface{})  { defaultStdio.Babble(args...) }
func Caution(args ...interface{}) { defaultStdio.Caution(args...) }
func Mention(args ...interface{}) { defaultStdio.Mention(args...) }

func Warnf(msg string, args ...interface{})    { defaultStdio.Warnf(msg, args...) }
func Babblef(msg string, args ...interface{})  { defaultStdio.Babblef(msg, args...) }
func Cautionf(msg string, args ...interface{}) { defaultStdio.Cautionf(msg, args...) }
func Mentionf(msg string, args ...interface{}) { defaultStdio.Mentionf(msg, args...) }
