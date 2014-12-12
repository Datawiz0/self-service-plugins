module Analyzer

  module AWS

    # Resources that have a plural name
    PLURAL_RESOURCE_NAMES = ['DhcpOptions']

    # A service resource, the main point of this class is to make sure that we can easily identify operations that apply
    # to either the resourcd or the collection.
    class Resource

      attr_reader :name

      # Initialize with resource name
      def initialize(name)
        @name               = name.underscore
        @orig_name          = name
        @actions            = {}
        @collection_actions = {}
        @custom_actions     = {}
      end

      # Register operation
      # OK, here is the trick:
      # name is the CamelCase name of the operation, this name ends with either the ResourceName or ResourceNames
      # for operations that apply to the collection. We detect which one it is and then infer the final operation
      # name and type (resource, collection or custom) from that.
      def add_operation(op)
        name = op['name']
        is_collection = name !~ /#{@orig_name}$/ # @orig_name is the singular version of ResourceName
        n = name.gsub(/(#{@orig_name}|#{@orig_name.pluralize})$/, '').underscore
        if n == 'describe'
          n = is_collection ? 'index' : 'show'
        end
        operation = to_operation(op, n)
        if is_collection
          if n == 'index'
            @actions['index'] = operation
          else
            @collection_actions[n] = operation
          end
        else
          if n == 'show'
            # Let's set the shape of the resource with the result of a describe
            @shape = op['output']['shape']
            if @shape.nil?
              raise "No shape for describe??? Resource: #{name}, Operation: #{op['name']}"
            end

          end
          if ['create', 'delete', 'update', 'show'].include?(n)
            @actions[n] = operation
          else
            @custom_actions[n] = operation
          end
        end
      end

      # Map raw JSON operation to analyzed YAML operation
      # e.g.
      #    - name: DescribeStackResource
      #      http:
      #        method: POST
      #        requestUri: "/"
      #      input:
      #        shape: DescribeStackResourceInput
      #      output:
      #        shape: DescribeStackResourceOutput
      #        resultWrapper: DescribeStackResourceResult
      # becomes
      #    - name: show
      #      verb: post
      #      path: "/"
      #      payload: describe_stack_resource_input
      #      params:
      #      response: describe_stack_resource_output
      def to_operation(op, name)
        { 'name'          => name,
          'original_name' => op['name'],
          'verb'          => op['http']['method'].downcase,
          'path'          => op['http']['requestUri'],
          'payload'       => op['input']['shape'].underscore,
          'params'        => [],
          'response'      => (out = op['output']) && out['shape'].underscore }
      end

      # Hashify
      def to_hash
        { 'name'               => @name,
          'shape'              => @shape,
          'primary_id'         => @primary_id,
          'secondary_ids'      => @secondary_ids,
          'actions'            => @actions,
          'custom_actions'     => @custom_actions,
          'collection_actions' => @collection_actions }
      end

    end

    # Registry of resources for a given service
    class ResourceRegistry

      # Resources indexed by name
      attr_reader :resources

      def initialize
        @resources = {}
      end

      # Add operation to resource
      # Create resource if non-existent, checks whether operation is collection or resource operation
      def add_resource_operation(res_name, op)
        canonical = canonical_name(res_name)
        res = @resources[canonical] ||= Resource.new(canonical)
        res.add_operation(op)
      end

      # Known resource names
      def resource_names
        @resources.keys
      end

      # Singularize name unless it's in the exception list
      def canonical_name(base_name)
        PLURAL_RESOURCE_NAMES.include?(base_name) ? base_name : base_name.singularize
      end

    end

  end

end
