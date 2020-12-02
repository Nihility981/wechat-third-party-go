# wechat-third-party-go
微信开放平台第三方平台，全网发布

需要注意的：

  1> 这个并不能开箱即用，只把关键的代码贴出来；有些代码比如：`tutils.HttpRequestPostJson`函数，需要自己补充
  
  2> APPID TOKEN ENCODINGAESKEY APPSECRET 在微信开放平台|第三方平台|第三方应用|详情里边，没有的话，需要自己设置
  
  3> component_verify_ticket与component_access_token需要保存起来，后续的请求需要用到
  
  4> 授权事件接收URL：对应`/auth/event`的处理函数（AuthEvent）；消息与事件接收URL：对应`/auth/wechat/event/:appid`的处理函数（EventAppid）
  
  5> 以上两个接口，只有勾选了 权限管理|消息管理权限时，微信才会发起回调
  
  6> 关键代码是能通过检测，发布成功的！
  
  N> 想到再写
