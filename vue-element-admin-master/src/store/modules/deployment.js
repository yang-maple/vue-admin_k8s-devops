import {
  getDeployment,
  getDeploymentDetail,
  createDeployment,
  updateDeployment,
  deleteDeployment,
  modifyDeployment,
  restartDeployment
} from '@/api/deployment'

const actions = {
  getDeployment({ commit }, data) {
    return new Promise((resolve, reject) => {
      getDeployment({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getDeploymentDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getDeploymentDetail({ namespace: data.namespace, deploy_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  createDeployment({ commit }, dataform) {
    return new Promise((resolve, reject) => {
      createDeployment({ data: dataform }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updateDeployment({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateDeployment({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteDeployment({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteDeployment({ namespace: data.namespace, deploy_name: data.deploy_name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  modifyDeployment({ commit }, data) {
    return new Promise((resolve, reject) => {
      modifyDeployment({ namespace: data.namespace, deploy_name: data.deploy_name, replicas: Number(data.replicas) }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  restartDeployment({ commit }, data) {
    return new Promise((resolve, reject) => {
      restartDeployment({ namespace: data.namespace, deploy_name: data.deploy_name }).then(response => {
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
