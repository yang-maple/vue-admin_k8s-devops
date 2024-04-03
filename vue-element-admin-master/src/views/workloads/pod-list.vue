<template>
  <div> <el-row style="padding-bottom: 5px;">
          <el-col :span="24">
            <el-card shadow="always" style="width: 100%;">
              <span>
                <div>
                  <svg class="header-icon" aria-hidden="true">
                    <use xlink:href="#icon-gitpod" />
                  </svg>
                  <span
                    style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
                  >容器组</span>
                  <br>
                  <span style="font-size: 12px;color: #79879c!important">
                    容器组（Pod）是 Kubernetes 应用程序的基本执行单元，是您创建或部署的 Kubernetes 对象模型中最小和最简单的单元。
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
            @clear="getPod()"
            @keyup.native.enter="getPod()"
          >
            <template #prepend>
              <el-button icon="el-icon-search" @click="getPod()" />
            </template>
          </el-input>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="header-grid-content" style="text-align: right;">
          <el-button type="info" icon="el-icon-refresh" round @click="getPod()" />
          <!-- <el-button type="info" round @click="createDeployment()">
            创建
          </el-button> -->
        </div>
      </el-col>
    </el-row>

    <div class="table-bg-purple">
      <el-table :data="podItem" :header-cell-style="{ background: '#e6e7e9' }" size="small" empty-text="抱歉，暂无数据">
        <el-table-column label="名称" width="200">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <span style="font-size: 14px;color: blueviolet;">
                {{ scope.row.name }}
              </span>
            </div>
            <div v-for="(v, k) in scope.row.image " :key="k">
              <svg class="icon-pods" aria-hidden="true">
                <use xlink:href="#icon-Docker" />
              </svg>
              <span style="margin-left: 5px">{{ v }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="命名空间" prop="namespace" width="150" />
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
        <el-table-column label="状态" prop="status" width="150">
          <template #default="scope">
            <div v-if="scope.row.status == 'Running' || scope.row.status == 'Completed'" style="color: #00a2ae">
              <el-tag type="success" size="small" effect="plain" round>
                {{ scope.row.status }}
              </el-tag>
            </div>
            <div v-if="scope.row.status != 'Running' && scope.row.status != 'Completed'" style="color: #f56c6c">
              <el-tag type="danger" size="small" effect="plain" round>
                {{ scope.row.status }}
              </el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="节点" prop="node" width="110">
          <template #default="scope">{{ scope.row.node ? scope.row.node : '---' }}</template>
        </el-table-column>
        <el-table-column label="重启次数" prop="restart" width="70" />
        <el-table-column label="请求CPU" prop="cpu" width="70" />
        <el-table-column label="请求内存" prop="memory" width="70" />

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
                    :command="beforeCommand('log', scope.row)"
                  >日志</el-dropdown-item>
                  <el-dropdown-item
                    icon="el-icon-connection"
                    :command="beforeCommand('shell', scope.row)"
                  >连接</el-dropdown-item>
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

<script>
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
      podItem: [],
      content: '',
      nslist: [],
      total: 0,
      page_size: [1, 10, 20, 50, 100]
    }
  },
  mounted() {
    this.getnsselect()
  },
  created() {
    this.getPod()
  },
  methods: {
    getPod() {
      this.serachInfo.namespace = this.$store.state.namespace.namespace
      this.$store.dispatch('pod/getPod', this.serachInfo).then((res) => {
        this.total = res.data.total
        this.podItem = res.data.item
        for (let i = 0; i < res.data.item.length; i++) {
          this.maxitem.push(3)
        }
      })
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
    handleCommand(command) {
      switch (command.command) {
        case 'edit':
          this.handleEdit(command.row.namespace, command.row.name)
          break
        case 'log':
          console.log('log')
          break
        case 'shell':
          console.log('shell')
          break
        case 'delete':
          this.handleDelete(command.row.namespace, command.row.name)
          break
      }
    },
    beforeCommand(command, name) {
      return {
        'command': command,
        'row': name
      }
    },
    changenamespace() {
      this.$store.dispatch('namespace/setNamespce', this.serachInfo.namespace)
      this.getPod()
    },
    handleSizeChange(limit) {
      this.serachInfo.limit = limit
      this.serachInfo.page = 1
      this.getPod()
    },
    handleCurrentChange(page) {
      this.serachInfo.page = page
      this.getPod()
    },
    handleEdit(namespace, name) {
      this.$store.dispatch('pod/getPodDetail', { namespace: namespace, name: name }).then((res) => {
        this.editdialog = true
        this.detailnamespace = namespace
        this.content = yaml.dump(res.data)
      })
    },
    handleDelete(namespace, name) {
      this.$confirm(`是否删除实例 ${name}`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('pod/deletePod', {
          namespace: namespace,
          pod_name: name
        }).then((res) => {
          this.$message({
            showClose: true,
            message: res.msg,
            type: 'warning'
          })
          setTimeout(() => {
            this.getPod()
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
      this.$store.dispatch('pod/updatePod', { namespace: this.detailnamespace, data: datajson }).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
        this.editdialog = false
        setTimeout(() => {
          this.getPod()
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
      }).catch(_ => { })
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

.icon-pods {
    width: 1.5em;
    height: 1.5em;
    padding-left: 3px;
    vertical-align: -0.65em;
    fill: currentColor;
    overflow: hidden;
}
</style>
