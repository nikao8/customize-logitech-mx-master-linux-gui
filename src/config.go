package src

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Action struct {
	Type     string
	Keys     []string
	Gestures []Gesture
	DPIs     []int
	Inc      int
	Sensor   int
	Host     string
}

type Gesture struct {
	Direction      string
	Mode           string
	Threshold      int
	Interval       int
	Axis           string
	AxisMultiplier int
	Action         Action
}

type ButtonConfig struct {
	CID    uint32
	Action Action
}

type ThumbwheelConfig struct {
	Divert bool
	Invert bool
	Left   Action
	Right  Action
	Tap    Action
}

type SmartShiftConfig struct {
	On        bool
	Threshold int
}

type HiResScrollConfig struct {
	Hires  bool
	Invert bool
	Target bool
}

type Config struct {
	Name        string
	DPI         int
	SmartShift  SmartShiftConfig
	HiResScroll HiResScrollConfig
	Thumbwheel  ThumbwheelConfig
	Buttons     []ButtonConfig
}

func DefaultConfig() Config {
	return Config{
		Name: "MX Master 3S",
		DPI:  1300,
		SmartShift: SmartShiftConfig{
			On:        true,
			Threshold: 30,
		},
		HiResScroll: HiResScrollConfig{
			Hires:  false,
			Invert: false,
			Target: false,
		},
		Thumbwheel: ThumbwheelConfig{
			Divert: true,
			Invert: false,
		},
		Buttons: []ButtonConfig{
			{CID: 0x00c3, Action: Action{Type: "Gestures"}},
			{CID: 0x00c4, Action: Action{Type: "ToggleSmartShift"}},
			{CID: 0x0053, Action: Action{Type: "Keypress", Keys: []string{"KEY_BACK"}}},
			{CID: 0x0056, Action: Action{Type: "Keypress", Keys: []string{"KEY_FORWARD"}}},
			{CID: 0x0052, Action: Action{Type: "Keypress", Keys: []string{"KEY_ENTER"}}},
			{CID: 0x0050, Action: Action{Type: "None"}},
			{CID: 0x0051, Action: Action{Type: "None"}},
		},
	}
}

func (c *Config) Generate() string {
	var b strings.Builder

	b.WriteString("devices: (\n{\n")

	b.WriteString(fmt.Sprintf("    name: \"%s\";\n", c.Name))

	b.WriteString("    smartshift:\n    {\n")
	b.WriteString(fmt.Sprintf("        on: %s;\n", boolStr(c.SmartShift.On)))
	b.WriteString(fmt.Sprintf("        threshold: %d;\n", c.SmartShift.Threshold))
	b.WriteString("    };\n")

	b.WriteString("    hiresscroll:\n    {\n")
	b.WriteString(fmt.Sprintf("        hires: %s;\n", boolStr(c.HiResScroll.Hires)))
	b.WriteString(fmt.Sprintf("        invert: %s;\n", boolStr(c.HiResScroll.Invert)))
	b.WriteString(fmt.Sprintf("        target: %s;\n", boolStr(c.HiResScroll.Target)))
	b.WriteString("    };\n")

	if c.Thumbwheel.Divert {
		b.WriteString("    thumbwheel:\n    {\n")
		b.WriteString(fmt.Sprintf("        divert: %s;\n", boolStr(c.Thumbwheel.Divert)))
		b.WriteString(fmt.Sprintf("        invert: %s;\n", boolStr(c.Thumbwheel.Invert)))

		if c.Thumbwheel.Left.Type != "" && c.Thumbwheel.Left.Type != "None" {
			b.WriteString("        left:\n        {\n")
			writeGestureContent(&b, c.Thumbwheel.Left, "        ")
			b.WriteString("        };\n")
		}
		if c.Thumbwheel.Right.Type != "" && c.Thumbwheel.Right.Type != "None" {
			b.WriteString("        right:\n        {\n")
			writeGestureContent(&b, c.Thumbwheel.Right, "        ")
			b.WriteString("        };\n")
		}
		if c.Thumbwheel.Tap.Type != "" && c.Thumbwheel.Tap.Type != "None" {
			b.WriteString("        tap:\n        {\n")
			writeActionContent(&b, c.Thumbwheel.Tap, "        ")
			b.WriteString("        };\n")
		}

		b.WriteString("    };\n")
	}

	b.WriteString(fmt.Sprintf("    dpi: %d;\n", c.DPI))

	if len(c.Buttons) > 0 {
		b.WriteString("\n    buttons: (\n")
		for i, btn := range c.Buttons {
			if btn.Action.Type == "" || btn.Action.Type == "None" {
				continue
			}
			if i > 0 {
				b.WriteString(",\n")
			}
			b.WriteString(fmt.Sprintf("        {\n            cid: 0x%04x;\n", btn.CID))
			b.WriteString("            action =\n            {\n")
			writeActionContent(&b, btn.Action, "            ")
			b.WriteString("            };\n        }")
		}
		b.WriteString("\n    );\n")
	}

	b.WriteString("}\n);\n")
	return b.String()
}

