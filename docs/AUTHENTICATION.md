# Authentication Documentation

Authentication for Contineo is handled by an external web server, so auth can be entirely modular. This in turn allows the admin to decide what the controller authenticates against. The design for the webserver will be minimal as so that it can easily be implemented for custom systems where admins want users to be able to log in with a preexisting username and password.

## Authentication Server

The authetntication server only needs to implement 1 method, `/authenticate`, which must accept a POST requet with a JSON payload. The payload must then contain a authentication token which is then processed. Upon processing, if valid, a session token is generated and returned. If invalid, no session token is returned. (Possibly may implement a system with a auth token and a refresh token)

### Example Use Case

For example, if implementing a client into a website, a login system for the site may already be in place. Rather than faffing with converting user account to the Contineo server and keeping them up to date (i,e, managing two login systems), a authetncation server can be created (or even better, just a new endpoint developed into an existing web API used by the website. Contineo will then be configured to use this endpoint as it's authentication server. In the case of an existing login system, this endpoint could then take the token provided and authenticate it with the existing login system, therefore verifying it's validity. The user can then be provided a session token for Contineo.