<template>
  <div>
    <el-row style="padding-bottom: 5px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-jiqun" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >集群config配置</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">集在 Kubernetes 集群中，config 资源用于存储集群的配置信息，包括 API
                Server 的地址、认证信息、授权信息等。通过 config 资源，可以方便地连接和管理 Kubernetes 集群。
              </span>
            </div>
          </span>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="10">
      <el-col v-for="(value,index) in clusteritem" :key="index" :span="5" class="cluster-card-content">
        <el-card :body-style="{ padding: '0px' }" shadow="always">
          <img :src="imageUrls(value.type)" class="image">
          <div style="padding: 14px">
            <span>{{ value.cluster_name }}</span>
            <div class="bottom">
              <time class="time">{{ value.create_time }}</time>
              <div class="status">
                <el-tag
                  v-if="value.status == false"
                  :type="value.status == false ? 'danger' : 'success'"
                >未应用</el-tag>
                <el-tag
                  v-if="value.status == true"
                  :type="value.status == true ? 'success' : 'danger'"
                >已应用</el-tag>
              </div>
            </div>
            <div>
              <el-button
                v-if="value.status == false"
                type="primary"
                class="button"
                plain
                @click="change(value.cluster_name)"
              >应用</el-button>
              <el-button
                v-if="value.status == false"
                class="button"
                type="danger"
                @click="remove(value.cluster_name)"
              >删除</el-button>
              <el-button
                type="info"
                class="button"
                @click="edit(value.cluster_name)"
              >编辑</el-button>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-tooltip content=" 点 击 添 加 集 群 " placement="top" effect="light">
          <div class="config-add" @click="dialogaddcluster = true">
            <div class="avatar-uploader-icon">
              <i class="el-icon-plus" />
            </div>
          </div>
        </el-tooltip>
      </el-col>
    </el-row>

    <el-dialog :visible.sync="dialogaddcluster" title="新增集群" center @close="closeaddDialog">
      <el-form ref="ruleForm" :model="ruleForm" label-width="25%">
        <el-form-item label="集群名称" prop="cluster_name">
          <el-input v-model="ruleForm.cluster_name" autocomplete="off" style="width: 80%;" />
        </el-form-item>
        <el-form-item label="集群类型" prop="cluster_type">
          <el-select v-model="ruleForm.cluster_type" placeholder="请选择集群类型">
            <el-option label="私有云K8S" value="k8s" />
            <el-option label="阿里云ACK" value="ack" />
            <el-option label="腾讯云TKE" value="tke" />
            <el-option label="华为云CCE" value="cce" />
          </el-select>
        </el-form-item>
        <el-form-item label="集群config" prop="cluster_file">
          <file-upload ref="uploadfile" :form-data="ruleForm" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span>
          <el-button
            type="primary"
            size="small"
            style="margin-right:50px;padding-left:20px;padding-right:20px"
            @click="submitForm()"
          >上传并新建集群</el-button>
          <el-button
            type="danger"
            size="small"
            style="margin-right:50px;padding-left:20px;padding-right:20px"
            @click="closeaddDialog('ruleForm')"
          >
            取消
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog :visible.sync="dialogupdatecluster" title="编辑集群" center>
      <el-form :model="updataForm" label-width="25%">
        <el-form-item label="集群名称">
          <el-input v-model="updataForm.cluster_name" autocomplete="off" style="width: 80%;" />
        </el-form-item>
        <el-form-item label="集群类型">
          <el-select v-model="updataForm.type" placeholder="请选择集群类型" :value-key="updataForm.type">
            <el-option label="私有云K8S" value="k8s" />
            <el-option label="阿里云ACK" value="ack" />
            <el-option label="腾讯云TKE" value="tke" />
            <el-option label="华为云CCE" value="cce" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button
            type="primary"
            style="margin-right:50px;padding-left:20px;padding-right:20px"
            @click="updataCluster()"
          >更新</el-button>
          <el-button
            style="margin-right:50px;padding-left:20px;padding-right:20px"
            @click="dialogcluster = false"
          >
            取消
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>

</template>

<script>
import FileUpload from '../components/FileUpload.vue'
export default {
  name: 'ClusterConfig',
  components: {
    FileUpload
  },
  data() {
    return {
      dialogaddcluster: false,
      dialogupdatecluster: false,
      isDisabled: true,
      currentDate: new Date(),
      k8s: require('@/assets/images/k8s.png'),
      ack: require('@/assets/images/ack.png'),
      tke: require('@/assets/images/tke.png'),
      cce: require('@/assets/images/cce.png'),
      clusteritem: [],
      k8stotal: [],
      updataForm: {
        id: null,
        cluster_name: '',
        type: ''
      },
      ruleForm: {
        cluster_name: '',
        cluster_type: ''
      }
    }
  },
  created() {
    this.getCluster()
  },
  methods: {
    getCluster() {
      this.$store.dispatch('cluster/getCluster').then((res) => {
        this.clusteritem = res.data.item
      })
    },
    notify(status, info, msg) {
      this.$message({
        title: info,
        message: msg,
        type: status
      })
    },
    imageUrls(type) {
      switch (type) {
        case 'k8s':
          return this.k8s
        case 'ack':
          return this.ack
        case 'tke':
          return this.tke
        case 'cce':
          return this.cce
        default:
          return ''
      }
    },
    submitForm() {
      if (this.$refs.uploadfile.submitUpload()) {
        this.notify('success', '成功', '集群添加成功')
        this.resetForm('ruleForm')
        this.dialogaddcluster = false
        this.getCluster()
      }
    },
    change(clustername) {
      this.$store.dispatch('cluster/changeCluster', clustername).then((res) => {
        this.notify('success', '成功', '集群应用成功')
        this.getCluster()
      })
    },
    remove(clustername) {
      this.$store.dispatch('cluster/deleteCluster', clustername).then((res) => {
        this.notify('warning', '成功', '集群删除成功')
        this.getCluster()
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    },
    closeaddDialog() {
      this.dialogaddcluster = false
      this.$refs.uploadfile.clearList()
      this.resetForm('ruleForm')
    },
    edit(clustername) {
      this.dialogupdatecluster = true
      this.$store.dispatch('cluster/detailCluster', clustername).then(res => {
        this.updataForm.cluster_name = res.data.cluster_name
        this.updataForm.type = res.data.type
        this.updataForm.id = res.data.id
      })
    },
    updataCluster() {
      this.$store.dispatch('cluster/updateCluster', this.updataForm).then(res => {
        this.notify('success', '成功', '集群更新成功')
        this.dialogupdatecluster = false
        this.getCluster()
      })
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";

.cluster-card-content {
    position: relative;
    border-radius: 6px;
}
.time {
    font-size: 13px;
    color: #999;
  }

  .bottom {
    margin-top: 13px;
    line-height: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .button {
    padding: 5px;
    min-height: auto;
}

  .time {
    font-size: 12px;
    color: #999;
}
  .image {
    width: 100%;
  }

  .clearfix:before,
  .clearfix:after {
      display: table;
      content: "";
  }

  .clearfix:after {
      clear: both
  }

  .config-add {
    border: 1px black;
    border-radius: 4px;
    border-style: dashed;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);
    border-color: var(--el-color-primary);
    width: 180px;
    height: 210px;
    .avatar-uploader-icon {
    font-size: 48px;
    color: #8c939d;
    width: 178px;
    height: 210px;
    text-align: center;
    cursor: pointer;
    .el-icon-plus{
       padding-top: 70px;
    }
}
}
</style>
