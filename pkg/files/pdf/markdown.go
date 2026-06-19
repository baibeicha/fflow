package pdf

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"
)

// TextSegment stores a piece of text and its style.
type TextSegment struct {
	Text   string
	Bold   bool
	Italic bool
	Link   string
}

func parseMarkdownLine(line string) []TextSegment {
	line = StripEmojis(line)

	var segments []TextSegment
	i, n := 0, len(line)
	var buf strings.Builder
	buf.Grow(len(line))

	flush := func() {
		if buf.Len() > 0 {
			segments = append(segments, TextSegment{Text: buf.String()})
			buf.Reset()
		}
	}

	for i < n {
		if i+1 < n && line[i] == '*' && line[i+1] == '*' {
			end := strings.Index(line[i+2:], "**")
			if end != -1 {
				flush()
				segments = append(segments, TextSegment{Text: line[i+2 : i+2+end], Bold: true})
				i += 4 + end
				continue
			}
		}

		if line[i] == '*' {
			end := -1
			for j := i + 1; j < n; j++ {
				if line[j] == '*' {
					if j+1 < n && line[j+1] == '*' {
						continue
					}
					end = j
					break
				}
			}
			if end != -1 && end > i+1 {
				flush()
				segments = append(segments, TextSegment{Text: line[i+1 : end], Italic: true})
				i = end + 1
				continue
			}
		}

		if line[i] == '[' {
			closeBracket := strings.Index(line[i:], "](")
			if closeBracket != -1 {
				closeParen := strings.Index(line[i+closeBracket+2:], ")")
				if closeParen != -1 {
					flush()
					text := line[i+1 : i+closeBracket]
					url := line[i+closeBracket+2 : i+closeBracket+2+closeParen]
					segments = append(segments, TextSegment{Text: text, Link: url})
					i += closeBracket + 2 + closeParen + 1
					continue
				}
			}
		}

		buf.WriteByte(line[i])
		i++
	}

	flush()
	return segments
}

func renderMarkdownParagraph(pdf *fpdf.Fpdf, segments []TextSegment, cfg *Config, maxW float64, cache *WidthCache) {
	lineHeight := cfg.FontSize * 0.5
	startX := cfg.Margin
	if pdf.GetX() > startX {
		pdf.Ln(lineHeight)
	}

	spaceW := cache.Width(" ")

	for _, seg := range segments {
		style := ""
		if seg.Bold {
			style += "B"
		}
		if seg.Italic {
			style += "I"
		}

		isLink := seg.Link != ""
		if isLink {
			style += "U"
			pdf.SetTextColor(0, 102, 204)
		} else {
			pdf.SetTextColor(0, 0, 0)
		}

		pdf.SetFont(cfg.FontName, style, cfg.FontSize)

		words := strings.Split(seg.Text, " ")

		for i, word := range words {
			if word == "" {
				continue
			}

			w := cache.Width(word)

			if pdf.GetX()+w > startX+maxW && pdf.GetX() > startX {
				pdf.Ln(lineHeight)
				pdf.SetX(startX)
			}

			pdf.CellFormat(w, lineHeight, word, "", 0, "L", false, 0, seg.Link)

			if i < len(words)-1 {
				if pdf.GetX()+spaceW <= startX+maxW {
					if isLink {
						pdf.SetFont(cfg.FontName, strings.ReplaceAll(style, "U", ""), cfg.FontSize)
					}
					pdf.CellFormat(spaceW, lineHeight, " ", "", 0, "L", false, 0, "")

					if isLink {
						pdf.SetFont(cfg.FontName, style, cfg.FontSize)
					}
				}
			}
		}
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont(cfg.FontName, "", cfg.FontSize)
	pdf.Ln(lineHeight * 1.5)
}

func renderMDTextFile(pdf *fpdf.Fpdf, path string, pageWidth float64, cfg *Config, cache *WidthCache) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderMDTextReader(pdf, f, pageWidth, cfg, cache)
}

func renderMDTextReader(pdf *fpdf.Fpdf, r io.Reader, pageWidth float64, cfg *Config, cache *WidthCache) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			pdf.Ln(cfg.FontSize * 0.5)
			continue
		}

		switch {
		case strings.HasPrefix(line, "### "):
			text := strings.TrimPrefix(line, "### ")
			pdf.SetFont(cfg.FontName, "B", cfg.FontSize+2.0)
			pdf.MultiCell(pageWidth, (cfg.FontSize+2)*0.5, text, "", "L", false)
			pdf.Ln(2)

		case strings.HasPrefix(line, "## "):
			text := strings.TrimPrefix(line, "## ")
			pdf.SetFont(cfg.FontName, "B", cfg.FontSize+4.0)
			pdf.MultiCell(pageWidth, (cfg.FontSize+4)*0.5, text, "", "L", false)
			pdf.Ln(2)

		case strings.HasPrefix(line, "# "):
			text := strings.TrimPrefix(line, "# ")
			pdf.SetFont(cfg.FontName, "B", cfg.FontSize+6.0)
			pdf.MultiCell(pageWidth, (cfg.FontSize+6)*0.5, text, "", "L", false)
			pdf.Ln(2)

		default:
			segments := parseMarkdownLine(line)
			renderMarkdownParagraph(pdf, segments, cfg, pageWidth, cache)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	pdf.Ln(5)
	return nil
}
