1. 执行./make.sh完成编译，将会在build目录下生成服务端程序nsd和客户端程序nscli。
2. 编译完成后执行./init.sh实现创建账户和添加token功能。
3. 测试
```bash
cd build
./nsd&
./nscli tx nameservice buy-name jack.id 5nametoken --from jack #买入name
./nscli tx nameservice set_can_auction jack.id true --from jack #设置可拍卖状态
./nscli tx nameservice offset jack.id 5nametoken --from alice #竞拍
```
