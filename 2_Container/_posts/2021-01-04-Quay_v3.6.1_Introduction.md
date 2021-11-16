

Introduction
------------

Despite the fact that we can build Image from Dockerfile before running a container, it is better to keep these Images on "Image Registry" and pull these images before running the container.

There are many Image Registry solutions, the most famous is Dockerhub. You can find numerous public registries and images on Dockerhub. However, if you want to build a local Image Registry on your private development environment or in your datacenter, I would suggest you consider the Red Hat solution - Quay.

Red Hat Quay container and application registry provides secure storage, distribution, and deployment of containers on any infrastructure. It is available as a standalone component or in conjunction with OpenShift.

Compare Harbor and Quay
----------------------

There is another similar and popular project - [Harbor](https://github.com/goharbor/harbor). It is an open source trusted cloud native registry project and hosted by the Cloud Native Computing Foundation (CNCF). Following is the comparison between Harbor and Quay.


| Product | Quay | Harbor |
|---------|------|--------|
| Language                      | Python | Golang |
| Type                          | Public (quay.io), private | private |
| Authentication                | LDAP, OIDC, Google, Github  | LDAP, OIDC, DB |
| Robot Account                 | O | O |
| Permission Management         | O | O |
| Security Scan                 | O | O |
| Image Signing                 | 0 | 0 |
| Image Clearning               | O | O |
| Helm Application Management   | O | O |
| Image Mirroring               | O | O |
| Notification                  | Webhook, Email, Slack | Webhook |

You can see there are not much different between Quay and Harbor.
If you are planing to run containers on OpenShift, I would recommend you using Quay, because it is supported by Red Hat and highly integrated with OpenShift.


Install Quay on Local Machine
-----------------------------

To install open source version, please refer https://github.com/quay/quay/blob/master/docs/getting_started.md. This article will focus on set up Red Hat Quay version 3.2.1 on a RHEL7 server for testing Environment.


### Prerequisites

1. docker service

```
# yum install docker
# systemctl enable docker
# systemctl start docker
```

2. Get username and password for pulling Red Hat Quay version 3 images.

Set up authentication to Quay.io, so you can pull the quay and Clair images, as described in [Accessing Red Hat Quay](https://access.redhat.com/solutions/3533201).


3. Pull required images

| Service | Image |
|---------|-------|
| PostgreSQL | registry.redhat.io/hscl/postgresql-96-rhel7:1 |
| Redis      | registry.redhat.io/rhscl/redis-5-rhel7:5 |
| Quay       | quay.io/redhat/quay:v3.2.1 |
| Clair      | quay.io/redhat/clair-jwt:v3.2.1 |


### Set up Postgresql Container

```
$ mkdir -p /opt/pgsql
$ chmod 777 /opt/pgsql
$ export POSTGRESQL_USER=quayuser
$ export POSTGRESQL_PASSWORD=JzxCTamgFBmHRhcGFtoPHFkrx1BH2vwQ
$ export POSTGRESQL_DATABASE=quaytestdb

$ docker run -d \
    --restart=always \
    --env POSTGRESQL_USER=${POSTGRESQL_USER} \
    --env POSTGRESQL_PASSWORD=${POSTGRESQL_PASSWORD} \
    --env POSTGRESQL_DATABASE=${POSTGRESQL_DATABASE} \
    --name pgsql \
    --privileged=true \
    --publish 5432:5432 \
    -v /opt/pgsql:/var/lib/psql/data:Z \
    registry.redhat.io/hscl/postgresql-96-rhel7:1
```


### Set up Redis 

```
# mkdir -p /opt/redis
# chmod 777 /opt/redis
# docker run -d --restart=always -p 6379:6379 \
    --privileged=true \
    -v /opt/redis:/var/lib/redis/data:Z \
    registry.redhat.io/rhscl/redis-5-rhel7:5
```

### Set up Quay

1. Run Quay Configuration mode

2. Download & Update Quay Configuration

3. Run container for Quay


### Set up Clair

1. add database for Clair in pgsql container

```

# create clair db user
$ docker exec -i pgsql /bin/bash -c 'createuser  {{ clair_pgsql_user }}'

# Giving the user a password
$ docker exec -i pgsql /bin/bash -c "echo \"ALTER USER {{ clair_pgsql_user }} WITH PASSWORD '{{ clair_pgsql_password }}';\" | psql"

# Creating Database
$ docker exec -i pgsql /bin/bash -c 'createdb -O {{ clair_pgsql_user }} {{ clair_pgsql_db }} '

# Granting privileges on database
$ docker exec -i pgsql /bin/bash -c 'echo "grant all privileges on database {{ clair_pgsql_db }} to {{ clair_pgsql_user }};" | psql'

```

2. add clair config file

``
clair:
  database:
    type: pgsql
    options:
      # A PostgreSQL Connection string pointing to the Clair Postgres database.
      # Documentation on the format can be found at: http://www.postgresql.org/docs/9.4/static/libpq-connect.html
      source: postgresql://{{clair_pgsql_user}}:{{clair_pgsql_password}}@{{pgsql_host}}:5432/{{clair_pgsql_db}}?statement_timeout=60000&sslmode=disable
      cachesize: 16384
  api:
    # The port at which Clair will report its health status. For example, if Clair is running at
    # https://clair.mycompany.com, the health will be reported at
    # http://clair.mycompany.com:6061/health.
    healthport: 6061

    port: 6062
    timeout: 900s

    # paginationkey can be any random set of characters. *Must be the same across all Clair instances*.
    paginationkey:

  updater:
    # interval defines how often Clair will check for updates from its upstream vulnerability databases.
    interval: 6h
    notifier:
      attempts: 3
      renotifyinterval: 1h
      http:
        # QUAY_ENDPOINT defines the endpoint at which Quay is running.
        # For example: https://myregistry.mycompany.com
        endpoint: https://quay/secscan/notify
        proxy: http://localhost:6063

jwtproxy:
  signer_proxy:
    enabled: true
    listen_addr: :6063
    ca_key_file: /certificates/mitm.key # Generated internally, do not change.
    ca_crt_file: /certificates/mitm.crt # Generated internally, do not change.
    signer:
      issuer: security_scanner
      expiration_time: 5m
      max_skew: 1m
      nonce_length: 32
      private_key:
        type: preshared
        options:
          # The ID of the service key generated for Clair. The ID is returned when setting up
          # the key in [Quay Setup](security-scanning.md)
          # replace key_id from the one you generated from quay UI
          key_id: 7488413ebdec76137deccb91a390ca157584caa7d3256a149d045db107c5c4d4
          private_key_path: /config/security_scanner.pem

  verifier_proxies:
  - enabled: true
    # The port at which Clair will listen.
    listen_addr: :6060

    # If Clair is to be served via TLS, uncomment these lines. See the "Running Clair under TLS"
    # section below for more information.
    # key_file: /clair/config/clair.key
    # crt_file: /clair/config/clair.crt

    verifier:
      # CLAIR_ENDPOINT is the endpoint at which this Clair will be accessible. Note that the port
      # specified here must match the listen_addr port a few lines above this.
      # Example: https://myclair.mycompany.com:6060
      audience: http://clair:6060

      upstream: http://localhost:6062
      key_server:
        type: keyregistry
        options:
          # QUAY_ENDPOINT defines the endpoint at which Quay is running.
          # Example: https://myregistry.mycompany.com
          registry: https://quay/keys/`

```

3. set up Clair container



Quay Basic Usage
-------------------

1. Create an organization and a new repository


2. Push Images to Quay


3. Check Scan Result


4. Pull Images from Quay


Robot Account Introduce
---------------------------


Application API
----------------


Expired Image
-------------

Notification
------------

Mirror Repository
-----------------




Conclusion
-----------

