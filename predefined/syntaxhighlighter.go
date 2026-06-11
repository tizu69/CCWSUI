package predefined

import (
	"strings"

	"g.tizu.dev/CCWSUI/components"
)

type tokKind int

const (
	tokIdent tokKind = iota
	tokKeyword
	tokTypeName
	tokBuiltin
	tokString
	tokComment
	tokNumber
	tokOperator
)

type token struct {
	kind tokKind
	val  string
}

var keywords = map[string]bool{
	"function": true, "end": true, "if": true, "then": true, "else": true,
	"elseif": true, "for": true, "in": true, "do": true, "while": true,
	"repeat": true, "until": true, "local": true, "return": true, "break": true,
	"nil": true, "true": true, "false": true,
}

var typeNames = map[string]bool{
	"number": true, "string": true, "boolean": true, "table": true,
	"thread": true, "userdata": true, "lightuserdata": true,
}

var builtinNames = map[string]bool{
	"print": true, "assert": true, "error": true, "ipairs": true, "pairs": true,
	"next": true, "type": true, "tonumber": true, "tostring": true, "unpack": true,
	"select": true, "getmetatable": true, "setmetatable": true, "rawget": true,
	"rawset": true, "rawequal": true, "rawlen": true, "pcall": true, "xpcall": true,
	"load": true, "loadfile": true, "dofile": true, "collectgarbage": true,
	"require": true, "module": true, "package": true, "io": true, "os": true,
	"debug": true, "math": true, "table": true, "string": true, "coroutine": true,
	"utf8": true, "bit32": true, "jit": true, "ffi": true, "buffer": true,
	"arg": true, "_VERSION": true, "_G": true, "_ENV": true,
}

var tokenColors = map[tokKind]string{
	tokKeyword:  "#B266E5",
	tokTypeName: "#CC4C4C",
	tokBuiltin:  "#F2B233",
	tokString:   "#57A64E",
	tokComment:  "#999999",
	tokNumber:   "#F2B233",
	tokOperator: "#999999",
	tokIdent:    "#F0F0F0",
}

func isIdentStart(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isIdentPart(r rune) bool {
	return isIdentStart(r) || (r >= '0' && r <= '9')
}

func highlightLua(code string) components.Literal {
	code = strings.ReplaceAll(code, "\t", "    ")
	code = strings.ReplaceAll(code, "--[[br]]", "") // line break
	l := components.MkLiteral("")
	for _, tok := range tokenizeLua(code) {
		l = l.WithText(tok.val).WithHexColor(tokenColors[tok.kind])
	}
	return l
}

func tokenizeLua(code string) []token {
	var tokens []token
	runes := []rune(code)
	i := 0
	n := len(runes)

	for i < n {
		r := runes[i]

		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			start := i
			for i < n && (runes[i] == ' ' || runes[i] == '\t' || runes[i] == '\n' || runes[i] == '\r') {
				i++
			}
			tokens = append(tokens, token{tokIdent, string(runes[start:i])})
			continue
		}

		if r == '-' && i+1 < n && runes[i+1] == '-' {
			start := i
			i += 2
			if i < n && runes[i] == '[' {
				eq := 0
				i++
				for i < n && runes[i] == '=' {
					eq++
					i++
				}
				if i < n && runes[i] == '[' {
					i++
					closers := "]"
					for j := 0; j < eq; j++ {
						closers += "="
					}
					closers += "]"
					for i+len(closers) <= n && string(runes[i:i+len(closers)]) != closers {
						i++
					}
					if i+len(closers) <= n {
						i += len(closers)
					}
				}
			} else {
				for i < n && runes[i] != '\n' {
					i++
				}
			}
			tokens = append(tokens, token{tokComment, string(runes[start:i])})
			continue
		}

		if r == '"' || r == '\'' || r == '`' {
			quote := r
			start := i
			i++
			for i < n && runes[i] != quote {
				if runes[i] == '\\' && i+1 < n {
					i += 2
				} else {
					i++
				}
			}
			if i < n {
				i++
			}
			tokens = append(tokens, token{tokString, string(runes[start:i])})
			continue
		}

		if r == '[' && i+1 < n && runes[i+1] == '[' {
			start := i
			i += 2
			for i+1 < n && (runes[i] != ']' || runes[i+1] != ']') {
				if runes[i] == '\n' {
					i++
				} else {
					i++
				}
			}
			if i+1 < n {
				i += 2
			}
			tokens = append(tokens, token{tokString, string(runes[start:i])})
			continue
		}

		if r == '.' && i+1 < n && isDigit(runes[i+1]) {
			start := i
			i++
			for i < n && isDigit(runes[i]) {
				i++
			}
			if i < n && (runes[i] == 'e' || runes[i] == 'E') {
				i++
				if i < n && (runes[i] == '+' || runes[i] == '-') {
					i++
				}
				for i < n && isDigit(runes[i]) {
					i++
				}
			}
			tokens = append(tokens, token{tokNumber, string(runes[start:i])})
			continue
		}

		if isDigit(r) {
			start := i
			if runes[i] == '0' && i+1 < n && (runes[i+1] == 'x' || runes[i+1] == 'X') {
				i += 2
				for i < n && (isDigit(runes[i]) || (runes[i] >= 'a' && runes[i] <= 'f') || (runes[i] >= 'A' && runes[i] <= 'F')) {
					i++
				}
			} else if runes[i] == '0' && i+1 < n && (runes[i+1] == 'b' || runes[i+1] == 'B') {
				i += 2
				for i < n && (runes[i] == '0' || runes[i] == '1') {
					i++
				}
			} else {
				for i < n && isDigit(runes[i]) {
					i++
				}
				if i < n && runes[i] == '.' {
					i++
					for i < n && isDigit(runes[i]) {
						i++
					}
				}
				if i < n && (runes[i] == 'e' || runes[i] == 'E') {
					i++
					if i < n && (runes[i] == '+' || runes[i] == '-') {
						i++
					}
					for i < n && isDigit(runes[i]) {
						i++
					}
				}
			}
			tokens = append(tokens, token{tokNumber, string(runes[start:i])})
			continue
		}

		if isIdentStart(r) {
			start := i
			for i < n && isIdentPart(runes[i]) {
				i++
			}
			val := string(runes[start:i])
			switch {
			case keywords[val]:
				tokens = append(tokens, token{tokKeyword, val})
			case typeNames[val]:
				tokens = append(tokens, token{tokTypeName, val})
			case builtinNames[val]:
				tokens = append(tokens, token{tokBuiltin, val})
			default:
				tokens = append(tokens, token{tokIdent, val})
			}
			continue
		}

		start := i
		i++
		if i < n {
			two := string(runes[start]) + string(runes[i])
			if two == ".." || two == "==" || two == "~=" || two == "<=" || two == ">=" || two == "::" {
				i++
				if two == ".." && i < n && runes[i] == '.' {
					i++
				}
			}
		}
		tokens = append(tokens, token{tokOperator, string(runes[start:i])})
	}

	return tokens
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
