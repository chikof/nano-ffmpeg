package screens

import "testing"

func TestScreenIDEnum(t *testing.T) {
	// These values travel over NavigateMsg; ordering changes would break
	// persisted or forwarded messages. Keep this guard in place.
	cases := []struct {
		name string
		got  ScreenID
		want ScreenID
	}{
		{"home", ScreenHome, 0},
		{"file picker", ScreenFilePicker, 1},
		{"operations", ScreenOperations, 2},
		{"settings", ScreenSettings, 3},
		{"progress", ScreenProgress, 4},
		{"result", ScreenResult, 5},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Fatalf("%s: got %d, want %d", c.name, c.got, c.want)
		}
	}
}

func TestScreenIDsAreUnique(t *testing.T) {
	values := []ScreenID{
		ScreenHome, ScreenFilePicker, ScreenOperations,
		ScreenSettings, ScreenProgress, ScreenResult,
	}
	seen := make(map[ScreenID]struct{}, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			t.Fatalf("duplicate screen id %d", v)
		}
		seen[v] = struct{}{}
	}
}

func TestNavigateMsg_CarriesPayload(t *testing.T) {
	msg := NavigateMsg{Screen: ScreenFilePicker, Payload: "pick-me"}
	if msg.Screen != ScreenFilePicker {
		t.Fatalf("screen: got %d, want %d", msg.Screen, ScreenFilePicker)
	}
	payload, ok := msg.Payload.(string)
	if !ok || payload != "pick-me" {
		t.Fatalf("payload lost: %+v", msg.Payload)
	}
}

func TestStatusMsg_HoldsText(t *testing.T) {
	m := StatusMsg{Text: "encoding complete"}
	if m.Text != "encoding complete" {
		t.Fatalf("StatusMsg text lost: %q", m.Text)
	}
}

func TestBackMsg_IsZeroValue(t *testing.T) {
	// BackMsg is intentionally empty; constructing it should be a no-op.
	_ = BackMsg{}
}
