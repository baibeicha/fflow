// Package presets provides ready-to-use configuration templates for PDF generation.
// Grouped by document type and localized for different languages.
// Always call the Clone() method before modifying a preset to prevent mutating the global variables.
package presets

import "github.com/baibeicha/fflow/pkg/files/pdf"

// PresetBundle groups configuration presets by localization.
type PresetBundle struct {
	EN *pdf.Config
	RU *pdf.Config
}

var (
	// A4Portrait provides a standard A4 layout for general purpose lists and text.
	A4Portrait = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      10.0,
			FontName:    "Arial",
			FontSize:    10.0,
			OutputPath:  "output.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignCenter,
				CaptionMargin:     3.0,
				CaptionPrefix:     "File: ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    3.0,
				CaptionPrefix:    "Image: ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      9.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      10.0,
			FontName:    "Arial",
			FontSize:    10.0,
			OutputPath:  "output.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignCenter,
				CaptionMargin:     3.0,
				CaptionPrefix:     "Файл: ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    3.0,
				CaptionPrefix:    "Рис: ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      9.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
	}

	// A4Landscape provides a landscape A4 layout, best suited for wide CSV tables.
	A4Landscape = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA4,
			Margin:      10.0,
			FontName:    "Arial",
			FontSize:    10.0,
			OutputPath:  "output_landscape.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     3.0,
				CaptionPrefix:     "Data: ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignLeft,
				CaptionMargin:    3.0,
				CaptionPrefix:    "Image: ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      9.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA4,
			Margin:      10.0,
			FontName:    "Arial",
			FontSize:    10.0,
			OutputPath:  "output_landscape.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     3.0,
				CaptionPrefix:     "Данные: ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignLeft,
				CaptionMargin:    3.0,
				CaptionPrefix:    "Рис: ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      9.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
	}

	// AcademicReport is configured for formal academic papers, thesis drafts, and university laboratory reports.
	AcademicReport = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      20.0,
			FontName:    "Arial",
			FontSize:    14.0,
			OutputPath:  "academic_report.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     5.0,
				CaptionPrefix:     "Table - ",
				StretchTableWidth: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    5.0,
				CaptionPrefix:    "Figure - ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  false,
				DisableHeader: true,
				FontName:      "Courier",
				FontSize:      12.0,
				BgColor:       [3]int{255, 255, 255},
				TextColor:     [3]int{0, 0, 0},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      20.0,
			FontName:    "Arial",
			FontSize:    14.0,
			OutputPath:  "academic_report.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     5.0,
				CaptionPrefix:     "Таблица - ",
				StretchTableWidth: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    5.0,
				CaptionPrefix:    "Рисунок - ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  false,
				DisableHeader: true,
				FontName:      "Courier",
				FontSize:      12.0,
				BgColor:       [3]int{255, 255, 255},
				TextColor:     [3]int{0, 0, 0},
			},
		},
	}

	// ScientificArticle is strictly tailored for publications (e.g., VAK, RINC indexing).
	ScientificArticle = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      25.0,
			FontName:    "Arial",
			FontSize:    12.0,
			OutputPath:  "article_draft.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignRight,
				CaptionMargin:     4.0,
				CaptionPrefix:     "Table ",
				StretchTableWidth: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    4.0,
				CaptionPrefix:    "Fig. ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  false,
				DisableHeader: true,
				FontName:      "Courier",
				FontSize:      10.0,
				BgColor:       [3]int{255, 255, 255},
				TextColor:     [3]int{0, 0, 0},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      25.0,
			FontName:    "Arial",
			FontSize:    14.0,
			OutputPath:  "article_draft.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignRight,
				CaptionMargin:     4.0,
				CaptionPrefix:     "Таблица ",
				StretchTableWidth: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    4.0,
				CaptionPrefix:    "Рис. ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  false,
				DisableHeader: true,
				FontName:      "Courier",
				FontSize:      12.0,
				BgColor:       [3]int{255, 255, 255},
				TextColor:     [3]int{0, 0, 0},
			},
		},
	}

	// SourceCode is optimized for printing backend codebase structures, configs, or raw scripts.
	SourceCode = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      15.0,
			FontName:    "Arial",
			FontSize:    9.0,
			OutputPath:  "source_code.pdf",
			Table: pdf.TableConfig{
				EnableCaption: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      9.0,
				BgColor:       [3]int{240, 240, 240},
				TextColor:     [3]int{0, 0, 0},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      15.0,
			FontName:    "Arial",
			FontSize:    9.0,
			OutputPath:  "source_code.pdf",
			Table: pdf.TableConfig{
				EnableCaption: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      9.0,
				BgColor:       [3]int{240, 240, 240},
				TextColor:     [3]int{0, 0, 0},
			},
		},
	}

	// ServerLogs provides a dense layout for exporting extensive server logs or trace events.
	ServerLogs = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA4,
			Margin:      5.0,
			FontName:    "Arial",
			FontSize:    7.0,
			OutputPath:  "server_logs.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     2.0,
				CaptionPrefix:     "Log: ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      7.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA4,
			Margin:      5.0,
			FontName:    "Arial",
			FontSize:    7.0,
			OutputPath:  "server_logs.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     2.0,
				CaptionPrefix:     "Лог: ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      7.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
	}

	// Cheatsheet maximizes data density on a single page, perfect for quick reference guides.
	Cheatsheet = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      3.0,
			FontName:    "Arial",
			FontSize:    6.0,
			OutputPath:  "cheatsheet.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     false,
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: true,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: true,
				FontName:      "Courier",
				FontSize:      6.0,
				BgColor:       [3]int{250, 250, 250},
				TextColor:     [3]int{0, 0, 0},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Portrait,
			PageSize:    pdf.SizeA4,
			Margin:      3.0,
			FontName:    "Arial",
			FontSize:    6.0,
			OutputPath:  "cheatsheet.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     false,
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: true,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: true,
				FontName:      "Courier",
				FontSize:      6.0,
				BgColor:       [3]int{250, 250, 250},
				TextColor:     [3]int{0, 0, 0},
			},
		},
	}

	// A3DataHeavy provides an A3 landscape layout for massive tables with numerous columns.
	A3DataHeavy = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA3,
			Margin:      15.0,
			FontName:    "Arial",
			FontSize:    8.0,
			OutputPath:  "big_data.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     2.0,
				CaptionPrefix:     "[EXPORT] ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      8.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA3,
			Margin:      15.0,
			FontName:    "Arial",
			FontSize:    8.0,
			OutputPath:  "big_data.pdf",
			Table: pdf.TableConfig{
				Separator:         0,
				EnableCaption:     true,
				CaptionPosition:   pdf.CaptionTop,
				CaptionAlignment:  pdf.AlignLeft,
				CaptionMargin:     2.0,
				CaptionPrefix:     "[ВЫГРУЗКА] ",
				StretchTableWidth: true,
			},
			Image: pdf.ImageConfig{
				EnableCaption: false,
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      8.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
	}

	// PhotoAlbum is optimized for rendering image galleries.
	PhotoAlbum = PresetBundle{
		EN: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA4,
			Margin:      5.0,
			FontName:    "Arial",
			FontSize:    12.0,
			OutputPath:  "album.pdf",
			Table: pdf.TableConfig{
				EnableCaption: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    5.0,
				CaptionPrefix:    "Photo: ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      10.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
		RU: &pdf.Config{
			Orientation: pdf.Landscape,
			PageSize:    pdf.SizeA4,
			Margin:      5.0,
			FontName:    "Arial",
			FontSize:    12.0,
			OutputPath:  "album.pdf",
			Table: pdf.TableConfig{
				EnableCaption: false,
			},
			Image: pdf.ImageConfig{
				EnableCaption:    true,
				CaptionPosition:  pdf.CaptionBottom,
				CaptionAlignment: pdf.AlignCenter,
				CaptionMargin:    5.0,
				CaptionPrefix:    "Фото: ",
			},
			Text: pdf.TextConfig{
				RenderMarkdownForAll: false,
			},
			Code: pdf.CodeConfig{
				StyleAsBlock:  true,
				DisableHeader: false,
				FontName:      "Courier",
				FontSize:      10.0,
				BgColor:       [3]int{245, 246, 248},
				TextColor:     [3]int{36, 41, 47},
			},
		},
	}
)
