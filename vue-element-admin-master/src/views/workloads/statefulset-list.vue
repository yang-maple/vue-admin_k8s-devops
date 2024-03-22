<template>
  <el-row style="padding-bottom: 10px;">
    <el-col :span="24">
      <el-card shadow="always" style="width: 100%;">
        <span>
          <div>
            <svg class="header-icon" aria-hidden="true">
              <use xlink:href="#icon-statefulset" />
            </svg>
            <span
              style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
            >有状态副本集</span>
            <br>
            <span style="font-size: 12px;color: #79879c!important">有状态副本集（StatefulSet）是指在 Kubernetes
              集群中运行的软件应用程序，这些应用程序在执行期间会保持一定的状态，以便在多个 Pod 之间维护稳定的身份标识和持久性存储。
            </span>
          </div>
        </span>
      </el-card>
    </el-col>
  </el-row>
</template>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
