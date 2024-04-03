<template>
  <div>
    <el-row style="padding-bottom: 10px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-cunchuleixingpeizhi1" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >存储类</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">存储类（StorageClass）支持动态卷供应，使管理员能够按需创建新的卷。</span>
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
            @clear="getStorageClass()"
            @keyup.native.enter="getStorageClass()"
          >
            <template #prepend>
              <el-button icon="el-icon-search" @click="getStorageClass()" />
            </template>
          </el-input>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="header-grid-content" style="text-align: right;">
          <el-button type="info" icon="el-icon-refresh" round @click="getStorageClass()" />
        </div>
      </el-col>
    </el-row>

    <div class="table-bg-purple">
      <el-table
        :data="storageClassItem"
        :header-cell-style="{ background: '#e6e7e9' }"
        size="small"
        tooltip-effect="light"
        empty-text="抱歉，暂无数据"
      >
        <el-table-column label="名称">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <svg class="table-icon-small" aria-hidden="true">
                <use xlink:href="#icon-cunchuleixingpeizhi1" />
              </svg>
              <span style="margin-left:3px" size="small">{{
                scope.row.name
              }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="提供者" prop="provisioner" width="280" />
        <custom-tag-table-column ref="tagTableColumn" :column="{ label: '参数', prop: 'parameters', width: '300px', tagLength: tagLengths }" />
        <el-table-column label="回收策略" prop="reclaim_policy" width="150" />
        <el-table-column label="创建时间" width="200" prop="age" />
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
import CustomTagTableColumn from '../components/CustomTagTableColumn'
import yaml from 'js-yaml'
export default {
  components: {
    CodeEdit,
    CustomTagTableColumn
  },
  data() {
    return {
      serachInfo: {
        filter_name: '',
        limit: 10,
        page: 1
      },
      editdialog: false,
      storageClassItem: [],
      tagLengths: 0,
      content: '',
      total: 0,
      page_size: [1, 10, 20, 50, 100]
    }
  },
  mounted() {
  },
  created() {
    this.getStorageClass()
  },
  methods: {
    getStorageClass() {
      this.$store.dispatch('storageclass/getStorageClass', this.serachInfo).then((res) => {
        this.total = res.data.total
        this.storageClassItem = res.data.item
        for (let i = 0; i < res.data.item.length; i++) {
          this.$refs.tagTableColumn.maxitem.push(3)
        }
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
    beforeCommand(command, name) {
      return {
        'command': command,
        'row': name
      }
    },
    handleSizeChange(limit) {
      this.serachInfo.limit = limit
      this.serachInfo.page = 1
      this.getStorageClass()
    },
    handleCurrentChange(page) {
      this.serachInfo.page = page
      this.getStorageClass()
    },
    handleEdit(name) {
      this.$store.dispatch('storageclass/getStorageClassDetail', { name: name }).then((res) => {
        this.editdialog = true
        this.content = yaml.dump(res.data)
      })
    },

    handleDelete(name) {
      this.$confirm(`是否删除实例 ${name}`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('storageclass/deleteStorageClass', {
          storageclass_name: name
        }).then((res) => {
          this.$message({
            showClose: true,
            message: res.msg,
            type: 'warning'
          })
          setTimeout(() => {
            this.getStorageClass()
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
      this.$store.dispatch('storageclass/updateStorageClass', { data: datajson }).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'success'
        })
        this.editdialog = false
        setTimeout(() => {
          this.getStorageClass()
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
