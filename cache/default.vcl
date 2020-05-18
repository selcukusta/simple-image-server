vcl 4.0;

backend
default {
  .host = "backend";
  .port = "8080";
}

acl purge {
  "172.18.0.1";
  "localhost";
  "127.0.0.1";
}

sub vcl_backend_response {
  set beresp.ttl = 14d;
  set beresp.http.Cache-Control = "public, max-age=1209600";
  return (deliver);
}

sub vcl_recv {
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