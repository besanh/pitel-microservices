package fingerprint

import (
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/danielgtaylor/huma/v2"
)

// The code points in the surrogate range are not valid for UTF-8.
const (
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
)

// Sanitize replaces control codes by the tofu symbol
// and invalid UTF-8 codes by the replacement character.
// Sanitize can be used to prevent log injection.
//
// Inspired from:
// - https://wikiless.org/wiki/Replacement_character#Replacement_character
// - https://graphicdesign.stackexchange.com/q/108297
func Sanitize(slice ...string) string {
	// most common case: one single string
	if len(slice) == 1 {
		return sanitize(slice[0])
	}

	// other cases: zero or multiple strings => use the slice representation
	str := strings.Join(slice, ", ")
	return "[" + sanitize(str) + "]"
}

func sanitize(str string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\t':
			return ' '
		case surrogateMin <= r && r <= surrogateMax, r > utf8.MaxRune:
			// The replacement character U+FFFD indicates an invalid UTF-8 character.
			return '�'
		case unicode.IsPrint(r):
			return r
		default: // r < 32, r == 127
			// The empty box (tofu) symbolizes the .notdef character
			// indicating a valid but not rendered character.
			return '􏿮'
		}
	}, str)
}

// SafeHeader stringifies a safe list of HTTP header values.
func SafeHeader(r *http.Request, header string) string {
	values := r.Header.Values(header)

	if len(values) == 0 {
		return ""
	}

	if len(values) == 1 {
		return Sanitize(values[0])
	}

	str := "["
	for i := range values {
		if i > 0 {
			str += " "
		}
		str += Sanitize(values[i])
	}
	str += "]"

	return str
}

func headerTxt(r *http.Request, header, key, skip string) string {
	v := SafeHeader(r, header)
	if v == skip {
		return ""
	}
	return " " + key + v
}

func headerMD(r *http.Request, header string) string {
	v := SafeHeader(r, header)
	if v == "" {
		return ""
	}
	return "\n" + "- **" + header + "**: " + v
}

func headerHuma(ctx huma.Context, header string) string {
	v := ctx.Header(header)
	if v == "" {
		return ""
	}
	return "\n" + "- **" + header + "**: " + v
}

func RequestFingerprint(r *http.Request) string {
	// double space after "in" is for padding with "out" logs
	line := " " +
		// 1. Accept-Language, the language preferred by the user.
		SafeHeader(r, "Accept-Language") + " " +
		// 2. User-Agent, name and version of the browser and OS.
		SafeHeader(r, "User-Agent") +
		// 3. R=Referer, the website from which the request originated.
		headerTxt(r, "Referer", "R=", "") +
		// 4. A=Accept, the content types the browser prefers.
		headerTxt(r, "Accept", "A=", "") +
		// 5. E=Accept-Encoding, the compression formats the browser supports.
		headerTxt(r, "Accept-Encoding", "E=", "") +
		// 6. Connection, can be empty, "keep-alive" or "close".
		headerTxt(r, "Connection", "", "") +
		// 7. Cache-Control, how the browser is caching data.
		headerTxt(r, "Cache-Control", "", "") +
		// 8. Upgrade-Insecure-Requests, the browser can upgrade from HTTP to HTTPS
		headerTxt(r, "Upgrade-Insecure-Requests", "UIR=", "1") +
		// 9. Via avoids request loops and identifies protocol capabilities
		headerTxt(r, "Via", "Via=", "") +
		// 10. Authorization and Cookie: both should not be present at the same time
		headerTxt(r, "Authorization", "", "") +
		headerTxt(r, "Cookie", "", "")

	// 11, DNT (Do Not Track) is being dropped by web browsers.
	if r.Header.Get("DNT") != "" {
		line += " DNT"
	}

	return line
}

func FingerprintMD(r *http.Request) string {
	return "\n" + "- **IP**: " + Sanitize(r.RemoteAddr) +
		headerMD(r, "Accept-Language") + // language preferred by the user
		headerMD(r, "User-Agent") + // name and version of browser and OS
		headerMD(r, "Referer") + // URL from which the request originated
		headerMD(r, "Accept") + // content types the browser prefers
		headerMD(r, "Accept-Encoding") + // compression formats the browser supports
		headerMD(r, "Connection") + // can be: empty, "keep-alive" or "close"
		headerMD(r, "Cache-Control") + // how the browser is caching data
		headerMD(r, "DNT") + // "Do Not Track" is being dropped by web standards and browsers
		headerMD(r, "Via") + // avoid request loops and identify protocol capabilities
		headerMD(r, "Cookie") // Attention: may contain confidential data
}

func FingerprintHuma(ctx huma.Context) string {
	return "\n" + "- **IP**: " + Sanitize(ctx.RemoteAddr()) +
		headerHuma(ctx, "Accept-Language") + // language preferred by the user
		headerHuma(ctx, "User-Agent") + // name and version of browser and OS
		headerHuma(ctx, "Referer") + // URL from which the request originated
		headerHuma(ctx, "Accept") + // content types the browser prefers
		headerHuma(ctx, "Accept-Encoding") + // compression formats the browser supports
		headerHuma(ctx, "Connection") + // can be: empty, "keep-alive" or "close"
		headerHuma(ctx, "Cache-Control") + // how the browser is caching data
		headerHuma(ctx, "DNT") + // "Do Not Track" is being dropped by web standards and browsers
		headerHuma(ctx, "Via") + // avoid request loops and identify protocol capabilities
		headerHuma(ctx, "Cookie") // Attention: may contain confidential data
}

func IPMethodURL(r *http.Request) string {
	// double space after "in" is for padding with "out" logs
	return "--> " + r.RemoteAddr + " " + r.Method + " " + r.RequestURI
}
