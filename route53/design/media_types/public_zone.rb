module V1
  module MediaTypes
    class PublicZone < Praxis::MediaType

      identifier 'application/vnd.rightscale.public_zone+json'

      attributes do
        attribute :id, Attributes::Route53Id
        attribute :href, String
        attribute :name, String
        attribute :caller_reference, String
        attribute :config do
          attribute :comment, String
          attribute :private_zone, String
        end
        attribute :resource_record_set_count, Integer
        attribute :change, Change

        links do
          link :change
        end
      end

      view :default do
        attribute :id
        attribute :href
        attribute :name
        attribute :caller_reference
        attribute :config
        attribute :resource_record_set_count
        attribute :links
      end

      view :link do
        attribute :href
      end

      def href()
        V1::ApiResources::PublicZone.prefix+'/'+id
      end
    end
  end
end
