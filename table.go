package go_orm

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
	}
}

func (t *Table) Where(where map[string]interface{}) *Table {
	t.whereCondition = where
	return t
}

func (t *Table) WhereCompare(fieldName string, condition string, fieldValue interface{}) *Table {
	if condition == "=" {
		t.whereCondition[fieldName] = fieldValue
	} else {
		t.compareCondition = append(t.compareCondition, []interface{}{fieldName, condition, fieldValue})
	}
	return t
}

func (t *Table) WhereBetween(fieldName string, left string, right string) *Table {
	t.whereBetweenCondition = append(t.whereBetweenCondition, []string{fieldName, left, right})
	return t
}

func (t *Table) WhereIn(fieldName string, scope []interface{}) *Table {
	t.whereInCondition[fieldName] = scope
	return t
}

func (t *Table) WhereNotIn(fieldName string, scope []interface{}) *Table {
	t.whereNotInCondition[fieldName] = scope
	return t
}

func (t *Table) WhereNull(nullField ...string) *Table {
	t.whereNullField = nullField
	return t
}

func (t *Table) WhereNotNull(nullField ...string) *Table {
	t.whereNotNullField = nullField
	return t
}

func (t *Table) OrderBy(fieldName string, sort bool) *Table {
	t.orderByMap[fieldName] = sort
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

func (t *Table) Select(arg ...string) *Table {
	t.selectField = arg
	return t
}

func (t *Table) LeftJoin(tableName string, foreignKey string, primaryKey string) *Table {
	t.leftJoinTable = append(t.leftJoinTable, []string{tableName, foreignKey, primaryKey})
	return t
}
