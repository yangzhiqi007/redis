package model

type Server struct {
	Address string
	Name    string
}

type ServerManager struct {
	ServerList []*Server
}

func (self *ServerManager) GetAddress(name string) string {

	for _, def := range self.ServerList {
		if def.Name == name {
			return def.Address
		}
	}
	return ""
}

func (self *ServerManager) AddServer(sv *Server) {
	self.ServerList = append(self.ServerList, sv)
}

func (self *ServerManager) UpdateServer(name string, insv *Server) {
	for index, task := range self.ServerList {
		if task.Name == name {
			self.ServerList[index] = insv
			break
		}
	}
}

func (self *ServerManager) DeleteServer(name string) {
	for index, task := range self.ServerList {
		if task.Name == name {
			self.ServerList = append(self.ServerList[:index], self.ServerList[index+1:]...)
			break
		}
	}
}

func (self *ServerManager) ServerByName(name string) *Server {
	for _, sv := range self.ServerList {
		if sv.Name == name {
			return sv
		}
	}

	return nil
}
