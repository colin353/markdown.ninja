# Kubernetes configuration

If you want to run this application in a kubernetes cluster, you can use the
sample replication controller definitions in this folder.

I'm using Google Compute Engine, which makes setting this up pretty straightforward. If you're also on GCE, there is one critical part to making this work: you must set up your cluster to use the `container-vm` image, *not* `gci`, which is the default.

The reason is that `gci` doesn't contain drivers for glusterfs. You can't easily install them on `gci`, either. So make sure you pick that when setting up your nodes.

Also, if you're not on GCE, you'll need to figure out how to get traffic into the cluster. I'm using a service with type=LoadBalancer, which probably only works with some managed clusters.

Once you've created all of the services and controllers in this folder, you'll need to follow the steps in `gluster-setup.md` to get the cluster working.
