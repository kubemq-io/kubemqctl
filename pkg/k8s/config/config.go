package config

var EnvTemplate = `
        - env:
            - name: KUBEMQ_TOKEN
              value: {{.Token}}
            - name: CLUSTER_ROUTES
              value: '{{.Name}}:5228'
            - name: CLUSTER_PORT
              value: '5228'
            - name: CLUSTER_ENABLE
              value: 'true'
            - name: GRPC_PORT
              value: '50000'
            - name: REST_PORT
              value: '9090'
            - name: KUBEMQ_PORT
              value: '8080'
			{{.Config}}
`
