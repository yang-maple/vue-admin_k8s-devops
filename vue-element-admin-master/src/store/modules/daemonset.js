import {
  getDaemonset,
  getDaemonsetDetail,
  updateDaemonset,
  deleteDaemonset
} from '@/api/daemonset'

const actions = {
  getDaemonset({ commit }, data) {
    return new Promise((resolve, reject) => {
      getDaemonset({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getDaemonsetDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getDaemonsetDetail({ namespace: data.namespace, daemon_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  updateDaemonset({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateDaemonset({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteDaemonset({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteDaemonset({ namespace: data.namespace, stateful_set_name: data.stateful_set_name }).then(response => {
        resolve(response)
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
