package main

import "sort"

var DeviceNames = []string{
	"MX Master 3",
	"MX Master 3S",
	"MX Master 4",
	"MX Master 4S",

	"Wireless Mouse MX Master 3",
	"Wireless Mouse MX Master 3S",
	"Wireless Mouse MX Master 4",
	"Wireless Mouse MX Master 4S",
	"Wireless Mouse MX Master",
}

type Lang string

const (
	LangEN Lang = "en"
	LangPT Lang = "pt"
)

var currentLang Lang = LangEN

type ButtonInfo struct {
	CID  uint32
	Name string
}

var MXMasterButtons = []ButtonInfo{
	{CID: 0x0050, Name: "left"},
	{CID: 0x0051, Name: "right"},
	{CID: 0x0052, Name: "middle"},
	{CID: 0x0053, Name: "back"},
	{CID: 0x0056, Name: "forward"},
	{CID: 0x00c3, Name: "gesture"},
	{CID: 0x00c4, Name: "smartshift_toggle"},
}

var ButtonNameEN = map[string]string{
	"left":              "Left Button",
	"right":             "Right Button",
	"middle":            "Middle Button",
	"back":              "Back Button",
	"forward":           "Forward Button",
	"gesture":           "Gesture Button (Thumb)",
	"smartshift_toggle": "SmartShift Toggle",
	"thumbwheel_left":   "Thumbwheel Left",
	"thumbwheel_right":  "Thumbwheel Right",
	"thumbwheel_tap":    "Thumbwheel Tap",
}

var ButtonNamePT = map[string]string{
	"left":              "Botão Esquerdo",
	"right":             "Botão Direito",
	"middle":            "Botão do Meio",
	"back":              "Botão Voltar",
	"forward":           "Botão Avançar",
	"gesture":           "Botão de Gesto (Polegar)",
	"smartshift_toggle": "Alternar SmartShift",
	"thumbwheel_left":   "Roda Lateral Esquerda",
	"thumbwheel_right":  "Roda Lateral Direita",
	"thumbwheel_tap":    "Toque na Roda Lateral",
}

var ActionTypesEN = map[string]string{
	"None":              "None (disabled)",
	"Keypress":          "Key Press",
	"Gestures":          "Gestures",
	"ToggleSmartShift":  "Toggle SmartShift",
	"ToggleHiresScroll": "Toggle HiRes Scroll",
	"CycleDPI":          "Cycle DPI",
	"ChangeDPI":         "Change DPI",
	"ChangeHost":        "Change Host",
}

var ActionTypesPT = map[string]string{
	"None":              "Nenhum (desativado)",
	"Keypress":          "Tecla",
	"Gestures":          "Gestos",
	"ToggleSmartShift":  "Alternar SmartShift",
	"ToggleHiresScroll": "Alternar Scroll HD",
	"CycleDPI":          "Alternar DPI",
	"ChangeDPI":         "Mudar DPI",
	"ChangeHost":        "Mudar Host",
}

var GestureDirectionsEN = map[string]string{
	"Up": "Up", "Down": "Down", "Left": "Left", "Right": "Right", "None": "None",
}

var GestureDirectionsPT = map[string]string{
	"Up": "Cima", "Down": "Baixo", "Left": "Esquerda", "Right": "Direita", "None": "Nenhum",
}

var GestureModesEN = map[string]string{
	"OnRelease":   "On Release",
	"OnInterval":  "On Interval",
	"OnThreshold": "On Threshold",
	"Axis":        "Axis",
	"NoPress":     "No Press",
}

var GestureModesPT = map[string]string{
	"OnRelease":   "Ao Soltar",
	"OnInterval":  "No Intervalo",
	"OnThreshold": "No Limiar",
	"Axis":        "Eixo",
	"NoPress":     "Sem Pressão",
}

var ActionTypes = func(lang Lang) map[string]string {
	if lang == LangPT {
		return ActionTypesPT
	}
	return ActionTypesEN
}

var GestureDirections = func(lang Lang) map[string]string {
	if lang == LangPT {
		return GestureDirectionsPT
	}
	return GestureDirectionsEN
}

var GestureModes = func(lang Lang) map[string]string {
	if lang == LangPT {
		return GestureModesPT
	}
	return GestureModesEN
}

const DeviceNameInfoEN = "Check the device name in the terminal with \"logid -v\". If it does not appear, turn the mouse off and on while logid is running."
const DeviceNameInfoPT = "Verificar o nome do dispositivo no terminal com o comando \"logid -v\". Se não conseguir ver, com o logid rodando, desligue e ligue o mouse."

