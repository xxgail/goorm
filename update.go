package goorm

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

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
