package countingwriter

import (
	"bytes"
	"io"
	"sync"
	"unicode"
	"unicode/utf8"
)

// CountingWriter wraps io.Writer and calculates statistics
type CountingWriter struct {
	writer io.Writer

	countBytes      bool
	countLines      bool
	countWords      bool
	countChars      bool
	countCharsNoSpc bool
	countEmptyLines bool
	countLineLength bool
	threadSafe      bool

	isWordSeparator func(r rune) bool
	isSpace         func(r rune) bool

	stats Stats

	inWord         bool
	currentLineLen int
	lineStarted    bool

	mu sync.Mutex
}

// New creates a new CountingWriter with the specified options
// If no options are specified, the default is to count bytes only
func New(writer io.Writer, opts ...Option) *CountingWriter {
	cw := &CountingWriter{
		writer:          writer,
		isWordSeparator: unicode.IsSpace,
		isSpace:         unicode.IsSpace,
		stats: Stats{
			MinLineLength: -1,
		},
	}

	for _, opt := range opts {
		opt(cw)
	}

	if !cw.countBytes && !cw.countLines && !cw.countWords &&
		!cw.countChars && !cw.countCharsNoSpc && !cw.countEmptyLines &&
		!cw.countLineLength {
		cw.countBytes = true
	}

	return cw
}

// Write implements io.Writer
func (cw *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = cw.writer.Write(p)
	if err != nil {
		return n, err
	}

	if cw.threadSafe {
		cw.mu.Lock()
		defer cw.mu.Unlock()
	}

	cw.processBytes(p[:n])

	return n, nil
}

func (cw *CountingWriter) processBytes(p []byte) {
	if cw.countBytes {
		cw.stats.Bytes += int64(len(p))
	}

	if cw.countLines && !cw.countEmptyLines && !cw.countLineLength {
		cw.stats.Lines += int64(bytes.Count(p, []byte{'\n'}))
	}

	if !cw.countWords && !cw.countChars && !cw.countCharsNoSpc &&
		!cw.countEmptyLines && !cw.countLineLength {
		return
	}

	for len(p) > 0 {
		b := p[0]

		if b < utf8.RuneSelf {
			p = p[1:]
			cw.processASCIIByte(b)
			continue
		}

		r, size := utf8.DecodeRune(p)
		p = p[size:]
		cw.processUnicodeRune(r)
	}
}

func (cw *CountingWriter) processASCIIByte(b byte) {
	isSep := false
	isSpace := false

	if b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\v' || b == '\f' {
		isSpace = true
		isSep = true
	}

	if cw.countLines && b == '\n' && (cw.countEmptyLines || cw.countLineLength) {
		cw.stats.Lines++
		cw.finalizeLine()
		return
	}

	if cw.countWords {
		if isSep {
			cw.inWord = false
		} else if !cw.inWord {
			cw.inWord = true
			cw.stats.Words++
		}
	}

	if cw.countChars {
		cw.stats.Characters++
	}

	if cw.countCharsNoSpc && !isSpace {
		cw.stats.CharactersNoSpc++
	}

	if cw.countLineLength && !isSpace {
		cw.currentLineLen++
		cw.lineStarted = true
	}
}

func (cw *CountingWriter) processUnicodeRune(r rune) {
	isSep := cw.isWordSeparator(r)
	isSpace := cw.isSpace(r)

	if cw.countLines && r == '\n' && (cw.countEmptyLines || cw.countLineLength) {
		cw.stats.Lines++
		cw.finalizeLine()
		return
	}

	if cw.countWords {
		if isSep {
			cw.inWord = false
		} else if !cw.inWord {
			cw.inWord = true
			cw.stats.Words++
		}
	}

	if cw.countChars {
		cw.stats.Characters++
	}

	if cw.countCharsNoSpc && !isSpace {
		cw.stats.CharactersNoSpc++
	}

	if cw.countLineLength && !isSpace {
		cw.currentLineLen++
		cw.lineStarted = true
	}
}

func (cw *CountingWriter) finalizeLine() {
	if cw.countEmptyLines {
		if cw.lineStarted {
			cw.stats.NonEmptyLines++
		} else {
			cw.stats.EmptyLines++
		}
	}

	if cw.countLineLength {
		lineLen := cw.currentLineLen

		if lineLen > cw.stats.MaxLineLength {
			cw.stats.MaxLineLength = lineLen
		}

		if lineLen > 0 {
			if cw.stats.MinLineLength < 0 || lineLen < cw.stats.MinLineLength {
				cw.stats.MinLineLength = lineLen
			}
		}

		cw.currentLineLen = 0
		cw.lineStarted = false
	}
}

// Stats returns the current statistics
func (cw *CountingWriter) Stats() Stats {
	if cw.threadSafe {
		cw.mu.Lock()
		defer cw.mu.Unlock()
	}
	return cw.stats
}

// StatsPtr returns a pointer to statistics
func (cw *CountingWriter) StatsPtr() *Stats {
	if cw.threadSafe {
		cw.mu.Lock()
		defer cw.mu.Unlock()
	}
	return &cw.stats
}

// Reset resets all counters
func (cw *CountingWriter) Reset() {
	if cw.threadSafe {
		cw.mu.Lock()
		defer cw.mu.Unlock()
	}

	cw.stats.Reset()
	cw.inWord = false
	cw.currentLineLen = 0
	cw.lineStarted = false
}

// Flush calls Flush on the internal writer if it supports io.WriterTo
func (cw *CountingWriter) Flush() error {
	if flusher, ok := cw.writer.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}

// Writer returns the original writer
func (cw *CountingWriter) Writer() io.Writer {
	return cw.writer
}
