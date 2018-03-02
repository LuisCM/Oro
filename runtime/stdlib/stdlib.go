// Copyright 2011 The LuisCM. All rights reserved.
// Use of this source code is license that can be found in the LICENSE file.

// Package stdlib implements functions to standard library.
package stdlib

var Modules = []string{

	`module Enum

  val size = fn (array: Array) -> Int
    var count = 0
    repeat v in array
      count += 1
    end
    count
  end

  val empty? = fn (array: Array) -> Boolean
    size(array) == 0
  end

  val reverse = fn (array: Array) -> Array
    var reversed = []
    repeat i in size(array)-1..0
      reversed[] = array[i]
    end
    reversed
  end

  val first = fn (array: Array)
    array[0]
  end

  val last = fn array: Array
    array[size(array) - 1]
  end

  val insert = fn (array: Array, el) -> Array
    array[] = el
  end

  val delete = fn (array: Array, index) -> Array
    var purged = []
    repeat i, v in array
      if i != index
        purged[] = v
      end
    end
    purged
  end

  val map = fn (array: Array, fun: Function) -> Array
    repeat v in array
      fun(v)
    end
  end

  val filter = fn (array: Array, fun: Function) -> Array
    var filtered = []
    repeat v in array
      if fun(v)
        filtered[] = v
      end
    end
    filtered
  end

  val reduce = fn (array: Array, start, fun: Function)
    var acc = start
    repeat v in array
      acc = fun(v, acc)
    end
    return acc
  end

  val find = fn (array: Array, fun: Function)
    repeat v in array
      if fun(v)
        return v
      end
    end
    nil
  end

  val contains? = fn (array: Array, search) -> Boolean
    repeat v in array
      if v == search
        return true
      end
    end
    false
  end

  val unique = fn (array: Array) -> Array
    var filtered = []
    var hash = [=>]
    repeat i, v in array
      if hash[v] == nil
        hash[v] = i
        filtered[] = v
      end
    end
    filtered
  end

  val random = fn (array: Array)
    var rnd = runtime_rand(0, size(array) - 1)
    array[rnd]
  end

end`,

	`module Math

  val pi = 3.14159265359
  val e = 2.718281828459

  val floor = fn (nr: Float) -> Int
    Int(nr - nr % 1)
  end

  val ceil = fn (nr: Float) -> Int
    val rem = nr % 1
    if rem == 0
      return Int(nr)
    end
    nr > 0 ? Int(nr + (1 - rem)) : Int(nr - (1 + rem))
  end

  val max = fn (nr1, nr2)
    if !Type.isNumber?(nr1) || !Type.isNumber?(nr2)
      panic("Math.max() expects a Float or Int")
    end
    return nr1 > nr2 ? nr1 : nr2
  end

  val min = fn (nr1, nr2)
    if !Type.isNumber?(nr1) || !Type.isNumber?(nr2)
      panic("Math.min() expects a Float or Int")
    end
    return nr1 > nr2 ? nr2 : nr1
  end

  val random = fn (min: Int, max: Int) -> Int
    runtime_rand(min, max)
  end

  val abs = fn (nr)
    if !Type.isNumber?(nr)
      panic("Math.abs() expects a Float or Int")
    end
    if nr < 0
      return -nr
    end
    nr
  end

  val pow = fn (nr, exp)
    if !Type.isNumber?(nr) || !Type.isNumber?(exp)
      panic("Math.pow() expects a Float or Int")
    end
    nr ** exp
  end

end`,

	`module Type

  val of = fn x
    typeof(x)
  end

  val isNumber? = fn x
    if typeof(x) == "Float" || typeof(x) == "Int"
      return true
    end
    false
  end

  val toString = fn x
    String(x)
  end

  val toInt = fn x
    Int(x)
  end

  val toFloat = fn x
    Float(x)
  end

  val toArray = fn x
    Array(x)
  end

end`,

	`module Dictionary

  val size = fn (dict: Dictionary) -> Int
    var count = 0
    repeat v in dict
      count += 1
    end
    count
  end

  val contains? = fn (dict: Dictionary, key) -> Boolean
    repeat k, v in dict
      if k == key
        return true
      end
    end
    false
  end

  val empty? = fn (dict: Dictionary) -> Boolean
    size(dict) == 0
  end

  val insert = fn (dict: Dictionary, key, value) -> Dictionary
    if dict[key] != nil
      panic("Dictionary key '" + String(key) + "' already exists")
    end
    dict[key] = value
  end

  val update = fn (dict: Dictionary, key, value) -> Dictionary
    if dict[key] == nil
      panic("Dictionary key '" + String(key) + "' doesn't exist")
    end
    dict[key] = value
  end

  val delete = fn (dict: Dictionary, key) -> Dictionary
    if dict[key] == nil
      panic("Dictionary key '" + String(key) + "' doesn't exist")
    end
    var purged = [=>]
    repeat k, v in dict
      if k != key
        purged[k] = v
      end
    end
    purged
  end

end`,

	`module String

  val count = fn (str: String) -> Int
    var cnt = 0
    repeat v in str
      cnt += 1
    end
    cnt
  end

  val first = fn (str: String) -> String
    str[0]
  end

  val last = fn (str: String) -> String
    str[String.count(str) - 1]
  end

  val lower = fn (str: String) -> String
    runtime_tolower(str)
  end

  val upper = fn (str: String) -> String
    runtime_toupper(str)
  end

  val capitalize = fn (str: String) -> String
    var title = str
    repeat i, v in str
      if i == 0 || str[i - 1] != nil && str[i - 1] == " "
        title[i] = String.upper(v)
      end
    end
    title
  end

  val reverse = fn (str: String) -> String
    var reversed = ""
    repeat i in String.count(str)-1..0
      reversed += str[i]
    end
    reversed
  end

  val slice = fn (str: String, start: Int, length: Int) -> String
    if start < 0 || length < 0
      panic("String.slice() expects positive start and length parameters")
    end
    var sliced = ""
    var chars = 0
    repeat i, v in str
      if i >= start && chars < length
        sliced += v
        chars += 1
      end
    end
    sliced
  end

  val trim = fn (str: String, subset: String) -> String
    var trimmed = String.trimRight(String.trimLeft(str, subset), subset)
    trimmed
  end

  val trimLeft = fn (str: String, subset: String) -> String
    var trimmed = str
    repeat v in subset
      if trimmed[0] == v
        trimmed = String.slice(trimmed, 1, String.count(trimmed))
        continue
      end
    end
    trimmed
  end

  val trimRight = fn (str: String, subset: String) -> String
    var trimmed = str
    repeat v in subset
      if String.last(trimmed) == v
        trimmed = String.slice(trimmed, 0, String.count(trimmed) - 1)
        continue
      end
    end
    trimmed
  end

  val join = fn (array: Array, glue: String) -> String
    var glued = ""
    repeat v in array
      glued += v + glue
    end
    if String.count(glued) > String.count(glue)
      return String.slice(glued, 0, String.count(glued) - String.count(glue))
    end
    glued
  end

  val split = fn (str: String, separator: String) -> Array
    val count_sep = String.count(separator)
    var array = []
    var last_index = 0
    repeat i, v in str
        if String.slice(str, i, count_sep) == separator
          var curr = String.slice(str, last_index, i - last_index)
          if curr != ""
            array[] = curr
          end
          last_index = i + count_sep
        end
    end
    array[] = String.slice(str, last_index, String.count(str))
    array
  end

  val starts? = fn (str: String, prefix: String) -> Boolean
    if String.count(str) < String.count(prefix)
      return false
    end
    if String.slice(str, 0, String.count(prefix)) == prefix
      return true
    end
    false
  end

  val ends? = fn (str: String, suffix: String) -> Boolean
    if String.count(str) < String.count(suffix)
      return false
    end
    if String.slice(str, String.count(str) - String.count(suffix), String.count(str)) == suffix
      return true
    end
    false
  end

  val contains? = fn (str: String, search: String) -> Boolean
    repeat i, v in str
      if String.slice(str, i, String.count(search)) == search
        return true
      end
    end
    false
  end

  val replace = fn (str: String, search: String, replace: String) -> String
    val count_search = String.count(search)
    var rpl = ""
    var last_index = 0
    repeat i, v in str
      if String.slice(str, i, count_search) == search
        rpl = rpl + String.slice(str, last_index, i - last_index) + replace
        last_index = i + count_search
      end
    end
    rpl + String.slice(str, last_index, String.count(str))
  end

  val match? = fn (str: String, regex: String) -> Boolean
    runtime_regex_match(str, regex)
  end

end`,
}
