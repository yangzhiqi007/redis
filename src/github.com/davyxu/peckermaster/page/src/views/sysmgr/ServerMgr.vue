<template>
  <div>
    <el-card style="margin-top: 10px">
      <el-row>
        <el-button type="success" @click="onPopAdd">添加</el-button>
      </el-row>

      <el-table :data="tabData" stripe style="width: 100%">
        <el-table-column prop="Name" label="名称" width="100"/>
        <el-table-column prop="Address" label="地址" width="300" />
        <el-table-column  label="操作" width="200">
          <template slot-scope="scope">
            <el-button type="warning" round size="small" @click="onModify(scope.row)">修改</el-button>
            <el-button type="danger"  round size="small" @click="onDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" :close-on-click-modal='closeOnModel' width="30%">
      <el-form :model="dialogForm">
        <el-form-item label="名称:" label-width="100px">
          <el-input type="text" size="medium" placeholder="请输入服务器名称" v-model="dialogForm.name"></el-input>
        </el-form-item>
        <el-form-item label="地址:" label-width="100px">
          <el-input type="text" size="medium" placeholder="请输入服务器地址" v-model="dialogForm.address"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible=false">取 消</el-button>
        <el-button type="primary" @click="onAdd">确 定</el-button>
      </span>
    </el-dialog>
  </div>

</template>

<script>
  import {SystemRequest} from '@/common/api';

export default {
  name: 'ServerMgr',
  data () {
    return {
      closeOnModel: false,
      tabData: [],
      dialogTitle: "",
      dialogVisible: false,
      dialogForm: {
        name: "",
        address: "",
      },
    }
  },

  mounted(){

    this.queryServerList()
  },

  methods:{

    onModify( row ){
      this.dialogTitle = "修改服务器信息"
      this.dialogForm.name = row.Name
      this.dialogForm.address = row.Address
      this.dialogForm.originName = row.Name
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

        SystemRequest('post', "/server_delete", {
          Name: row.Name,
        }, function (res) {

          self.queryTask()
        })

      })

    },

    onPopAdd( ){
      this.dialogTitle = "添加服务器"
      this.dialogVisible=true
    },

    onAdd(){
      this.dialogVisible=false
      let self = this

      if ( this.dialogTitle === "修改服务器信息" ){
        SystemRequest('post', "/server_update", {

          'Name': self.dialogForm.originName,
          'Server':{
            'Name': self.dialogForm.name,
            'Address': self.dialogForm.address,
          },

        }, function (res) {

          self.queryServerList()
        })
      }else{
        SystemRequest('post', "/server_add", {
          'Name': self.dialogForm.name,
          'Address': self.dialogForm.address,
        }, function (res) {

          self.queryServerList()
        })
      }

    },


    queryServerList( ){

      let self = this
      SystemRequest('post', "/server_querydetail", null, function (res) {

        self.serverList = []
        self.tabData = []
        res.data.forEach(d =>{
          self.tabData.push(d)
          self.serverList.push({
            Name: d,
            ID: d,
          })
        })
      })

    },


  }
}
</script>
