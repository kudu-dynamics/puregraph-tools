s3cli
=====
Provide a simple statically-compiled s3cli utility that can be post-inserted
into any Docker image to download and upload objects to/from an S3-compatible
API.

This is necessary as the semantics for most s3 tools perform additional checks
like getting object stat information or listing parent buckets that is too
costly to perform in our environment (specifically a MinIO deployment with
millions of files in 1 bucket).

### Tools that time out when trying to retrieve a file

Nomad specifically uses go-getter and has issues downloading files due to an
extra set of costly API calls.

- awscli (Python)
- go-getter (Go)
- mc (Go)

Distribution Statement "A" (Approved for Public Release, Distribution
Unlimited).
