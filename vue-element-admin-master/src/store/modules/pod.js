import {
  getPod,
  getPodDetail,
  updatePod,
  deletePod,
  getPodLog,
  getPodContainer
} from '@/api/pod'

const actions = {
  getPod({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPod({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  getPodDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPodDetail({ namespace: data.namespace, pod_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  updatePod({ commit }, data) {
    return new Promise((resolve, reject) => {
      updatePod({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deletePod({ commit }, data) {
    return new Promise((resolve, reject) => {
      deletePod({ namespace: data.namespace, pod_name: data.pod_name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  getPodLog({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPodLog({ namespace: data.namespace, pod_name: data.pod_name, container_name: data.container_name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  getPodContainer({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPodContainer({ namespace: data.namespace, pod_name: data.pod_name }).then(response => {
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
