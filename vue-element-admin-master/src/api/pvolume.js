import request from '@/utils/request'

export function getPVolume(data) {
  return request({
    url: '/pv/list',
    method: 'get',
    params: data
  })
}

export function getPVolumeDetail(data) {
  return request({
    url: '/pv/detail',
    method: 'get',
    params: data
  })
}

export function createPVolume(data) {
  return request({
    url: '/pv/create',
    method: 'post',
    data
  })
}

export function updatePVolume(data) {
  return request({
    url: '/pv/update',
    method: 'put',
    data
  })
}

export function deletePVolume(data) {
  return request({
    url: '/pv/delete',
    method: 'delete',
    data
  })
}
