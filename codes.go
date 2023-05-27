package main

const (
	Bel           = "\a"
	Backspace     = "\b"
	HorizontalTab = "\t"
	LineFeed      = "\n"
	VerticalTab   = "\v"
	FormFeed      = "\f"
	CarrageReturn = "\r"
	Esc           = "\x1B"
	Del           = "\x7F"

	// Cursor handling
	Hide    = Esc + "[?25l"
	Show    = Esc + "[?25h"
	Home    = Esc + "[H"
	Save    = Esc + "7" // Also save attributes
	Restore = Esc + "8"

	// Screen handling
	SaveScreen      = Esc + "[?47h"
	RestoreScreen   = Esc + "[?47l"
	EnAltScreenBuf  = Esc + "[?1049h"
	DisAltScreenBuf = Esc + "[?1049l"

	// Clear line
	ClearEnd       = Esc + "[0K" // From cursor to end of line
	ClearBeginning = Esc + "[1K" // From cursor to beginning of line
	ClearLine      = Esc + "[2K"
	ClearScreen    = Esc + "[2J"

	// Control text attributes
	Reset              = Esc + "[0m"
	Bold               = Esc + "[1m"
	Dim                = Esc + "[2m"
	Italic             = Esc + "[3m"
	Underscore         = Esc + "[4m"
	Blink              = Esc + "[5m"
	Inverse            = Esc + "[7m"
	Hidden             = Esc + "[8m"
	CrossedOut         = Esc + "[9m"
	Fraktur            = Esc + "[20m"
	DoubleUnderscore   = Esc + "[21m"
	NoBold             = Esc + "[22m"
	NoItalic           = Esc + "[23m"
	NoDoubleUnderscore = Esc + "[24m"
	NoBlink            = Esc + "[25m"
	NoInverse          = Esc + "[27m"
	Reveal             = Esc + "[28m"
	NoCrossedOut       = Esc + "[29m"
	Framed             = Esc + "[51m"
	Encircled          = Esc + "[53m"

	// Control foreground coloring
	FgBlack         = Esc + "[30m"
	FgRed           = Esc + "[31m"
	FgGreen         = Esc + "[32m"
	FgYellow        = Esc + "[33m"
	FgBlue          = Esc + "[34m"
	FgMagenta       = Esc + "[35m"
	FgCyan          = Esc + "[36m"
	FgWhite         = Esc + "[37m"
	FgDefault       = Esc + "[39m"
	FgBrightBlack   = Esc + "[90m"
	FgBrightRed     = Esc + "[91m"
	FgBrightGreen   = Esc + "[92m"
	FgBrightYellow  = Esc + "[93m"
	FgBrightBlue    = Esc + "[94m"
	FgBrightMagenta = Esc + "[95m"
	FgBrightCyan    = Esc + "[96m"
	FgBrightWhite   = Esc + "[97m"

	// Control background coloring
	BgBlack         = Esc + "[40m"
	BgRed           = Esc + "[41m"
	BgGreen         = Esc + "[42m"
	BgYellow        = Esc + "[43m"
	BgBlue          = Esc + "[44m"
	BgMagenta       = Esc + "[45m"
	BgCyan          = Esc + "[46m"
	BgWhite         = Esc + "[47m"
	BgDefault       = Esc + "[49m"
	BgBrightBlack   = Esc + "[100m"
	BgBrightRed     = Esc + "[101m"
	BgBrightGreen   = Esc + "[102m"
	BgBrightYellow  = Esc + "[103m"
	BgBrightBlue    = Esc + "[104m"
	BgBrightMagenta = Esc + "[105m"
	BgBrightCyan    = Esc + "[106m"
	BgBrightWhite   = Esc + "[107m"
)
