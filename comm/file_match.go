// A revised https://github.com/sabhiram/go-gitignore/blob/master/ignore.go

/*
	The MIT License (MIT)

	# Copyright (c) 2015 Shaba Abhiram

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.

/*
ignore is a library which returns a new ignorer object which can
test against various paths. This is particularly useful when trying
to filter files based on a .gitignore document

The rules for parsing the input file are the same as the ones listed
in the Git docs here: http://git-scm.com/docs/gitignore

The summarized version of the same has been copied here:

 1. A blank line matches no files, so it can serve as a separator
    for readability.
 2. A line starting with # serves as a comment. Put a backslash ("\")
    in front of the first hash for patterns that begin with a hash.
 3. Trailing spaces are ignored unless they are quoted with backslash ("\").
 4. An optional prefix "!" which negates the pattern; any matching file
    excluded by a previous pattern will become included again. It is not
    possible to re-include a file if a parent directory of that file is
    excluded. Git doesn’t list excluded directories for performance reasons,
    so any patterns on contained files have no effect, no matter where they
    are defined. Put a backslash ("\") in front of the first "!" for
    patterns that begin with a literal "!", for example, "\!important!.txt".
 5. If the pattern ends with a slash, it is removed for the purpose of the
    following description, but it would only find a match with a directory.
    In other words, foo/ will match a directory foo and paths underneath it,
    but will not match a regular file or a symbolic link foo (this is
    consistent with the way how pathspec works in general in Git).
 6. If the pattern does not contain a slash /, Git treats it as a shell glob
    pattern and checks for a match against the pathname relative to the
    location of the .gitignore file (relative to the toplevel of the work
    tree if not from a .gitignore file).
 7. Otherwise, Git treats the pattern as a shell glob suitable for
    consumption by fnmatch(3) with the FNM_PATHNAME flag: wildcards in the
    pattern will not match a / in the pathname. For example,
    "Documentation/*.html" matches "Documentation/git.html" but not
    "Documentation/ppc/ppc.html" or "tools/perf/Documentation/perf.html".
 8. A leading slash matches the beginning of the pathname. For example,
    "/*.c" matches "cat-file.c" but not "mozilla-sha1/sha1.c".
 9. Two consecutive asterisks ("**") in patterns matched against full
    pathname may have special meaning:
    i.   A leading "**" followed by a slash means match in all directories.
    For example, "** /foo" matches file or directory "foo" anywhere,
    the same as pattern "foo". "** /foo/bar" matches file or directory
    "bar" anywhere that is directly under directory "foo".
    ii.  A trailing "/**" matches everything inside. For example, "abc/**"
    matches all files inside directory "abc", relative to the location
    of the .gitignore file, with infinite depth.
    iii. A slash followed by two consecutive asterisks then a slash matches
    zero or more directories. For example, "a/** /b" matches "a/b",
    "a/x/b", "a/x/y/b" and so on.
    iv.  Other consecutive asterisks are considered invalid.
*/
package comm

