<template>
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
            <span
              style="font-size: 12px;color: #79879c!important"
            >持久卷声明（Persistent Volume Claim, PVC）是Kubernetes中用于定义和管理应用程序对其它持久卷（Persistent Volume, PV）的请求和绑定的组件。</span>
          </div>
        </span>
      </el-card>
    </el-col>
  </el-row>
</template>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
