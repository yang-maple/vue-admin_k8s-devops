import request from '@/utils/request'

export function getDeployment(data) {
  return request({
    url: '/deployment/list',
    method: 'get',
    params: data
  })
}

export function getDeploymentDetail(data) {
  return request({
    url: '/deployment/detail',
    method: 'get',
    params: data
  })
}

export function createDeployment(data) {
  return request({
    url: '/deployment/create',
    method: 'post',
    data
  })
}

export function updateDeployment(data) {
  return request({
    url: '/deployment/update',
    method: 'put',
    data
  })
}

export function deleteDeployment(data) {
  return request({
    url: '/deployment/delete',
    method: 'delete',
    data
  })
}

export function modifyDeployment(data) {
  return request({
    url: '/deployment/modify',
    method: 'put',
    data
  })
}

export function restartDeployment(data) {
  return request({
    url: '/deployment/restart',
    method: 'post',
    data
  })
}
