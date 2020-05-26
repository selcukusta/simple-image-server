vcl 4.0;

backend
default {
  .host = "web";
  .port = "8080";
}

acl purge {
  "172.18.0.1";
  "localhost";
  "127.0.0.1";
}

sub vcl_backend_response {
  if (bereq.url == "/version") {
      set beresp.uncacheable = true;
      set beresp.ttl = 1m;
      return(deliver);
  }
  
  # https://github.com/mattiasgeniar/varnish-4.0-configuration-templates/issues/24
  # Happens after we have read the response headers from the backend.
  # Here you clean the response headers and other mistakes your backend does.

  # This block will make sure that if the upstream returns a 5xx, but we have the response in the cache (even if it's expired),
  # we fall back to the cached value (until the grace period is over).
  if (beresp.status == 500 || beresp.status == 502 || beresp.status == 503 || beresp.status == 504)
  {
      # This check is important. If is_bgfetch is true, it means that we've found and returned the cached object to the client,
      # and triggered an asynchoronus background update. In that case, if it was a 5xx, we have to abandon, otherwise the previously cached object
      # would be erased from the cache (even if we set uncacheable to true).
      if (bereq.is_bgfetch)
      {
          return (abandon);
      }

      # Even if we couldn't send a previous successful response from the cache, we should never cache a 5xx response.
      set beresp.uncacheable = true;
  }
  # Set the grace time to 1 hour.
  # Gives the posibility to serve backend traffic from cache if backend is down
  set beresp.grace = 1h;
  # The object's remaining time to live.
  set beresp.ttl = 14d;
  # Set response header for browser caching.
  set beresp.http.Cache-Control = "public, max-age=1209600";
  return (deliver);
}

sub vcl_recv {
  if (req.url == "/version") {
        return(pass);
  }

  if (req.method == "PURGE") {
    if (!client.ip ~ purge) {
      return (synth(405, "Not allowed."));
    }
    return (purge);
  }

  unset req.http.Cache-Control;
  unset req.http.Max-Age;
  unset req.http.Pragma;
  unset req.http.Cookie;
  return (hash);
}

sub vcl_deliver {
  set resp.http.Hits = obj.hits;
  if (obj.hits > 0) {
    set resp.http.X-Cache = "HIT";
  }
  else {
    set resp.http.X-Cache = "MISS";
  }
  return (deliver);
}