<template>
  <div>
    <el-row style="padding-bottom: 5px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-cunchujuan1" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >持久卷</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">持久卷（Persistent
                Volume,PV）是Kubernetes中用于定义和管理集群中持久性存储资源的组件。</span>
            </div>
          </span>
        </el-card>
      </el-col>
    </el-row>

    <el-row>
      <el-col :span="20">
        <div class="header-grid-content">
          <el-input
            v-model="serachInfo.filter_name"
            placeholder="请输入搜索资源的名称"
            clearable
            @clear="getService()"
            @keyup.native.enter="getService()"
          >
            <template #prepend>
              <el-button icon="el-icon-search" @click="getService()" />
            </template>
          </el-input>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="header-grid-content" style="text-align: right;">
          <el-button type="info" icon="el-icon-refresh" round @click="getService()" />
          <el-button type="info" @click="handleCreate">
            创建
          </el-button>
        </div>
      </el-col>
    </el-row>

    <div class="table-bg-purple">
      <el-table :data="pvolumeItem" :header-cell-style="{ background: '#e6e7e9' }" size="small" empty-text="抱歉，暂无数据">
        <el-table-column label="名称" width="200">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <svg class="table-icon-small" aria-hidden="true">
                <use xlink:href="#icon-cunchujuan1" />
              </svg>
              <span style="margin-left:3px" size="small">{{
                scope.row.name
              }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="容量" prop="capacity.storage" width="100" align="center" />
        <el-table-column label="访问模式" prop="access_mode" width="120" />
        <el-table-column label="回收策略" prop="reclaim_policy" width="100" align="center" />
        <el-table-column label="状态" prop="status.phase" width="100">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <span v-if="scope.row.status.phase == 'Bound'" slot="reference">
                <el-tooltip placement="bottom" effect="light"><template #content>Bound</template>
                  <i class="dotClass" style="background-color: springgreen" /></el-tooltip>
              </span>
              <span v-if="scope.row.status.phase != 'Bound'" slot="reference">
                <el-tooltip placement="bottom" effect="light"><template #content> Available </template>
                  <i class="dotClass" style="background-color: red" /></el-tooltip>
              </span>
              <span size="small" style="margin-left: 10px">{{
                scope.row.status.phase
              }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="绑定声明" prop="claim">
          <template #default="scope">{{ scope.row.claim ? scope.row.claim : '---' }}</template>
        </el-table-column>
        <el-table-column label="存储类型" prop="storage_class" width="140">
          <template #default="scope">{{ scope.row.storage_class ? scope.row.storage_class : '---' }}</template>
        </el-table-column>
        <el-table-column label="创建时间" width="140" prop="age" />
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
        filter_name: '',
        limit: 10,
        page: 1
      },
      editdialog: false,
      pvolumeItem: [],
      content: '',
      total: 0,
      page_size: [1, 10, 20, 50, 100]
    }
  },
  mounted() {
  },
  created() {
    this.getPVolume()
  },
  methods: {
    getPVolume() {
      this.$store.dispatch('pvolume/getPVolume', this.serachInfo).then((res) => {
        this.total = res.data.total
        this.pvolumeItem = res.data.item
      })
    },
    handleCommand(command) {
      switch (command.command) {
        case 'edit':
          this.handleEdit(command.row.name)
          break
        case 'delete':
          this.handleDelete(command.row.name)
          break
      }
    },
    handleCreate() {
      this.$router.push({
        path: 'persistenvolume/create'
      })
    },
    beforeCommand(command, name) {
      return {
        'command': command,
        'row': name
      }
    },
    handleSizeChange(limit) {
      this.serachInfo.limit = limit
      this.serachInfo.page = 1
      this.getPVolume()
    },
    handleCurrentChange(page) {
      this.serachInfo.page = page
      this.getPVolume()
    },
    handleEdit(name) {
      this.$store.dispatch('pvolume/getPVolumeDetail', { name: name }).then((res) => {
        this.editdialog = true
        this.content = yaml.dump(res.data.detail)
      })
    },

    handleDelete(name) {
      this.$confirm(`是否删除实例 ${name}`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('pvolume/deletePVolume', {
          pvolume_name: name
        }).then((res) => {
          this.$message({
            showClose: true,
            message: res.msg,
            type: 'warning'
          })
          setTimeout(() => {
            this.getPVolume()
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
      this.$store.dispatch('pvolume/updatePVolume', { data: datajson }).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
        this.editdialog = false
        setTimeout(() => {
          this.getPVolume()
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
