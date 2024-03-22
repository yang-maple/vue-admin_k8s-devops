<template>
  <el-empty
    description="还没有添加任何集群哦~"
    :image="emptyImageUrl"
    :image-size="400"
  >
    <el-button type="primary" @click="addCluster">按钮</el-button>
  </el-empty>
</template>

<script>
export default {
  data() {
    return {
      emptyImageUrl: require('@/assets/images/empty.png')
    }
  },
  methods: {
    addCluster() {
      this.$router.push('/clusterinfo/index')
    }
  }
}
</script>
