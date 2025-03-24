# etcdctl --endpoints=http://127.0.0.1:2379,http://127.0.0.1:2389,http://127.0.0.1:2399 endpoint status -w table

# 启动etcd集群
.PHONY: etcd etcd-stop etcd-status
etcd:
	etcd --config-file ./etcd/etcd-2379.conf > /tmp/etcd1.log 2>&1 &
	etcd --config-file ./etcd/etcd-2389.conf > /tmp/etcd2.log 2>&1 &
	etcd --config-file ./etcd/etcd-2399.conf > /tmp/etcd3.log 2>&1 &
	@echo "etcd cluster is starting..."
	@echo "Check logs in /tmp/etcd*.log"
	@sleep 2

etcd-stop:
	@echo "Stopping etcd cluster..."
	-pkill -15 etcd
	@rm -f /tmp/etcd*.log
	@echo "etcd cluster stopped"

etcd-status:
	etcdctl --endpoints=http://127.0.0.1:2379,http://127.0.0.1:2389,http://127.0.0.1:2399 endpoint status -w table

