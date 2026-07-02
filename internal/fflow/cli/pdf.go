package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/baibeicha/fflow/pkg/telemetry"
	"github.com/spf13/cobra"

	"github.com/baibeicha/fflow/internal/fflow/locale"
	"github.com/baibeicha/fflow/internal/fflow/ui"
	"github.com/baibeicha/fflow/pkg/files"
	"github.com/baibeicha/fflow/pkg/files/pdf"
	"github.com/baibeicha/fflow/pkg/files/pdf/presets"
)

var (
	pdfPaths      []string
	pdfOutput     string
	pdfPreset     string
	pdfLang       string
	pdfRecursive  bool
	pdfExtensions []string
	pdfBlacklist  []string
	pdfCategories []string
	pdfMinSize    string
	pdfMaxSize    string
	pdfSortBy     []string
	pdfLimit      string
	pdfOffset     string

	pdfOrientation string
	pdfPageSize    string
	pdfMargin      float64
	pdfFontName    string
	pdfFontSize    float64
	pdfFontPath    string

	pdfTableSeparator        string
	pdfTableEnableCaption    bool
	pdfTableCaptionPosition  string
	pdfTableCaptionAlignment string
	pdfTableCaptionMargin    float64
	pdfTableCaptionPrefix    string
	pdfTableStretchWidth     bool

	pdfImageEnableCaption    bool
	pdfImageCaptionPosition  string
	pdfImageCaptionAlignment string
	pdfImageCaptionMargin    float64
	pdfImageCaptionPrefix    string

	pdfTextRenderMarkdown bool

	pdfCodeStyleAsBlock  bool
	pdfCodeDisableHeader bool
	pdfCodeFontName      string
	pdfCodeFontSize      float64
	pdfCodeBgColor       string
	pdfCodeTextColor     string
)

var pdfCmd = &cobra.Command{
	Use:     "pdf",
	Short:   locale.T("commands.pdf.short"),
	Long:    locale.T("commands.pdf.long"),
	Example: "",
	RunE:    runPDF,
}

