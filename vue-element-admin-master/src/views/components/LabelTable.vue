<template>
  <div style="padding: 0%;">
    <el-button type="primary" icon="el-icon-plus" style="margin-bottom: 10px;" @click="add">新增</el-button>
    <el-table
      v-if="showtable"
      :data="tableData"
      stripe
      style="width: 100%"
      :header-cell-style="{ background: '#e6e7e9',padding:0 }"
      size="small"
    >
      <el-table-column prop="key" label="名称" center width="190">
        <template slot-scope="{row}">
          <template>
            <el-input v-model="row.key" class="edit-input" size="small" />
          </template>
        </template>
      </el-table-column>
      <el-table-column prop="value" label="值" center width="190">
        <template slot-scope="{row}">
          <template>
            <el-input v-model="row.value" class="edit-input" size="small" />
          </template>
        </template>
      </el-table-column>
      <el-table-column align="center" label="Actions" width="120">
        <template slot-scope="{row,$index}">
          <el-button type="danger" size="small" icon="el-icon-delete" @click="deleteTable($index)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showtable: false,
      tableData: []
    }
  },
  methods: {
    add() {
      this.showtable = true
      this.tableData.push({
        key: '',
        value: ''
      })
    },
    deleteTable(row) {
      this.tableData.splice(row, 1)
      if (this.tableData.length === 0) {
        this.showtable = false
      }
    }
  }
}
</script>
