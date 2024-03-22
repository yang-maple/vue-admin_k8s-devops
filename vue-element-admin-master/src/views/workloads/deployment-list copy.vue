<template>
  <div>
    <!-- <el-row style="padding-bottom: 10px;">
      <el-col :span="24">
        <el-card shadow="always" style="width: 100%;">
          <span>
            <div>
              <svg class="header-icon" aria-hidden="true">
                <use xlink:href="#icon-deployment-copy" />
              </svg>
              <span
                style="font-size: 24px; color: #242e42;text-shadow: 0 4px 8px rgba(36,46,66,.1);font-weight: 600;"
              >无状态副本集</span>
              <br>
              <span style="font-size: 12px;color: #79879c!important">
                无状态副本集（Deployment）是指在 Kubernetes 集群中运行的不保存任何状态或数据软件应用程序，它们将所有需要保存的状态存储在外部系统中，并且在重启或迁移时可以很快恢复。
              </span>
            </div>
          </span>
        </el-card>
      </el-col>
    </el-row> -->
    <div>
      <div class="container">
        <editor
          ref="aceEditor"
          v-model="content"
          width="700"
          height="600"
          lang="json"
          :theme="theme"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 6,
            fontSize: 14,
            showPrintMargin: false, //去除编辑器里的竖线
          }"
          @init="editorInit"
        />
        <el-button type="primary" size="small" @click="getValue">获 取</el-button>
        <el-button type="primary" size="small" @click="pre">上一个主题</el-button>
        <el-button type="primary" size="small" @click="next">下一个主题</el-button>
      </div>
    </div>
  </div></template>

<script scoped>
export default {
  components: {
    editor: require('vue2-ace-editor')
  },
  data() {
    return {
      content: '',
      theme: '',
      num: 0,
      arr: [ // 将brace/theme文件夹下的所有主题名字拷贝出来
        'ambiance',
        'chaos',
        'chrome',
        'clouds',
        'clouds_midnight',
        'cobalt',
        'crimson_editor',
        'dawn',
        'dracula',
        'dreamweaver',
        'eclipse',
        'github',
        'gob',
        'gruvbox',
        'idle_fingers',
        'iplastic',
        'katzenmilch',
        'kr_theme',
        'kuroir',
        'merbivore',
        'merbivore_soft',
        'monokai',
        'mono_industrial',
        'pastel_on_dark',
        'solarized_dark',
        'solarized_light',
        'sqlserver',
        'terminal',
        'textmate',
        'tomorrow',
        'tomorrow_night',
        'tomorrow_night_blue',
        'tomorrow_night_bright',
        'tomorrow_night_eighties',
        'twilight',
        'vibrant_ink',
        'xcode'
      ]
    }
  },
  mounted() {
    this.editorInit()
    this.theme = this.arr[0]
    console.log(this.$refs.aceEditor.editor.setValue('设置的初始值'))
  },
  methods: {
    editorInit() { // 初始化
      require('brace/ext/language_tools') // language extension prerequsite...
      require('brace/mode/javascript') // language
      require('brace/mode/json')
      // require("brace/theme/tomorrow_night");
      for (let i = 0; i < this.arr.length; i++) { // 方便看哪个主题好看，循坏加载了所有主题，通过点击按钮切换
        require('brace/theme/' + this.arr[i])
      }
      require('brace/snippets/javascript') // snippet
    },
    getValue() { // 获取编辑器中的值
      console.log('编辑器中的值：' + this.$refs.aceEditor.editor.getValue())
      console.log('编辑器中第一个换行符的位置：' + this.$refs.aceEditor.editor.getValue().indexOf('\n'))
    },
    pre() { // 切换到上一个主题
      if (this.num === 0) {
        return
      }
      this.num--
      this.theme = this.arr[this.num]
      console.log('主题' + this.num + '__' + this.arr[this.num])
    },
    next() { // 切换到下一个主题
      if (this.num === this.arr.length - 1) {
        return
      }
      this.num++
      this.theme = this.arr[this.num]
      console.log('主题' + this.num + '__' + this.arr[this.num])
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
