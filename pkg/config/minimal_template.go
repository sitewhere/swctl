/**
 * Copyright © 2014-2021 The SiteWhere Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

const minimalTemplate string = `microservices:
- functionalarea: asset-management
  name: Asset Management
  description: Provides APIs for managing assets associated with device assignments
  icon: devices_other
  replicas: 1
  multitenant: true
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8006
    jmxport: 1106
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.asset
      level: info
  configuration: null
- functionalarea: command-delivery
  name: Command Delivery
  description: Manages delivery of commands in various formats based on invocation events
  icon: call_made
  replicas: 1
  multitenant: true
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8012
    jmxport: 1112
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.commands
      level: info
  configuration: null
- functionalarea: device-management
  name: Device Management
  description: Provides APIs for managing the device object model
  icon: developer_board
  replicas: 1
  multitenant: true
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8004
    jmxport: 1104
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.device
      level: info
  configuration: null
- functionalarea: event-management
  name: Event Management
  description: Provides APIs for persisting and accessing events generated by devices
  icon: dynamic_feed
  replicas: 1
  multitenant: true
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8005
    jmxport: 1105
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.event
      level: info
  configuration: null
- functionalarea: event-sources
  name: Event Sources
  description: Handles inbound device data from various sources, protocols, and formats
  icon: forward
  replicas: 1
  multitenant: true
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8008
    jmxport: 1108
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.sources
      level: info
  configuration: null
- functionalarea: inbound-processing
  name: Inbound Processing
  description: Common processing logic applied to enrich and direct inbound events
  icon: input
  replicas: 1
  multitenant: true
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8007
    jmxport: 1107
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.inbound
      level: info
  configuration: null
- functionalarea: instance-management
  name: Instance Management
  description: Handles APIs for managing global aspects of an instance
  icon: language
  replicas: 2
  multitenant: false
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 8080
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    - name: http-rest
      protocol: TCP
      appprotocol: null
      port: 8080
      targetport:
        type: 0
        intval: 8080
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8001
    jmxport: 1101
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.instance
      level: info
    - logger: com.sitewhere.web
      level: info
  configuration:
    raw:
    - 123
    - 34
    - 117
    - 115
    - 101
    - 114
    - 77
    - 97
    - 110
    - 97
    - 103
    - 101
    - 109
    - 101
    - 110
    - 116
    - 34
    - 58
    - 123
    - 34
    - 115
    - 121
    - 110
    - 99
    - 111
    - 112
    - 101
    - 72
    - 111
    - 115
    - 116
    - 34
    - 58
    - 34
    - 115
    - 105
    - 116
    - 101
    - 119
    - 104
    - 101
    - 114
    - 101
    - 45
    - 115
    - 121
    - 110
    - 99
    - 111
    - 112
    - 101
    - 46
    - 115
    - 105
    - 116
    - 101
    - 119
    - 104
    - 101
    - 114
    - 101
    - 45
    - 115
    - 121
    - 115
    - 116
    - 101
    - 109
    - 46
    - 99
    - 108
    - 117
    - 115
    - 116
    - 101
    - 114
    - 46
    - 108
    - 111
    - 99
    - 97
    - 108
    - 34
    - 44
    - 34
    - 115
    - 121
    - 110
    - 99
    - 111
    - 112
    - 101
    - 80
    - 111
    - 114
    - 116
    - 34
    - 58
    - 56
    - 48
    - 56
    - 48
    - 44
    - 34
    - 106
    - 119
    - 116
    - 69
    - 120
    - 112
    - 105
    - 114
    - 97
    - 116
    - 105
    - 111
    - 110
    - 73
    - 110
    - 77
    - 105
    - 110
    - 117
    - 116
    - 101
    - 115
    - 34
    - 58
    - 54
    - 48
    - 125
    - 125
    object: null
- functionalarea: outbound-connectors
  name: Outbound Connectors
  description: Allows event streams to be delivered to external systems for additional processing
  icon: label
  replicas: 1
  multitenant: true
  podspec:
    annotations: {}
    name: ""
    dockerspec:
      registry: {{ .Registry }}
      repository: {{ .Repository }}
      tag: "{{ .Tag }}"
    imagepullpolicy: IfNotPresent
    ports:
    - name: ""
      hostport: 0
      containerport: 9000
      protocol: TCP
      hostip: ""
    - name: ""
      hostport: 0
      containerport: 9090
      protocol: TCP
      hostip: ""
    env:
    - name: sitewhere.config.k8s.name
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.name
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.namespace
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: metadata.namespace
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.k8s.pod.ip
      value: ""
      valuefrom:
        fieldref:
          apiversion: v1
          fieldpath: status.podIP
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref: null
    - name: sitewhere.config.product.id
      value: {{ .InstanceName }}
      valuefrom: null
    - name: sitewhere.config.keycloak.service.name
      value: sitewhere-keycloak-http
      valuefrom: null
    - name: sitewhere.config.keycloak.api.port
      value: "80"
      valuefrom: null
    - name: sitewhere.config.keycloak.realm
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.realm
      value: master
      valuefrom: null
    - name: sitewhere.config.keycloak.master.username
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.master.password
      value: sitewhere
      valuefrom: null
    - name: sitewhere.config.keycloak.oidc.secret
      value: ""
      valuefrom:
        fieldref: null
        resourcefieldref: null
        configmapkeyref: null
        secretkeyref:
          localobjectreference:
            name: {{ .InstanceName }}
          key: client-secret
          optional: null
    resources: null
    livenessprobe: null
    readinessprobe: null
  serivcespec:
    ports:
    - name: grpc-api
      protocol: TCP
      appprotocol: null
      port: 9000
      targetport:
        type: 0
        intval: 9000
        strval: ""
      nodeport: 0
    - name: http-metrics
      protocol: TCP
      appprotocol: null
      port: 9090
      targetport:
        type: 0
        intval: 9090
        strval: ""
      nodeport: 0
    type: ClusterIP
  debug:
    enabled: false
    jdwpport: 8016
    jmxport: 1116
  logging:
    overrides:
    - logger: com.sitewhere
      level: info
    - logger: com.sitewhere.grpc.client
      level: info
    - logger: com.sitewhere.microservice.grpc
      level: info
    - logger: com.sitewhere.microservice.kafka
      level: info
    - logger: org.redisson
      level: info
    - logger: com.sitewhere.asset
      level: info
    - logger: com.sitewhere.connectors
      level: info
  configuration: null
`
