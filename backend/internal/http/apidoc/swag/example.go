// Package swag contains annotation examples for swaggo/swag. The annotations
// are not active without `make swag-gen`. Copy patterns from this file when
// annotating real handlers in internal/http/handler/.
//
// All responses follow the lintas envelope (see internal/http/response).
package swag

// ExampleListItems documents the GET /v1/items endpoint pattern.
//
//	@Summary     List items
//	@Description Returns a paged list of items, optionally filtered.
//	@Tags        items
//	@Accept      json
//	@Produce     json
//	@Param       owner_id  query    string false "Filter by owner ID"
//	@Param       status    query    string false "Filter by status (active|archived)"
//	@Param       page      query    int    false "Page number (default 1)"
//	@Param       page_size query    int    false "Page size (default 10)"
//	@Success     200       {object} PaginatedItemsEnvelope
//	@Failure     401       {object} ErrorEnvelope
//	@Router      /v1/items [get]
func ExampleListItems() {}

// ExampleCreateItem documents the POST /v1/items endpoint pattern.
//
//	@Summary     Create item
//	@Tags        items
//	@Accept      json
//	@Produce     json
//	@Param       body body     CreateItemRequest      true "Item to create"
//	@Success     201  {object} ItemEnvelope
//	@Failure     400  {object} ErrorEnvelope
//	@Failure     401  {object} ErrorEnvelope
//	@Failure     422  {object} ErrorEnvelope
//	@Router      /v1/items [post]
func ExampleCreateItem() {}

// CreateItemRequest is a placeholder to make the swag annotations resolve.
// Replace references with the real DTO from internal/http/handler/dto.
type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int64  `json:"price"`
}

// ItemResponse is a placeholder.
type ItemResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

// PaginationMeta mirrors pkg/pagination.Page for swag refs.
type PaginationMeta struct {
	Page      uint64 `json:"page"`
	PageSize  uint64 `json:"page_size"`
	TotalData uint64 `json:"total_data"`
	TotalPage uint64 `json:"total_page"`
}

// ItemEnvelope mirrors internal/http/response.SuccessEnvelope for swag refs.
type ItemEnvelope struct {
	Success   bool         `json:"success"`
	Code      int          `json:"code"`
	Message   string       `json:"message"`
	Data      ItemResponse `json:"data"`
	Timestamp string       `json:"timestamp"`
}

// PaginatedItemsEnvelope mirrors internal/http/response.PaginatedEnvelope for swag refs.
type PaginatedItemsEnvelope struct {
	Success    bool           `json:"success"`
	Code       int            `json:"code"`
	Message    string         `json:"message"`
	Data       []ItemResponse `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
	Timestamp  string         `json:"timestamp"`
}

// ErrorEnvelope mirrors internal/http/response.ErrorEnvelope for swag refs.
type ErrorEnvelope struct {
	Success   bool                   `json:"success"`
	Code      int                    `json:"code"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
