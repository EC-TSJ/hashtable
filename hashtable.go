package hashtable

import (
	"ec-tsj/core"
	"fmt"
	"strings"
)

const (
	_MINLOADFACTOR   = 0.25
	_MAXLOADFACTOR   = 0.75
	_DEFAULTABLESIZE = 3
)

type (
	// Datos
	Value = core.T
	Key   = core.T

	//
	// Esquema:
	//
	//        HashTable
	//        ---------
	//    *int       		*Hash
	//     ---        	 ----
	//           error  			[]*Record
	//           -----    			-------
	//                   string Value *Record
	//                   ------ ----- -------
	Record struct {
		key   string
		value Value
		next  *Record
	}

	Hash struct {
		records []*Record
		_error_ error
	}

	HashTable struct {
		table   *Hash
		records *int
	}

	//...
	IHashTable interface {
		Add(Key, ...Value)
		Set(string, Value)
		Get(string) Value
		Remove(string) (Record, bool)
		Keys() []string
		Values() []Value
		Size() int
		Capacity() int
		LoadFactor() float64
		GetHash(string) int
		ContainsKey(string) bool
		ContainsValue(Value) bool
		Clone() IHashTable
		String() string
	}
)

// createHashTable: Llamdo por checkLoadFactorAndUpdate cuando se crea un nuevo hash
// @param {int}
// @return {HashTable}
func _createHashTable(tableSize int) HashTable {
	num := 0
	hash := Hash{make([]*Record, tableSize), &core.BaseError{-1, "", ""}}
	return HashTable{table: &hash, records: &num}
}

// NewHashTable: Llamado por el usuario para crear una hashtable.
// @return {HashTable}
func NewHashTable() IHashTable {
	num := 0
	hash := Hash{make([]*Record, _DEFAULTABLESIZE), &core.BaseError{-1, "", ""}}
	return &HashTable{table: &hash, records: &num}
}

// hash: Método de Horner, calcula el hash
// @param {string}
// @return {int}
func _hash(key string) int {
	h := 0
	for i := 0; i < len(key); i++ {
		h = 31*h + int(key[i])
	}

	return int(h)
}

// Método hash
//
// func (sha)_hash(key string) uint32 {
// 	c := fnv.New32a()
// 	c.Write([]byte(key))
// 	return c.Sum32()
// }

// index: Usado para calcular el índice registro dentro del slice.
// @param{string}
// @param {int}
// @return {int}
func _index(key string, size int) int {
	return _hash(key) % size
}

// // Debug: Imprime la hashtable tal y cual es
// func (h *HashTable) Debug() {
// 	fmt.Printf("----------%d elementos-------\n", *h.records)
// 	for i, node := range h.table.records {
// 		fmt.Printf("%d :", i)
// 		for node != nil {
// 			fmt.Printf("[%s, %d]->", node.key, node.value)
// 			node = node.next
// 		}
// 		fmt.Println("nil")
// 	}
// }

// String: Interface Stringer
func (h *HashTable) String() string {
	items := make([]string, 0)
	for _, node := range h.table.records {
		for node != nil {
			items = append(items, fmt.Sprintf("%s,{%#v}", node.key, node.value))
			node = node.next
		}
	}

	return fmt.Sprintf("((%s))", strings.Join(items, ", "))
}

// Add: Inserta un par Clave:Valor
// @param {string}
// @param {Value}
// @returm {bool}
func (h *HashTable) _add(key string, value Value) bool {
	index := _index(key, len(h.table.records))
	iterator := h.table.records[index]
	node := Record{key, value, nil}
	if iterator == nil {
		h.table.records[index] = &node
	} else {
		prev := &Record{"", 0, nil}
		for iterator != nil {
			if iterator.key == key { // Key already exists
				//iterator.value = value // FIX: operación desechada
				return false
			}
			prev = iterator
			iterator = iterator.next
		}
		prev.next = &node
	}
	*h.records += 1

	return true
}

// Add: Inserta un par Clave:Valor
// @param {string}
// @param {Value}
func (h *HashTable) Add(key Key, value ...Value) { // FIX: original versión: Add(string, Value)
	if v, ok := key.(map[string]Value); ok {
		for key, val := range v {
			sizeChanged := h._add(key, val)
			if sizeChanged /*== true */ {
				h._checkLoadFactorAndUpdate()
			}
		}
	} else if v, ok := key.(string); ok {
		//  FIX: Código original
		//  FIX: from
		sizeChanged := h._add(v, value[0])
		if sizeChanged /*== true */ {
			h._checkLoadFactorAndUpdate()
		}
		// FIX: to
	}
}

// Set: Cambia el valor a 'key' y le da el valor 'value'
// @param {string}
// @param {Value}
func (h *HashTable) Set(key string, value Value) {
	for _, node := range h.table.records {
		for node != nil {
			if key == node.key {
				node.value = value
			}
			node = node.next
		}
	}
}

