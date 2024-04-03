import {
  getService,
  getServiceDetail,
  createService,
  updateService,
  deleteService
} from '@/api/service'

const actions = {
  getService({ commit }, data) {
    return new Promise((resolve, reject) => {
      getService({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getServiceDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getServiceDetail({ namespace: data.namespace, service_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  createService({ commit }, dataform) {
    return new Promise((resolve, reject) => {
      createService({ data: dataform }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updateService({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateService({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteService({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteService({ namespace: data.namespace, service_name: data.service_name }).then(response => {
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
