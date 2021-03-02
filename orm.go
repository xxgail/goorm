package go_orm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func (t *Table) InsertGetId(create map[string]interface{}) (lastInsertId int64, err error) {
	if len(create) == 0 {
		return 0, errors.New("")
	}
	create["created_at"] = time.Now()
	create["updated_at"] = time.Now()
	var fieldsKey []string
	var fieldsKey_ []string
	var fieldsValue []interface{}
	for k, v := range create {
		fieldsKey = append(fieldsKey, "`"+k+"`")
		fieldsValue = append(fieldsValue, v)
		fieldsKey_ = append(fieldsKey_, "?")
	}

	tx, err := t.DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}

	query := "INSERT INTO " + t.TableName + " ( " + strings.Join(fieldsKey, ",") + " ) VALUES ( " + strings.Join(fieldsKey_, ",") + " )"
	fmt.Println("Insert query---", query)
	stem, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("InsertLog----Prepare fail" + t.TableName)
	}
	res, err := stem.Exec(fieldsValue...)
	if err != nil {
		fmt.Println("InsertLog----Exec fail" + t.TableName)
	}
	tx.Commit()
	lastInsertId, err = res.LastInsertId()
	return
}

func (t *Table) Update(update map[string]interface{}) (rowsAffected int64, err error) {
	if len(update) == 0 {
		return 0, errors.New("")
	}
	update["updated_at"] = time.Now()
	var fieldsKey []string
	var fieldsValue []interface{}
	for k, v := range update {
		fieldsKey = append(fieldsKey, k+" = ?")
		fieldsValue = append(fieldsValue, v)
	}

	whereState, fieldsValue_ := t.getUpdateWhereStatement()
	fieldsValue = append(fieldsValue, fieldsValue_...)
	fmt.Println(fieldsValue)
	tx, err := t.DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}

	query := "UPDATE " + t.TableName + " SET " + strings.Join(fieldsKey, ",") + whereState
	fmt.Println("update query", query)
	stmtUp, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("Prepare fail")
	}
	//设置参数以及执行sql语句
	res, err := stmtUp.Exec(fieldsValue...)
	if err != nil {
		fmt.Println("Exec fail")
	}
	tx.Commit()
	rowsAffected, err = res.RowsAffected()
	return
}

func (t *Table) Delete() (rowsAffected int64, err error) {
	whereState, fieldsValue := t.getUpdateWhereStatement()
	fieldsValue = append([]interface{}{time.Now()}, fieldsValue...)

	tx, err := t.DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}

	query := "UPDATE " + t.TableName + " SET deleted_at = ?" + whereState
	//query := "DELETE FROM " + t.TableName + whereState
	fmt.Println("delete query", query)
	stmtUp, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("Prepare fail")
	}
	//设置参数以及执行sql语句
	res, err := stmtUp.Exec(fieldsValue...)
	if err != nil {
		fmt.Println("Exec fail")
	}
	tx.Commit()
	rowsAffected, err = res.RowsAffected()
	return
}

func (t *Table) FindQuery(primaryName string, primaryKey string) string {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + " WHERE " + primaryName + " = " + primaryKey + " AND `deleted_at` IS NULL"
	fmt.Println(query)
	return query
}

func (t *Table) GetQuery() string {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.orderBy() + t.limitAndOffset()
	log.Println("SELECT----" + query)
	return query
}

func (t *Table) FirstQuery() string {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.orderBy() + " LIMIT 1"
	log.Println("SELECT----" + query)
	return query
}

func (t *Table) Find(primaryName string, primaryKey string) (info []string, err error) {
	if len(t.selectField) == 0 {
		err = errors.New("请填写select")
		return
	}
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + " WHERE " + primaryName + " = " + primaryKey + " AND `deleted_at` IS NULL"
	values := make([]sql.RawBytes, len(t.selectField))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	err = t.DB.QueryRow(query).Scan(scanArgs...)
	if err != nil {
		fmt.Println("找不到信息")
		return
	}
	info = []string{}
	for _, v := range values {
		if v == nil {
			info = append(info, "")
		} else {
			info = append(info, string(v))
		}
	}
	return
}

