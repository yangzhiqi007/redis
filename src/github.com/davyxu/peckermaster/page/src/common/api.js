import axios from 'axios';
import global from '@/common/global';
import local from '@/common/local'

export const SystemRequest = ( method, url, data, callback ) => {

  let headers = {
    'Content-Type': 'application/json',
  }

  let user = sessionStorage.getItem('user');
  if (user) {
    let profile = JSON.parse(user);
    headers['Access-Token'] = profile.Token
  }

  axios.request(
    {
      method: method,
      url: local.Host + url,
      data: JSON.stringify(data),
      headers: headers
    }
  ).then(function (res) {

    if (res.data === null){
      return
    }

    callback(res)

  }).catch(function (res) { // js错误也会捕获

    if (typeof res.response !== "undefined" && res.response.status !== 200){
      global.bus.$message.error(res.response.statusText+ " " + url +" " + res.response.data);
    }else{
      console.log(res);
      global.bus.$message.error(res);
    }

  })
}
