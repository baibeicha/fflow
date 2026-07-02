package files

// Paginate returns the desired page of files and the total number of elements on the page.
func Paginate(files []FileInfo, fsc *FolderSearchConfig) ([]FileInfo, uint64) {
	total := uint64(len(files))

	offset := uint64(0)
	if fsc.Offset != nil {
		offset = *fsc.Offset
	}

	if offset >= total {
		return []FileInfo{}, 0
	}

	limit := total - offset
	if fsc.Limit != nil {
		limit = *fsc.Limit
	}

	end := offset + limit

	if end > total {
		end = total
	}

	result := append([]FileInfo(nil), files[offset:end]...)

	return result, uint64(len(result))
}
