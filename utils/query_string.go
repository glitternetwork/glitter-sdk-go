package utils

import (
	"fmt"
	"strings"

	glittertypes "github.com/glitternetwork/glitter.proto/golang/glitter_proto/index/types"
)

type QueryString struct {
	Querys    []string
	DB        string
	Table     string
	Fields    string
	Highlight string
	Limit     int64
	Offset    int64
}

func NewCandyQueryString(db, table, fields string, limit, offset int64) *QueryString {
	return &QueryString{
		Querys: make([]string, 0),
		DB:     db,
		Table:  table,
		Fields: fields,
		Limit:  limit,
		Offset: offset,
	}
}

func NewQueryString() *QueryString {
	return &QueryString{Querys: make([]string, 0)}
}

func (q *QueryString) GetQueryString() string {
	return strings.Join(q.Querys, " ")
}

func (q *QueryString) Row(sql string, args ...interface{}) (string, []*glittertypes.Argument, error) {
	glitterArgs, err := toGlitterArguments(0, args)
	return sql, glitterArgs, err
}

func (q *QueryString) Build() (string, []*glittertypes.Argument, error) {
	glitterArg, err := toGlitterArgument(q.GetQueryString())
	if err != nil {
		return "", nil, err
	}
	if q.Limit > 0 {
		_sql := fmt.Sprintf("select %s %s from %s.%s where query_string(%s) limit %d, %d", q.Highlight, q.Fields, q.DB, q.Table, "?", q.Offset, q.Limit)
		return _sql, []*glittertypes.Argument{glitterArg}, nil
	} else {
		_sql := fmt.Sprintf("select %s %s from %s.%s where query_string(%s)", q.Highlight, q.Fields, q.DB, q.Table, "?")
		return _sql, []*glittertypes.Argument{glitterArg}, nil
	}
}

func (q *QueryString) AddHighLight(fields []string) error {
	_highlight := HighlightHint(fields)
	q.Highlight = _highlight
	return nil
}

func (q *QueryString) Add(query string) {
	q.Querys = append(q.Querys, query)
}

func (q *QueryString) AddMatchQuery(field, query string, boost float64) {
	q.Add(MatchQuery(field, query, boost))
}

func (q *QueryString) AddMatchPhraseQuery(field, query string, boost float64) {
	q.Add(MatchPhraseQuery(field, query, boost))
}

func (q *QueryString) AddRegexpQuery(field, query string, boost float64) {
	q.Add(RegexpQuery(field, query, boost))
}
func (q *QueryString) AddNumericRangeQuery(field, operator string, value int, boost float64) {
	q.Add(NumericRangeQuery(field, operator, value, boost))
}
func (q *QueryString) AddDateRangeQuery(field, operator, value string, boost float64) {
	q.Add(DateRangeQuery(field, operator, value, boost))
}

func HighlightHint(fields []string) string {
	arr := make([]string, 0)
	_fields := make([]any, 0)
	for i := 0; i < len(fields); i++ {
		arr = append(arr, `"%s"`)
		_fields = append(_fields, fields[i])
	}
	option := `/*+ SET_VAR(full_text_option='{"highlight":{ "style":"html","fields":[` + strings.Join(arr, ",") + `]}}')*/`
	return fmt.Sprintf(option, _fields...)
}

func MatchQuery(field, query string, boost float64) string {
	replacer := strings.NewReplacer(
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"&", "\\&",
		"|", "\\|",
		">", "\\>",
		"<", "\\<",
		"!", "\\!",
		"(", "\\(",
		")", "\\)",
		"{", "\\{",
		"}", "\\}",
		"[", "\\[",
		"]", "\\]",
		"^", "\\^",
		"~", "\\~",
		"*", "\\*",
		"?", "\\?",
		":", "\\:",
		"\\", "\\\\",
		`/`, `\/`,
		` `, `\ `,
		`"`, `\"`,
	)
	return fmt.Sprintf("%s:%s^%f", field, replacer.Replace(query), boost)
}

func MatchPhraseQuery(field, query string, boost float64) string {
	return fmt.Sprintf("%s:\"%s\"^%f", field, query, boost)
}

func RegexpQuery(field, query string, boost float64) string {
	return fmt.Sprintf("%s:/%s/^%f", field, query, boost)
}

func NumericRangeQuery(field, operator string, value int, boost float64) string {
	return fmt.Sprintf("%s:%s%d^%f", field, operator, value, boost)
}

func DateRangeQuery(field, operator, value string, boost float64) string {
	return fmt.Sprintf("%s:%s\"%s\"^%f", field, operator, value, boost)
}
