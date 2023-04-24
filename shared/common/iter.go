package common

type (
	//Iterator represents iterator
	Iterator func(pair Pair) error
	//Pair represents a pair
	Pair func(key string, value interface{}) error
)

//ToMap coverts iterator to map
func (r Iterator) ToMap() (map[string]interface{}, error) 