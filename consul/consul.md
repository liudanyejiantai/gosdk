# consul集群搭建
### consul 14,11,10做server集群, 10以客户端的方式加入集群
***

## 在14上启动14的consul，会自动创建目录consuldata
consul agent -server -bootstrap-expect 1 -data-dir consuldata -node=node14 -bind=192.168.177.14 -ui-dir dist -datacenter=dc1 -client=0.0.0.0


## 在11上开启
consul agent -server -bootstrap-expect 2 -data-dir consuldata -node=node11 -bind=192.168.177.11 -ui-dir dist -datacenter=dc1 -client=0.0.0.0

## 在10上开启
consul agent -server -bootstrap-expect 2 -data-dir consuldata -node=node10 -bind=192.168.177.10 -ui-dir dist -datacenter=dc1 -client=0.0.0.0

## 在14上做添加
consul join 192.168.177.11

consul join 192.168.177.10


## 12以客户端的方式加入集群组中
consul agent  -data-dir consuldata -node=node12 -bind=192.168.177.12 -datacenter=dc1 -client=0.0.0.0
consul join 192.168.177.10



