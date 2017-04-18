
cors gives you easy control over Cross Origin Resource Sharing for your site.

It allows you to whitelist particular domains per route, or to simply allow all domains `*` If desired you may customize nearly every aspect of the specification.

### Syntax

```
cors [path] [domains...] {
	origin            [origin]
	methods           [methods]
	allow_credentials [allowCredentials]
	max_age           [maxAge]
	allowed_headers   [allowedHeaders]
	exposed_headers   [exposedHeaders]
}
```

*   **path** is the file or directory this applies to (default is /).
*   **domains** is a space-seperated list of domains to allow. If ommitted, all domains will be granted access.
*   **origin** is a domain to grant access to. May be specified multiple times or ommitted.
*   **methods** is set of http methods to allow. Default is these: POST,GET,OPTIONS,PUT,DELETE.
*   **allow_credentials** sets the value of the Access-Control-Allow-Credentials header. Can be true or false. By default, header will not be included.
*   **max_age** is the length of time in seconds to cache preflight info. Not set by default.
*   **allowed_headers** is a comma-seperated list of request headers a client may send.
*   **exposed_headers** is a comma-seperated list of response headers a client may access.

### Examples

Simply allow all domains to request any path:

<code class="block"><span class="hl-directive">cors</span></code>

Protect specific paths only, and only allow a few domains:

<code class="block"><span class="hl-directive">cors</span> <span class="hl-arg">/foo http://mysite.com http://anothertrustedsite.com</span></code>

Full configuration:

<code class="block"><span class="hl-directive">cors</span> <span class="hl-arg">/</span> {
    <span class="hl-subdirective">origin</span>            http://allowedSite.com
	<span class="hl-subdirective">origin</span>            http://anotherSite.org https://anotherSite.org
    <span class="hl-subdirective">methods</span>           POST,PUT
    <span class="hl-subdirective">allow_credentials</span> false
	<span class="hl-subdirective">max_age</span>           3600
	<span class="hl-subdirective">allowed_headers</span>   X-Custom-Header,X-Foobar
	<span class="hl-subdirective">exposed_headers</span>   X-Something-Special,SomethingElse
}</code>
