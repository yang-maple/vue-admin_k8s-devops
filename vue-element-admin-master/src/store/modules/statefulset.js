import {
  getStatefulset,
  getStatefulsetDetail,
  updateStatefulset,
  deleteStatefulset,
  modifyStatefulset
} from '@/api/statefulset'

const actions = {
  getStatefulset({ commit }, data) {
    return new Promise((resolve, reject) => {
      getStatefulset({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getStatefulsetDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getStatefulsetDetail({ namespace: data.namespace, stateful_set_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  updateStatefulset({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateStatefulset({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteStatefulset({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteStatefulset({ namespace: data.namespace, stateful_set_name: data.stateful_set_name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  modifyStatefulset({ commit }, data) {
    return new Promise((resolve, reject) => {
      modifyStatefulset({ namespace: data.namespace, stateful_set_name: data.stateful_set_name, replicas: Number(data.replicas) }).then(response => {
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
