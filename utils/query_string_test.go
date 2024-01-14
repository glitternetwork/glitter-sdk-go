package utils

import (
	"fmt"
	"testing"
)

func TestMatchPhraseQuery(t *testing.T) {
	field := "title"
	boost := 0.1
	m1 := MatchPhraseQuery(field, "aaa+bbb", boost)

	fmt.Printf("m1=%s\n", m1)
}

func TestMatchQuery(t *testing.T) {
	fmt.Println(MatchQuery("title", "aa+bb", 1))
	fmt.Println(MatchQuery("title", "aa\\\\bb", 1))
	fmt.Println(MatchQuery("title", "aa\"bb", 1))
	fmt.Println(MatchQuery("title", "aa\bb", 1))
	fmt.Println(MatchQuery("title", "aa'bb", 1))
	fmt.Println(MatchQuery("title", "aa\\bb", 1))
	fmt.Println(MatchQuery("title", "aa/bb", 1))
	fmt.Println(MatchQuery("title", "aa bb", 1))
	fmt.Println(MatchQuery("title", "aa^bb", 1))
}

func TestRegexpQuery(t *testing.T) {
	m1 := RegexpQuery("title", "aaa+bbb", 1)
	fmt.Printf("m1=%s\n", m1)
}

func TestNumericRangeQuery(t *testing.T) {
	m1 := NumericRangeQuery("title", "aaa+bbb", 1, 1)
	fmt.Printf("m1=%s\n", m1)
}

func TestDateRangeQuery(t *testing.T) {
	m1 := DateRangeQuery("title", "aaa+bbb", "d1", 1)
	fmt.Printf("m1=%s\n", m1)
}

func TestCandyQueryString(t *testing.T) {
	db := "library"
	table := "ebook_v4"

	qs := NewCandyQueryString(db, table, "_score,title,author", 0, 0)
	qs.AddMatchQuery("title", "哈利波特", 0.7)
	qs.AddMatchQuery("author", "哈利波特", 0.6)
	qs.AddMatchPhraseQuery("title", "哈利波特", 0.8)
	qs.AddDateRangeQuery("title", "aaa", "bbb", 1)
	qs.AddRegexpQuery("title", "ccc", 1)
	qs.AddNumericRangeQuery("title", "ddd", 12, 1)
	qs.AddHighLight([]string{"title", "author"})
	_sql, args, err := qs.Build()
	fmt.Printf("sql=%s,args=%+v, err=%+v\n", _sql, args, err)
}
