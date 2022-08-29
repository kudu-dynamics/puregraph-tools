certcat
=======
HashiCorp Vault's PKI integration with Nomad only allows for the issuing of a
single file that contains the CA, private key, and certificate.

certcat properly deconstructs this bundle even when used in Docker scratch
containers which may not have a shell environment (commonly used by proxy
software).

### Containers using scratch

- traefik
- fabio

Distribution Statement "A" (Approved for Public Release, Distribution
Unlimited).
