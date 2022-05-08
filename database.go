package glittersdk

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Database struct {
	client *Client
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type CreateSchemaBody struct {
	SchemaName string          `json:"schema_name"`
	Data       json.RawMessage `json:"data"`
}

// CreateSchema create schema
func (d *Database) CreateSchema(schemaName, schema string) (string, error) {
	body := CreateSchemaBody{
		SchemaName: schemaName,
		Data:       []byte(schema),
	}
	resp := new(response)
	err := d.client.post(urlCreateSchema, body, resp)
	if err != nil {
		return "", err
	}
	return resp.TX.String(), err
}

type Schema struct {
	Name      string
	ESMapping string
}

// GetSchema get schema by name
func (d *Database) GetSchema(name string) (string, error) {
	var s json.RawMessage
	req := map[string]string{
		"schema_name": name,
	}
	err := d.client.get(urlGetSchema, req, &s)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

// ListSchema list all glitter schemas
func (d *Database) ListSchema() (string, error) {
	var s json.RawMessage
	err := d.client.get(urlListSchema, nil, &s)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

type putDocRequest struct {
	SchemaName string      `json:"schema_name"`
	DocData    interface{} `json:"doc_data"`
}

// PutDoc put a document to glitter
// returns transcation id
func (d *Database) PutDoc(scheaName string, document interface{}) (string, error) {
	req := &putDocRequest{
		SchemaName: scheaName,
		DocData:    document,
	}
	resp := new(response)
	err := d.client.post(urlPutDoc, req, resp)
	if err != nil {
		return "", err
	}
	return resp.TX.String(), nil
}

var _ json.Marshaler = (*Document)(nil)
var _ json.Unmarshaler = (*Document)(nil)

type Document []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m Document) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *Document) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (d Document) Unmarshal(to interface{}) error {
	return json.Unmarshal(d, to)
}

func (d Document) String() string {
	buf := bytes.NewBuffer(nil)
	json.Indent(buf, []byte(d), "", "\t")
	return buf.String()
}

type getDocsReq struct {
	SchemaName string   `json:"schema_name"`
	DocIDs     []string `json:"doc_ids"`
}

type GetDocsResult struct {
	Total     int64               `json:"total"`
	Documents map[string]Document `json:"hits"`
}

// GetDocs get documents by id list
func (d *Database) GetDocs(scheaName string, docIDs []string) (*GetDocsResult, error) {
	req := &getDocsReq{
		SchemaName: scheaName,
		DocIDs:     docIDs,
	}
	resp := &GetDocsResult{}
	err := d.client.post(urlGetDocs, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type complexSearchReq struct {
	Index      string   `json:"index"`
	Query      string   `json:"query"`
	Filters    []Filter `json:"filters"`
	QueryField []string `json:"query_field"`
	AggsField  []string `json:"aggs_field"`

	OrderBy string `json:"order_by"`
	Limit   int    `json:"limit"`
	Page    int    `json:"page"`
}

// Search search from glitter database with given condition
func (d *Database) Search(cond *SearchCond) (*SearchResult, error) {
	req := &complexSearchReq{
		Index:      cond.schema,
		QueryField: cond.selects,
		AggsField:  cond.aggsFields,
		Query:      cond.query,
		Filters:    cond.filters,
		OrderBy:    cond.orderby,
		Limit:      cond.limit,
		Page:       cond.page,
	}

	resp := &SearchResult{}
	err := d.client.post(urlSearch, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type SearchResult struct {
	SearchTime    int64               `json:"search_time"`
	Index         string              `json:"index"`
	Meta          Meta                `json:"meta"`
	Items         []Item              `json:"items"`
	SortedByField []SortedByField     `json:"sorted_by_field"`
	Facet         map[string][]Filter `json:"facet"`
}

type Filter struct {
	Type     string      `json:"type"`
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
	From     float64     `json:"from"`
	To       float64     `json:"to"`
	DocCount int64       `json:"doc_count"`
}

type Item struct {
	Highlight map[string][]string `json:"highlight"`
	Data      Document            `json:"data"`
}

type Meta struct {
	Page PageInfo `json:"page"`
}

type PageInfo struct {
	CurrentPage  int64  `json:"current_page"`
	TotalPages   int64  `json:"total_pages"`
	TotalResults int64  `json:"total_results"`
	Size         int64  `json:"size"`
	SortedBy     string `json:"sorted_by"`
}

type SortedByField struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

type SearchCond struct {
	schema     string
	orderby    string
	query      string
	filters    []Filter
	limit      int
	page       int
	selects    []string `form:"query_field"`
	aggsFields []string `form:"aggs_field"`
}

func NewSearchCond() *SearchCond {
	return &SearchCond{}
}

func (s *SearchCond) Schema(schema string) *SearchCond {
	s.schema = schema
	return s
}

func (s *SearchCond) Select(fields ...string) *SearchCond {
	s.selects = fields
	return s
}

func (s *SearchCond) Query(q string) *SearchCond {
	s.query = q
	return s
}

func (s *SearchCond) Filter(f Filter) *SearchCond {
	s.filters = append(s.filters, f)
	return s
}

func (s *SearchCond) OrderBy(expr string) *SearchCond {
	s.orderby = expr
	return s
}

func (s *SearchCond) Limit(n int) *SearchCond {
	s.limit = n
	return s
}

func (s *SearchCond) Page(p int) *SearchCond {
	s.page = p
	return s
}

func (s *SearchCond) AggregateBy(fields ...string) *SearchCond {
	s.aggsFields = fields
	return s
}
