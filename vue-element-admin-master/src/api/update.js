import request from '@/utils/request'

export function uploadYaml(data) {
  return request({
    url: '/upload/uploadYaml',
    method: 'post',
    data
  })
}
