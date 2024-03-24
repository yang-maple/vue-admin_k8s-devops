import { uploadYaml } from '@/api/update'

const actions = {
  uploadYaml({ commit }, data) {
    return new Promise((resolve, reject) => {
      uploadYaml({ yamlContent: data }).then(res => {
        resolve(res)
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
