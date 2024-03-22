<template>
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
            <span
              style="font-size: 12px;color: #79879c!important"
            >存储类（StorageClass）支持动态卷供应，使管理员能够按需创建新的卷。</span>
          </div>
        </span>
      </el-card>
    </el-col>
  </el-row>
</template>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
