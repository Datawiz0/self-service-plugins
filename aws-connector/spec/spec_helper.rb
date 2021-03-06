require 'sinatra'
require 'rack/test'
require 'logger'
require 'sinatra/json'
require 'yajl'
require 'multi_json'
require 'base64'
require "net/http"
require "net/https"
require "uri"
require "aws-sdk-core"

ENV["RACK_ENV"] ||= "development"

$logger = ::Logger.new(STDERR)
$logger.info "Hello Logger"

require_relative "../application.rb"
require_relative "../app/proxy"

RSpec.configure do |config|
  config.include Rack::Test::Methods

  def app
    Proxy
  end
end
