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
          <el-button type="info" round icon="el-icon-refresh" />
          <el-button type="info" round @click="dialogcreatens = true">
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
            <el-dropdown trigger="click" @command="handleEdit(scope.row.name)">
              <el-button type="text">
                <i class="el-icon-s-operation" :style="{ fontSize: '18px' }" />
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    icon="el-icon-edit"
                    command="a"
                  >编辑</el-dropdown-item>
                  <el-dropdown-item
                    icon="el-icon-delete"
                    @click="messageboxOperate(scope.row, 'delete')"
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
  </div>
</template>

<script scoped>
import Vue from 'vue'
export default {
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
      dialogFormVisible: false,
      dialogcreatens: false,
      form: {
        name: '',
        region: '',
        date1: '',
        date2: '',
        delivery: false,
        type: [],
        resource: '',
        desc: '',
        newnamespaces: ''
      },
      // filter_name: 'nginx',
      total: 0,
      page_size: [1, 10, 20, 50, 100],
      // limit: 10,
      // page: 1,
      content: ''
    }
  },
  created() {
    this.getNamespaces()
  },
  methods: {
    handleEdit(name) {
      this.editdialog = true
      console.log(name)
      this.$store.dispatch('namespace/getNamespaceDeatil', name).then((res) => {
        console.log(res)
      })
      // this.$ajax({
      //   method: 'get',
      //   url: '/namespaces/detail',
      //   params: {
      //     namespace_name: name
      //   }
      // }).then((res) => {
      //   this.activeName = 'json'
      //   this.aceConfig.lang = 'json'
      //   this.content = JSON.stringify(res.data.detail, null, 2)
      // }).catch((res) => {
      //   console.log(res.data)
      // })
    },
    // handleUpdate() {
    //   let data = this.content
    //   let datajson = {}
    //   try {
    //     if (this.aceConfig.lang == 'yaml') {
    //       data = JSON.stringify(yaml.load(data), null, 2)
    //     }
    //     datajson = JSON.parse(data)
    //   } catch (e) {
    //     this.$message({
    //       showClose: true,
    //       message: '格式错误,请检查格式',
    //       type: 'error'
    //     })
    //     return
    //   }
    //   this.$ajax.put(
    //     '/namespaces/update',
    //     {
    //       data: datajson
    //     }
    //   ).then((res) => {
    //     this.dialogFormVisible = false,
    //     this.$message({
    //       showClose: true,
    //       message: res.msg,
    //       type: 'success'
    //     })
    //   }).catch((res) => {
    //     this.$message({
    //       showClose: true,
    //       message: res.reason,
    //       type: 'error'
    //     })
    //   })
    //   this.reload()
    // },
    handleDelete(row) {
      this.$ajax({
        method: 'delete',
        url: '/namespaces/delete',
        params: {
          namespace_name: row.name
        }
      }
      ).then((res) => {
        this.$message({
          showClose: true,
          message: res.msg,
          type: 'warning'
        })
      }).catch((res) => {
        console.log(res)
      })
      this.reload()
    },
    messageboxOperate(row, name) {
      this.$confirm(`是否${name}实例${row.name}`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.handleDelete(row)
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    getNamespaces() {
      this.$store.dispatch('namespace/getNamespace', this.serachInfo).then((res) => {
        this.namespacesItem = res.data.item
        this.total = res.data.total
        for (let i = 0; i < res.data.item.length; i++) {
          this.maxitem.push(3)
        }
      })

      // this.$ajax({
      //   method: 'get',
      //   url: '/namespaces/list',
      //   params: {
      //     page: this.page,
      //     limit: this.limit,
      //     filter_name: this.filter_name
      //   }
      // }).then((res) => {
      //   this.total = res.data.total
      //   this.namespacesItem = res.data.item
      //   for (let i = 0; i < res.data.item.length; i++) {
      //     this.maxitem.push(3)
      //   }
      // }).catch((res) => {
      //   console.log(res)
      // })
    },
    createNamespace() {
      this.$ajax.post(
        '/namespaces/create',
        {
          namespace_name: this.form.newnamespaces
        }
      ).then((res) => {
        this.$message({
          message: res.msg,
          type: 'success'
        })
      }).catch((res) => {
        console.log(res)
      })
      this.reload()
    },
    Refresh() {
      setTimeout(() => {
        this.reload()
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
    handle(row) {
      this.$router.push({
        path: '/namespaces/namespaces_detail',
        name: 'namespaces详情',
        query: {
          namespace_name: row.name
        }
      })
    },
    // yamlFormat() {
    //   if (this.aceConfig.lang == 'yaml') {
    //     return
    //   }
    //   this.aceConfig.lang = 'yaml'
    //   this.content = yaml.dump(JSON.parse(this.content))
    // },
    // jsonFormat() {
    //   if (this.aceConfig.lang == 'json') {
    //     return
    //   }
    //   this.aceConfig.lang = 'json'
    //   this.content = JSON.stringify(yaml.load(this.content), null, 2)
    // },
    // handleClick(tab) {
    //   if (tab.props.name == 'yaml') {
    //     this.yamlFormat()
    //     return
    //   }
    //   this.jsonFormat()
    // },
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

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
