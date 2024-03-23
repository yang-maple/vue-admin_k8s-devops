<template>
  <div>
    <el-row style="padding-bottom: 10px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-DaemonSet" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >守护进程集</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">守护进程（DaemonSet）是一种用来管理 Pod 的控制器，它可以确保在每个 Node
                上运行一个或多个 Pod 副本。
              </span>
            </div>
          </span>
        </el-card>
      </el-col>
    </el-row>
    <div>
      <!--只需要双向绑定代码块即可-->
      <el-tabs type="border-card">
        <el-tab-pane label="用户管理">
          <CodeEdit ref="editor" v-model="content" language="json" /></el-tab-pane>
        <el-tab-pane label="配置管理">
          <CodeEdit ref="editor" v-model="content" language="yaml" /></el-tab-pane>
        <el-tab-pane label="角色管理">角色管理</el-tab-pane>
        <el-tab-pane label="定时任务补偿">定时任务补偿</el-tab-pane>
      </el-tabs>
      <button @click="getCode">获取代码</button>
    </div>
  </div>
</template>

<script>
import CodeEdit from '../components/Editor'
export default {
  components: { CodeEdit },
  data() {
    return {
      content: '11111' // 代码块
    }
  },
  methods: {
    // 获取代码
    getCode() {
      console.log(this.$refs.editor.saveconnect)
    },
    getchildren(data) {
      console.log(data)
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
