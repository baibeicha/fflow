package pdf

import (
	"embed"
)

//go:embed fonts
var embeddedFonts embed.FS

// Orientation type for page orientation
type Orientation string

const (
	Portrait  Orientation = "P"
	Landscape Orientation = "L"
)

// PageSize type for page size
type PageSize string

const (
	SizeA4 PageSize = "A4"
	SizeA3 PageSize = "A3"
	SizeA5 PageSize = "A5"
)

// CaptionPosition - position of the caption relative to the table
type CaptionPosition string

const (
	CaptionTop    CaptionPosition = "TOP"
	CaptionBottom CaptionPosition = "BOTTOM"
)

// CaptionAlignment caption alignment
type CaptionAlignment string

const (
	AlignLeft   CaptionAlignment = "L"
	AlignCenter CaptionAlignment = "C"
	AlignRight  CaptionAlignment = "R"
)

// TableConfig stores settings for table rendering
type TableConfig struct {
	Separator         rune
	EnableCaption     bool
	CaptionPosition   CaptionPosition
	CaptionAlignment  CaptionAlignment
	CaptionMargin     float64
	CaptionPrefix     string
	StretchTableWidth bool
}

// ImageConfig stores settings for image rendering
type ImageConfig struct {
	EnableCaption    bool
	CaptionPosition  CaptionPosition
	CaptionAlignment CaptionAlignment
	CaptionMargin    float64
	CaptionPrefix    string
}

// TextConfig stores settings for text rendering
type TextConfig struct {
	RenderMarkdownForAll bool
}

// CodeConfig stores settings for source code rendering
type CodeConfig struct {
	StyleAsBlock  bool
	DisableHeader bool
	FontName      string
	FontSize      float64
	BgColor       [3]int
	TextColor     [3]int
}

// Config stores settings for PDF generation.
// If FontPath is empty, the program will try to use the built-in font

type Config struct {
	Orientation Orientation
	PageSize    PageSize
	Margin      float64
	FontName    string
	FontSize    float64
	FontPath    string
	OutputPath  string
	Table       TableConfig
	Image       ImageConfig
	Text        TextConfig
	Code        CodeConfig
}
