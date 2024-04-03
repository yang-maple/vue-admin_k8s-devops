<template>
  <el-table-column :label="column.label" :width="column.width">
    <template #default="scope">
      <div v-for="(v, k, index) in scope.row[column.prop]" :key="k">
        <div v-if="index < maxitem[scope.$index]">
          <el-tooltip :content="k + ': ' + v" placement="top-start" effect="light" :disabled="!isTooltipNeeded(k, v)">
            <el-tag
              type="info"
              style=" max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"
              effect="plain"
              size="small"
              round
            >
              {{ k }}: {{ v }}
            </el-tag>
          </el-tooltip>
        </div>
      </div>
      <div v-if="scope.row[column.prop] == null">---</div>
      <div v-if="scope.row[column.prop] != null && Object.keys(scope.row[column.prop]).length >3"><el-button
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
</template>

<script>
import Vue from 'vue'
export default {
  props: {
    column: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      maxitem: []
    }
  },
  methods: {
    isTooltipNeeded(k, v) {
      const tagContent = `${k}: ${v}`
      const tagLength = tagContent.length * 7 // Assuming an average character width of 7px (adjust as needed)

      return tagLength > 300 // Adjust the constant value to match your desired maximum width
    },
    hasParameters(parameters) {
      return Array.isArray(parameters) && parameters.length > 0
    },
    showLabels(index) {
      if (this.maxitem[index] === 3) {
        Vue.set(this.maxitem, index, 99)
      } else {
        Vue.set(this.maxitem, index, 3)
      }
    }
  }
}
</script>
