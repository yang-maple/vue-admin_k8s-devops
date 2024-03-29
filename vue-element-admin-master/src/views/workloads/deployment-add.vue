<template>
  <div>
    <div style="padding:20px;font-size: 24px;">
      <el-page-header title="返回" content="创建无状态副本" @back="goBack" />
    </div>
    <div v-if="succeed" class="form-create">
      <el-form ref="ruleForm" :model="ruleForm" :rules="rules" label-width="100px" label-position="right" size="small">
        <el-form-item label="名称" prop="name">
          <el-input v-model="ruleForm.name" />
        </el-form-item>
        <el-form-item label="命名空间" prop="namespace">
          <el-select
            v-model="ruleForm.namespace"
            remote
            popper-append-to-body
            :remote-method="remoteMethod"
            :loading="loading"
            placeholder="请选择命名空间"
            @focus="remoteMethod('')"
          >
            <el-option v-for="(item,index) in options" :key="index" :label="item.label" :value="item.namespace" />
          </el-select>
        </el-form-item>
        <el-form-item label="标签" prop="labels">
          <ltable ref="lable" />
        </el-form-item>
        <el-form-item label="副本" prop="replicas">
          <el-input v-model.number="ruleForm.replicas" />
        </el-form-item>
        <el-form-item label="容器名称" prop="container.container_name">
          <el-input v-model="ruleForm.container.container_name" />
        </el-form-item>
        <el-form-item label="镜像" prop="container.image">
          <el-input v-model="ruleForm.container.image" />
        </el-form-item>
        <!-- <el-form-item label="Cpu需求">
                <el-input v-model="ruleForm.container.cpu" placeholder="单位为 m"></el-input>m
            </el-form-item>
            <el-form-item label="Memory需求">
                <el-input v-model="ruleForm.container.memory" placeholder="单位为 Mi"></el-input>Mi
            </el-form-item> -->
        <el-form-item label="端口名称" prop="container.container_port.port_name">
          <el-input v-model="ruleForm.container.container_port.port_name" />
        </el-form-item>
        <el-form-item label="端口" prop="container.container_port.container_port">
          <el-input v-model.number="ruleForm.container.container_port.container_port" />
        </el-form-item>
        <el-form-item label="协议" prop="container.container_port.protocol">
          <el-radio-group v-model="ruleForm.container.container_port.protocol">
            <el-radio label="TCP" />
            <el-radio label="UDP" />
          </el-radio-group>
        </el-form-item>

        <el-form-item label="健康检查" prop="delivery">
          <el-switch v-model="ruleForm.health_check" />
        </el-form-item>
        <el-form-item
          v-show="ruleForm.health_check == true"
          label="路径"
          :prop="ruleForm.health_check ? 'health_path' : 'nocheck'"
        >
          <el-input v-model="ruleForm.health_path" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitForm('ruleForm')">立即创建</el-button>
          <el-button @click="resetForm('ruleForm')">立即重置</el-button>
          <el-button type="danger" @click="dialogcreatens = false, resetForm('ruleForm')">
            取消
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
      duration: 3,
      ruleForm: {
        name: '',
        namespace: '',
        replicas: null,
        labels: '',
        container:
        {
          container_name: '',
          image: '',
          cpu: '0m',
          memory: '0Mi',
          container_port:
          {
            port_name: '',
            container_port: null,
            protocol: ''
          }
        },
        value: '',
        health_check: false,
        health_path: ''
      },
      options: [],
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'blur' }
        ],
        namespace: [
          { required: true, message: '请选择命名空间', trigger: 'change' }
        ],
        replicas: [
          { required: true, message: '请输入副本数量', trigger: 'change' },
          { type: 'number', min: 0, max: 99, message: '副本数范围为0-99', trigger: 'change' }
        ],
        'container.container_name': [
          { required: true, message: '请输入容器名称', trigger: 'blur' }
        ],
        'container.image': [
          { required: true, message: '请输入镜像名称', trigger: 'blur' }
        ],
        'container.container_port.container_port': [
          { required: true, message: '请输入端口', trigger: 'change' },
          { type: 'number', min: 1, max: 65535, message: '端口范围为1-65535', trigger: 'change' }
        ],
        'container.container_port.protocol': [
          { required: true, message: '请选择一种协议', trigger: 'change' }
        ],
        health_path: [
          { required: true, message: '请输入健康检查路径', trigger: 'blur' }
        ],
        nocheck: []
      }
    }
  },
  unmounted() {
    clearInterval(this.timer)
  },
  methods: {
    goBack() {
      this.$router.push({
        path: '/workload/deployment'
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
      this.$store.dispatch('deployment/createDeployment', this.ruleForm).then((res) => {
        this.succeed = false
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
      }, 200)
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
