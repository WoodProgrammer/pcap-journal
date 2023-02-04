# kdump-artifact

This is a project that converts all the network activity as a solid artifact and keeps it as a pcap files on your target object storage (compatible with S3)

## How basically it works
It basically track the fs events and uploads the pcap outputs to the target S3 bucket.