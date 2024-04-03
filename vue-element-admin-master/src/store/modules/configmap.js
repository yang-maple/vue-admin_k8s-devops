import {
  getConfigmap,
  getConfigmapDetail,
  createConfigmap,
  updateConfigmap,
  deleteConfigmap
} from '@/api/configmap'

const actions = {
  getConfigmap({ commit }, data) {
    return new Promise((resolve, reject) => {
      getConfigmap({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getConfigmapDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getConfigmapDetail({ namespace: data.namespace, configmap_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  createConfigmap({ commit }, dataform) {
    return new Promise((resolve, reject) => {
      createConfigmap({ data: dataform }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updateConfigmap({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateConfigmap({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteConfigmap({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteConfigmap({ namespace: data.namespace, configmap_name: data.configmap_name }).then(response => {
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
