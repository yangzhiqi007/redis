package pbmeta

import (
	"fmt"

	pbprotos "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

type DescriptorPool struct {

	// 缓存所有的文件描述符
	fileMap map[string]*FileDescriptor

	// 对应的文件描述符数组
	fileArray []*FileDescriptor

	// 全局消息表, key: .package.msgname
	msgMap map[string]*Descriptor

	// 全局枚举表, key: .package.enumname
	enumMap map[string]*EnumDescriptor
}

func NewDescriptorPool(fds *pbprotos.FileDescriptorSet) *DescriptorPool {

	self := &DescriptorPool{
		fileMap: make(map[string]*FileDescriptor),
		msgMap:  make(map[string]*Descriptor),
		enumMap: make(map[string]*EnumDescriptor),
	}

	self.fileArray = make([]*FileDescriptor, len(fds.GetFile()))

	for i, def := range fds.GetFile() {

		newFD := newFileDescriptor(def, self)
		self.fileMap[def.GetName()] = newFD
		self.fileArray[i] = newFD

		// 注册到全局
		for m := 0; m < newFD.MessageCount(); m++ {
			msg := newFD.Message(m)

			self.registerMessage(newFD, msg)
		}

		for e := 0; e < newFD.EnumCount(); e++ {
			en := newFD.Enum(e)

			self.registerEnum(newFD, en)
		}
	}

	return self
}

// 获取文件描述符
func (self *DescriptorPool) FileByName(name string) *FileDescriptor {
	if v, ok := self.fileMap[name]; ok {
		return v
	}

	return nil
}

// 取文件描述符
func (self *DescriptorPool) File(index int) *FileDescriptor {
	return self.fileArray[index]
}

// 文件描述符数量
func (self *DescriptorPool) FileCount() int {
	return len(self.fileArray)
}

func (self *DescriptorPool) registerMessage(fd *FileDescriptor, md *Descriptor) {
	key := fmt.Sprintf("%s.%s", fd.Define.GetPackage(), md.Name())
	self.msgMap[key] = md

	//log.Debugf("reg msg %s", key)

}

func (self *DescriptorPool) registerEnum(fd *FileDescriptor, ed *EnumDescriptor) {
	key := fmt.Sprintf("%s.%s", fd.Define.GetPackage(), ed.Name())
	//log.Debugf("reg enum %s", key)
	self.enumMap[key] = ed
}

func normalizeFullName(name string) string {

	if len(name) == 0 {
		return ""
	}

	if name[0:1] == "." {
		return name[1:]
	}

	return name

}

// 通过全名取消息
func (self *DescriptorPool) MessageByFullName(fullname string) *Descriptor {

	if v, ok := self.msgMap[fullname]; ok {
		return v
	}

	return nil
}

// 通过全名取枚举
func (self *DescriptorPool) EnumByFullName(fullname string) *EnumDescriptor {
	if v, ok := self.enumMap[fullname]; ok {
		return v
	}

	return nil
}
