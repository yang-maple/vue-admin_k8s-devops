import {
  getPVClaim,
  getPVClaimDetail,
  createPVClaim,
  updatePVClaim,
  deletePVClaim
} from '@/api/pvclaim'

const actions = {
  getPVClaim({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPVClaim({ filter_name: data.filter_name, namespace: data.namespace, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getPVClaimDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPVClaimDetail({ namespace: data.namespace, persistent_volume_claim_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  createPVClaim({ commit }, dataform) {
    return new Promise((resolve, reject) => {
      createPVClaim({ data: dataform }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updatePVClaim({ commit }, data) {
    return new Promise((resolve, reject) => {
      updatePVClaim({ namespace: data.namespace, data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deletePVClaim({ commit }, data) {
    return new Promise((resolve, reject) => {
      deletePVClaim({ namespace: data.namespace, persistent_volume_claim_name: data.persistent_volume_claim_name }).then(response => {
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