func init() {
	pdfCmd.Flags().StringSliceVarP(&pdfPaths, "path", "p", []string{"."}, locale.T("flags.path"))
	pdfCmd.Flags().StringVarP(&pdfOutput, "output", "o", "", locale.T("flags.output"))
	pdfCmd.Flags().StringVar(&pdfPreset, "preset", "A4Portrait", locale.T("flags.preset"))
	pdfCmd.Flags().StringVar(&pdfLang, "lang", "ru", locale.T("flags.lang"))
	pdfCmd.Flags().BoolVarP(&pdfRecursive, "recursive", "r", false, locale.T("flags.recursive"))
	pdfCmd.Flags().StringSliceVarP(&pdfExtensions, "extensions", "e", nil, locale.T("flags.extensions"))
	pdfCmd.Flags().StringSliceVar(&pdfBlacklist, "blacklist", nil, locale.T("flags.blacklist"))
	pdfCmd.Flags().StringSliceVarP(&pdfCategories, "categories", "c", nil, locale.T("flags.categories"))
	pdfCmd.Flags().StringVar(&pdfMinSize, "min-size", "", locale.T("flags.min_size"))
	pdfCmd.Flags().StringVar(&pdfMaxSize, "max-size", "", locale.T("flags.max_size"))
	pdfCmd.Flags().StringSliceVar(&pdfSortBy, "sort-by", []string{"name:asc"}, locale.T("flags.sort_by"))
	pdfCmd.Flags().StringVar(&pdfLimit, "limit", "", locale.T("flags.limit"))
	pdfCmd.Flags().StringVar(&pdfOffset, "offset", "", locale.T("flags.offset"))

	pdfCmd.Flags().StringVar(&pdfOrientation, "orientation", "", locale.T("flags.orientation"))
	pdfCmd.Flags().StringVar(&pdfPageSize, "page-size", "", locale.T("flags.page_size"))
	pdfCmd.Flags().Float64Var(&pdfMargin, "margin", 0, locale.T("flags.margin"))
	pdfCmd.Flags().StringVar(&pdfFontName, "font-name", "", locale.T("flags.font_name"))
	pdfCmd.Flags().Float64Var(&pdfFontSize, "font-size", 0, locale.T("flags.font_size"))
	pdfCmd.Flags().StringVar(&pdfFontPath, "font-path", "", locale.T("flags.font_path"))

	pdfCmd.Flags().StringVar(&pdfTableSeparator, "table-separator", "", locale.T("flags.table_separator"))
	pdfCmd.Flags().BoolVar(&pdfTableEnableCaption, "table-enable-caption", false, locale.T("flags.table_enable_caption"))
	pdfCmd.Flags().StringVar(&pdfTableCaptionPosition, "table-caption-position", "", locale.T("flags.table_caption_position"))
	pdfCmd.Flags().StringVar(&pdfTableCaptionAlignment, "table-caption-alignment", "", locale.T("flags.table_caption_alignment"))
	pdfCmd.Flags().Float64Var(&pdfTableCaptionMargin, "table-caption-margin", 0, locale.T("flags.table_caption_margin"))
	pdfCmd.Flags().StringVar(&pdfTableCaptionPrefix, "table-caption-prefix", "", locale.T("flags.table_caption_prefix"))
	pdfCmd.Flags().BoolVar(&pdfTableStretchWidth, "table-stretch-width", false, locale.T("flags.table_stretch_width"))

	pdfCmd.Flags().BoolVar(&pdfImageEnableCaption, "image-enable-caption", false, locale.T("flags.image_enable_caption"))
	pdfCmd.Flags().StringVar(&pdfImageCaptionPosition, "image-caption-position", "", locale.T("flags.image_caption_position"))
	pdfCmd.Flags().StringVar(&pdfImageCaptionAlignment, "image-caption-alignment", "", locale.T("flags.image_caption_alignment"))
	pdfCmd.Flags().Float64Var(&pdfImageCaptionMargin, "image-caption-margin", 0, locale.T("flags.image_caption_margin"))
	pdfCmd.Flags().StringVar(&pdfImageCaptionPrefix, "image-caption-prefix", "", locale.T("flags.image_caption_prefix"))

	pdfCmd.Flags().BoolVar(&pdfTextRenderMarkdown, "text-render-markdown", false, locale.T("flags.text_render_markdown"))

	pdfCmd.Flags().BoolVar(&pdfCodeStyleAsBlock, "code-style-as-block", false, locale.T("flags.code_style_as_block"))
	pdfCmd.Flags().BoolVar(&pdfCodeDisableHeader, "code-disable-header", false, locale.T("flags.code_disable_header"))
	pdfCmd.Flags().StringVar(&pdfCodeFontName, "code-font-name", "", locale.T("flags.code_font_name"))
	pdfCmd.Flags().Float64Var(&pdfCodeFontSize, "code-font-size", 0, locale.T("flags.code_font_size"))
	pdfCmd.Flags().StringVar(&pdfCodeBgColor, "code-bg-color", "", locale.T("flags.code_bg_color"))
	pdfCmd.Flags().StringVar(&pdfCodeTextColor, "code-text-color", "", locale.T("flags.code_text_color"))
}

