<template>
  <div>
    <el-row style="padding-bottom: 5px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-namespace1" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >命名空间</span>
              <br>
              <span
                style="font-size: 12px;color: #79879c!important"
              >命名空间（Namespace）提供了一种逻辑分隔的方式，可以将集群中的资源划分到不同的命名空间中，以实现不同用户或应用之间的资源隔离和管理。
              </span>
            </div>
          </span>
        </el-card>
      </el-col>
    </el-row>
    <el-row>
      <el-col :span="8">
        <div class="header-grid-content">
          <el-input
            v-model="serachInfo.filter_name"
            placeholder="请输入搜索资源的名称"
            clearable
            @clear="getNamespaces()"
            @keyup.native.enter="getNamespaces()"
          >
            <template #prepend>
              <el-button icon="el-icon-search" @click="getNamespaces()" />
            </template>
          </el-input>
        </div>
      </el-col>
      <el-col :span="10">
        <div class="header-grid-content" style="height: 46px;" />
      </el-col>
      <el-col :span="6">
        <div class="header-grid-content" style="text-align: right;">
          <el-button type="info" round icon="el-icon-refresh" @click="Refresh()" />
          <el-button type="info" round @click="adddialog = true">
            创建
          </el-button>
        </div>
      </el-col>
    </el-row>
    <div class="table-bg-purple">
      <el-table :data="namespacesItem" :header-cell-style="{ background: '#e6e7e9' }" size="small" empty-text="抱歉，暂无数据">
        <el-table-column label="名称" width="200">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <svg class="table-icon-small" aria-hidden="true">
                <use xlink:href="#icon-namespace1" />
              </svg>
              <span style="margin-left:3px" size="small">{{
                scope.row.name
              }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="标签">
          <template #default="scope">
            <div v-for="(v, k, index) in scope.row.labels " :key="k">
              <div v-if="index < maxitem[scope.$index]">
                <el-tag type="info" style="margin-left: 5px;" size="small" effect="plain" round>
                  {{ k }}:{{ v }}</el-tag>
              </div>
            </div>
            <div v-if="scope.row.labels == null" align="center">---</div>
            <div v-if="scope.row.labels != null && Object.keys(scope.row.labels).length > 3"><el-button
              size="small"
              type="text"
              @click="showLabels(scope.$index)"
            >{{
              maxitem[scope.$index] == 3 ?
                '展开' : '收起'
            }}</el-button></div>
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="status" width="300">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <span v-if="scope.row.status == 'Active'" slot="reference">
                <el-tooltip placement="top" effect="light"><template #content>Active</template>
                  <i class="dotClass" style="background-color: springgreen" /></el-tooltip>
              </span>
              <span v-if="scope.row.status != 'Active'" slot="reference">
                <el-tooltip placement="bottom" effect="light"><template #content> UnActive </template>
                  <i class="dotClass" style="background-color: red" /></el-tooltip>
              </span>
              <span style="margin-left: 3px;color: green;"> {{ scope.row.status }} </span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="age" width="140" />
        <el-table-column label="操作" width="70" align="center">
          <template #default="scope">
            <el-dropdown trigger="click" @command="handleCommand">
              <el-button type="text">
                <i class="el-icon-s-operation" :style="{ fontSize: '18px' }" />
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    icon="el-icon-edit"
                    :command="beforeCommand('edit',scope.row.name)"
                  >编辑</el-dropdown-item>
                  <el-dropdown-item
                    icon="el-icon-delete"
                    :command="beforeCommand('delete',scope.row.name)"
                  >删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      <el-row>
        <el-col>
          <el-pagination
            :current-page.sync="serachInfo.page"
            :page-size.sync="serachInfo.limit"
            :page-sizes="page_size"
            background
            layout="total, sizes, prev, pager, next"
            :total="total"
            :pager-count="5"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </el-col>
      </el-row>
    </div>

    <div class="dialog-edit">
      <el-dialog :visible.sync="editdialog" :before-close="handleClose" title="编辑Yaml" top="5vh">
        <CodeEdit ref="editor" v-model="content" language="yaml" />
        <template #footer>
          <span class="dialog-footer">
            <el-button type="primary" @click="handleUpdate()">更新</el-button>
            <el-button @click="handleClose">
              取消
            </el-button>
          </span>
        </template>
      </el-dialog>
    </div>

    <div class="dialog-add">
      <el-dialog :visible.sync="adddialog" title="创建命名空间">
        <el-form :model="form" label-width="15%">
          <el-form-item label="名称">
            <el-input v-model="form.namespace" autocomplete="off" style="width: 90%;" />
            <div><span :font-family="'Arial, sans-serif'" style="font-size: 12px;">长度为 1 ~ 63
              个字符，只能包含数字、小写字母和中划线（-），且首尾只能是字母或数字</span></div>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button type="primary" @click="handleCreate">确认创建</el-button>
            <el-button type="danger" @click="adddialog = false">
              取消
            </el-button>
          </span>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script scoped>
import CodeEdit from '../components/Editor'
import yaml from 'js-yaml'
import Vue from 'vue'
export default {
  components: { CodeEdit },
  data() {
    return {
      serachInfo: {
        filter_name: '',
        limit: 10,
        page: 1
      },
      maxitem: [],
      namespacesItem: [],
      editdialog: false,
      adddialog: false,
      form: {
        namespace: '',
        newnamespace: ''
      },
      total: 0,
      page_size: [1, 10, 20, 50, 100],
      content: ''
    }
  },
  created() {
    this.getNamespaces()
  },
  methods: {
    handleEdit(name) {
      this.editdialog = true
      this.$store.dispatch('namespace/getNamespaceDeatil', name).then((res) => {
        this.content = yaml.dump(res.data.detail)
      })
    },
    handleUpdate() {
      let data = this.$refs.editor.saveconnect
      let datajson = {}
      try {
        data = JSON.stringify(yaml.load(data), null, 2)
        datajson = JSON.parse(data)
      } catch (e) {
        this.$message({
          showClose: true,
          message: '格式错误,请检查格式',
          type: 'error'
        })
        return
      }
      this.$store.dispatch('namespace/updateNamespace', datajson).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
        this.editdialog = false
        setTimeout(() => {
          this.getNamespaces()
        }, 1000)
      })
    },
    handleClose() {
      this.$confirm('未保存退出, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(_ => {
        this.editdialog = false
      }).catch(_ => {})
    },
    handleDelete(name) {
      this.$confirm(`是否删除命名空间 ${name}`, '提示', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('namespace/deleteNamespace', name).then((res) => {
          this.$message({
            type: 'success',
            message: res.msg
          })
          setTimeout(() => {
            this.getNamespaces()
          }, 1000)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    handleCommand(command) {
      switch (command.command) {
        case 'edit':
          this.handleEdit(command.name)
          break
        case 'delete':
          this.handleDelete(command.name)
          break
      }
    },
    beforeCommand(command, name) {
      return {
        'command': command,
        'name': name
      }
    },
    getNamespaces() {
      this.$store.dispatch('namespace/getNamespacelist', this.serachInfo).then((res) => {
        this.namespacesItem = res.data.item
        this.total = res.data.total
        for (let i = 0; i < res.data.item.length; i++) {
          this.maxitem.push(3)
        }
      })
    },
    handleCreate() {
      this.$store.dispatch('namespace/createNamespace', this.form.namespace).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
        this.adddialog = false
        this.form.namespace = ''
        setTimeout(() => {
          this.getNamespaces()
        }, 1000)
      })
    },
    Refresh() {
      setTimeout(() => {
        this.getNamespaces()
      }, 1000)
    },
    handleSizeChange(limit) {
      this.serachInfo.limit = limit
      this.serachInfo.page = 1
      this.getNamespaces()
    },
    handleCurrentChange(page) {
      this.serachInfo.page = page
      this.getNamespaces()
    },
    showLabels(index) {
      if (this.maxitem[index] === 3) {
        Vue.set(this.maxitem, index, 99)
      } else {
        Vue.set(this.maxitem, index, 3)
      }
    }
  }
}
</script>

<style lang="scss">
  @import "~@/styles/anticon.scss";
.dialog-edit{
  .el-dialog{
    .el-dialog__body{
      padding: 20px;
      border-bottom: 1px solid #ededed;
    }
  }
}

.dialog-add{
  .el-dialog{
    width: 650px;
  }
}

</style>
