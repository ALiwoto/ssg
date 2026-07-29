package legacyQs

// Ss will generate a new StrongString
// with the specified non-encoded string value.
func Ss(s string) StrongString {
	_strong := StrongString{}
	_strong._setValue(s)
	return _strong
}

// Qss will generate a new QString
// with the specified non-encoded string value.
func Qss(s string) QString {
	str := Ss(s)
	return &str
}

// Sb will generate a new StrongString
// with the specified non-encoded bytes value.
func Sb(b []byte) StrongString {
	return Ss(string(b))
}

// QSb will generate a new QString
// with the specified non-encoded bytes value.
func Qsb(b []byte) QString {
	str := Ss(string(b))
	return &str
}

// SS will generate a new StrongString
// with the specified non-encoded string value.
func SsPtr(s string) *StrongString {
	strong := StrongString{}
	strong._setValue(s)
	return &strong
}

func ToStrSlice(qs []QString) []string {
	tmp := make([]string, len(qs))
	for i, current := range qs {
		tmp[i] = current.GetValue()
	}
	return tmp
}

func ToQSlice(strs []string) []QString {
	tmp := make([]QString, len(strs))
	for i, current := range strs {
		tmp[i] = SsPtr(current)
	}
	return tmp
}

func repairString(value string) string {
	entered := false
	ignoreNext := false
	final := EMPTY
	last := len(value) - BaseIndex
	next := BaseIndex
	for i, current := range value {
		if ignoreNext {
			ignoreNext = false
			continue
		}

		if current == CHAR_STR {
			if !entered {
				entered = true
			} else {
				entered = false
			}

			final += string(current)
			continue
		} else {
			if !entered {
				final += string(current)
				continue
			}

			if isSpecial(current) {
				final += BackSlash + string(current)
				continue
			} else {
				if current == LineChar {
					if i != last {
						next = i + BaseOneIndex
						if value[next] == LineChar {
							final += BackSlash +
								string(current) + string(current)
							ignoreNext = true
							continue
						}
					}
				}
			}
		}

		final += string(current)
	}

	return final
}

func isSpecial(r rune) bool {
	switch r {
	case EqualChar, DPointChar:
		return true
	default:
		return false
	}
}
