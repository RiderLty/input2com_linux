package main

import "strconv"

var SpecialKeysMap = map[byte]byte{
	KeyLeftCtrl:   byte(1 << 0),
	KeyLeftShift:  byte(1 << 1),
	KeyLeftAlt:    byte(1 << 2),
	KeyLeftGui:    byte(1 << 3),
	KeyRightCtrl:  byte(1 << 4),
	KeyRightShift: byte(1 << 5),
	KeyRightAlt:   byte(1 << 6),
	KeyRightGui:   byte(1 << 7),
}

const (
	MouseBtnLeft    = byte(1 << 0) // 左键
	MouseBtnRight   = byte(1 << 1) // 右键
	MouseBtnMiddle  = byte(1 << 2) // 中键
	MouseBtnBack    = byte(1 << 3) // 后退键
	MouseBtnForward = byte(1 << 4) // 前进键

	KeyLeftCtrl   = byte(0xe0)
	KeyLeftShift  = byte(0xe1)
	KeyLeftAlt    = byte(0xe2)
	KeyLeftGui    = byte(0xe3)
	KeyRightCtrl  = byte(0xe4)
	KeyRightShift = byte(0xe5)
	KeyRightAlt   = byte(0xe6)
	KeyRightGui   = byte(0xe8)

	KeyA           = byte(0x04)
	KeyB           = byte(0x05)
	KeyC           = byte(0x06)
	KeyD           = byte(0x07)
	KeyE           = byte(0x08)
	KeyF           = byte(0x09)
	KeyG           = byte(0x0A)
	KeyH           = byte(0x0B)
	KeyI           = byte(0x0C)
	KeyJ           = byte(0x0D)
	KeyK           = byte(0x0E)
	KeyL           = byte(0x0F)
	KeyM           = byte(0x10)
	KeyN           = byte(0x11)
	KeyO           = byte(0x12)
	KeyP           = byte(0x13)
	KeyQ           = byte(0x14)
	KeyR           = byte(0x15)
	KeyS           = byte(0x16)
	KeyT           = byte(0x17)
	KeyU           = byte(0x18)
	KeyV           = byte(0x19)
	KeyW           = byte(0x1A)
	KeyX           = byte(0x1B)
	KeyY           = byte(0x1C)
	KeyZ           = byte(0x1D)
	Key1           = byte(0x1E)
	Key2           = byte(0x1F)
	Key3           = byte(0x20)
	Key4           = byte(0x21)
	Key5           = byte(0x22)
	Key6           = byte(0x23)
	Key7           = byte(0x24)
	Key8           = byte(0x25)
	Key9           = byte(0x26)
	Key0           = byte(0x27)
	KeyReturn      = byte(0x28)
	KeyEnter       = byte(0x28)
	KeyEsc         = byte(0x29)
	KeyEscape      = byte(0x29)
	KeyBckspc      = byte(0x2A)
	KeyBackspace   = byte(0x2A)
	KeyTab         = byte(0x2B)
	KeySpace       = byte(0x2C)
	KeyMinus       = byte(0x2D)
	KeyDash        = byte(0x2D)
	KeyEquals      = byte(0x2E)
	KeyEqual       = byte(0x2E)
	KeyLbracket    = byte(0x2F)
	KeyRbracket    = byte(0x30)
	KeyBackslash   = byte(0x31)
	KeyHash        = byte(0x32)
	KeyNumber      = byte(0x32)
	KeySemicolon   = byte(0x33)
	KeyQuote       = byte(0x34)
	KeyBackquote   = byte(0x35)
	KeyTilde       = byte(0x35)
	KeyComma       = byte(0x36)
	KeyPeriod      = byte(0x37)
	KeyStop        = byte(0x37)
	KeySlash       = byte(0x38)
	KeyCapsLock    = byte(0x39)
	KeyCapslock    = byte(0x39)
	KeyF1          = byte(0x3A)
	KeyF2          = byte(0x3B)
	KeyF3          = byte(0x3C)
	KeyF4          = byte(0x3D)
	KeyF5          = byte(0x3E)
	KeyF6          = byte(0x3F)
	KeyF7          = byte(0x40)
	KeyF8          = byte(0x41)
	KeyF9          = byte(0x42)
	KeyF10         = byte(0x43)
	KeyF11         = byte(0x44)
	KeyF12         = byte(0x45)
	KeyPrint       = byte(0x46)
	KeyScrollLock  = byte(0x47)
	KeyScrolllock  = byte(0x47)
	KeyPause       = byte(0x48)
	KeyInsert      = byte(0x49)
	KeyHome        = byte(0x4A)
	KeyPageup      = byte(0x4B)
	KeyPgup        = byte(0x4B)
	KeyDel         = byte(0x4C)
	KeyDelete      = byte(0x4C)
	KeyEnd         = byte(0x4D)
	KeyPagedown    = byte(0x4E)
	KeyPgdown      = byte(0x4E)
	KeyRight       = byte(0x4F)
	KeyLeft        = byte(0x50)
	KeyDown        = byte(0x51)
	KeyUp          = byte(0x52)
	KeyNumLock     = byte(0x53)
	KeyNumlock     = byte(0x53)
	KeyKpDivide    = byte(0x54)
	KeyKpMultiply  = byte(0x55)
	KeyKpMinus     = byte(0x56)
	KeyKpPlus      = byte(0x57)
	KeyKpEnter     = byte(0x58)
	KeyKpReturn    = byte(0x58)
	KeyKp1         = byte(0x59)
	KeyKp2         = byte(0x5A)
	KeyKp3         = byte(0x5B)
	KeyKp4         = byte(0x5C)
	KeyKp5         = byte(0x5D)
	KeyKp6         = byte(0x5E)
	KeyKp7         = byte(0x5F)
	KeyKp8         = byte(0x60)
	KeyKp9         = byte(0x61)
	KeyKp0         = byte(0x62)
	KeyKpPeriod    = byte(0x63)
	KeyKpStop      = byte(0x63)
	KeyApplication = byte(0x65)
	KeyPower       = byte(0x66)
	KeyKpEquals    = byte(0x67)
	KeyKpEqual     = byte(0x67)
	KeyF13         = byte(0x68)
	KeyF14         = byte(0x69)
	KeyF15         = byte(0x6A)
	KeyF16         = byte(0x6B)
	KeyF17         = byte(0x6C)
	KeyF18         = byte(0x6D)
	KeyF19         = byte(0x6E)
	KeyF20         = byte(0x6F)
	KeyF21         = byte(0x70)
	KeyF22         = byte(0x71)
	KeyF23         = byte(0x72)
	KeyF24         = byte(0x73)
	KeyExecute     = byte(0x74)
	KeyHelp        = byte(0x75)
	KeyMenu        = byte(0x76)
	KeySelect      = byte(0x77)
	KeyCancel      = byte(0x78)
	KeyRedo        = byte(0x79)
	KeyUndo        = byte(0x7A)
	KeyCut         = byte(0x7B)
	KeyCopy        = byte(0x7C)
	KeyPaste       = byte(0x7D)
	KeyFind        = byte(0x7E)
	KeyMute        = byte(0x7F)
	KeyVolumeUp    = byte(0x80)
	KeyVolumeDown  = byte(0x81)
)

