package goorm

import (
	"errors"
	"fmt"
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
