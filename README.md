# r2

r2 is a command line interface for working with Cloudflare's
[R2 Storage](https://www.cloudflare.com/products/r2/).

Cloudflare's R2 implements the
[S3 API](https://docs.aws.amazon.com/AmazonS3/latest/API/Welcome.html),
attempting to allow users and their applications to migrate easily, but
importantly lacks the key, simple-to-use features provided by the AWS CLI's
[s3 subcommand](https://docs.aws.amazon.com/cli/latest/reference/s3/), as
opposed to the more complex and verbose API calls of the
[s3api subcommand](https://docs.aws.amazon.com/cli/latest/reference/s3api/index.html).
This CLI fills that gap.

## Progress

We're working to implement all the functionality of the AWS CLI's s3 subcommand.
Here's where we stand:

- [x] configure
- [ ] cp
- [x] ls
- [ ] mb
- [ ] mv
- [ ] presign
- [ ] rb
- [ ] rm
- [ ] sync
- [ ] website

If you'd like to contribute, to our project, we'd much appreciate your help!