func runPDF(cmd *cobra.Command, args []string) (err error) {
	start := time.Now()
	var recordedFiles, recordedBytes int64

	defer func() {
		telemetry.Record("pdf", recordedFiles, recordedBytes, start, err)
	}()

	ui.Title(locale.T("commands.pdf.short"))

	presetCfg, err := getPreset(pdfPreset, pdfLang)
	if err != nil {
		return err
	}

	if pdfOutput != "" {
		presetCfg.OutputPath = pdfOutput
	}

	applyOverrides(cmd, presetCfg)

	fsc := files.NewFolderSearchConfig(pdfRecursive, pdfPaths...)

	if len(pdfCategories) > 0 {
		categories := parseCategories(pdfCategories)
		pdf.ExportTypes(fsc, categories...)
	} else if len(pdfExtensions) > 0 {
		fsc.SearchForExtensions(pdfExtensions...)
	} else {
		pdf.ExportTypes(fsc, pdf.CategoryAll)
	}

	if len(pdfBlacklist) > 0 {
		fsc.AddToBlackList(pdfBlacklist...)
	}

	if pdfMinSize != "" {
		size, unit := parseSize(pdfMinSize)
		fsc.SetMinSize(files.SizeFromUnit(size, unit))
	}
	if pdfMaxSize != "" {
		size, unit := parseSize(pdfMaxSize)
		fsc.SetMaxSize(files.SizeFromUnit(size, unit))
	}
	if pdfLimit != "" {
		limit, err := strconv.ParseUint(pdfLimit, 10, 64)
		if err != nil {
			fsc.SetLimit(limit)
		}
	}
	if pdfOffset != "" {
		offset, err := strconv.ParseUint(pdfOffset, 10, 64)
		if err != nil {
			fsc.SetLimit(offset)
		}
	}

	sorter := buildSorter(pdfSortBy)

	start = time.Now()

	spinner := ui.NewSpinner(locale.T("messages.progress.collecting"))
	fileInfos, err := files.CollectFiles(fsc)
	fileInfos, amount := files.Paginate(fileInfos, fsc)
	spinner.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.generating_pdf"), err)
	}

	if len(fileInfos) == 0 {
		ui.Warn(locale.T("messages.errors.no_files_found"))
		return nil
	}

	sorter.SortFiles(fileInfos)

	bar := ui.NewProgressBar(int64(amount), locale.T("messages.progress.generating_pdf"))

	err = pdf.GeneratePDF(fileInfos, presetCfg, func() {
		bar.Add(1)
	})
	bar.Finish()

	if err != nil {
		return fmt.Errorf(locale.T("messages.errors.generating_pdf"), err)
	}

	recordedFiles = int64(len(fileInfos))

	elapsed := time.Since(start)
	printPDFResult(presetCfg, elapsed)

	return nil
}

func applyOverrides(cmd *cobra.Command, cfg *pdf.Config) {
	if cmd.Flags().Changed("orientation") {
		cfg.Orientation = pdf.Orientation(pdfOrientation)
	}
	if cmd.Flags().Changed("page-size") {
		cfg.PageSize = pdf.PageSize(pdfPageSize)
	}
	if cmd.Flags().Changed("margin") {
		cfg.Margin = pdfMargin
	}
	if cmd.Flags().Changed("font-name") {
		cfg.FontName = pdfFontName
	}
	if cmd.Flags().Changed("font-size") {
		cfg.FontSize = pdfFontSize
	}
	if cmd.Flags().Changed("font-path") {
		cfg.FontPath = pdfFontPath
	}

	if cmd.Flags().Changed("table-separator") {
		if len(pdfTableSeparator) > 0 {
			cfg.Table.Separator = rune(pdfTableSeparator[0])
		}
	}
	if cmd.Flags().Changed("table-enable-caption") {
		cfg.Table.EnableCaption = pdfTableEnableCaption
	}
	if cmd.Flags().Changed("table-caption-position") {
		cfg.Table.CaptionPosition = pdf.CaptionPosition(pdfTableCaptionPosition)
	}
	if cmd.Flags().Changed("table-caption-alignment") {
		cfg.Table.CaptionAlignment = pdf.CaptionAlignment(pdfTableCaptionAlignment)
	}
	if cmd.Flags().Changed("table-caption-margin") {
		cfg.Table.CaptionMargin = pdfTableCaptionMargin
	}
	if cmd.Flags().Changed("table-caption-prefix") {
		cfg.Table.CaptionPrefix = pdfTableCaptionPrefix
	}
	if cmd.Flags().Changed("table-stretch-width") {
		cfg.Table.StretchTableWidth = pdfTableStretchWidth
	}

	if cmd.Flags().Changed("image-enable-caption") {
		cfg.Image.EnableCaption = pdfImageEnableCaption
	}
	if cmd.Flags().Changed("image-caption-position") {
		cfg.Image.CaptionPosition = pdf.CaptionPosition(pdfImageCaptionPosition)
	}
	if cmd.Flags().Changed("image-caption-alignment") {
		cfg.Image.CaptionAlignment = pdf.CaptionAlignment(pdfImageCaptionAlignment)
	}
	if cmd.Flags().Changed("image-caption-margin") {
		cfg.Image.CaptionMargin = pdfImageCaptionMargin
	}
	if cmd.Flags().Changed("image-caption-prefix") {
		cfg.Image.CaptionPrefix = pdfImageCaptionPrefix
	}

	if cmd.Flags().Changed("text-render-markdown") {
		cfg.Text.RenderMarkdownForAll = pdfTextRenderMarkdown
	}

	if cmd.Flags().Changed("code-style-as-block") {
		cfg.Code.StyleAsBlock = pdfCodeStyleAsBlock
	}
	if cmd.Flags().Changed("code-disable-header") {
		cfg.Code.DisableHeader = pdfCodeDisableHeader
	}
	if cmd.Flags().Changed("code-font-name") {
		cfg.Code.FontName = pdfCodeFontName
	}
	if cmd.Flags().Changed("code-font-size") {
		cfg.Code.FontSize = pdfCodeFontSize
	}
	if cmd.Flags().Changed("code-bg-color") {
		if color, err := parseColor(pdfCodeBgColor); err == nil {
			cfg.Code.BgColor = color
		}
	}
	if cmd.Flags().Changed("code-text-color") {
		if color, err := parseColor(pdfCodeTextColor); err == nil {
			cfg.Code.TextColor = color
		}
	}
}

