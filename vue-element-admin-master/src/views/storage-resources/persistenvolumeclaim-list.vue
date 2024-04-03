<template>
  <div>
    <el-row style="padding-bottom: 10px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-cunchujuan" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >持久卷声明</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">持久卷声明（Persistent Volume Claim,
                PVC）是Kubernetes中用于定义和管理应用程序对其它持久卷（Persistent Volume, PV）的请求和绑定的组件。</span>
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
            <el-option-group v-for="(group, index) in nslist" :key="index" :label="group.label">
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
            @clear="getPVClaim()"
            @keyup.native.enter="getPVClaim()"
          >
            <template #prepend>
              <el-button icon="el-icon-search" @click="getPVClaim()" />
            </template>
          </el-input>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="header-grid-content" style="text-align: right;">
          <el-button type="info" icon="el-icon-refresh" round @click="getPVClaim()" />
          <el-button type="primary" icon="el-icon-plus" @click="createResource()">创建</el-button>
        </div>
      </el-col>
    </el-row>

    <div class="table-bg-purple">
      <el-table :data="claimsItem" :header-cell-style="{ background: '#e6e7e9' }" size="small" empty-text="抱歉，暂无数据">
        <el-table-column label="名称" width="240">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <svg class="table-icon-small" aria-hidden="true">
                <use xlink:href="#icon-cunchujuan" />
              </svg>
              <span style="margin-left:3px" size="small">{{
                scope.row.name
              }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="命名空间" prop="namespace" width="150" />
        <el-table-column label="容量" prop="claim" width="70" align="center" />
        <el-table-column label="持久卷" prop="volume">
          <template #default="scope">{{ scope.row.volume ? scope.row.volume : '---' }}</template>
        </el-table-column>
        <el-table-column label="状态" prop="status" width="100">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <span v-if="scope.row.status == 'Bound'" slot="reference">
                <el-tooltip placement="bottom" effect="light"><template #content>Bound</template>
                  <i class="dotClass" style="background-color: springgreen" /></el-tooltip>
              </span>
              <span v-if="scope.row.status != 'Bound'" slot="reference">
                <el-tooltip placement="bottom" effect="light"><template #content> Available </template>
                  <i class="dotClass" style="background-color: red" /></el-tooltip>
              </span>
              <span size="small" style="margin-left: 10px">{{
                scope.row.status
              }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="访问模式" width="150">
          <template #default="scope">
            <div v-for="(v, k) in scope.row.access_modes " :key="k">{{ v }}<br></div>
          </template>
        </el-table-column>
        <el-table-column label="存储类型" prop="storage_class" width="100">
          <template #default="scope">{{ scope.row.storage_class ? scope.row.storage_class : '---' }}</template>
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
      detailnamespace: null,
      editdialog: false,
      claimsItem: [],
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
    this.getPVClaim()
  },
  methods: {
    getPVClaim() {
      this.serachInfo.namespace = this.$store.state.namespace.namespace
      this.$store.dispatch('pvclaim/getPVClaim', this.serachInfo).then((res) => {
        this.total = res.data.total
        this.claimsItem = res.data.item
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
        case 'delete':
          this.handleDelete(command.row.namespace, command.row.name)
          break
      }
    },
    createResource() {
      this.$router.push({
        path: 'persistenvolumeclaim/create'
      })
    },
    beforeCommand(command, name) {
      return {
        'command': command,
        'row': name
      }
    },
    changenamespace() {
      this.$store.dispatch('namespace/setNamespce', this.serachInfo.namespace)
      this.getPVClaim()
    },
    handleSizeChange(limit) {
      this.serachInfo.limit = limit
      this.serachInfo.page = 1
      this.getPVClaim()
    },
    handleCurrentChange(page) {
      this.serachInfo.page = page
      this.getPVClaim()
    },
    handleEdit(namespace, name) {
      this.$store.dispatch('pvclaim/getPVClaimDetail', { namespace: namespace, name: name }).then((res) => {
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
        this.$store.dispatch('pvclaim/deletePVClaim', {
          namespace: namespace,
          persistent_volume_claim_name: name
        }).then((res) => {
          this.$message({
            showClose: true,
            message: res.msg,
            type: 'warning'
          })
          setTimeout(() => {
            this.getPVClaim()
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
      this.$store.dispatch('pvclaim/updatePVClaim', { namespace: this.detailnamespace, data: datajson }).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
        this.editdialog = false
        setTimeout(() => {
          this.getPVClaim()
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
