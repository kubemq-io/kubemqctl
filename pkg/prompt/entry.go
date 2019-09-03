package prompt

type Entry interface {
	Ask() error
}
