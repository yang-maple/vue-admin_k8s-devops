import { getNode, getNodeDetail, changeNodeSchedule } from '@/api/node'

// const state = {}

// const mutations = {}

const actions = {
  getNode({ commit }) {
    return new Promise((resolve, reject) => {
      getNode().then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  getNodeDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getNodeDetail({ node_name: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  changeNodeSchedule({ commit }, data) {
    return new Promise((resolve, reject) => {
      changeNodeSchedule({ node_name: data.node_name, status: data.status }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  }
}

export default {
  namespaced: true,
  // state,
  // mutations,
  actions
}
