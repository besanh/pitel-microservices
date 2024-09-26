package constant

import "time"

const (
	// Version
	VERSION      string = "1"
	VERSION_NAME string = "Chat(BSS)"
)

var (
	SOURCE_AUTH = []string{"aaa", "authentication"}
)

const OBJECT_EXPIRE_TIME time.Duration = time.Second * 60 * 60 * 24 * 7
