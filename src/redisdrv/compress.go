package redisdrv

import "github.com/davyxu/cellnet/util"

// 实现了该方法的对象, 将自动开启对象字节压缩
type ObjectCompress interface {
	EnableObjectCompress()
}

var (
	compressHeader = []byte("ZIPCOMPRESS_HEADER")
)

func isCompressedData(data []byte) (ret bool, retData []byte) {
	if len(data) < len(compressHeader) {
		return false, data
	}

	for index, c := range compressHeader {
		if data[index] != c {
			return false, data
		}
	}

	return true, data[len(compressHeader):]
}

func encodeData(obj interface{}, data []byte) ([]byte, error) {
	if _, ok := obj.(ObjectCompress); ok {
		ret, err := util.CompressBytes(data)
		if err != nil {
			return data, err
		}

		return append(compressHeader, ret...), nil

	} else {
		return data, nil
	}
}

func decodeData(data []byte) ([]byte, error) {
	if ok, rawData := isCompressedData(data); ok {
		return util.DecompressBytes(rawData)
	} else {
		return data, nil
	}
}