func DeviceNameInfo(lang Lang) string {
	if lang == LangPT {
		return DeviceNameInfoPT
	}
	return DeviceNameInfoEN
}

func ButtonName(id string, lang Lang) string {
	if lang == LangPT {
		if v, ok := ButtonNamePT[id]; ok {
			return v
		}
	}
	if v, ok := ButtonNameEN[id]; ok {
		return v
	}
	return id
}

func Translate(pt, en string, lang Lang) string {
	if lang == LangPT {
		return pt
	}
	return en
}

var KeyCodes = map[string]uint16{
	"KEY_ESC": 1, "KEY_1": 2, "KEY_2": 3, "KEY_3": 4, "KEY_4": 5,
	"KEY_5": 6, "KEY_6": 7, "KEY_7": 8, "KEY_8": 9, "KEY_9": 10,
	"KEY_0": 11, "KEY_MINUS": 12, "KEY_EQUAL": 13, "KEY_BACKSPACE": 14,
	"KEY_TAB": 15, "KEY_Q": 16, "KEY_W": 17, "KEY_E": 18, "KEY_R": 19,
	"KEY_T": 20, "KEY_Y": 21, "KEY_U": 22, "KEY_I": 23, "KEY_O": 24,
	"KEY_P": 25, "KEY_BRACELEFT": 26, "KEY_BRACERIGHT": 27,
	"KEY_ENTER": 28, "KEY_LEFTCTRL": 29, "KEY_A": 30,
	"KEY_S": 31, "KEY_D": 32, "KEY_F": 33, "KEY_G": 34, "KEY_H": 35,
	"KEY_J": 36, "KEY_K": 37, "KEY_L": 38, "KEY_SEMICOLON": 39,
	"KEY_APOSTROPHE": 40, "KEY_GRAVE": 41, "KEY_LEFTSHIFT": 42,
	"KEY_BACKSLASH": 43, "KEY_Z": 44, "KEY_X": 45, "KEY_C": 46,
	"KEY_V": 47, "KEY_B": 48, "KEY_N": 49, "KEY_M": 50,
	"KEY_COMMA": 51, "KEY_DOT": 52, "KEY_SLASH": 53,
	"KEY_RIGHTSHIFT": 54, "KEY_KPASTERISK": 55, "KEY_LEFTALT": 56,
	"KEY_SPACE": 57, "KEY_CAPSLOCK": 58, "KEY_F1": 59, "KEY_F2": 60,
	"KEY_F3": 61, "KEY_F4": 62, "KEY_F5": 63, "KEY_F6": 64,
	"KEY_F7": 65, "KEY_F8": 66, "KEY_F9": 67, "KEY_F10": 68,
	"KEY_F11": 87, "KEY_F12": 88, "KEY_HOME": 102, "KEY_UP": 103,
	"KEY_PAGEUP": 104, "KEY_LEFT": 105, "KEY_RIGHT": 106,
	"KEY_END": 107, "KEY_DOWN": 108, "KEY_PAGEDOWN": 109,
	"KEY_INSERT": 110, "KEY_DELETE": 111, "KEY_MUTE": 113,
	"KEY_VOLUMEDOWN": 114, "KEY_VOLUMEUP": 115,
	"KEY_LEFTMETA": 125, "KEY_RIGHTMETA": 126,
	"KEY_RIGHTALT": 100, "KEY_RIGHTCTRL": 97,
	"KEY_PAUSE": 119, "KEY_SCROLLLOCK": 70, "KEY_NUMLOCK": 69,
	"KEY_KP0": 82, "KEY_KP1": 79, "KEY_KP2": 80, "KEY_KP3": 81,
	"KEY_KP4": 75, "KEY_KP5": 76, "KEY_KP6": 77, "KEY_KP7": 71,
	"KEY_KP8": 72, "KEY_KP9": 73, "KEY_KPDOT": 83, "KEY_KPSLASH": 98,
	"KEY_KPMINUS": 74, "KEY_KPPLUS": 78, "KEY_KPENTER": 96,
	"KEY_PREVIOUSSONG": 165, "KEY_NEXTSONG": 163,
	"KEY_PLAYPAUSE": 164, "KEY_STOP": 166,
	"KEY_EMAIL": 172, "KEY_CALC": 140, "KEY_COMPUTER": 157,
	"KEY_BRIGHTNESSDOWN": 224, "KEY_BRIGHTNESSUP": 225,
	"KEY_SLEEP": 142, "KEY_WAKEUP": 143,
	"KEY_BACK": 158, "KEY_FORWARD": 159, "KEY_REFRESH": 173,
	"KEY_HOMEPAGE": 172, "KEY_SEARCH": 217, "KEY_BOOKMARKS": 220,
	"KEY_ZOOMIN": 418, "KEY_ZOOMOUT": 419,
	"KEY_PRINT": 210, "KEY_SCREENLOCK": 152,
}