var Linux2hid = map[uint8]uint8{
	30:  4,
	48:  5,
	46:  6,
	32:  7,
	18:  8,
	33:  9,
	34:  10,
	35:  11,
	23:  12,
	36:  13,
	37:  14,
	38:  15,
	50:  16,
	49:  17,
	24:  18,
	25:  19,
	16:  20,
	19:  21,
	31:  22,
	20:  23,
	22:  24,
	47:  25,
	17:  26,
	45:  27,
	21:  28,
	44:  29,
	2:   30,
	3:   31,
	4:   32,
	5:   33,
	6:   34,
	7:   35,
	8:   36,
	9:   37,
	10:  38,
	11:  39,
	28:  40,
	1:   41,
	14:  42,
	15:  43,
	57:  44,
	12:  45,
	13:  46,
	26:  47,
	27:  48,
	43:  49,
	39:  51,
	40:  52,
	41:  53,
	51:  54,
	52:  55,
	53:  56,
	58:  57,
	59:  58,
	60:  59,
	61:  60,
	62:  61,
	63:  62,
	64:  63,
	65:  64,
	66:  65,
	67:  66,
	68:  67,
	87:  68,
	88:  69,
	99:  70,
	70:  71,
	119: 72,
	110: 73,
	102: 74,
	104: 75,
	111: 76,
	107: 77,
	109: 78,
	106: 79,
	105: 80,
	108: 81,
	103: 82,
	69:  83,
	98:  84,
	55:  85,
	74:  86,
	78:  87,
	96:  88,
	79:  89,
	80:  90,
	81:  91,
	75:  92,
	76:  93,
	77:  94,
	71:  95,
	72:  96,
	73:  97,
	82:  98,
	83:  99,
	86:  100,
	127: 101,
	29:  224,
	42:  225,
	56:  226,
	125: 227,
	97:  228,
	54:  229,
	100: 230,
	126: 232,
}

