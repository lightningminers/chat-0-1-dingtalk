import axios from 'axios';
import Cookies from 'js-cookie';

const TYPE = {
  HTTP: 'http'
}

let MOCK = false;

function csrfSafeMethod(method) {
  return (/^(GET|HEAD|OPTIONS|TRACE|POST)$/.test(method));
}

export default function request(opt){
  return new Promise(function(resolve, reject){
    const { ServiceName, MethodName, body, type='http', meta } = opt;
    if (MOCK){
      const url = `http://127.0.0.1:8999${ServiceName || ''}${MethodName || ''}`;
      const xhr = new XMLHttpRequest();
      xhr.onreadystatechange = function(){
        if (xhr.readyState === 4){
          if (xhr.status === 200){
            const res = JSON.parse(xhr.responseText);
            resolve({ data: res});
          }
        }
      }
      const { method='get', host='' } = meta;
      xhr.open(method, url, true);
      xhr.setRequestHeader("Content-Type","application/x-www-form-urlencoded;");
      xhr.send(method === 'get' ? null : JSON.stringify(body));
    } else {
      switch(type){
        case TYPE.HTTP:
          console.info(['Call HTTP ServiceNam: ' + ServiceName + ' MethodName: ' + MethodName + ' body: ' + JSON.stringify(body)]);
          const { method='GET', host='' } = meta;
          const RequestConfig = {
            method: method,
            url: `${host}${ServiceName || ''}${MethodName || ''}`,
          }
          if (method === 'GET'){
            RequestConfig.params = body;
          } else {
            RequestConfig.data = body;
          }
          if (csrfSafeMethod(method)){
            const csrfToken = Cookies.get('csrfToken');
            RequestConfig.headers = {
              'x-csrf-token': csrfToken
            }
          }
          axios(RequestConfig).then(res => {
            resolve(res);

          }).catch(err => {
            reject(err);
          });
          break;
        default:
          console.info(['No Request'], LogType.WARNING);
          break;
      }
    }
  });
}