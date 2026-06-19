package countingwriter

// Option functional option for the CountingWriter
type Option func(*CountingWriter)

// WithBytes enables counting bytes
func WithBytes() Option {
	return func(cw *CountingWriter) {
		cw.countBytes = true
	}
}

// WithLines enables counting lines
func WithLines() Option {
	return func(cw *CountingWriter) {
		cw.countLines = true
	}
}

// WithWords enables counting words
func WithWords() Option {
	return func(cw *CountingWriter) {
		cw.countWords = true
	}
}

// WithCharacters enables counting runes
func WithCharacters() Option {
	return func(cw *CountingWriter) {
		cw.countChars = true
	}
}

// WithCharactersNoSpace enables counting runes without spaces
func WithCharactersNoSpace() Option {
	return func(cw *CountingWriter) {
		cw.countCharsNoSpc = true
	}
}

// WithEmptyLines enables counting empty and not empty lines
func WithEmptyLines() Option {
	return func(cw *CountingWriter) {
		cw.countEmptyLines = true
		// Автоматически включаем подсчет строк
		cw.countLines = true
	}
}

// WithLineLength enables counting min and max line length
func WithLineLength() Option {
	return func(cw *CountingWriter) {
		cw.countLineLength = true
		// Автоматически включаем подсчет строк
		cw.countLines = true
	}
}

// WithAll enables counting all metrics
func WithAll() Option {
	return func(cw *CountingWriter) {
		cw.countBytes = true
		cw.countLines = true
		cw.countWords = true
		cw.countChars = true
		cw.countCharsNoSpc = true
		cw.countEmptyLines = true
		cw.countLineLength = true
	}
}

// WithThreadSafe uses mutex
func WithThreadSafe() Option {
	return func(cw *CountingWriter) {
		cw.threadSafe = true
	}
}

// WithCustomWordSeparator uses custom word separator, default is unicode.IsSpace
func WithCustomWordSeparator(isSep func(r rune) bool) Option {
	return func(cw *CountingWriter) {
		cw.isWordSeparator = isSep
	}
}

// WithCustomSpace uses custom space finder, default is unicode.IsSpace
func WithCustomSpace(isSpace func(r rune) bool) Option {
	return func(cw *CountingWriter) {
		cw.isSpace = isSpace
	}
}
