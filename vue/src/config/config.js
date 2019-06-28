const baseUrl = 'http://127.0.0.1:10001'
//const baseUrl = 'https://api.hyahm.coms:10001'

const g = {
// cookie 过期时间，单位分，与后端保持一致
  expirament: 120,
  downloadUrl: baseUrl + '/share/down', // 下载用到的地址
  username: 'admin@qq.com',
  password: 'admin',
  uploadUrl: baseUrl + '/uploadimg',
  headImgUrl: baseUrl + '/upload/headimg',
  shareUpload: baseUrl + '/share/upload', // 共享文件夹上传接口
  apiUrl: baseUrl,
  pubkey: `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDTMBzg/rduoIaGrxJtGxhvjJJJ
VLdQhiTu+LzJb/7mhvK2wExiPE5+rNyx5t9bGrsOdcSOY+KtxYiqqn3pxc5qeVHI
ptHH4XdMN4deNTAjbKRp6w54pvHWGejuak0MPV2g6ml+sYE2LtyFh8z932YqzWd6
FWo+L2+C8bOLanM97wIDAQAB
-----END PUBLIC KEY-----`
}

export default g