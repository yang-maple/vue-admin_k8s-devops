<template>
  <div class="container">
    <el-card>
      <editor
        :value="content"
        :lang="language"
        :theme="theme"
        :options="{
          enableBasicAutocompletion: true,
          enableSnippets: true,
          enableLiveAutocompletion: true,
          // showInvisibles: true, // 显示不可见字符（如换行符和制表符）
          wrap: true, //启动软换行
          tabSize: 6,
          fontSize: 14,
          showPrintMargin: false, //去除编辑器里的竖线
        }"
        width="100%"
        height="500px"
        @input="handleInput"
        @init="editorInit"
      />
    </el-card>
    <!-- <el-button type="primary" size="small" @click="pre">上一个主题</el-button>
    <el-button type="primary" size="small" @click="next">下一个主题</el-button> -->
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
    },
    language: {
      type: String,
      reqiired: true,
      default: 'json'
    }
  },
  data() {
    return {
      saveconnect: '',
      theme: '',
      num: 0,
      arr: [ // 将brace/theme文件夹下的所有主题名字拷贝出来
        'chrome',
        'clouds',
        'clouds_midnight',
        'eclipse',
        'github',
        'tomorrow',
        'tomorrow_night',
        'tomorrow_night_blue',
        'tomorrow_night_bright',
        'tomorrow_night_eighties',
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
      // require('brace/mode/javascript') // language
      require('brace/mode/json')
      require('brace/mode/yaml')
      // require("brace/theme/tomorrow_night");
      for (let i = 0; i < this.arr.length; i++) { // 方便看哪个主题好看，循坏加载了所有主题，通过点击按钮切换
        require('brace/theme/' + this.arr[i])
      }
      // require('brace/snippets/javascript') // snippet
      require('brace/snippets/yaml')
      require('brace/snippets/json')
    },
    // getValue() { // 获取编辑器中的值
    //   console.log('编辑器中的值：' + this.$refs.aceEditor.editor.getValue())
    //   console.log('编辑器中第一个换行符的位置：' + this.$refs.aceEditor.editor.getValue().indexOf('\n'))
    // },
    handleInput(e) {
      this.saveconnect = e
      // this.$emit('change', e) // 这里e是每次子组件修改的值，change就是上面event中声明的，不要变
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

.container{
  .el-card__body{
    padding: 0
  }
}
</style>
