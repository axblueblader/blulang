package main

import (
	"unicode"
)

type TokenType int

type Token struct {
	name  TokenType
	value string
}

const (
	TkBinaryOperator = iota
	TkNumber
	TkString
	TkIdentifier
	TkDeclareVar
	TkDeclareFunc
	TkIf
	TkElse
	TkWhile
	TkComma
	TkDot
	TkNot
	TkOpenRound
	TkCloseRound
	TkOpenCurly
	TkCloseCurly
	TKOpenSquare
	TkCloseSquare
)

var Keywords = map[string]TokenType{
	"let":   TkDeclareVar,
	"cho":   TkDeclareVar,
	"fn":    TkDeclareFunc,
	"hàm":   TkDeclareFunc,
	"if":    TkIf,
	"nếu":   TkIf,
	"else":  TkElse,
	"hay":   TkElse,
	"while": TkWhile,
	"khi":   TkWhile,
}

func NewToken(name TokenType, value string) Token {
	return Token{
		name:  name,
		value: value,
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isAlpha(ch rune) bool {
	return unicode.IsLetter(rune(ch))
	//return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isIgnored(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\t'
}

func isOneCharBinaryOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '=' || ch == '>' || ch == '<'
}

func isTwoCharBinaryOperator(first rune, second rune) bool {
	if first == '=' && second == '=' {
		return true
	}
	if first == '!' && second == '=' {
		return true
	}
	if first == '>' && second == '=' {
		return true
	}
	if first == '<' && second == '=' {
		return true
	}
	if first == '&' && second == '&' {
		return true
	}
	if first == '|' && second == '|' {
		return true
	}
	return false
}

func Tokenize(source string) []Token {
	var tokens []Token
	runeArr := []rune(source)
	for i := 0; i < len(runeArr); i++ {
		ch := runeArr[i]
		if isIgnored(ch) {
			continue
		}

		if i+1 < len(runeArr) && isTwoCharBinaryOperator(ch, runeArr[i+1]) {
			tokens = append(tokens, NewToken(TkBinaryOperator, string(ch)+string(runeArr[i+1])))
			i++
			continue
		}

		if ch == ';' {
			// comments
			for i+1 < len(runeArr) && runeArr[i+1] != '\n' {
				i++
			}
		}

		if ch == '[' {
			tokens = append(tokens, NewToken(TKOpenSquare, string(ch)))
			continue
		}
		if ch == ']' {
			tokens = append(tokens, NewToken(TkCloseSquare, string(ch)))
			continue
		}
		if ch == '(' {
			tokens = append(tokens, NewToken(TkOpenRound, string(ch)))
			continue
		}
		if ch == ')' {
			tokens = append(tokens, NewToken(TkCloseRound, string(ch)))
			continue
		}
		if ch == '{' {
			tokens = append(tokens, NewToken(TkOpenCurly, string(ch)))
			continue
		}
		if ch == '}' {
			tokens = append(tokens, NewToken(TkCloseCurly, string(ch)))
			continue
		}
		if ch == ',' {
			tokens = append(tokens, NewToken(TkComma, string(ch)))
			continue
		}
		if ch == '.' {
			tokens = append(tokens, NewToken(TkDot, string(ch)))
			continue
		}
		if ch == '"' {
			str := ""
			for i+1 < len(runeArr) && (runeArr[i+1] != '"' || (runeArr[i] == '\\' && runeArr[i+1] == '"')) {
				i++
				// skip escape character
				if runeArr[i] != '\\' {
					str += string(runeArr[i])
				}
			}
			i++ // skip closing quote
			tokens = append(tokens, NewToken(TkString, str))
			continue
		}
		if ch == '!' {
			tokens = append(tokens, NewToken(TkNot, string(ch)))
			continue
		}

		if isOneCharBinaryOperator(ch) {
			tokens = append(tokens, NewToken(TkBinaryOperator, string(ch)))
			continue
		}

		if isDigit(ch) {
			numStr := string(ch)
			for i+1 < len(runeArr) && isDigit(runeArr[i+1]) {
				i++
				numStr = numStr + string(runeArr[i])
			}
			tokens = append(tokens, NewToken(TkNumber, numStr))
			continue
		}

		if isAlpha(ch) {
			word := string(ch)
			for i+1 < len(runeArr) && (isAlpha(runeArr[i+1]) || isDigit(runeArr[i+1])) {
				i++
				word = word + string(runeArr[i])
			}
			keywordType, found := Keywords[word]
			if found {
				tokens = append(tokens, NewToken(keywordType, word))
			} else {
				tokens = append(tokens, NewToken(TkIdentifier, word))
			}
			continue
		}
	}
	return tokens
}
