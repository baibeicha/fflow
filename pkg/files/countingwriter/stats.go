package countingwriter

// Stats contains collected stats
type Stats struct {
	Bytes           int64
	Lines           int64
	Words           int64
	Characters      int64
	CharactersNoSpc int64
	EmptyLines      int64
	NonEmptyLines   int64
	MaxLineLength   int
	MinLineLength   int
	CurrentLineLen  int
}

func (s *Stats) Reset() {
	*s = Stats{
		MinLineLength: -1,
	}
}

func (s *Stats) Add(other *Stats) {
	s.Bytes += other.Bytes
	s.Lines += other.Lines
	s.Words += other.Words
	s.Characters += other.Characters
	s.CharactersNoSpc += other.CharactersNoSpc
	s.EmptyLines += other.EmptyLines
	s.NonEmptyLines += other.NonEmptyLines

	if other.MaxLineLength > s.MaxLineLength {
		s.MaxLineLength = other.MaxLineLength
	}

	if other.MinLineLength >= 0 {
		if s.MinLineLength < 0 || other.MinLineLength < s.MinLineLength {
			s.MinLineLength = other.MinLineLength
		}
	}
}