import (
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

////////////////////////////////////////////////////////////

// FileMatchParser is an interface with `MatchesPaths`.
type FileMatchParser interface {
	MatchesPath(f string) bool
	MatchesPathHow(f string) (bool, FileMatchPattern)
}

////////////////////////////////////////////////////////////

// This function pretty much attempts to mimic the parsing rules
// listed above at the start of this file
func getPatternFromLine(line string) (*regexp.Regexp, bool) {
	// Trim OS-specific carriage returns.
	line = strings.TrimRight(line, "\r")

	// Strip comments [Rule 2]
	if strings.HasPrefix(line, `#`) {
		return nil, false
	}

	// Trim string [Rule 3]
	// TODO: Handle [Rule 3], when the " " is escaped with a \
	line = strings.Trim(line, " ")

	// Exit for no-ops and return nil which will prevent us from
	// appending a pattern against this line
	if line == "" {
		return nil, false
	}

	// TODO: Handle [Rule 4] which negates the match for patterns leading with "!"
	negatePattern := false
	if line[0] == '!' {
		negatePattern = true
		line = line[1:]
	}

	// Handle [Rule 2, 4], when # or ! is escaped with a \
	// Handle [Rule 4] once we tag negatePattern, strip the leading ! char
	if regexp.MustCompile(`^(\#|\!)`).MatchString(line) {
		line = line[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) && line[0] != '/' {
		line = "/" + line
	}

	// Handle escaping the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.Replace(line, "?", `\?`, -1)

	line = strings.Replace(line, magicStar, "*", -1)

	// Temporary regex
	expr := ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "(|/.*)$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern, negatePattern
}

////////////////////////////////////////////////////////////

// FileMatchPatternT encapsulates a pattern and if it is a negated pattern.
type FileMatchPatternT struct {
	Pattern *regexp.Regexp
	Negate  bool
	LineNo  int
	Line    string
}

type FileMatchPattern = *FileMatchPatternT

// FileMatchT wraps a list of match pattern.
type FileMatchT struct {
	patterns []FileMatchPattern
}

type FileMatch = *FileMatchT

// CompileMatchLines accepts a variadic set of strings, and returns a FileMatchT
// instance which converts and appends the lines in the input to regexp.Regexp
// patterns held within the FileMatchT objects "patterns" field.
func CompileMatchLines(parent FileMatch, lines ...string) FileMatch {
	var fm FileMatch

	if parent == nil {
		fm = &FileMatchT{patterns: []FileMatchPattern{}}
	} else {
		parent_size := len(parent.patterns)
		fm = &FileMatchT{patterns: make([]FileMatchPattern, parent_size, parent_size*2)}
		copy(fm.patterns, parent.patterns)
	}

	for i, line := range lines {
		pattern, negatePattern := getPatternFromLine(line)
		if pattern != nil {
			// LineNo is 1-based numbering to match `git check-ignore -v` output
			ip := &FileMatchPatternT{pattern, negatePattern, i + 1, line}
			fm.patterns = append(fm.patterns, ip)
		}
	}
	return fm
}

// CompileMatchFile uses an match file as the input, parses the lines out of
// the file and invokes the CompileMatchLines method.
func CompileMatchFile(fs afero.Fs, parent FileMatch, fpath string) FileMatch {
	if !FileExistsP(fs, fpath) {
		return parent
	}
	lines := ReadFileLinesP(fs, fpath)
	return CompileMatchLines(parent, lines...)
}

func CompileGitIgnoreFile(fs afero.Fs, parent FileMatch, dirPath string) FileMatch {
	fpath := path.Join(dirPath, ".gitignore")
	return CompileMatchFile(fs, parent, fpath)
}

// CompileMatchFileAndLines accepts a match file as the input, parses the
// lines out of the file and invokes the CompileMatchLines method with
// additional lines.
func CompileMatchFileAndLines(fs afero.Fs, fpath string, lines ...string) FileMatch {
	if !FileExistsP(fs, fpath) {
		return CompileMatchLines(nil, lines...)
	}
	newLines := ReadFileLinesP(fs, fpath)
	return CompileMatchLines(nil, append(newLines, lines...)...)
}

func CompileGitIgnoreFileAndLines(fs afero.Fs, dirPath string, lines ...string) FileMatch {
	fpath := path.Join(dirPath, ".gitignore")
	return CompileMatchFileAndLines(fs, fpath, lines...)
}

////////////////////////////////////////////////////////////

// MatchesPath returns true if the given FileMatchT structure would target
// a given path string `f`.
func (me FileMatch) MatchesPath(f string) bool {
	matchesPath, _ := me.MatchesPathHow(f)
	return matchesPath
}

// MatchesPathHow returns true, `pattern` if the given FileMatchT structure would target
// a given path string `f`.
// The FileMatchPatternT has the Line, LineNo fields.
func (me FileMatch) MatchesPathHow(f string) (bool, FileMatchPattern) {
	// Replace OS-specific path separator.
	f = strings.Replace(f, string(os.PathSeparator), "/", -1)

	matchesPath := false
	var mip FileMatchPattern
	for _, ip := range me.patterns {
		if ip.Pattern.MatchString(f) {
			// If this is a regular target (not negated with a gitignore
			// exclude "!" etc)
			if !ip.Negate {
				matchesPath = true
				mip = ip
			} else if matchesPath {
				// Negated pattern, and matchesPath is already set
				matchesPath = false
			}
		}
	}
	return matchesPath, mip
}

////////////////////////////////////////////////////////////
