package main

type wsMsg struct {
	Type    string
	Content interface{}
}

type story struct {
	Source,
	Icon,
	Link,
	Text string
	//Pictures []string
}

type alert struct {
	Text string
}

func storyMsg(s *story) *wsMsg {
	return &wsMsg{
		"story",
		s,
	}
}

func alertMsg(a *alert) *wsMsg {
	return &wsMsg{
		"alert",
		a,
	}
}
