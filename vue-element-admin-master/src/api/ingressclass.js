import request from '@/utils/request'

export function getIngressClass(data) {
  return request({
    url: '/ingressClass/list',
    method: 'get',
    params: data
  })
}

export function getIngressClassDetail(data) {
  return request({
    url: '/ingressClass/detail',
    method: 'get',
    params: data
  })
}

export function updateIngressClass(data) {
  return request({
    url: '/ingressClass/update',
    method: 'put',
    data
  })
}

export function deleteIngressClass(data) {
  return request({
    url: '/ingressClass/delete',
    method: 'delete',
    data
  })
}
