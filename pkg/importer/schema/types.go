package schema

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"strconv"
	"strings"
)

func decodeHexFromXML(d *xml.Decoder, start xml.StartElement) ([]byte, error) {
	var value string
	if err := d.DecodeElement(&value, &start); err != nil {
		return nil, err
	}
	decoded, err := hex.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

type HexUint16 uint16

func (i *HexUint16) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	decoded, err := decodeHexFromXML(d, start)
	if err != nil {
		return err
	}
	*i = HexUint16(binary.BigEndian.Uint16(decoded))
	return nil
}

type HexUint32 uint32

func (i *HexUint32) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	decoded, err := decodeHexFromXML(d, start)
	if err != nil {
		return err
	}
	*i = HexUint32(binary.BigEndian.Uint32(decoded))
	return nil
}

type HexUint64 uint64

func (i *HexUint64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	decoded, err := decodeHexFromXML(d, start)
	if err != nil {
		return err
	}
	*i = HexUint64(binary.BigEndian.Uint64(decoded))
	return nil
}

type TrimUint32 uint32

func (i *TrimUint32) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	if err := d.DecodeElement(&value, &start); err != nil {
		return err
	}
	value = strings.TrimSpace(value)
	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return err
	}
	*i = TrimUint32(parsed)
	return nil
}
