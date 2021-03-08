package goorm

import (
	"database/sql"
)

type Table struct {
	TableName             string
	DB                    *sql.DB
	whereCondition        map[string]interface{}
	compareCondition      [][]interface{}
	whereBetweenCondition [][]string
	whereInCondition      map[string][]interface{}
	whereNotInCondition   map[string][]interface{}
	orderByMap            map[string]bool
	limitNum              int
	offsetNum             int
	selectField           []string
	whereNullField        []string
	whereNotNullField     []string
	leftJoinTable         [][]string
	groupByMap            []string
	orWhere               []string
}

func NewTable(tableName string, db *sql.DB) *Table {
	return &Table{
		TableName:             tableName,
		DB:                    db,
		whereCondition:        map[string]interface{}{},
		compareCondition:      [][]interface{}{},
		whereBetweenCondition: [][]string{},
		whereInCondition:      map[string][]interface{}{},
		whereNotInCondition:   map[string][]interface{}{},
		orderByMap:            map[string]bool{},
		selectField:           []string{},
		whereNullField:        []string{},
		whereNotNullField:     []string{},
		leftJoinTable:         [][]string{},
		groupByMap:            []string{},
		orWhere:               []string{},
	}
}

func (t *Table) Where(where map[string]interface{}) *Table {
	t.whereCondition = where
	return t
}

func (t *Table) WhereCompare(field string, compare string, value interface{}) *Table {
	if compare == "=" {
		t.whereCondition[field] = value
	} else {
		t.compareCondition = append(t.compareCondition, []interface{}{field, compare, value})
	}
	return t
}

func (t *Table) WhereBetween(field string, left string, right string) *Table {
	t.whereBetweenCondition = append(t.whereBetweenCondition, []string{field, left, right})
	return t
}

func (t *Table) WhereIn(field string, scope []interface{}) *Table {
	t.whereInCondition[field] = scope
	return t
}

func (t *Table) WhereNotIn(field string, scope []interface{}) *Table {
	t.whereNotInCondition[field] = scope
	return t
}

func (t *Table) WhereNull(fields ...string) *Table {
	t.whereNullField = append(t.whereNullField, fields...)
	return t
}

func (t *Table) WhereNotNull(fields ...string) *Table {
	t.whereNotNullField = append(t.whereNotNullField, fields...)
	return t
}

func (t *Table) OrderBy(field string, sort bool) *Table {
	t.orderByMap[field] = sort
	return t
}

func (t *Table) Limit(limit int) *Table {
	t.limitNum = limit
	return t
}

func (t *Table) Offset(offset int) *Table {
	t.offsetNum = offset
	return t
}

func (t *Table) Select(fields ...string) *Table {
	t.selectField = append(t.selectField, fields...)
	return t
}

func (t *Table) LeftJoin(tableName string, foreignKey string, primaryKey string) *Table {
	t.leftJoinTable = append(t.leftJoinTable, []string{tableName, foreignKey, primaryKey})
	return t
}

func (t *Table) GroupBy(fields ...string) *Table {
	t.groupByMap = append(t.groupByMap, fields...)
	return t
}

func (t *Table) OrWhere(condition ...string) *Table {
	t.orWhere = append(t.orWhere, condition...)
	return t
}
