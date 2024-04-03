<template>
  <div>
    <div style="padding:20px;font-size: 24px;">
      <el-page-header title="返回" content="创建持久卷" @back="goBack" />
    </div>
    <div v-if="succeed" class="form-create">
      <el-form ref="ruleForm" :model="ruleForm" :rules="rules" label-width="100px" label-position="right" size="small">
        <el-form-item label="资源名称" prop="name">
          <el-input v-model="ruleForm.name" />
        </el-form-item>
        <el-form-item label="资源标签" prop="labels">
          <ltable ref="lable" />
        </el-form-item>

        <el-form-item label="资源容量" prop="storages">
          <el-input v-model.number="ruleForm.storages" class="add-select">
            <template slot="append">
              <el-select v-model="storage_type" placeholder="Select" style="width: 100px">
                <el-option label="Gi" value="Gi" />
                <el-option label="Mi" value="Mi" />
              </el-select>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="卷模式" prop="volume_mode">
          <el-radio-group v-model="ruleForm.volume_mode">
            <el-radio label="Filesystem" />
            <el-radio label="Block" />
          </el-radio-group>
        </el-form-item>

        <el-form-item label="访问模式" prop="access_mode">
          <el-checkbox-group v-model="ruleForm.access_mode" style="width: 150%;">
            <el-checkbox label="ReadWriteOnce" name="access_mode" />
            <el-checkbox label="ReadOnlyMany" name="access_mode" />
            <el-checkbox label="ReadWriteMany" name="access_mode" />
            <el-checkbox label="ReadWriteOncePod" name="access_mode" />
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="回收策略" prop="recycle_mode">
          <el-radio-group v-model="ruleForm.recycle_mode">
            <el-radio label="Retain" />
            <el-radio label="Recycle" />
            <el-radio label="Delete" />
          </el-radio-group>
        </el-form-item>
        <el-form-item label="存储类型" prop="class_name">
          <el-select
            v-model="ruleForm.class_name"
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
            <el-option v-for="item in classtypelist" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="资源类型" prop="type">
          <el-select v-model="ruleForm.type" placeholder=" " class="add-select">
            <el-option label="NFS" value="NFS" />
            <el-option label="HostPATH" value="HostPATH" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="ruleForm.type=='NFS'" label="服务地址" :prop="ruleForm.type == 'NFS' ? 'server' : 'nocheck'">
          <el-input v-model="ruleForm.server" autocomplete="off" />
        </el-form-item>
        <el-form-item label="路径" prop="path">
          <el-input v-model="ruleForm.path" autocomplete="off" />
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
      showEl: false,
      text: '',
      timer: null,
      duration: 2,
      storage_type: 'Gi',
      ruleForm: {
        name: '',
        labels: null,
        storage: null,
        storages: null,
        type: '',
        path: '',
        server: '',
        access_mode: [],
        volume_mode: '',
        recycle_mode: '',
        class_name: ''
      },
      classtypelist: [],
      options: [],
      rules: {
        name: [
          { required: true, message: '请输入资源名称', trigger: 'blur' }
        ],
        storages: [
          { required: true, message: '请输入资源容量', trigger: 'blur' },
          { type: 'number', message: '请输入数字', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择资源类型', trigger: 'blur' }
        ],
        volume_mode: [
          { required: true, message: '请选择一种卷模式', trigger: 'blur' }
        ],
        access_mode: [
          { required: true, message: '至少选择一种访问模式', trigger: 'blur' }
        ],
        server: [
          { required: true, message: '请输入服务地址', trigger: 'blur' }
        ],
        path: [
          { required: true, message: '请输入路径', trigger: 'blur' }
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
        path: '/storage/persistenvolume'
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
      this.$store.dispatch('pvolume/createPVolume', this.ruleForm).then((res) => {
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
    remoteClass(query) {
      this.loading = true
      this.classtypelist = []
      setTimeout(() => {
        this.loading = false
        this.$store.dispatch('storageclass/getStorageClass', {}).then((res) => {
          res.data.item.forEach(v => {
            this.classtypelist.push({ 'value': v.name, 'label': v.name })
          })
        })
      }, 500)
    }
  }
}
</script>

<style  lang="scss">
@import "~@/styles/anticon.scss";
.form-create{
  width: 600px;
}

.label{
  width: 150px;
  padding-right: 10px;
}
</style>
