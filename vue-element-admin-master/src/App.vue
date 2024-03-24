<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}
</script>

<style>
/* ::-webkit-scrollbar{
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track{
  width: 6px;
  background: rgba(#101F1C,0.1);
  -webkit-border-radius: 2em;
  -moz-border-radius: 2em;
  border-radius: 2em;
}

::-webkit-scrollbar-thumb {
  background-color: rgba(144, 147,153,.5);
  background-clip: padding-box;
  min-height: 28px;
  -webkit-border-radius: 2em;
  -moz-border-radius: 2em;
  border-radius: 2em;
  transition: background-color .3s;
  cursor: pointer;
}
*/
::-webkit-scrollbar-thumb:hover{
  background-color: rgba(144, 147,153,.3);
}
 ::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 0;
}

/* 滚动条整体部分 */
::-webkit-scrollbar {
  -webkit-appearance: none;
  width: 6px;
  height: 6px;
}

/*滚动条里面的小方块，能向上向下移动（或往左往右移动，取决于是垂直滚动条还是水平滚动条） */
::-webkit-scrollbar-thumb {
  cursor: pointer;
  border-radius: 5px;
  background: rgba(0, 0, 0, 0.15);
  transition: color 0.2s ease;
}
</style>
