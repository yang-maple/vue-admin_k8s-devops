<template>
  <div>
    <div style="padding:20px;font-size: 24px;">
      <el-page-header title="返回" content="创建持久卷声明" @back="goBack" />
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
            :loading="loading"
            class="add-select"
            placeholder="请选择命名空间"
            @focus="remoteMethod('')"
          >
            <el-option v-for="(item,index) in options" :key="index" :label="item.label" :value="item.namespace" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源标签" prop="labels">
          <ltable ref="lable" />
        </el-form-item>
        <el-form-item label="访问模式" prop="access_modes">
          <el-checkbox-group v-model="ruleForm.access_modes" style="width: 150%;">
            <el-checkbox label="ReadWriteOnce" name="access_modes" />
            <el-checkbox label="ReadOnlyMany" name="access_modes" />
            <el-checkbox label="ReadWriteMany" name="access_modes" />
            <el-checkbox label="ReadWriteOncePod" name="access_modes" />
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="存储容量" prop="storages">
          <div>
            <el-input v-model="ruleForm.storages" placeholder="" class="add-select">
              <template slot="append">
                <el-select v-model="storage_type" style="width: 100px">
                  <el-option label="Gi" value="Gi" />
                  <el-option label="Mi" value="Mi" />
                </el-select>
              </template>
            </el-input>
          </div>
        </el-form-item>
        <el-form-item label="存储类型">
          <el-select
            v-model="ruleForm.storage_classname"
            remote
            popper-append-to-body
            :remote-method="remoteClass"
            :loading="loading"
            allow-create
            filterable
            placeholder=" "
            class="add-select"
            @focus="remoteClass('')"
          >
            <el-option v-for="item in storageclasslist" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
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
      duration: 2,
      storage_type: 'Gi',
      ruleForm: {
        name: null,
        namespace: '',
        access_modes: [],
        labels: null,
        storage: null,
        storages: null,
        storage_classname: null
      },
      options: [],
      storageclasslist: [],
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'blur' }
        ],
        namespace: [
          { required: true, message: '请选择命名空间', trigger: 'change' }
        ],
        access_modes: [
          { required: true, message: '请选择访问模式', trigger: 'change' }
        ],
        storages: [
          { required: true, message: '请输入存储容量', trigger: 'change' }
        ],
        resource: [
          { message: '请选择一种协议', trigger: 'change' }
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
        path: '/storage/persistenvolumeclaim'
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
          this.ruleForm.storage = String(this.ruleForm.storages) + this.storage_type
          this.create()
        } else {
          return false
        }
      })
    },
    create() {
      this.$store.dispatch('pvclaim/createPVClaim', this.ruleForm).then((res) => {
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
      this.storageclasslist = []
      setTimeout(() => {
        this.loading = false
        this.$store.dispatch('storageclass/getStorageClass', {}).then((res) => {
          res.data.item.forEach(v => {
            this.storageclasslist.push({ 'value': v.name, 'label': v.name })
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
