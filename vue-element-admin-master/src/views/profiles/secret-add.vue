<template>
  <div>
    <div style="padding:20px;font-size: 24px;">
      <el-page-header title="返回" content="创建保密字典" @back="goBack" />
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
            popper-append-to-body
            :remote-method="remoteMethod"
            style="width: 300px;"
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
        <el-form-item label="保密类型">
          <el-select v-model="ruleForm.type" placeholder=" " style="width: 300px;">
            <el-option v-for="item in secretType" :key="item.label" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="ruleForm.type == 'Opaque'" label="数据字典">
          <ltable ref="dataOpaque" buttontype="textarea" :buttonsize="true" />
        </el-form-item>
        <el-form-item v-if="ruleForm.type == 'kubernetes.io/tls'">
          <el-form-item
            label="证书"
            label-width="50px"
            style="padding-bottom:18px;"
            :prop="ruleForm.type == 'kubernetes.io/tls' ? 'tls_crt' : 'nocheck'"
          >
            <el-input
              v-model="ruleForm.tls_crt"
              type="textarea"
              :autosize="{ minRows: 4, maxRows: 6 }"
              style="width: 400px"
            />
          </el-form-item>
          <el-form-item
            label="私钥"
            label-width="50px"
            style="padding-bottom:18px"
            :prop="ruleForm.type == 'kubernetes.io/tls' ? 'tls_key' : 'nocheck'"
          >
            <el-input
              v-model="ruleForm.tls_key"
              type="textarea"
              :autosize="{ minRows: 4, maxRows: 6 }"
              style="width: 400px"
            />
          </el-form-item>
        </el-form-item>
        <el-form-item v-if="ruleForm.type == 'kubernetes.io/basic-auth'">
          <el-form-item
            label="账号"
            label-width="50px"
            style="padding-bottom:18px"
            :prop="ruleForm.type == 'kubernetes.io/basic-auth' ? 'username' : 'nocheck'"
          >
            <el-input v-model="ruleForm.username" type="text" style="width: 400px" />
          </el-form-item>
          <el-form-item
            label="密码"
            label-width="50px"
            style="padding-bottom:18px"
            :prop="ruleForm.type == 'kubernetes.io/basic-auth' ? 'password' : 'nocheck'"
          >
            <el-input v-model="ruleForm.password" type="password" show-password style="width: 400px" />
          </el-form-item>
        </el-form-item>
        <el-form-item v-if="ruleForm.type == 'kubernetes.io/dockerconfigjson'">
          <el-form-item
            label="仓库地址"
            label-width="80px"
            style="padding-bottom:18px"
            :prop="ruleForm.type == 'kubernetes.io/dockerconfigjson' ? 'docker_registry' : 'nocheck'"
          >
            <el-input v-model="ruleForm.docker_registry" style="width: 400px">
              <template slot="prepend">Http://</template>
            </el-input>
          </el-form-item>
          <el-form-item
            label="账号"
            label-width="80px"
            style="padding-bottom:18px"
            :prop="ruleForm.type == 'kubernetes.io/dockerconfigjson' ? 'docker_username' : 'nocheck'"
          >
            <el-input v-model="ruleForm.docker_username" style="width: 400px" />
          </el-form-item>
          <el-form-item
            label="密码"
            label-width="80px"
            style="padding-bottom:18px"
            :prop="ruleForm.type == 'kubernetes.io/dockerconfigjson' ? 'docker_password' : 'nocheck'"
          >
            <el-input v-model="ruleForm.docker_password" type="password" show-password style="width: 400px" />
          </el-form-item>
          <el-form-item label="邮箱" label-width="80px" style="padding-bottom:18px">
            <el-input v-model="ruleForm.docker_email" type="email" style="width: 400px" />
          </el-form-item>
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
      secretType: [{
        label: '默认',
        value: 'Opaque'
      },
      {
        label: 'TLS 信息',
        value: 'kubernetes.io/tls'
      },
      {
        label: '镜像服务信息',
        value: 'kubernetes.io/dockerconfigjson'
      },
      {
        label: '用户名和密码',
        value: 'kubernetes.io/basic-auth'
      }
      ],
      ruleForm: {
        name: null,
        namespace: '',
        labels: null,
        data: null,
        immutable: false,
        type: 'Opaque',
        tls_crt: null,
        tls_key: null,
        username: null,
        password: null,
        docker_registry: null,
        docker_username: null,
        docker_password: null,
        docker_email: null
      },
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'blur' }
        ],
        namespace: [
          { required: true, message: '请选择命名空间', trigger: 'change' }
        ],
        tls_crt: [
          { required: true, message: '请输入证书', trigger: 'blur' }
        ],
        tls_key: [
          { required: true, message: '请输入私钥', trigger: 'blur' }
        ],
        username: [
          { required: true, message: '请输入账号', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ],
        docker_registry: [
          { required: true, message: '请输入仓库地址', trigger: 'blur' }
        ],
        docker_username: [
          { required: true, message: '请输入账号', trigger: 'blur' }
        ],
        docker_password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ],
        nocheck: []
      },
      options: []
    }
  },
  unmounted() {
    clearInterval(this.timer)
  },
  methods: {
    goBack() {
      this.$router.push({
        path: '/profiles/secret'
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
          if (this.ruleForm.type === 'Opaque') {
            const data = this.$refs.dataOpaque.tableData
            const dataOpaque = {}
            data.forEach(element => {
              dataOpaque[element.key] = element.value
            })
            this.ruleForm.data = dataOpaque
          }

          this.create()
        } else {
          return false
        }
      })
    },
    create() {
      this.$store.dispatch('secret/createSecret', this.ruleForm).then((res) => {
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
