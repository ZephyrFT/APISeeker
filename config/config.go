package config

var (
	Risks = []string{"remove", "delete", "insert", "update", "logout", "logOut"}

	JsFind = []string{
		"(https{0,1}:[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"[\"'‘“`]\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
		"=\\s{0,6}[\",',’,”]{0,1}\\s{0,6}(/{0,1}[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{2,250}?[-a-zA-Z0-9（）@:%_\\+.~#?&//=]{3}[.]js)",
	}
	UrlFind = []string{
		"[\"'‘“`]\\s{0,6}(https{0,1}:[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250}?)\\s{0,6}[\"'‘“`]",
		"=\\s{0,6}(https{0,1}:[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})",
		"[\"'‘“`]\\s{0,6}([#,.]{0,2}/[-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250}?)\\s{0,6}[\"'‘“`]",
		"\"([-a-zA-Z0-9()@:%_\\+.~#?&//={}]+?[/]{1}[-a-zA-Z0-9()@:%_\\+.~#?&//={}]+?)\"",
		"href\\s{0,6}=\\s{0,6}[\"'‘“`]{0,1}\\s{0,6}([-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})|action\\s{0,6}=\\s{0,6}[\"'‘“`]{0,1}\\s{0,6}([-a-zA-Z0-9()@:%_\\+.~#?&//={}]{2,250})",
	}

	JsFiler = []string{
		"www\\.w3\\.org",
		"example\\.com",
	}
	UrlFiler = []string{
		"(?i)\\.ttf\\?|\\.mjs\\?|\\.gz\\?|\\.xls\\?|\\.go\\?|\\.mp3\\?|\\.flv\\?|\\.xlsx\\?|\\.zip\\?|\\.js\\?|\\.css\\?|\\.pptx\\?|\\.rar\\?|\\.jpeg\\?|\\.cgi\\?|\\.jxl\\?|\\.webp\\?|\\.pdf\\?|\\.docx\\?|\\.bmp\\?|\\.aspx\\?|\\.jpg\\?|\\.png\\?|.gif\\?|www\\.w3\\.org|example\\.com|\\<|\\>|\\{|\\}|\\[|\\]|\\||\\^|;|/js/|\\.src|\\.replace|\\.url|\\.att|\\.href|location\\.href|javascript:|location:|application/x-www-form-urlencoded|\\.createObject|:location|\\.path|\\*#__PURE__\\*|\\*\\$0\\*|\\n",
		"(?i).*\\.ttf$|.*\\.mjs$|.*\\.gz$|.*\\.xls$|.*\\.go$|.*\\.mp3$|.*\\.flv$|.*\\.xlsx$|.*\\.zip$|.*\\.js$|.*\\.css$|.*\\.pptx$|.*\\.rar$|.*\\.scss$|.*,$|.*\\.jpeg$|.*\\.jxl$|.*\\.webp$|.*\\.cgi$|.*\\.docx$|.*\\.pdf$|.*\\.bmp$|.*\\.aspx$|.*\\.jpg$|.*\\.png$|.*\\.gif$|.*\\.ico$|.*\\.svg$|.*\\.vue$|.*\\.ts$|.*\\.pdf$|.*\\.swf$|.*\\.doc$|.*\\.ppt$|.*\\.mp4$|.*\\.mov$",
	}

	FrontendExtensions = []string{
		".ttf", ".mjs", ".html", ".js", ".css", ".jpg", ".png", ".gif", ".svg", ".htm", ".mp3", ".flv", ".jpeg", ".cgi", ".pdf", ".doc", ".docx", ".bmp", "aspx", "mp4", ".pot", ".gz", ".xls",
	}

	FileSignatures = map[string]string{
		"mp3":  "\xFF\xFB",
		"flv":  "\x46\x4C\x56\x01",
		"xlsx": "\x50\x4B\x03\x04",
		"docx": "\x50\x4B\x03\x04",
		"pptx": "\x50\x4B\x03\x04",
		"pdf":  "%PDF-",
		"zip":  "\x50\x4B\x03\x04",
		"rar":  "\x52\x61\x72\x21",
		"jpeg": "\xFF\xD8\xFF",
		"png":  "\x89\x50\x4E\x47",
		"gif":  "GIF8",
		"bmp":  "\x42\x4D",
		"webp": "RIFF", // RIFF + WEBP 需要进一步验证
		"jxl":  "\x00\x00\x00\x0C\x4A\x58\x4C\x20",
	}

	FrontendPoint = []string{
		"/static/css",
	}

	BuiltinTypesLower = map[string]string{
		"image/avif":       "image/avif",
		"text/css":         "text/css; charset=utf-8",
		"image/gif":        "image/gif",
		"text/html":        "text/html; charset=utf-8",
		"image/jpeg":       "image/jpeg",
		"text/javascript":  "text/javascript; charset=utf-8",
		"application/json": "application/json",
		"application/pdf":  "application/pdf",
		"image/png":        "image/png",
		"image/svg+xml":    "image/svg+xml",
		"application/wasm": "application/wasm",
		"image/webp":       "image/webp",
		"text/xml":         "text/xml; charset=utf-8",
		"image/x-icon":     "image/x-icon",
		"application/l10n": "application/l10n",
	}

	ParamsInPath = []string{
		"\\/\\d+\\/",
		"\\/[0-9a-fA-F]+\\/",
	}
)
