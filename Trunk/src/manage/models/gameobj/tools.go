package gameobj

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
)

func CompressData(src []byte) []byte {
	//压缩数据
	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.LittleEndian, uint32(COMPRESSED_DATA_VERSION))
	binary.Write(b, binary.LittleEndian, uint32(len(src)))
	w := zlib.NewWriter(b)
	w.Write(src)
	w.Close()
	return b.Bytes()
}

func OutputJson(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(data)
}
