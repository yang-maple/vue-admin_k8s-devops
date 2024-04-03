import {
  getSecret,
  getSecretDetail,
  createSecret,
  updateSecret,
  deleteSecret
} from '@/api/secret'

const actions = {
  getSecret({ commit }, data) {
    return new Promise((resolve, reject) => {
      getSecret({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getSecretDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getSecretDetail({ namespace: data.namespace, secret_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  createSecret({ commit }, dataform) {
    return new Promise((resolve, reject) => {
      createSecret({ data: dataform }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updateSecret({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateSecret({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteSecret({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteSecret({ namespace: data.namespace, secret_name: data.secret_name }).then(response => {
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
