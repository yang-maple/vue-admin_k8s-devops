import request from '@/utils/request'

export function getPod(data) {
  return request({
    url: '/pod/list',
    method: 'get',
    params: data
  })
}

export function getPodDetail(data) {
  return request({
    url: '/pod/detail',
    method: 'get',
    params: data
  })
}

export function updatePod(data) {
  return request({
    url: '/pod/update',
    method: 'put',
    data
  })
}

export function deletePod(data) {
  return request({
    url: '/pod/delete',
    method: 'delete',
    data
  })
}

export function getPodLog(data) {
  return request({
    url: '/pod/container/log',
    method: 'get',
    params: data
  })
}

export function getPodContainer(data) {
  return request({
    url: '/pod/container/list',
    method: 'get',
    params: data
  })
}