func writeActionContent(b *strings.Builder, a Action, indent string) {
	if a.Type == "" {
		a.Type = "None"
	}
	b.WriteString(fmt.Sprintf("%s    type: \"%s\";\n", indent, a.Type))

	switch a.Type {
	case "Keypress":
		if len(a.Keys) > 0 {
			keys := make([]string, len(a.Keys))
			for i, k := range a.Keys {
				keys[i] = fmt.Sprintf("\"%s\"", k)
			}
			b.WriteString(fmt.Sprintf("%s    keys: [%s];\n", indent, strings.Join(keys, ", ")))
		}
	case "Gestures":
		if len(a.Gestures) > 0 {
			b.WriteString(fmt.Sprintf("%s    gestures: (\n", indent))
			for j, g := range a.Gestures {
				if j > 0 {
					b.WriteString(",\n")
				}
				b.WriteString(fmt.Sprintf("%s        {\n", indent))
				b.WriteString(fmt.Sprintf("%s            direction: \"%s\";\n", indent, g.Direction))
				if g.Mode != "" {
					b.WriteString(fmt.Sprintf("%s            mode: \"%s\";\n", indent, g.Mode))
				}
				if g.Threshold > 0 {
					b.WriteString(fmt.Sprintf("%s            threshold: %d;\n", indent, g.Threshold))
				}
				if g.Interval > 0 {
					b.WriteString(fmt.Sprintf("%s            interval: %d;\n", indent, g.Interval))
				}
				if g.Axis != "" {
					b.WriteString(fmt.Sprintf("%s            axis: \"%s\";\n", indent, g.Axis))
				}
				if g.AxisMultiplier > 0 {
					b.WriteString(fmt.Sprintf("%s            axis_multiplier: %d;\n", indent, g.AxisMultiplier))
				}
				b.WriteString(fmt.Sprintf("%s            action =\n%s            {\n", indent, indent))
				writeActionContent(b, g.Action, indent+"            ")
				b.WriteString(fmt.Sprintf("%s            };\n", indent))
				b.WriteString(fmt.Sprintf("%s        }", indent))
			}
			b.WriteString(fmt.Sprintf("\n%s    );\n", indent))
		}
	case "CycleDPI":
		if len(a.DPIs) > 0 {
			dpiStrs := make([]string, len(a.DPIs))
			for i, d := range a.DPIs {
				dpiStrs[i] = fmt.Sprintf("%d", d)
			}
			b.WriteString(fmt.Sprintf("%s    dpis: [%s];\n", indent, strings.Join(dpiStrs, ", ")))
		}
	case "ChangeDPI":
		if a.Inc != 0 {
			b.WriteString(fmt.Sprintf("%s    inc: %d;\n", indent, a.Inc))
		}
	case "ChangeHost":
		if a.Host != "" {
			b.WriteString(fmt.Sprintf("%s    host: \"%s\";\n", indent, a.Host))
		}
	}
}

func writeGestureContent(b *strings.Builder, a Action, indent string) {
	if a.Type == "" {
		a.Type = "None"
	}
	b.WriteString(fmt.Sprintf("%s    mode: \"OnInterval\";\n", indent))
	b.WriteString(fmt.Sprintf("%s    interval: 2;\n", indent))
	b.WriteString(fmt.Sprintf("%s    action =\n%s    {\n", indent, indent))
	writeActionContent(b, a, indent+"    ")
	b.WriteString(fmt.Sprintf("%s    };\n", indent))
}

func LoadConfigFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	return ParseConfig(string(data)), nil
}

