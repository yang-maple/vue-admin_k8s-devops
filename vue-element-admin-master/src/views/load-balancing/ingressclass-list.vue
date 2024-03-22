<template>
  <el-row style="padding-bottom: 10px;">
    <el-col :span="24">
      <el-card shadow="always" style="width: 100%;">
        <span>
          <div>
            <svg class="header-icon" aria-hidden="true">
              <use xlink:href="#icon-ingress_nginx" />
            </svg>
            <span
              style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
            >应用路由类</span>
            <br>
            <span
              style="font-size: 12px;color: #79879c!important"
            >应用路由类（IngressClass）一种用于控制Ingress资源的访问权限的机制</span>
          </div>
        </span>
      </el-card>
    </el-col>
  </el-row>
</template>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