var MouseValidKeys = map[string]bool{
	strconv.FormatUint(uint64(MouseBtnLeft), 10):    true,
	strconv.FormatUint(uint64(MouseBtnRight), 10):   true,
	strconv.FormatUint(uint64(MouseBtnMiddle), 10):  true,
	strconv.FormatUint(uint64(MouseBtnBack), 10):    true,
	strconv.FormatUint(uint64(MouseBtnForward), 10): true,
}

var KeyboardValidKeys = map[string]bool{
	strconv.FormatUint(uint64(KeyLeftCtrl), 10):    true,
	strconv.FormatUint(uint64(KeyLeftShift), 10):   true,
	strconv.FormatUint(uint64(KeyLeftAlt), 10):     true,
	strconv.FormatUint(uint64(KeyLeftGui), 10):     true,
	strconv.FormatUint(uint64(KeyRightCtrl), 10):   true,
	strconv.FormatUint(uint64(KeyRightShift), 10):  true,
	strconv.FormatUint(uint64(KeyRightAlt), 10):    true,
	strconv.FormatUint(uint64(KeyRightGui), 10):    true,
	strconv.FormatUint(uint64(KeyA), 10):           true,
	strconv.FormatUint(uint64(KeyB), 10):           true,
	strconv.FormatUint(uint64(KeyC), 10):           true,
	strconv.FormatUint(uint64(KeyD), 10):           true,
	strconv.FormatUint(uint64(KeyE), 10):           true,
	strconv.FormatUint(uint64(KeyF), 10):           true,
	strconv.FormatUint(uint64(KeyG), 10):           true,
	strconv.FormatUint(uint64(KeyH), 10):           true,
	strconv.FormatUint(uint64(KeyI), 10):           true,
	strconv.FormatUint(uint64(KeyJ), 10):           true,
	strconv.FormatUint(uint64(KeyK), 10):           true,
	strconv.FormatUint(uint64(KeyL), 10):           true,
	strconv.FormatUint(uint64(KeyM), 10):           true,
	strconv.FormatUint(uint64(KeyN), 10):           true,
	strconv.FormatUint(uint64(KeyO), 10):           true,
	strconv.FormatUint(uint64(KeyP), 10):           true,
	strconv.FormatUint(uint64(KeyQ), 10):           true,
	strconv.FormatUint(uint64(KeyR), 10):           true,
	strconv.FormatUint(uint64(KeyS), 10):           true,
	strconv.FormatUint(uint64(KeyT), 10):           true,
	strconv.FormatUint(uint64(KeyU), 10):           true,
	strconv.FormatUint(uint64(KeyV), 10):           true,
	strconv.FormatUint(uint64(KeyW), 10):           true,
	strconv.FormatUint(uint64(KeyX), 10):           true,
	strconv.FormatUint(uint64(KeyY), 10):           true,
	strconv.FormatUint(uint64(KeyZ), 10):           true,
	strconv.FormatUint(uint64(Key1), 10):           true,
	strconv.FormatUint(uint64(Key2), 10):           true,
	strconv.FormatUint(uint64(Key3), 10):           true,
	strconv.FormatUint(uint64(Key4), 10):           true,
	strconv.FormatUint(uint64(Key5), 10):           true,
	strconv.FormatUint(uint64(Key6), 10):           true,
	strconv.FormatUint(uint64(Key7), 10):           true,
	strconv.FormatUint(uint64(Key8), 10):           true,
	strconv.FormatUint(uint64(Key9), 10):           true,
	strconv.FormatUint(uint64(Key0), 10):           true,
	strconv.FormatUint(uint64(KeyReturn), 10):      true,
	strconv.FormatUint(uint64(KeyEnter), 10):       true,
	strconv.FormatUint(uint64(KeyEsc), 10):         true,
	strconv.FormatUint(uint64(KeyEscape), 10):      true,
	strconv.FormatUint(uint64(KeyBckspc), 10):      true,
	strconv.FormatUint(uint64(KeyBackspace), 10):   true,
	strconv.FormatUint(uint64(KeyTab), 10):         true,
	strconv.FormatUint(uint64(KeySpace), 10):       true,
	strconv.FormatUint(uint64(KeyMinus), 10):       true,
	strconv.FormatUint(uint64(KeyDash), 10):        true,
	strconv.FormatUint(uint64(KeyEquals), 10):      true,
	strconv.FormatUint(uint64(KeyEqual), 10):       true,
	strconv.FormatUint(uint64(KeyLbracket), 10):    true,
	strconv.FormatUint(uint64(KeyRbracket), 10):    true,
	strconv.FormatUint(uint64(KeyBackslash), 10):   true,
	strconv.FormatUint(uint64(KeyHash), 10):        true,
	strconv.FormatUint(uint64(KeyNumber), 10):      true,
	strconv.FormatUint(uint64(KeySemicolon), 10):   true,
	strconv.FormatUint(uint64(KeyQuote), 10):       true,
	strconv.FormatUint(uint64(KeyBackquote), 10):   true,
	strconv.FormatUint(uint64(KeyTilde), 10):       true,
	strconv.FormatUint(uint64(KeyComma), 10):       true,
	strconv.FormatUint(uint64(KeyPeriod), 10):      true,
	strconv.FormatUint(uint64(KeyStop), 10):        true,
	strconv.FormatUint(uint64(KeySlash), 10):       true,
	strconv.FormatUint(uint64(KeyCapsLock), 10):    true,
	strconv.FormatUint(uint64(KeyCapslock), 10):    true,
	strconv.FormatUint(uint64(KeyF1), 10):          true,
	strconv.FormatUint(uint64(KeyF2), 10):          true,
	strconv.FormatUint(uint64(KeyF3), 10):          true,
	strconv.FormatUint(uint64(KeyF4), 10):          true,
	strconv.FormatUint(uint64(KeyF5), 10):          true,
	strconv.FormatUint(uint64(KeyF6), 10):          true,
	strconv.FormatUint(uint64(KeyF7), 10):          true,
	strconv.FormatUint(uint64(KeyF8), 10):          true,
	strconv.FormatUint(uint64(KeyF9), 10):          true,
	strconv.FormatUint(uint64(KeyF10), 10):         true,
	strconv.FormatUint(uint64(KeyF11), 10):         true,
	strconv.FormatUint(uint64(KeyF12), 10):         true,
	strconv.FormatUint(uint64(KeyPrint), 10):       true,
	strconv.FormatUint(uint64(KeyScrollLock), 10):  true,
	strconv.FormatUint(uint64(KeyScrolllock), 10):  true,
	strconv.FormatUint(uint64(KeyPause), 10):       true,
	strconv.FormatUint(uint64(KeyInsert), 10):      true,
	strconv.FormatUint(uint64(KeyHome), 10):        true,
	strconv.FormatUint(uint64(KeyPageup), 10):      true,
	strconv.FormatUint(uint64(KeyPgup), 10):        true,
	strconv.FormatUint(uint64(KeyDel), 10):         true,
	strconv.FormatUint(uint64(KeyDelete), 10):      true,
	strconv.FormatUint(uint64(KeyEnd), 10):         true,
	strconv.FormatUint(uint64(KeyPagedown), 10):    true,
	strconv.FormatUint(uint64(KeyPgdown), 10):      true,
	strconv.FormatUint(uint64(KeyRight), 10):       true,
	strconv.FormatUint(uint64(KeyLeft), 10):        true,
	strconv.FormatUint(uint64(KeyDown), 10):        true,
	strconv.FormatUint(uint64(KeyUp), 10):          true,
	strconv.FormatUint(uint64(KeyNumLock), 10):     true,
	strconv.FormatUint(uint64(KeyNumlock), 10):     true,
	strconv.FormatUint(uint64(KeyKpDivide), 10):    true,
	strconv.FormatUint(uint64(KeyKpMultiply), 10):  true,
	strconv.FormatUint(uint64(KeyKpMinus), 10):     true,
	strconv.FormatUint(uint64(KeyKpPlus), 10):      true,
	strconv.FormatUint(uint64(KeyKpEnter), 10):     true,
	strconv.FormatUint(uint64(KeyKpReturn), 10):    true,
	strconv.FormatUint(uint64(KeyKp1), 10):         true,
	strconv.FormatUint(uint64(KeyKp2), 10):         true,
	strconv.FormatUint(uint64(KeyKp3), 10):         true,
	strconv.FormatUint(uint64(KeyKp4), 10):         true,
	strconv.FormatUint(uint64(KeyKp5), 10):         true,
	strconv.FormatUint(uint64(KeyKp6), 10):         true,
	strconv.FormatUint(uint64(KeyKp7), 10):         true,
	strconv.FormatUint(uint64(KeyKp8), 10):         true,
	strconv.FormatUint(uint64(KeyKp9), 10):         true,
	strconv.FormatUint(uint64(KeyKp0), 10):         true,
	strconv.FormatUint(uint64(KeyKpPeriod), 10):    true,
	strconv.FormatUint(uint64(KeyKpStop), 10):      true,
	strconv.FormatUint(uint64(KeyApplication), 10): true,
	strconv.FormatUint(uint64(KeyPower), 10):       true,
	strconv.FormatUint(uint64(KeyKpEquals), 10):    true,
	strconv.FormatUint(uint64(KeyKpEqual), 10):     true,
	strconv.FormatUint(uint64(KeyF13), 10):         true,
	strconv.FormatUint(uint64(KeyF14), 10):         true,
	strconv.FormatUint(uint64(KeyF15), 10):         true,
	strconv.FormatUint(uint64(KeyF16), 10):         true,
	strconv.FormatUint(uint64(KeyF17), 10):         true,
	strconv.FormatUint(uint64(KeyF18), 10):         true,
	strconv.FormatUint(uint64(KeyF19), 10):         true,
	strconv.FormatUint(uint64(KeyF20), 10):         true,
	strconv.FormatUint(uint64(KeyF21), 10):         true,
	strconv.FormatUint(uint64(KeyF22), 10):         true,
	strconv.FormatUint(uint64(KeyF23), 10):         true,
	strconv.FormatUint(uint64(KeyF24), 10):         true,
	strconv.FormatUint(uint64(KeyExecute), 10):     true,
	strconv.FormatUint(uint64(KeyHelp), 10):        true,
	strconv.FormatUint(uint64(KeyMenu), 10):        true,
	strconv.FormatUint(uint64(KeySelect), 10):      true,
	strconv.FormatUint(uint64(KeyCancel), 10):      true,
	strconv.FormatUint(uint64(KeyRedo), 10):        true,
	strconv.FormatUint(uint64(KeyUndo), 10):        true,
	strconv.FormatUint(uint64(KeyCut), 10):         true,
	strconv.FormatUint(uint64(KeyCopy), 10):        true,
	strconv.FormatUint(uint64(KeyPaste), 10):       true,
	strconv.FormatUint(uint64(KeyFind), 10):        true,
	strconv.FormatUint(uint64(KeyMute), 10):        true,
	strconv.FormatUint(uint64(KeyVolumeUp), 10):    true,
	strconv.FormatUint(uint64(KeyVolumeDown), 10):  true,
}
var MouseKeyDown = map[byte]string{
	MouseBtnLeft:    "km.left(1)\r",
	MouseBtnRight:   "km.right(1)\r",
	MouseBtnMiddle:  "km.middle(1)\r",
	MouseBtnBack:    "km.side1(1)\r",
	MouseBtnForward: "km.side2(1)\r",
}
var MouseKeyUp = map[byte]string{
	MouseBtnLeft:    "km.left(0)\r",
	MouseBtnRight:   "km.right(0)\r",
	MouseBtnMiddle:  "km.middle(0)\r",
	MouseBtnBack:    "km.side1(0)\r",
	MouseBtnForward: "km.side2(0)\r",
}
