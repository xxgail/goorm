package goorm

import (
	"fmt"
	"time"
)

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
