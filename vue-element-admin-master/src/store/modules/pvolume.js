import {
  getPVolume,
  getPVolumeDetail,
  createPVolume,
  updatePVolume,
  deletePVolume
} from '@/api/pvolume'

const actions = {
  getPVolume({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPVolume({ filter_name: data.filter_name, limit: data.limit, page: data.page }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  getPVolumeDetail({ commit }, data) {
    return new Promise((resolve, reject) => {
      getPVolumeDetail({ persistent_volume_name: data.name }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },

  createPVolume({ commit }, dataform) {
    return new Promise((resolve, reject) => {
      createPVolume({ data: dataform }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  updatePVolume({ commit }, data) {
    return new Promise((resolve, reject) => {
      updatePVolume({ data: data.data }).then(response => {
        resolve(response)
      }).catch(error => {
        reject(error)
      })
    })
  },
  deletePVolume({ commit }, data) {
    return new Promise((resolve, reject) => {
      deletePVolume({ persistent_volume_name: data.pvolume_name }).then(response => {
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
