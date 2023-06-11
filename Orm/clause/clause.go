package clause

import (
	"fmt"
	"strings"
)

// sql语句
type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

type Type int
type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

// 各种数据库操作
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

func genBindvars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")

}

func _insert(values ...interface{}) (string, []interface{}) {
	// Insert into $tablename ($fileds)
	tablename := values[0]
	// 拼接插入的字段
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("insert into %s (%v)", tablename, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// values ($v1),($v2),...
	// 一次插入多条记录
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("values ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindvars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(",")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	// select $filed from $tablename
	tablename := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("select %s from %s", fields, tablename), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// limit $num
	return "limit ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// where $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("where %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("orderby %s", values[0]), []interface{}{}
}

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
