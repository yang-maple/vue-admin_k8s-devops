<template>
  <div class="codeEditBox" :style="{ height: height + 'px' }">
    <editor
      ref="aceEditor"
      v-model="options.value"
      :lang="editorOptions.language"
      :options="editorOptions"
      theme="tomorrow_night_blue"
      @init="editorInit"
      @input="codeChange"
      @setCompletions="setCompletions"
    />
  </div>
</template>

<script scoped>
import Editor from 'vue2-ace-editor'
import ace from 'brace'
export default {
  name: 'AceEditor',
  components: {
    Editor
  },
  props: {
    options: {
      type: Object,
      default() {
        return {
          language: 'json'
        }
      }
    },
    height: {
      type: Number,
      default: 500
    }
  },
  data() {
    return {
      defaultOpts: {
        value: '',
        language: 'sql',
        // 设置代码编辑器的样式
        enableBasicAutocompletion: true, // 启用基本自动完成
        enableSnippets: true, // 启用代码段
        enableLiveAutocompletion: true, // 启用实时自动完成
        tabSize: 2, // 标签大小
        fontSize: 14, // 设置字号
        showPrintMargin: false // 去除编辑器里的竖线
      },
      languageObj: {
        javascript: ['mode', 'snippets'],
        css: ['mode', 'snippets'],
        json: ['mode', 'snippets']
      }
    }
  },
  computed: {},
  watch: {},
  created() {
    this.editorOptions = Object.assign(this.defaultOpts, this.options)
  },
  mounted() {
  },
  methods: {
    codeChange(val) {
      this.$emit('change', val)
    },
    editorInit() {
      const that = this
      require('brace/ext/searchbox') // 添加搜索功能
      require('brace/theme/clouds') // 添加风格
      require('brace/ext/language_tools') // language extension prerequsite...
      require('brace/mode/' + this.editorOptions.language)
      require('brace/snippets/' + this.editorOptions.language)
      // 添加自定义提示
      const completer = {
        getCompletions: function (editors, session, pos, prefix, callback) {
          that.setCompleteData(editors, session, pos, prefix, callback)
        }
      }
      const lnTools = ace.acequire('ace/ext/language_tools')
      lnTools.addCompleter(completer)
    },
    getVal() {
      return this.$refs.aceEditor.editor.getValue()
    },
    setCompleteData(editor, session, pos, prefix, callback) {
      const data = [
        {
          meta: '表名', // 描述
          caption: 'sonic', // 展示出的名字（一般与value值相同）
          value: 'sonic', // 数据值
          score: 1 // 权重 数值越大 提示越靠前
        },
        {
          meta: '库名',
          caption: 'sonww',
          value: 'sonww',
          score: 2
        }
      ]
      if (prefix.length === 0) {
        return callback(null, [])
      } else {
        return callback(null, data)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
  @import "~@/styles/anticon.scss";
</style>
