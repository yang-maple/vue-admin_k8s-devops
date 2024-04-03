import request from '@/utils/request'

export function getStorageClass(data) {
  return request({
    url: '/storageClass/list',
    method: 'get',
    params: data
  })
}

export function getStorageClassDetail(data) {
  return request({
    url: '/storageClass/detail',
    method: 'get',
    params: data
  })
}

export function updateStorageClass(data) {
  return request({
    url: '/storageClass/update',
    method: 'put',
    data
  })
}

export function deleteStorageClass(data) {
  return request({
    url: '/storageClass/delete',
    method: 'delete',
    data
  })
}
