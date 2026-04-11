package screens

// ScreenID identifies which screen is active.
type ScreenID int

const (
	ScreenHome ScreenID = iota
	ScreenFilePicker
	ScreenOperations
	ScreenSettings
	ScreenProgress
	ScreenResult
)

// NavigateMsg tells the app to switch screens.
type NavigateMsg struct {
	Screen  ScreenID
	Payload interface{}
}

// StatusMsg updates the persistent status line.
type StatusMsg struct {
	Text string
}

// BackMsg requests going back one screen.
type BackMsg struct{}
