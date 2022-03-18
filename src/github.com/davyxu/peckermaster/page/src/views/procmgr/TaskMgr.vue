<template>
  <div>
    <el-card style="margin-top: 10px">
      <el-row>
        <span>过滤器:</span>
        <el-select placeholder="任务分组" v-model="taskGroup" @change="onTaskGroupChanged" style="margin-right: 10px">
          <el-option
            v-for="item in taskGroupList"
            :key="item.ID"
            :label="item.Name"
            :value="item.ID">
          </el-option>
        </el-select>
        <el-button type="success" :loading="opratingProc" @click="onPopAdd">添加</el-button>
      </el-row>

      <el-table ref="multipleTable" :data="tabData" stripe style="width: 100%" @selection-change="onTaskSelectChanged">
        <el-table-column type="selection" width="55"/>
        <el-table-column prop="Name" label="名称" width="200"/>
        <el-table-column prop="Group" label="分组" width="200"/>
        <el-table-column prop="Server" label="服务器" width="100" />
        <el-table-column prop="Code" label="代码" width="500"/>
        <el-table-column  label="操作" width="200">
          <template slot-scope="scope">
            <el-button type="warning" round size="small" @click="onModify(scope.row)">修改</el-button>
            <el-button type="danger"  round size="small" @click="onDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-row style="margin-top: 10px">
        <el-button type="warning" :loading="opratingProc" @click="onExecute">执行</el-button>
      </el-row>
    </el-card>
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" :close-on-click-modal='closeOnModel' width="50%">
      <el-form :model="dialogForm">
        <el-form-item label="名称:" label-width="100px">
          <el-input type="text" size="medium" placeholder="请输入任务名称" v-model="dialogForm.name"></el-input>
        </el-form-item>
        <el-form-item label="服务器:" label-width="100px">
          <el-select placeholder="请选择服务器" filterable v-model="dialogForm.server">
            <el-option
              v-for="item in serverList"
              :key="item.ID"
              :label="item.Name"
              :value="item.ID">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="分组:" label-width="100px">
          <el-select placeholder="任务分组" filterable allow-create v-model="dialogForm.group" >
            <el-option
              v-for="item in rawTaskGroupList"
              :key="item.ID"
              :label="item.Name"
              :value="item.ID">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="代码:" label-width="100px">
          <el-input type="textarea" :autosize="{ minRows: 4}" placeholder="请输入shell代码" v-model="dialogForm.code"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible=false">取 消</el-button>
        <el-button type="primary" @click="onAdd">确 定</el-button>
      </span>
    </el-dialog>
    <el-dialog title="日志" :visible.sync="logDialogVisible" width="50%">
      <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 50}" readonly v-model="taskLog"></el-input>
    </el-dialog>
  </div>

</template>

<script>
  import {SystemRequest} from '@/common/api';
  const allGroup = {Name:"所有任务", ID: "@all" }

export default {
  name: 'TaskMgr',
  data () {
    return {
      closeOnModel: false,
      rawData: [],
      tabData: [],
      taskGroup:"",
      taskGroupList: [],
      rawTaskGroupList: [],
      serverList: [],
      selectedTask: [],
      opratingProc: false,
      dialogTitle: "",
      dialogVisible: false,
      dialogForm: {
        name: "",
        group: "",
        server: "",
        code: "",
      },
      taskLog:"",
      logDialogVisible: false,
    }
  },

  mounted(){

    this.queryServerList()
    this.queryTask()
  },

  methods:{
    onTaskSelectChanged( selection ){
      this.selectedTask = selection
    },

    onTaskGroupChanged( ){
      this.refreshTaskList()
    },



    onModify( row ){
      this.dialogTitle = "修改任务"
      this.dialogForm.name = row.Name
      this.dialogForm.group = row.Group
      this.dialogForm.server = row.Server
      this.dialogForm.code = row.Code
      this.dialogForm.originName = "" + row.Name
      this.dialogVisible=true
    },

    onDelete( row ){

      let self = this

      self.$confirm('确认?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        center: true
      }).then(() => {

        SystemRequest('post', "/task_delete", {
          Name: row.Name,
        }, function (res) {

          self.queryTask()
        })

      })

    },

    onPopAdd( ){
      this.refreshFilter( )
      this.dialogTitle = "创建任务"
      this.dialogVisible=true

    },

    onAdd(){
      this.dialogVisible=false
      let self = this


      if ( this.dialogTitle === "修改任务" ){

        SystemRequest('post', "/task_update", {
          'Name': self.dialogForm.originName,
          'Task':{
            'Name': self.dialogForm.name,
            'Group': self.dialogForm.group,
            'Server': self.dialogForm.server,
            'Code': self.dialogForm.code,
          }
        }, function (res) {

          self.queryTask()
        })
      }else{

        SystemRequest('post', "/task_create", {
          'Name': self.dialogForm.name,
          'Group': self.dialogForm.group,
          'Server': self.dialogForm.server,
          'Code': self.dialogForm.code,
        }, function (res) {

          self.queryTask()
        })
      }

    },


    onExecute( ){
      this.taskLog = ""
      if (this.selectedTask.length === 0 ){
        return
      }

      let self = this

      self.$confirm('确认?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        center: true
      }).then(() => {

        this.logDialogVisible = true

        this.execTask(0)

      })
    },

    execTask( index ){

      let self = this

      let d = self.selectedTask[index]

      self.taskLog += "执行任务: '" + d.Name + "' 服务器:'"+ d.Server + "'\n"

      SystemRequest('post', "/task_exec", {
        'Name':d.Name,
      }, function (res) {

        self.taskLog += "返回: \n" + res.data + "\n"

        if (index < self.selectedTask.length ){
          self.execTask(index + 1)
        }

      })

    },

    refreshTaskList( ){

      let self = this

      // 显示所有
      if (self.taskGroup === allGroup.ID){

        self.tabData = self.rawData

      }else{

        self.tabData =[]
        self.rawData.forEach(d =>{

            // 只放入与过滤器匹配的
            if (d.Group === self.taskGroup ) {
              self.tabData.push(d)
            }

        })

      }

    },

    refreshFilter( ){

      let self = this
      self.rawTaskGroupList = []

      self.rawData.forEach(d =>{

        let foundTask = false
        for ( let k in self.rawTaskGroupList ){
          if (self.rawTaskGroupList[k].ID === d.Group){
            foundTask = true
          }
        }

        if (!foundTask && d.Group !== ""){
          self.rawTaskGroupList.push({
            Name: d.Group,
            ID: d.Group,
          })
        }

      })
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
      })

    },

    queryTask( ){


      let self = this
      SystemRequest('post', "/task_query", null, function (res) {

        self.rawData = res.data

        self.refreshFilter( )

        self.taskGroupList = [ allGroup ]
        self.taskGroupList = self.taskGroupList.concat(self.rawTaskGroupList)



        if (self.taskGroup === ""){
          self.taskGroup = allGroup.ID
        }


        self.refreshTaskList()

      })

    },


  }
}
</script>
