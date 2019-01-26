package test

type testClipboard struct {
}

func (c *testClipboard) Content() string {
	return "clipboard content"
}

func (c *testClipboard) SetContent(content string) {
}
