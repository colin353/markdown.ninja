# Gluster setup

When you boot up the daemonset and the service, you'll get three symmetric glusterfs
nodes. But they don't automatically communicate with each other and form a cluster.
You'll have to manually ssh into each node and set it up that way.

## On the first machine

The first machine defines the volume and invites others to join the cluster.

```
gluster peer probe `getent hosts gluster.default | awk '{ print $1 }' | grep -ve \`hostname -i\` | head -n 1`
setfattr -x trusted.glusterfs.volume-id /mnt/brick1/media
setfattr -x trusted.gfid /mnt/brick1/media/
gluster volume create media `hostname -i`:/mnt/brick1/media
gluster volume start media
```

## Subsequent machines

You'll have to go into a cluster member and send an invite to the new machine:

```
gluster peer probe <NEW_MEMBER_IP>
```

and then once it joins the cluster, you need to jump over to the new machine and run

```
setfattr -x trusted.glusterfs.volume-id /mnt/brick1/media
setfattr -x trusted.gfid /mnt/brick1/media/
gluster volume add-brick media replica 2 `hostname -i`:/mnt/brick1/media
```

## setfattr

The setfattr commands undo some attributes set by gluster and allow you to remake
the cluster on a node. I think it might cause file corruption with your old stuff
on there, though. So probably don't do it unless you messed up during testing and
want to scrap it and restart.
