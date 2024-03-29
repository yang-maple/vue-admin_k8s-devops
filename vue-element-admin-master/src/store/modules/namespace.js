import { getNamespacelist, getNamespaceDeatil, createNamespace, deleteNamespace, updateNamespace } from '@/api/namespace'

const state = {
  namespace: ''
}

const mutations = {
  SET_NAMESPACE: (state, data) => {
    state.namespace = data
  }
}

const actions = {
  getNamespacelist({ commit }, data) {
    return new Promise((resolve, reject) => {
      getNamespacelist({ filter_name: data.filter_name, limit: data.limit, page: data.page }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  getNamespaceDeatil({ commit }, data) {
    return new Promise((resolve, reject) => {
      getNamespaceDeatil({ namespace_name: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  createNamespace({ commit }, data) {
    return new Promise((resolve, reject) => {
      createNamespace({ namespace_name: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteNamespace({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteNamespace({ namespace_name: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updateNamespace({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateNamespace({ data: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  setNamespce({ commit }, data) {
    commit('SET_NAMESPACE', data)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
