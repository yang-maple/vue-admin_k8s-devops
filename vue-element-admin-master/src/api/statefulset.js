import request from '@/utils/request'

export function getStatefulset(data) {
  return request({
    url: '/stateful/list',
    method: 'get',
    params: data
  })
}

export function getStatefulsetDetail(data) {
  return request({
    url: '/stateful/detail',
    method: 'get',
    params: data
  })
}

export function updateStatefulset(data) {
  return request({
    url: '/stateful/update',
    method: 'put',
    data
  })
}

export function deleteStatefulset(data) {
  return request({
    url: '/stateful/delete',
    method: 'delete',
    data
  })
}

export function modifyStatefulset(data) {
  return request({
    url: '/stateful/modify',
    method: 'put',
    data
  })
}
