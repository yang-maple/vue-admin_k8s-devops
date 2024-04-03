<template>
  <div>
    <div style="padding:20px;font-size: 24px;">
      <el-page-header title="返回" content="创建服务资源" @back="goBack" />
    </div>
    <div v-if="succeed" class="form-create">
      <el-form ref="ruleForm" :model="ruleForm" :rules="rules" label-width="100px" label-position="right" size="small">
        <el-form-item label="服务名称" prop="name">
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
        <el-form-item label="服务标签" prop="labels">
          <ltable ref="lable" />
        </el-form-item>
        <el-form-item label="服务类型" prop="type">
          <el-select v-model="ruleForm.type" placeholder="请选择服务类型" class="add-select">
            <el-option v-for="item in typelist" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="端口名称" prop="container.container_port.port_name">
          <el-input v-model="ruleForm.service_ports.port_name" />
        </el-form-item>
        <el-form-item label="容器端口" prop="service_ports.port">
          <el-input v-model.number="ruleForm.service_ports.port" />
        </el-form-item>
        <el-form-item label="访问协议" prop="service_ports.protocol">
          <el-radio-group v-model="ruleForm.service_ports.protocol">
            <el-radio label="TCP" />
            <el-radio label="UDP" />
          </el-radio-group>
        </el-form-item>

        <el-form-item label="暴露端口" prop="service_ports.target_port">
          <el-input v-model.number="ruleForm.service_ports.target_port" />
        </el-form-item>
        <el-form-item v-show="ruleForm.type == 'NodePort'" label="访问端口" prop="service_ports.node_port">
          <el-input v-model.number="ruleForm.service_ports.node_port" />
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
      typelist: [
        {
          value: 'ClusterIP',
          label: 'ClusterIP'
        },
        {
          value: 'NodePort',
          label: 'NodePort'
        },
        {
          value: 'LoadBalancer',
          label: 'LoadBalancer'
        }
      ],
      ruleForm: {
        name: null,
        namespace: null,
        labels: null,
        type: null,
        service_ports: {
          port_name: null,
          port: null,
          protocol: null,
          target_port: null,
          node_port: null
        }
      },
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'change' }
        ],
        namespace: [
          { required: true, message: '请选择命名空间', trigger: 'change' }
        ],
        resource: [
          { message: '请选择一种协议', trigger: 'change' }
        ],
        type: [
          { required: true, message: '请选择一种服务类型', trigger: 'change' }
        ],
        'service_ports.port': [
          { required: true, message: '请输入容器端口', trigger: 'blur' },
          { type: 'number', message: '必须为数字值', trigger: 'blur' }
        ],
        'service_ports.target_port': [
          { type: 'number', message: '必须为数字值', trigger: 'blur' }
        ],
        'service_ports.node_port': [
          { type: 'number', min: 30000, max: 32767, message: '必须为数字值，范围为 30000-32767', trigger: 'blur' }
        ]
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
        path: '/loadbalancing/service'
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
      this.$store.dispatch('service/createService', this.ruleForm).then((res) => {
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
