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
    attach_function :xdo_get_window_size, [:pointer, :long, :pointer, :pointer], :int
    attach_function :xdo_window_search, [:pointer, :pointer, :pointer, :pointer], :int
  end

  class XdoSearch < FFI::Struct
    layout :title, :pointer,
           :winclass, :pointer,
           :winclassname, :pointer,
           :winname, :pointer,
           :pid, :int,
           :max_depth, :long,
           :only_visible, :int,
           :screen, :int,
           :require, :int,
           :searchmask, :uint,
           :desktop, :long
  end # class XdoSearch

  def initialize(config)
    super(config)
    @xdo = LibXdo::xdo_new(nil)
    if @xdo.null?
       raise "xdo_new failed"
    end

    search = XdoSearch.new
    search[:searchmask] = 1 << 2 # SEARCH_NAME, from xdo.h
    search[:max_depth] = 0
    search[:winname] = FFI::MemoryPointer.new(:char, 3)
    search[:winname].put_string(0, ".*")
    ptr_nwindows = FFI::MemoryPointer.new(:ulong, 1)
    ptr_winlist = FFI::MemoryPointer.new(:pointer, 1)
    LibXdo::xdo_window_search(@xdo, search, ptr_winlist, ptr_nwindows)
    nwindows = ptr_nwindows.read_long
    @rootwin = ptr_winlist.read_pointer.read_array_of_long(nwindows)[0]

    ptr_x = FFI::MemoryPointer.new(:int, 1)
    ptr_y = FFI::MemoryPointer.new(:int, 1)

    LibXdo::xdo_get_window_size(@xdo, @rootwin, ptr_x, ptr_y)
    @screen_x = ptr_x.read_int
    @screen_y = ptr_y.read_int
  end

  def mousemove_relative(x, y)
    return LibXdo::xdo_mousemove_relative(@xdo, x, y)
  end

  def mousemove_absolute(px, py)
    # Edges may be hard to hit on some devices, so inflate things a bit.
    xbuf = @screen_x * 0.1
    ybuf = @screen_y * 0.1
    x = (((@screen_x + xbuf) * px) - (xbuf / 2)).to_i
    y = (((@screen_y + ybuf) * py) - (ybuf / 2)).to_i

    return LibXdo::xdo_mousemove(@xdo, x, y, 0)
  end

  def click(button)
    return LibXdo::xdo_click(@xdo, 0, button.to_i)
  end

  def mousedown(button)
    return LibXdo::xdo_mousedown(@xdo, 0, button.to_i)
  end

  def mouseup(button)
    return LibXdo::xdo_mouseup(@xdo, 0, button.to_i)
  end

  def type(string)
    return LibXdo::xdo_type(@xdo, 0, string, 12000)
  end

  def keypress(key)
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
        LibXdo::xdo_type(@xdo, 0, key.chr, 12000)
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
    return nil
  end # def keypress
end # class FingerPoken::Target::Xdo 
