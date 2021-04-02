package goorm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func (t *Table) FindQuery(id string) string {
	// 查主键
	queryPrimary := "SELECT column_name FROM INFORMATION_SCHEMA.`KEY_COLUMN_USAGE` WHERE table_name='" + t.TableName + "' AND constraint_name='PRIMARY';"
	var primaryKey string
	if err := t.DB.QueryRow(queryPrimary).Scan(&primaryKey); err != nil {
		fmt.Println("没有主键")
		err = errors.New("no primary key")
		return ""
	}
	// 查主键 = id 的记录
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + " WHERE " + primaryKey + " = " + id + " AND `deleted_at` IS NULL LIMIT 1"
	return query
}

func (t *Table) GetQuery() string {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.getGroupBy() + t.orderBy() + t.limitAndOffset()
	log.Println("SELECT----" + query)
	return query
}

func (t *Table) FirstQuery() string {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.getGroupBy() + t.orderBy() + " LIMIT 1"
	log.Println("SELECT----" + query)
	return query
}

func (t *Table) Find(id string) (info map[string]interface{}, err error) {
	// 查主键
	queryPrimary := "SELECT column_name FROM INFORMATION_SCHEMA.`KEY_COLUMN_USAGE` WHERE table_name='" + t.TableName + "' AND constraint_name='PRIMARY';"
	var primaryKey string
	if err = t.DB.QueryRow(queryPrimary).Scan(&primaryKey); err != nil {
		fmt.Println("没有主键")
		err = errors.New("no primary key")
		return
	}
	// 查主键 = id 的记录
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + " WHERE " + primaryKey + " = " + id + " AND `deleted_at` IS NULL LIMIT 1"
	fmt.Println("Find() ----- query", query)
	rows, err := t.DB.Query(query)
	if err != nil {
		fmt.Println("no this sql_table", err)
		return
	}
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	defer rows.Close()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	info = map[string]interface{}{}
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range values {
			if v == nil {
				info[columns[k]] = ""
			} else {
				info[columns[k]] = string(v)
			}
		}
	}
	if len(info) == 0 {
		err = errors.New("no this info")
		return
	}
	return
}

func (t *Table) Get() (list []map[string]interface{}, err error) {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.getGroupBy() + t.orderBy() + t.limitAndOffset()
	fmt.Println("GET() ----- query", query)
	rows, err := t.DB.Query(query)
	if err != nil {
		fmt.Println("no this sql_table", err)
		return
	}
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	defer rows.Close()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list = []map[string]interface{}{}
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			fmt.Println(err)
			return
		}
		var item map[string]interface{}
		item = map[string]interface{}{}
		for k, v := range values {
			if v == nil {
				item[columns[k]] = ""
			} else {
				item[columns[k]] = string(v)
			}
		}
		list = append(list, item)
	}
	return
}

func (t *Table) First() (info map[string]interface{}, err error) {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.getGroupBy() + t.orderBy() + " LIMIT 1"
	fmt.Println("First() ----- query", query)
	rows, err := t.DB.Query(query)
	if err != nil {
		fmt.Println("no this sql_table", err)
		return
	}
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	defer rows.Close()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	//if rows.Next() == false {
	//	err = errors.New("no info")
	//	return
	//}
	info = map[string]interface{}{}
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range values {
			if v == nil {
				info[columns[k]] = ""
			} else {
				info[columns[k]] = string(v)
			}
		}
	}
	if len(info) == 0 {
		err = errors.New("no info")
		return
	}
	return
}

func (t *Table) Pluck(field string) (value string, err error) {
	query := "SELECT " + field + " FROM " + t.TableName + t.getWhereStatement() + t.orderBy() + " LIMIT 1"
	err = t.DB.QueryRow(query).Scan(&value)
	if err != nil {
		fmt.Println("找不到信息")
		err = errors.New("该信息不存在")
		return
	}
	return
}

func (t *Table) Count() (count int) {
	query := "SELECT COUNT(*)" + " FROM " + t.TableName + t.getWhereStatement()
	err := t.DB.QueryRow(query).Scan(&count)
	if err != nil {
		fmt.Println("无记录：", err)
	}
	return
}

func (t *Table) Sum(fields ...string) (sum map[string]string) {
	return t.math("SUM", fields...)
}

func (t *Table) Max(fields ...string) (max map[string]string) {
	return t.math("MAX", fields...)
}

func (t *Table) Min(fields ...string) (min map[string]string) {
	return t.math("MIN", fields...)
}

func (t *Table) Avg(fields ...string) (avg map[string]float64) {
	res := t.math("AVG", fields...)
	avg = map[string]float64{} // 保留两位小数的
	for k, v := range res {
		vv, _ := strconv.ParseFloat(v, 64)
		avg[k], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", vv), 64)
	}
	return
}

