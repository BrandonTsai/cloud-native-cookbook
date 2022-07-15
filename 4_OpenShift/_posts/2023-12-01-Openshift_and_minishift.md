---
title: "#1 Openshift Introduction"
author: Brandon Tsai
---

Openshift Introduction
---------------------

Red Hat OpenShift (also called OCP) is a container application platform based on the Kubernetes for enterprise application development and deployment. But why do we need another enterprise grade of kubernetes to run our containerize applications? I list 3 benefits to show why Openshift overweight Kubernetes.


1. OCP is a product for the organizations that need more reliable and stable solution while K8S is just a project for every one. OpenShift allows the users to install the product that offers paid support with a subscription.

2. OCP integrated some useful feature as they are required to maintain the enterprise-level platform, while k8s need to set up these add-on manually.

3. The security policies of OpenShift are stricter than the Kubernetes. It is easy for people to run simple apps on the Kubernetes, but the security policies of OpenShift restrict them to do so. A certain level of permissions is required to maintain the minimum security level, which can be provided by OpenShift. With the use of OpenShift, users do not have any choice but have to learn the policies to deploy more apps.



Install and Try OpenShift 4 on Mac via Red Hat CodeReady Containers
--------------------------------------------------------------------------


**Prerequisites:**

CodeReady Containers requires the following minimum system resources to run Red Hat OpenShift:

- 4 virtual CPUs (vCPUs)
- 8 GB of memory
- 35 GB of storage space

You will also require the native hypervisor for your host operating system.
CodeReady Containers currently supports libvirt for Linux, HyperKit for macOS, and Hyper-V for Windows.


1. Install hyperkit

```
$ brew install hyperkit
$ brew link --overwrite hyperkit
```

2. download CodeReady Containers

Download CodeReady Containers archive from the [Red Hat CodeReady Containers product page](the Red Hat CodeReady Containers product page)

3. extract the CodeReady Containers archive for your operating system and place the binary in your $PATH

```
$ cp crc /usr/local/bin/
```

4. setup cluster
Once CodeReady Containers has been installed, set up your host environment with the `crc setup`command. This command must be run before starting the OpenShift cluster.


5. Start your OpenShift 4.x cluster

After your host environment has been set up with the crc setup command, you can start the OpenShift cluster with the `crc start` command.

When prompted, supply your user pull secret for the cluster. Your user pull secret can be copied or downloaded from the the Red Hat CodeReady Containers product page under the Pull Secret section. A Red Hat account is required to access the user pull secret.

When start process finished, it will show the user credentials that you can use to login OpenShift

