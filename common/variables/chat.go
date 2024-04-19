package variables

var (
	CONNECTION_TYPE []string = []string{"facebook", "zalo", "telegram"}

	EVENT_NAME []string = []string{
		"user_send_text",
		"user_send_image",
		"user_send_link",
		"user_send_sticker",
		"user_send_gif",
		"user_send_audio",
		"user_send_video",
		"user_send_file",
		"user_received_message",
		"user_seen_message",
		"oa_connection",
		"submit_info",
	}

	EVENT_READ_MESSAGE []string = []string{
		"user_seen_message",
	}

	DIRECTION map[string]string = map[string]string{
		"send":    "send",
		"receive": "receive",
	}

	ATTACHMENT_TYPE_MAP []string = []string{
		"image",
		"audio",
		"video",
		"file",
		"link",
		"sticker",
		"gif",
	}

	ATTACHMENT_TYPE []string = []string{
		"image",
		"audio",
		"video",
		"file",
		"link",
		"sticker",
		"gif",
	}

	CHAT_ROUTING []string = []string{
		"random",
		"min_conversation",
		"round_robin_online",
	}

	EVENT_NAME_SEND_MESSAGE []string = []string{
		"text",
		"image",
		"audio",
		"video",
		"file",
		"link",
		"sticker",
		"gif",
	}

	EVENT_NAME_EXCLUDE []string = []string{
		"oa_connection",
		"submit_info",
		"ask_info",
		"seen",
		"received",
	}

	EVENT_CHAT map[string]string = map[string]string{
		"oa_connection":           "oa_connection",
		"submit_info":             "submit_info",
		"ask_info":                "ask_info",
		"message_created":         "message_created",
		"conversation_created":    "conversation_created",
		"conversation_done":       "conversation_done",
		"conversation_assigned":   "conversation_assigned",
		"conversation_unassigned": "conversation_unassigned",
		"conversation_removed":    "conversation_removed",
		"conversation_reopen":     "conversation_reopen",
	}

	STATUS_CONVERSATION []string = []string{
		"reopen",
		"done",
	}
)
