package pdf

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/baibeicha/fflow/pkg/files"

	"github.com/go-pdf/fpdf"
)

// WidthCache caches string width measurements for the current font
type WidthCache struct {
	cache map[string]float64
	pdf   *fpdf.Fpdf
}

// NewWidthCache creates a new width cache bound to the given PDF instance
func NewWidthCache(pdf *fpdf.Fpdf) *WidthCache {
	return &WidthCache{
		cache: make(map[string]float64),
		pdf:   pdf,
	}
}

// Width returns the width of the string in the current font, using the cache when available
func (c *WidthCache) Width(s string) float64 {
	if w, ok := c.cache[s]; ok {
		return w
	}
	w := c.pdf.GetStringWidth(s)
	c.cache[s] = w
	return w
}

var bufPool = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

// GeneratePDF creates a PDF based on the passed configuration
func GeneratePDF(filesList []files.FileInfo, cfg *Config, onProgress func()) error {
	validateConfig(cfg)
	pdf := fpdf.New(string(cfg.Orientation), "mm", string(cfg.PageSize), " ")
	pdf.SetMargins(cfg.Margin, cfg.Margin, cfg.Margin)

	if err := setupFonts(pdf, cfg); err != nil {
		return fmt.Errorf("error setting font: %w", err)
	}
	pdf.SetFont(cfg.FontName, "", cfg.FontSize)

	w, h := pdf.GetPageSize()
	pageWidth := w - (cfg.Margin * 2)
	pageHeight := h - (cfg.Margin * 2)

	cache := NewWidthCache(pdf)

	type fileData struct {
		info     files.FileInfo
		category FileCategory
		content  []byte
		err      error
	}

	const maxParallelSize = 10 * 1024 * 1024

	workers := runtime.NumCPU()
	if len(filesList) > 0 && workers > len(filesList) {
		workers = len(filesList)
	}
	if workers < 1 {
		workers = 1
	}

	data := make([]fileData, len(filesList))

	var wg sync.WaitGroup
	jobs := make(chan int, len(filesList))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				fi := filesList[idx]
				cat := GetCategory(fi.Name)

				if fi.Size > maxParallelSize || cat == CategoryImage || cat == CategoryTable {
					data[idx] = fileData{info: fi, category: cat}
					continue
				}

				content, err := os.ReadFile(fi.Path)
				data[idx] = fileData{
					info:     fi,
					category: cat,
					content:  content,
					err:      err,
				}
			}
		}()
	}

	for i := range filesList {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	for _, fd := range data {
		pdf.AddPage()

		switch fd.category {
		case CategoryImage:
			renderImage(pdf, fd.info.Path, pageWidth, pageHeight, cfg)

		case CategoryText:
			ext := strings.ToLower(filepath.Ext(fd.info.Path))
			if ext == ".md" || cfg.Text.RenderMarkdownForAll {
				var err error
				if fd.content != nil {
					err = renderMDTextReader(pdf, bytes.NewReader(fd.content), pageWidth, cfg, cache)
				} else {
					err = renderMDTextFile(pdf, fd.info.Path, pageWidth, cfg, cache)
				}
				if err != nil {
					renderText(pdf, fmt.Sprintf("Error reading markdown: %v", err))
				}
			} else {
				var err error
				if fd.content != nil {
					err = renderTextReader(pdf, bytes.NewReader(fd.content))
				} else {
					err = renderTextFile(pdf, fd.info.Path)
				}
				if err != nil {
					renderText(pdf, fmt.Sprintf("Error reading text: %v", err))
				}
			}

		case CategoryCode:
			var err error
			if fd.content != nil {
				err = renderCodeReader(pdf, bytes.NewReader(fd.content), fd.info.Path, pageWidth, cfg)
			} else {
				err = renderCodeFile(pdf, fd.info.Path, pageWidth, cfg)
			}
			if err != nil {
				renderText(pdf, fmt.Sprintf("Error reading code: %v", err))
			}

		case CategoryTable:
			if err := renderTable(pdf, fd.info.Path, pageWidth, cfg, cache); err != nil {
				renderText(pdf, fmt.Sprintf("Error reading table: %v", err))
			}

		default:
			renderText(pdf, fmt.Sprintf("File: %s\nSize: %d bytes\n(Not supported for a direct rendering)",
				fd.info.Name, fd.info.Size))
		}

		if onProgress != nil {
			onProgress()
		}
	}

	return pdf.OutputFileAndClose(cfg.OutputPath)
}

