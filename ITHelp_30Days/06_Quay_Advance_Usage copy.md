Quay Advance Usage
===================

Tag Expiration
-------------

Images can be set to expire from a Red Hat Quay repository at a chosen date and time using a feature called tag expiration. Here are a few things to know about about tag expiration:

- When a tag expires, the tag is deleted from the repository. If it is the last tag for a specific image, the image is set to be deleted.
- Expiration is set on a per-tag basis, not for a repository on the whole.
- When a tag expires or is deleted, it is not immediately removed from the registry. The value of Time Machine (in User settings) defines when the deleted tag is actually removed and garbage collected. By default, that value is 14 days. Up until that time, a tag can be repointed to an expired or deleted image.
- The Red Hat Quay superuser has no special privilege related to deleting expired images from user repositories. There is no central mechanism for the superuser to gather information and act on user repositories. It is up to the owners of each repository to manage expiration and ultimate deletion of their images.

Tag expiration is very useful when you just want to build and push some images for temporary testing. Tag expiration can be set in two ways:


-  Setting tag expiration from a Dockerfile
Adding a label like ``quay.expires-after=20h`` via the Dockerfile LABEL command will cause a tag to automatically expire after the time indicated. The time values could be something like 1h, 2d, 3w for hours, days, and weeks, respectively, from the time the image is built.

-  Setting tag expiration from the repository
On the Repository Tag page there is a UI column titled EXPIRES that indicates when a tag will expire. Users can set this by clicking on the time that it will expire or by clicking the Settings button (gear icon) on the right and choosing Change Expiration.



Mirror Repository
-----------------

If you have multiple distinct Quay servers on different regions, or you want to sychronize the latest official images from DockerHub or Red Hat Registry to your local private registry, the mirror repository feature is for you.


To mirror an external repository from an external container registry, do the following:


1. Create a robot account to pull images for the mirrored repository:

2. Select Create New Repository and give it a name.

3. Select the Settings button and change the repository state to MIRROR.

![](images/03_quay/m01.png)

4. Open the new repository and select the Mirroring button in the left column. Fill in the fields to identify the repository you are mirroring in your new repository:

- Registry URL: Location of the container registry you want to mirror from.
- Tags: This field is required. You may enter a comma-separated list of individual tags (1-1,1-2,latest) or tag patterns (1-*). At least one Tag must be explicitly entered (ie. not a tag pattern) or the tag "latest" must exist in the remote repository.
- Sync Interval: Defaults to syncing every 24 hours. You can change that based on hours or days.
- Robot User: Select the robot account you created earlier to do the mirroring.
- Username: The username for logging into the external registry holding the repository you are mirroring.
- Password: The password associated with the Username. Note that the password cannot include characters that require an escape character (\).
- Start Date: The date on which mirroring begins. The current date and time used by default.
- Verify TLS: Check this box if you want to verify the authenticity of the external registry. Uncheck this box if, for example, you set up Red Hat Quay for testing with a self-signed certificate or no certificate.
- HTTP Proxy: Identify the proxy server needed to access the remote site, if one is required.

![](images/03_quay/m02.png)

5. Check the Sychronize status from ``Usage Logs``

![](images/03_quay/m03.png)

Notification
------------

Quay supports adding notifications to a repository for various events that occur in the repository’s lifecycle. It is a pretty useful feature that can notify user or call CI/CD pipeline when some events happen.

Quay support following event notification:
- Repository Push: A successful push of one or more images was made to the repository.
- Vulnerability Detected: A new vulnerability was detected in the exist images.

When an event happen, you can use the following methods to notify users:

- E-mail: An e-mail will be sent to the specified address describing the event that occurred.
- Flowdock Notification: Posts a message to Flowdock.
- Hipchat Notification: Posts a message to HipChat.
- Slack Notification: Posts a message to Slack.
- Webhook POST: An HTTP POST call will be made to the specified URL with the event’s data. This method can be used to trigger CI/CD pipeline for some automation.



Oauth Application
----------------

Red Hat Quay offers programmatic access via an OAuth 2 compatible API. It is very useful when you want to do some automation to set up and manage Quay, such as set up ``Mirror-Repository`` or ``Notification`` automatically. Generation of an OAuth access token must normally be done via either an OAuth web approval flow, or via the Generate Token tab in the Application settings with Quay's UI.

1. Create a new Oauth Application

![](images/03_quay/09.png)

2. Create Access Token

![](images/03_quay/10.png)


3. Give permission for this access token

![](images/03_quay/12.png)

4. obtain the access token, it will show only once.

![](images/03_quay/13.png)
![](images/03_quay/14.png)

5. use the access token to access Quay API. You can refer https://docs.quay.io/api/swagger/ to check the Quay API usage. for example:

```
curl -s -X GET -k https://quay-eu-uat/api/v1/repository/${QUAY_REPO}/image/${IMAGE_ID}/security?vulnerabilities=true -H "Authorization: Bearer ${QUAY_ACCESS_TOKEN}" -H "Content-Type: application/json"
```

Conclusion
-----------

As a private images registries solution, Quay support many useful functionalities for daily usage. We will discuss more features for the integration of Quay and OpenShift later.