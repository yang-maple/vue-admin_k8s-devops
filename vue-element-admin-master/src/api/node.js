import request from '@/utils/request'

export function getNode() {
  return request({
    url: '/node/list',
    method: 'get'
  })
}

export function getNodeDetail(data) {
  return request({
    url: '/node/detail',
    method: 'get',
    params: data
  })
}

export function changeNodeSchedule(data) {
  return request({
    url: '/node/schedule',
    method: 'post',
    data
  })
}

