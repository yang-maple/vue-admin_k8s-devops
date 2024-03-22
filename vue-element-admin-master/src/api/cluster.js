import request from '@/utils/request'

export function getCluster() {
  return request({
    url: '/cluster/list',
    method: 'get'
  })
}

export function changeCluster(data) {
  return request({
    url: '/cluster/change',
    method: 'post',
    data
  })
}

export function deleteCluster(data) {
  return request({
    url: '/cluster/delete',
    method: 'delete',
    data
  })
}

export function updateCluster(data) {
  return request({
    url: '/cluster/update',
    method: 'put',
    data
  })
}

export function detailCluster(data) {
  return request({
    url: '/cluster/detail',
    method: 'get',
    params: data
  })
}
