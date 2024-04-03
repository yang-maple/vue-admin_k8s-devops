<template>
  <div>
    <div style="padding:20px;font-size: 24px;">
      <el-page-header title="返回" content="创建配置字典" @back="goBack" />
    </div>
    <div v-if="succeed" class="form-create">
      <el-form ref="ruleForm" :model="ruleForm" :rules="rules" label-width="100px" label-position="right" size="small">
        <el-form-item label="资源名称" prop="name">
          <el-input v-model="ruleForm.name" />
        </el-form-item>
        <el-form-item label="命名空间" prop="namespace">
          <el-select
            v-model="ruleForm.namespace"
            remote
            style="width: 300px;"
            popper-append-to-body
            :remote-method="remoteMethod"
            :loading="loading"
            placeholder="请选择命名空间"
            @focus="remoteMethod('')"
          >
            <el-option v-for="(item,index) in options" :key="index" :label="item.label" :value="item.namespace" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务标签" prop="labels">
          <ltable ref="lable" />
        </el-form-item>
        <el-form-item label="数据字典" prop="datamap">
          <div class="datamap">
            <ltable ref="datamap" buttontype="textarea" :buttonsize="true" />
          </div>
        </el-form-item>

        <el-form-item label="禁止修改">
          <el-switch v-model="ruleForm.modify" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm('ruleForm')">立即创建</el-button>
          <el-button @click="resetForm('ruleForm')">立即重置</el-button>
          <el-button type="danger" @click="goBack">
            立即返回
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div v-else>
      <el-result icon="success" title="创建成功">
        <template slot="subTitle">
          {{ text }}
        </template>
        <template slot="extra">
          <el-button type="primary" size="medium" @click="goBack">立即返回</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script>
import ltable from '../components/LabelTable'
export default {
  components: { ltable },
  data() {
    return {
      succeed: true,
      loading: false,
      text: '',
      timer: null,
      duration: 2,
      options: [],
      ruleForm: {
        name: null,
        namespace: null,
        labels: null,
        data: null,
        modify: false
      },
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'blur' }
        ],
        namespace: [
          { required: true, message: '请选择命名空间', trigger: 'change' }
        ]
      }
    }
  },
  unmounted() {
    clearInterval(this.timer)
  },
  methods: {
    goBack() {
      this.$router.push({
        path: '/profiles/configmap'
      })
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          const array = this.$refs.lable.tableData
          const label = {}
          array.forEach(element => {
            label[element.key] = element.value
          })
          this.ruleForm.labels = label

          const data = this.$refs.datamap.tableData
          const datamap = {}
          data.forEach(element => {
            datamap[element.key] = element.value
          })
          this.ruleForm.data = datamap
          this.create()
        } else {
          return false
        }
      })
    },
    create() {
      this.$store.dispatch('configmap/createConfigmap', this.ruleForm).then((res) => {
        this.succeed = false
        this.text = `3秒后自动返回`
        this.timer = setInterval(() => {
          const tmp = this.duration--
          this.text = `${tmp}秒后自动返回`
          if (tmp <= 0) {
            // 清除掉定时器
            clearInterval(this.timer)
            this.goBack()
          }
        }, 1000)
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    },
    remoteMethod(query) {
      this.loading = true
      this.options = []
      setTimeout(() => {
        this.loading = false
        this.$store.dispatch('namespace/getNamespacelist', {}).then((res) => {
          res.data.item.forEach(v => {
            this.options.push({ 'namespace': v.name, 'label': v.name })
          })
        })
      }, 500)
    }
  }
}
</script>

<style>
.form-create{
  width: 600px;
}

.label{
  width: 150px;
  padding-right: 10px;
}
</style>
