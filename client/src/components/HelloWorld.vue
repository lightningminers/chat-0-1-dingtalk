<template>
  <div class="container">
    <div class="user-container">
      <div class="avatar"><img :src="userInfo.avatar" /></div>
      <div class="name">{{ userInfo.name }}</div>
    </div>
    <div class="c-container">
      <div class="c-item">
        <div class="i-item"><input type="text" class="money" placeholder="输入税前薪资"/></div>
        <button class="b-item">计算</button>
      </div>
      <div class="table">
        <table>
          <tr>
            <th>种类</th>
            <th>个人缴纳</th>
            <th>单位缴纳</th>
          </tr>
          <tr>
            <td>老保</td>
            <td>0(8%)</td>
            <td>0(19%)</td>
          </tr>
          <tr>
            <td>医保</td>
            <td>0(2%)</td>
            <td>0(10%)</td>
          </tr>
          <tr>
            <td>失业保险</td>
            <td>0(0.2%)</td>
            <td>0(0.8%)</td>
          </tr>
          <tr>
            <td>公积金</td>
            <td>0(12%)</td>
            <td>0(12%)</td>
          </tr>
          <tr>
            <td>工伤</td>
            <td>0</td>
            <td>0(0.4%)</td>
          </tr>
          <tr>
            <td>生育</td>
            <td>0</td>
            <td>0(0.8%)</td>
          </tr>
          <tr>
            <td>小结</td>
            <td>0</td>
            <td>0</td>
          </tr>
          <tr>
            <td>应税收入</td>
            <td>0</td>
            <td>N</td>
          </tr>
          <tr>
            <td>税收</td>
            <td>0</td>
            <td>N</td>
          </tr>
          <tr>
            <td>实际收入</td>
            <td>0</td>
            <td>N</td>
          </tr>
          <tr>
            <td>单位支出</td>
            <td>0</td>
            <td>N</td>
          </tr>
        </table>
      </div>
    </div>
  </div>
</template>

<script>

import request from '../http/index'

export default {
  name: 'HelloWorld',
  data () {
    return {
      CorpID: '',
      code: '',
      userInfo: {
        avatar: 'http://images.w3crange.com/3321837.jpeg',
        name: 'icepy'
      }
    }
  },
  created(){
    this.initApp();
  },
  watch:{
    code(){
      this.getUserInfo();
    }
  },
  methods:{
    initApp(){
      const Request = {
        ServiceName: '',
        MethodName: 'get_config',
        meta: {
          method: 'GET'
        },
        body: {}
      }
      request(Request).then((res) => {
        const data = res.data;
        this.CorpID = data.CorpID;
        this.getCode()
      }).catch((err) => {
        alert(JSON.stringify(err));
      })
    },
    getUserInfo(){
      const Request = {
        ServiceName: '',
        MethodName: 'user',
        meta: {
          method: 'GET'
        },
        body: {
          code: this.code
        }
      }
      request(Request).then((res) => {
        this.userInfo = res.data;
      }).catch((err) => {
        alert(JSON.stringify(err))
      })
    },
    getCode(){
      dd.ready(() => {
        dd.runtime.permission.requestAuthCode({
          corpId: this.CorpID,
          onSuccess: (result) => {
            this.code = result.code;
          },
          onFail : (err) => {
            console.log('err',err)
          }
        })
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="less" scoped>
  .container{
    width: 100%;
    font-size: 14px;
    .user-container{
      width: 100%;
      height: 40px;
      padding: 5px 0 13px 0;
      display: flex;
      flex-direction: row;
      border-bottom: 1px solid #e4e4e4;
      .avatar{
        width: 40px;
        height: 40px;
        margin-right: 10px;
        img{
          width: 40px;
          height: 40px;
          border-radius: 20px;
        }
      }
      .name{
        display: flex;
        justify-content: center;
        align-items: center;
      }
    }
    .c-container{
      width: 100%;
      padding: 20px 0 0 0;
      .c-item{
        width: 100%;
        display: flex;
        flex-direction: row;
        .i-item{
          margin-right: 10px;
          .money{
            width: 150px;
            height: 20px;
            border: 1px solid #e4e4e4;
            outline: none;
            border-radius: 5px;
            font-size: 12px;
            padding-left: 5px;
            color:#666; 
            -webkit-appearance: none;
          }
        }
        .b-item{
          margin: 0px;
          padding: 0px;
          color: #fff;
          background: #42b983;
          border: none;
          width: 45px;
          height: 25px;
          outline: none;
          border-radius: 5px;
        }
        
      }
      .table{
        margin-top: 20px;
        table{
          border-collapse: collapse;
          margin: 0 auto;
          width: 100%;
          td,th{
            border: 1px solid #e4e4e4;
            height: 25px;
            padding-left: 5px;
          }
        }
      }
    }
  }
</style>