func validateConfig(cfg *Config) {
	if cfg == nil {
		return
	}
	if cfg.Orientation == "" {
		cfg.Orientation = Portrait
	}
	if cfg.PageSize == "" {
		cfg.PageSize = SizeA4
	}
	if cfg.Margin == 0 {
		cfg.Margin = 10.0
	}
	if cfg.FontName == "" {
		cfg.FontName = "Arial"
	}
	if cfg.FontSize == 0 {
		cfg.FontSize = 10.0
	}
	if cfg.Code.FontName == "" {
		cfg.Code.FontName = "Courier"
	}
	if cfg.Code.FontSize == 0 {
		cfg.Code.FontSize = 9.0
	}

	if cfg.Code.StyleAsBlock && cfg.Code.BgColor == [3]int{0, 0, 0} && cfg.Code.TextColor == [3]int{0, 0, 0} {
		cfg.Code.BgColor = [3]int{245, 246, 248}
		cfg.Code.TextColor = [3]int{36, 41, 47}
	}
}

func renderImage(pdf *fpdf.Fpdf, path string, pageWidth, pageHeight float64, cfg *Config) {
	imgType := strings.ToUpper(strings.TrimPrefix(filepath.Ext(path), "."))
	if imgType == "JPG" {
		imgType = "JPEG"
	}
	opt := fpdf.ImageOptions{ImageType: imgType, ReadDpi: true}
	info := pdf.RegisterImageOptions(path, opt)

	fileName := filepath.Base(path)

	renderCaption := func() {
		if !cfg.Image.EnableCaption {
			return
		}
		captionText := cfg.Image.CaptionPrefix + fileName
		pdf.CellFormat(pageWidth, 5, captionText, "", 1,
			string(cfg.Image.CaptionAlignment), false, 0, "")
	}

	availableHeight := pageHeight
	if cfg.Image.EnableCaption {
		availableHeight -= 5.0 + cfg.Image.CaptionMargin
	}

	if cfg.Image.EnableCaption && cfg.Image.CaptionPosition == CaptionTop {
		renderCaption()
		pdf.SetY(pdf.GetY() + cfg.Image.CaptionMargin)
	}

	x, y := pdf.GetXY()
	targetW, targetH := pageWidth, availableHeight

	if info != nil {
		imgW, imgH := info.Width(), info.Height()
		ratio := imgW / imgH

		targetW = pageWidth
		targetH = targetW / ratio

		if targetH > availableHeight {
			targetH = availableHeight
			targetW = targetH * ratio
		}
	}

	xOffset := x + (pageWidth-targetW)/2

	pdf.ImageOptions(path, xOffset, y, targetW, targetH, false, opt, 0, "")

	pdf.SetY(y + targetH)

	if cfg.Image.EnableCaption && cfg.Image.CaptionPosition == CaptionBottom {
		pdf.SetY(pdf.GetY() + cfg.Image.CaptionMargin)
		renderCaption()
	}
}

func renderTextFile(pdf *fpdf.Fpdf, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderTextReader(pdf, f)
}

func renderTextReader(pdf *fpdf.Fpdf, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	for scanner.Scan() {
		text := StripEmojis(scanner.Text())
		pdf.MultiCell(0, 5, text, "", "", false)
	}
	return scanner.Err()
}

func renderText(pdf *fpdf.Fpdf, text string) {
	pdf.MultiCell(0, 5, text, "", "", false)
}

func seekPastBOM(f *os.File) {
	_, err := f.Seek(0, 0)
	if err != nil {
		return
	}
	buf := make([]byte, 3)
	n, _ := f.Read(buf)
	if n == 3 && buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		return
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return
	}
}