![](https://lh3.googleusercontent.com/Cj8dsFgSD20czZbS0mQkgEWw5kYM4Uhx0th_Ox_NyUwlo_YMTl040P-4U4ZEkGPVxxbjEMj4VP9cPK8dNGgq2AC357PW-pgm2-Up2IUHjTkNu4GdiMJUhXszB7n2RKQZl124w3Lu8nAHR5D3QaxP0DIVyzbAenQFMwbHLWUENdwArjOq4elkSOvmtJfVuCmwHMvdO0-lxPtxdgIbt-kBNSp3cPKGcUu6h8rO6u22_zxQj5DvOQ6-ZfwS6N4uTSWxzO7Xq585zTrljyMuCcOAWrbjQBp3GvAF5z-PrXxhVbCdtfY_LsqeLXdC3ebqdof5cI3wXaEQpdyaeqw2Ln2AurL0mMQ3EakagWRO3kDKbD_pun1Pii5TBVKl6qPvmCA8z3Vjk3NHM3jB0j8nH86VaYpNxccKxw0BcSUQc84z_Yjgd8lH1YppqYYM4jKAaJlpSuJRaxDvk5Agr5_6bsW9A-5bcx0ZvabuHxIqgKRwVBwMuLjwHUhEqJTtKKXalasgTN3-Kz14TBAou7LRt3987J575mq1L6CqJV0YHkXalaz1bOIX7-vRFSNHGwzgGXC0xmwqrj8qrJXTk9Q1XK4In5n2YrXySEXojAXC0WLDkEphEEAe8wXKb6Pu2XfL4gpAML0V8Bhd7wUc9Bwk8P4rJNri4yfyfnYSo1YTzhgVPNF-5dBkvEkpcJ_bU1YjkQ=w2400-h238-no?authuser=0)


5. Using your cluster

To access the cluster, first set up your environment by following `crc oc-env` instructions.

```
$ Brandons-MacBook-Pro:crc-macos-1.16.0-amd64 brandon$ crc oc-env
export PATH="/Users/brandon/.crc/bin/oc:$PATH"
# Run this command to configure your shell:
# eval $(crc oc-env)
```

login OpenShift as developer

```
$ oc login -u developer -p developer https://api.crc.testing:6443
Login successful.

You don't have any projects. You can try to create a new project, by running

    oc new-project <projectname>
```


Create your first project/namespace.

```
$ oc new-project myproject
```

> A `project` is a Kubernetes `namespace` with additional annotations, and is the central vehicle by which access to resources for regular users is managed. A project allows a community of users to organize and manage their content in isolation from other communities. Users must be given access to projects by administrators, or if allowed to create projects, automatically have access to their own projects.
> Unlike Kubernetes, A normal user can not create pods in the default project/namespace. You must create the first project before appling the new pod.


list all projects/namespaces

```
$ oc projects
You have one project on this server: "myproject".

Using project "myproject" on server "https://api.crc.testing:6443".
```

In Kebernetes, you need to install plugin `kubens` to switch between namespaces. However, this is build-in feature in OpenShift, you just need to run `oc project <my-project>`.


Switch to different project/namespace

```
$ oc project myproject
Already on project "myproject" on server "https://api.crc.testing:6443".
```


6. Access the Web Console

You can now run `crc console` and use these credentials to access the OpenShift web console.

6.1 login as "Developer"

choose "htpasswd_provider"

![](https://lh3.googleusercontent.com/QT2Ubxy_cJVaisklaAcnqavWNQ7e7BWIw4uizE5QMjwQkIV_3Ny82U5gE27hDNDy3UKJixkbTgNgBcmxgtLHbhJ2g1rtYgN1RysmJamORA8T5ZRhGtNwF_JW0dfz3bzZjjblpKrMkYjdcmvXt90qu_nRDhd5ItqWiqr8jWWoTO02sOWf2Wrm7t-rKHT2WdWJE448-CJ7hax8GkmwMcTB8iKPDkfRxVWzycXRmEHTI1veBudfMTWpnULYtmO-1h3Ofk41xX7Ojuu0J63IKue3bIOUNaTqRliZK3zHqBnS6u28misoNfGYk49ozlLTXJx4yk-7MDwBHnBLSBrON8txVW4QstW8WGK4AF-O-mh7sNdVeiWDXxHJNPCQGXQ4LIbiU2RfMWaIZB0LNYx9DqQ8Djfpdcn7dQQbJmLe0aQuFfcPV_7HEMAJdj6IEQpTJGb6v_nMXgExnei6nd6lxZm7_TWeVWv-BGimABa9nwWCf-6aBe372YCxOpXlrNNsExerJnb5MRInT-o3dm3Iy3_pF6UY0x8PPiLMcNLwhMV-i00z2OkDPKzqv1ngR9I-zua6D8LyjsswOgnjqPWF4_0DuEn3zBLc8hjQzTTZen-p3eNLvX5dvZas2UfNYVkHsN_VeGJS1RgpzwigrPbDN17ERgRg7kJiXsux9Y9Chb5fKQHMTFsToH1KIewMyhEbhQ=w2208-h998-no?authuser=0)

Input developer password

![](https://lh3.googleusercontent.com/vQkCt_eXPfKrQi2qKCUou8jAodjfkCGyhPGQCDZpM1kweuwwbiGpYQewcF_fO_YT6SbH1ZAFBHqGvnc61zSSOj2k1LbcIp9Q6FKAV3JIfBsg5aGsSonJSrFoKirjYusSYMC0CTrezpZqXKXjM9VJcad0H_HyWosS-0neuxYH0UtHsJ-j8RGW3YGW_XGQO4OUJoXYaXM8hjOLShUsCnfhL6KSdMFiXpfqHQ2yUzbTg6xHOC4hK6Upj8mKpltK6dAETay8jkwamF3cWUgNg4d1qdeIMT6WHkRS-XxmNWLc2_CMvPB2mEvrCugNrcK0LWPSPiTEGea-GvlBR0CgJZdUN23KewBIcME-5xFxC6CT_VGgIb2dH4SKf09bYNQ9qkuNO69Nd3HBMXes2ep5cyIy_Pk-uhArGb_TQHnob2DmdoC-RQv6C45KXDaxcf7Q8dczlkJ3a-Vb90sGU0kqveXf0EEt-t5mnaTNj37a-v63qWl95obTmLlATmSVxpg6bS8CmjUca-Y-dRXaPWy4bhDIbBa3pMLjzCQhSHhrAJkEcCYu_DWcdWjoWA6hxM3nOyqu94ek25VoWP88w1_WSd7n5c_1Mq2eR9QH_L3Pi5IzBFvSyBBo2ZqPMtMa35-ZjvJiCRZtgtHA-EyIbkX26DcOd-CXjCJJtdbzI2pDgYQRufh3Z7S5ji_FYpSvmFGOew=w2146-h1104-no?authuser=0)


Check the Web concole for developer

![](https://lh3.googleusercontent.com/somhA2CDnW9xAk_U_f5EtGRBoeAWZu8mimH6pMzQfCFzpBRAJLYErWEhd0BfFfYNgQMGOauik1Hd7F6QJPuSjbGtET1nzUd9OCCBQfKJnYcN22wS1n4fS_qWN20tFcoDugNZFuG6JaNqNcPN_gGyGAwlxb7x6sjbyqGBlTf4Ra5XmypJEsILodt4aV4kLjGffH1GdmFQkkc2ewGzyASxcHbqxXdqyw6S7pmOo5_2lOw0fEdfLoMKerW5l1JTRvWeOZFto7JCWjug_mnCRS-jYyO9R8w0X_mPmuhqx7JqcpA5OWzomK2bUod0vTYFXnmVTngRCcWo3Rsf1voRmwd6wixRtZN70JnQ5h8fFd87utQXaqqc89W3eiZciORHp_wGvH2DXKsTj8h3e2MWceDTqQP6PClDzH5m2oGoFw3TwwHT3-_aoNirFnE0srJdo94FG9OJBxSfjx7APWl5_oHiJw5c96ZqBJiBZkmMcC_27iN8V4YAT9rYi2ev2uhpvYWVvcBC5KzkHEzPvtju2MvpqF8VvEwqu7UXQ8wE-2hmv6g9VAsdsPbstpl1NQd6joQLivElFurC_tvEFSAjY-LzV-mNDqFynuZjoV7Xq6b2-qGo1JqKSilOVhtsCXXolZfOPaWEbsYtbE0Mz0J0w1kF5edSwKdcc6zl0A8OnZYc2pPVLoOrbei9cLr28u1ROw=w2858-h1406-no?authuser=0)



6.1 login as "kubeadmin"

choose "kube:admin"

![](https://lh3.googleusercontent.com/lI-VTKyem3h253ZFPvqh7kdQCTjO3lxPwT4FSoj1FgsakVWpfvmI8taBv-zA1GB76hJLWWH3x9mRw__ZrZT9j0uLGlIH0Rs8rzfmZR3XHfk87SSFZspERhYYIwXceo3W2rqcFmvmuOdW45tl5cyzH7ta2A61ELpVv5sUraGcY2G46g5Bq_dfw1ZNgsDBcpSExjqS_SPx08xKojlSzM-uMHUr_-XxA9BNUtX24InW8ugCEbYJKdumYpMoBC6arzZ9rpvtu_CLK9yLVPO-hBd74C5WfsT9z2ljYPVbZ3XMYkfHNs4Kd6YAt4BfH15WzGRRDBHNdnma-pbw2Aycov49CE01hbC-6vyaLaadjt9gcVs_vQ6Kc9W44bODi3bLz2yO_YM1E-EPMFAywD-pJhCSo1R8W7CZdkKy5S0TBpDYVIWEmdvT0Diufyvb66bH3Orz-B2hhC8N6oGLIgb7gGAoJH9h-llUspGTi7V0yfOEONNDxoVoW_szzTbaJ3JPkDPAzJ4bF-k-rR7esUbCx-FnL8UDqXUKr2sM_bYv1O30o_fAwVwB55WCVRDV-Nb8jS_zg48awaOGPKvwpPw4NTbsUlSvGBC0ZdJGbJbj3uXyXmzkhJ_hNyqMBnNVgLIe1JB5x1Qc-b9XKQsYgRBcZAiDzOqBxPaaimeN2Idx5j4VdvpgD732OTNmXOYQKjYXMQ=w2208-h998-no?authuser=0)

Input admin password

![](https://lh3.googleusercontent.com/REMQYQhvgkfH_TmDNmlueaOz0TKFz1VTVpwC_eAKJYZuod92R5GJbHFBcekII1gLClzZ9lgA7SAZkQZ5U6wFBhyIHLpsAMIs0wzdP98tb0iGTCz_iPEgYo58mhsSLtzixLI57hq8prCE3kk5rznWztBzHVL2K-q_5qP7Z6BFQae2rqCDMEgrQ_awd-S81EmmmHXqVYX5-QeizHKP6_4fOqJJOT7Kkc9Ce3sh8SWxItwJE2TE7KZEybjlZ86c9ilvnPs07hLvCEg5REi9jSfx2Sf8mN6t35xYseNzqfi6h_ShF3IJbGC4M8uUvDHNTkLIDjQL00fYscW32u4WroKl5bM9q_gH-kffqIi0TknRvMtZ2Je13nNuc5VfgH7BzEiQJ8jt22N7EbJHy3Cim9cyjQMI2vwgonr64BMWPy0OPfE-ilHZ00Dna1iWiY_Tt4Q9ZCS1di7rOhARvU3DG0WspldJHSpuJU_vpTuJyKehAkhAqec9u4FqHAc32GejA5IgECvvVf0XfWUVKLZPrIyCTYGXZbwGYdjQDBT1ly3vZRIHDrtaR71lizQuuEz4e_mll2EJWPpDoSGHZ0aoGcB7z_zYlnal7e4uOe2E_D9HogWOXGR9-WfJnWZeFj-lHzbIClECdOeNFPDIwwpmEnhHwrIQat8JBIysiFJmzJDFatyEwO1xyQ-X_FBU_yYqmA=w711-h367-no?authuser=0)


Check the Web concole for admin

![](https://lh3.googleusercontent.com/RArp4lUsRydUzHyHFID3H-UuOd057k0CdWbHxPK8DvWpQXIgFiZ69FY5m13vRCxI0xcIS6MDQ_JfGDNxbb59Y--KWfkLUC_SDNqunAdNQE95jg89XIR6hplsmCRAG6YZyewugF34Ysb9l3cmntbHxjuXcX4GJO2Y9916DTN8-LXLJ81CpuFJZSOqRLYMLiB9F0RZZHYbdkX4j3KQX4BNKoqJFJ5E8_9ZC85ssV1gaWPapMRbKpVcG3HGoVIeBc2Rd6IMqG1pX6vE_2x_Tykw_4y1NCQoFW3ucD_89GrcxC_YMm7DqonAULLPFE-BFL4RKh-zAf3FGMEd_FQKRJvSkisBjHgGdUbA-m0lfwNiJpu3Ig3TL19PX1oIryEn9zuAkF_PWX7OKQXUrtbNHmR1WjveG8y9B4GpmPfcd_ktUnPIEOkY71lrIULjWqA7pi9lht4SARdiuGgNAvPvUunqkZLI7WaYWLCNnX3vGJFJRjWctV1U9usi-2xhejfIEz_PvMWs53K9Hb-KtmZ6dUFC_pSNP3TRB7M_okPMpSairyvII5gIGEm-pDdDsWH60Fsbzj12b42NfVg2TBVr9Hs3L0iaTwbNBsMWRG_3u3G-tnVRc2HefPRpDfJAwe98WDdQxZtDaf-XEXMu-lXOiM1mCVU6ekNR6sSJVE2MQ_LnjHTKgMOjg_FhQGCeEUOnnQ=w2872-h1558-no?authuser=0)


7. Stop cluster

```
$ crc stop
```

8. Check cluster status

```
$ crc status
CRC VM:          Stopped
OpenShift:       Stopped
Disk Usage:      0B of 0B (Inside the CRC VM)
Cache Usage:     13.05GB
Cache Directory: /Users/brandon/.crc/cache
```


Conclusion
==========

Since Kubernetes forms the base of OpenShift, one can find a lot of common aspects between the two. In OpenShift vs Kubernetes. Above mentioned are some of the significant differences available. While Kubernetes remains a container platform, OpenShift comes into being and keeps a tab on the needs of different enterprises.