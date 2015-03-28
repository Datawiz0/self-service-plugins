name "Configurator - Single Chef Node"
rs_ca_ver 20131202
short_description "Configurator - Single Chef Node"
long_description "This CAT uses the configurator plugin to launch a raw image as a RightScale server configured by a pre-existing Chef installation"

###########
# Namespace
###########

# For clarity sake
::IP_REGEXP = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"
::HOST_REGEXP = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$"

namespace "cm" do
  service do
    host  "54.184.12.120:8000" #"cm.test.rightscale.com"
    path "/cm/accounts/:account_id"
    headers do {
      "X-Api-Version" => "1.0",
      "X-Secret" => "R9Bt4GMqQBT3UZREoW7aaACUnWLWdGqO"
    } end
  end

  type "chef_configuration" do
    provision "provision_chef_configuration"
    delete "delete_resource"
    fields do
      field "chef_server_url" do
        type "string"
        required true
      end
      field "node_name" do
        type "string"
        required true
      end
      field "validation_client_name" do
        type "string"
        required true
      end
      field "validation_key" do
        type "string"
        required true
      end
      field "chef_environment" do
        type "string"
        required true
      end
      field "run_list" do
        type "array"
        required true
      end
      field "first_attributes" do
        type "composite"
      end
    end
  end

  type "booter" do
    provision "provision_booter"
    delete "delete_resource"
    fields do
      field "host" do
        type "string"
        required "true"
        regexp "(?:#{::IP_REGEXP}|#{::HOST_REGEXP})"
      end
      field "ssh_key" do
        type "resource"
      end
    end
  end
end

# Define the RCL definitions to create and destroy the resource
define provision_chef_configuration(@raw_conf) return @conf do
  $obj = to_object(@raw_conf)
  $fields = $obj["fields"]
  @conf = cm.chef_configuration.create($fields) # Calls .create on the API resource
end

define provision_booter(@raw_booter) return @booter do
  $obj = to_object(@raw_booter)
  $fields = $obj["fields"]
  @conf = cm.booter.create($fields) # Calls .create on the API resource
end

define delete_resource(@resource) do
  @resource.destroy() # Calls .delete on the API resource
end

#########
# Parameters
#########
parameter "chef_server_url" do
  type "string"
  label "Chef Server URL"
  category "Chef"
  default "https://api.opscode.com/organizations/rs-st-dev"
end

parameter "validation_client_name" do
  type "string"
  label "Chef Validator Name"
  description "The name of the Chef validator"
  default "rs-st-dev-validator"
end

parameter "validation_key" do
  type "string"
  label "Chef Validation Key"
  description "Name of RightScale credential holding the Chef server validation key"
  default "KM_CHEF_VALIDATION_KEY"
end

parameter "environment" do
  type "string"
  label "Chef environment"
  category "Chef"
  default "_default"
end

parameter "run_list" do
  type "list"
  label "Boot run list"
  category "Chef"
end

# parameter "first_attributes" do
#   type "string"
#   label "Attributes used to run initial configuration"
#   category "Chef"
# end

#########
# Resources
#########
resource "chef_cm", type: "cm.chef_configuration" do
  node_name              "test_node"
  chef_server_url        $chef_server_url
  validation_client_name $validation_client_name
  validation_key         $validation_key
  run_list               $run_list
  chef_environment       $environment
  first_attributes       {}
end

resource "cm_server", type: "server" do
  name "cm_server"
  cloud_href "/api/clouds/6"
  instance_type "m1.small"
  ssh_key "default"
  user_data "@chef_cm.bootstrap_script" # server must be tagged with 'rs_agent:userdata=mime'
  server_template find('RightLink 10.0.3 Linux Base') # Could be anything
end