var CommonKeyNames []string

func init() {
	for k := range KeyCodes {
		CommonKeyNames = append(CommonKeyNames, k)
	}
	sort.Slice(CommonKeyNames, func(i, j int) bool {
		return KeyCodes[CommonKeyNames[i]] < KeyCodes[CommonKeyNames[j]]
	})
}

var KeyNamesEN = map[string]string{
	"KEY_ESC":            "Esc",
	"KEY_TAB":            "Tab",
	"KEY_CAPSLOCK":       "Caps Lock",
	"KEY_LEFTSHIFT":      "Left Shift",
	"KEY_RIGHTSHIFT":     "Right Shift",
	"KEY_LEFTCTRL":       "Left Ctrl",
	"KEY_RIGHTCTRL":      "Right Ctrl",
	"KEY_LEFTALT":        "Left Alt",
	"KEY_RIGHTALT":       "Right Alt",
	"KEY_LEFTMETA":       "Left Meta (Win/Command)",
	"KEY_RIGHTMETA":      "Right Meta (Win/Command)",
	"KEY_SPACE":          "Space",
	"KEY_ENTER":          "Enter",
	"KEY_BACKSPACE":      "Backspace",
	"KEY_DELETE":         "Delete",
	"KEY_INSERT":         "Insert",
	"KEY_HOME":           "Home",
	"KEY_END":            "End",
	"KEY_PAGEUP":         "Page Up",
	"KEY_PAGEDOWN":       "Page Down",
	"KEY_UP":             "Arrow Up",
	"KEY_DOWN":           "Arrow Down",
	"KEY_LEFT":           "Arrow Left",
	"KEY_RIGHT":          "Arrow Right",
	"KEY_PRINT":          "Print Screen",
	"KEY_PAUSE":          "Pause/Break",
	"KEY_SCROLLLOCK":     "Scroll Lock",
	"KEY_NUMLOCK":        "Num Lock",
	"KEY_BACKSLASH":      "Backslash",
	"KEY_SLASH":          "Slash",
	"KEY_DOT":            "Period",
	"KEY_COMMA":          "Comma",
	"KEY_SEMICOLON":      "Semicolon",
	"KEY_APOSTROPHE":     "Apostrophe",
	"KEY_GRAVE":          "Grave Accent",
	"KEY_MINUS":          "Minus",
	"KEY_EQUAL":          "Equal",
	"KEY_BRACELEFT":      "Left Brace",
	"KEY_BRACERIGHT":     "Right Brace",
	"KEY_KPASTERISK":     "Numpad *",
	"KEY_KPMINUS":        "Numpad -",
	"KEY_KPPLUS":         "Numpad +",
	"KEY_KPDOT":          "Numpad .",
	"KEY_KPSLASH":        "Numpad /",
	"KEY_KPENTER":        "Numpad Enter",
	"KEY_MUTE":           "Mute",
	"KEY_VOLUMEDOWN":     "Volume Down",
	"KEY_VOLUMEUP":       "Volume Up",
	"KEY_PLAYPAUSE":      "Play/Pause",
	"KEY_STOP":           "Stop",
	"KEY_PREVIOUSSONG":   "Previous Track",
	"KEY_NEXTSONG":       "Next Track",
	"KEY_EMAIL":          "Email",
	"KEY_CALC":           "Calculator",
	"KEY_COMPUTER":       "Computer",
	"KEY_BRIGHTNESSDOWN": "Brightness Down",
	"KEY_BRIGHTNESSUP":   "Brightness Up",
	"KEY_SLEEP":          "Sleep",
	"KEY_WAKEUP":         "Wake Up",
	"KEY_BACK":           "Browser Back",
	"KEY_FORWARD":        "Browser Forward",
	"KEY_REFRESH":        "Refresh",
	"KEY_HOMEPAGE":       "Homepage",
	"KEY_SEARCH":         "Search",
	"KEY_BOOKMARKS":      "Bookmarks",
	"KEY_ZOOMIN":         "Zoom In",
	"KEY_ZOOMOUT":        "Zoom Out",
	"KEY_SCREENLOCK":     "Screen Lock",
	"KEY_0":              "0",
	"KEY_1":              "1",
	"KEY_2":              "2",
	"KEY_3":              "3",
	"KEY_4":              "4",
	"KEY_5":              "5",
	"KEY_6":              "6",
	"KEY_7":              "7",
	"KEY_8":              "8",
	"KEY_9":              "9",
	"KEY_A":              "A",
	"KEY_B":              "B",
	"KEY_C":              "C",
	"KEY_D":              "D",
	"KEY_E":              "E",
	"KEY_F":              "F",
	"KEY_G":              "G",
	"KEY_H":              "H",
	"KEY_I":              "I",
	"KEY_J":              "J",
	"KEY_K":              "K",
	"KEY_L":              "L",
	"KEY_M":              "M",
	"KEY_N":              "N",
	"KEY_O":              "O",
	"KEY_P":              "P",
	"KEY_Q":              "Q",
	"KEY_R":              "R",
	"KEY_S":              "S",
	"KEY_T":              "T",
	"KEY_U":              "U",
	"KEY_V":              "V",
	"KEY_W":              "W",
	"KEY_X":              "X",
	"KEY_Y":              "Y",
	"KEY_Z":              "Z",
	"KEY_F1":             "F1",
	"KEY_F2":             "F2",
	"KEY_F3":             "F3",
	"KEY_F4":             "F4",
	"KEY_F5":             "F5",
	"KEY_F6":             "F6",
	"KEY_F7":             "F7",
	"KEY_F8":             "F8",
	"KEY_F9":             "F9",
	"KEY_F10":            "F10",
	"KEY_F11":            "F11",
	"KEY_F12":            "F12",
	"KEY_KP0":            "Numpad 0",
	"KEY_KP1":            "Numpad 1",
	"KEY_KP2":            "Numpad 2",
	"KEY_KP3":            "Numpad 3",
	"KEY_KP4":            "Numpad 4",
	"KEY_KP5":            "Numpad 5",
	"KEY_KP6":            "Numpad 6",
	"KEY_KP7":            "Numpad 7",
	"KEY_KP8":            "Numpad 8",
	"KEY_KP9":            "Numpad 9",
}

