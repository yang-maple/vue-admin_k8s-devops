<template>
  <div class="container">
    <editor
      :value="content"
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
      width="100%"
      height="400px"
      @input="handleInput"
      @init="editorInit"
    />
    <!-- <el-button type="primary" size="small" @click="getValue">获 取</el-button> -->
    <el-button type="primary" size="small" @click="pre">上一个主题</el-button>
    <el-button type="primary" size="small" @click="next">下一个主题</el-button>
  </div>
</template>

<script>
export default {
  components: {
    editor: require('vue2-ace-editor')
  },
  model: {
    prop: 'content',
    event: 'change'
  },
  props: {
    content: {
      type: String,
      required: true,
      default: '' // 设置默认为空字符串或其他适合的默认值
    }
  },
  data() {
    return {
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
    // console.log(this.$refs.aceEditor.editor.setValue('设置的初始值'))
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
    // getValue() { // 获取编辑器中的值
    //   console.log('编辑器中的值：' + this.$refs.aceEditor.editor.getValue())
    //   console.log('编辑器中第一个换行符的位置：' + this.$refs.aceEditor.editor.getValue().indexOf('\n'))
    // },
    handleInput(e) {
      console.log('编辑器中的值：' + e)
      this.$emit('change', e) // 这里e是每次子组件修改的值，change就是上面event中声明的，不要变
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

<style scoped>
.settingTitle {
  font-size: larger;
}

.settingDescription {
  font-size: small;
  color: #a8a8af
}
</style>
