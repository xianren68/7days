package schema

import (
	"Orm/dialect"
	"go/ast"
	"reflect"
)

// 数据库字段结构体
type Filed struct {
	Name string // 字段名
	Type string // 字段类型
	Tag  string // 字段约束
}

// 存储映射的表字段和对应结构体
type Schema struct {
	Model      interface{}       // 原型
	Name       string            // 对应的表名
	Fileds     []*Filed          // 所有的字段
	FiledNames []string          // 字段名
	filedMap   map[string]*Filed // 字段与字段名映射表

}

// 获取字段属性
func (s *Schema) GetFiled(name string) *Filed {
	return s.filedMap[name]
}

// 将表结构体转化为对应的Schema
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		Name:  modelType.Name(),
		// slice虽然零值为nil，但go语言对其进行了特殊处理，不用手动定义
		filedMap: make(map[string]*Filed),
	}
	for i := 0; i < modelType.NumField(); i++ {
		// 获取结构体每一项
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			filed := &Filed{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// 将tag转化为字段约束
			if v, ok := p.Tag.Lookup("Orm"); ok {
				filed.Tag = v
			}
			schema.Fileds = append(schema.Fileds, filed)
			schema.FiledNames = append(schema.FiledNames, p.Name)
			schema.filedMap[p.Name] = filed
		}

	}
	return schema
}
