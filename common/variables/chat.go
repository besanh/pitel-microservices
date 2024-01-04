package variables

var (
	CONNECTION_TYPE []string          = []string{"facebook", "zalo", "telegram"}
	EVENT_NAME      map[string]string = map[string]string{
		"user_send_text":        "user_send_text",
		"user_send_image":       "user_send_image",
		"user_send_link":        "user_send_link",
		"user_send_sticker":     "user_send_sticker",
		"user_send_gif":         "user_send_gif",
		"user_send_audio":       "user_send_audio",
		"user_send_video":       "user_send_video",
		"user_send_file":        "user_send_file",
		"user_received_message": "user_received_message",
		"user_seen_message":     "user_seen_message",
	}
	EVENT_READ_MESSAGE []string = []string{
		"user_seen_message",
	}
	DIRECTION map[string]string = map[string]string{
		"send":    "send",
		"receive": "receive",
	}
	ATTACHMENT_TYPE map[string]string = map[string]string{
		"text":    "text",
		"image":   "image",
		"audio":   "audio",
		"video":   "video",
		"file":    "file",
		"link":    "link",
		"sticker": "sticker",
		"gif":     "gif",
	}
)
