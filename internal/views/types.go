package views

type Data map[string]any
type Errors map[string]string

func (d Data) GetIntFromData(key string) (int, bool) {
	i, ok := d[key].(int)
	if !ok {
		return 0, ok
	}
	return i, ok
}

func (d Data) GetErrorsFromData() (Errors, bool) {
	errors, ok := d["errors"].(Errors)
	if !ok {
		return nil, ok
	}
	return errors, ok
}
