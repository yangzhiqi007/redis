<template>
  <div>
    <el-card style="margin-top: 10px">
      <span>服务器:</span>
      <el-select placeholder="请选择服务器" v-model="currServer" @change="onServerChanged" style="margin-right: 10px">
        <el-option
          v-for="item in serverList"
          :key="item.ID"
          :label="item.Name"
          :value="item.ID">
        </el-option>
      </el-select>

      <span>过滤器:</span>
      <el-select placeholder="进程组过滤" v-model="procGroup" @change="onProcGroupChanged" style="margin-right: 10px">
        <el-option
          v-for="item in procGroupList"
          :key="item.ID"
          :label="item.Name"
          :value="item.ID">
        </el-option>
      </el-select>
      <el-table ref="multipleTable" :data="tabData" stripe style="width: 100%" @selection-change="onProcSelectChanged">
        <el-table-column type="selection" width="55"/>
        <el-table-column prop="Name" label="进程" width="300"/>
        <el-table-column prop="Status" label="状态" width="100" />
        <el-table-column prop="Desc" label="信息" width="300"/>
        <el-table-column  label="日志" width="200">
          <template slot-scope="scope">
            <el-button type="info" round size="small" @click="onTail(scope.row)">日志</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-row style="margin-top: 10px">
        <el-button type="warning" :loading="opratingProc" @click="onRestart">重启</el-button>
        <el-button type="success" :loading="opratingProc" @click="onStart">启动</el-button>
        <el-button type="danger"  :loading="opratingProc" @click="onStop">停止</el-button>
      </el-row>
    </el-card>
    <el-dialog title="进程日志" :visible.sync="dialogVisible" width="50%" :close-on-click-modal='closeOnModel'>
      <el-input type="textarea" :autosize="{ minRows: 4}" readonly v-model="procLog"></el-input>
    </el-dialog>

  </div>

</template>

<script>
  import {SystemRequest} from '@/common/api';

  const allGroup = {Name:"所有进程", ID: "@all" }
  const noneGroup = {Name:"全局组", ID: "@global" }

export default {
  name: 'ProcMgr',
  data () {
    return {
      closeOnModel: false,
      tabData: [],
      rawData: [],
      procGroupList: [],
      procGroup: "",
      selectedProc: [],
      currServer: "",
      serverList: [],
      opratingProc: false,
      procLog: "",
      dialogVisible: false,
    }
  },

  mounted(){

    this.queryServerList()
  },

  methods:{

    onTail( row ){
      this.tailProc(row.Name)
    },

    onProcSelectChanged( selection ){
      this.selectedProc = selection
    },

    onServerChanged( ){
      this.procGroup = ""
      this.queryProcList()
    },

    onProcGroupChanged( ){
      this.refreshList( )

      // 保存选择信息
      this.storeSelProc()
    },

    onStart( ){
      this.execCommand("start")
    },

    onRestart( ){
      this.execCommand("restart")
    },

    onStop( ){
      this.execCommand("stop")
    },

    refreshList( ){

      let self = this

      // 显示所有
      if (self.procGroup === allGroup.ID){

        self.tabData = self.rawData

      }else{

        self.tabData =[]
        self.rawData.forEach(d =>{

          let namePair = d.Name.split(":")
          if (namePair.length === 2){

            let groupName = namePair[0]

            // 只放入与过滤器匹配的
            if (groupName === self.procGroup ) {
              self.tabData.push(d)
            }


          }else if ( self.procGroup === noneGroup.ID ){
            self.tabData.push(d)
          }

        })

      }

    },

    queryProcList( ){
      let self = this

      self.rawData = []
      SystemRequest('post', "/proc_query", {
        "Server": self.currServer,
      }, function (res) {

        if (res.data.length === 0 ){
          self.procGroupList = []
          self.rawData = []
          self.tabData = []
          self.selectedProc = []
          self.procGroup = ""
          return
        }

        self.procGroupList = [allGroup, noneGroup]

        res.data.forEach(d => {

          // 保存原始数据
          self.rawData.push(d)

          let namePair = d.Name.split(":")

          if (namePair.length === 2){

            let groupName = namePair[0]

            let found = false
            for ( let k in self.procGroupList ){
              if (self.procGroupList[k].ID === groupName){
                found = true
                break
              }
            }

            if (!found){
              // 更新组列表
              self.procGroupList.push({
                Name: groupName,
                ID: groupName,
              })
            }
          }
        })

        if (self.procGroup === ""){
          self.procGroup = allGroup.ID

          // 保存选择信息
          self.storeSelProc()
        }


        self.refreshList()

      })
    },

    storeSelProc( ){

      let data = {
        SelectedServer: this.currServer,
        SelectedGroup: this.procGroup,
      }

      sessionStorage.setItem("selproc", JSON.stringify(data))

    },

    queryServerList( ){

      let self = this
      SystemRequest('post', "/server_query", null, function (res) {

        self.serverList = []
        res.data.forEach(d =>{
          self.serverList.push({
            Name: d,
            ID: d,
          })
        })

        if (self.firstLoad === undefined){

          // 第一次加载时, 查看之前保存的选择信息
          let selproc =sessionStorage.getItem("selproc")
          if (typeof selproc === 'string' ){
            selproc = JSON.parse(selproc)
            self.currServer = selproc.SelectedServer
            self.procGroup  = selproc.SelectedGroup

            // 根据信息马上刷新列表
            self.queryProcList()
          }

          self.firstLoad = true
        }


      })

    },


    tailProc( procName ){

      let self = this
      SystemRequest('post', "/proc_log", {
        "Server": self.currServer,
        "Name": procName,
      }, function (res) {

        self.procLog = res.data
        self.dialogVisible = true

      })
    },


    execCommand( cmd ){

      if (this.selectedProc.length === 0 ){
        return
      }

      let nameList = []


      // 选择指定组(非全局或所有)时, 如果全选, 用快速组重启
      if (this.procGroup !== allGroup.ID &&
          this.procGroup !== noneGroup.ID &&
          this.selectedProc.length === this.tabData.length){

        nameList = [ this.procGroup+ ":"]

      }else{

        this.selectedProc.forEach(d =>{
          nameList.push(d.Name)
        })
      }

      let self = this
      self.opratingProc = true
      SystemRequest('post', "/proc_ctl", {
        "Server": self.currServer,
        "NameList": nameList,
        "Command": cmd,
      }, function (res) {
        self.opratingProc = false
        self.queryProcList( )
      })
    },

  }
}
</script>
