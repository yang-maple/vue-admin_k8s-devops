import request from '@/utils/request'
export function sendmail(data) {
  return request({
    url: '/user/register/email',
    method: 'post',
    data
  })
}
