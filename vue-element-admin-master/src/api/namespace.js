import request from '@/utils/request'

export function getNamespace(data) {
  return request({
    url: '/namespaces/list',
    method: 'get',
    params: data
  })
}

export function getNamespaceDeatil(data) {
  return request({
    url: '/namespaces/detail',
    method: 'get',
    params: { data }
  })
}

export function deleteNamespace(data) {
  return request({
    url: '/namespaces/delete',
    method: 'delete',
    data
  })
}

export function updateNamespace(data) {
  return request({
    url: '/namespaces/update',
    method: 'put',
    data
  })
}

export function createNamespace(data) {
  return request({
    url: '/namespaces/create',
    method: 'post',
    data
  })
}
