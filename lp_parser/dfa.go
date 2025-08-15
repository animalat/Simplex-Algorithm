package parser

import "fmt"

var allTokens = []TokenType{
	TokenLet,
	TokenSubjectTo,
	TokenMin,
	TokenMax,
	TokenId,
	TokenNumber,
	TokenDecimal,
	TokenSemiColon,
	TokenEqual,
	TokenLessEqual,
	TokenGreaterEqual,
	TokenPlus,
	TokenMinus,
	TokenAsterisk,
	TokenDivide,
	TokenLParen,
	TokenRParen,
}

func (dfa *DFA) initAlphabet() {
	for letter := 'a'; letter <= 'z'; letter++ {
		dfa.AlphabetSymbols[letter] = true
	}

	for letter := 'A'; letter <= 'Z'; letter++ {
		dfa.AlphabetSymbols[letter] = true
	}

	for number := '0'; number <= '9'; number++ {
		dfa.AlphabetSymbols[number] = true
	}

	const operators = "=<>+-*/()"
	for _, operator := range operators {
		dfa.AlphabetSymbols[operator] = true
	}

	// we use '.' in "s.t."
	dfa.AlphabetSymbols['.'] = true
}

func (dfa *DFA) initStates() {
	for _, state := range allTokens {
		dfa.States[string(state)] = true
	}

	dfa.States["<"] = true
	dfa.States[">"] = true
}

func (dfa *DFA) initFinalStates() {
	for _, state := range allTokens {
		dfa.FinalStates[string(state)] = state
	}
}

func addFallbackToId(dfa *DFA, fromState string, toState string, exclude rune) {
	for letter := 'a'; letter <= 'z'; letter++ {
		if letter == exclude {
			continue
		}

		key := TransitionKey{fromState, letter}
		if _, ok := dfa.Transitions[key]; ok {
			continue
		}

		dfa.Transitions[key] = toState
	}
}

func addWordTransitions(dfa *DFA, keyword string, token TokenType) {
	curr := StartingState
	runes := []rune(keyword)

	for i, ch := range runes {
		isLast := i+1 == len(runes)

		var next string

		if isLast {
			next = string(token)
		} else {
			next = string(runes[:i+1])
		}

		dfa.States[next] = true
		dfa.Transitions[TransitionKey{curr, ch}] = next

		// add fallbacks (e.g. less -> ID if another follows less)
		exclude := rune(0)
		if !isLast {
			exclude = runes[i+1]
		}
		addFallbackToId(dfa, next, string(TokenId), exclude)

		curr = next
	}
}

func (dfa *DFA) initTransitions() {
	// ID transitions
	for letter := 'a'; letter <= 'z'; letter++ {
		// start -> ID and ID -> ID with rune
		dfa.Transitions[TransitionKey{string(TokenId), letter}] = string(TokenId)

		// skip 'l' and 's' as we need those for LET and S.T.
		if letter == 'l' || letter == 's' {
			continue
		}
		dfa.Transitions[TransitionKey{StartingState, letter}] = string(TokenId)
	}

	for letter := '0'; letter <= '9'; letter++ {
		// ID -> ID with non-starting character as a number
		dfa.Transitions[TransitionKey{string(TokenId), letter}] = string(TokenId)
	}

	// NUMBER/DECIMAL transitions
	for number := '0'; number <= '9'; number++ {
		// start -> NUMBER and NUMBER -> NUMBER with number
		dfa.Transitions[TransitionKey{StartingState, number}] = string(TokenNumber)
		dfa.Transitions[TransitionKey{string(TokenNumber), number}] = string(TokenNumber)
		dfa.Transitions[TransitionKey{string(TokenDecimal), number}] = string(TokenDecimal)
	}

	addWordTransitions(dfa, "let", TokenLet)
	addWordTransitions(dfa, "s.t.", TokenSubjectTo)
	addWordTransitions(dfa, "min", TokenMin)
	addWordTransitions(dfa, "max", TokenMax)

	// OPERATOR and SYMBOL transitions
	dfa.Transitions[TransitionKey{StartingState, ';'}] = string(TokenSemiColon)

	dfa.Transitions[TransitionKey{StartingState, '<'}] = "<"
	dfa.Transitions[TransitionKey{StartingState, '>'}] = ">"
	dfa.Transitions[TransitionKey{StartingState, '='}] = string(TokenEqual)
	dfa.Transitions[TransitionKey{"<", '='}] = string(TokenLessEqual)
	dfa.Transitions[TransitionKey{">", '='}] = string(TokenGreaterEqual)

	dfa.Transitions[TransitionKey{StartingState, '+'}] = string(TokenPlus)
	dfa.Transitions[TransitionKey{StartingState, '-'}] = string(TokenMinus)
	dfa.Transitions[TransitionKey{StartingState, '*'}] = string(TokenAsterisk)
	dfa.Transitions[TransitionKey{StartingState, '/'}] = string(TokenDivide)
	dfa.Transitions[TransitionKey{StartingState, '('}] = string(TokenLParen)
	dfa.Transitions[TransitionKey{StartingState, ')'}] = string(TokenRParen)
}

func NewDFA() *DFA {
	dfa := &DFA{
		AlphabetSymbols: make(map[rune]bool),
		States:          make(map[string]bool),
		FinalStates:     make(map[string]TokenType),
		Transitions:     make(map[TransitionKey]string),
		StartState:      StartingState,
	}
	dfa.initAlphabet()
	dfa.initStates()
	dfa.initFinalStates()
	dfa.initTransitions()
	return dfa
}

// Token is the state that has been read, int is how many characters were consumed
func (dfa *DFA) Run(input []rune, lineNumber int) (Token, int, error) {
	curState := dfa.StartState
	const startPos = -1
	prevFinalPos := startPos
	var prevFinalType TokenType

	for i := 0; i < len(input); i++ {
		key := TransitionKey{curState, input[i]}
		nextState, ok := dfa.Transitions[key]
		if !ok {
			break
		}
		curState = nextState
		if finalType, ok := dfa.FinalStates[curState]; ok {
			prevFinalPos = i
			prevFinalType = finalType
		}
	}

	if prevFinalPos != startPos {
		return Token{Type: prevFinalType, Value: string(input[:prevFinalPos+1]), Line: lineNumber}, prevFinalPos + 1, nil
	}

	return Token{}, 0, fmt.Errorf("no valid token recognized from input %q", string(input))
}
