<template>
  <div>
    <div style="padding:20px;font-size: 24px;">
      <el-page-header title="返回" content="创建路由资源" @back="goBack" />
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
            class="add-select"
            popper-append-to-body
            :remote-method="remoteMethod"
            :loading="loading"
            placeholder="请选择命名空间"
            @focus="remoteMethod('')"
          >
            <el-option v-for="(item,index) in options" :key="index" :label="item.label" :value="item.namespace" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源标签" prop="labels">
          <ltable ref="lable" />
        </el-form-item>
        <el-form-item label="控制器" prop="ingress_class_name">
          <el-select
            v-model="ruleForm.ingress_class_name"
            remote
            class="add-select"
            popper-append-to-body
            :remote-method="remoteClass"
            :loading="loading"
            placeholder="请选择路由控制器"
            @focus="remoteClass('')"
          >

            <el-option v-for="item in ingressclasslist" :key="item.value" :label="item.label" :value="item.value" />

          </el-select></el-form-item>
        <el-form-item label="访问域名" prop="port_name">
          <el-input v-model="ruleForm.rules[0].host" />
        </el-form-item>
        <el-form-item label="访问路径" prop="rules[0].http_ingress_rule_values[0].path">
          <el-input v-model="ruleForm.rules[0].http_ingress_rule_values[0].path" />
        </el-form-item>
        <el-form-item label="路径类型" prop="rules[0].http_ingress_rule_values[0].path_type">
          <el-select
            v-model="ruleForm.rules[0].http_ingress_rule_values[0].path_type"
            class="add-select"
            placeholder="请选择路径类型"
          >
            <el-option v-for="item in path_type" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务名称" prop="rules[0].http_ingress_rule_values[0].service_name">
          <el-select
            v-model="ruleForm.rules[0].http_ingress_rule_values[0].service_name"
            filterable
            remote
            style="width: 300px;"
            reserve-keyword
            placeholder="请输入服务名称"
            :remote-method="remoteService"
            :loading="loading"
            @focus="remoteService('')"
          >
            <el-option v-for="(item,index) in servicelist" :key="index" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务端口" prop="rules[0].http_ingress_rule_values[0].service_port">
          <el-input v-model.number="ruleForm.rules[0].http_ingress_rule_values[0].service_port" />
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
      ingressclasslist: [],
      servicelist: [],
      path_type: [
        {
          'value': 'ImplementationSpecific',
          'label': 'ImplementationSpecific'
        },
        {
          'value': 'Exact',
          'label': 'Exact'
        },
        {
          'value': 'Prefix',
          'label': 'Prefix'
        }
      ],
      options: [],
      ruleForm: {
        name: null,
        namespace: null,
        labels: null,
        ingress_class_name: null,
        rules: [
          {
            host: null,
            http_ingress_rule_values: [
              {
                path: null,
                path_type: null,
                service_name: null,
                service_port: null
              }
            ]
          }
        ]
      },
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'change' }
        ],
        namespace: [
          { required: true, message: '请选择命名空间', trigger: 'change' }
        ],
        ingress_class_name: [
          { required: true, message: '请选择路由类型', trigger: 'change' }
        ],
        resource: [
          { message: '请选择一种协议', trigger: 'change' }
        ],
        type: [
          { required: true, message: '请选择一种服务类型', trigger: 'change' }
        ],
        'rules[0].http_ingress_rule_values[0].path': [
          { required: true, message: '请输入服务访问路径', trigger: 'change' }
        ],
        'rules[0].http_ingress_rule_values[0].path_type': [
          { required: true, message: '请选择路径类型', trigger: 'change' }
        ],
        'rules[0].http_ingress_rule_values[0].service_name': [
          { required: true, message: '请输入服务名称', trigger: 'change' }
        ],
        'rules[0].http_ingress_rule_values[0].service_port': [
          { required: true, message: '请输入服务访问端口', trigger: 'change' },
          { type: 'number', min: 1, max: 65535, message: '请输入1-65535数字' }
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
        path: '/loadbalancing/ingress'
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
          this.create()
        } else {
          return false
        }
      })
    },
    create() {
      this.$store.dispatch('ingress/createIngress', this.ruleForm).then((res) => {
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
    },
    remoteClass(query) {
      this.loading = true
      this.ingressclasslist = []
      setTimeout(() => {
        this.loading = false
        this.$store.dispatch('ingressclass/getIngressClass', {}).then((res) => {
          res.data.item.forEach(v => {
            this.ingressclasslist.push({ 'value': v.name, 'label': v.name })
          })
        })
      }, 500)
    },
    remoteService(query) {
      this.loading = true
      this.servicelist = []
      setTimeout(() => {
        this.loading = false
        this.$store.dispatch('service/getService', { filter_name: query }).then((res) => {
          res.data.item.forEach(v => {
            this.servicelist.push({ 'value': v.name, 'label': v.name })
          })
        })
      }, 500)
    }
  }
}
</script>

<style lang="scss">
@import "~@/styles/anticon.scss";
.form-create{
  width: 600px;
}

.label{
  width: 150px;
  padding-right: 10px;
}
</style>
