<template>
  <el-row style="padding-bottom: 10px;">
    <el-col :span="24">
      <el-card shadow="always" style="width: 100%;">
        <span>
          <div>
            <svg class="header-icon" aria-hidden="true">
              <use xlink:href="#icon-a-nav_secretkey-copy" />
            </svg>
            <span
              style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
            >保密字典</span>
            <br>
            <span
              style="font-size: 12px;color: #79879c!important"
            >保密字典（Secret）是一种包含少量敏感信息的资源对象，例如密码、令牌、保密字典等，以键值对形式保存并且可以在容器组中使用。</span>
          </div>
        </span>
      </el-card>
    </el-col>
  </el-row>
</template>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
