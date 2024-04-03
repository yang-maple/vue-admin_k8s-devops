<template>
  <div>
    <el-row style="padding-bottom: 5px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-jisuanjiedian" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >集群节点</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">集群节点（Node）是指在 Kubernetes 集群中运行的 worker 节点。它们是实际运行
                Pod
                的物理或虚拟机器。
              </span>
            </div>
          </span>
        </el-card>
      </el-col>
    </el-row>
    <div class="table-bg-purple">
      <el-table
        :data="nodeItem"
        :header-cell-style="{ background: '#e6e7e9' }"
        style="width: 100%"
        size="small"
        :default-sort="{ prop: 'date', order: 'descending' }"
        empty-text="抱歉，暂无数据"
      >
        <el-table-column label="名称" width="200">
          <template #default="scope">
            <el-row>
              <el-col :span="4"><svg class="table-icon" aria-hidden="true">
                <use xlink:href="#icon-jisuanjiedian" />
              </svg></el-col>
              <el-col :span="20">
                <span style="margin-left: 5px;font-size: 14px;font-weight: bold;">{{
                  scope.row.name
                }}</span>
                <div style="margin-top: -5px;">
                  <span size="small" style="margin-left: 5px;">{{
                    scope.row.nodeIp
                  }}</span>
                </div>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
        <el-table-column label="标签">
          <template #default="scope">
            <div v-for="(v, k, index) in scope.row.labels" :key="k">
              <div v-if="index < maxitem[scope.$index]">
                <el-tag

                  type="info"
                  style="margin-left: 5px;"
                  size="small"
                  effect="plain"
                  round
                >
                  {{ k }}:{{ v }}
                </el-tag>
                <!-- <el-tooltip v-else effect="light" placement="top">
                  <el-tag type="info" style="margin-left: 5px;" size="small" effect="plain" round>
                    {{ k }}:{{ v }}
                  </el-tag>
                  <template #content>
                    {{ k }}:{{ v }}
                  </template>
                </el-tooltip> -->
              </div>
            </div>
            <div v-if="scope.row.labels == null">---</div>
            <div v-if="scope.row.labels != null && Object.keys(scope.row.labels).length > 3"><el-button
              size="small"
              type="text"
              link
              @click="showLabels(scope.$index)"
            >{{
              maxitem[scope.$index] == 3 ?
                '展开' : '收起'
            }}</el-button></div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="150">
          <template #default="scope">
            <div style="display: flex; align-items: center">
              <span v-if="scope.row.status == 'Ready' && scope.row.unschedulable == false" slot="reference">
                <el-tooltip placement="top" effect="light"><template #content>Ready</template>
                  <i class="dotClass" style="background-color: springgreen" /></el-tooltip>
              </span>
              <span v-if="scope.row.status == 'Ready' && scope.row.unschedulable == true" slot="reference">
                <el-tooltip placement="top" effect="light"><template #content>Unschedulable</template>
                  <i class="dotClass" style="background-color: orange" /></el-tooltip>
              </span>
              <span v-if="scope.row.status != 'Ready'" slot="reference">
                <el-tooltip placement="top" effect="light"><template #content> NotReady </template>
                  <i class="dotClass" style="background-color: red" /></el-tooltip>
              </span>
              <span>
                <span
                  v-if="scope.row.status === 'Ready' && scope.row.unschedulable === false"
                  style="margin-left: 3px;color: green;"
                >
                  运行中
                </span>
                <span
                  v-if="scope.row.status == 'Ready' && scope.row.unschedulable == true"
                  style="margin-left: 3px;color: orangered;"
                >
                  不可调度
                </span>
                <span v-if="scope.row.status != 'Ready'" style="margin-left: 3px;color: red;">
                  异常
                </span>
                <span v-if="scope.row.taints != null">
                  <el-tooltip class="box-item" effect="light" placement="top">
                    <template #content>
                      污点：<br>
                      <span v-for="(v, k) in scope.row.taints" :key="k">
                        {{ v.key }}:{{ v.effect }}<br>
                      </span>
                    </template>
                    <el-button size="small" type="info">
                      <template #default>
                        {{ scope.row.taints.length }}
                      </template>
                    </el-button>
                  </el-tooltip>
                </span>
              </span>

            </div>
          </template>
        </el-table-column>
        <el-table-column prop="cpu" label="CPU（m）" width="100">
          <template #default="scope">
            <div>
              <el-tooltip class="box-item" effect="light" placement="top">
                <template #content>
                  <span>{{ resource["cpu_request"][scope.row.name] }} m</span>
                </template>
                <span>使用：{{ toNumber(resource["cpu_request"][scope.row.name] * 100 / scope.row.cpu_total) }}% </span>
              </el-tooltip>
            </div>
            <div>
              <el-tooltip class="box-item" effect="light" placement="top">
                <template #content>
                  <span>{{ resource["cpu_limit"][scope.row.name] }} m</span>
                </template>
                <span>限制：{{ toNumber(resource["cpu_limit"][scope.row.name] * 100 / scope.row.cpu_total) }}% </span>
              </el-tooltip>
            </div>
            <div>
              <span>总量：{{ scope.row.cpu_total }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="memory" label="内存（Mi）" width="100">
          <template #default="scope">
            <div>
              <el-tooltip class="box-item" effect="light" placement="top">
                <template #content>
                  <span> {{ resource["memory_request"][scope.row.name] }} Mi</span>
                </template>
                <span>使用：{{ toNumber(resource["memory_request"][scope.row.name] * 100 / scope.row.memory_total) }}%
                </span>
              </el-tooltip>
            </div>
            <div>
              <el-tooltip class="box-item" effect="light" placement="top">
                <template #content>
                  <span> {{ resource["memory_limit"][scope.row.name] }} Mi</span>
                </template>
                <span>限制：{{ toNumber(resource["memory_limit"][scope.row.name] * 100 / scope.row.memory_total) }}%
                </span>
              </el-tooltip>
            </div>
            <div>
              <span>总量：{{ scope.row.memory_total }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="pods" align="center" label="pod 数量" width="80" />
        <el-table-column prop="create_time" label="创建时间" width="140" />
        <el-table-column prop="kubelet_version" label="版本" width="80" align="center" />
        <el-table-column label="操作" width="50" align="center">
          <template #default="scope">
            <el-dropdown trigger="click" @command="setNodeSchedule(scope.row)">
              <el-button type="text">
                <i class="el-icon-s-operation" :style="{ fontSize: '18px' }" />
              </el-button>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item command="a">
                  <i :class="scope.row.unschedulable == false ? 'el-icon-video-pause' : 'el-icon-video-play'" />
                  {{ scope.row.unschedulable == false ? '停止调度' : '启动调度' }}</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </div>

  </div>
</template>

<script>
import Vue from 'vue'
export default {
  data() {
    return {
      maxitem: [],
      nodeItem: [],
      resource: {}
    }
  },
  created() {
    this.getNodelist()
  },
  methods: {
    getNodelist() {
      this.$store.dispatch('node/getNode').then((res) => {
        this.nodeItem = res.data.Item
        this.resource = res.data.resource
        for (let i = 0; i < res.data.Item.length; i++) {
          this.maxitem.push(3)
        }
      })
    },
    handle(row) {
      this.$router.push({
        path: '/cluster/node_detail',
        name: 'node详情',
        query: {
          nodename: row.metadata.name
        }
      })
    },
    showLabels(index) {
      if (this.maxitem[index] === 3) {
        Vue.set(this.maxitem, index, 99)
      } else {
        Vue.set(this.maxitem, index, 3)
      }
    },
    test() {
      console.log('111')
    },
    setNodeSchedule(row) {
      const data = {
        node_name: row.name,
        status: !row.unschedulable
      }
      this.$store.dispatch('node/changeNodeSchedule', data).then((res) => {
        if (data.status !== false) {
          this.$message({
            message: `节点 ${row.name} 已不可调度`,
            type: 'warning'
          })
        } else {
          this.$message({
            message: `节点 ${row.name} 已可调度`,
            type: 'success'
          })
        }
        setTimeout(() => {
          this.getNodelist()
        }, 1000)
      })
    },
    toNumber(val) {
      return Math.floor(val)
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";

  .table-bg-purple {
  padding-right: 2px;
  padding-left: 2px;
  border-radius: 4px;
  background: #f0f2f5;
}

</style>
