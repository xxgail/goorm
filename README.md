# goorm
[![PkgGoDev](https://pkg.go.dev/badge/github.com/xxgail/goorm)](https://pkg.go.dev/github.com/xxgail/goorm)

goorm 自用封装包（仿照laravel）

🍬 Add the library to your $GOPATH/src

`go get github.com/xxgail/goorm`

### 目前实现的方法有：

#### 初始化
- NewTable("表名",db) : 初始化方法 db 目前仅有*sql.DB

#### 查询
1. 查询条件
- Select(fields ...string) : 筛选字段
- Where(map[string]interface{}) : where数组查询
- WhereCompare(field string, compare string, value interface{}) : where比较查询
- WhereIN(field string, scope []interface{})
- WhereNotIn(field string, scope []interface{})
- WhereBetween(field string, left string, right string)
- WhereNull(fields ...string) : 筛选为空字段
- WhereNotNull(fields ...string) : 筛选不为空字段
- OrderBy(field string, sort bool) : sort-true:ASC, sort-false:DESC
- Limit(limit int)
- Offset(offset int)
- LeftJoin(tableName string, foreignKey string, primaryKey string) : _因为写的太乱，目前还不太实用_
- GroupBy(fields ...string)

2. 查询语句
- Find(id string) map,err : 返回查询主键单条信息，只需要传主键的值（前提是只有一个主键，待优化吧..
- Get() []map,err : 获取多条信息，可进行分页（待测试
- First() []map,err : 获取单条信息
- FindQuery(id string) string : 返回查询主键单条信息的query语句
- GetQuery() string : 返回查询多条信息的query语句
- FirstQuery() string : 返回查询单条信息的query语句

#### 增加
- InsertGetId(create map[string]interface{}) (lastInsertId int64, err error) : 新增单条信息，并返回新增加的ID

#### 修改
- Update(update map[string]interface{}) (rowsAffected int64, err error) : 修改信息，返回修改的记录数量

#### 删除
- Delete() (rowsAffected int64, err error) : 删除信息，返回删除的记录数量

#### 其他操作
- Pluck(field string) : 获取单个字段的信息
- Count() : 返回记录的数量
- Sum(fields ...string) map[string]string
- Max(fields ...string) map[string]string
- Min(fields ...string) map[string]string
- Avg(fields ...string) map[string]float64
- Increment(field string, num int) : 字段field的值增加num
- Decrement(field string, num int) : 字段field的值减少num

### example.go
```go
package main
import (
    "database/sql"
    "fmt"
    gorm "github.com/xxgail/goorm"
    "strings"
)

func main()  {
    var DB *sql.DB
    username := "root"
    password := "root"
    host := "127.0.0.1"
    port := "3306"
    database := "test"
    path := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=utf8&loc=Asia%2FShanghai&parseTime=true"}, "")
    fmt.Println(path)
    DB, _ = sql.Open("mysql", path) 
    DB.SetConnMaxLifetime(100)
    DB.SetMaxIdleConns(10)

    // 查询
    list,err := gorm.NewTable("students",DB).WhereCompare("name","=","gai").OrderBy("age",true).Get()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(list)
    
    // 插入
    var insert map[string]interface{}
    insert = map[string]interface{}{
        "name" : "gai1",
        "age" : 18,
    }
    newId,err := gorm.NewTable("students",DB).InsertGetId(insert)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(newId)
    
    // 修改
    var where map[string]interface{}
    where = map[string]interface{}{
        "name" : "gai1",
    }
    var update map[string]interface{}
    update = map[string]interface{}{
    	"age" : 19,
    }
    affect,err := gorm.NewTable("students",DB).Where(where).Update(update)
    if err != nil {
    	fmt.Println(err)
    }
    fmt.Println(affect)
}
```