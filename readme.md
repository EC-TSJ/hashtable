# Set

[![Home](https://godoc.org/github.com/gookit/event?status.svg)](file:///D:/EC-TSJ/Documents/CODE/SOURCE/Go/pkg/lib/cli)
[![Build Status](https://travis-ci.org/gookit/event.svg?branch=master)](https://travis-ci.org/)
[![Coverage Status](https://coveralls.io/repos/github/gookit/event/badge.svg?branch=master)](https://coveralls.io/github/)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/event)](https://goreportcard.com/report/github.com/)

> **[EN README](README.md)**

HashTable es una librería para manipular la DataStructure hashtable.

## GoDoc

- [godoc for github](https://godoc.org/github.com/)

## Funciones Principales
--- 


Tiene los objetos siguientes:
- Type `interface` ***IHashTable*** con métodos:


	- *`Add(key Key, value ...Value)`* 
	- *`Set(key string, value Value)`* 
	- *`Get(key string) Value`* 
	- *`Remove(key string) (Record, bool)`* 
	- *`Keys() []string`* 
	- *`Values() []Value`* 
	- *`Size() int`* 
	- *`Capacity() int`* 
	- *`LoadFactor() float64`* 
	- *`GetHash(key string) int`* 
	- *`ContainsKey(key string) bool`* 
	- *`ContainsValue(value Value) bool`* 
	- *`Clone() HashTable`* 
	- *`String() string`*

- Type `struct` ***HashTable***, con los mismos métodos que el anterior.

- Type `struct` ***NotificationError***, que implementa las interfaz Error:
  - *`Error() string`*

- Funciones:

	- *`NewHashTable() HashTable`*





## Ejemplos
```go

	hash := hashtable.NewHashTable()
	vare := map[string]hashtable.Value{
		"corcon": 123,
		"val":    345,
		"kkk":    "lllllll",
	}
	hash.Add(vare)
	rdf, ok := hash.Remove("corcon")
	hash.Add("jojo", 8596)
	fda := hash.Keys()
	fmt.Println(rdf, ok, fda)



```
## Notas





<!-- - [gookit/ini](https://github.com/gookit/ini) INI配置读取管理，支持多文件加载，数据覆盖合并, 解析ENV变量, 解析变量引用
-->
## LICENSE

**[MIT](LICENSE)**
