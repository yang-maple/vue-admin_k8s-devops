import {
  getIngressClass,
  getIngressClassDetail,
  updateIngressClass,
  deleteIngressClass
} from '@/api/ingressclass'

const actions = {
  getIngressClass({ commit }, data) {
    return new Promise((resolve, reject) => {
      getIngressClass({ filter_name: data.filter_name, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getIngressClassDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getIngressClassDetail({ Name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  updateIngressClass({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateIngressClass({ data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteIngressClass({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteIngressClass({ Name: data.ingressClass_name }).then(response => {
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
