import request from '@/utils/request'

export function getService(data) {
  return request({
    url: '/svc/list',
    method: 'get',
    params: data
  })
}

export function getServiceDetail(data) {
  return request({
    url: '/svc/detail',
    method: 'get',
    params: data
  })
}

export function createService(data) {
  return request({
    url: '/svc/create',
    method: 'post',
    data
  })
}

export function updateService(data) {
  return request({
    url: '/svc/update',
    method: 'put',
    data
  })
}

export function deleteService(data) {
  return request({
    url: '/svc/delete',
    method: 'delete',
    data
  })
}
