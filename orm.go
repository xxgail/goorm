package goorm

import (
	"fmt"
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

	//query := "SELECT id," + field + " FROM " + t.TableName + t.getWhereStatement() + t.limitAndOffset()
	//rows, err := t.DB.Query(query)
	//if err != nil {
	//	fmt.Println("no this sql_table", err)
	//	return
	//}
	//if err != nil {
	//	return
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var id, fieldValue int
	//	if err = rows.Scan(&id, &fieldValue); err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	tx, err := t.DB.Begin()
	//	if err != nil {
	//		fmt.Println("tx fail")
	//	}
	//	query1 := "UPDATE " + t.TableName + " SET " + field + "=? WHERE id=?"
	//	fmt.Println("update query", query1)
	//	stmtUp, err := tx.Prepare(query1)
	//	if err != nil {
	//		fmt.Println("Prepare fail")
	//	}
	//	//设置参数以及执行sql语句
	//	res, err := stmtUp.Exec(fieldValue+num, id)
	//	if err != nil {
	//		fmt.Println("Exec fail")
	//	}
	//	tx.Commit()
	//	rowsAffected, err := res.RowsAffected()
	//	fmt.Println(rowsAffected)
	//}
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
	//query := "SELECT id," + field + " FROM " + t.TableName + t.getWhereStatement() + t.limitAndOffset()
	//rows, err := t.DB.Query(query)
	//if err != nil {
	//	fmt.Println("no this sql_table", err)
	//	return
	//}
	//if err != nil {
	//	return
	//}
	//defer rows.Close()
	//for rows.Next() {
	//	var id, fieldValue int
	//	if err = rows.Scan(&id, &fieldValue); err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	tx, err := t.DB.Begin()
	//	if err != nil {
	//		fmt.Println("tx fail")
	//	}
	//	query1 := "UPDATE " + t.TableName + " SET " + field + "=? WHERE id=?"
	//	fmt.Println("update query", query1)
	//	stmtUp, err := tx.Prepare(query1)
	//	if err != nil {
	//		fmt.Println("Prepare fail")
	//	}
	//	//设置参数以及执行sql语句
	//	res, err := stmtUp.Exec(fieldValue-num, id)
	//	if err != nil {
	//		fmt.Println("Exec fail")
	//	}
	//	tx.Commit()
	//	rowsAffected, err := res.RowsAffected()
	//	fmt.Println(rowsAffected)
	//}
}
