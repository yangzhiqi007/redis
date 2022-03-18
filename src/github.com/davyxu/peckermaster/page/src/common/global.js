import Vue from 'vue';

// 使用 Event Bus
const bus = new Vue();

function toPercent(point){
  let str=Number(point*100).toFixed(2);
  str+="%";
  return str;
}


export default {
  bus,
  toPercent,
}