func detectSeparator(path string) (rune, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	seekPastBOM(f)

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, err
		}
		return ',', nil
	}
	firstLine := scanner.Text()

	semicolons := strings.Count(firstLine, ";")
	commas := strings.Count(firstLine, ",")
	tabs := strings.Count(firstLine, "\t")

	if semicolons > commas && semicolons > tabs {
		return ';', nil
	}
	if tabs > commas && tabs > semicolons {
		return '\t', nil
	}
	return ',', nil
}

func computeColumnWidths(path string, sep rune, pageWidth float64, cfg *Config, cache *WidthCache) ([]float64, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()
	seekPastBOM(f)

	reader := csv.NewReader(f)
	reader.Comma = sep
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.ReuseRecord = true

	var colWidths []float64
	textPadding := 1.0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, err
		}

		if len(record) > len(colWidths) {
			colWidths = append(colWidths, make([]float64, len(record)-len(colWidths))...)
		}
		for i, col := range record {
			w := cache.Width(col) + (textPadding * 2)
			if w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	if len(colWidths) == 0 {
		return nil, 0, nil
	}

	var totalWidth float64
	for _, w := range colWidths {
		totalWidth += w
	}

	if totalWidth > pageWidth || cfg.Table.StretchTableWidth {
		scale := pageWidth / totalWidth
		for i := range colWidths {
			colWidths[i] *= scale
		}
	}

	return colWidths, len(colWidths), nil
}

func renderTable(pdf *fpdf.Fpdf, path string, pageWidth float64, cfg *Config, cache *WidthCache) error {
	sep := cfg.Table.Separator
	if sep == 0 {
		var err error
		sep, err = detectSeparator(path)
		if err != nil {
			return err
		}
	}

	colWidths, colsCount, err := computeColumnWidths(path, sep, pageWidth, cfg, cache)
	if err != nil {
		return err
	}
	if colsCount == 0 {
		return nil
	}

	renderCaption := func() {
		if !cfg.Table.EnableCaption {
			return
		}
		captionText := cfg.Table.CaptionPrefix + filepath.Base(path)
		pdf.CellFormat(pageWidth, 5, captionText, "", 1, string(cfg.Table.CaptionAlignment), false, 0, "")
	}

	if cfg.Table.EnableCaption && cfg.Table.CaptionPosition == CaptionTop {
		renderCaption()
		pdf.SetY(pdf.GetY() + cfg.Table.CaptionMargin)
	}

	lineHeight := 5.0
	_, pageHeight := pdf.GetPageSize()
	textPadding := 1.0

	calcRowHeight := func(row []string) float64 {
		maxLines := 1
		for i, col := range row {
			if i < colsCount {
				w := colWidths[i] - (textPadding * 2)
				if w <= 0 {
					w = 1.0
				}

				lines := pdf.SplitLines([]byte(col), w)
				if len(lines) > maxLines {
					maxLines = len(lines)
				}
			}
		}
		return float64(maxLines) * lineHeight
	}

	drawRow := func(row []string, h float64) {
		startX, startY := pdf.GetXY()
		for i, col := range row {
			if i < colsCount {
				w := colWidths[i]

				pdf.Rect(startX, startY, w, h, "D")

				pdf.SetXY(startX+textPadding, startY+textPadding)

				pdf.MultiCell(w-(textPadding*2), lineHeight, col, "", "L", false)

				startX += w
			}
		}

		pdf.SetXY(cfg.Margin, startY+h)
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	seekPastBOM(f)

	reader := csv.NewReader(f)
	reader.Comma = sep
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.ReuseRecord = true

	var headers []string
	firstRow := true
	rowIndex := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		for i, r := range record {
			record[i] = StripEmojis(r)
		}

		if firstRow {
			headers = make([]string, len(record))
			copy(headers, record)
			firstRow = false
		}

		rowH := calcRowHeight(record) + (textPadding * 2)

		if pdf.GetY()+rowH > pageHeight-cfg.Margin {
			pdf.AddPage()
			if rowIndex > 0 {
				headerH := calcRowHeight(headers) + (textPadding * 2)
				drawRow(headers, headerH)
			}
		}

		drawRow(record, rowH)
		rowIndex++
	}

	if cfg.Table.EnableCaption && cfg.Table.CaptionPosition == CaptionBottom {
		pdf.SetY(pdf.GetY() + cfg.Table.CaptionMargin)
		renderCaption()
	}

	pdf.Ln(5)
	return nil
}

