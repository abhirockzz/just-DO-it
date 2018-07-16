## Thumbs up for !

- Simple and intuitive UI
- Well documented REST API
- Great documentation & tutorials
- Affordable pricing and usage tiers, trials
- and India data center !

## "nice to have"

- Automate SSH keys generation (in addition to allowing users to use existing keys)
- Provide JSON payload for resources (e.g. Droplets) via the UI to make it super easy to use copy-paste them to create resources programatically (I realize that it can be accessed via API doc)
- One-click app suggestions (might be a little biased!) - Redis, NATS, Prometheus, Kafka
-  Categorize [community tutorials](https://www.digitalocean.com/community/tutorials) and allow searching by DO products (e.g. Droplets etc.)

## General observations

### Droplets

- Droplet Access Control - by default, it is possible to freely access HTTP/TCP endpoints on the Droplet. I realize that it 'makes things easy' but would it make sense to have explicit (user provided) access/firewall rules (except port 22 for SSH maybe) ?

### Spaces

- Programmatic Space creation only works with (hardcoded) `us-east-1` region (regardless of the actual DO region). This confused me since I thought that the Go based sample in the [Spaces API doc](https://developers.digitalocean.com/documentation/spaces/#introduction) was a typo
- Spaces created programmatically don't show up in main Dashboard/Control Panel, but they do in the **Spaces** menu
- Providing an example of a complete S3 URL (e.g. `s3://<space_key>:<space_secret>@<region>.digitaloceanspaces.com/<region>/<space_name>?ssl=true`) would be helpful (maybe in the Spaces API doc?) 

### Monitoring

After activating Monitoring for an existing Droplet, all the existing graphs disappeared for a while. This was a bit of a surprise since I thought that the default graphs (bandwidth, CPU, disk I/O) would continue to show up and the user would be notified about the fact that additional graphs are being churned