# Node Labels Copier
Небольшой контроллер, который копирует значение лейбла `role` в лейбл `node-role.kubernetes.io/<ROLE>` на нодах.

Начиная с 1.16 версии kubernetes, kubelet (опция `--node-labels`) запрещает создавать лейблы с доменами k8s.io и kubernetes.io в именах. То есть при инициализации ноды  это теперь нельзя сделать. В качестве решения можно вместо этого создавать лейбл `role` и данный контроллер скопирует это в `node-role.kubernetes.io/<ROLE>` тем самым решив проблему.

Плюс в том что это способ уневерсальный, как для managed кластеров так и для self-hosted. Достаточно добавить лейбл `role` на ноды и задеплоить контроллер:
```bash
kubectl apply -f deploy/deploy.yaml
```

В yandex managed kubernetes просто добавьте лейбл `role` при создании группы нод:
```hcl
resource "yandex_kubernetes_node_group" "workers" {
  cluster_id  = "${yandex_kubernetes_cluster.my_cluster.id}"
  name        = "name"
  description = "description"
  version     = "1.19"

  labels = {
    "role" = "worker"
  }
```

В self-hosted нодах добавьте --node-labels:
```bash
cat > /etc/systemd/system/kubelet.service.d/25-extra-labels.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--node-labels role=worker"
EOF
systemctl daemon-reload
systemctl restart kubelet
```
Либо используйте конфиг kind: JoinConfiguration: `nodeRegistration.kubeletExtraArgs.role: worker`. Если используете kubeadm для инициализации ноды.