func parseColor(s string) ([3]int, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return [3]int{}, fmt.Errorf("invalid color format, expected R,G,B")
	}

	var color [3]int
	for i, part := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return [3]int{}, fmt.Errorf("invalid color value: %s", part)
		}
		if val < 0 || val > 255 {
			return [3]int{}, fmt.Errorf("color value must be between 0 and 255: %d", val)
		}
		color[i] = val
	}

	return color, nil
}

func getPreset(presetName, lang string) (*pdf.Config, error) {
	var bundle interface{}

	switch presetName {
	case "A4Portrait":
		bundle = presets.A4Portrait
	case "A4Landscape":
		bundle = presets.A4Landscape
	case "AcademicReport":
		bundle = presets.AcademicReport
	case "ScientificArticle":
		bundle = presets.ScientificArticle
	case "SourceCode":
		bundle = presets.SourceCode
	case "ServerLogs":
		bundle = presets.ServerLogs
	case "Cheatsheet":
		bundle = presets.Cheatsheet
	case "A3DataHeavy":
		bundle = presets.A3DataHeavy
	case "PhotoAlbum":
		bundle = presets.PhotoAlbum
	default:
		return nil, fmt.Errorf(locale.T("messages.errors.unknown_preset"), presetName)
	}

	switch b := bundle.(type) {
	case presets.PresetBundle:
		switch lang {
		case "en":
			return cloneConfig(b.EN), nil
		case "ru":
			return cloneConfig(b.RU), nil
		default:
			return cloneConfig(b.RU), nil
		}
	default:
		return nil, fmt.Errorf(locale.T("messages.errors.invalid_preset_type"))
	}
}

func cloneConfig(cfg *pdf.Config) *pdf.Config {
	clone := *cfg
	return &clone
}

func parseCategories(categories []string) []pdf.FileCategory {
	var result []pdf.FileCategory
	for _, c := range categories {
		switch c {
		case "image":
			result = append(result, pdf.CategoryImage)
		case "text":
			result = append(result, pdf.CategoryText)
		case "table":
			result = append(result, pdf.CategoryTable)
		case "code":
			result = append(result, pdf.CategoryCode)
		case "all":
			result = append(result, pdf.CategoryAll)
		}
	}
	return result
}

func printPDFResult(cfg *pdf.Config, elapsed time.Duration) {
	ui.PrintSection(locale.T("messages.results.pdf_title"))

	data := map[string]string{
		locale.T("messages.labels.output_file"): cfg.OutputPath,
		locale.T("messages.labels.preset"):      pdfPreset,
		locale.T("messages.labels.language"):    pdfLang,
		locale.T("messages.labels.orientation"): string(cfg.Orientation),
		locale.T("messages.labels.page_size"):   string(cfg.PageSize),
		locale.T("messages.labels.elapsed"):     elapsed.Round(time.Millisecond).String(),
	}

	ui.PrintStatsTable(data)
	ui.Success(fmt.Sprintf(locale.T("messages.success.pdf_generated"), cfg.OutputPath))
}
