<template>
  <div>
    <el-row style="padding-bottom: 5px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-deployment-copy" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >无状态副本集</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">
                无状态副本集（Deployment）是指在 Kubernetes 集群中运行的不保存任何状态或数据软件应用程序，它们将所有需要保存的状态存储在外部系统中，并且在重启或迁移时可以很快恢复。
              </span>
            </div>
          </span>
        </el-card>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="6">
        <div class="header-grid-content">
          <el-select
            v-model="serachInfo.namespace"
            filterable
            placeholder="命名空间（默认ALL）"
            clearable
            @change="changenamespace()"
          >
            <el-option-group v-for="(group,index) in nslist" :key="index" :label="group.label">
              <el-option
                v-for="item in group.options"
                :key="item.namespace"
                :label="item.label"
                :value="item.namespace"
              />
            </el-option-group>
          </el-select>
        </div>
      </el-col>
      <el-col :span="14">
        <div class="header-grid-content">
          <el-input
            v-model="serachInfo.filter_name"
            placeholder="请输入搜索资源的名称"
            clearable
            @clear="getDeployment()"
            @keyup.native.enter="getDeployment()"
          >
            <template #prepend>
              <el-button icon="el-icon-search" @click="getDeployment()" />
            </template>
          </el-input>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="header-grid-content" style="text-align: right;">
          <el-button type="info" icon="el-icon-refresh" round @click="getDeployment()" />
          <el-button type="info" round @click="createDeployment()">
            创建
          </el-button>
        </div>
      </el-col>
    </el-row>

    <div class="table-bg-purple">
      <el-table :data="deploymentItem" :header-cell-style="{ background: '#e6e7e9' }" size="small" empty-text="抱歉，暂无数据">
        <el-table-column label="名称" width="200">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <svg class="table-icon-small" aria-hidden="true">
                <use xlink:href="#icon-deployment-copy" />
              </svg>
              <span style="margin-left:3px" size="small">{{
                scope.row.name
              }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="命名空间" prop="namespaces" width="150" />
        <el-table-column label="标签">
          <template #default="scope">
            <div v-for="(v, k, index) in scope.row.labels " :key="k">
              <div v-if="index < maxitem[scope.$index]">
                <el-tag type="info" style="margin-left: 5px;" size="small" effect="plain" round>
                  {{ k }}:{{ v }}</el-tag>
              </div>
            </div>
            <div v-if="scope.row.labels == null">---</div>
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
        <el-table-column label="容器组数量" prop="pods" width="100" align="center" />
        <el-table-column label="状态" prop="status" width="100">
          <template #default="scope">
            <div v-if="scope.row.status == 'Running'" style="color: #00a2ae">
              <el-tag type="success" size="small" effect="plain" round>
                Running
              </el-tag>
            </div>
            <div v-if="scope.row.status != 'Running'" style="color: #f56c6c">
              <el-tag type="danger" size="small" effect="plain" round>
                {{ scope.row.status }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="镜像" width="200">
          <template #default="scope">
            <div v-for="(v, k) in scope.row.image " :key="k">{{ v }}<br></div>
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
                    :command="beforeCommand('edit', scope.row)"
                  >编辑</el-dropdown-item>
                  <el-dropdown-item
                    icon="el-icon-document"
                    :command="beforeCommand('replica', scope.row)"
                  >
                    副本
                  </el-dropdown-item>
                  <el-dropdown-item
                    icon="el-icon-refresh-right"
                    :command="beforeCommand('restart', scope.row)"
                  >重启</el-dropdown-item>
                  <el-dropdown-item
                    icon="el-icon-delete"
                    :command="beforeCommand('delete', scope.row)"
                  >删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
      <el-row>
        <el-col>
          <div class="demo-pagination-block" style="padding-bottom: 5px;">
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
          </div>
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
        namespace: '',
        filter_name: '',
        limit: 10,
        page: 1
      },
      maxitem: [],
      detailnamespace: null,
      editdialog: false,
      deploymentItem: [],
      content: '',
      nslist: [],
      total: 0,
      page_size: [1, 10, 20, 50, 100],
      ruleForm: {
        name: '',
        namespace: '',
        replicas: 1,
        labels: null,
        container:
        {
          container_name: '',
          image: '',
          cpu: '0m',
          memory: '0Mi',
          container_port:
          {
            port_name: '',
            container_port: null,
            protocol: ''
          }
        },
        health_check: false,
        health_path: ''
      },
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'blur' }
        ],
        namespace: [
          { required: true, message: '请选择命名空间', trigger: 'change' }
        ],
        labels: [
          { required: true, message: '请输入标签', trigger: 'blur' }
        ],
        replica: [
          { required: true, message: '请输入副本数量', trigger: 'change' }
        ],
        'container.container_name': [
          { required: true, message: '请输入容器名称', trigger: 'blur' }
        ],
        'container.image': [
          { required: true, message: '请输入镜像名称', trigger: 'blur' }
        ],
        'container.container_port.container_port': [
          { required: true, message: '请输入端口', trigger: 'change' },
          { type: 'number', min: 1, max: 65535, message: '端口范围为1-65535', trigger: 'change' }
        ],
        'container.container_port.protocol': [
          { required: true, message: '请选择一种协议', trigger: 'change' }
        ],
        health_path: [
          { required: true, message: '请输入健康检查路径', trigger: 'blur' }
        ],
        nocheck: []
      }
    }
  },
  mounted() {
    this.getnsselect()
  },
  created() {
    this.getDeployment()
  },
  methods: {
    getDeployment() {
      this.serachInfo.namespace = this.$store.state.namespace.namespace
      this.$store.dispatch('deployment/getDeployment', this.serachInfo).then((res) => {
        this.total = res.data.total
        this.deploymentItem = res.data.item
        for (let i = 0; i < res.data.item.length; i++) {
          this.maxitem.push(3)
        }
      })
    },
    Refresh() {
      setTimeout(() => {
        this.getDeployment()
      }, 1000)
    },
    async getnsselect() {
      if (this.nslist.length === 0) {
        this.nslist.push({
          label: '',
          options: [
            {
              namespace: '',
              label: '全部空间'
            }
          ]
        })
        this.$store.dispatch('namespace/getNamespacelist', {}).then((res) => {
          this.nslist.push({
            label: '',
            options: []
          })
          res.data.item.forEach(v => {
            this.nslist[1].options.push({ 'namespace': v.name, 'label': v.name })
          })
        })
      }
    },
    changenamespace() {
      this.$store.dispatch('namespace/setNamespce', this.serachInfo.namespace)
      this.getDeployment()
    },
    handleSizeChange(limit) {
      this.serachInfo.limit = limit
      this.serachInfo.page = 1
      this.getDeployment()
    },
    handleCurrentChange(page) {
      this.serachInfo.page = page
      this.getDeployment()
    },
    handleEdit(namespace, name) {
      this.$store.dispatch('deployment/getDeploymentDetail', { namespace: namespace, name: name }).then((res) => {
        this.editdialog = true
        this.detailnamespace = namespace
        this.content = yaml.dump(res.data)
      })
    },
    handleReplica(namespace, name, replicas) {
      this.$store.dispatch('deployment/modifyDeployment', {
        deploy_name: name,
        namespace: namespace,
        replicas: replicas
      }).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
      })
      setTimeout(() => {
        this.getDeployment()
      }, 1000)
    },
    handleRestart(namespace, name) {
      this.$confirm(`是否重启实例 ${name}`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('deployment/restartDeployment', {
          deploy_name: name,
          namespace: namespace
        }).then((res) => {
          this.$message({
            showClose: true,
            message: res.msg,
            type: 'success'
          })
          setTimeout(() => {
            this.getDeployment()
          }, 1000)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消重启'
        })
      })
    },
    handleDelete(namespace, name) {
      this.$confirm(`是否删除实例 ${name}`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('deployment/deleteDeployment', {
          namespace: namespace,
          deploy_name: name
        }).then((res) => {
          this.$message({
            showClose: true,
            message: res.msg,
            type: 'warning'
          })
          setTimeout(() => {
            this.getDeployment()
          }, 1000)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
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
      this.$store.dispatch('deployment/updateDeployment', { namespace: this.detailnamespace, data: datajson }).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
        this.editdialog = false
        setTimeout(() => {
          this.getDeployment()
        }, 1000)
      })
    },
    beforeCommand(command, name) {
      return {
        'command': command,
        'row': name
      }
    },
    handleCommand(command) {
      switch (command.command) {
        case 'edit':
          this.handleEdit(command.row.namespaces, command.row.name)
          break
        case 'delete':
          this.handleDelete(command.row.namespaces, command.row.name)
          break
        case 'restart':
          this.handleRestart(command.row.namespaces, command.row.name)
          break
        case 'replica':
          this.changeReplica(command.row)
          break
      }
    },
    handleClose() {
      this.$confirm('未保存退出, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(_ => {
        this.editdialog = false
      }).catch(_ => { })
    },
    changeReplica(row) {
      this.$prompt('修改的副本数为', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputType: 'tel',
        inputPattern: /^([1-9]|[1-9][0-9]|0?[0-9])$/,
        inputErrorMessage: '副本的数量范围 0-99'
      }).then(({ value }) => {
        this.handleReplica(row.namespaces, row.name, parseInt(value))
        console.log(value)
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '取消输入'
        })
      })
    },
    messageboxOperate(row, name) {
      this.$confirm(`是否${name}实例${row.name}`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        if (name === '重启') {
          this.handleRestart(row.namespaces, row.name)
        } else if (name === '删除') {
          this.handleDelete(row.namespaces, row.name)
        }
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    createDeployment() {
      this.$router.push({
        path: '/workload/create'
      })
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
.dialog-edit {
  .el-dialog {
    .el-dialog__body {
      padding: 20px;
      border-bottom: 1px solid #ededed;
    }
  }
}
</style>
