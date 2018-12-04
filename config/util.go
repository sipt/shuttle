package config

import "fmt"

func EmojiDecode(data []byte) int {
	index, length := 0, len(data)
	offset := 0
	for index < length {
		if data[index] == '\\' && data[index+1] == 'U' {
			index += 2
			decodeEmoji(data[offset:offset], data[index:index+8])
			offset += 4
			index += 8
		} else {
			if index != offset {
				data[offset] = data[index]
			}
			offset ++
			index ++
		}
	}
	return offset
}

func decodeEmoji(dst []byte, src []byte) (err error) {
	const code_length = 8
	var value int
	for k := 0; k < code_length; k++ {
		if !is_hex(src, k) {
			err = fmt.Errorf("is not hex :%v", src[k])
			return
		}
		value = (value << 4) + as_hex(src, k)
	}

	// Check the value and write the character.
	if (value >= 0xD800 && value <= 0xDFFF) || value > 0x10FFFF {
		err = fmt.Errorf("is not hex :%v", value)
		return
	}
	if value <= 0x7F {
		dst = append(dst, byte(value))
	} else if value <= 0x7FF {
		dst = append(dst, byte(0xC0+(value>>6)))
		dst = append(dst, byte(0x80+(value&0x3F)))
	} else if value <= 0xFFFF {
		dst = append(dst, byte(0xE0+(value>>12)))
		dst = append(dst, byte(0x80+((value>>6)&0x3F)))
		dst = append(dst, byte(0x80+(value&0x3F)))
	} else {
		dst = append(dst, byte(0xF0+(value>>18)))
		dst = append(dst, byte(0x80+((value>>12)&0x3F)))
		dst = append(dst, byte(0x80+((value>>6)&0x3F)))
		dst = append(dst, byte(0x80+(value&0x3F)))
	}
	return nil
}
func is_hex(b []byte, i int) bool {
	return b[i] >= '0' && b[i] <= '9' || b[i] >= 'A' && b[i] <= 'F' || b[i] >= 'a' && b[i] <= 'f'
}

// Get the value of a hex-digit.
func as_hex(b []byte, i int) int {
	bi := b[i]
	if bi >= 'A' && bi <= 'F' {
		return int(bi) - 'A' + 10
	}
	if bi >= 'a' && bi <= 'f' {
		return int(bi) - 'a' + 10
	}
	return int(bi) - '0'
}
