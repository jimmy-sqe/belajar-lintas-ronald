package pagination

func SQLOffset(page, perPage uint64) uint64 {
	if page == 0 {
		page = 1
	}
	return (page - 1) * perPage
}

func TotalPage(totalRow, perPage uint64) uint64 {
	return (totalRow + perPage - 1) / perPage
}

// Page matches the lintas API contract Pagination shape.
type Page struct {
	Page      uint64 `json:"page"`
	PageSize  uint64 `json:"page_size"`
	TotalData uint64 `json:"total_data"`
	TotalPage uint64 `json:"total_page"`
}

// Build returns a Page from raw inputs. Page < 1 is normalised to 1.
func Build(page, pageSize, totalData uint64) Page {
	if page == 0 {
		page = 1
	}
	return Page{
		Page:      page,
		PageSize:  pageSize,
		TotalData: totalData,
		TotalPage: TotalPage(totalData, pageSize),
	}
}
