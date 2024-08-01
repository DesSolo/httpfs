package entities

// Hash хэш
type Hash string

// String возвращает хэш в виде строки
func (h Hash) String() string {
	return string(h)
}