func (t *Table) Get() (list []map[string]interface{}, err error) {
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.orderBy() + t.limitAndOffset()
	rows, err := t.DB.Query(query)
	columns, err := rows.Columns()
	if err != nil {

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
	query := "SELECT " + t.getSelect() + " FROM " + t.TableName + t.leftJoin() + t.getWhereStatement() + t.orderBy() + " LIMIT 1"
	rows, err := t.DB.Query(query)
	columns, err := rows.Columns()
	if err != nil {

	}
	defer rows.Close()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			fmt.Println(err)
			return
		}
		info = map[string]interface{}{}
		for k, v := range values {
			if v == nil {
				info[columns[k]] = ""
			} else {
				info[columns[k]] = string(v)
			}
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
	if len(t.orderByMap) != 0 {
		var orderByArr []string
		for k, v := range t.orderByMap {
			if v {
				orderByArr = append(orderByArr, k+" ASC")
			} else {
				orderByArr = append(orderByArr, k+" DESC")
			}
		}
		orderBy = " ORDER BY " + strings.Join(orderByArr, ",")
	}
	return
}

func (t *Table) getWhereStatement() (whereState string) {
	var whereStateArr []string
	if len(t.whereCondition) != 0 {
		for k, v := range t.whereCondition {
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
	if len(t.compareCondition) != 0 {
		for _, v := range t.compareCondition {
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
	if len(t.leftJoinTable) == 0 {
		whereState = " WHERE deleted_at IS NULL"
	} else {
		whereState = " WHERE t.deleted_at IS NULL"
	}
	if len(whereStateArr) != 0 {
		whereState += " AND " + strings.Join(whereStateArr, " AND ")
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

func (t *Table) getUpdateWhereStatement() (string, []interface{}) {
	var whereState string
	var fieldsValue []interface{}
	var whereStateArr []string
	if len(t.whereCondition) != 0 {
		for k, v := range t.whereCondition {
			whereStateArr = append(whereStateArr, k+"=?")
			fieldsValue = append(fieldsValue, v)
		}
	}
	if len(t.compareCondition) != 0 {
		for _, v := range t.compareCondition {
			whereStateArr = append(whereStateArr, v[0].(string)+v[1].(string)+" ?")
			fieldsValue = append(fieldsValue, v[2])
		}
	}
	if len(t.whereInCondition) != 0 {
		for k, v := range t.whereInCondition {
			replace := make([]string, len(v))
			for k, _ := range replace {
				replace[k] = "?"
			}
			whereStateArr = append(whereStateArr, k+" IN ("+strings.Join(replace, ",")+")")
			fieldsValue = append(fieldsValue, v...)
		}
	}
	if len(t.whereNotInCondition) != 0 {
		for k, v := range t.whereNotInCondition {
			replace := make([]string, len(v))
			for k, _ := range replace {
				replace[k] = "?"
			}
			whereStateArr = append(whereStateArr, k+" NOT IN ("+strings.Join(replace, ",")+")")
			fieldsValue = append(fieldsValue, v...)
		}
	}
	if len(t.whereBetweenCondition) != 0 {
		for _, v := range t.whereBetweenCondition {
			whereStateArr = append(whereStateArr, v[0]+" BETWEEN "+"?"+" AND "+"?")
			fieldsValue = append(fieldsValue, []interface{}{v[1], v[2]}...)
		}
	}
	if len(t.whereNullField) != 0 {
		for _, v := range t.whereNullField {
			whereStateArr = append(whereStateArr, v+" IS NULL ")
		}
	}
	if len(t.whereNotNullField) != 0 {
		for _, v := range t.whereNotNullField {
			whereStateArr = append(whereStateArr, v+" IS NOT NULL ")
		}
	}
	whereState = " WHERE deleted_at IS NULL"
	if len(whereStateArr) != 0 {
		whereState += " AND " + strings.Join(whereStateArr, " AND ")
	}
	return whereState, fieldsValue
}