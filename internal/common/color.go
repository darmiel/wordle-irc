package common

type Color struct {
	prefix rune
	value  string
}

func colorOf(value string) Color {
	return Color{
		prefix: 3,
		value:  value,
	}
}

func toggle(prefix rune) Color {
	return Color{
		prefix: prefix,
		value:  "",
	}
}

var (
	ColorGreenBG  = colorOf("0,3")
	ColorYellowBG = colorOf("0,7")
	ColorGreyBG   = colorOf("0,14")
	ColorCyanBG   = colorOf("0,10")
	ColorCyan     = colorOf("10")
	ColorRedBG    = colorOf("0,04")
)

//goland:noinspection GoUnusedGlobalVariable
var (
	StyleReset         = toggle(0x0F)
	StyleBold          = toggle(0x02)
	StyleItalics       = toggle(0x1D)
	StyleUnderline     = toggle(0x1F)
	StyleStrikethrough = toggle(0x1E)
	StyleMonospace     = toggle(0x11)
)

func (c Color) String() string {
	return string(c.prefix) + c.value
}

func (c Color) Enclose(msg string) string {
	return c.String() + " " + msg + " " + StyleReset.String()
}
