package views

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
)

type Alert struct {
	Level   string
	Message string
}

// Top level structure that views expect
type Data struct {
	Alert *Alert
	Yield interface{}
}