// Get: Recupera un valor para una clave
// @param {string}
// @return {Value}
func (h *HashTable) Get(key string) Value {
	index := _index(key, len(h.table.records))
	iterator := h.table.records[index]
	for iterator != nil {
		if iterator.key == key { // Key already exists
			return iterator.value
		}
		iterator = iterator.next
	}

	return nil
}

// Remove: Remueve una clave (y valor)
// @param {string}
// @return {Value}
func (h *HashTable) _remove(key string) (Record, bool) {
	index := _index(key, len(h.table.records))
	iterator := h.table.records[index]
	if iterator == nil {
		return Record{"", nil, nil}, false
	}
	if iterator.key == key {
		dev := h.table.records[index]
		h.table.records[index] = iterator.next
		*h.records--
		return *dev, true
	} else {
		prev := iterator
		iterator = iterator.next
		for iterator != nil {
			if iterator.key == key {
				dev := h.table.records[index]
				prev.next = iterator.next
				*h.records--
				return *dev, true
			}
			prev = iterator
			iterator = iterator.next
		}
		return Record{"", nil, nil}, false
	}
}

// Remove: Remueve una clave (y valor)
// @param {string}
// @return {bool}
func (h *HashTable) Remove(key string) (Record, bool) {
	val, sizeChanged := h._remove(key)
	if sizeChanged {
		h._checkLoadFactorAndUpdate()
	}

	return val, sizeChanged
}

// getLoadFactor: calcula el loadfactor
// Calculado como, número de registros almacenados / ancho Remove slice subyacente
// @return {float64}
func (h *HashTable) _getLoadFactor() float64 {
	return float64(*h.records) / float64(len(h.table.records))
}

// checkLoadFactorAndUpdate: si 0.25 > loadfactor o 0.75 < loadfactor,
// actualiza el slice subyacente para tener loadfactor cercano a 0.5
func (h *HashTable) _checkLoadFactorAndUpdate() {
	if *h.records == 0 {
		return
	} else {
		loadFactor := h._getLoadFactor()
		if loadFactor < _MINLOADFACTOR {
			hash := _createHashTable(len(h.table.records) / 2)
			for _, record := range h.table.records {
				for record != nil {
					hash._add(record.key, record.value)
					record = record.next
				}
			}
			h.table = hash.table
			h.table._error_ = core.NewWarning("UserWarning", "** Loadfactor por debajo del límite, aumentando tamaño **")
		} else if loadFactor > _MAXLOADFACTOR {
			hash := _createHashTable(*h.records * 2)
			for _, record := range h.table.records {
				for record != nil {
					hash._add(record.key, record.value)
					record = record.next
				}
			}
			h.table = hash.table
			h.table._error_ = core.NewWarning("UserWarning", "** Loadfactor superior al límite, reduciendo tamaño **")
		}
	}
}

// Keys: Nos da una lista de claves hash
// @return {[]string}
func (h *HashTable) Keys() []string {
	var ret []string

	for _, node := range h.table.records {
		for node != nil {
			ret = append(ret, node.key)
			node = node.next
		}
	}

	return ret
}

// Values: Nos da una lista de valores hash
// @return {Values}
func (h *HashTable) Values() []Value {
	var ret []Value

	for _, node := range h.table.records {
		for node != nil {
			ret = append(ret, node.value)
			node = node.next
		}
	}

	return ret
}

// Size: Nos dice el número de elementos que tiene el hash
// @return {int}
func (h *HashTable) Size() int {
	return *h.records
}

// Capacity: Nos dice la capacidad del HashTable
// @return {int}
func (h *HashTable) Capacity() int {
	return len(h.table.records)
}

// LoadFactor: Nos dice el loadfactor del HashTaable
// @return {float64}
func (h *HashTable) LoadFactor() float64 {
	return h._getLoadFactor()
}

// GetHash: Nos dice el número de hash de la clave 'key'
// @param {string}
// @param {int}
func (h *HashTable) GetHash(key string) int {
	return _hash(key)
}

// ContainsKey: Nos dice si existe la clave 'key'
// @param {string}
// @return {bool}
func (h *HashTable) ContainsKey(key string) bool {
	for _, node := range h.table.records {
		for node != nil {
			if key == node.key {
				return true
			}
			node = node.next
		}
	}

	return false
}

// ContainsValue: Nos dice si existe el valor 'value'
// @param {Value}
// @return {bool}
func (h *HashTable) ContainsValue(value Value) bool {
	for _, node := range h.table.records {
		for node != nil {
			if value == node.value {
				return true
			}
			node = node.next
		}
	}

	return false
}

// Clone: Nos da una copia HashTable
// @return HashTable
func (h *HashTable) Clone() IHashTable {
	copy := NewHashTable()

	for _, node := range h.table.records {
		for node != nil {
			copy.Add(node.key, node.value)
			node = node.next
		}
	}

	return copy
}
