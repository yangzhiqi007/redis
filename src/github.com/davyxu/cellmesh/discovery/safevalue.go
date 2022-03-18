package discovery

import (
	"fmt"

	"reflect"
)

//KV中的Value最大不超过512K,
const (
	// 不能直接保存二进制，底层用Json转base64，base64的二进制比原二进制要大最终二进制不到512K就会达到限制
	PackedValueSize = 300 * 1024
)

type rawGetter interface {
	// 获取原始值
	GetRawValue(key string) ([]byte, error)
	GetValueDirect(key string, valuePtr interface{}) error
}

func getMultiKey(sd rawGetter, key string) (ret []string, err error) {

	mainKey := key

	ret = append(ret, mainKey)

	for i := 1; i < 1000; i++ {

		key = fmt.Sprintf("%s.%d", mainKey, i)

		_, err = sd.GetRawValue(key)
		if err != nil {

			if err.Error() == "value not exists" {
				return ret, nil
			}

			return nil, err

		}

		ret = append(ret, key)
	}

	return ret, nil
}

type Option struct {
	PrettyPrint   bool
	Compress      bool //压缩处理
	DisableNotify bool // 禁用同步到每个服务发现
}

// compress value按 key, key.1, key.2 ... 保存
func SafeSetValue(sd Discovery, key string, value interface{}, opt Option) error {
	if opt.Compress {

		// 先删除key
		keys, err := getMultiKey(sd.(rawGetter), key)
		if err != nil {
			return err
		}
		for _, multiKey := range keys {

			err := sd.DeleteValue(multiKey)
			if err != nil {
				fmt.Printf("delete kv error, %s\n", err)
			}
		}

		cData := value.([]byte)

		if len(cData) >= PackedValueSize {

			var pos = PackedValueSize

			err = sd.SetValue(key, cData[:pos], opt)
			if err != nil {
				return err
			}

			index := 1
			for len(cData)-pos > PackedValueSize {

				multiKey := fmt.Sprintf("%s.%d", key, index)
				err = sd.SetValue(multiKey, cData[pos:pos+PackedValueSize])
				if err != nil {
					return err
				}
				pos += PackedValueSize
				index++
			}

			if len(cData)-pos > 0 {
				multiKey := fmt.Sprintf("%s.%d", key, index)
				err = sd.SetValue(multiKey, cData[pos:])
				if err != nil {
					return err
				}
			}

			return nil

		} else {
			return sd.SetValue(key, cData)
		}

	} else {
		return sd.SetValue(key, value)
	}
}

func SafeGetValue(sd Discovery, key string, valuePtr interface{}, opt Option) error {

	rg := sd.(rawGetter)

	var (
		finalData []byte
		err       error
	)

	if opt.Compress {

		keys, err := getMultiKey(rg, key)
		if err != nil {
			return err
		}

		var data []byte
		for _, multiKey := range keys {

			var partData []byte
			err := rg.GetValueDirect(multiKey, &partData)
			if err != nil {
				return err
			}

			data = append(data, partData...)
		}
		finalData = data

		//
		//finalData, err = util.DecompressBytes(data)
		//
		//if err != nil {
		//	return err
		//}

	} else {
		err = rg.GetValueDirect(key, &finalData)

		if err != nil {
			return err
		}

	}

	reflect.ValueOf(valuePtr).Elem().Set(reflect.ValueOf(finalData))

	return nil
}