var (
	reName = regexp.MustCompile(`name:\s*"([^"]+)"`)
	reDPI  = regexp.MustCompile(`dpi:\s*(\d+)`)
	reBool = regexp.MustCompile(`\b(on|hires|invert|target|divert):\s*(true|false)`)
	reType = regexp.MustCompile(`type:\s*"([^"]+)"`)
	reKeys = regexp.MustCompile(`keys:\s*\[([^\]]*)\]`)
	reDPIs = regexp.MustCompile(`dpis:\s*\[([^\]]*)\]`)
	reHost = regexp.MustCompile(`host:\s*"([^"]+)"`)
	reDir  = regexp.MustCompile(`direction:\s*"([^"]+)"`)
	reMode = regexp.MustCompile(`mode:\s*"([^"]+)"`)
	reAxis = regexp.MustCompile(`axis:\s*"([^"]+)"`)

	reSmartShift  = regexp.MustCompile(`smartshift:\s*\{([^{}]*)\}`)
	reHiResScroll = regexp.MustCompile(`hiresscroll:\s*\{([^{}]*)\}`)
)

func extractBlock(data string, re *regexp.Regexp) string {
	m := re.FindStringSubmatch(data)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func extractBool(data string, re *regexp.Regexp, key string) bool {
	m := re.FindStringSubmatch(data)
	if m == nil {
		return false
	}
	for i := 1; i < len(m); i += 2 {
		if m[i] == key {
			return m[i+1] == "true"
		}
	}
	return false
}

func extractInt(data, key string) int {
	re := regexp.MustCompile(key + `:\s*(\d+)`)
	m := re.FindStringSubmatch(data)
	if m != nil {
		v, _ := strconv.Atoi(m[1])
		return v
	}
	return 0
}

func parseAction(data string) Action {
	a := Action{}
	if m := reType.FindStringSubmatch(data); m != nil {
		a.Type = m[1]
	}
	if m := reKeys.FindStringSubmatch(data); m != nil {
		for _, p := range strings.Split(m[1], ",") {
			p = strings.Trim(strings.TrimSpace(p), "\"")
			if p != "" {
				a.Keys = append(a.Keys, p)
			}
		}
	}
	if m := reDPIs.FindStringSubmatch(data); m != nil {
		for _, p := range strings.Split(m[1], ",") {
			p = strings.TrimSpace(p)
			if v, err := strconv.Atoi(p); err == nil {
				a.DPIs = append(a.DPIs, v)
			}
		}
	}
	a.Inc = extractInt(data, "inc")
	if m := reHost.FindStringSubmatch(data); m != nil {
		a.Host = m[1]
	}
	if a.Type == "Gestures" {
		a.Gestures = parseGestures(data)
	}
	return a
}

func parseGestures(data string) []Gesture {
	reGestStart := regexp.MustCompile(`gestures:\s*\(`)
	loc := reGestStart.FindStringIndex(data)
	if loc == nil {
		return nil
	}
	parenStart := loc[1] - 1
	parenEnd := matchParen(data, parenStart)
	if parenEnd < 0 {
		return nil
	}
	block := data[parenStart+1 : parenEnd]

	var gestures []Gesture
	var depth int
	var start int

	for i := 0; i < len(block); i++ {
		switch block[i] {
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			depth--
			if depth == 0 {
				g := parseSingleGesture(block[start : i+1])
				if g.Direction != "" {
					gestures = append(gestures, g)
				}
			}
		}
	}
	return gestures
}

func parseSingleGesture(block string) Gesture {
	var g Gesture
	if d := reDir.FindStringSubmatch(block); d != nil {
		g.Direction = d[1]
	}
	if d := reMode.FindStringSubmatch(block); d != nil {
		g.Mode = d[1]
	}
	g.Threshold = extractInt(block, "threshold")
	g.Interval = extractInt(block, "interval")
	if d := reAxis.FindStringSubmatch(block); d != nil {
		g.Axis = d[1]
	}
	g.AxisMultiplier = extractInt(block, "axis_multiplier")

	reActStart := regexp.MustCompile(`action\s*=\s*\{`)
	if actLoc := reActStart.FindStringIndex(block); actLoc != nil {
		actStart := actLoc[1] - 1
		actEnd := matchBrace(block, actStart)
		if actEnd >= 0 {
			g.Action = parseAction(block[actStart+1 : actEnd])
		}
	}
	return g
}

func parseThumbSide(data string) (left, right, tap Action) {
	reSide := regexp.MustCompile(`(left|right|tap)\s*:\s*\{`)
	for _, loc := range reSide.FindAllStringSubmatchIndex(data, -1) {
		sideName := data[loc[2]:loc[3]]
		braceStart := loc[1] - 1
		braceEnd := matchBrace(data, braceStart)
		if braceEnd < 0 {
			continue
		}
		sideBlock := data[braceStart+1 : braceEnd]
		a := parseAction(sideBlock)
		switch sideName {
		case "left":
			left = a
		case "right":
			right = a
		case "tap":
			tap = a
		}
	}
	return
}

func ParseConfig(data string) Config {
	var cfg Config

	if m := reName.FindStringSubmatch(data); m != nil {
		cfg.Name = m[1]
	}
	if m := reDPI.FindStringSubmatch(data); m != nil {
		if v, err := strconv.Atoi(m[1]); err == nil {
			cfg.DPI = v
		}
	}

	if block := extractBlock(data, reSmartShift); block != "" {
		for _, m := range reBool.FindAllStringSubmatch(block, -1) {
			if m[1] == "on" {
				cfg.SmartShift.On = m[2] == "true"
			}
		}
		cfg.SmartShift.Threshold = extractInt(block, "threshold")
	}

	if block := extractBlock(data, reHiResScroll); block != "" {
		for _, m := range reBool.FindAllStringSubmatch(block, -1) {
			switch m[1] {
			case "hires":
				cfg.HiResScroll.Hires = m[2] == "true"
			case "invert":
				cfg.HiResScroll.Invert = m[2] == "true"
			case "target":
				cfg.HiResScroll.Target = m[2] == "true"
			}
		}
	}

	// Find and parse thumbwheel block
	reTwStart := regexp.MustCompile(`thumbwheel\s*:\s*\{`)
	if twLoc := reTwStart.FindStringIndex(data); twLoc != nil {
		twStart := twLoc[1] - 1
		twEnd := matchBrace(data, twStart)
		if twEnd >= 0 {
			twBlock := data[twStart+1 : twEnd]
			for _, m := range reBool.FindAllStringSubmatch(twBlock, -1) {
				switch m[1] {
				case "divert":
					cfg.Thumbwheel.Divert = m[2] == "true"
				case "invert":
					cfg.Thumbwheel.Invert = m[2] == "true"
				}
			}
			cfg.Thumbwheel.Left, cfg.Thumbwheel.Right, cfg.Thumbwheel.Tap = parseThumbSide(twBlock)
		}
	}

	// Find button entries using brace matching
	reBtnStart := regexp.MustCompile(`\{\s*cid:\s*(0x[0-9a-fA-F]+)`)
	for _, loc := range reBtnStart.FindAllStringSubmatchIndex(data, -1) {
		start := loc[0]
		cid := parseHex(data[loc[2]:loc[3]])

		// Find the closing } for this button block
		btnEnd := matchBrace(data, start)
		if btnEnd < 0 {
			continue
		}
		btnBlock := data[start : btnEnd+1]

		// Find the action block using brace matching
		reActStart := regexp.MustCompile(`action\s*=\s*\{`)
		if actLoc := reActStart.FindStringIndex(btnBlock); actLoc != nil {
			actStart := actLoc[1] - 1 // position of {
			actEnd := matchBrace(btnBlock, actStart)
			if actEnd >= 0 {
				actionContent := btnBlock[actStart+1 : actEnd]
				action := parseAction(actionContent)
				cfg.Buttons = append(cfg.Buttons, ButtonConfig{CID: cid, Action: action})
			}
		}
	}

	return cfg
}

func matchBrace(data string, start int) int {
	if start < 0 || start >= len(data) || data[start] != '{' {
		return -1
	}
	depth := 0
	for i := start; i < len(data); i++ {
		switch data[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func matchParen(data string, start int) int {
	if start < 0 || start >= len(data) || data[start] != '(' {
		return -1
	}
	depth := 0
	for i := start; i < len(data); i++ {
		switch data[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func parseHex(s string) uint32 {
	s = strings.TrimPrefix(s, "0x")
	if v, err := strconv.ParseUint(s, 16, 32); err == nil {
		return uint32(v)
	}
	return 0
}

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
