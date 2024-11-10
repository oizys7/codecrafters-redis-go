package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"strconv"
)

type Decoder struct {
	input     *bufio.Reader
	readCount int
	buffer    []byte
}

func NewDecoder(reader io.Reader) *Decoder {
	parser := new(Decoder)
	parser.input = bufio.NewReader(reader)
	parser.buffer = make([]byte, 8)
	return parser
}

func (dec *Decoder) readByte() (byte, error) {
	b, err := dec.input.ReadByte()
	if err != nil {
		return 0, err
	}
	dec.readCount++
	return b, nil
}

func (dec *Decoder) readFull(buf []byte) error {
	n, err := io.ReadFull(dec.input, buf)
	if err != nil {
		return err
	}
	dec.readCount += n
	return nil
}

func (dec *Decoder) getInt64ExpireMs() (expireMs int64) {
	return int64(binary.LittleEndian.Uint64(dec.buffer))
}

func (dec *Decoder) decodeLength() (int, error) {
	num, err := dec.readByte()
	if err != nil {
		return 0, err
	}

	switch {
	case num <= 63: // leading bits 00
		// Remaining 6 bits are the length.
		return int(num & 0b00111111), nil
	case num <= 127: // leading bits 01
		// Remaining 6 bits plus next byte are the length
		nextNum, err := dec.readByte()
		if err != nil {
			return 0, err
		}
		length := binary.BigEndian.Uint16([]byte{num & 0b00111111, nextNum})
		return int(length), nil
	case num <= 191: // leading bits 10
		// Next 4 getBytes are the length
		getBytes := make([]byte, 4)
		_, err := dec.readByte()
		if err != nil {
			return 0, err
		}
		length := binary.BigEndian.Uint32(getBytes)
		return int(length), nil
	case num <= 255: // leading bits 11
		// Next 6 bits indicate the format of the encoded object.
		// TODO: This will result in problems on the next read, possibly.
		valueType := num & 0b00111111
		return int(valueType), nil
	default:
		return 0, err
	}
}

func (dec *Decoder) parseRDB() ([]string, error) {
	var result []string
	reader := dec.input

	for {
		opcode, err := reader.ReadByte()
		if err != nil {
			return result, err
		}

		switch opcode {
		case opCodeSelectDB:
			// Following byte(s) is the db number.
			dbNum, err := dec.decodeLength()
			if err != nil {
				return result, err
			}
			logger.Debug("DB number: " + strconv.Itoa(dbNum))
		case opCodeAux:
			// Length prefixed key and value strings follow.
			var kv [][]byte
			for i := 0; i < 2; i++ {
				length, err := dec.decodeLength()
				if err != nil {
					return result, err
				}
				data := make([]byte, length)
				if _, err = reader.Read(data); err != nil {
					return result, err
				}
				kv = append(kv, data)
			}
		case opCodeResizeDB:
			// Implement
			hashTableNum, err := dec.decodeLength()
			if err != nil {
				return result, err
			}
			_, _ = reader.ReadByte()
			logger.Debug("Hash table resize: " + strconv.Itoa(hashTableNum))
		case opCodeExpireTimeMs:
			_, _ = reader.ReadByte()
		case typeString:
			var kv [][]byte
			for i := 0; i < 2; i++ {
				length, err := dec.decodeLength()
				if err != nil {
					return result, err
				}
				data := make([]byte, length)
				if _, err = reader.Read(data); err != nil {
					return result, err
				}
				kv = append(kv, data)
			}
			result = append(result, string(kv[0]), string(kv[1]))
		case opCodeEOF:
			// Get the 8-byte checksum after this
			checksum := make([]byte, 8)
			_, err := reader.Read(checksum)
			if err != nil {
				return result, err
			}
			return result, nil
		default:
			// Handle any other tags.
		}
	}
}

//func sliceIndex(data []byte, sep byte) int {
//	for i, b := range data {
//		if b == sep {
//			return i
//		}
//	}
//	return -1
//}
//func parseTable(bytes []byte) []byte {
//	start := sliceIndex(bytes, opCodeResizeDB)
//	end := sliceIndex(bytes, opCodeEOF)
//	return bytes[start:end]
//}