var KeyNamesPT = map[string]string{
	"KEY_ESC":            "Esc",
	"KEY_TAB":            "Tab",
	"KEY_CAPSLOCK":       "Caps Lock",
	"KEY_LEFTSHIFT":      "Shift Esquerdo",
	"KEY_RIGHTSHIFT":     "Shift Direito",
	"KEY_LEFTCTRL":       "Ctrl Esquerdo",
	"KEY_RIGHTCTRL":      "Ctrl Direito",
	"KEY_LEFTALT":        "Alt Esquerdo",
	"KEY_RIGHTALT":       "Alt Direito",
	"KEY_LEFTMETA":       "Meta Esquerdo (Win/Command)",
	"KEY_RIGHTMETA":      "Meta Direito (Win/Command)",
	"KEY_SPACE":          "Espaço",
	"KEY_ENTER":          "Enter",
	"KEY_BACKSPACE":      "Backspace",
	"KEY_DELETE":         "Delete",
	"KEY_INSERT":         "Insert",
	"KEY_HOME":           "Home",
	"KEY_END":            "End",
	"KEY_PAGEUP":         "Page Up",
	"KEY_PAGEDOWN":       "Page Down",
	"KEY_UP":             "Seta Cima",
	"KEY_DOWN":           "Seta Baixo",
	"KEY_LEFT":           "Seta Esquerda",
	"KEY_RIGHT":          "Seta Direita",
	"KEY_PRINT":          "Print Screen",
	"KEY_PAUSE":          "Pause/Break",
	"KEY_SCROLLLOCK":     "Scroll Lock",
	"KEY_NUMLOCK":        "Num Lock",
	"KEY_BACKSLASH":      "Barra Invertida",
	"KEY_SLASH":          "Barra",
	"KEY_DOT":            "Ponto",
	"KEY_COMMA":          "Vírgula",
	"KEY_SEMICOLON":      "Ponto e Vírgula",
	"KEY_APOSTROPHE":     "Apóstrofo",
	"KEY_GRAVE":          "Acento Grave",
	"KEY_MINUS":          "Menos",
	"KEY_EQUAL":          "Igual",
	"KEY_BRACELEFT":      "Chave Esquerda",
	"KEY_BRACERIGHT":     "Chave Direita",
	"KEY_KPASTERISK":     "Teclado Numérico *",
	"KEY_KPMINUS":        "Teclado Numérico -",
	"KEY_KPPLUS":         "Teclado Numérico +",
	"KEY_KPDOT":          "Teclado Numérico .",
	"KEY_KPSLASH":        "Teclado Numérico /",
	"KEY_KPENTER":        "Teclado Numérico Enter",
	"KEY_MUTE":           "Mudo",
	"KEY_VOLUMEDOWN":     "Volume Baixo",
	"KEY_VOLUMEUP":       "Volume Alto",
	"KEY_PLAYPAUSE":      "Tocar/Pausar",
	"KEY_STOP":           "Parar",
	"KEY_PREVIOUSSONG":   "Faixa Anterior",
	"KEY_NEXTSONG":       "Próxima Faixa",
	"KEY_EMAIL":          "E-mail",
	"KEY_CALC":           "Calculadora",
	"KEY_COMPUTER":       "Computador",
	"KEY_BRIGHTNESSDOWN": "Reduzir Brilho",
	"KEY_BRIGHTNESSUP":   "Aumentar Brilho",
	"KEY_SLEEP":          "Suspender",
	"KEY_WAKEUP":         "Ativar",
	"KEY_BACK":           "Voltar (Navegador)",
	"KEY_FORWARD":        "Avançar (Navegador)",
	"KEY_REFRESH":        "Atualizar",
	"KEY_HOMEPAGE":       "Página Inicial",
	"KEY_SEARCH":         "Pesquisar",
	"KEY_BOOKMARKS":      "Favoritos",
	"KEY_ZOOMIN":         "Aumentar Zoom",
	"KEY_ZOOMOUT":        "Reduzir Zoom",
	"KEY_SCREENLOCK":     "Bloquear Tela",
	"KEY_0":              "0",
	"KEY_1":              "1",
	"KEY_2":              "2",
	"KEY_3":              "3",
	"KEY_4":              "4",
	"KEY_5":              "5",
	"KEY_6":              "6",
	"KEY_7":              "7",
	"KEY_8":              "8",
	"KEY_9":              "9",
	"KEY_A":              "A",
	"KEY_B":              "B",
	"KEY_C":              "C",
	"KEY_D":              "D",
	"KEY_E":              "E",
	"KEY_F":              "F",
	"KEY_G":              "G",
	"KEY_H":              "H",
	"KEY_I":              "I",
	"KEY_J":              "J",
	"KEY_K":              "K",
	"KEY_L":              "L",
	"KEY_M":              "M",
	"KEY_N":              "N",
	"KEY_O":              "O",
	"KEY_P":              "P",
	"KEY_Q":              "Q",
	"KEY_R":              "R",
	"KEY_S":              "S",
	"KEY_T":              "T",
	"KEY_U":              "U",
	"KEY_V":              "V",
	"KEY_W":              "W",
	"KEY_X":              "X",
	"KEY_Y":              "Y",
	"KEY_Z":              "Z",
	"KEY_F1":             "F1",
	"KEY_F2":             "F2",
	"KEY_F3":             "F3",
	"KEY_F4":             "F4",
	"KEY_F5":             "F5",
	"KEY_F6":             "F6",
	"KEY_F7":             "F7",
	"KEY_F8":             "F8",
	"KEY_F9":             "F9",
	"KEY_F10":            "F10",
	"KEY_F11":            "F11",
	"KEY_F12":            "F12",
	"KEY_KP0":            "Teclado Numérico 0",
	"KEY_KP1":            "Teclado Numérico 1",
	"KEY_KP2":            "Teclado Numérico 2",
	"KEY_KP3":            "Teclado Numérico 3",
	"KEY_KP4":            "Teclado Numérico 4",
	"KEY_KP5":            "Teclado Numérico 5",
	"KEY_KP6":            "Teclado Numérico 6",
	"KEY_KP7":            "Teclado Numérico 7",
	"KEY_KP8":            "Teclado Numérico 8",
	"KEY_KP9":            "Teclado Numérico 9",
}

func KeyDisplayName(key string, lang Lang) string {
	var m map[string]string
	if lang == LangPT {
		m = KeyNamesPT
	} else {
		m = KeyNamesEN
	}
	if v, ok := m[key]; ok {
		return v
	}
	return key
}

func SortedKeyDisplayNames(lang Lang) []string {
	var m map[string]string
	if lang == LangPT {
		m = KeyNamesPT
	} else {
		m = KeyNamesEN
	}
	names := make([]string, 0, len(CommonKeyNames))
	for _, k := range CommonKeyNames {
		if v, ok := m[k]; ok {
			names = append(names, v)
		} else {
			names = append(names, k)
		}
	}
	return names
}

func KeyCodeFromDisplay(display string, lang Lang) string {
	var m map[string]string
	if lang == LangPT {
		m = KeyNamesPT
	} else {
		m = KeyNamesEN
	}
	for code, name := range m {
		if name == display {
			return code
		}
	}
	// Fallback: search in CommonKeyNames
	for _, k := range CommonKeyNames {
		if k == display {
			return k
		}
	}
	return display
}
