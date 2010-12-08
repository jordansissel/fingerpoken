#!/usr/bin/env ruby

require "rubygems"
require "ffi"
require "fingerpoken/target"

class FingerPoken::Target::Xdo < FingerPoken::Target
  module LibXdo
    extend FFI::Library
    ffi_lib "libxdo.so"

    attach_function :xdo_new, [:string], :pointer
    attach_function :xdo_mousemove, [:pointer, :int, :int, :int], :int
    attach_function :xdo_mousemove_relative, [:pointer, :int, :int], :int
    attach_function :xdo_click, [:pointer, :long, :int], :int
    attach_function :xdo_mousedown, [:pointer, :long, :int], :int
    attach_function :xdo_mouseup, [:pointer, :long, :int], :int
    attach_function :xdo_type, [:pointer, :long, :string, :long], :int
    attach_function :xdo_keysequence, [:pointer, :long, :string, :long], :int
  end

  def initialize(config)
    super(config)
    @xdo = LibXdo::xdo_new(nil)
    if @xdo.null?
       raise "xdo_new failed"
    end
    register
  end

  def mousemove_relative(x, y)
    @logger.info("move #{x},#{y}")
    return LibXdo::xdo_mousemove_relative(@xdo, x, y)
  end

  def click(button)
    LibXdo::xdo_click(@xdo, 0, button)
  end

  def mousedown(button)
    LibXdo::xdo_mousedown(@xdo, 0, button)
  end

  def mouseup(button)
    LibXdo::xdo_mouseup(@xdo, 0, button)
  end

  def type(string)
    LibXdo::xdo_type(@xdo, 0, string, 12000)
  end

  def keypress(key)
    key = request["key"]
    if key.is_a?(String)
      if key.length == 1
        # Assume letter
        LibXdo::xdo_type(@xdo, 0, key, 12000)
      else
        # Assume keysym
        LibXdo::xdo_keysequence(@xdo, 0, key, 12000)
      end
    else
      # type printables, key others.
      if 32.upto(127).include?(key)
        LibXdo::xdo_type(@xdo, 0, request["key"].chr, 12000)
      else
        case key
          when 8 
            LibXdo::xdo_keysequence(@xdo, 0, "BackSpace", 12000)
          when 13
            LibXdo::xdo_keysequence(@xdo, 0, "Return", 12000)
          else
            puts "I don't know how to type web keycode '#{key}'"
          end # case key
      end # if 32.upto(127).include?(key)
    end # if key.is_a?String
  end # def keypress
end # class FingerPoken::Target::Xdo 
