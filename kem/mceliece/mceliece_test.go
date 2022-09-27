package mceliece

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func FindTestDataByte(searchKey string) ([]byte, error) {
	file, err := os.Open("testdata/testdata.txt")
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil, errors.New(fmt.Sprintf("key %s not found", searchKey))
		}
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		var key, value string
		for i, v := range strings.Split(line, "=") {
			switch i {
			case 0:
				key = strings.TrimSpace(v)
			case 1:
				value = strings.TrimSpace(v)
			default:
				break
			}
		}

		if value == "" {
			return nil, errors.New(fmt.Sprintf("value is nil for key %s", key))
		}

		if key != searchKey {
			continue
		}

		data, err := hex.DecodeString(value)
		return data, err
	}
}

func FindTestDataU16(searchKey string) ([]uint16, error) {
	data, err := FindTestDataByte(searchKey)
	if err != nil {
		return nil, err
	}

	if len(data)%2 != 0 {
		return nil, errors.New("data length does not align")
	}

	out := make([]uint16, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		out[i] = binary.BigEndian.Uint16(data)
		data = data[2:]
	}

	return out, nil
}

func FindTestDataI16(searchKey string) ([]int16, error) {
	data, err := FindTestDataU16(searchKey)
	if err != nil {
		return nil, err
	}

	out := make([]int16, len(data))
	for i := 0; i < len(data); i++ {
		out[i] = int16(data[i])
	}

	return out, nil
}

func FindTestDataU32(searchKey string) ([]uint32, error) {
	data, err := FindTestDataByte(searchKey)
	if err != nil {
		return nil, err
	}

	if len(data)%4 != 0 {
		return nil, errors.New("data length does not align")
	}

	out := make([]uint32, len(data)/4)
	for i := 0; i < len(data); i += 4 {
		out[i] = binary.BigEndian.Uint32(data)
		data = data[4:]
	}

	return out, nil
}

func FindTestDataU64(searchKey string) ([]uint64, error) {
	data, err := FindTestDataByte(searchKey)
	if err != nil {
		return nil, err
	}

	if len(data)%8 != 0 {
		return nil, errors.New("data length does not align")
	}

	out := make([]uint64, len(data)/8)
	for i := 0; i < len(data); i += 8 {
		out[i] = binary.BigEndian.Uint64(data)
		data = data[8:]
	}

	return out, nil
}

func Test(t *testing.T) {
	data, err := FindTestDataByte("controlbits_kat8_mceliece348864_pi")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(data)
}
