package goorm

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func (t *Table) Increment(field string, num int) {
	tx, err := t.DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	query := "UPDATE " + t.TableName + " SET " + field + "=" + field + "+?" + t.getWhereStatement()
	fmt.Println("update query", query)
	stmtUp, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("Prepare fail")
	}
	//设置参数以及执行sql语句
	res, err := stmtUp.Exec(num)
	if err != nil {
		fmt.Println("Exec fail")
	}
	tx.Commit()
	rowsAffected, err := res.RowsAffected()
	fmt.Println(rowsAffected)
}

func (t *Table) Decrement(field string, num int) {
	tx, err := t.DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}
	query := "UPDATE " + t.TableName + " SET " + field + "=" + field + "-?" + t.getWhereStatement()
	fmt.Println("update query", query)
	stmtUp, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("Prepare fail")
	}
	//设置参数以及执行sql语句
	res, err := stmtUp.Exec(num)
	if err != nil {
		fmt.Println("Exec fail")
	}
	tx.Commit()
	rowsAffected, err := res.RowsAffected()
	fmt.Println(rowsAffected)
}

func (t *Table) FirstOrCreate(attributes map[string]interface{}, values map[string]interface{}) (lastInsertId int64, err error) {
	lastInsertId, err = t.first(attributes)
	if err != nil {
		return
	} else if lastInsertId == 0 {
		create := mergeMaps(attributes, values)
		lastInsertId, err = t.create(create)
	}
	return
}

func (t *Table) first(wheres map[string]interface{}) (existId int64, err error) {
	t.wheresCondition = wheres
	query := "SELECT id FROM " + t.TableName + t.getWhereStatement() + " LIMIT 1"
	fmt.Println("First() ----- query", query)
	rows, err := t.DB.Query(query)
	if err != nil {
		fmt.Println("no this sql_table", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&existId); err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

func (t *Table) create(create map[string]interface{}) (lastInsertId int64, err error) {
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

func (t *Table) firstOrCreate(attributes map[string]interface{}, values map[string]interface{}) (lastInsertId int64, err error) {
	lastInsertId, err = t.first(attributes)
	if err != nil {
		return
	} else if lastInsertId == 0 {
		create := mergeMaps(attributes, values)
		lastInsertId, err = t.create(create)
	}
	return
}

func (t *Table) UpdateOrCreate(attributes map[string]interface{}, values map[string]interface{}) (err error) {
	var values_ map[string]interface{}
	id, err := t.firstOrCreate(attributes, values_)
	if err != nil {
		return
	}
	return t.update(id, values)
}

func (t *Table) update(id int64, update map[string]interface{}) (err error) {
	if len(update) == 0 {
		return errors.New("")
	}
	update["updated_at"] = time.Now()
	var fieldsKey []string
	var fieldsValue []interface{}
	for k, v := range update {
		fieldsKey = append(fieldsKey, k+" = ?")
		fieldsValue = append(fieldsValue, v)
	}

	fieldsValue = append(fieldsValue, id)
	fmt.Println(fieldsValue)
	tx, err := t.DB.Begin()
	if err != nil {
		fmt.Println("tx fail")
	}

	query := "UPDATE " + t.TableName + " SET " + strings.Join(fieldsKey, ",") + " WHERE id=?"
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
	_, err = res.RowsAffected()
	return
}