var languageMap = map[string]string{
	".go": "Go", ".mod": "Go Module", ".js": "JavaScript", ".ts": "TypeScript",
	".py": "Python", ".rb": "Ruby", ".php": "PHP", ".java": "Java",
	".c": "C", ".cpp": "C++", ".cs": "C#", ".rs": "Rust", ".swift": "Swift",
	".kt": "Kotlin", ".dart": "Dart", ".lua": "Lua", ".pl": "Perl",
	".sh": "Bash", ".bash": "Bash", ".zsh": "Zsh", ".bat": "Batch",
	".html": "HTML", ".css": "CSS", ".scss": "SCSS", ".sql": "SQL",
	".json": "JSON", ".yaml": "YAML", ".yml": "YAML", ".xml": "XML",
	".md": "Markdown", "dockerfile": "Dockerfile", "makefile": "Makefile",
}

func getLanguageName(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	base := strings.ToLower(filepath.Base(path))
	if ext == "" {
		if name, exists := languageMap[base]; exists {
			return name
		}
		return strings.ToUpper(base)
	}

	if name, exists := languageMap[ext]; exists {
		return name
	}

	return strings.ToUpper(strings.TrimPrefix(ext, "."))
}

func renderCodeFile(pdf *fpdf.Fpdf, path string, pageWidth float64, cfg *Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderCodeReader(pdf, f, path, pageWidth, cfg)
}

func renderCodeReader(pdf *fpdf.Fpdf, r io.Reader, path string, pageWidth float64, cfg *Config) error {
	fill := false
	if cfg.Code.StyleAsBlock {
		fill = true
		pdf.SetFillColor(cfg.Code.BgColor[0], cfg.Code.BgColor[1], cfg.Code.BgColor[2])
		pdf.Ln(2)
	}

	startX, startY := pdf.GetXY()

	if !cfg.Code.DisableHeader {
		pdf.SetFont(cfg.FontName, "B", 8)
		pdf.SetTextColor(120, 120, 120)

		fileName := filepath.Base(path)
		langName := getLanguageName(path)

		pdf.SetXY(startX, startY)
		pdf.CellFormat(pageWidth/2, 6, "  "+fileName, "", 0, "L", fill, 0, "")

		pdf.SetXY(startX+pageWidth/2, startY)
		pdf.CellFormat(pageWidth/2, 6, langName+"  ", "", 1, "R", fill, 0, "")

		lineY := pdf.GetY()
		if fill {
			r := cfg.Code.BgColor[0] - 15
			g := cfg.Code.BgColor[1] - 15
			b := cfg.Code.BgColor[2] - 15

			if r < 0 {
				r = 0
			}
			if g < 0 {
				g = 0
			}
			if b < 0 {
				b = 0
			}
			pdf.SetDrawColor(r, g, b)
		} else {
			pdf.SetDrawColor(200, 200, 200)
		}

		pdf.SetLineWidth(0.3)
		pdf.Line(startX, lineY, startX+pageWidth, lineY)

		pdf.SetXY(startX, lineY)
		pdf.CellFormat(pageWidth, 2, " ", "", 1, "L", fill, 0, "")
	}

	pdf.SetFont(cfg.Code.FontName, "", cfg.Code.FontSize)
	pdf.SetTextColor(cfg.Code.TextColor[0], cfg.Code.TextColor[1], cfg.Code.TextColor[2])

	lineHeight := cfg.Code.FontSize*0.5 + 1.0

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "\t", "    ")
		line = StripEmojis(line)
		pdf.MultiCell(pageWidth, lineHeight, line, "", "L", fill)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont(cfg.FontName, "", cfg.FontSize)
	pdf.Ln(5)

	return nil
}