func (t *Table) math(operate string, args ...string) (res map[string]string) {
	var sumArr []string
	sumArr = []string{}
	for _, v := range args {
		sumArr = append(sumArr, operate+"("+v+")")
	}
	values := make([]sql.RawBytes, len(args))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	query := "SELECT " + strings.Join(sumArr, ",") + " FROM " + t.TableName + t.getWhereStatement()
	rows, err := t.DB.Query(query)
	if err != nil {
		fmt.Println("no this sql_table", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			fmt.Println("no info", err)
			return
		}
		res = map[string]string{}
		for k, v := range values {
			res[args[k]+operate] = string(v)
		}
	}
	return
}

func (t *Table) getSelect() (selectField string) {
	selectField = "*"
	if len(t.selectField) != 0 {
		selectField = strings.Join(t.selectField, ",")
	}
	return
}

func (t *Table) limitAndOffset() (limit string) {
	if t.limitNum == 0 {
		t.limitNum = 100
	}
	return " LIMIT " + strconv.Itoa(t.offsetNum) + "," + strconv.Itoa(t.limitNum)
}

func (t *Table) orderBy() (orderBy string) {
	var orderByArr []string
	if len(t.orderByMap) != 0 {
		for k, v := range t.orderByMap {
			if v {
				orderByArr = append(orderByArr, k+" ASC")
			} else {
				orderByArr = append(orderByArr, k+" DESC")
			}
		}
	}
	if len(t.orderByRawMap) != 0 {
		for k, v := range t.orderByRawMap {
			if v {
				orderByArr = append(orderByArr, "FIELD("+k+") ASC")
			} else {
				orderByArr = append(orderByArr, "FIELD("+k+") DESC")
			}
		}
	}
	if len(orderByArr) != 0 {
		orderBy = " ORDER BY " + strings.Join(orderByArr, ",")
	}
	return
}

func (t *Table) getWhereStatement() (whereState string) {
	var whereStateArr []string
	if len(t.wheresCondition) != 0 {
		for k, v := range t.wheresCondition {
			var tmp string
			switch v.(type) {
			case string:
				tmp = k + "=" + "'" + v.(string) + "'"
				break
			case int:
				tmp = k + "=" + strconv.Itoa(v.(int))
				break
			}
			whereStateArr = append(whereStateArr, tmp)
		}
	}
	if len(t.whereCondition) != 0 {
		for _, v := range t.whereCondition {
			tmp := v[0].(string) + v[1].(string) + v[2].(string)
			whereStateArr = append(whereStateArr, tmp)
		}
	}
	if len(t.whereInCondition) != 0 {
		for k, v := range t.whereInCondition {
			var inCon []string
			for _, vv := range v {
				inCon = append(inCon, vv.(string))
			}
			whereStateArr = append(whereStateArr, k+" IN ("+strings.Join(inCon, ",")+")")
		}
	}
	if len(t.whereNotInCondition) != 0 {
		for k, v := range t.whereNotInCondition {
			var inCon []string
			for _, vv := range v {
				inCon = append(inCon, vv.(string))
			}
			whereStateArr = append(whereStateArr, k+" NOT IN ("+strings.Join(inCon, ",")+")")
		}
	}
	if len(t.whereBetweenCondition) != 0 {
		for _, v := range t.whereBetweenCondition {
			whereStateArr = append(whereStateArr, v[0]+" BETWEEN '"+v[1]+"' AND '"+v[2]+"'")
		}
	}
	if len(t.whereNullField) != 0 {
		for _, v := range t.whereNullField {
			whereStateArr = append(whereStateArr, v+" IS NULL")
		}
	}
	if len(t.whereNotNullField) != 0 {
		for _, v := range t.whereNotNullField {
			whereStateArr = append(whereStateArr, v+" IS NOT NULL")
		}
	}
	if len(t.whereQuery) != 0 {
		for _, v := range t.whereQuery {
			whereStateArr = append(whereStateArr, "("+v+")")
		}
	}
	if len(t.leftJoinTable) == 0 {
		whereState = " WHERE deleted_at IS NULL"
	} else {
		whereState = " WHERE t.deleted_at IS NULL"
	}
	if len(whereStateArr) != 0 {
		whereState += " AND (" + strings.Join(whereStateArr, " AND ")
		if len(t.orWhere) != 0 {
			whereState += " OR " + strings.Join(t.orWhere, " OR ")
		}
		whereState += ")"
	}
	return whereState
}

func (t *Table) leftJoin() (leftJoin string) {
	if len(t.leftJoinTable) != 0 {
		var leftJoinArr []string
		leftJoinArr = []string{}
		leftJoin = " t "
		i := 97
		for k, v := range t.leftJoinTable {
			alias := string(rune(i + k))
			leftJoinArr = append(leftJoinArr, "LEFT JOIN "+v[0]+" "+alias+" ON "+"t."+v[1]+" = "+alias+"."+v[2])
		}
		leftJoin += strings.Join(leftJoinArr, " ")
	}
	return
}

func (t *Table) getGroupBy() (groupByField string) {
	if len(t.groupByMap) != 0 {
		groupByField = " GROUP BY " + strings.Join(t.groupByMap, ",")
	}
	return
}
