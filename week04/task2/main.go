package main

import (
	"bytes"
	"fmt"
	"reflect"
)

// 自定义错误类型
type UnsupportedTypeError struct {
	Type reflect.Type
}

// 实现error接口的Error方法
func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("不支持的类型: %v", e.Type)
}

func Marshal(v interface{}) ([]byte, error) {
	//获取类型
	rv := reflect.ValueOf(v)
	//存字符串
	var buf bytes.Buffer
	//把"{"提前输入
	buf.WriteString("{")
	//获取整个传入结构体的字段数
	numFields := rv.NumField()
	for i := 0; i < numFields; i++ {
		//不是第一个就加上","
		if i != 0 {
			buf.WriteString(",")
		}
		//获取该字段
		field := rv.Type().Field(i)
		//获取对应的json标签
		jsonKey := field.Tag.Get("json")
		//不存在的话自动填充
		if jsonKey == "" {
			jsonKey = field.Name
		}
		//写入jsonKey
		buf.WriteString(fmt.Sprintf("\"%s\":", jsonKey))
		//调用值处理函数
		value, err := marshalValue(rv.Field(i))
		if err != nil {
			return nil, err
		}
		//最终写入值
		buf.Write(value)
	}
	//输入"}"到buf里
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// 值处理函数
func marshalValue(rv reflect.Value) ([]byte, error) {
	//对整个传入值进行处理
	switch rv.Kind() {
	case reflect.String:
		return []byte(fmt.Sprintf("\"%s\"", rv.String())), nil
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return []byte(fmt.Sprintf("%d", rv.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16:
		return []byte(fmt.Sprintf("%d", rv.Uint())), nil
	case reflect.Struct, reflect.Map, reflect.Slice:
		return Marshal(rv.Interface())
	default:
		return nil, &UnsupportedTypeError{Type: rv.Type()}
	}
}

func main() {
	type Student struct {
		Name string `json:"name"`
		Age  uint   `json:"age"`
	}
	s := Student{
		Name: "孙俊楠",
		Age:  19,
	}
	jsonStr, err := Marshal(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonStr))
}
