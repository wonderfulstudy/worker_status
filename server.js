// 引入express
const express = require('express')
//创建一个app实例对象
const app = express()
 
//配置后端路由
//req 就是 request，就是请求过来的数据。
//res 就是 response，就是回调信息，要响应回去的数据。
app.get('/',(req,res)=>{
    res.send({
        "latestBlockNumber": 829805,
        "latestBlockHash": "0x0e0d0168762d73738bbe9c8c8e7226ad9e1a6afd1d0bac5f1e90aa6c070378fd",
        "finalizedBlockNumber": 829802,
        "finalizedBlockHash": "0x482b8edac250f0066587dadf5065be70c5c9b6eafd2bc569e1458b964c908efa",
        "workerAddress": "5FUFoFNpfsgUEX1CWLYodmfSSMKfbpMb2xfA47UmE1ZWn4B6",
        "workerPublicKey": "0x96aade25a73f217d3af05389292d64bf399e7aa3e4e98075401240edf0f6c63e",
        "workerStatus": "Online",
        "attestedAt": 1735319028,
        "appVersion": "v0.0.1-dev",
        "implVersion": 1
    })
})
 
// listen()端口监听
app.listen(8080,(err)=>{
    if(!err){
        console.log('服务器启动成功了')
    }
})
 
 