import request from '@/utils/request'

export function getDaemonset(data) {
  return request({
    url: '/daemon/list',
    method: 'get',
    params: data
  })
}

export function getDaemonsetDetail(data) {
  return request({
    url: '/daemon/detail',
    method: 'get',
    params: data
  })
}

export function updateDaemonset(data) {
  return request({
    url: '/daemon/update',
    method: 'put',
    data
  })
}

export function deleteDaemonset(data) {
  return request({
    url: '/daemon/delete',
    method: 'delete',
    data
  })
}
