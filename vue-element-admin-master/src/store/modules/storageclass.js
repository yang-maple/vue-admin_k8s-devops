import {
  getStorageClass,
  getStorageClassDetail,
  updateStorageClass,
  deleteStorageClass
} from '@/api/storageclass'

const actions = {
  getStorageClass({ commit }, data) {
    return new Promise((resolve, reject) => {
      getStorageClass({ filter_name: data.filter_name, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getStorageClassDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getStorageClassDetail({ Name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  updateStorageClass({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateStorageClass({ data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteStorageClass({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteStorageClass({ Name: data.storageclass_name }).then(response => {
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
