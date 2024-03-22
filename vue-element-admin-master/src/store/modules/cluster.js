import { getCluster, changeCluster, deleteCluster, updateCluster, detailCluster } from '@/api/cluster'

const state = {
  Cluster_Config: {}
}

const mutations = {
  setCluster_Config(state, payload) {
    state.Cluster_Config = payload
  }
}

const actions = {
  getCluster({ commit }) {
    return new Promise((resolve, reject) => {
      getCluster().then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  changeCluster({ commit }, data) {
    return new Promise((resolve, reject) => {
      changeCluster({ cluster_name: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deleteCluster({ commit }, data) {
    return new Promise((resolve, reject) => {
      deleteCluster({ cluster_name: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updateCluster({ commit }, data) {
    return new Promise((resolve, reject) => {
      updateCluster({ cluster_id: data.id, cluster_name: data.cluster_name, cluster_type: data.type }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  },
  detailCluster({ commit }, data) {
    console.log(data)
    return new Promise((resolve, reject) => {
      detailCluster({ cluster_name: data }).then(res => {
        resolve(res)
      }).catch(error => {
        reject(error)
      })
    })
  }

}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
