import request from '@/utils/request'

export function getConfigmap(data) {
  return request({
    url: '/cm/list',
    method: 'get',
    params: data
  })
}

export function getConfigmapDetail(data) {
  return request({
    url: '/cm/detail',
    method: 'get',
    params: data
  })
}

export function createConfigmap(data) {
  return request({
    url: '/cm/create',
    method: 'post',
    data
  })
}

export function updateConfigmap(data) {
  return request({
    url: '/cm/update',
    method: 'put',
    data
  })
}

export function deleteConfigmap(data) {
  return request({
    url: '/cm/delete',
    method: 'delete',
    data
  })
}
