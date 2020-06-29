# jupyter-notebooks-operator

oc apply -f deploy/crds/cache.example.com_jupyternotebooks_crd.yaml
oc create -f deploy/role.yaml
oc create -f deploy/role_binding.yaml
oc create -f deploy/service_account.yaml
oc create -f deploy/operator.yaml
oc apply -f deploy/crds/cache.example.com_v1alpha1_jupyternotebooks_cr.yaml
