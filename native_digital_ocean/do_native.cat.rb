name "Digital Ocean Self-Service Namespace"
rs_ca_ver 20150729
short_description "Digital Ocean integration"
long_description "This CAT defines a namespace and accompanying definitions that makes it " +
                 "possible to interact with all the resources exposed by the Digital Ocean " +
                 "platform."

namespace "do" do

  # Name of RightScale credential used to authenticate with the Digital Ocean APIs
  auth_credentials "DO_TOKEN"

  # Description of the Digital Ocean APIs endpoint
  service do
    host "https://api.digitalocean.com"
    path "/v2"
    headers do {
      "Authorization" => "Bearer $DO_TOKEN"
    } end
    no_cert_check false # Do check the endpoint SSL cert (default)
  end

  type "droplet" do

    # `prefix` defines the path prefix for the resource href and actions.
    # Individual actions may override the prefix by speciyfing an "absolute" path.
    # The prefix may define variables that get substituted using values passed as argument
    # to the actions. See the "record" type below for an example.
    prefix "/droplets"

    # `href` defines the patterns used to build one or more hrefs from a response JSON.
    # Each pattern may use JMESPath expressions delimited with the ':' character.
    # The pattern JMESPaths are matched against the response and the first pattern for which all
    # JMESPath sections match is used to build the response resource collection.
    # Individual actions may override the pattern, see `action` below.
    # There are two use cases where the href pattern is used (in both cases the action is either
    # `get`, `create` or a custom action):
    #   - actions where the resource type is known: `@res.action()` or `namespace.type.action()`.
    #     The type href pattern is matched against the response body to build the resulting
    #     collection if any. An error is returned if the action is `get` and no pattern matches.
    #   - `get` on the namespace: `namespace.get()`. The namespace type href patterns are matched
    #     against the response. The longest pattern that matches "wins" and the corresponding type
    #      is used to build the resource collection. An error is returned if none matches.
    href "/:droplet.id:", "/:droplets[*].id:"

    # `provision` sets the name of the definition invoked by the `provision` RCL function.
    # If not specified the generic provision definition is used.
    provision "provision_droplet"

    # `delete` sets the name of the definition invoked by the `delete` RCL function
    # If not specified the generic delete definition is used.
    delete "delete_droplet"

    # `field` declares a resource field.
    # resource fields are used in CATs to declare resources and correspond to fields sent to
    # create a resource. Not to be confused with `output` (see below).
    field "name" do
      type "string"
      required true
    end
    field "region" do
      type "string"
      required true
    end
    field "size" do
      type "string"
      required true
    end
    field "image" do
      type "number"
      required true
    end
    field "ssh_keys" do
      type "array"
    end
    field "backups" do
      type "boolean"
    end
    field "ipv6" do
      type "boolean"
    end
    field "private_networking" do
      type "boolean"
    end
    field "user_data" do
      type "string"
    end

    # `output` define the fields that can be used to define:
    #   - CAT outputs
    #   - other CAT resource fields (dependent resources)
    #   - the action `href` patterns
    output "id", "name", "memory", "vcpus", "disk", "locked", "created_at", "status", "features",
           "region", "image", "size", "size_slug", "networks", "kernel", "next_backup_window"

    # `action` defines a supported action including:
    #   - the HTTP verb (defaults to GET for list and show, PUT for update, DELETE for destroy
    #     and POST for everything else)
    #   - the URL path (defaults to href for show, update and destroy, collection href for create
    #     and list, and '/actions/<name>' for all others). The path may include any of the type
    #     output fields. Fields are specified by prefixing their names with the ':' character.
    #   - the type returned by the action if any (defaults to this type for list, show and create)
    #   - the pattern used to build the href if type is set and the pattern differs from the type
    #     default patterns.
    #  If none of the CRUD actions is listed then the system assumes all are supported. CRUD
    #  actions are list, show, create, update and destroy.
    action "list", "show", "create", "destroy"
    action "operate" do
      path "/droplets/:id/actions"
      type "action"
    end
    action "kernels" do
      verb "GET"
      path "/droplets/:id/kernels"
    end

    # `link` defines a link to other resources.
    # Links can be use in RCL e.g. `@droplet.actions()`
    # `link` specifies both the type of the resource(s) being linked to and the pattern used to
    # build the corresponding href(s). Links may also specify an absolute URL instead of an href
    # via the `url` field, see the "next" and "last" links below as an example.
    # Note that links are always retrieved via a GET request and must always return resource
    # collections (might be empty, the point is no raw value).
    link "actions" do
      type "action"
      href "/droplets/:id/actions"
    end
    link "snapshots" do
      type "image"
      href "/droplets/:id/snapshots"
    end
    link "backups" do
      type "image"
      href "/droplets/:id/backups"
    end
    link "neighbors" do
      type "droplet"
      href "/droplets/:id/neighbors"
    end
    link "last" do
      type "droplet"
      url "links.pages.last"
    end
    link "next" do
      type "droplet"
      url "links.pages.next"
    end

  end

  type "domain" do
    prefix "/domains"
    href "/:domain.name:", "/:domains[*].name:"

    field "name" do
      type "string"
      required true
    end

    field "ip_address" do
      type "string"
      required true
      regex "(([A-Fa-f0-9]{1,4}:){7}[A-Fa-f0-9]{1,4}|([0-9]{1,3}\.){3}[0-9]{1,3})"
    end

    output "name", "ttl", "zone_file"

    action "list", "show", "create", "delete"
  end

  type "record" do
    prefix "/:domain_name/records" # all actions must pass a "domain_name" argument
    href "/:domain_record.id:", "/:domain_records[*].id:"

    field "type" do
      type "string"
      required true
      enum "A", "MX", "CNAME", "TXT", "AAAA", "SRV", "NS"
    end

    field "name" do
      type "string"
    end

    field "data" do
      type "string"
    end

    field "priority" do
      type "number"
    end

    field "port" do
      type "number"
    end

    field "weight" do
      type "number"
    end

    output "id", "type", "name", "data", "priority", "port", "weight"
  end

  type "action" do
    prefix "/actions"
    href "/:action.id:", "/:actions[*].id:"
    output "id", "status", "type", "started_at", "completed_at", "resource_id", "resource_type",
           "region", "region_slug"
    action "list", "show"

    link "last" do
      type "action"
      url "links.pages.last"
    end
    link "next" do
      type "action"
      url "links.pages.next"
    end

  end

  type "image" do
    prefix "/images"
    href "/:images[*].id:"
    output "id", "name", "type", "distribution", "slug", "public", "regions", "min_disk_size"

    action "list", "show", "update", "delete"
    action "operate" do
      path "/images/:id/actions"
      type "action"
    end

    link "last" do
      type "image"
      url "links.pages.last"
    end
    link "next" do
      type "image"
      url "links.pages.next"
    end

  end

  type "ssh_key" do
    prefix "/account/keys"
    href "/:ssh_key.id:", "/:ssh_keys[*].id:"

    field "name" do
      type "string"
      required true
    end
    field "public_key" do
      type "string"
      required true
    end

    output "id", "name", "fingerprint", "public_key"

    link "last" do
      type "ssh_key"
      url "links.pages.last"
    end
    link "next" do
      type "ssh_key"
      url "links.pages.next"
    end
  end

  type "region" do
    prefix "/regions"
    href "/:regions[*].slug:"
    output "slug", "name", "sizes", "features", "available", "sizes"
    action "list"
  end

  type "size" do
    prefix "/sizes"
    href "/:sizes[*].slug:"
    output "slug", "memory", "vcpus", "disk", "transfer", "price_monthly", "price_hourly",
           "available", "regions"
    action "list"
  end

end
