import etcd3
import sys
import json



# Function to put data into etcd
def put_etcd_data(key, value):
    client.put(key, json.dumps(value))
    print(f"Data written to etcd: {key}")

if __name__ == "__main__":
    # Connect to etcd
    client = etcd3.client(host=sys.argv[0].split(":")[0], port=int(sys.argv[0].split(":")[1]))

    # Create nodes
    for i in range(1, 5):
        node_id = f"node-{i}"
        node_data = {
            "metadata": {
                "name": node_id,
                "labels": {
                    "kubernetes.io/role": "worker"
                }
            },
            "status": {
                "capacity": {
                    "cpu": "4",
                    "memory": "16Gi"
                },
                "conditions": [
                    {
                        "type": "Ready",
                        "status": "True"
                    }
                ]
            }
        }
        key = f"/registry/nodes/{node_id}"
        put_etcd_data(key, node_data)

    # Create roles
    role_data = {
        "metadata": {
            "namespace": "default",
            "name": "admin-role"
        },
        "rules": [
            {
                "apiGroups": [""],
                "resources": ["pods"],
                "verbs": ["create", "delete", "get", "list", "watch"]
            }
        ]
    }
    role_key = "/registry/roles/default/admin-role"
    put_etcd_data(role_key, role_data)

    # Create rolebindings
    rolebinding_data = {
        "metadata": {
            "namespace": "default",
            "name": "admin-binding"
        },
        "subjects": [
            {
                "kind": "User",
                "name": "admin-user"
            }
        ],
        "roleRef": {
            "kind": "Role",
            "name": "admin-role",
            "apiGroup": "rbac.authorization.k8s.io"
        }
    }
    rolebinding_key = "/registry/rolebindings/default/admin-binding"
    put_etcd_data(rolebinding_key, rolebinding_data)

    print("Data population complete.")
