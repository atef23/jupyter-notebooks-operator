
# Jupyter-Notebooks-Operator

This operator deploys machine learning workflows notebooks: https://github.com/willb/openshift-ml-workflows-workshop

**Requirements:**
- Openshift CLI

Run the following commands to deploy the operator and custom resource:

      
    oc apply -f deploy/crds/cache.example.com_jupyternotebooks_crd.yaml
    oc create -f deploy/role.yaml
    oc create -f deploy/role_binding.yaml
    oc create -f deploy/service_account.yaml
    oc create -f deploy/operator.yaml
    oc apply -f deploy/crds/cache.example.com_v1alpha1_jupyternotebooks_cr.yaml
