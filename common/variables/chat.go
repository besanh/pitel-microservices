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
		"oa_connection":                   "oa_connection",
		"submit_info":                     "submit_info",
		"ask_info":                        "ask_info",
		"message_created":                 "message_created",
		"conversation_created":            "conversation_created",
		"conversation_done":               "conversation_done",
		"conversation_assigned":           "conversation_assigned",
		"conversation_unassigned":         "conversation_unassigned",
		"conversation_removed":            "conversation_removed",
		"conversation_reopen":             "conversation_reopen",
		"conversation_add_labels":         "conversation_add_labels",
		"conversation_remove_labels":      "conversation_remove_labels",
		"conversation_user_put_major":     "conversation_user_put_major",
		"conversation_user_put_following": "conversation_user_put_following",
		"conversation_note_created":       "conversation_note_created",
		"conversation_note_updated":       "conversation_note_updated",
		"conversation_note_removed":       "conversation_note_removed",
	}

	STATUS_CONVERSATION []string = []string{
		"reopen",
		"done",
	}

	CHAT_AUTO_SCRIPT_EVENT = []string{
		"keyword",
		"offline",
	}

	CHAT_SCRIPT_TYPE = []string{
		"text",
		"image",
		"file",
		"other",
	}

	CHAT_LABEL_ACTION = []string{
		"create",
		"update",
		"delete",
	}

	PERSONALIZATION_KEYWORDS = []string{
		"{{page_name}}",
		"{{customer_name}}",
	}

	PREFERENCE_EVENT map[string]string = map[string]string{
		"major":     "conversation_user_put_major",
		"following": "conversation_user_put_following",
	}

	EVENT_USER_STATUS map[string]string = map[string]string{
		"user_status_updated": "user_status_updated",
	}

	USER_STATUSES []string = []string{
		USER_STATUS_ONLINE,
		USER_STATUS_OFFLINE,
	}
)

const (
	USER_STATUS_ONLINE  string = "active"
	USER_STATUS_OFFLINE string = "inactive"
)
