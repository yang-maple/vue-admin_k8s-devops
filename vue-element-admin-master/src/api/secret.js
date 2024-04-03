import request from '@/utils/request'

export function getSecret(data) {
  return request({
    url: '/secret/list',
    method: 'get',
    params: data
  })
}

export function getSecretDetail(data) {
  return request({
    url: '/secret/detail',
    method: 'get',
    params: data
  })
}

export function createSecret(data) {
  return request({
    url: '/secret/create',
    method: 'post',
    data
  })
}

export function updateSecret(data) {
  return request({
    url: '/secret/update',
    method: 'put',
    data
  })
}

export function deleteSecret(data) {
  return request({
    url: '/secret/delete',
    method: 'delete',
    data
  })
}
