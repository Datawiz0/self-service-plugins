module Analyzer

  # Generic analyzer
  #
  # Example usage:
  #
  # a = Analyzer::Analyzer.new(cloud: 'aws')
  # definition = a.service_definition('cloud_formation')
  class Analyzer

    # Initialize
    #
    # @option [String] :cloud Service cloud provider, 'aws' or 'gce'
    def initialize(options)
      cloud = options.delete(:cloud) || 'aws'
      klass = ::Analyzer.const_get(cloud.upcase).const_get('Analyzer') rescue nil
      if klass.nil?
        puts "No analyzer implemented yet for #{cloud.inspect} - exiting..."
        exit 1
      end
      @analyzer = klass.new(options)
    end

    # Return the service definition for the given service
    def service_definition(service)
      @analyzer.analyze(service)
    end

    # Analyze service from definition in hash
    # @options :force
    # @options :resource_only
    def analyze_service(service, options)
      force = options[:force]
      begin
        definition = service_definition(service)
        errors = @analyzer.errors
      rescue Exception => e
        if e.is_a? SystemExit
          raise
        end
        errors = [e.message + " from\n" + e.backtrace.join("\n")]
      end
      if errors && !errors.empty?
        if errors.size == 1
          puts "ERROR: #{errors.first}"
        else
          puts "ERRORS:\n#{errors.join("\n")}"
        end
        exit 1 if !force
        puts
      end
      return if definition.nil? || definition.empty?
      if options[:resource_only]
        hash = definition.to_hash
        hash.delete('shapes')
        puts YAML.dump(hash)
      else
        puts definition.to_yaml
      end
    end

    # YAML representation
    def to_yaml
      YAML.dump({ name:               @name,
                  shape:              @shape.to_hash,
                  primary_id:         @primary_id,
                  secondary_ids:      @secondary_ids,
                  actions:            @actions.map(&:to_hash),
                  custom_actions:     @custom_actions.map(&:to_hash),
                  collection_actions: @collection_actions.map(&:to_hash),
                  links:              @links })
    end
    alias :to_s :to_yaml

  end

end
