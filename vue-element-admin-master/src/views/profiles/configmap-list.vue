<template>
  <el-row style="padding-bottom: 10px;">
    <el-col :span="24">
      <el-card shadow="always" style="width: 100%;">
        <span>
          <div>
            <svg class="header-icon" aria-hidden="true">
              <use xlink:href="#icon-configmap" />
            </svg>
            <span
              style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
            >配置字典</span>
            <br>
            <span
              style="font-size: 12px;color: #79879c!important"
            >配置字典（ConfigMap）常用于存储工作负载所需的配置信息，许多应用程序会从配置文件、命令行参数或环境变量中读取配置信息。</span>
          </div>
        </span>
      </el-card>
    </el-col>
  </el-row>
</template>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
