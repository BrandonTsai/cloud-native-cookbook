

OpenShift 系列 2： 懶人救星 - Source-To-Image (S2I)
=================================================

Using S2I Base Images in OpenShift
====================================

If your application is not complicate, instead of implement Dockerfile and work hard to figure out how to fix the permission isse, another way is using the Source-To-Image (S2I) image.


The main reasons one might be interested in using source builds are:

Speed - with S2I, the assemble process can perform a large number of complex operations without creating a new layer at each step, resulting in a fast process.
Patchability - S2I allows you to rebuild the application consistently if an underlying image needs a patch due to a security issue.
User efficiency - S2I prevents developers from performing arbitrary yum install type operations during their application build, which results in slow development iteration.
Ecosystem - S2I encourages a shared ecosystem of images where you can leverage best practices for your applications.
This article is about creating a simple S2I builder image.

S2I generates a new Docker image using source code and a  builder Docker image. The S2I project includes a number of commonly used Docker builder images, such as Python or Ruby, you can also extend S2I with your own custom scripts.




Reference
---------

- https://www.openshift.com/blog/create-s2i-builder-image