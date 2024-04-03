import request from '@/utils/request'

export function getPVClaim(data) {
  return request({
    url: '/pvc/list',
    method: 'get',
    params: data
  })
}

export function getPVClaimDetail(data) {
  return request({
    url: '/pvc/detail',
    method: 'get',
    params: data
  })
}

export function createPVClaim(data) {
  return request({
    url: '/pvc/create',
    method: 'post',
    data
  })
}

export function updatePVClaim(data) {
  return request({
    url: '/pvc/update',
    method: 'put',
    data
  })
}

export function deletePVClaim(data) {
  return request({
    url: '/pvc/delete',
    method: 'delete',
    data
  })
}
