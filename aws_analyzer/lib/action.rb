module Analyzer

  # A service resource
  class Resource

    # [String] Resource name (e.g. "Stack")
    attr_reader :name

    # [Shape] Resource shape (structure definition)
    attr_reader :shape

    # [String] Resource primary id field (e.g. "StackId")
    attr_reader :primary_id

    # [Array<String>] Resource secondary ids field (e.g. ["StackName"])
    attr_reader :secondary_ids

    # [Array<ResourceAction>] Resource CRUD actions (index, show, update, create, delete)
    attr_reader :actions

    # [Array<ResourceAction>] Resource custom actions (e.g. cancel_update)
    attr_reader :custom_actions

    # [Array<ResourceAction>] Resource collection custom actions (e.g. list)
    attr_reader :collection_actions

    # [Hash<String, String>] Linked resource names indexed by link field name (e.g. { "stack_id" => "Stack" })
    attr_reader :links

    # Initialize
    def initialize
    end

    # YAML representation
    def to_yaml
      YAML.dump({ name:               @name,
                  shape:              @shape.to_hash,
                  primary_id:         @primary_id,
                  secondary_ids:      @secondary_ids,
                  actions:            @actions.map(&:to_hash),
                  custom_actions:     @custom_actiohns.map(&:to_hash),
                  collection_actions: @collection_actions.map(&:to_hash),
                  links:              @links })
    end
    alias :to_s :to_yaml

  end

end
