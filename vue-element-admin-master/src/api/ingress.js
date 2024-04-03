import request from '@/utils/request'

export function getIngress(data) {
  return request({
    url: '/ing/list',
    method: 'get',
    params: data
  })
}

export function getIngressDetail(data) {
  return request({
    url: '/ing/detail',
    method: 'get',
    params: data
  })
}

export function createIngress(data) {
  return request({
    url: '/ing/create',
    method: 'post',
    data
  })
}

export function updateIngress(data) {
  return request({
    url: '/ing/update',
    method: 'put',
    data
  })
}

export function deleteIngress(data) {
  return request({
    url: '/ing/delete',
    method: 'delete',
    data
  })
}
