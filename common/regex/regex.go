package regex

const REGEX_URL string = `^(http|https):\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`
const REGEX_NORMAL string = `^[A-Za-z0-9_-]*$`
const REGEX_TEXT_ACCENT string = `^[[:alnum:]\p{L}\p{M}\s]+$`
const REGEX_EMAIL string = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
