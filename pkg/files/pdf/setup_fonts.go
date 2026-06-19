package pdf

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-pdf/fpdf"
)

func setupFonts(pdf *fpdf.Fpdf, cfg *Config) error {
	if cfg.FontPath != "" {
		pdf.AddUTF8Font(cfg.FontName, "", cfg.FontPath)

		boldPath := injectStyleSuffix(cfg.FontPath, "-Bold")
		if checkFileExists(boldPath) {
			pdf.AddUTF8Font(cfg.FontName, "B", boldPath)
		} else {
			pdf.AddUTF8Font(cfg.FontName, "B", cfg.FontPath)
		}

		italicPath := injectStyleSuffix(cfg.FontPath, "-Italic")
		if checkFileExists(italicPath) {
			pdf.AddUTF8Font(cfg.FontName, "I", italicPath)
		} else {
			pdf.AddUTF8Font(cfg.FontName, "I", cfg.FontPath)
		}

		boldItalicPath := injectStyleSuffix(cfg.FontPath, "-BoldItalic")
		if checkFileExists(boldItalicPath) {
			pdf.AddUTF8Font(cfg.FontName, "BI", boldItalicPath)
		} else {
			pdf.AddUTF8Font(cfg.FontName, "BI", cfg.FontPath)
		}

		return nil
	}

	baseFontBytes, err := embeddedFonts.ReadFile("fonts/Arial.ttf")
	if err != nil {
		return err
	}
	pdf.AddUTF8FontFromBytes(cfg.FontName, "", baseFontBytes)

	if boldBytes, err := embeddedFonts.ReadFile("fonts/Arial-Bold.ttf"); err == nil {
		pdf.AddUTF8FontFromBytes(cfg.FontName, "B", boldBytes)
	} else {
		pdf.AddUTF8FontFromBytes(cfg.FontName, "B", baseFontBytes)
	}

	if italicBytes, err := embeddedFonts.ReadFile("fonts/Arial-Italic.ttf"); err == nil {
		pdf.AddUTF8FontFromBytes(cfg.FontName, "I", italicBytes)
	} else {
		pdf.AddUTF8FontFromBytes(cfg.FontName, "I", baseFontBytes)
	}

	if boldItalicBytes, err := embeddedFonts.ReadFile("fonts/Arial-BoldItalic.ttf"); err == nil {
		pdf.AddUTF8FontFromBytes(cfg.FontName, "BI", boldItalicBytes)
	} else {
		pdf.AddUTF8FontFromBytes(cfg.FontName, "BI", baseFontBytes)
	}

	return nil
}

// "/fonts/Arial.ttf" + "-Bold" -> "/fonts/Arial-Bold.ttf"
func injectStyleSuffix(path, suffix string) string {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	return base + suffix + ext
}

func checkFileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
