<template>
  <el-upload
    ref="upload"
    class="upload-demo"
    action="/v1/api/cluster/create"
    name="config_file"
    :headers="config"
    :data="formData"
    :on-preview="handlePreview"
    :on-remove="handleRemove"
    :on-change="handleChange"
    :before-remove="beforeRemove"
    :file-list="fileList"
    :auto-upload="false"
  >
    <el-button slot="trigger" size="small" type="primary">选取集群文件</el-button>
    <div slot="tip" class="el-upload__tip">只能上传yaml文件，且不超过500kb</div>
  </el-upload>
</template>

<script>
export default {
  props: {
    // 这个 prop 属性接收父组件传递进来的值
    formData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      limit: 1,
      fileList: [],
      uploadfile: []
    }
  },
  computed: {
    config() {
      return {
        Authorization: this.$store.state.user.token
      }
    }
  },
  methods: {
    submitUpload() {
      if (!this.formData.cluster_name || !this.formData.cluster_type) {
        this.$message({
          message: '请填写完整信息',
          type: 'error'
        })
        return false
      }
      console.log(this.uploadfile.length)
      if (this.uploadfile.length < 1) {
        this.$message({
          message: '请上传文件',
          type: 'error'
        })
        return false
      }
      this.$refs.upload.submit()
      this.clearList()
      return true
    },
    handleRemove(file, fileList) {
      console.log(file, fileList)
    },
    handlePreview(file) {
      console.log(file)
    },
    handleChange(file, fileList) {
      this.uploadfile = fileList
    },
    beforeRemove(file, fileList) {
      return this.$confirm(`确定移除 ${file.name}？`)
    },
    clearList() {
      this.fileList = []
      this.uploadfile = []
    }
  }
}
</script>
