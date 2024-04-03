import {
  getIngress,
  getIngressDetail,
  createIngress,
  updateIngress,
  deleteIngress
} from '@/api/ingress'

const actions = {
  getIngress({ commit }, data) {
    return new Promise((resolve, reject) => {
      getIngress({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getIngressDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getIngressDetail({ namespace: data.namespace, ingress_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  createIngress({ commit }, dataform) {
    return new Promise((resolve, reject) => {
      createIngress({ data: dataform }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updateIngress({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateIngress({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteIngress({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteIngress({ namespace: data.namespace, ingress_name: data.ingress_name }).then(response => {
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
