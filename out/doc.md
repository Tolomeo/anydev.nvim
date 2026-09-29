# ArrayIter

 Special case implementations for iterators on list tables.

## __call


```lua
function ArrayIter.__call(self: ArrayIter)
  -> unknown
  2. unknown|nil
  3. unknown|nil
  4. unknown|nil
  5. unknown|nil
  6. unknown|nil
  7. unknown|nil
  8. unknown|nil
  9. unknown|nil
 10. unknown|nil
```

## __index


```lua
ArrayIter|Iter
```

 Special case implementations for iterators on list tables.

## _head


```lua
number
```

Index to the front of a table iterator

## _table


```lua
table
```

Underlying table data

## _tail


```lua
number
```

Index to the end of a table iterator (exclusive)

## all


```lua
(method) Iter:all(pred: fun(...any):boolean)
  -> boolean
```

 Returns true if all items in the iterator match the given predicate.

@*param* `pred` — Predicate function. Takes all values returned from the previous

                          stage in the pipeline as arguments and returns true if the
                          predicate matches.

## any


```lua
(method) Iter:any(pred: fun(...any):boolean)
  -> boolean
```

 Returns true if any of the items in the iterator match the given predicate.

@*param* `pred` — Predicate function. Takes all values returned from the previous

                          stage in the pipeline as arguments and returns true if the
                          predicate matches.

## each


```lua
(method) ArrayIter:each(f: fun(...any))
```

## enumerate


```lua
(method) ArrayIter:enumerate()
  -> ArrayIter
```

## filter


```lua
(method) ArrayIter:filter(f: fun(...any):boolean)
  -> ArrayIter
```

## find


```lua
(method) Iter:find(f: any)
  -> any
```

 Find the first value in the iterator that satisfies the given predicate.

 Advances the iterator. Returns nil and drains the iterator if no value is found.

 Examples:

 ```lua

 local it = vim.iter({ 3, 6, 9, 12 })
 it:find(12)
 -- 12

 local it = vim.iter({ 3, 6, 9, 12 })
 it:find(20)
 -- nil

 local it = vim.iter({ 3, 6, 9, 12 })
 it:find(function(v) return v % 4 == 0 end)
 -- 12

 ```

## flatten


```lua
(method) ArrayIter:flatten(depth: number)
  -> ArrayIter
```

## fold


```lua
(method) ArrayIter:fold(init: any, f: fun(acc: <A>, ...any):<A>)
  -> unknown
```

## join


```lua
(method) Iter:join(delim: string)
  -> string
```

 Collect the iterator into a delimited string.

 Each element in the iterator is joined into a string separated by {delim}.

 Consumes the iterator.

@*param* `delim` — Delimiter

## last


```lua
(method) ArrayIter:last()
  -> unknown|nil
```

## map


```lua
(method) ArrayIter:map(f: fun(...any):...any)
  -> ArrayIter
```

## new


```lua
function ArrayIter.new(t: table)
  -> Iter
```

 Create a new ArrayIter

@*param* `t` — Array-like table. Caller guarantees that this table is a valid array. Can have

               holes (nil values).

## next


```lua
(method) ArrayIter:next()
  -> unknown
  2. unknown|nil
  3. unknown|nil
  4. unknown|nil
  5. unknown|nil
  6. unknown|nil
  7. unknown|nil
  8. unknown|nil
  9. unknown|nil
 10. unknown|nil
```

## nth


```lua
(method) Iter:nth(n: number)
  -> any
```

 Gets the nth value of an iterator (and advances to it).

 If `n` is negative, offsets from the end of a |list-iterator|.

 Example:

 ```lua
 local it = vim.iter({ 3, 6, 9, 12 })
 it:nth(2)
 -- 6
 it:nth(2)
 -- 12

 local it2 = vim.iter({ 3, 6, 9, 12 })
 it2:nth(-2)
 -- 9
 it2:nth(-2)
 -- 3
 ```

@*param* `n` — Index of the value to return. May be negative if the source is a |list-iterator|.

## peek


```lua
(method) ArrayIter:peek()
  -> unknown
```

## pop


```lua
(method) ArrayIter:pop()
  -> unknown
```

 @nodoc

## rev


```lua
(method) ArrayIter:rev()
  -> ArrayIter
```

## rfind


```lua
(method) ArrayIter:rfind(f: any)
  -> unknown
  2. unknown|nil
  3. unknown|nil
  4. unknown|nil
  5. unknown|nil
  6. unknown|nil
  7. unknown|nil
  8. unknown|nil
  9. unknown|nil
 10. unknown|nil
```

## rpeek


```lua
(method) ArrayIter:rpeek()
  -> unknown
```

## rskip


```lua
(method) ArrayIter:rskip(n: number)
  -> ArrayIter
```

## skip


```lua
(method) ArrayIter:skip(n: number)
  -> ArrayIter
```

## slice


```lua
(method) ArrayIter:slice(first: number, last: number)
  -> ArrayIter
```

## take


```lua
(method) ArrayIter:take(n: integer)
  -> ArrayIter
```

## totable


```lua
(method) ArrayIter:totable()
  -> table
```


---

# ElementType

 This loop reads each line, putting them into stacks with some extra data since
 rendering each line requires understanding what is beneath it.


---

# Iter

## __call


```lua
function Iter.__call(self: Iter)
  -> unknown
```

## __index


```lua
Iter
```

## all


```lua
(method) Iter:all(pred: fun(...any):boolean)
  -> boolean
```

 Returns true if all items in the iterator match the given predicate.

@*param* `pred` — Predicate function. Takes all values returned from the previous

                          stage in the pipeline as arguments and returns true if the
                          predicate matches.

## any


```lua
(method) Iter:any(pred: fun(...any):boolean)
  -> boolean
```

 Returns true if any of the items in the iterator match the given predicate.

@*param* `pred` — Predicate function. Takes all values returned from the previous

                          stage in the pipeline as arguments and returns true if the
                          predicate matches.

## each


```lua
(method) Iter:each(f: fun(...any))
```

 Calls a function once for each item in the pipeline, draining the iterator.

 For functions with side effects. To modify the values in the iterator, use |Iter:map()|.

@*param* `f` — Function to execute for each item in the pipeline.

                  Takes all of the values returned by the previous stage
                  in the pipeline as arguments.

## enumerate


```lua
(method) Iter:enumerate()
  -> Iter
```

 Yields the item index (count) and value for each item of an iterator pipeline.

 For list tables, this is more efficient:

 ```lua
 vim.iter(ipairs(t))
 ```

 instead of:

 ```lua
 vim.iter(t):enumerate()
 ```

 Example:

 ```lua

 local it = vim.iter(vim.gsplit('abc', '')):enumerate()
 it:next()
 -- 1	'a'
 it:next()
 -- 2	'b'
 it:next()
 -- 3	'c'

 ```

## filter


```lua
(method) Iter:filter(f: fun(...any):boolean)
  -> Iter
```

 Filters an iterator pipeline.

 Example:

 ```lua
 local bufs = vim.iter(vim.api.nvim_list_bufs()):filter(vim.api.nvim_buf_is_loaded)
 ```

@*param* `f` — Takes all values returned from the previous stage

                       in the pipeline and returns false or nil if the
                       current iterator element should be removed.

## find


```lua
(method) Iter:find(f: any)
  -> any
```

 Find the first value in the iterator that satisfies the given predicate.

 Advances the iterator. Returns nil and drains the iterator if no value is found.

 Examples:

 ```lua

 local it = vim.iter({ 3, 6, 9, 12 })
 it:find(12)
 -- 12

 local it = vim.iter({ 3, 6, 9, 12 })
 it:find(20)
 -- nil

 local it = vim.iter({ 3, 6, 9, 12 })
 it:find(function(v) return v % 4 == 0 end)
 -- 12

 ```

## flatten


```lua
(method) Iter:flatten(depth?: number)
  -> Iter
```

 Flattens a |list-iterator|, un-nesting nested values up to the given {depth}.
 Errors if it attempts to flatten a dict-like value.

 Examples:

 ```lua
 vim.iter({ 1, { 2 }, { { 3 } } }):flatten():totable()
 -- { 1, 2, { 3 } }

 vim.iter({1, { { a = 2 } }, { 3 } }):flatten():totable()
 -- { 1, { a = 2 }, 3 }

 vim.iter({ 1, { { a = 2 } }, { 3 } }):flatten(math.huge):totable()
 -- error: attempt to flatten a dict-like table
 ```

@*param* `depth` — Depth to which |list-iterator| should be flattened

                        (defaults to 1)
 luacheck: no unused args

## fold


```lua
(method) Iter:fold(init: <A>, f: fun(acc: <A>, ...any):<A>)
  -> <A>
```

 Folds ("reduces") an iterator into a single value. [Iter:reduce()](file:///usr/local/share/nvim/runtime/lua/vim)

 Examples:

 ```lua
 -- Create a new table with only even values
 vim.iter({ a = 1, b = 2, c = 3, d = 4 })
   :filter(function(k, v) return v % 2 == 0 end)
   :fold({}, function(acc, k, v)
     acc[k] = v
     return acc
   end) --> { b = 2, d = 4 }

 -- Get the "maximum" item of an iterable.
 vim.iter({ -99, -4, 3, 42, 0, 0, 7 })
   :fold({}, function(acc, v)
     acc.max = math.max(v, acc.max or v)
     return acc
   end) --> { max = 42 }
 ```


@*param* `init` — Initial value of the accumulator.

@*param* `f` — Accumulation function.

## join


```lua
(method) Iter:join(delim: string)
  -> string
```

 Collect the iterator into a delimited string.

 Each element in the iterator is joined into a string separated by {delim}.

 Consumes the iterator.

@*param* `delim` — Delimiter

## last


```lua
(method) Iter:last()
  -> any
```

 Drains the iterator and returns the last item.

 Example:

 ```lua

 local it = vim.iter(vim.gsplit('abcdefg', ''))
 it:last()
 -- 'g'

 local it = vim.iter({ 3, 6, 9, 12, 15 })
 it:last()
 -- 15

 ```

## map


```lua
(method) Iter:map(f: fun(...any):...any)
  -> Iter
```

 Maps the items of an iterator pipeline to the values returned by `f`.

 If the map function returns nil, the value is filtered from the iterator.

 Example:

 ```lua
 local it = vim.iter({ 1, 2, 3, 4 }):map(function(v)
   if v % 2 == 0 then
     return v * 3
   end
 end)
 it:totable()
 -- { 6, 12 }
 ```

@*param* `f` — Mapping function. Takes all values returned from

                      the previous stage in the pipeline as arguments
                      and returns one or more new values, which are used
                      in the next pipeline stage. Nil return values
                      are filtered from the output.

## new


```lua
function Iter.new(src: function|table, ...any)
  -> Iter
```

 Creates a new Iter object from a table or other |iterable|.

@*param* `src` — Table or iterator to drain values from

## next


```lua
function Iter.next()
  -> unknown
```

## nth


```lua
(method) Iter:nth(n: number)
  -> any
```

 Gets the nth value of an iterator (and advances to it).

 If `n` is negative, offsets from the end of a |list-iterator|.

 Example:

 ```lua
 local it = vim.iter({ 3, 6, 9, 12 })
 it:nth(2)
 -- 6
 it:nth(2)
 -- 12

 local it2 = vim.iter({ 3, 6, 9, 12 })
 it2:nth(-2)
 -- 9
 it2:nth(-2)
 -- 3
 ```

@*param* `n` — Index of the value to return. May be negative if the source is a |list-iterator|.

## peek


```lua
(method) Iter:peek()
  -> any
```

 Gets the next value in a |list-iterator| without consuming it.

 Example:

 ```lua

 local it = vim.iter({ 3, 6, 9, 12 })
 it:peek()
 -- 3
 it:peek()
 -- 3
 it:next()
 -- 3

 ```

## pop


```lua
(method) Iter:pop()
  -> any
```

 "Pops" a value from a |list-iterator| (gets the last value and decrements the tail).

 Example:

 ```lua
 local it = vim.iter({1, 2, 3, 4})
 it:pop()
 -- 4
 it:pop()
 -- 3
 ```

## rev


```lua
(method) Iter:rev()
  -> Iter
```

 Reverses a |list-iterator| pipeline.

 Example:

 ```lua

 local it = vim.iter({ 3, 6, 9, 12 }):rev()
 it:totable()
 -- { 12, 9, 6, 3 }

 ```

## rfind


```lua
(method) Iter:rfind(f: any)
  -> any
```

 Gets the first value satisfying a predicate, from the end of a |list-iterator|.

 Advances the iterator. Returns nil and drains the iterator if no value is found.

 Examples:

 ```lua

 local it = vim.iter({ 1, 2, 3, 2, 1 }):enumerate()
 it:rfind(1)
 -- 5	1
 it:rfind(1)
 -- 1	1

 ```


 luacheck: no unused args

## rpeek


```lua
(method) Iter:rpeek()
  -> any
```

 Gets the last value of a |list-iterator| without consuming it.

 Example:

 ```lua
 local it = vim.iter({1, 2, 3, 4})
 it:rpeek()
 -- 4
 it:rpeek()
 -- 4
 it:pop()
 -- 4
 ```

## rskip


```lua
(method) Iter:rskip(n: number)
  -> Iter
```

 Discards `n` values from the end of a |list-iterator| pipeline.

 Example:

 ```lua
 local it = vim.iter({ 1, 2, 3, 4, 5 }):rskip(2)
 it:next()
 -- 1
 it:pop()
 -- 3
 ```

@*param* `n` — Number of values to skip.

 luacheck: no unused args

## skip


```lua
(method) Iter:skip(n: number)
  -> Iter
```

 Skips `n` values of an iterator pipeline.

 Example:

 ```lua

 local it = vim.iter({ 3, 6, 9, 12 }):skip(2)
 it:next()
 -- 9

 ```

@*param* `n` — Number of values to skip.

## slice


```lua
(method) Iter:slice(first: number, last: number)
  -> Iter
```

 Sets the start and end of a |list-iterator| pipeline.

 Equivalent to `:skip(first - 1):rskip(len - last + 1)`.

 luacheck: no unused args

## take


```lua
(method) Iter:take(n: integer)
  -> Iter
```

 Transforms an iterator to yield only the first n values.

 Example:

 ```lua
 local it = vim.iter({ 1, 2, 3, 4 }):take(2)
 it:next()
 -- 1
 it:next()
 -- 2
 it:next()
 -- nil
 ```

## totable


```lua
(method) Iter:totable()
  -> table
```

 Collect the iterator into a table.

 The resulting table depends on the initial source in the iterator pipeline.
 Array-like tables and function iterators will be collected into an array-like
 table. If multiple values are returned from the final stage in the iterator
 pipeline, each value will be included in a table.

 Examples:

 ```lua
 vim.iter(string.gmatch('100 20 50', '%d+')):map(tonumber):totable()
 -- { 100, 20, 50 }

 vim.iter({ 1, 2, 3 }):map(function(v) return v, 2 * v end):totable()
 -- { { 1, 2 }, { 2, 4 }, { 3, 6 } }

 vim.iter({ a = 1, b = 2, c = 3 }):filter(function(k, v) return v % 2 ~= 0 end):totable()
 -- { { 'a', 1 }, { 'c', 3 } }
 ```

 The generated table is an array-like table with consecutive, numeric indices.
 To create a map-like table with arbitrary keys, use |Iter:fold()|.


---

# IterMod


---

# Iterator


---

# LuaLS


---

# Man.Attribute


---

# MarkInfo


---

# ParserEntry

## index


```lua
integer
```

runtime path index (unique)

## name


```lua
string
```

## path


```lua
string
```


---

# ParserThreadState


---

# QueryLinterLanguageContext

 Contains language dependent context for the query linter

## is_first_lang


```lua
boolean
```

Whether this is the first language of a linter run checking queries for multiple `langs`

## lang


```lua
string?
```

Current `lang` of the targeted parser

## parser_info


```lua
table?
```

Parser info returned by vim.treesitter.language.inspect


---

# QueryLinterNormalizedOpts

## clear


```lua
boolean
```

## langs


```lua
string[]
```


---

# Range


---

# Range2

## [1]


```lua
integer
```

start row

## [2]


```lua
integer
```

end row


---

# Range4

## [1]


```lua
integer
```

start row

## [2]


```lua
integer
```

start column

## [3]


```lua
integer
```

end row

## [4]


```lua
integer
```

end column


---

# Range6

## [1]


```lua
integer
```

start row

## [2]


```lua
integer
```

start column

## [3]


```lua
integer
```

start bytes

## [4]


```lua
integer
```

end row

## [5]


```lua
integer
```

end column

## [6]


```lua
integer
```

end bytes


---

# STActiveRequest


## request_id


```lua
integer?
```

the LSP request ID of the most recent request sent to the server

## version


```lua
integer?
```

the document version associated with the most recent request


---

# STClientState


## active_request


```lua
STActiveRequest
```


## current_result


```lua
STCurrentResult
```


## namespace


```lua
integer
```


---

# STCurrentResult


## highlights


```lua
STTokenRange[]?
```

cache of highlight ranges for this document version

## namespace_cleared


```lua
boolean?
```

whether the namespace was cleared for this result yet

## result_id


```lua
string?
```

resultId from the server; used with delta requests

## tokens


```lua
integer[]?
```

raw token array as received by the server. used for calculating delta responses

## version


```lua
integer?
```

document version associated with this result


---

# STHighlighter

## active


```lua
table<integer, STHighlighter>
```

## attach


```lua
(method) STHighlighter:attach(client_id: any)
```

## augroup


```lua
integer
```

augroup for buffer events

## bufnr


```lua
integer
```

## client_state


```lua
table<integer, STClientState>
```

## debounce


```lua
integer
```

milliseconds to debounce requests for new tokens

## destroy


```lua
(method) STHighlighter:destroy()
```

## detach


```lua
(method) STHighlighter:detach(client_id: any)
```

## mark_dirty


```lua
(method) STHighlighter:mark_dirty(client_id: integer)
```

 Mark a client's results as dirty. This method will cancel any active
 requests to the server and pause new highlights from being added
 in the on_win callback. The rest of the current results are saved
 in case the server supports delta requests.

## new


```lua
function STHighlighter.new(bufnr: integer)
  -> STHighlighter
```

 Construct a new STHighlighter for the buffer

## on_change


```lua
(method) STHighlighter:on_change()
```

## on_win


```lua
(method) STHighlighter:on_win(topline: integer, botline: integer)
```

 on_win handler for the decoration provider (see |nvim_set_decoration_provider|)

 If there is a current result for the buffer and the version matches the
 current document version, then the tokens are valid and can be applied. As
 the buffer is drawn, this function will add extmark highlights for every
 token in the range of visible lines. Once a highlight has been added, it
 sticks around until the document changes and there's a new set of matching
 highlight tokens available.

 If this is the first time a buffer is being drawn with a new set of
 highlights for the current document version, the namespace is cleared to
 remove extmarks from the last version. It's done here instead of the response
 handler to avoid the "blink" that occurs due to the timing between the
 response handler and the actual redraw.

## process_response


```lua
(method) STHighlighter:process_response(response: lsp.SemanticTokens|lsp.SemanticTokensDelta, client: any, version: any)
```

 This function will parse the semantic token responses and set up the cache
 (current_result). It also performs document synchronization by checking the
 version of the document associated with the resulting request_id and only
 performing work if the response is not out-of-date.

 Delta edits are applied if necessary, and new highlight ranges are calculated
 and stored in the buffer state.

 Finally, a redraw command is issued to force nvim to redraw the screen to
 pick up changed highlight tokens.

## reset


```lua
(method) STHighlighter:reset()
```

 Reset the buffer's highlighting state and clears the extmark highlights.

## reset_timer


```lua
(method) STHighlighter:reset_timer()
```

## send_request


```lua
(method) STHighlighter:send_request()
```

 This is the entry point for getting all the tokens in a buffer.

 For the given clients (or all attached, if not provided), this sends a request
 to ask for semantic tokens. If the server supports delta requests, that will
 be prioritized if we have a previous requestId and token array.

 This function will skip servers where there is an already an active request in
 flight for the same version. If there is a stale request in flight, that is
 cancelled prior to sending a new one.

 Finally, if the request was successful, the requestId and document version
 are saved to facilitate document synchronization in the response.

## timer


```lua
table
```

uv_timer for debouncing requests for new tokens


---

# STTokenRange

## end_col


```lua
integer
```

end column 0-based

## line


```lua
integer
```

line number 0-based

## marked


```lua
boolean
```

whether this token has had extmarks applied

## modifiers


```lua
table<string, boolean>
```

token modifiers as a set. E.g., { static = true, readonly = true }

## start_col


```lua
integer
```

start column 0-based

## type


```lua
string
```

token type as string


---

# STTokenRangeInspect

 @nodoc

## client_id


```lua
integer
```

## end_col


```lua
integer
```

end column 0-based

## line


```lua
integer
```

line number 0-based

## marked


```lua
boolean
```

whether this token has had extmarks applied

## modifiers


```lua
table<string, boolean>
```

token modifiers as a set. E.g., { static = true, readonly = true }

## start_col


```lua
integer
```

start column 0-based

## type


```lua
string
```

token type as string


---

# TS.FoldInfo

Treesitter folding is done in two steps:
(1) compute the fold levels with the syntax tree and cache the result (`compute_folds_levels`)
(2) evaluate foldexpr for each window, which reads from the cache (`foldupdate`)

## __index


```lua
TS.FoldInfo
```

Treesitter folding is done in two steps:
(1) compute the fold levels with the syntax tree and cache the result (`compute_folds_levels`)
(2) evaluate foldexpr for each window, which reads from the cache (`foldupdate`)

## add_range


```lua
(method) TS.FoldInfo:add_range(srow: integer, erow: integer)
```

@*param* `erow` — 0-indexed, exclusive

## do_foldupdate


```lua
(method) TS.FoldInfo:do_foldupdate(bufnr: any)
```

## foldupdate


```lua
(method) TS.FoldInfo:foldupdate(bufnr: any, srow: integer, erow: integer)
```

 Update the folds in the windows that contain the buffer and use expr foldmethod (assuming that
 the user doesn't use different foldexpr for the same buffer).

 Nvim usually automatically updates folds when text changes, but it doesn't work here because
 FoldInfo update is scheduled. So we do it manually.

@*param* `erow` — 0-indexed, exclusive

## foldupdate_range


```lua
Range2?
```


The range on which to evaluate foldexpr.
When in insert mode, the evaluation is deferred to InsertLeave.

## levels


```lua
string[]
```

the cached foldexpr result for each line

## levels0


```lua
integer[]
```

the cached raw fold levels

## new


```lua
function TS.FoldInfo.new(bufnr: integer)
  -> TS.FoldInfo
```

## on_bytes_range


```lua
Range2?
```


The range edited since the last invocation of the callback scheduled in on_bytes.
Should compute fold levels in this range.

## parser


```lua
(vim.treesitter.LanguageTree)?
```


The treesitter parser associated with this buffer.

## remove_range


```lua
(method) TS.FoldInfo:remove_range(srow: integer, erow: integer)
```

@*param* `erow` — 0-indexed, exclusive


---

# TS.Heading

## bufnr


```lua
integer
```

## level


```lua
integer
```

## lnum


```lua
integer
```

## text


```lua
string
```


---

# TSCallbackName


---

# TSCallbackNameOn


---

# TSDirective


---

# TSLangInfo


## _wasm


```lua
boolean
```

## abi_version


```lua
integer
```

## fields


```lua
string[]
```

## metadata


```lua
TSLangMetadata?
```

ABI 15 only

## state_count


```lua
integer
```

## supertypes


```lua
table<string, string[]>
```

## symbols


```lua
table<string, boolean>
```


---

# TSLangMetadata


## major_version


```lua
integer
```

## minor_version


```lua
integer
```

## patch_version


```lua
integer
```


---

# TSLoggerCallback


---

# TSNode

## __has_ancestor


```lua
(method) TSNode:__has_ancestor(node_types: string[])
  -> boolean
```

 Check if the node has any of the given node types as its ancestor.

## byte_length


```lua
(method) TSNode:byte_length()
  -> integer
```

 Return the number of bytes spanned by this node.

## child


```lua
(method) TSNode:child(index: integer)
  -> TSNode?
```

 Get the node's child at the given {index}, where zero represents the first
 child.

## child_count


```lua
(method) TSNode:child_count()
  -> integer
```

 Get the node's number of children.

## child_with_descendant


```lua
(method) TSNode:child_with_descendant(descendant: TSNode)
  -> TSNode?
```

 Get the node's child that contains {descendant} (includes {descendant}).

 For example, with the following node hierarchy:

 ```
 a -> b -> c

 a:child_with_descendant(c) == b
 a:child_with_descendant(b) == b
 a:child_with_descendant(a) == nil
 ```

## descendant_for_range


```lua
(method) TSNode:descendant_for_range(start_row: integer, start_col: integer, end_row: integer, end_col: integer)
  -> TSNode?
```

 Get the smallest node within this node that spans the given range of (row,
 column) positions

## end_


```lua
(method) TSNode:end_()
  -> integer
  2. integer
  3. integer
```

 Get the node's end position. Return three values: the row, column and
 total byte count (all zero-based).

## equal


```lua
(method) TSNode:equal(node: TSNode)
  -> boolean
```

 Check if {node} refers to the same node within the same tree.

## extra


```lua
(method) TSNode:extra()
  -> boolean
```

 Check if the node is extra. Extra nodes represent things like comments,
 which are not required by the grammar but can appear anywhere.

## field


```lua
(method) TSNode:field(name: string)
  -> TSNode[]
```

 Returns a list of all the node's children that have the given field name.

## has_changes


```lua
(method) TSNode:has_changes()
  -> boolean
```

 Check if a syntax node has been edited.

## has_error


```lua
(method) TSNode:has_error()
  -> boolean
```

 Check if the node is a syntax error or contains any syntax errors.

## id


```lua
(method) TSNode:id()
  -> string
```

 Get a unique identifier for the node inside its own tree.

 No guarantees are made about this identifier's internal representation,
 except for being a primitive Lua type with value equality (so not a
 table). Presently it is a (non-printable) string.

 Note: The `id` is not guaranteed to be unique for nodes from different
 trees.

## iter_children


```lua
(method) TSNode:iter_children()
  -> fun():TSNode, string
```

 Iterates over all the direct children of {TSNode}, regardless of whether
 they are named or not.
 Returns the child node plus the eventual field name corresponding to this
 child node.

## missing


```lua
(method) TSNode:missing()
  -> boolean
```

 Check if the node is missing. Missing nodes are inserted by the parser in
 order to recover from certain kinds of syntax errors.

## named


```lua
(method) TSNode:named()
  -> boolean
```

 Check if the node is named. Named nodes correspond to named rules in the
 grammar, whereas anonymous nodes correspond to string literals in the
 grammar.

## named_child


```lua
(method) TSNode:named_child(index: integer)
  -> TSNode?
```

 Get the node's named child at the given {index}, where zero represents the
 first named child.

## named_child_count


```lua
(method) TSNode:named_child_count()
  -> integer
```

 Get the node's number of named children.

## named_children


```lua
(method) TSNode:named_children()
  -> TSNode[]
```

 Returns a list of the node's named children.

## named_descendant_for_range


```lua
(method) TSNode:named_descendant_for_range(start_row: integer, start_col: integer, end_row: integer, end_col: integer)
  -> TSNode?
```

 Get the smallest named node within this node that spans the given range of
 (row, column) positions

## next_named_sibling


```lua
(method) TSNode:next_named_sibling()
  -> TSNode?
```

 Get the node's next named sibling.

## next_sibling


```lua
(method) TSNode:next_sibling()
  -> TSNode?
```

 Get the node's next sibling.

## parent


```lua
(method) TSNode:parent()
  -> TSNode?
```

 Get the node's immediate parent.
 Prefer |TSNode:child_with_descendant()|
 for iterating over the node's ancestors.

## prev_named_sibling


```lua
(method) TSNode:prev_named_sibling()
  -> TSNode?
```

 Get the node's previous named sibling.

## prev_sibling


```lua
(method) TSNode:prev_sibling()
  -> TSNode?
```

 Get the node's previous sibling.

## range


```lua
(method) TSNode:range(include_bytes?: false)
  -> integer
  2. integer
  3. integer
  4. integer
```

 Get the range of the node.

 Return four or six values:

 - start row
 - start column
 - start byte (if {include_bytes} is `true`)
 - end row
 - end column
 - end byte (if {include_bytes} is `true`)

---

```lua
include_bytes:
    | false
```

## sexpr


```lua
(method) TSNode:sexpr()
  -> string
```

 Get an S-expression representing the node as a string.

## start


```lua
(method) TSNode:start()
  -> integer
  2. integer
  3. integer
```

 Get the node's start position. Return three values: the row, column and
 total byte count (all zero-based).

## symbol


```lua
(method) TSNode:symbol()
  -> integer
```

 Get the node's type as a numerical id.

## tree


```lua
(method) TSNode:tree()
  -> TSTree
```

 Get the |TSTree| of the node.

## type


```lua
(method) TSNode:type()
  -> string
```

 Get the node's type as a string.


---

# TSParser

## _logger


```lua
fun(self: TSParser):fun(logtype: 'lex'|'parse', msg: string)
```

## _set_logger


```lua
fun(self: TSParser, lex: boolean, parse: boolean, cb: fun(logtype: 'lex'|'parse', msg: string))
```

## included_ranges


```lua
fun(self: TSParser, include_bytes?: boolean):integer[]
```

## parse


```lua
fun(self: TSParser, tree?: TSTree, source: string|integer, include_bytes: boolean):TSTree, (Range4|Range6)[]
```

## reset


```lua
fun(self: TSParser)
```

## set_included_ranges


```lua
fun(self: TSParser, ranges: (Range6|TSNode)[])
```

## set_timeout


```lua
fun(self: TSParser, timeout: integer)
```

## timeout


```lua
fun(self: TSParser):integer
```


---

# TSPredicate


---

# TSQuery

 Reference to an object held by the treesitter library that is used as a
 component of the |vim.treesitter.Query| for language feature support.
 See |treesitter-query| for more about queries or |vim.treesitter.query.parse()|
 for an example of how to obtain a query object.


## disable_capture


```lua
(method) TSQuery:disable_capture(capture_name: string)
```

 Disable a specific capture in this query; once disabled the capture cannot be re-enabled.
 {capture_name} should not include a leading "@".

 Example: To disable the `@variable.parameter` capture from the vimdoc highlights query:
 ```lua
 local query = vim.treesitter.query.get('vimdoc', 'highlights')
 query.query:disable_capture("variable.parameter")
 vim.treesitter.get_parser():parse()
 ```

## disable_pattern


```lua
(method) TSQuery:disable_pattern(pattern_index: integer)
```

 Disable a specific pattern in this query; once disabled the pattern cannot be re-enabled.
 The {pattern_index} for a particular match can be obtained with |:Inspect!|, or by reading
 the source of the query (i.e. from |vim.treesitter.query.get_files()|).

 Example: To disable `|` links in vimdoc but keep other `@markup.link`s highlighted:
 ```lua
 local link_pattern = 9 -- from :Inspect!
 local query = vim.treesitter.query.get('vimdoc', 'highlights')
 query.query:disable_pattern(link_pattern)
 local tree = vim.treesitter.get_parser():parse()[1]
 ```

## inspect


```lua
(method) TSQuery:inspect()
  -> TSQueryInfo
```

 Get information about the query's patterns and captures.


---

# TSQueryCursor

## next_capture


```lua
(method) TSQueryCursor:next_capture()
  -> capture: integer
  2. captured_node: TSNode
  3. match: TSQueryMatch
```

## next_match


```lua
(method) TSQueryCursor:next_match()
  -> match: TSQueryMatch
```

## remove_match


```lua
fun(self: TSQueryCursor, id: integer)
```


---

# TSQueryInfo

## captures


```lua
string[]
```

## patterns


```lua
table<integer, (string|integer)[][]>
```


---

# TSQueryMatch

## captures


```lua
fun(self: TSQueryMatch):table<integer, TSNode[]>
```

## info


```lua
(method) TSQueryMatch:info()
  -> match_id: integer
  2. pattern_index: integer
```


---

# TSTree

## copy


```lua
(method) TSTree:copy()
  -> TSTree
```

 Returns a copy of the `TSTree`.

## edit


```lua
(method) TSTree:edit(start_byte: integer, end_byte_old: integer, end_byte_new: integer, start_row: integer, start_col: integer, end_row_old: integer, end_col_old: integer, end_row_new: integer, end_col_new: integer)
  -> TSTree
```

 stylua: ignore

## included_ranges


```lua
(method) TSTree:included_ranges(include_bytes: true)
  -> Range6[]
```

```lua
include_bytes:
    | true
```

## root


```lua
(method) TSTree:root()
  -> TSNode
```

 Return the root node of this tree.


---

# TermFeatures

## osc52


```lua
boolean?
```


---

# _


```lua
integer
```


```lua
integer
```


---

# copcall


```lua
function
```


```lua
function _G.copcall(f: any, ...any)
  -> boolean|unknown
  2. unknown
  3. unknown
```


---

# coxpcall


```lua
function
```


```lua
function _G.coxpcall(f: any, err: any, ...any)
  -> boolean|unknown
  2. unknown
  3. unknown
```


---

# decimal


---

# elem_or_list


---

# export.gatherGlobals


```lua
function export.gatherGlobals()
  -> unknown
```


---

# export.getLualsConfig


```lua
function export.getLualsConfig()
  -> unknown
```


---

# export.makeDocObject.INIT


```lua
function export.makeDocObject.INIT(source: any, has_seen: any)
  -> unknown
```


---

# export.makeDocObject.doc.class


```lua
function (source: any, obj: any, has_seen: any)
  -> unknown
```


---

# export.makeDocs


```lua
function export.makeDocs(globals: any, callback: any)
  -> unknown
```


---

# export.serializeAndExport


```lua
function export.serializeAndExport(docs: any, outputDir: any)
  -> unknown
```


---

# lsp.AnnotatedTextEdit

A special text edit with an additional change annotation.


## annotationId


```lua
string
```


The actual identifier of the change annotation

## newText


```lua
string
```


The string to be inserted. For delete operations use an
empty string.

## range


```lua
lsp.Range
```


The range of the text document to be manipulated. To insert
text into a document create a range where start === end.


---

# lsp.ApplyKind

Defines how values from a set of defaults and an individual item will be
merged.



---

# lsp.ApplyWorkspaceEditParams

The parameters passed via an apply workspace edit request.

## edit


```lua
lsp.WorkspaceEdit
```


The edits to apply.

## label


```lua
string?
```


An optional label of the workspace edit. This label is
presented in the user interface for example on an undo
stack to undo the workspace edit.

## metadata


```lua
(lsp.WorkspaceEditMetadata)?
```


Additional data about the edit.



---

# lsp.ApplyWorkspaceEditResult

The result returned from the apply workspace edit request.


## applied


```lua
boolean
```


Indicates whether the edit was applied or not.

## failedChange


```lua
integer?
```


Depending on the client's failure handling strategy `failedChange` might
contain the index of the change that failed. This property is only available
if the client signals a `failureHandlingStrategy` in its client capabilities.

## failureReason


```lua
string?
```


An optional textual description for why the edit was not applied.
This may be used by the server for diagnostic logging or to provide
a suitable error for a request that triggered the edit.


---

# lsp.BaseSymbolInformation

A base for all symbol information.

## containerName


```lua
string?
```


The name of the symbol containing this symbol. This information is for
user interface purposes (e.g. to render a qualifier in the user interface
if necessary). It can't be used to re-infer a hierarchy for the document
symbols.

## kind


```lua
1|10|11|12|13...(+21)
```


The kind of this symbol.

## name


```lua
string
```


The name of this symbol.

## tags


```lua
1[]?
```


Tags for this symbol.



---

# lsp.CallHierarchyClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
return value for the corresponding server capability as well.


---

# lsp.CallHierarchyIncomingCall

Represents an incoming call, e.g. a caller of a method or constructor.


## from


```lua
lsp.CallHierarchyItem
```


The item that makes the call.

## fromRanges


```lua
lsp.Range[]
```


The ranges at which the calls appear. This is relative to the caller
denoted by {@link CallHierarchyIncomingCall.from `this.from`}.


---

# lsp.CallHierarchyIncomingCallsParams

The parameter of a `callHierarchy/incomingCalls` request.


## item


```lua
lsp.CallHierarchyItem
```


## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.CallHierarchyItem

Represents programming constructs like functions or constructors in the context
of call hierarchy.


## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved between a call hierarchy prepare and
incoming calls or outgoing calls requests.

## detail


```lua
string?
```


More detail for this item, e.g. the signature of a function.

## kind


```lua
1|10|11|12|13...(+21)
```


The kind of this item.

## name


```lua
string
```


The name of this item.

## range


```lua
lsp.Range
```


The range enclosing this symbol not including leading/trailing whitespace but everything else, e.g. comments and code.

## selectionRange


```lua
lsp.Range
```


The range that should be selected and revealed when this symbol is being picked, e.g. the name of a function.
Must be contained by the {@link CallHierarchyItem.range `range`}.

## tags


```lua
1[]?
```


Tags for this item.

## uri


```lua
string
```


The resource identifier of this item.


---

# lsp.CallHierarchyOptions

Call hierarchy options used during static registration.


## workDoneProgress


```lua
boolean?
```



---

# lsp.CallHierarchyOutgoingCall

Represents an outgoing call, e.g. calling a getter from a method or a method from a constructor etc.


## fromRanges


```lua
lsp.Range[]
```


The range at which this item is called. This is the range relative to the caller, e.g the item
passed to {@link CallHierarchyItemProvider.provideCallHierarchyOutgoingCalls `provideCallHierarchyOutgoingCalls`}
and not {@link CallHierarchyOutgoingCall.to `this.to`}.

## to


```lua
lsp.CallHierarchyItem
```


The item that is called.


---

# lsp.CallHierarchyOutgoingCallsParams

The parameter of a `callHierarchy/outgoingCalls` request.


## item


```lua
lsp.CallHierarchyItem
```


## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.CallHierarchyPrepareParams

The parameter of a `textDocument/prepareCallHierarchy` request.


## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.CallHierarchyRegistrationOptions

Call hierarchy options used during static or dynamic registration.


## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.CancelParams

## id


```lua
string|integer
```


The request id to cancel.


---

# lsp.ChangeAnnotation

Additional information that describes document changes.


## description


```lua
string?
```


A human-readable string which is rendered less prominent in
the user interface.

## label


```lua
string
```


A human-readable string describing the actual change. The string
is rendered prominent in the user interface.

## needsConfirmation


```lua
boolean?
```


A flag which indicates that user confirmation is needed
before applying the change.


---

# lsp.ChangeAnnotationIdentifier

An identifier to refer to a change annotation stored with a workspace edit.


---

# lsp.ChangeAnnotationsSupportOptions

## groupsOnLabel


```lua
boolean?
```


Whether the client groups edits with equal labels into tree nodes,
for instance all edits labelled with "Changes in Strings" would
be a tree node.


---

# lsp.ClientCapabilities

Defines the capabilities provided by the client.

## experimental


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


Experimental client capabilities.

## general


```lua
(lsp.GeneralClientCapabilities)?
```


General client capabilities.


## notebookDocument


```lua
(lsp.NotebookDocumentClientCapabilities)?
```


Capabilities specific to the notebook document support.


## textDocument


```lua
(lsp.TextDocumentClientCapabilities)?
```


Text document specific client capabilities.

## window


```lua
(lsp.WindowClientCapabilities)?
```


Window specific client capabilities.

## workspace


```lua
(lsp.WorkspaceClientCapabilities)?
```


Workspace specific client capabilities.


---

# lsp.ClientCodeActionKindOptions

## valueSet


```lua
""|"notebook"|"quickfix"|"refactor"|"refactor.extract"...(+6)[]
```


The code action kind values the client supports. When this
property exists the client also guarantees that it will
handle values outside its set gracefully and falls back
to a default value when unknown.


---

# lsp.ClientCodeActionLiteralOptions

## codeActionKind


```lua
lsp.ClientCodeActionKindOptions
```


The code action kind is support with the following value
set.


---

# lsp.ClientCodeActionResolveOptions

## properties


```lua
string[]
```


The properties that a client can resolve lazily.


---

# lsp.ClientCodeLensResolveOptions

## properties


```lua
string[]
```


The properties that a client can resolve lazily.


---

# lsp.ClientCompletionItemInsertTextModeOptions

## valueSet


```lua
1|2[]
```



---

# lsp.ClientCompletionItemOptions

## commitCharactersSupport


```lua
boolean?
```


Client supports commit characters on a completion item.

## deprecatedSupport


```lua
boolean?
```


Client supports the deprecated property on a completion item.

## documentationFormat


```lua
"markdown"|"plaintext"[]?
```


Client supports the following content formats for the documentation
property. The order describes the preferred format of the client.

## insertReplaceSupport


```lua
boolean?
```


Client support insert replace edit to control different behavior if a
completion item is inserted in the text or should replace text.


## insertTextModeSupport


```lua
(lsp.ClientCompletionItemInsertTextModeOptions)?
```


The client supports the `insertTextMode` property on
a completion item to override the whitespace handling mode
as defined by the client (see `insertTextMode`).


## labelDetailsSupport


```lua
boolean?
```


The client has support for completion item label
details (see also `CompletionItemLabelDetails`).


## preselectSupport


```lua
boolean?
```


Client supports the preselect property on a completion item.

## resolveSupport


```lua
(lsp.ClientCompletionItemResolveOptions)?
```


Indicates which properties a client can resolve lazily on a completion
item. Before version 3.16.0 only the predefined properties `documentation`
and `details` could be resolved lazily.


## snippetSupport


```lua
boolean?
```


Client supports snippets as insert text.

A snippet can define tab stops and placeholders with `$1`, `$2`
and `${3:foo}`. `$0` defines the final tab stop, it defaults to
the end of the snippet. Placeholders with equal identifiers are linked,
that is typing in one will update others too.

## tagSupport


```lua
(lsp.CompletionItemTagOptions)?
```


Client supports the tag property on a completion item. Clients supporting
tags have to handle unknown tags gracefully. Clients especially need to
preserve unknown tags when sending a completion item back to the server in
a resolve call.



---

# lsp.ClientCompletionItemOptionsKind

## valueSet


```lua
1|10|11|12|13...(+20)[]?
```


The completion item kind values the client supports. When this
property exists the client also guarantees that it will
handle values outside its set gracefully and falls back
to a default value when unknown.

If this property is not present the client only supports
the completion items kinds from `Text` to `Reference` as defined in
the initial version of the protocol.


---

# lsp.ClientCompletionItemResolveOptions

## properties


```lua
string[]
```


The properties that a client can resolve lazily.


---

# lsp.ClientDiagnosticsTagOptions

## valueSet


```lua
1|2[]
```


The tags supported by the client.


---

# lsp.ClientFoldingRangeKindOptions

## valueSet


```lua
"comment"|"imports"|"region"[]?
```


The folding range kind values the client supports. When this
property exists the client also guarantees that it will
handle values outside its set gracefully and falls back
to a default value when unknown.


---

# lsp.ClientFoldingRangeOptions

## collapsedText


```lua
boolean?
```


If set, the client signals that it supports setting collapsedText on
folding ranges to display custom labels instead of the default text.



---

# lsp.ClientInfo

Information about the client


## name


```lua
string
```


The name of the client as defined by the client.

## version


```lua
string?
```


The client's version as defined by the client.


---

# lsp.ClientInlayHintResolveOptions

## properties


```lua
string[]
```


The properties that a client can resolve lazily.


---

# lsp.ClientSemanticTokensRequestFullDelta

## delta


```lua
boolean?
```


The client will send the `textDocument/semanticTokens/full/delta` request if
the server provides a corresponding handler.


---

# lsp.ClientSemanticTokensRequestOptions

## full


```lua
(boolean|lsp.ClientSemanticTokensRequestFullDelta)?
```


The client will send the `textDocument/semanticTokens/full` request if
the server provides a corresponding handler.

## range


```lua
(boolean|lsp._anonym2.range)?
```


The client will send the `textDocument/semanticTokens/range` request if
the server provides a corresponding handler.


---

# lsp.ClientShowMessageActionItemOptions

## additionalPropertiesSupport


```lua
boolean?
```


Whether the client supports additional attributes which
are preserved and send back to the server in the
request's response.


---

# lsp.ClientSignatureInformationOptions

## activeParameterSupport


```lua
boolean?
```


The client supports the `activeParameter` property on `SignatureInformation`
literal.


## documentationFormat


```lua
"markdown"|"plaintext"[]?
```


Client supports the following content formats for the documentation
property. The order describes the preferred format of the client.

## noActiveParameterSupport


```lua
boolean?
```


The client supports the `activeParameter` property on
`SignatureHelp`/`SignatureInformation` being set to `null` to
indicate that no parameter should be active.


## parameterInformation


```lua
(lsp.ClientSignatureParameterInformationOptions)?
```


Client capabilities specific to parameter information.


---

# lsp.ClientSignatureParameterInformationOptions

## labelOffsetSupport


```lua
boolean?
```


The client supports processing label offsets instead of a
simple label string.



---

# lsp.ClientSymbolKindOptions

## valueSet


```lua
1|10|11|12|13...(+21)[]?
```


The symbol kind values the client supports. When this
property exists the client also guarantees that it will
handle values outside its set gracefully and falls back
to a default value when unknown.

If this property is not present the client only supports
the symbol kinds from `File` to `Array` as defined in
the initial version of the protocol.


---

# lsp.ClientSymbolResolveOptions

## properties


```lua
string[]
```


The properties that a client can resolve lazily. Usually
`location.range`


---

# lsp.ClientSymbolTagOptions

## valueSet


```lua
1[]
```


The tags supported by the client.


---

# lsp.CodeAction

A code action represents a change that can be performed in code, e.g. to fix a problem or
to refactor code.

A CodeAction must set either `edit` and/or a `command`. If both are supplied, the `edit` is applied first, then the `command` is executed.

## command


```lua
(lsp.Command)?
```


A command this code action executes. If a code action
provides an edit and a command, first the edit is
executed and then the command.

## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved on a code action between
a `textDocument/codeAction` and a `codeAction/resolve` request.


## diagnostics


```lua
lsp.Diagnostic[]?
```


The diagnostics that this code action resolves.

## disabled


```lua
(lsp.CodeActionDisabled)?
```


Marks that the code action cannot currently be applied.

Clients should follow the following guidelines regarding disabled code actions:

  - Disabled code actions are not shown in automatic [lightbulbs](https://code.visualstudio.com/docs/editor/editingevolved#_code-action)
    code action menus.

  - Disabled actions are shown as faded out in the code action menu when the user requests a more specific type
    of code action, such as refactorings.

  - If the user has a [keybinding](https://code.visualstudio.com/docs/editor/refactoring#_keybindings-for-code-actions)
    that auto applies a code action and only disabled code actions are returned, the client should show the user an
    error message with `reason` in the editor.


## edit


```lua
(lsp.WorkspaceEdit)?
```


The workspace edit this code action performs.

## isPreferred


```lua
boolean?
```


Marks this as a preferred action. Preferred actions are used by the `auto fix` command and can be targeted
by keybindings.

A quick fix should be marked preferred if it properly addresses the underlying error.
A refactoring should be marked preferred if it is the most reasonable choice of actions to take.


## kind


```lua
(""|"notebook"|"quickfix"|"refactor"|"refactor.extract"...(+6))?
```


The kind of the code action.

Used to filter code actions.

## tags


```lua
1[]?
```


Tags for this code action.


## title


```lua
string
```


A short, human-readable, title for this code action.


---

# lsp.CodeActionClientCapabilities

The Client Capabilities of a {@link CodeActionRequest}.

## codeActionLiteralSupport


```lua
(lsp.ClientCodeActionLiteralOptions)?
```


The client support code action literals of type `CodeAction` as a valid
response of the `textDocument/codeAction` request. If the property is not
set the request can only return `Command` literals.


## dataSupport


```lua
boolean?
```


Whether code action supports the `data` property which is
preserved between a `textDocument/codeAction` and a
`codeAction/resolve` request.


## disabledSupport


```lua
boolean?
```


Whether code action supports the `disabled` property.


## documentationSupport


```lua
boolean?
```


Whether the client supports documentation for a class of
code actions.


## dynamicRegistration


```lua
boolean?
```


Whether code action supports dynamic registration.

## honorsChangeAnnotations


```lua
boolean?
```


Whether the client honors the change annotations in
text edits and resource operations returned via the
`CodeAction#edit` property by for example presenting
the workspace edit in the user interface and asking
for confirmation.


## isPreferredSupport


```lua
boolean?
```


Whether code action supports the `isPreferred` property.


## resolveSupport


```lua
(lsp.ClientCodeActionResolveOptions)?
```


Whether the client supports resolving additional code action
properties via a separate `codeAction/resolve` request.


## tagSupport


```lua
(lsp.CodeActionTagOptions)?
```


Client supports the tag property on a code action. Clients
supporting tags have to handle unknown tags gracefully.



---

# lsp.CodeActionContext

Contains additional diagnostic information about the context in which
a {@link CodeActionProvider.provideCodeActions code action} is run.

## diagnostics


```lua
lsp.Diagnostic[]
```


An array of diagnostics known on the client side overlapping the range provided to the
`textDocument/codeAction` request. They are provided so that the server knows which
errors are currently presented to the user for the given range. There is no guarantee
that these accurately reflect the error state of the resource. The primary parameter
to compute code actions is the provided range.

## only


```lua
""|"notebook"|"quickfix"|"refactor"|"refactor.extract"...(+6)[]?
```


Requested kind of actions to return.

Actions not of this kind are filtered out by the client before being shown. So servers
can omit computing them.

## triggerKind


```lua
(1|2)?
```


The reason why code actions were requested.



---

# lsp.CodeActionDisabled

Captures why the code action is currently disabled.


## reason


```lua
string
```


Human readable description of why the code action is currently disabled.

This is displayed in the code actions UI.


---

# lsp.CodeActionKind

A set of predefined code action kinds


---

# lsp.CodeActionKindDocumentation

Documentation for a class of code actions.


## command


```lua
lsp.Command
```


Command that is ued to display the documentation to the user.

The title of this documentation code action is taken from {@linkcode Command.title}

## kind


```lua
""|"notebook"|"quickfix"|"refactor"|"refactor.extract"...(+6)
```


The kind of the code action being documented.

If the kind is generic, such as `CodeActionKind.Refactor`, the documentation will be shown whenever any
refactorings are returned. If the kind if more specific, such as `CodeActionKind.RefactorExtract`, the
documentation will only be shown when extract refactoring code actions are returned.


---

# lsp.CodeActionOptions

Provider options for a {@link CodeActionRequest}.

## codeActionKinds


```lua
""|"notebook"|"quickfix"|"refactor"|"refactor.extract"...(+6)[]?
```


CodeActionKinds that this server may return.

The list of kinds may be generic, such as `CodeActionKind.Refactor`, or the server
may list out every specific kind they provide.

## documentation


```lua
lsp.CodeActionKindDocumentation[]?
```


Static documentation for a class of code actions.

Documentation from the provider should be shown in the code actions menu if either:

- Code actions of `kind` are requested by the editor. In this case, the editor will show the documentation that
  most closely matches the requested code action kind. For example, if a provider has documentation for
  both `Refactor` and `RefactorExtract`, when the user requests code actions for `RefactorExtract`,
  the editor will use the documentation for `RefactorExtract` instead of the documentation for `Refactor`.

- Any code actions of `kind` are returned by the provider.

At most one documentation entry should be shown per provider.


## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for a code action.


## workDoneProgress


```lua
boolean?
```



---

# lsp.CodeActionParams

The parameters of a {@link CodeActionRequest}.

## context


```lua
lsp.CodeActionContext
```


Context carrying additional information.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## range


```lua
lsp.Range
```


The range for which the command was invoked.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document in which the command was invoked.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.CodeActionRegistrationOptions

Registration options for a {@link CodeActionRequest}.

## codeActionKinds


```lua
""|"notebook"|"quickfix"|"refactor"|"refactor.extract"...(+6)[]?
```


CodeActionKinds that this server may return.

The list of kinds may be generic, such as `CodeActionKind.Refactor`, or the server
may list out every specific kind they provide.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## documentation


```lua
lsp.CodeActionKindDocumentation[]?
```


Static documentation for a class of code actions.

Documentation from the provider should be shown in the code actions menu if either:

- Code actions of `kind` are requested by the editor. In this case, the editor will show the documentation that
  most closely matches the requested code action kind. For example, if a provider has documentation for
  both `Refactor` and `RefactorExtract`, when the user requests code actions for `RefactorExtract`,
  the editor will use the documentation for `RefactorExtract` instead of the documentation for `Refactor`.

- Any code actions of `kind` are returned by the provider.

At most one documentation entry should be shown per provider.


## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for a code action.


## workDoneProgress


```lua
boolean?
```



---

# lsp.CodeActionTag

Code action tags are extra annotations that tweak the behavior of a code action.



---

# lsp.CodeActionTagOptions

## valueSet


```lua
1[]
```


The tags supported by the client.


---

# lsp.CodeActionTriggerKind

The reason why code actions were requested.



---

# lsp.CodeDescription

Structure to capture a description for an error code.


## href


```lua
string
```


An URI to open with more information about the diagnostic error.


---

# lsp.CodeLens

A code lens represents a {@link Command command} that should be shown along with
source text, like the number of references, a way to run tests, etc.

A code lens is _unresolved_ when no command is associated to it. For performance
reasons the creation of a code lens and resolving should be done in two stages.

## command


```lua
(lsp.Command)?
```


The command this code lens represents.

## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved on a code lens item between
a {@link CodeLensRequest} and a {@link CodeLensResolveRequest}

## range


```lua
lsp.Range
```


The range in which this code lens is valid. Should only span a single line.


---

# lsp.CodeLensClientCapabilities

The client capabilities  of a {@link CodeLensRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether code lens supports dynamic registration.

## resolveSupport


```lua
(lsp.ClientCodeLensResolveOptions)?
```


Whether the client supports resolving additional code lens
properties via a separate `codeLens/resolve` request.



---

# lsp.CodeLensOptions

Code Lens provider options of a {@link CodeLensRequest}.

## resolveProvider


```lua
boolean?
```


Code lens has a resolve provider as well.

## workDoneProgress


```lua
boolean?
```



---

# lsp.CodeLensParams

The parameters of a {@link CodeLensRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document to request code lens for.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.CodeLensRegistrationOptions

Registration options for a {@link CodeLensRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## resolveProvider


```lua
boolean?
```


Code lens has a resolve provider as well.

## workDoneProgress


```lua
boolean?
```



---

# lsp.CodeLensWorkspaceClientCapabilities

## refreshSupport


```lua
boolean?
```


Whether the client implementation supports a refresh request sent from the
server to the client.

Note that this event is global and will force the client to refresh all
code lenses currently shown. It should be used with absolute care and is
useful for situation where a server for example detect a project wide
change that requires such a calculation.


---

# lsp.Color

Represents a color in RGBA space.

## alpha


```lua
number
```


The alpha component of this color in the range [0-1].

## blue


```lua
number
```


The blue component of this color in the range [0-1].

## green


```lua
number
```


The green component of this color in the range [0-1].

## red


```lua
number
```


The red component of this color in the range [0-1].


---

# lsp.ColorInformation

Represents a color range from a document.

## color


```lua
lsp.Color
```


The actual color value for this color range.

## range


```lua
lsp.Range
```


The range in the document where this color appears.


---

# lsp.ColorPresentation

## additionalTextEdits


```lua
lsp.TextEdit[]?
```


An optional array of additional {@link TextEdit text edits} that are applied when
selecting this color presentation. Edits must not overlap with the main {@link ColorPresentation.textEdit edit} nor with themselves.

## label


```lua
string
```


The label of this color presentation. It will be shown on the color
picker header. By default this is also the text that is inserted when selecting
this color presentation.

## textEdit


```lua
(lsp.TextEdit)?
```


An {@link TextEdit edit} which is applied to a document when selecting
this presentation for the color.  When `falsy` the {@link ColorPresentation.label label}
is used.


---

# lsp.ColorPresentationParams

Parameters for a {@link ColorPresentationRequest}.

## color


```lua
lsp.Color
```


The color to request presentations for.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## range


```lua
lsp.Range
```


The range where the color would be inserted. Serves as a context.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.Command

Represents a reference to a command. Provides a title which
will be used to represent a command in the UI and, optionally,
an array of arguments which will be passed to the command handler
function when invoked.

## arguments


```lua
boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]?
```


Arguments that the command handler should be
invoked with.

## command


```lua
string
```


The identifier of the actual command handler.

## title


```lua
string
```


Title of the command, like `save`.

## tooltip


```lua
string?
```


An optional tooltip.



---

# lsp.CompletionClientCapabilities

Completion client capabilities

## completionItem


```lua
(lsp.ClientCompletionItemOptions)?
```


The client supports the following `CompletionItem` specific
capabilities.

## completionItemKind


```lua
(lsp.ClientCompletionItemOptionsKind)?
```


## completionList


```lua
(lsp.CompletionListCapabilities)?
```


The client supports the following `CompletionList` specific
capabilities.


## contextSupport


```lua
boolean?
```


The client supports to send additional context information for a
`textDocument/completion` request.

## dynamicRegistration


```lua
boolean?
```


Whether completion supports dynamic registration.

## insertTextMode


```lua
(1|2)?
```


Defines how the client handles whitespace and indentation
when accepting a completion item that uses multi line
text in either `insertText` or `textEdit`.



---

# lsp.CompletionContext

Contains additional information about the context in which a completion request is triggered.

## triggerCharacter


```lua
string?
```


The trigger character (a single character) that has trigger code complete.
Is undefined if `triggerKind !== CompletionTriggerKind.TriggerCharacter`

## triggerKind


```lua
1|2|3
```


How the completion was triggered.


---

# lsp.CompletionItem

A completion item represents a text snippet that is
proposed to complete text that is being typed.

## additionalTextEdits


```lua
lsp.TextEdit[]?
```


An optional array of additional {@link TextEdit text edits} that are applied when
selecting this completion. Edits must not overlap (including the same insert position)
with the main {@link CompletionItem.textEdit edit} nor with themselves.

Additional text edits should be used to change text unrelated to the current cursor position
(for example adding an import statement at the top of the file if the completion item will
insert an unqualified type).

## command


```lua
(lsp.Command)?
```


An optional {@link Command command} that is executed *after* inserting this completion. *Note* that
additional modifications to the current document should be described with the
{@link CompletionItem.additionalTextEdits additionalTextEdits}-property.

## commitCharacters


```lua
string[]?
```


An optional set of characters that when pressed while this completion is active will accept it first and
then type that character. *Note* that all commit characters should have `length=1` and that superfluous
characters will be ignored.

## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved on a completion item between a
{@link CompletionRequest} and a {@link CompletionResolveRequest}.

## deprecated


```lua
boolean?
```


Indicates if this item is deprecated.

## detail


```lua
string?
```


A human-readable string with additional information
about this item, like type or symbol information.

## documentation


```lua
(string|lsp.MarkupContent)?
```


A human-readable string that represents a doc-comment.

## filterText


```lua
string?
```


A string that should be used when filtering a set of
completion items. When `falsy` the {@link CompletionItem.label label}
is used.

## insertText


```lua
string?
```


A string that should be inserted into a document when selecting
this completion. When `falsy` the {@link CompletionItem.label label}
is used.

The `insertText` is subject to interpretation by the client side.
Some tools might not take the string literally. For example
VS Code when code complete is requested in this example
`con<cursor position>` and a completion item with an `insertText` of
`console` is provided it will only insert `sole`. Therefore it is
recommended to use `textEdit` instead since it avoids additional client
side interpretation.

## insertTextFormat


```lua
(1|2)?
```


The format of the insert text. The format applies to both the
`insertText` property and the `newText` property of a provided
`textEdit`. If omitted defaults to `InsertTextFormat.PlainText`.

Please note that the insertTextFormat doesn't apply to
`additionalTextEdits`.

## insertTextMode


```lua
(1|2)?
```


How whitespace and indentation is handled during completion
item insertion. If not provided the clients default value depends on
the `textDocument.completion.insertTextMode` client capability.


## kind


```lua
(1|10|11|12|13...(+20))?
```


The kind of this completion item. Based of the kind
an icon is chosen by the editor.

## label


```lua
string
```


The label of this completion item.

The label property is also by default the text that
is inserted when selecting this completion.

If label details are provided the label itself should
be an unqualified name of the completion item.

## labelDetails


```lua
(lsp.CompletionItemLabelDetails)?
```


Additional details for the label


## preselect


```lua
boolean?
```


Select this item when showing.

*Note* that only one completion item can be selected and that the
tool / client decides which item that is. The rule is that the *first*
item of those that match best is selected.

## sortText


```lua
string?
```


A string that should be used when comparing this item
with other items. When `falsy` the {@link CompletionItem.label label}
is used.

## tags


```lua
1[]?
```


Tags for this completion item.


## textEdit


```lua
(lsp.InsertReplaceEdit|lsp.TextEdit)?
```


An {@link TextEdit edit} which is applied to a document when selecting
this completion. When an edit is provided the value of
{@link CompletionItem.insertText insertText} is ignored.

Most editors support two different operations when accepting a completion
item. One is to insert a completion text and the other is to replace an
existing text with a completion text. Since this can usually not be
predetermined by a server it can report both ranges. Clients need to
signal support for `InsertReplaceEdits` via the
`textDocument.completion.insertReplaceSupport` client capability
property.

*Note 1:* The text edit's range as well as both ranges from an insert
replace edit must be a [single line] and they must contain the position
at which completion has been requested.
*Note 2:* If an `InsertReplaceEdit` is returned the edit's insert range
must be a prefix of the edit's replace range, that means it must be
contained and starting at the same position.


## textEditText


```lua
string?
```


The edit text used if the completion item is part of a CompletionList and
CompletionList defines an item default for the text edit range.

Clients will only honor this property if they opt into completion list
item defaults using the capability `completionList.itemDefaults`.

If not provided and a list's default range is provided the label
property is used as a text.



---

# lsp.CompletionItemApplyKinds

Specifies how fields from a completion item should be combined with those
from `completionList.itemDefaults`.

If unspecified, all fields will be treated as ApplyKind.Replace.

If a field's value is ApplyKind.Replace, the value from a completion item (if
provided and not `null`) will always be used instead of the value from
`completionItem.itemDefaults`.

If a field's value is ApplyKind.Merge, the values will be merged using the rules
defined against each field below.

Servers are only allowed to return `applyKind` if the client
signals support for this via the `completionList.applyKindSupport`
capability.


## commitCharacters


```lua
(1|2)?
```


Specifies whether commitCharacters on a completion will replace or be
merged with those in `completionList.itemDefaults.commitCharacters`.

If ApplyKind.Replace, the commit characters from the completion item will
always be used unless not provided, in which case those from
`completionList.itemDefaults.commitCharacters` will be used. An
empty list can be used if a completion item does not have any commit
characters and also should not use those from
`completionList.itemDefaults.commitCharacters`.

If ApplyKind.Merge the commitCharacters for the completion will be the
union of all values in both `completionList.itemDefaults.commitCharacters`
and the completion's own `commitCharacters`.


## data


```lua
(1|2)?
```


Specifies whether the `data` field on a completion will replace or
be merged with data from `completionList.itemDefaults.data`.

If ApplyKind.Replace, the data from the completion item will be used if
provided (and not `null`), otherwise
`completionList.itemDefaults.data` will be used. An empty object can
be used if a completion item does not have any data but also should
not use the value from `completionList.itemDefaults.data`.

If ApplyKind.Merge, a shallow merge will be performed between
`completionList.itemDefaults.data` and the completion's own data
using the following rules:

- If a completion's `data` field is not provided (or `null`), the
  entire `data` field from `completionList.itemDefaults.data` will be
  used as-is.
- If a completion's `data` field is provided, each field will
  overwrite the field of the same name in
  `completionList.itemDefaults.data` but no merging of nested fields
  within that value will occur.



---

# lsp.CompletionItemDefaults

In many cases the items of an actual completion result share the same
value for properties like `commitCharacters` or the range of a text
edit. A completion list can therefore define item defaults which will
be used if a completion item itself doesn't specify the value.

If a completion list specifies a default value and a completion item
also specifies a corresponding value, the rules for combining these are
defined by `applyKinds` (if the client supports it), defaulting to
ApplyKind.Replace.

Servers are only allowed to return default values if the client
signals support for this via the `completionList.itemDefaults`
capability.


## commitCharacters


```lua
string[]?
```


A default commit character set.


## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A default data value.


## editRange


```lua
(lsp.EditRangeWithInsertReplace|lsp.Range)?
```


A default edit range.


## insertTextFormat


```lua
(1|2)?
```


A default insert text format.


## insertTextMode


```lua
(1|2)?
```


A default insert text mode.



---

# lsp.CompletionItemKind

The kind of a completion entry.


---

# lsp.CompletionItemLabelDetails

Additional details for a completion item label.


## description


```lua
string?
```


An optional string which is rendered less prominently after {@link CompletionItem.detail}. Should be used
for fully qualified names and file paths.

## detail


```lua
string?
```


An optional string which is rendered less prominently directly after {@link CompletionItem.label label},
without any spacing. Should be used for function signatures and type annotations.


---

# lsp.CompletionItemTag

Completion item tags are extra annotations that tweak the rendering of a completion
item.



---

# lsp.CompletionItemTagOptions

## valueSet


```lua
1[]
```


The tags supported by the client.


---

# lsp.CompletionList

Represents a collection of {@link CompletionItem completion items} to be presented
in the editor.

## applyKind


```lua
(lsp.CompletionItemApplyKinds)?
```


Specifies how fields from a completion item should be combined with those
from `completionList.itemDefaults`.

If unspecified, all fields will be treated as ApplyKind.Replace.

If a field's value is ApplyKind.Replace, the value from a completion item
(if provided and not `null`) will always be used instead of the value
from `completionItem.itemDefaults`.

If a field's value is ApplyKind.Merge, the values will be merged using
the rules defined against each field below.

Servers are only allowed to return `applyKind` if the client
signals support for this via the `completionList.applyKindSupport`
capability.


## isIncomplete


```lua
boolean
```


This list it not complete. Further typing results in recomputing this list.

Recomputed lists have all their items replaced (not appended) in the
incomplete completion sessions.

## itemDefaults


```lua
(lsp.CompletionItemDefaults)?
```


In many cases the items of an actual completion result share the same
value for properties like `commitCharacters` or the range of a text
edit. A completion list can therefore define item defaults which will
be used if a completion item itself doesn't specify the value.

If a completion list specifies a default value and a completion item
also specifies a corresponding value, the rules for combining these are
defined by `applyKinds` (if the client supports it), defaulting to
ApplyKind.Replace.

Servers are only allowed to return default values if the client
signals support for this via the `completionList.itemDefaults`
capability.


## items


```lua
lsp.CompletionItem[]
```


The completion items.


---

# lsp.CompletionListCapabilities

The client supports the following `CompletionList` specific
capabilities.


## applyKindSupport


```lua
boolean?
```


Specifies whether the client supports `CompletionList.applyKind` to
indicate how supported values from `completionList.itemDefaults`
and `completion` will be combined.

If a client supports `applyKind` it must support it for all fields
that it supports that are listed in `CompletionList.applyKind`. This
means when clients add support for new/future fields in completion
items the MUST also support merge for them if those fields are
defined in `CompletionList.applyKind`.


## itemDefaults


```lua
string[]?
```


The client supports the following itemDefaults on
a completion list.

The value lists the supported property names of the
`CompletionList.itemDefaults` object. If omitted
no properties are supported.



---

# lsp.CompletionOptions

Completion options.

## allCommitCharacters


```lua
string[]?
```


The list of all possible characters that commit a completion. This field can be used
if clients don't support individual commit characters per completion item. See
`ClientCapabilities.textDocument.completion.completionItem.commitCharactersSupport`

If a server provides both `allCommitCharacters` and commit characters on an individual
completion item the ones on the completion item win.


## completionItem


```lua
(lsp.ServerCompletionItemOptions)?
```


The server supports the following `CompletionItem` specific
capabilities.


## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for a completion item.

## triggerCharacters


```lua
string[]?
```


Most tools trigger completion request automatically without explicitly requesting
it using a keyboard shortcut (e.g. Ctrl+Space). Typically they do so when the user
starts to type an identifier. For example if the user types `c` in a JavaScript file
code complete will automatically pop up present `console` besides others as a
completion item. Characters that make up identifiers don't need to be listed here.

If code complete should automatically be trigger on characters not being valid inside
an identifier (for example `.` in JavaScript) list them in `triggerCharacters`.

## workDoneProgress


```lua
boolean?
```



---

# lsp.CompletionParams

Completion parameters

## context


```lua
(lsp.CompletionContext)?
```


The completion context. This is only available it the client specifies
to send this using the client capability `textDocument.completion.contextSupport === true`

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.CompletionRegistrationOptions

Registration options for a {@link CompletionRequest}.

## allCommitCharacters


```lua
string[]?
```


The list of all possible characters that commit a completion. This field can be used
if clients don't support individual commit characters per completion item. See
`ClientCapabilities.textDocument.completion.completionItem.commitCharactersSupport`

If a server provides both `allCommitCharacters` and commit characters on an individual
completion item the ones on the completion item win.


## completionItem


```lua
(lsp.ServerCompletionItemOptions)?
```


The server supports the following `CompletionItem` specific
capabilities.


## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for a completion item.

## triggerCharacters


```lua
string[]?
```


Most tools trigger completion request automatically without explicitly requesting
it using a keyboard shortcut (e.g. Ctrl+Space). Typically they do so when the user
starts to type an identifier. For example if the user types `c` in a JavaScript file
code complete will automatically pop up present `console` besides others as a
completion item. Characters that make up identifiers don't need to be listed here.

If code complete should automatically be trigger on characters not being valid inside
an identifier (for example `.` in JavaScript) list them in `triggerCharacters`.

## workDoneProgress


```lua
boolean?
```



---

# lsp.CompletionTriggerKind

How a completion was triggered


---

# lsp.ConfigurationItem

## scopeUri


```lua
string?
```


The scope to get the configuration section for.

## section


```lua
string?
```


The configuration section asked for.


---

# lsp.ConfigurationParams

The parameters of a configuration request.

## items


```lua
lsp.ConfigurationItem[]
```



---

# lsp.CreateFile

Create file operation.

## annotationId


```lua
string?
```


An optional annotation identifier describing the operation.


## kind


```lua
"create"
```


A create

## options


```lua
(lsp.CreateFileOptions)?
```


Additional options

## uri


```lua
string
```


The resource to create.


---

# lsp.CreateFileOptions

Options to create a file.

## ignoreIfExists


```lua
boolean?
```


Ignore if exists.

## overwrite


```lua
boolean?
```


Overwrite existing file. Overwrite wins over `ignoreIfExists`


---

# lsp.CreateFilesParams

The parameters sent in notifications/requests for user-initiated creation of
files.


## files


```lua
lsp.FileCreate[]
```


An array of all files/folders created in this operation.


---

# lsp.Declaration

The declaration of a symbol representation as one or many {@link Location locations}.


---

# lsp.DeclarationClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether declaration supports dynamic registration. If this is set to `true`
the client supports the new `DeclarationRegistrationOptions` return value
for the corresponding server capability as well.

## linkSupport


```lua
boolean?
```


The client supports additional metadata in the form of declaration links.


---

# lsp.DeclarationLink

Information about where a symbol is declared.

Provides additional metadata over normal {@link Location location} declarations, including the range of
the declaring symbol.

Servers should prefer returning `DeclarationLink` over `Declaration` if supported
by the client.


---

# lsp.DeclarationOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.DeclarationParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DeclarationRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.Definition

The definition of a symbol represented as one or many {@link Location locations}.
For most programming languages there is only one location at which a symbol is
defined.

Servers should prefer returning `DefinitionLink` over `Definition` if supported
by the client.


---

# lsp.DefinitionClientCapabilities

Client Capabilities for a {@link DefinitionRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether definition supports dynamic registration.

## linkSupport


```lua
boolean?
```


The client supports additional metadata in the form of definition links.



---

# lsp.DefinitionLink

Information about where a symbol is defined.

Provides additional metadata over normal {@link Location location} definitions, including the range of
the defining symbol


---

# lsp.DefinitionOptions

Server Capabilities for a {@link DefinitionRequest}.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DefinitionParams

Parameters for a {@link DefinitionRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DefinitionRegistrationOptions

Registration options for a {@link DefinitionRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DeleteFile

Delete file operation

## annotationId


```lua
string?
```


An optional annotation identifier describing the operation.


## kind


```lua
"delete"
```


A delete

## options


```lua
(lsp.DeleteFileOptions)?
```


Delete options.

## uri


```lua
string
```


The file to delete.


---

# lsp.DeleteFileOptions

Delete file options

## ignoreIfNotExists


```lua
boolean?
```


Ignore the operation if the file doesn't exist.

## recursive


```lua
boolean?
```


Delete the content recursively if a folder is denoted.


---

# lsp.DeleteFilesParams

The parameters sent in notifications/requests for user-initiated deletes of
files.


## files


```lua
lsp.FileDelete[]
```


An array of all files/folders deleted in this operation.


---

# lsp.Diagnostic

Represents a diagnostic, such as a compiler error or warning. Diagnostic objects
are only valid in the scope of a resource.

## code


```lua
(string|integer)?
```


The diagnostic's code, which usually appear in the user interface.

## codeDescription


```lua
(lsp.CodeDescription)?
```


An optional property to describe the error code.
Requires the code field (above) to be present/not null.


## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved between a `textDocument/publishDiagnostics`
notification and `textDocument/codeAction` request.


## message


```lua
string
```


The diagnostic's message. It usually appears in the user interface

## range


```lua
lsp.Range
```


The range at which the message applies

## relatedInformation


```lua
lsp.DiagnosticRelatedInformation[]?
```


An array of related diagnostic information, e.g. when symbol-names within
a scope collide all definitions can be marked via this property.

## severity


```lua
(1|2|3|4)?
```


The diagnostic's severity. To avoid interpretation mismatches when a
server is used with different clients it is highly recommended that servers
always provide a severity value.

## source


```lua
string?
```


A human-readable string describing the source of this
diagnostic, e.g. 'typescript' or 'super lint'. It usually
appears in the user interface.

## tags


```lua
1|2[]?
```


Additional metadata about the diagnostic.



---

# lsp.DiagnosticClientCapabilities

Client capabilities specific to diagnostic pull requests.


## codeDescriptionSupport


```lua
boolean?
```


Client supports a codeDescription property


## dataSupport


```lua
boolean?
```


Whether code action supports the `data` property which is
preserved between a `textDocument/publishDiagnostics` and
`textDocument/codeAction` request.


## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
return value for the corresponding server capability as well.

## relatedDocumentSupport


```lua
boolean?
```


Whether the clients supports related documents for document diagnostic pulls.

## relatedInformation


```lua
boolean?
```


Whether the clients accepts diagnostics with related information.

## tagSupport


```lua
(lsp.ClientDiagnosticsTagOptions)?
```


Client supports the tag property to provide meta data about a diagnostic.
Clients supporting tags have to handle unknown tags gracefully.



---

# lsp.DiagnosticOptions

Diagnostic options.


## identifier


```lua
string?
```


An optional identifier under which the diagnostics are
managed by the client.

## interFileDependencies


```lua
boolean
```


Whether the language has inter file dependencies meaning that
editing code in one file can result in a different diagnostic
set in another file. Inter file dependencies are common for
most programming languages and typically uncommon for linters.

## workDoneProgress


```lua
boolean?
```


## workspaceDiagnostics


```lua
boolean
```


The server provides support for workspace diagnostics as well.


---

# lsp.DiagnosticRegistrationOptions

Diagnostic registration options.


## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## identifier


```lua
string?
```


An optional identifier under which the diagnostics are
managed by the client.

## interFileDependencies


```lua
boolean
```


Whether the language has inter file dependencies meaning that
editing code in one file can result in a different diagnostic
set in another file. Inter file dependencies are common for
most programming languages and typically uncommon for linters.

## workDoneProgress


```lua
boolean?
```


## workspaceDiagnostics


```lua
boolean
```


The server provides support for workspace diagnostics as well.


---

# lsp.DiagnosticRelatedInformation

Represents a related message and source code location for a diagnostic. This should be
used to point to code locations that cause or related to a diagnostics, e.g when duplicating
a symbol in a scope.

## location


```lua
lsp.Location
```


The location of this related diagnostic information.

## message


```lua
string
```


The message of this related diagnostic information.


---

# lsp.DiagnosticServerCancellationData

Cancellation data returned from a diagnostic request.


## retriggerRequest


```lua
boolean
```



---

# lsp.DiagnosticSeverity

The diagnostic's severity.


---

# lsp.DiagnosticTag

The diagnostic tags.



---

# lsp.DiagnosticWorkspaceClientCapabilities

Workspace client capabilities specific to diagnostic pull requests.


## refreshSupport


```lua
boolean?
```


Whether the client implementation supports a refresh request sent from
the server to the client.

Note that this event is global and will force the client to refresh all
pulled diagnostics currently shown. It should be used with absolute care and
is useful for situation where a server for example detects a project wide
change that requires such a calculation.


---

# lsp.DiagnosticsCapabilities

General diagnostics capabilities for pull and push model.

## codeDescriptionSupport


```lua
boolean?
```


Client supports a codeDescription property


## dataSupport


```lua
boolean?
```


Whether code action supports the `data` property which is
preserved between a `textDocument/publishDiagnostics` and
`textDocument/codeAction` request.


## relatedInformation


```lua
boolean?
```


Whether the clients accepts diagnostics with related information.

## tagSupport


```lua
(lsp.ClientDiagnosticsTagOptions)?
```


Client supports the tag property to provide meta data about a diagnostic.
Clients supporting tags have to handle unknown tags gracefully.



---

# lsp.DidChangeConfigurationClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Did change configuration notification supports dynamic registration.


---

# lsp.DidChangeConfigurationParams

The parameters of a change configuration notification.

## settings


```lua
boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1)
```


The actual changed settings


---

# lsp.DidChangeConfigurationRegistrationOptions

## section


```lua
(string|string[])?
```



---

# lsp.DidChangeNotebookDocumentParams

The params sent in a change notebook document notification.


## change


```lua
lsp.NotebookDocumentChangeEvent
```


The actual changes to the notebook document.

The changes describe single state changes to the notebook document.
So if there are two changes c1 (at array index 0) and c2 (at array
index 1) for a notebook in state S then c1 moves the notebook from
S to S' and c2 from S' to S''. So c1 is computed on the state S and
c2 is computed on the state S'.

To mirror the content of a notebook using change events use the following approach:
- start with the same initial content
- apply the 'notebookDocument/didChange' notifications in the order you receive them.
- apply the `NotebookChangeEvent`s in a single notification in the order
  you receive them.

## notebookDocument


```lua
lsp.VersionedNotebookDocumentIdentifier
```


The notebook document that did change. The version number points
to the version after all provided changes have been applied. If
only the text document content of a cell changes the notebook version
doesn't necessarily have to change.


---

# lsp.DidChangeTextDocumentParams

The change text document notification's parameters.

## contentChanges


```lua
lsp.TextDocumentContentChangePartial|lsp.TextDocumentContentChangeWholeDocument[]
```


The actual content changes. The content changes describe single state changes
to the document. So if there are two content changes c1 (at array index 0) and
c2 (at array index 1) for a document in state S then c1 moves the document from
S to S' and c2 from S' to S''. So c1 is computed on the state S and c2 is computed
on the state S'.

To mirror the content of a document using change events use the following approach:
- start with the same initial content
- apply the 'textDocument/didChange' notifications in the order you receive them.
- apply the `TextDocumentContentChangeEvent`s in a single notification in the order
  you receive them.

## textDocument


```lua
lsp.VersionedTextDocumentIdentifier
```


The document that did change. The version number points
to the version after all provided content changes have
been applied.


---

# lsp.DidChangeWatchedFilesClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Did change watched files notification supports dynamic registration. Please note
that the current protocol doesn't support static configuration for file changes
from the server side.

## relativePatternSupport


```lua
boolean?
```


Whether the client has support for {@link  RelativePattern relative pattern}
or not.



---

# lsp.DidChangeWatchedFilesParams

The watched files change notification's parameters.

## changes


```lua
lsp.FileEvent[]
```


The actual file events.


---

# lsp.DidChangeWatchedFilesRegistrationOptions

Describe options to be used when registered for text document change events.

## watchers


```lua
lsp.FileSystemWatcher[]
```


The watchers to register.


---

# lsp.DidChangeWorkspaceFoldersParams

The parameters of a `workspace/didChangeWorkspaceFolders` notification.

## event


```lua
lsp.WorkspaceFoldersChangeEvent
```


The actual workspace folder change event.


---

# lsp.DidCloseNotebookDocumentParams

The params sent in a close notebook document notification.


## cellTextDocuments


```lua
lsp.TextDocumentIdentifier[]
```


The text documents that represent the content
of a notebook cell that got closed.

## notebookDocument


```lua
lsp.NotebookDocumentIdentifier
```


The notebook document that got closed.


---

# lsp.DidCloseTextDocumentParams

The parameters sent in a close text document notification

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document that was closed.


---

# lsp.DidOpenNotebookDocumentParams

The params sent in an open notebook document notification.


## cellTextDocuments


```lua
lsp.TextDocumentItem[]
```


The text documents that represent the content
of a notebook cell.

## notebookDocument


```lua
lsp.NotebookDocument
```


The notebook document that got opened.


---

# lsp.DidOpenTextDocumentParams

The parameters sent in an open text document notification

## textDocument


```lua
lsp.TextDocumentItem
```


The document that was opened.


---

# lsp.DidSaveNotebookDocumentParams

The params sent in a save notebook document notification.


## notebookDocument


```lua
lsp.NotebookDocumentIdentifier
```


The notebook document that got saved.


---

# lsp.DidSaveTextDocumentParams

The parameters sent in a save text document notification

## text


```lua
string?
```


Optional the content when saved. Depends on the includeText value
when the save notification was requested.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document that was saved.


---

# lsp.DocumentColorClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `DocumentColorRegistrationOptions` return value
for the corresponding server capability as well.


---

# lsp.DocumentColorOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentColorParams

Parameters for a {@link DocumentColorRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentColorRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentDiagnosticParams

Parameters of the document diagnostic request.


## identifier


```lua
string?
```


The additional identifier  provided during registration.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## previousResultId


```lua
string?
```


The result id of a previous response if provided.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentDiagnosticReport

The result of a document diagnostic pull request. A report can
either be a full report containing all diagnostics for the
requested document or an unchanged report indicating that nothing
has changed in terms of diagnostics in comparison to the last
pull request.



---

# lsp.DocumentDiagnosticReportKind

The document diagnostic report kinds.



---

# lsp.DocumentDiagnosticReportPartialResult

A partial result for a document diagnostic report.


## relatedDocuments


```lua
table<string, lsp.FullDocumentDiagnosticReport|lsp.UnchangedDocumentDiagnosticReport>
```



---

# lsp.DocumentFilter

A document filter describes a top level text document or
a notebook cell document.



---

# lsp.DocumentFormattingClientCapabilities

Client capabilities of a {@link DocumentFormattingRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether formatting supports dynamic registration.


---

# lsp.DocumentFormattingOptions

Provider options for a {@link DocumentFormattingRequest}.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentFormattingParams

The parameters of a {@link DocumentFormattingRequest}.

## options


```lua
lsp.FormattingOptions
```


The format options.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document to format.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentFormattingRegistrationOptions

Registration options for a {@link DocumentFormattingRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentHighlight

A document highlight is a range inside a text document which deserves
special attention. Usually a document highlight is visualized by changing
the background color of its range.

## kind


```lua
(1|2|3)?
```


The highlight kind, default is {@link DocumentHighlightKind.Text text}.

## range


```lua
lsp.Range
```


The range this highlight applies to.


---

# lsp.DocumentHighlightClientCapabilities

Client Capabilities for a {@link DocumentHighlightRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether document highlight supports dynamic registration.


---

# lsp.DocumentHighlightKind

A document highlight kind.


---

# lsp.DocumentHighlightOptions

Provider options for a {@link DocumentHighlightRequest}.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentHighlightParams

Parameters for a {@link DocumentHighlightRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentHighlightRegistrationOptions

Registration options for a {@link DocumentHighlightRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentLink

A document link is a range in a text document that links to an internal or external resource, like another
text document or a web site.

## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved on a document link between a
DocumentLinkRequest and a DocumentLinkResolveRequest.

## range


```lua
lsp.Range
```


The range this link applies to.

## target


```lua
string?
```


The uri this link points to. If missing a resolve request is sent later.

## tooltip


```lua
string?
```


The tooltip text when you hover over this link.

If a tooltip is provided, is will be displayed in a string that includes instructions on how to
trigger the link, such as `{0} (ctrl + click)`. The specific instructions vary depending on OS,
user settings, and localization.



---

# lsp.DocumentLinkClientCapabilities

The client capabilities of a {@link DocumentLinkRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether document link supports dynamic registration.

## tooltipSupport


```lua
boolean?
```


Whether the client supports the `tooltip` property on `DocumentLink`.



---

# lsp.DocumentLinkOptions

Provider options for a {@link DocumentLinkRequest}.

## resolveProvider


```lua
boolean?
```


Document links have a resolve provider as well.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentLinkParams

The parameters of a {@link DocumentLinkRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document to provide document links for.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentLinkRegistrationOptions

Registration options for a {@link DocumentLinkRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## resolveProvider


```lua
boolean?
```


Document links have a resolve provider as well.

## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentOnTypeFormattingClientCapabilities

Client capabilities of a {@link DocumentOnTypeFormattingRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether on type formatting supports dynamic registration.


---

# lsp.DocumentOnTypeFormattingOptions

Provider options for a {@link DocumentOnTypeFormattingRequest}.

## firstTriggerCharacter


```lua
string
```


A character on which formatting should be triggered, like `{`.

## moreTriggerCharacter


```lua
string[]?
```


More trigger characters.


---

# lsp.DocumentOnTypeFormattingParams

The parameters of a {@link DocumentOnTypeFormattingRequest}.

## ch


```lua
string
```


The character that has been typed that triggered the formatting
on type request. That is not necessarily the last character that
got inserted into the document since the client could auto insert
characters as well (e.g. like automatic brace completion).

## options


```lua
lsp.FormattingOptions
```


The formatting options.

## position


```lua
lsp.Position
```


The position around which the on type formatting should happen.
This is not necessarily the exact position where the character denoted
by the property `ch` got typed.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document to format.


---

# lsp.DocumentOnTypeFormattingRegistrationOptions

Registration options for a {@link DocumentOnTypeFormattingRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## firstTriggerCharacter


```lua
string
```


A character on which formatting should be triggered, like `{`.

## moreTriggerCharacter


```lua
string[]?
```


More trigger characters.


---

# lsp.DocumentRangeFormattingClientCapabilities

Client capabilities of a {@link DocumentRangeFormattingRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether range formatting supports dynamic registration.

## rangesSupport


```lua
boolean?
```


Whether the client supports formatting multiple ranges at once.



---

# lsp.DocumentRangeFormattingOptions

Provider options for a {@link DocumentRangeFormattingRequest}.

## rangesSupport


```lua
boolean?
```


Whether the server supports formatting multiple ranges at once.


## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentRangeFormattingParams

The parameters of a {@link DocumentRangeFormattingRequest}.

## options


```lua
lsp.FormattingOptions
```


The format options

## range


```lua
lsp.Range
```


The range to format

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document to format.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentRangeFormattingRegistrationOptions

Registration options for a {@link DocumentRangeFormattingRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## rangesSupport


```lua
boolean?
```


Whether the server supports formatting multiple ranges at once.


## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentRangesFormattingParams

The parameters of a {@link DocumentRangesFormattingRequest}.


## options


```lua
lsp.FormattingOptions
```


The format options

## ranges


```lua
lsp.Range[]
```


The ranges to format

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document to format.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentSelector

A document selector is the combination of one or many document filters.

\@sample `let sel:DocumentSelector = [{ language: 'typescript' }, { language: 'json', pattern: '**∕tsconfig.json' }]`;

The use of a string as a document filter is deprecated @since 3.16.0.


---

# lsp.DocumentSymbol

Represents programming constructs like variables, classes, interfaces etc.
that appear in a document. Document symbols can be hierarchical and they
have two ranges: one that encloses its definition and one that points to
its most interesting range, e.g. the range of an identifier.

## children


```lua
lsp.DocumentSymbol[]?
```


Children of this symbol, e.g. properties of a class.

## deprecated


```lua
boolean?
```


Indicates if this symbol is deprecated.


## detail


```lua
string?
```


More detail for this symbol, e.g the signature of a function.

## kind


```lua
1|10|11|12|13...(+21)
```


The kind of this symbol.

## name


```lua
string
```


The name of this symbol. Will be displayed in the user interface and therefore must not be
an empty string or a string only consisting of white spaces.

## range


```lua
lsp.Range
```


The range enclosing this symbol not including leading/trailing whitespace but everything else
like comments. This information is typically used to determine if the clients cursor is
inside the symbol to reveal in the symbol in the UI.

## selectionRange


```lua
lsp.Range
```


The range that should be selected and revealed when this symbol is being picked, e.g the name of a function.
Must be contained by the `range`.

## tags


```lua
1[]?
```


Tags for this document symbol.



---

# lsp.DocumentSymbolClientCapabilities

Client Capabilities for a {@link DocumentSymbolRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether document symbol supports dynamic registration.

## hierarchicalDocumentSymbolSupport


```lua
boolean?
```


The client supports hierarchical document symbols.

## labelSupport


```lua
boolean?
```


The client supports an additional label presented in the UI when
registering a document symbol provider.


## symbolKind


```lua
(lsp.ClientSymbolKindOptions)?
```


Specific capabilities for the `SymbolKind` in the
`textDocument/documentSymbol` request.

## tagSupport


```lua
(lsp.ClientSymbolTagOptions)?
```


The client supports tags on `SymbolInformation`. Tags are supported on
`DocumentSymbol` if `hierarchicalDocumentSymbolSupport` is set to true.
Clients supporting tags have to handle unknown tags gracefully.



---

# lsp.DocumentSymbolOptions

Provider options for a {@link DocumentSymbolRequest}.

## label


```lua
string?
```


A human-readable string that is shown when multiple outlines trees
are shown for the same document.


## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentSymbolParams

Parameters for a {@link DocumentSymbolRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.DocumentSymbolRegistrationOptions

Registration options for a {@link DocumentSymbolRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## label


```lua
string?
```


A human-readable string that is shown when multiple outlines trees
are shown for the same document.


## workDoneProgress


```lua
boolean?
```



---

# lsp.DocumentUri


---

# lsp.DynamicCapabilities


---

# lsp.EditRangeWithInsertReplace

Edit range variant that includes ranges for insert and replace operations.


## insert


```lua
lsp.Range
```


## replace


```lua
lsp.Range
```



---

# lsp.ErrorCodes

Predefined error codes.


---

# lsp.ExecuteCommandClientCapabilities

The client capabilities of a {@link ExecuteCommandRequest}.

## dynamicRegistration


```lua
boolean?
```


Execute command supports dynamic registration.


---

# lsp.ExecuteCommandOptions

The server capabilities of a {@link ExecuteCommandRequest}.

## commands


```lua
string[]
```


The commands to be executed on the server

## workDoneProgress


```lua
boolean?
```



---

# lsp.ExecuteCommandParams

The parameters of a {@link ExecuteCommandRequest}.

## arguments


```lua
boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]?
```


Arguments that the command should be invoked with.

## command


```lua
string
```


The identifier of the actual command handler.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.ExecuteCommandRegistrationOptions

Registration options for a {@link ExecuteCommandRequest}.

## commands


```lua
string[]
```


The commands to be executed on the server

## workDoneProgress


```lua
boolean?
```



---

# lsp.ExecutionSummary

## executionOrder


```lua
integer
```


A strict monotonically increasing value
indicating the execution order of a cell
inside a notebook.

## success


```lua
boolean?
```


Whether the execution was successful or
not if known by the client.


---

# lsp.FailureHandlingKind


---

# lsp.FileChangeType

The file event type


---

# lsp.FileCreate

Represents information on a file/folder create.


## uri


```lua
string
```


A file:// URI for the location of the file/folder being created.


---

# lsp.FileDelete

Represents information on a file/folder delete.


## uri


```lua
string
```


A file:// URI for the location of the file/folder being deleted.


---

# lsp.FileEvent

An event describing a file change.

## type


```lua
1|2|3
```


The change type.

## uri


```lua
string
```


The file's uri.


---

# lsp.FileOperationClientCapabilities

Capabilities relating to events from file operations by the user in the client.

These events do not come from the file system, they come from user operations
like renaming a file in the UI.


## didCreate


```lua
boolean?
```


The client has support for sending didCreateFiles notifications.

## didDelete


```lua
boolean?
```


The client has support for sending didDeleteFiles notifications.

## didRename


```lua
boolean?
```


The client has support for sending didRenameFiles notifications.

## dynamicRegistration


```lua
boolean?
```


Whether the client supports dynamic registration for file requests/notifications.

## willCreate


```lua
boolean?
```


The client has support for sending willCreateFiles requests.

## willDelete


```lua
boolean?
```


The client has support for sending willDeleteFiles requests.

## willRename


```lua
boolean?
```


The client has support for sending willRenameFiles requests.


---

# lsp.FileOperationFilter

A filter to describe in which file operation requests or notifications
the server is interested in receiving.


## pattern


```lua
lsp.FileOperationPattern
```


The actual file operation pattern.

## scheme


```lua
string?
```


A Uri scheme like `file` or `untitled`.


---

# lsp.FileOperationOptions

Options for notifications/requests for user operations on files.


## didCreate


```lua
(lsp.FileOperationRegistrationOptions)?
```


The server is interested in receiving didCreateFiles notifications.

## didDelete


```lua
(lsp.FileOperationRegistrationOptions)?
```


The server is interested in receiving didDeleteFiles file notifications.

## didRename


```lua
(lsp.FileOperationRegistrationOptions)?
```


The server is interested in receiving didRenameFiles notifications.

## willCreate


```lua
(lsp.FileOperationRegistrationOptions)?
```


The server is interested in receiving willCreateFiles requests.

## willDelete


```lua
(lsp.FileOperationRegistrationOptions)?
```


The server is interested in receiving willDeleteFiles file requests.

## willRename


```lua
(lsp.FileOperationRegistrationOptions)?
```


The server is interested in receiving willRenameFiles requests.


---

# lsp.FileOperationPattern

A pattern to describe in which file operation requests or notifications
the server is interested in receiving.


## glob


```lua
string
```


The glob pattern to match. Glob patterns can have the following syntax:
- `*` to match one or more characters in a path segment
- `?` to match on one character in a path segment
- `**` to match any number of path segments, including none
- `{}` to group sub patterns into an OR expression. (e.g. `**/*.{ts,js}` matches all TypeScript and JavaScript files)
- `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
- `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)

## matches


```lua
("file"|"folder")?
```


Whether to match files or folders with this pattern.

Matches both if undefined.

## options


```lua
(lsp.FileOperationPatternOptions)?
```


Additional options used during matching.


---

# lsp.FileOperationPatternKind

A pattern kind describing if a glob pattern matches a file a folder or
both.



---

# lsp.FileOperationPatternOptions

Matching options for the file operation pattern.


## ignoreCase


```lua
boolean?
```


The pattern should be matched ignoring casing.


---

# lsp.FileOperationRegistrationOptions

The options to register for file operations.


## filters


```lua
lsp.FileOperationFilter[]
```


The actual filters.


---

# lsp.FileRename

Represents information on a file/folder rename.


## newUri


```lua
string
```


A file:// URI for the new location of the file/folder being renamed.

## oldUri


```lua
string
```


A file:// URI for the original location of the file/folder being renamed.


---

# lsp.FileSystemWatcher

## globPattern


```lua
string|lsp.RelativePattern
```


The glob pattern to watch. See {@link GlobPattern glob pattern} for more detail.


## kind


```lua
(1|2|4)?
```


The kind of events of interest. If omitted it defaults
to WatchKind.Create | WatchKind.Change | WatchKind.Delete
which is 7.


---

# lsp.FoldingRange

Represents a folding range. To be valid, start and end line must be bigger than zero and smaller
than the number of lines in the document. Clients are free to ignore invalid ranges.

## collapsedText


```lua
string?
```


The text that the client should show when the specified range is
collapsed. If not defined or not supported by the client, a default
will be chosen by the client.


## endCharacter


```lua
integer?
```


The zero-based character offset before the folded range ends. If not defined, defaults to the length of the end line.

## endLine


```lua
integer
```


The zero-based end line of the range to fold. The folded area ends with the line's last character.
To be valid, the end must be zero or larger and smaller than the number of lines in the document.

## kind


```lua
("comment"|"imports"|"region")?
```


Describes the kind of the folding range such as 'comment' or 'region'. The kind
is used to categorize folding ranges and used by commands like 'Fold all comments'.
See {@link FoldingRangeKind} for an enumeration of standardized kinds.

## startCharacter


```lua
integer?
```


The zero-based character offset from where the folded range starts. If not defined, defaults to the length of the start line.

## startLine


```lua
integer
```


The zero-based start line of the range to fold. The folded area starts after the line's last character.
To be valid, the end must be zero or larger and smaller than the number of lines in the document.


---

# lsp.FoldingRangeClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration for folding range
providers. If this is set to `true` the client supports the new
`FoldingRangeRegistrationOptions` return value for the corresponding
server capability as well.

## foldingRange


```lua
(lsp.ClientFoldingRangeOptions)?
```


Specific options for the folding range.


## foldingRangeKind


```lua
(lsp.ClientFoldingRangeKindOptions)?
```


Specific options for the folding range kind.


## lineFoldingOnly


```lua
boolean?
```


If set, the client signals that it only supports folding complete lines.
If set, client will ignore specified `startCharacter` and `endCharacter`
properties in a FoldingRange.

## rangeLimit


```lua
integer?
```


The maximum number of folding ranges that the client prefers to receive
per document. The value serves as a hint, servers are free to follow the
limit.


---

# lsp.FoldingRangeKind

A set of predefined range kinds.


---

# lsp.FoldingRangeOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.FoldingRangeParams

Parameters for a {@link FoldingRangeRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.FoldingRangeRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.FoldingRangeWorkspaceClientCapabilities

Client workspace capabilities specific to folding ranges


## refreshSupport


```lua
boolean?
```


Whether the client implementation supports a refresh request sent from the
server to the client.

Note that this event is global and will force the client to refresh all
folding ranges currently shown. It should be used with absolute care and is
useful for situation where a server for example detects a project wide
change that requires such a calculation.



---

# lsp.FormattingOptions

Value-object describing what options formatting should use.

## insertFinalNewline


```lua
boolean?
```


Insert a newline character at the end of the file if one does not exist.


## insertSpaces


```lua
boolean
```


Prefer spaces over tabs.

## tabSize


```lua
integer
```


Size of a tab in spaces.

## trimFinalNewlines


```lua
boolean?
```


Trim all newlines after the final newline at the end of the file.


## trimTrailingWhitespace


```lua
boolean?
```


Trim trailing whitespace on a line.



---

# lsp.FullDocumentDiagnosticReport

A diagnostic report with a full set of problems.


## items


```lua
lsp.Diagnostic[]
```


The actual items.

## kind


```lua
"full"
```


A full document diagnostic report.

## resultId


```lua
string?
```


An optional result id. If provided it will
be sent on the next diagnostic request for the
same document.


---

# lsp.GeneralClientCapabilities

General client capabilities.


## markdown


```lua
(lsp.MarkdownClientCapabilities)?
```


Client capabilities specific to the client's markdown parser.


## positionEncodings


```lua
"utf-16"|"utf-32"|"utf-8"[]?
```


The position encodings supported by the client. Client and server
have to agree on the same position encoding to ensure that offsets
(e.g. character position in a line) are interpreted the same on both
sides.

To keep the protocol backwards compatible the following applies: if
the value 'utf-16' is missing from the array of position encodings
servers can assume that the client supports UTF-16. UTF-16 is
therefore a mandatory encoding.

If omitted it defaults to ['utf-16'].

Implementation considerations: since the conversion from one encoding
into another requires the content of the file / line the conversion
is best done where the file is read which is usually on the server
side.


## regularExpressions


```lua
(lsp.RegularExpressionsClientCapabilities)?
```


Client capabilities specific to regular expressions.


## staleRequestSupport


```lua
(lsp.StaleRequestSupportOptions)?
```


Client capability that signals how the client
handles stale requests (e.g. a request
for which the client will not process the response
anymore since the information is outdated).



---

# lsp.GlobPattern

The glob pattern. Either a string pattern or a relative pattern.



---

# lsp.Handler


---

# lsp.HandlerContext

## bufnr


```lua
integer?
```

## client_id


```lua
integer
```

## method


```lua
string
```

## params


```lua
any
```

## version


```lua
integer?
```


---

# lsp.Hover

The result of a hover request.

## contents


```lua
string|lsp.MarkedStringWithLanguage|lsp.MarkupContent|string|lsp.MarkedStringWithLanguage[]
```


The hover's content

## range


```lua
(lsp.Range)?
```


An optional range inside the text document that is used to
visualize the hover, e.g. by changing the background color.


---

# lsp.HoverClientCapabilities

## contentFormat


```lua
"markdown"|"plaintext"[]?
```


Client supports the following content formats for the content
property. The order describes the preferred format of the client.

## dynamicRegistration


```lua
boolean?
```


Whether hover supports dynamic registration.


---

# lsp.HoverOptions

Hover options.

## workDoneProgress


```lua
boolean?
```



---

# lsp.HoverParams

Parameters for a {@link HoverRequest}.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.HoverRegistrationOptions

Registration options for a {@link HoverRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## workDoneProgress


```lua
boolean?
```



---

# lsp.ImplementationClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `ImplementationRegistrationOptions` return value
for the corresponding server capability as well.

## linkSupport


```lua
boolean?
```


The client supports additional metadata in the form of definition links.



---

# lsp.ImplementationOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.ImplementationParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.ImplementationRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.InitializeError

The data type of the ResponseError if the
initialize request fails.

## retry


```lua
boolean
```


Indicates whether the client execute the following retry logic:
(1) show the message provided by the ResponseError to the user
(2) user selects retry or cancel
(3) if user selected retry the initialize method is sent again.


---

# lsp.InitializeParams

## capabilities


```lua
lsp.ClientCapabilities
```


The capabilities provided by the client (editor or tool)

## clientInfo


```lua
(lsp.ClientInfo)?
```


Information about the client


## initializationOptions


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


User provided initialization options.

## locale


```lua
string?
```


The locale the client is currently showing the user interface
in. This must not necessarily be the locale of the operating
system.

Uses IETF language tags as the value's syntax
(See https://en.wikipedia.org/wiki/IETF_language_tag)


## processId


```lua
integer|nil
```


The process Id of the parent process that started
the server.

Is `null` if the process has not been started by another process.
If the parent process is not alive then the server should exit.

## rootPath


```lua
(string|nil)?
```


The rootPath of the workspace. Is null
if no folder is open.


## rootUri


```lua
string|nil
```


The rootUri of the workspace. Is null if no
folder is open. If both `rootPath` and `rootUri` are set
`rootUri` wins.


## trace


```lua
("messages"|"off"|"verbose")?
```


The initial trace setting. If omitted trace is disabled ('off').

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.

## workspaceFolders


```lua
(lsp.WorkspaceFolder[]|nil)?
```


The workspace folders configured in the client when the server starts.

This property is only available if the client supports workspace folders.
It can be `null` if the client supports workspace folders but none are
configured.



---

# lsp.InitializeResult

The result returned from an initialize request.

## capabilities


```lua
lsp.ServerCapabilities
```


The capabilities the language server provides.

## serverInfo


```lua
(lsp.ServerInfo)?
```


Information about the server.



---

# lsp.InitializedParams


---

# lsp.InlayHint

Inlay hint information.


## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved on an inlay hint between
a `textDocument/inlayHint` and a `inlayHint/resolve` request.

## kind


```lua
(1|2)?
```


The kind of this hint. Can be omitted in which case the client
should fall back to a reasonable default.

## label


```lua
string|lsp.InlayHintLabelPart[]
```


The label of this hint. A human readable string or an array of
InlayHintLabelPart label parts.

*Note* that neither the string nor the label part can be empty.

## paddingLeft


```lua
boolean?
```


Render padding before the hint.

Note: Padding should use the editor's background color, not the
background color of the hint itself. That means padding can be used
to visually align/separate an inlay hint.

## paddingRight


```lua
boolean?
```


Render padding after the hint.

Note: Padding should use the editor's background color, not the
background color of the hint itself. That means padding can be used
to visually align/separate an inlay hint.

## position


```lua
lsp.Position
```


The position of this hint.

If multiple hints have the same position, they will be shown in the order
they appear in the response.

## textEdits


```lua
lsp.TextEdit[]?
```


Optional text edits that are performed when accepting this inlay hint.

*Note* that edits are expected to change the document so that the inlay
hint (or its nearest variant) is now part of the document and the inlay
hint itself is now obsolete.

## tooltip


```lua
(string|lsp.MarkupContent)?
```


The tooltip text when you hover over this item.


---

# lsp.InlayHintClientCapabilities

Inlay hint client capabilities.


## dynamicRegistration


```lua
boolean?
```


Whether inlay hints support dynamic registration.

## resolveSupport


```lua
(lsp.ClientInlayHintResolveOptions)?
```


Indicates which properties a client can resolve lazily on an inlay
hint.


---

# lsp.InlayHintKind

Inlay hint kinds.



---

# lsp.InlayHintLabelPart

An inlay hint label part allows for interactive and composite labels
of inlay hints.


## command


```lua
(lsp.Command)?
```


An optional command for this label part.

Depending on the client capability `inlayHint.resolveSupport` clients
might resolve this property late using the resolve request.

## location


```lua
(lsp.Location)?
```


An optional source code location that represents this
label part.

The editor will use this location for the hover and for code navigation
features: This part will become a clickable link that resolves to the
definition of the symbol at the given location (not necessarily the
location itself), it shows the hover that shows at the given location,
and it shows a context menu with further code navigation commands.

Depending on the client capability `inlayHint.resolveSupport` clients
might resolve this property late using the resolve request.

## tooltip


```lua
(string|lsp.MarkupContent)?
```


The tooltip text when you hover over this label part. Depending on
the client capability `inlayHint.resolveSupport` clients might resolve
this property late using the resolve request.

## value


```lua
string
```


The value of this label part.


---

# lsp.InlayHintOptions

Inlay hint options used during static registration.


## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for an inlay hint item.

## workDoneProgress


```lua
boolean?
```



---

# lsp.InlayHintParams

A parameter literal used in inlay hint requests.


## range


```lua
lsp.Range
```


The document range for which inlay hints should be computed.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.InlayHintRegistrationOptions

Inlay hint options used during static or dynamic registration.


## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for an inlay hint item.

## workDoneProgress


```lua
boolean?
```



---

# lsp.InlayHintWorkspaceClientCapabilities

Client workspace capabilities specific to inlay hints.


## refreshSupport


```lua
boolean?
```


Whether the client implementation supports a refresh request sent from
the server to the client.

Note that this event is global and will force the client to refresh all
inlay hints currently shown. It should be used with absolute care and
is useful for situation where a server for example detects a project wide
change that requires such a calculation.


---

# lsp.InlineCompletionClientCapabilities

Client capabilities specific to inline completions.


## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration for inline completion providers.


---

# lsp.InlineCompletionContext

Provides information about the context in which an inline completion was requested.


## selectedCompletionInfo


```lua
(lsp.SelectedCompletionInfo)?
```


Provides information about the currently selected item in the autocomplete widget if it is visible.

## triggerKind


```lua
1|2
```


Describes how the inline completion was triggered.


---

# lsp.InlineCompletionItem

An inline completion item represents a text snippet that is proposed inline to complete text that is being typed.


## command


```lua
(lsp.Command)?
```


An optional {@link Command} that is executed *after* inserting this completion.

## filterText


```lua
string?
```


A text that is used to decide if this inline completion should be shown. When `falsy` the {@link InlineCompletionItem.insertText} is used.

## insertText


```lua
string|lsp.StringValue
```


The text to replace the range with. Must be set.

## range


```lua
(lsp.Range)?
```


The range to replace. Must begin and end on the same line.


---

# lsp.InlineCompletionList

Represents a collection of {@link InlineCompletionItem inline completion items} to be presented in the editor.


## items


```lua
lsp.InlineCompletionItem[]
```


The inline completion items


---

# lsp.InlineCompletionOptions

Inline completion options used during static registration.


## workDoneProgress


```lua
boolean?
```



---

# lsp.InlineCompletionParams

A parameter literal used in inline completion requests.


## context


```lua
lsp.InlineCompletionContext
```


Additional information about the context in which inline completions were
requested.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.InlineCompletionRegistrationOptions

Inline completion options used during static or dynamic registration.


## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.InlineCompletionTriggerKind

Describes how an {@link InlineCompletionItemProvider inline completion provider} was triggered.



---

# lsp.InlineValue

Inline value information can be provided by different means:
- directly as a text value (class InlineValueText).
- as a name to use for a variable lookup (class InlineValueVariableLookup)
- as an evaluatable expression (class InlineValueEvaluatableExpression)
The InlineValue types combines all inline value types into one type.



---

# lsp.InlineValueClientCapabilities

Client capabilities specific to inline values.


## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration for inline value providers.


---

# lsp.InlineValueContext

## frameId


```lua
integer
```


The stack frame (as a DAP Id) where the execution has stopped.

## stoppedLocation


```lua
lsp.Range
```


The document range where execution has stopped.
Typically the end position of the range denotes the line where the inline values are shown.


---

# lsp.InlineValueEvaluatableExpression

Provide an inline value through an expression evaluation.
If only a range is specified, the expression will be extracted from the underlying document.
An optional expression can be used to override the extracted expression.


## expression


```lua
string?
```


If specified the expression overrides the extracted expression.

## range


```lua
lsp.Range
```


The document range for which the inline value applies.
The range is used to extract the evaluatable expression from the underlying document.


---

# lsp.InlineValueOptions

Inline value options used during static registration.


## workDoneProgress


```lua
boolean?
```



---

# lsp.InlineValueParams

A parameter literal used in inline value requests.


## context


```lua
lsp.InlineValueContext
```


Additional information about the context in which inline values were
requested.

## range


```lua
lsp.Range
```


The document range for which inline values should be computed.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.InlineValueRegistrationOptions

Inline value options used during static or dynamic registration.


## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.InlineValueText

Provide inline value as text.


## range


```lua
lsp.Range
```


The document range for which the inline value applies.

## text


```lua
string
```


The text of the inline value.


---

# lsp.InlineValueVariableLookup

Provide inline value through a variable lookup.
If only a range is specified, the variable name will be extracted from the underlying document.
An optional variable name can be used to override the extracted name.


## caseSensitiveLookup


```lua
boolean
```


How to perform the lookup.

## range


```lua
lsp.Range
```


The document range for which the inline value applies.
The range is used to extract the variable name from the underlying document.

## variableName


```lua
string?
```


If specified the name of the variable to look up.


---

# lsp.InlineValueWorkspaceClientCapabilities

Client workspace capabilities specific to inline values.


## refreshSupport


```lua
boolean?
```


Whether the client implementation supports a refresh request sent from the
server to the client.

Note that this event is global and will force the client to refresh all
inline values currently shown. It should be used with absolute care and is
useful for situation where a server for example detects a project wide
change that requires such a calculation.


---

# lsp.InsertReplaceEdit

A special text edit to provide an insert and a replace operation.


## insert


```lua
lsp.Range
```


The range if the insert is requested

## newText


```lua
string
```


The string to be inserted.

## replace


```lua
lsp.Range
```


The range if the replace is requested.


---

# lsp.InsertTextFormat

Defines whether the insert text in a completion item should be interpreted as
plain text or a snippet.


---

# lsp.InsertTextMode

How whitespace and indentation is handled during completion
item insertion.



---

# lsp.ItemDefaults

 TODO(mariasolos): Remove this declaration once we figure out a better way to handle
 literal/anonymous types (see https://github.com/neovim/neovim/pull/27542/files#r1495259331).
 @nodoc

## data


```lua
any
```

## editRange


```lua
lsp.Range|{ insert: lsp.Range, replace: lsp.Range }|nil
```

A range in a text document expressed as (zero-based) start and end positions.

If you want to specify a range that contains a line including the line ending
character(s) then use an end position denoting the start of the next line.
For example:
```ts
{
    start: { line: 5, character: 23 }
    end : { line 6, character : 0 }
}
```

## insertTextFormat


```lua
(1|2)?
```

Defines whether the insert text in a completion item should be interpreted as
plain text or a snippet.

## insertTextMode


```lua
(1|2)?
```

How whitespace and indentation is handled during completion
item insertion.



---

# lsp.LSPAny

The LSP any type.
Please note that strictly speaking a property with the value `undefined`
can't be converted into JSON preserving the property name. However for
convenience it is allowed and assumed that all these properties are
optional as well.


---

# lsp.LSPArray

LSP arrays.


---

# lsp.LSPErrorCodes


---

# lsp.LSPObject

LSP object definition.


---

# lsp.LanguageKind

Predefined Language kinds


---

# lsp.LinkedEditingRangeClientCapabilities

Client capabilities for the linked editing range request.


## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
return value for the corresponding server capability as well.


---

# lsp.LinkedEditingRangeOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.LinkedEditingRangeParams

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.LinkedEditingRangeRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.LinkedEditingRanges

The result of a linked editing range request.


## ranges


```lua
lsp.Range[]
```


A list of ranges that can be edited together. The ranges must have
identical length and contain identical text content. The ranges cannot overlap.

## wordPattern


```lua
string?
```


An optional word pattern (regular expression) that describes valid contents for
the given ranges. If no pattern is provided, the client configuration's word
pattern will be used.


---

# lsp.Location

Represents a location inside a resource, such as a line
inside a text file.

## range


```lua
lsp.Range
```


## uri


```lua
string
```



---

# lsp.LocationLink

Represents the connection of two locations. Provides additional metadata over normal {@link Location locations},
including an origin range.

## originSelectionRange


```lua
(lsp.Range)?
```


Span of the origin of this link.

Used as the underlined span for mouse interaction. Defaults to the word range at
the definition position.

## targetRange


```lua
lsp.Range
```


The full target range of this link. If the target for example is a symbol then target range is the
range enclosing this symbol not including leading/trailing whitespace but everything else
like comments. This information is typically used to highlight the range in the editor.

## targetSelectionRange


```lua
lsp.Range
```


The range that should be selected and revealed when this link is being followed, e.g the name of a function.
Must be contained by the `targetRange`. See also `DocumentSymbol#range`

## targetUri


```lua
string
```


The target resource identifier of this link.


---

# lsp.LocationUriOnly

Location with only uri and does not include range.


## uri


```lua
string
```



---

# lsp.LogMessageParams

The log message parameters.

## message


```lua
string
```


The actual message.

## type


```lua
1|2|3|4|5
```


The message type. See {@link MessageType}


---

# lsp.LogTraceParams

## message


```lua
string
```


## verbose


```lua
string?
```



---

# lsp.MarkdownClientCapabilities

Client capabilities specific to the used markdown parser.


## allowedTags


```lua
string[]?
```


A list of HTML tags that the client allows / supports in
Markdown.


## parser


```lua
string
```


The name of the parser.

## version


```lua
string?
```


The version of the parser.


---

# lsp.MarkedString

MarkedString can be used to render human readable text. It is either a markdown string
or a code-block that provides a language and a code snippet. The language identifier
is semantically equal to the optional language identifier in fenced code blocks in GitHub
issues. See https://help.github.com/articles/creating-and-highlighting-code-blocks/#syntax-highlighting

The pair of a language and a value is an equivalent to markdown:
```${language}
${value}
```

Note that markdown strings will be sanitized - that means html will be escaped.


---

# lsp.MarkedStringWithLanguage

## language


```lua
string
```


## value


```lua
string
```



---

# lsp.MarkupContent

A `MarkupContent` literal represents a string value which content is interpreted base on its
kind flag. Currently the protocol supports `plaintext` and `markdown` as markup kinds.

If the kind is `markdown` then the value can contain fenced code blocks like in GitHub issues.
See https://help.github.com/articles/creating-and-highlighting-code-blocks/#syntax-highlighting

Here is an example how such a string can be constructed using JavaScript / TypeScript:
```ts
let markdown: MarkdownContent = {
 kind: MarkupKind.Markdown,
 value: [
   '# Header',
   'Some text',
   '```typescript',
   'someCode();',
   '```'
 ].join('\n')
};
```

*Please Note* that clients might sanitize the return markdown. A client could decide to
remove HTML from the markdown to avoid script execution.

## kind


```lua
"markdown"|"plaintext"
```


The type of the Markup

## value


```lua
string
```


The content itself


---

# lsp.MarkupKind

Describes the content type that a client supports in various
result literals like `Hover`, `ParameterInfo` or `CompletionItem`.

Please note that `MarkupKinds` must not start with a `$`. This kinds
are reserved for internal usage.


---

# lsp.MessageActionItem

## title


```lua
string
```


A short title like 'Retry', 'Open Log' etc.


---

# lsp.MessageType

The message type


---

# lsp.Moniker

Moniker definition to match LSIF 0.5 moniker definition.


## identifier


```lua
string
```


The identifier of the moniker. The value is opaque in LSIF however
schema owners are allowed to define the structure if they want.

## kind


```lua
("export"|"import"|"local")?
```


The moniker kind if known.

## scheme


```lua
string
```


The scheme of the moniker. For example tsc or .Net

## unique


```lua
"document"|"global"|"group"|"project"|"scheme"
```


The scope in which the moniker is unique


---

# lsp.MonikerClientCapabilities

Client capabilities specific to the moniker request.


## dynamicRegistration


```lua
boolean?
```


Whether moniker supports dynamic registration. If this is set to `true`
the client supports the new `MonikerRegistrationOptions` return value
for the corresponding server capability as well.


---

# lsp.MonikerKind

The moniker kind.



---

# lsp.MonikerOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.MonikerParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.MonikerRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## workDoneProgress


```lua
boolean?
```



---

# lsp.MultiHandler


---

# lsp.NotebookCell

A notebook cell.

A cell's document URI must be unique across ALL notebook
cells and can therefore be used to uniquely identify a
notebook cell or the cell's text document.


## document


```lua
string
```


The URI of the cell's text document
content.

## executionSummary


```lua
(lsp.ExecutionSummary)?
```


Additional execution summary information
if supported by the client.

## kind


```lua
1|2
```


The cell's kind

## metadata


```lua
table<string, lsp.LSPAny>?
```


Additional metadata stored with the cell.

Note: should always be an object literal (e.g. LSPObject)


---

# lsp.NotebookCellArrayChange

A change describing how to move a `NotebookCell`
array from state S to S'.


## cells


```lua
lsp.NotebookCell[]?
```


The new cells, if any

## deleteCount


```lua
integer
```


The deleted cells

## start


```lua
integer
```


The start oftest of the cell that changed.


---

# lsp.NotebookCellKind

A notebook cell kind.



---

# lsp.NotebookCellLanguage

## language


```lua
string
```



---

# lsp.NotebookCellTextDocumentFilter

A notebook cell text document filter denotes a cell text
document by different properties.


## language


```lua
string?
```


A language id like `python`.

Will be matched against the language id of the
notebook cell document. '*' matches every language.

## notebook


```lua
string|lsp.NotebookDocumentFilterNotebookType|lsp.NotebookDocumentFilterPattern|lsp.NotebookDocumentFilterScheme
```


A filter that matches against the notebook
containing the notebook cell. If a string
value is provided it matches against the
notebook type. '*' matches every notebook.


---

# lsp.NotebookDocument

A notebook document.


## cells


```lua
lsp.NotebookCell[]
```


The cells of a notebook.

## metadata


```lua
table<string, lsp.LSPAny>?
```


Additional metadata stored with the notebook
document.

Note: should always be an object literal (e.g. LSPObject)

## notebookType


```lua
string
```


The type of the notebook.

## uri


```lua
string
```


The notebook document's uri.

## version


```lua
integer
```


The version number of this document (it will increase after each
change, including undo/redo).


---

# lsp.NotebookDocumentCellChangeStructure

Structural changes to cells in a notebook document.


## array


```lua
lsp.NotebookCellArrayChange
```


The change to the cell array.

## didClose


```lua
lsp.TextDocumentIdentifier[]?
```


Additional closed cell text documents.

## didOpen


```lua
lsp.TextDocumentItem[]?
```


Additional opened cell text documents.


---

# lsp.NotebookDocumentCellChanges

Cell changes to a notebook document.


## data


```lua
lsp.NotebookCell[]?
```


Changes to notebook cells properties like its
kind, execution summary or metadata.

## structure


```lua
(lsp.NotebookDocumentCellChangeStructure)?
```


Changes to the cell structure to add or
remove cells.

## textContent


```lua
lsp.NotebookDocumentCellContentChanges[]?
```


Changes to the text content of notebook cells.


---

# lsp.NotebookDocumentCellContentChanges

Content changes to a cell in a notebook document.


## changes


```lua
lsp.TextDocumentContentChangePartial|lsp.TextDocumentContentChangeWholeDocument[]
```


## document


```lua
lsp.VersionedTextDocumentIdentifier
```



---

# lsp.NotebookDocumentChangeEvent

A change event for a notebook document.


## cells


```lua
(lsp.NotebookDocumentCellChanges)?
```


Changes to cells

## metadata


```lua
table<string, lsp.LSPAny>?
```


The changed meta data if any.

Note: should always be an object literal (e.g. LSPObject)


---

# lsp.NotebookDocumentClientCapabilities

Capabilities specific to the notebook document support.


## synchronization


```lua
lsp.NotebookDocumentSyncClientCapabilities
```


Capabilities specific to notebook document synchronization



---

# lsp.NotebookDocumentFilter

A notebook document filter denotes a notebook document by
different properties. The properties will be match
against the notebook's URI (same as with documents)



---

# lsp.NotebookDocumentFilterNotebookType

A notebook document filter where `notebookType` is required field.


## notebookType


```lua
string
```


The type of the enclosing notebook.

## pattern


```lua
(string|lsp.RelativePattern)?
```


A glob pattern.

## scheme


```lua
string?
```


A Uri {@link Uri.scheme scheme}, like `file` or `untitled`.


---

# lsp.NotebookDocumentFilterPattern

A notebook document filter where `pattern` is required field.


## notebookType


```lua
string?
```


The type of the enclosing notebook.

## pattern


```lua
string|lsp.RelativePattern
```


A glob pattern.

## scheme


```lua
string?
```


A Uri {@link Uri.scheme scheme}, like `file` or `untitled`.


---

# lsp.NotebookDocumentFilterScheme

A notebook document filter where `scheme` is required field.


## notebookType


```lua
string?
```


The type of the enclosing notebook.

## pattern


```lua
(string|lsp.RelativePattern)?
```


A glob pattern.

## scheme


```lua
string
```


A Uri {@link Uri.scheme scheme}, like `file` or `untitled`.


---

# lsp.NotebookDocumentFilterWithCells

## cells


```lua
lsp.NotebookCellLanguage[]
```


The cells of the matching notebook to be synced.

## notebook


```lua
(string|lsp.NotebookDocumentFilterNotebookType|lsp.NotebookDocumentFilterPattern|lsp.NotebookDocumentFilterScheme)?
```


The notebook to be synced If a string
value is provided it matches against the
notebook type. '*' matches every notebook.


---

# lsp.NotebookDocumentFilterWithNotebook

## cells


```lua
lsp.NotebookCellLanguage[]?
```


The cells of the matching notebook to be synced.

## notebook


```lua
string|lsp.NotebookDocumentFilterNotebookType|lsp.NotebookDocumentFilterPattern|lsp.NotebookDocumentFilterScheme
```


The notebook to be synced If a string
value is provided it matches against the
notebook type. '*' matches every notebook.


---

# lsp.NotebookDocumentIdentifier

A literal to identify a notebook document in the client.


## uri


```lua
string
```


The notebook document's uri.


---

# lsp.NotebookDocumentSyncClientCapabilities

Notebook specific client capabilities.


## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is
set to `true` the client supports the new
`(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
return value for the corresponding server capability as well.

## executionSummarySupport


```lua
boolean?
```


The client supports sending execution summary data per cell.


---

# lsp.NotebookDocumentSyncOptions

Options specific to a notebook plus its cells
to be synced to the server.

If a selector provides a notebook document
filter but no cell selector all cells of a
matching notebook document will be synced.

If a selector provides no notebook document
filter but only a cell selector all notebook
document that contain at least one matching
cell will be synced.


## notebookSelector


```lua
(lsp.NotebookDocumentFilterWithCells|lsp.NotebookDocumentFilterWithNotebook)[]
```


The notebooks to be synced

## save


```lua
boolean?
```


Whether save notification should be forwarded to
the server. Will only be honored if mode === `notebook`.


---

# lsp.NotebookDocumentSyncRegistrationOptions

Registration options specific to a notebook.


## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## notebookSelector


```lua
(lsp.NotebookDocumentFilterWithCells|lsp.NotebookDocumentFilterWithNotebook)[]
```


The notebooks to be synced

## save


```lua
boolean?
```


Whether save notification should be forwarded to
the server. Will only be honored if mode === `notebook`.


---

# lsp.OptionalVersionedTextDocumentIdentifier

A text document identifier to optionally denote a specific version of a text document.

## uri


```lua
string
```


The text document's uri.

## version


```lua
integer|nil
```


The version number of this document. If a versioned text document identifier
is sent from the server to the client and the file is not open in the editor
(the server has not received an open notification before) the server can send
`null` to indicate that the version is unknown and the content on disk is the
truth (as specified with document content ownership).


---

# lsp.ParameterInformation

Represents a parameter of a callable-signature. A parameter can
have a label and a doc-comment.

## documentation


```lua
(string|lsp.MarkupContent)?
```


The human-readable doc-comment of this parameter. Will be shown
in the UI but can be omitted.

## label


```lua
string|[integer, integer]
```


The label of this parameter information.

Either a string or an inclusive start and exclusive end offsets within its containing
signature label. (see SignatureInformation.label). The offsets are based on a UTF-16
string representation as `Position` and `Range` does.

To avoid ambiguities a server should use the [start, end] offset value instead of using
a substring. Whether a client support this is controlled via `labelOffsetSupport` client
capability.

*Note*: a label of type string should be a substring of its containing signature label.
Its intended use case is to highlight the parameter label part in the `SignatureInformation.label`.


---

# lsp.PartialResultParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.


---

# lsp.Pattern

The glob pattern to watch relative to the base path. Glob patterns can have the following syntax:
- `*` to match one or more characters in a path segment
- `?` to match on one character in a path segment
- `**` to match any number of path segments, including none
- `{}` to group conditions (e.g. `**/*.{ts,js}` matches all TypeScript and JavaScript files)
- `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
- `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)



---

# lsp.Position

Position in a text document expressed as zero-based line and character
offset. Prior to 3.17 the offsets were always based on a UTF-16 string
representation. So a string of the form `a𐐀b` the character offset of the
character `a` is 0, the character offset of `𐐀` is 1 and the character
offset of b is 3 since `𐐀` is represented using two code units in UTF-16.
Since 3.17 clients and servers can agree on a different string encoding
representation (e.g. UTF-8). The client announces it's supported encoding
via the client capability [`general.positionEncodings`](https://microsoft.github.io/language-server-protocol/specifications/specification-current/#clientCapabilities).
The value is an array of position encodings the client supports, with
decreasing preference (e.g. the encoding at index `0` is the most preferred
one). To stay backwards compatible the only mandatory encoding is UTF-16
represented via the string `utf-16`. The server can pick one of the
encodings offered by the client and signals that encoding back to the
client via the initialize result's property
[`capabilities.positionEncoding`](https://microsoft.github.io/language-server-protocol/specifications/specification-current/#serverCapabilities). If the string value
`utf-16` is missing from the client's capability `general.positionEncodings`
servers can safely assume that the client supports UTF-16. If the server
omits the position encoding in its initialize result the encoding defaults
to the string value `utf-16`. Implementation considerations: since the
conversion from one encoding into another requires the content of the
file / line the conversion is best done where the file is read which is
usually on the server side.

Positions are line end character agnostic. So you can not specify a position
that denotes `\r|\n` or `\n|` where `|` represents the character offset.


## character


```lua
integer
```


Character offset on a line in a document (zero-based).

The meaning of this offset is determined by the negotiated
`PositionEncodingKind`.

## line


```lua
integer
```


Line position in a document (zero-based).


---

# lsp.PositionEncodingKind

A set of predefined position encoding kinds.



---

# lsp.PrepareRenameDefaultBehavior

## defaultBehavior


```lua
boolean
```



---

# lsp.PrepareRenameParams

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.PrepareRenamePlaceholder

## placeholder


```lua
string
```


## range


```lua
lsp.Range
```



---

# lsp.PrepareRenameResult


---

# lsp.PrepareSupportDefaultBehavior


---

# lsp.PreviousResultId

A previous result id in a workspace pull request.


## uri


```lua
string
```


The URI for which the client knowns a
result id.

## value


```lua
string
```


The value of the previous result id.


---

# lsp.ProgressParams

## token


```lua
string|integer
```


The progress token provided by the client or server.

## value


```lua
boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1)
```


The progress data.


---

# lsp.ProgressToken


---

# lsp.PublishDiagnosticsClientCapabilities

The publish diagnostic client capabilities.

## codeDescriptionSupport


```lua
boolean?
```


Client supports a codeDescription property


## dataSupport


```lua
boolean?
```


Whether code action supports the `data` property which is
preserved between a `textDocument/publishDiagnostics` and
`textDocument/codeAction` request.


## relatedInformation


```lua
boolean?
```


Whether the clients accepts diagnostics with related information.

## tagSupport


```lua
(lsp.ClientDiagnosticsTagOptions)?
```


Client supports the tag property to provide meta data about a diagnostic.
Clients supporting tags have to handle unknown tags gracefully.


## versionSupport


```lua
boolean?
```


Whether the client interprets the version property of the
`textDocument/publishDiagnostics` notification's parameter.



---

# lsp.PublishDiagnosticsParams

The publish diagnostic notification's parameters.

## diagnostics


```lua
lsp.Diagnostic[]
```


An array of diagnostic information items.

## uri


```lua
string
```


The URI for which diagnostic information is reported.

## version


```lua
integer?
```


Optional the version number of the document the diagnostics are published for.



---

# lsp.Range

A range in a text document expressed as (zero-based) start and end positions.

If you want to specify a range that contains a line including the line ending
character(s) then use an end position denoting the start of the next line.
For example:
```ts
{
    start: { line: 5, character: 23 }
    end : { line 6, character : 0 }
}
```

## end


```lua
lsp.Position
```


The range's end position.

## start


```lua
lsp.Position
```


The range's start position.


---

# lsp.ReferenceClientCapabilities

Client Capabilities for a {@link ReferencesRequest}.

## dynamicRegistration


```lua
boolean?
```


Whether references supports dynamic registration.


---

# lsp.ReferenceContext

Value-object that contains additional information when
requesting references.

## includeDeclaration


```lua
boolean
```


Include the declaration of the current symbol.


---

# lsp.ReferenceOptions

Reference options.

## workDoneProgress


```lua
boolean?
```



---

# lsp.ReferenceParams

Parameters for a {@link ReferencesRequest}.

## context


```lua
lsp.ReferenceContext
```


## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.ReferenceRegistrationOptions

Registration options for a {@link ReferencesRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## workDoneProgress


```lua
boolean?
```



---

# lsp.Registration

General parameters to register for a notification or to register a provider.

## id


```lua
string
```


The id used to register the request. The id can be used to deregister
the request again.

## method


```lua
string
```


The method / capability to register for.

## registerOptions


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


Options necessary for the registration.


---

# lsp.RegistrationParams

## registrations


```lua
lsp.Registration[]
```



---

# lsp.RegularExpressionEngineKind


---

# lsp.RegularExpressionsClientCapabilities

Client capabilities specific to regular expressions.


## engine


```lua
string
```


The engine's name.

## version


```lua
string?
```


The engine's version.


---

# lsp.RelatedFullDocumentDiagnosticReport

A full diagnostic report with a set of related documents.


## items


```lua
lsp.Diagnostic[]
```


The actual items.

## kind


```lua
"full"
```


A full document diagnostic report.

## relatedDocuments


```lua
table<string, lsp.FullDocumentDiagnosticReport|lsp.UnchangedDocumentDiagnosticReport>?
```


Diagnostics of related documents. This information is useful
in programming languages where code in a file A can generate
diagnostics in a file B which A depends on. An example of
such a language is C/C++ where marco definitions in a file
a.cpp and result in errors in a header file b.hpp.


## resultId


```lua
string?
```


An optional result id. If provided it will
be sent on the next diagnostic request for the
same document.


---

# lsp.RelatedUnchangedDocumentDiagnosticReport

An unchanged diagnostic report with a set of related documents.


## kind


```lua
"unchanged"
```


A document diagnostic report indicating
no changes to the last result. A server can
only return `unchanged` if result ids are
provided.

## relatedDocuments


```lua
table<string, lsp.FullDocumentDiagnosticReport|lsp.UnchangedDocumentDiagnosticReport>?
```


Diagnostics of related documents. This information is useful
in programming languages where code in a file A can generate
diagnostics in a file B which A depends on. An example of
such a language is C/C++ where marco definitions in a file
a.cpp and result in errors in a header file b.hpp.


## resultId


```lua
string
```


A result id which will be sent on the next
diagnostic request for the same document.


---

# lsp.RelativePattern

A relative pattern is a helper to construct glob patterns that are matched
relatively to a base URI. The common value for a `baseUri` is a workspace
folder root, but it can be another absolute URI as well.


## baseUri


```lua
string|lsp.WorkspaceFolder
```


A workspace folder or a base URI to which this pattern will be matched
against relatively.

## pattern


```lua
string
```


The actual glob pattern;


---

# lsp.RenameClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether rename supports dynamic registration.

## honorsChangeAnnotations


```lua
boolean?
```


Whether the client honors the change annotations in
text edits and resource operations returned via the
rename request's workspace edit by for example presenting
the workspace edit in the user interface and asking
for confirmation.


## prepareSupport


```lua
boolean?
```


Client supports testing for validity of rename operations
before execution.


## prepareSupportDefaultBehavior


```lua
1?
```


Client supports the default behavior result.

The value indicates the default behavior used by the
client.



---

# lsp.RenameFile

Rename file operation

## annotationId


```lua
string?
```


An optional annotation identifier describing the operation.


## kind


```lua
"rename"
```


A rename

## newUri


```lua
string
```


The new location.

## oldUri


```lua
string
```


The old (existing) location.

## options


```lua
(lsp.RenameFileOptions)?
```


Rename options.


---

# lsp.RenameFileOptions

Rename file options

## ignoreIfExists


```lua
boolean?
```


Ignores if target exists.

## overwrite


```lua
boolean?
```


Overwrite target if existing. Overwrite wins over `ignoreIfExists`


---

# lsp.RenameFilesParams

The parameters sent in notifications/requests for user-initiated renames of
files.


## files


```lua
lsp.FileRename[]
```


An array of all files/folders renamed in this operation. When a folder is renamed, only
the folder will be included, and not its children.


---

# lsp.RenameOptions

Provider options for a {@link RenameRequest}.

## prepareProvider


```lua
boolean?
```


Renames should be checked and tested before being executed.


## workDoneProgress


```lua
boolean?
```



---

# lsp.RenameParams

The parameters of a {@link RenameRequest}.

## newName


```lua
string
```


The new name of the symbol. If the given name is not valid the
request must return a {@link ResponseError} with an
appropriate message set.

## position


```lua
lsp.Position
```


The position at which this request was sent.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document to rename.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.RenameRegistrationOptions

Registration options for a {@link RenameRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## prepareProvider


```lua
boolean?
```


Renames should be checked and tested before being executed.


## workDoneProgress


```lua
boolean?
```



---

# lsp.ResourceOperation

A generic resource operation.

## annotationId


```lua
string?
```


An optional annotation identifier describing the operation.


## kind


```lua
string
```


The resource operation kind.


---

# lsp.ResourceOperationKind


---

# lsp.ResponseError

## code


```lua
integer
```

## data


```lua
boolean|string|number|table|table[]...(+1)
```

## message


```lua
string
```


---

# lsp.SaveOptions

Save options.

## includeText


```lua
boolean?
```


The client is supposed to include the content on save.


---

# lsp.SelectedCompletionInfo

Describes the currently selected completion item.


## range


```lua
lsp.Range
```


The range that will be replaced if this completion item is accepted.

## text


```lua
string
```


The text the range will be replaced with if this completion is accepted.


---

# lsp.SelectionRange

A selection range represents a part of a selection hierarchy. A selection range
may have a parent selection range that contains it.

## parent


```lua
(lsp.SelectionRange)?
```


The parent selection range containing this range. Therefore `parent.range` must contain `this.range`.

## range


```lua
lsp.Range
```


The {@link Range range} of this selection range.


---

# lsp.SelectionRangeClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration for selection range providers. If this is set to `true`
the client supports the new `SelectionRangeRegistrationOptions` return value for the corresponding server
capability as well.


---

# lsp.SelectionRangeOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.SelectionRangeParams

A parameter literal used in selection range requests.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## positions


```lua
lsp.Position[]
```


The positions inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.SelectionRangeRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.SemanticTokenModifiers

A set of predefined token modifiers. This set is not fixed
an clients can specify additional token types via the
corresponding client capabilities.



---

# lsp.SemanticTokenTypes

A set of predefined token types. This set is not fixed
an clients can specify additional token types via the
corresponding client capabilities.



---

# lsp.SemanticTokens

## data


```lua
integer[]
```


The actual tokens.

## resultId


```lua
string?
```


An optional result id. If provided and clients support delta updating
the client will include the result id in the next semantic token request.
A server can then instead of computing all semantic tokens again simply
send a delta.


---

# lsp.SemanticTokensClientCapabilities

## augmentsSyntaxTokens


```lua
boolean?
```


Whether the client uses semantic tokens to augment existing
syntax tokens. If set to `true` client side created syntax
tokens and semantic tokens are both used for colorization. If
set to `false` the client only uses the returned semantic tokens
for colorization.

If the value is `undefined` then the client behavior is not
specified.


## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
return value for the corresponding server capability as well.

## formats


```lua
"relative"[]
```


The token formats the clients supports.

## multilineTokenSupport


```lua
boolean?
```


Whether the client supports tokens that can span multiple lines.

## overlappingTokenSupport


```lua
boolean?
```


Whether the client supports tokens that can overlap each other.

## requests


```lua
lsp.ClientSemanticTokensRequestOptions
```


Which requests the client supports and might send to the server
depending on the server's capability. Please note that clients might not
show semantic tokens or degrade some of the user experience if a range
or full request is advertised by the client but not provided by the
server. If for example the client capability `requests.full` and
`request.range` are both set to true but the server only provides a
range provider the client might not render a minimap correctly or might
even decide to not show any semantic tokens at all.

## serverCancelSupport


```lua
boolean?
```


Whether the client allows the server to actively cancel a
semantic token request, e.g. supports returning
LSPErrorCodes.ServerCancelled. If a server does the client
needs to retrigger the request.


## tokenModifiers


```lua
string[]
```


The token modifiers that the client supports.

## tokenTypes


```lua
string[]
```


The token types that the client supports.


---

# lsp.SemanticTokensDelta

## edits


```lua
lsp.SemanticTokensEdit[]
```


The semantic token edits to transform a previous result into a new result.

## resultId


```lua
string?
```



---

# lsp.SemanticTokensDeltaParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## previousResultId


```lua
string
```


The result id of a previous response. The result Id can either point to a full response
or a delta response depending on what was received last.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.SemanticTokensDeltaPartialResult

## edits


```lua
lsp.SemanticTokensEdit[]
```



---

# lsp.SemanticTokensEdit

## data


```lua
integer[]?
```


The elements to insert.

## deleteCount


```lua
integer
```


The count of elements to remove.

## start


```lua
integer
```


The start offset of the edit.


---

# lsp.SemanticTokensFullDelta

Semantic tokens options to support deltas for full documents


## delta


```lua
boolean?
```


The server supports deltas for full documents.


---

# lsp.SemanticTokensLegend

## tokenModifiers


```lua
string[]
```


The token modifiers a server uses.

## tokenTypes


```lua
string[]
```


The token types a server uses.


---

# lsp.SemanticTokensOptions

## full


```lua
(boolean|lsp.SemanticTokensFullDelta)?
```


Server supports providing semantic tokens for a full document.

## legend


```lua
lsp.SemanticTokensLegend
```


The legend used by the server

## range


```lua
(boolean|lsp._anonym1.range)?
```


Server supports providing semantic tokens for a specific range
of a document.

## workDoneProgress


```lua
boolean?
```



---

# lsp.SemanticTokensParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.SemanticTokensPartialResult

## data


```lua
integer[]
```



---

# lsp.SemanticTokensRangeParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## range


```lua
lsp.Range
```


The range the semantic tokens are requested for.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.SemanticTokensRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## full


```lua
(boolean|lsp.SemanticTokensFullDelta)?
```


Server supports providing semantic tokens for a full document.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## legend


```lua
lsp.SemanticTokensLegend
```


The legend used by the server

## range


```lua
(boolean|lsp._anonym1.range)?
```


Server supports providing semantic tokens for a specific range
of a document.

## workDoneProgress


```lua
boolean?
```



---

# lsp.SemanticTokensWorkspaceClientCapabilities

## refreshSupport


```lua
boolean?
```


Whether the client implementation supports a refresh request sent from
the server to the client.

Note that this event is global and will force the client to refresh all
semantic tokens currently shown. It should be used with absolute care
and is useful for situation where a server for example detects a project
wide change that requires such a calculation.


---

# lsp.ServerCapabilities

Defines the capabilities provided by a language
server.

## callHierarchyProvider


```lua
(boolean|lsp.CallHierarchyOptions|lsp.CallHierarchyRegistrationOptions)?
```


The server provides call hierarchy support.


## codeActionProvider


```lua
(boolean|lsp.CodeActionOptions)?
```


The server provides code actions. CodeActionOptions may only be
specified if the client states that it supports
`codeActionLiteralSupport` in its initial `initialize` request.

## codeLensProvider


```lua
(lsp.CodeLensOptions)?
```


The server provides code lens.

## colorProvider


```lua
(boolean|lsp.DocumentColorOptions|lsp.DocumentColorRegistrationOptions)?
```


The server provides color provider support.

## completionProvider


```lua
(lsp.CompletionOptions)?
```


The server provides completion support.

## declarationProvider


```lua
(boolean|lsp.DeclarationOptions|lsp.DeclarationRegistrationOptions)?
```


The server provides Goto Declaration support.

## definitionProvider


```lua
(boolean|lsp.DefinitionOptions)?
```


The server provides goto definition support.

## diagnosticProvider


```lua
(lsp.DiagnosticOptions|lsp.DiagnosticRegistrationOptions)?
```


The server has support for pull model diagnostics.


## documentFormattingProvider


```lua
(boolean|lsp.DocumentFormattingOptions)?
```


The server provides document formatting.

## documentHighlightProvider


```lua
(boolean|lsp.DocumentHighlightOptions)?
```


The server provides document highlight support.

## documentLinkProvider


```lua
(lsp.DocumentLinkOptions)?
```


The server provides document link support.

## documentOnTypeFormattingProvider


```lua
(lsp.DocumentOnTypeFormattingOptions)?
```


The server provides document formatting on typing.

## documentRangeFormattingProvider


```lua
(boolean|lsp.DocumentRangeFormattingOptions)?
```


The server provides document range formatting.

## documentSymbolProvider


```lua
(boolean|lsp.DocumentSymbolOptions)?
```


The server provides document symbol support.

## executeCommandProvider


```lua
(lsp.ExecuteCommandOptions)?
```


The server provides execute command support.

## experimental


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


Experimental server capabilities.

## foldingRangeProvider


```lua
(boolean|lsp.FoldingRangeOptions|lsp.FoldingRangeRegistrationOptions)?
```


The server provides folding provider support.

## hoverProvider


```lua
(boolean|lsp.HoverOptions)?
```


The server provides hover support.

## implementationProvider


```lua
(boolean|lsp.ImplementationOptions|lsp.ImplementationRegistrationOptions)?
```


The server provides Goto Implementation support.

## inlayHintProvider


```lua
(boolean|lsp.InlayHintOptions|lsp.InlayHintRegistrationOptions)?
```


The server provides inlay hints.


## inlineCompletionProvider


```lua
(boolean|lsp.InlineCompletionOptions)?
```


Inline completion options used during static registration.


## inlineValueProvider


```lua
(boolean|lsp.InlineValueOptions|lsp.InlineValueRegistrationOptions)?
```


The server provides inline values.


## linkedEditingRangeProvider


```lua
(boolean|lsp.LinkedEditingRangeOptions|lsp.LinkedEditingRangeRegistrationOptions)?
```


The server provides linked editing range support.


## monikerProvider


```lua
(boolean|lsp.MonikerOptions|lsp.MonikerRegistrationOptions)?
```


The server provides moniker support.


## notebookDocumentSync


```lua
(lsp.NotebookDocumentSyncOptions|lsp.NotebookDocumentSyncRegistrationOptions)?
```


Defines how notebook documents are synced.


## positionEncoding


```lua
("utf-16"|"utf-32"|"utf-8")?
```


The position encoding the server picked from the encodings offered
by the client via the client capability `general.positionEncodings`.

If the client didn't provide any position encodings the only valid
value that a server can return is 'utf-16'.

If omitted it defaults to 'utf-16'.


## referencesProvider


```lua
(boolean|lsp.ReferenceOptions)?
```


The server provides find references support.

## renameProvider


```lua
(boolean|lsp.RenameOptions)?
```


The server provides rename support. RenameOptions may only be
specified if the client states that it supports
`prepareSupport` in its initial `initialize` request.

## selectionRangeProvider


```lua
(boolean|lsp.SelectionRangeOptions|lsp.SelectionRangeRegistrationOptions)?
```


The server provides selection range support.

## semanticTokensProvider


```lua
(lsp.SemanticTokensOptions|lsp.SemanticTokensRegistrationOptions)?
```


The server provides semantic tokens support.


## signatureHelpProvider


```lua
(lsp.SignatureHelpOptions)?
```


The server provides signature help support.

## textDocumentSync


```lua
(0|1|2|lsp.TextDocumentSyncOptions)?
```


Defines how text documents are synced. Is either a detailed structure
defining each notification or for backwards compatibility the
TextDocumentSyncKind number.

## typeDefinitionProvider


```lua
(boolean|lsp.TypeDefinitionOptions|lsp.TypeDefinitionRegistrationOptions)?
```


The server provides Goto Type Definition support.

## typeHierarchyProvider


```lua
(boolean|lsp.TypeHierarchyOptions|lsp.TypeHierarchyRegistrationOptions)?
```


The server provides type hierarchy support.


## workspace


```lua
(lsp.WorkspaceOptions)?
```


Workspace specific server capabilities.

## workspaceSymbolProvider


```lua
(boolean|lsp.WorkspaceSymbolOptions)?
```


The server provides workspace symbol support.


---

# lsp.ServerCompletionItemOptions

## labelDetailsSupport


```lua
boolean?
```


The server has support for completion item label
details (see also `CompletionItemLabelDetails`) when
receiving a completion item in a resolve call.



---

# lsp.ServerInfo

Information about the server


## name


```lua
string
```


The name of the server as defined by the server.

## version


```lua
string?
```


The server's version as defined by the server.


---

# lsp.SetTraceParams

## value


```lua
"messages"|"off"|"verbose"
```



---

# lsp.ShowDocumentClientCapabilities

Client capabilities for the showDocument request.


## support


```lua
boolean
```


The client has support for the showDocument
request.


---

# lsp.ShowDocumentParams

Params to show a resource in the UI.


## external


```lua
boolean?
```


Indicates to show the resource in an external program.
To show, for example, `https://code.visualstudio.com/`
in the default WEB browser set `external` to `true`.

## selection


```lua
(lsp.Range)?
```


An optional selection range if the document is a text
document. Clients might ignore the property if an
external program is started or the file is not a text
file.

## takeFocus


```lua
boolean?
```


An optional property to indicate whether the editor
showing the document should take focus or not.
Clients might ignore this property if an external
program is started.

## uri


```lua
string
```


The uri to show.


---

# lsp.ShowDocumentResult

The result of a showDocument request.


## success


```lua
boolean
```


A boolean indicating if the show was successful.


---

# lsp.ShowMessageParams

The parameters of a notification message.

## message


```lua
string
```


The actual message.

## type


```lua
1|2|3|4|5
```


The message type. See {@link MessageType}


---

# lsp.ShowMessageRequestClientCapabilities

Show message request client capabilities

## messageActionItem


```lua
(lsp.ClientShowMessageActionItemOptions)?
```


Capabilities specific to the `MessageActionItem` type.


---

# lsp.ShowMessageRequestParams

## actions


```lua
lsp.MessageActionItem[]?
```


The message action items to present.

## message


```lua
string
```


The actual message.

## type


```lua
1|2|3|4|5
```


The message type. See {@link MessageType}


---

# lsp.SignatureHelp

Signature help represents the signature of something
callable. There can be multiple signature but only one
active and only one active parameter.

## activeParameter


```lua
(integer|nil)?
```


The active parameter of the active signature.

If `null`, no parameter of the signature is active (for example a named
argument that does not match any declared parameters). This is only valid
if the client specifies the client capability
`textDocument.signatureHelp.noActiveParameterSupport === true`

If omitted or the value lies outside the range of
`signatures[activeSignature].parameters` defaults to 0 if the active
signature has parameters.

If the active signature has no parameters it is ignored.

In future version of the protocol this property might become
mandatory (but still nullable) to better express the active parameter if
the active signature does have any.

## activeSignature


```lua
integer?
```


The active signature. If omitted or the value lies outside the
range of `signatures` the value defaults to zero or is ignored if
the `SignatureHelp` has no signatures.

Whenever possible implementors should make an active decision about
the active signature and shouldn't rely on a default value.

In future version of the protocol this property might become
mandatory to better express this.

## signatures


```lua
lsp.SignatureInformation[]
```


One or more signatures.


---

# lsp.SignatureHelpClientCapabilities

Client Capabilities for a {@link SignatureHelpRequest}.

## contextSupport


```lua
boolean?
```


The client supports to send additional context information for a
`textDocument/signatureHelp` request. A client that opts into
contextSupport will also support the `retriggerCharacters` on
`SignatureHelpOptions`.


## dynamicRegistration


```lua
boolean?
```


Whether signature help supports dynamic registration.

## signatureInformation


```lua
(lsp.ClientSignatureInformationOptions)?
```


The client supports the following `SignatureInformation`
specific properties.


---

# lsp.SignatureHelpContext

Additional information about the context in which a signature help request was triggered.


## activeSignatureHelp


```lua
(lsp.SignatureHelp)?
```


The currently active `SignatureHelp`.

The `activeSignatureHelp` has its `SignatureHelp.activeSignature` field updated based on
the user navigating through available signatures.

## isRetrigger


```lua
boolean
```


`true` if signature help was already showing when it was triggered.

Retriggers occurs when the signature help is already active and can be caused by actions such as
typing a trigger character, a cursor move, or document content changes.

## triggerCharacter


```lua
string?
```


Character that caused signature help to be triggered.

This is undefined when `triggerKind !== SignatureHelpTriggerKind.TriggerCharacter`

## triggerKind


```lua
1|2|3
```


Action that caused signature help to be triggered.


---

# lsp.SignatureHelpOptions

Server Capabilities for a {@link SignatureHelpRequest}.

## retriggerCharacters


```lua
string[]?
```


List of characters that re-trigger signature help.

These trigger characters are only active when signature help is already showing. All trigger characters
are also counted as re-trigger characters.


## triggerCharacters


```lua
string[]?
```


List of characters that trigger signature help automatically.

## workDoneProgress


```lua
boolean?
```



---

# lsp.SignatureHelpParams

Parameters for a {@link SignatureHelpRequest}.

## context


```lua
(lsp.SignatureHelpContext)?
```


The signature help context. This is only available if the client specifies
to send this using the client capability `textDocument.signatureHelp.contextSupport === true`


## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.SignatureHelpRegistrationOptions

Registration options for a {@link SignatureHelpRequest}.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## retriggerCharacters


```lua
string[]?
```


List of characters that re-trigger signature help.

These trigger characters are only active when signature help is already showing. All trigger characters
are also counted as re-trigger characters.


## triggerCharacters


```lua
string[]?
```


List of characters that trigger signature help automatically.

## workDoneProgress


```lua
boolean?
```



---

# lsp.SignatureHelpTriggerKind

How a signature help was triggered.



---

# lsp.SignatureInformation

Represents the signature of something callable. A signature
can have a label, like a function-name, a doc-comment, and
a set of parameters.

## activeParameter


```lua
(integer|nil)?
```


The index of the active parameter.

If `null`, no parameter of the signature is active (for example a named
argument that does not match any declared parameters). This is only valid
if the client specifies the client capability
`textDocument.signatureHelp.noActiveParameterSupport === true`

If provided (or `null`), this is used in place of
`SignatureHelp.activeParameter`.


## documentation


```lua
(string|lsp.MarkupContent)?
```


The human-readable doc-comment of this signature. Will be shown
in the UI but can be omitted.

## label


```lua
string
```


The label of this signature. Will be shown in
the UI.

## parameters


```lua
lsp.ParameterInformation[]?
```


The parameters of this signature.


---

# lsp.SnippetTextEdit

An interactive text edit.


## annotationId


```lua
string?
```


The actual identifier of the snippet edit.

## range


```lua
lsp.Range
```


The range of the text document to be manipulated.

## snippet


```lua
lsp.StringValue
```


The snippet to be inserted.


---

# lsp.StaleRequestSupportOptions

## cancel


```lua
boolean
```


The client will actively cancel the request.

## retryOnContentModified


```lua
string[]
```


The list of requests for which the client
will retry the request if it receives a
response with error code `ContentModified`


---

# lsp.StaticRegistrationOptions

Static registration options to be returned in the initialize
request.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.


---

# lsp.StringValue

A string value used as a snippet is a template which allows to insert text
and to control the editor cursor when insertion happens.

A snippet can define tab stops and placeholders with `$1`, `$2`
and `${3:foo}`. `$0` defines the final tab stop, it defaults to
the end of the snippet. Variables are defined with `$name` and
`${name:default value}`.


## kind


```lua
"snippet"
```


The kind of string value.

## value


```lua
string
```


The snippet string.


---

# lsp.SymbolInformation

Represents information about programming constructs like variables, classes,
interfaces etc.

## containerName


```lua
string?
```


The name of the symbol containing this symbol. This information is for
user interface purposes (e.g. to render a qualifier in the user interface
if necessary). It can't be used to re-infer a hierarchy for the document
symbols.

## deprecated


```lua
boolean?
```


Indicates if this symbol is deprecated.


## kind


```lua
1|10|11|12|13...(+21)
```


The kind of this symbol.

## location


```lua
lsp.Location
```


The location of this symbol. The location's range is used by a tool
to reveal the location in the editor. If the symbol is selected in the
tool the range's start information is used to position the cursor. So
the range usually spans more than the actual symbol's name and does
normally include things like visibility modifiers.

The range doesn't have to denote a node range in the sense of an abstract
syntax tree. It can therefore not be used to re-construct a hierarchy of
the symbols.

## name


```lua
string
```


The name of this symbol.

## tags


```lua
1[]?
```


Tags for this symbol.



---

# lsp.SymbolKind

A symbol kind.


---

# lsp.SymbolTag

Symbol tags are extra annotations that tweak the rendering of a symbol.



---

# lsp.TextDocumentChangeRegistrationOptions

Describe options to be used when registered for text document change events.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## syncKind


```lua
0|1|2
```


How documents are synced to the server.


---

# lsp.TextDocumentClientCapabilities

Text document specific client capabilities.

## callHierarchy


```lua
(lsp.CallHierarchyClientCapabilities)?
```


Capabilities specific to the various call hierarchy requests.


## codeAction


```lua
(lsp.CodeActionClientCapabilities)?
```


Capabilities specific to the `textDocument/codeAction` request.

## codeLens


```lua
(lsp.CodeLensClientCapabilities)?
```


Capabilities specific to the `textDocument/codeLens` request.

## colorProvider


```lua
(lsp.DocumentColorClientCapabilities)?
```


Capabilities specific to the `textDocument/documentColor` and the
`textDocument/colorPresentation` request.


## completion


```lua
(lsp.CompletionClientCapabilities)?
```


Capabilities specific to the `textDocument/completion` request.

## declaration


```lua
(lsp.DeclarationClientCapabilities)?
```


Capabilities specific to the `textDocument/declaration` request.


## definition


```lua
(lsp.DefinitionClientCapabilities)?
```


Capabilities specific to the `textDocument/definition` request.

## diagnostic


```lua
(lsp.DiagnosticClientCapabilities)?
```


Capabilities specific to the diagnostic pull model.


## documentHighlight


```lua
(lsp.DocumentHighlightClientCapabilities)?
```


Capabilities specific to the `textDocument/documentHighlight` request.

## documentLink


```lua
(lsp.DocumentLinkClientCapabilities)?
```


Capabilities specific to the `textDocument/documentLink` request.

## documentSymbol


```lua
(lsp.DocumentSymbolClientCapabilities)?
```


Capabilities specific to the `textDocument/documentSymbol` request.

## filters


```lua
(lsp.TextDocumentFilterClientCapabilities)?
```


Defines which filters the client supports.


## foldingRange


```lua
(lsp.FoldingRangeClientCapabilities)?
```


Capabilities specific to the `textDocument/foldingRange` request.


## formatting


```lua
(lsp.DocumentFormattingClientCapabilities)?
```


Capabilities specific to the `textDocument/formatting` request.

## hover


```lua
(lsp.HoverClientCapabilities)?
```


Capabilities specific to the `textDocument/hover` request.

## implementation


```lua
(lsp.ImplementationClientCapabilities)?
```


Capabilities specific to the `textDocument/implementation` request.


## inlayHint


```lua
(lsp.InlayHintClientCapabilities)?
```


Capabilities specific to the `textDocument/inlayHint` request.


## inlineCompletion


```lua
(lsp.InlineCompletionClientCapabilities)?
```


Client capabilities specific to inline completions.


## inlineValue


```lua
(lsp.InlineValueClientCapabilities)?
```


Capabilities specific to the `textDocument/inlineValue` request.


## linkedEditingRange


```lua
(lsp.LinkedEditingRangeClientCapabilities)?
```


Capabilities specific to the `textDocument/linkedEditingRange` request.


## moniker


```lua
(lsp.MonikerClientCapabilities)?
```


Client capabilities specific to the `textDocument/moniker` request.


## onTypeFormatting


```lua
(lsp.DocumentOnTypeFormattingClientCapabilities)?
```


Capabilities specific to the `textDocument/onTypeFormatting` request.

## publishDiagnostics


```lua
(lsp.PublishDiagnosticsClientCapabilities)?
```


Capabilities specific to the `textDocument/publishDiagnostics` notification.

## rangeFormatting


```lua
(lsp.DocumentRangeFormattingClientCapabilities)?
```


Capabilities specific to the `textDocument/rangeFormatting` request.

## references


```lua
(lsp.ReferenceClientCapabilities)?
```


Capabilities specific to the `textDocument/references` request.

## rename


```lua
(lsp.RenameClientCapabilities)?
```


Capabilities specific to the `textDocument/rename` request.

## selectionRange


```lua
(lsp.SelectionRangeClientCapabilities)?
```


Capabilities specific to the `textDocument/selectionRange` request.


## semanticTokens


```lua
(lsp.SemanticTokensClientCapabilities)?
```


Capabilities specific to the various semantic token request.


## signatureHelp


```lua
(lsp.SignatureHelpClientCapabilities)?
```


Capabilities specific to the `textDocument/signatureHelp` request.

## synchronization


```lua
(lsp.TextDocumentSyncClientCapabilities)?
```


Defines which synchronization capabilities the client supports.

## typeDefinition


```lua
(lsp.TypeDefinitionClientCapabilities)?
```


Capabilities specific to the `textDocument/typeDefinition` request.


## typeHierarchy


```lua
(lsp.TypeHierarchyClientCapabilities)?
```


Capabilities specific to the various type hierarchy requests.



---

# lsp.TextDocumentContentChangeEvent

An event describing a change to a text document. If only a text is provided
it is considered to be the full content of the document.


---

# lsp.TextDocumentContentChangePartial

## range


```lua
lsp.Range
```


The range of the document that changed.

## rangeLength


```lua
integer?
```


The optional length of the range that got replaced.


## text


```lua
string
```


The new text for the provided range.


---

# lsp.TextDocumentContentChangeWholeDocument

## text


```lua
string
```


The new text of the whole document.


---

# lsp.TextDocumentContentClientCapabilities

Client capabilities for a text document content provider.


## dynamicRegistration


```lua
boolean?
```


Text document content provider supports dynamic registration.


---

# lsp.TextDocumentContentOptions

Text document content provider options.


## schemes


```lua
string[]
```


The schemes for which the server provides content.


---

# lsp.TextDocumentContentParams

Parameters for the `workspace/textDocumentContent` request.


## uri


```lua
string
```


The uri of the text document.


---

# lsp.TextDocumentContentRefreshParams

Parameters for the `workspace/textDocumentContent/refresh` request.


## uri


```lua
string
```


The uri of the text document to refresh.


---

# lsp.TextDocumentContentRegistrationOptions

Text document content provider registration options.


## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## schemes


```lua
string[]
```


The schemes for which the server provides content.


---

# lsp.TextDocumentContentResult

Result of the `workspace/textDocumentContent` request.


## text


```lua
string
```


The text content of the text document. Please note, that the content of
any subsequent open notifications for the text document might differ
from the returned content due to whitespace and line ending
normalizations done on the client


---

# lsp.TextDocumentEdit

Describes textual changes on a text document. A TextDocumentEdit describes all changes
on a document version Si and after they are applied move the document to version Si+1.
So the creator of a TextDocumentEdit doesn't need to sort the array of edits or do any
kind of ordering. However the edits must be non overlapping.

## edits


```lua
(lsp.AnnotatedTextEdit|lsp.SnippetTextEdit|lsp.TextEdit)[]
```


The edits to be applied.

client capability.

client capability.

## textDocument


```lua
lsp.OptionalVersionedTextDocumentIdentifier
```


The text document to change.


---

# lsp.TextDocumentFilter

A document filter denotes a document by different properties like
the {@link TextDocument.languageId language}, the {@link Uri.scheme scheme} of
its resource, or a glob-pattern that is applied to the {@link TextDocument.fileName path}.

Glob patterns can have the following syntax:
- `*` to match one or more characters in a path segment
- `?` to match on one character in a path segment
- `**` to match any number of path segments, including none
- `{}` to group sub patterns into an OR expression. (e.g. `**/*.{ts,js}` matches all TypeScript and JavaScript files)
- `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
- `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)

\@sample A language filter that applies to typescript files on disk: `{ language: 'typescript', scheme: 'file' }`
\@sample A language filter that applies to all package.json paths: `{ language: 'json', pattern: '**package.json' }`



---

# lsp.TextDocumentFilterClientCapabilities

## relativePatternSupport


```lua
boolean?
```


The client supports Relative Patterns.



---

# lsp.TextDocumentFilterLanguage

A document filter where `language` is required field.


## language


```lua
string
```


A language id, like `typescript`.

## pattern


```lua
(string|lsp.RelativePattern)?
```


A glob pattern, like **/*.{ts,js}. See TextDocumentFilter for examples.

relative patterns depends on the client capability
`textDocuments.filters.relativePatternSupport`.

## scheme


```lua
string?
```


A Uri {@link Uri.scheme scheme}, like `file` or `untitled`.


---

# lsp.TextDocumentFilterPattern

A document filter where `pattern` is required field.


## language


```lua
string?
```


A language id, like `typescript`.

## pattern


```lua
string|lsp.RelativePattern
```


A glob pattern, like **/*.{ts,js}. See TextDocumentFilter for examples.

relative patterns depends on the client capability
`textDocuments.filters.relativePatternSupport`.

## scheme


```lua
string?
```


A Uri {@link Uri.scheme scheme}, like `file` or `untitled`.


---

# lsp.TextDocumentFilterScheme

A document filter where `scheme` is required field.


## language


```lua
string?
```


A language id, like `typescript`.

## pattern


```lua
(string|lsp.RelativePattern)?
```


A glob pattern, like **/*.{ts,js}. See TextDocumentFilter for examples.

relative patterns depends on the client capability
`textDocuments.filters.relativePatternSupport`.

## scheme


```lua
string
```


A Uri {@link Uri.scheme scheme}, like `file` or `untitled`.


---

# lsp.TextDocumentIdentifier

A literal to identify a text document in the client.

## uri


```lua
string
```


The text document's uri.


---

# lsp.TextDocumentItem

An item to transfer a text document from the client to the
server.

## languageId


```lua
"abap"|"bat"|"bibtex"|"c"|"clojure"...(+55)
```


The text document's language identifier.

## text


```lua
string
```


The content of the opened text document.

## uri


```lua
string
```


The text document's uri.

## version


```lua
integer
```


The version number of this document (it will increase after each
change, including undo/redo).


---

# lsp.TextDocumentPositionParams

A parameter literal used in requests to pass a text document and a position inside that
document.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.


---

# lsp.TextDocumentRegistrationOptions

General text document registration options.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.


---

# lsp.TextDocumentSaveReason

Represents reasons why a text document is saved.


---

# lsp.TextDocumentSaveRegistrationOptions

Save registration options.

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## includeText


```lua
boolean?
```


The client is supposed to include the content on save.


---

# lsp.TextDocumentSyncClientCapabilities

## didSave


```lua
boolean?
```


The client supports did save notifications.

## dynamicRegistration


```lua
boolean?
```


Whether text document synchronization supports dynamic registration.

## willSave


```lua
boolean?
```


The client supports sending will save notifications.

## willSaveWaitUntil


```lua
boolean?
```


The client supports sending a will save request and
waits for a response providing text edits which will
be applied to the document before it is saved.


---

# lsp.TextDocumentSyncKind

Defines how the host (editor) should sync
document changes to the language server.


---

# lsp.TextDocumentSyncOptions

## change


```lua
(0|1|2)?
```


Change notifications are sent to the server. See TextDocumentSyncKind.None, TextDocumentSyncKind.Full
and TextDocumentSyncKind.Incremental. If omitted it defaults to TextDocumentSyncKind.None.

## openClose


```lua
boolean?
```


Open and close notifications are sent to the server. If omitted open close notification should not
be sent.

## save


```lua
(boolean|lsp.SaveOptions)?
```


If present save notifications are sent to the server. If omitted the notification should not be
sent.

## willSave


```lua
boolean?
```


If present will save notifications are sent to the server. If omitted the notification should not be
sent.

## willSaveWaitUntil


```lua
boolean?
```


If present will save wait until requests are sent to the server. If omitted the request should not be
sent.


---

# lsp.TextEdit

A text edit applicable to a text document.

## newText


```lua
string
```


The string to be inserted. For delete operations use an
empty string.

## range


```lua
lsp.Range
```


The range of the text document to be manipulated. To insert
text into a document create a range where start === end.


---

# lsp.TokenFormat


---

# lsp.TraceValue


---

# lsp.TypeDefinitionClientCapabilities

Since 3.6.0

## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `TypeDefinitionRegistrationOptions` return value
for the corresponding server capability as well.

## linkSupport


```lua
boolean?
```


The client supports additional metadata in the form of definition links.

Since 3.14.0


---

# lsp.TypeDefinitionOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.TypeDefinitionParams

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.TypeDefinitionRegistrationOptions

## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.TypeHierarchyClientCapabilities

## dynamicRegistration


```lua
boolean?
```


Whether implementation supports dynamic registration. If this is set to `true`
the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
return value for the corresponding server capability as well.


---

# lsp.TypeHierarchyItem

## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved between a type hierarchy prepare and
supertypes or subtypes requests. It could also be used to identify the
type hierarchy in the server, helping improve the performance on
resolving supertypes and subtypes.

## detail


```lua
string?
```


More detail for this item, e.g. the signature of a function.

## kind


```lua
1|10|11|12|13...(+21)
```


The kind of this item.

## name


```lua
string
```


The name of this item.

## range


```lua
lsp.Range
```


The range enclosing this symbol not including leading/trailing whitespace
but everything else, e.g. comments and code.

## selectionRange


```lua
lsp.Range
```


The range that should be selected and revealed when this symbol is being
picked, e.g. the name of a function. Must be contained by the
{@link TypeHierarchyItem.range `range`}.

## tags


```lua
1[]?
```


Tags for this item.

## uri


```lua
string
```


The resource identifier of this item.


---

# lsp.TypeHierarchyOptions

Type hierarchy options used during static registration.


## workDoneProgress


```lua
boolean?
```



---

# lsp.TypeHierarchyPrepareParams

The parameter of a `textDocument/prepareTypeHierarchy` request.


## position


```lua
lsp.Position
```


The position inside the text document.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The text document.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.TypeHierarchyRegistrationOptions

Type hierarchy options used during static or dynamic registration.


## documentSelector


```lua
lsp.NotebookCellTextDocumentFilter|lsp.TextDocumentFilterLanguage|lsp.TextDocumentFilterPattern|lsp.TextDocumentFilterScheme[]|nil
```


A document selector to identify the scope of the registration. If set to null
the document selector provided on the client side will be used.

## id


```lua
string?
```


The id used to register the request. The id can be used to deregister
the request again. See also Registration#id.

## workDoneProgress


```lua
boolean?
```



---

# lsp.TypeHierarchySubtypesParams

The parameter of a `typeHierarchy/subtypes` request.


## item


```lua
lsp.TypeHierarchyItem
```


## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.TypeHierarchySupertypesParams

The parameter of a `typeHierarchy/supertypes` request.


## item


```lua
lsp.TypeHierarchyItem
```


## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.URI


---

# lsp.UnchangedDocumentDiagnosticReport

A diagnostic report indicating that the last returned
report is still accurate.


## kind


```lua
"unchanged"
```


A document diagnostic report indicating
no changes to the last result. A server can
only return `unchanged` if result ids are
provided.

## resultId


```lua
string
```


A result id which will be sent on the next
diagnostic request for the same document.


---

# lsp.UniquenessLevel

Moniker uniqueness level to define scope of the moniker.



---

# lsp.Unregistration

General parameters to unregister a request or notification.

## id


```lua
string
```


The id used to unregister the request or notification. Usually an id
provided during the register request.

## method


```lua
string
```


The method to unregister for.


---

# lsp.UnregistrationParams

## unregisterations


```lua
lsp.Unregistration[]
```



---

# lsp.VersionedNotebookDocumentIdentifier

A versioned notebook document identifier.


## uri


```lua
string
```


The notebook document's uri.

## version


```lua
integer
```


The version number of this notebook document.


---

# lsp.VersionedTextDocumentIdentifier

A text document identifier to denote a specific version of a text document.

## uri


```lua
string
```


The text document's uri.

## version


```lua
integer
```


The version number of this document.


---

# lsp.WatchKind


---

# lsp.WillSaveTextDocumentParams

The parameters sent in a will save text document notification.

## reason


```lua
1|2|3
```


The 'TextDocumentSaveReason'.

## textDocument


```lua
lsp.TextDocumentIdentifier
```


The document that will be saved.


---

# lsp.WindowClientCapabilities

## showDocument


```lua
(lsp.ShowDocumentClientCapabilities)?
```


Capabilities specific to the showDocument request.


## showMessage


```lua
(lsp.ShowMessageRequestClientCapabilities)?
```


Capabilities specific to the showMessage request.


## workDoneProgress


```lua
boolean?
```


It indicates whether the client supports server initiated
progress using the `window/workDoneProgress/create` request.

The capability also controls Whether client supports handling
of progress notifications. If set servers are allowed to report a
`workDoneProgress` property in the request specific server
capabilities.



---

# lsp.WorkDoneProgressBegin

## cancellable


```lua
boolean?
```


Controls if a cancel button should show to allow the user to cancel the
long running operation. Clients that don't support cancellation are allowed
to ignore the setting.

## kind


```lua
"begin"
```


## message


```lua
string?
```


Optional, more detailed associated progress message. Contains
complementary information to the `title`.

Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
If unset, the previous progress message (if any) is still valid.

## percentage


```lua
integer?
```


Optional progress percentage to display (value 100 is considered 100%).
If not provided infinite progress is assumed and clients are allowed
to ignore the `percentage` value in subsequent in report notifications.

The value should be steadily rising. Clients are free to ignore values
that are not following this rule. The value range is [0, 100].

## title


```lua
string
```


Mandatory title of the progress operation. Used to briefly inform about
the kind of operation being performed.

Examples: "Indexing" or "Linking dependencies".


---

# lsp.WorkDoneProgressCancelParams

## token


```lua
string|integer
```


The token to be used to report progress.


---

# lsp.WorkDoneProgressCreateParams

## token


```lua
string|integer
```


The token to be used to report progress.


---

# lsp.WorkDoneProgressEnd

## kind


```lua
"end"
```


## message


```lua
string?
```


Optional, a final message indicating to for example indicate the outcome
of the operation.


---

# lsp.WorkDoneProgressOptions

## workDoneProgress


```lua
boolean?
```



---

# lsp.WorkDoneProgressParams

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.WorkDoneProgressReport

## cancellable


```lua
boolean?
```


Controls enablement state of a cancel button.

Clients that don't support cancellation or don't support controlling the button's
enablement state are allowed to ignore the property.

## kind


```lua
"report"
```


## message


```lua
string?
```


Optional, more detailed associated progress message. Contains
complementary information to the `title`.

Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
If unset, the previous progress message (if any) is still valid.

## percentage


```lua
integer?
```


Optional progress percentage to display (value 100 is considered 100%).
If not provided infinite progress is assumed and clients are allowed
to ignore the `percentage` value in subsequent in report notifications.

The value should be steadily rising. Clients are free to ignore values
that are not following this rule. The value range is [0, 100]


---

# lsp.WorkspaceClientCapabilities

Workspace specific client capabilities.

## applyEdit


```lua
boolean?
```


The client supports applying batch edits
to the workspace by supporting the request
'workspace/applyEdit'

## codeLens


```lua
(lsp.CodeLensWorkspaceClientCapabilities)?
```


Capabilities specific to the code lens requests scoped to the
workspace.


## configuration


```lua
boolean?
```


The client supports `workspace/configuration` requests.


## diagnostics


```lua
(lsp.DiagnosticWorkspaceClientCapabilities)?
```


Capabilities specific to the diagnostic requests scoped to the
workspace.


## didChangeConfiguration


```lua
(lsp.DidChangeConfigurationClientCapabilities)?
```


Capabilities specific to the `workspace/didChangeConfiguration` notification.

## didChangeWatchedFiles


```lua
(lsp.DidChangeWatchedFilesClientCapabilities)?
```


Capabilities specific to the `workspace/didChangeWatchedFiles` notification.

## executeCommand


```lua
(lsp.ExecuteCommandClientCapabilities)?
```


Capabilities specific to the `workspace/executeCommand` request.

## fileOperations


```lua
(lsp.FileOperationClientCapabilities)?
```


The client has support for file notifications/requests for user operations on files.

Since 3.16.0

## foldingRange


```lua
(lsp.FoldingRangeWorkspaceClientCapabilities)?
```


Capabilities specific to the folding range requests scoped to the workspace.


## inlayHint


```lua
(lsp.InlayHintWorkspaceClientCapabilities)?
```


Capabilities specific to the inlay hint requests scoped to the
workspace.


## inlineValue


```lua
(lsp.InlineValueWorkspaceClientCapabilities)?
```


Capabilities specific to the inline values requests scoped to the
workspace.


## semanticTokens


```lua
(lsp.SemanticTokensWorkspaceClientCapabilities)?
```


Capabilities specific to the semantic token requests scoped to the
workspace.


## symbol


```lua
(lsp.WorkspaceSymbolClientCapabilities)?
```


Capabilities specific to the `workspace/symbol` request.

## textDocumentContent


```lua
(lsp.TextDocumentContentClientCapabilities)?
```


Capabilities specific to the `workspace/textDocumentContent` request.


## workspaceEdit


```lua
(lsp.WorkspaceEditClientCapabilities)?
```


Capabilities specific to `WorkspaceEdit`s.

## workspaceFolders


```lua
boolean?
```


The client has support for workspace folders.



---

# lsp.WorkspaceDiagnosticParams

Parameters of the workspace diagnostic request.


## identifier


```lua
string?
```


The additional identifier provided during registration.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## previousResultIds


```lua
lsp.PreviousResultId[]
```


The currently known diagnostic reports with their
previous result ids.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.WorkspaceDiagnosticReport

A workspace diagnostic report.


## items


```lua
lsp.WorkspaceFullDocumentDiagnosticReport|lsp.WorkspaceUnchangedDocumentDiagnosticReport[]
```



---

# lsp.WorkspaceDiagnosticReportPartialResult

A partial result for a workspace diagnostic report.


## items


```lua
lsp.WorkspaceFullDocumentDiagnosticReport|lsp.WorkspaceUnchangedDocumentDiagnosticReport[]
```



---

# lsp.WorkspaceDocumentDiagnosticReport

A workspace diagnostic document report.



---

# lsp.WorkspaceEdit

A workspace edit represents changes to many resources managed in the workspace. The edit
should either provide `changes` or `documentChanges`. If documentChanges are present
they are preferred over `changes` if the client can handle versioned document edits.

Since version 3.13.0 a workspace edit can contain resource operations as well. If resource
operations are present clients need to execute the operations in the order in which they
are provided. So a workspace edit for example can consist of the following two changes:
(1) a create file a.txt and (2) a text document edit which insert text into file a.txt.

An invalid sequence (e.g. (1) delete file a.txt and (2) insert text into file a.txt) will
cause failure of the operation. How the client recovers from the failure is described by
the client capability: `workspace.workspaceEdit.failureHandling`

## changeAnnotations


```lua
table<string, lsp.ChangeAnnotation>?
```


A map of change annotations that can be referenced in `AnnotatedTextEdit`s or create, rename and
delete file / folder operations.

Whether clients honor this property depends on the client capability `workspace.changeAnnotationSupport`.


## changes


```lua
table<string, lsp.TextEdit[]>?
```


Holds changes to existing resources.

## documentChanges


```lua
(lsp.CreateFile|lsp.DeleteFile|lsp.RenameFile|lsp.TextDocumentEdit)[]?
```


Depending on the client capability `workspace.workspaceEdit.resourceOperations` document changes
are either an array of `TextDocumentEdit`s to express changes to n different text documents
where each text document edit addresses a specific version of a text document. Or it can contain
above `TextDocumentEdit`s mixed with create, rename and delete file / folder operations.

Whether a client supports versioned document edits is expressed via
`workspace.workspaceEdit.documentChanges` client capability.

If a client neither supports `documentChanges` nor `workspace.workspaceEdit.resourceOperations` then
only plain `TextEdit`s using the `changes` property are supported.


---

# lsp.WorkspaceEditClientCapabilities

## changeAnnotationSupport


```lua
(lsp.ChangeAnnotationsSupportOptions)?
```


Whether the client in general supports change annotations on text edits,
create file, rename file and delete file changes.


## documentChanges


```lua
boolean?
```


The client supports versioned document changes in `WorkspaceEdit`s

## failureHandling


```lua
("abort"|"textOnlyTransactional"|"transactional"|"undo")?
```


The failure handling strategy of a client if applying the workspace edit
fails.


## metadataSupport


```lua
boolean?
```


Whether the client supports `WorkspaceEditMetadata` in `WorkspaceEdit`s.


## normalizesLineEndings


```lua
boolean?
```


Whether the client normalizes line endings to the client specific
setting.
If set to `true` the client will normalize line ending characters
in a workspace edit to the client-specified new line
character.


## resourceOperations


```lua
"create"|"delete"|"rename"[]?
```


The resource operations the client supports. Clients should at least
support 'create', 'rename' and 'delete' files and folders.


## snippetEditSupport


```lua
boolean?
```


Whether the client supports snippets as text edits.



---

# lsp.WorkspaceEditMetadata

Additional data about a workspace edit.


## isRefactoring


```lua
boolean?
```


Signal to the editor that this edit is a refactoring.


---

# lsp.WorkspaceFolder

A workspace folder inside a client.

## name


```lua
string
```


The name of the workspace folder. Used to refer to this
workspace folder in the user interface.

## uri


```lua
string
```


The associated URI for this workspace folder.


---

# lsp.WorkspaceFoldersChangeEvent

The workspace folder change event.

## added


```lua
lsp.WorkspaceFolder[]
```


The array of added workspace folders

## removed


```lua
lsp.WorkspaceFolder[]
```


The array of the removed workspace folders


---

# lsp.WorkspaceFoldersInitializeParams

## workspaceFolders


```lua
(lsp.WorkspaceFolder[]|nil)?
```


The workspace folders configured in the client when the server starts.

This property is only available if the client supports workspace folders.
It can be `null` if the client supports workspace folders but none are
configured.



---

# lsp.WorkspaceFoldersServerCapabilities

## changeNotifications


```lua
(boolean|string)?
```


Whether the server wants to receive workspace folder
change notifications.

If a string is provided the string is treated as an ID
under which the notification is registered on the client
side. The ID can be used to unregister for these events
using the `client/unregisterCapability` request.

## supported


```lua
boolean?
```


The server has support for workspace folders


---

# lsp.WorkspaceFullDocumentDiagnosticReport

A full document diagnostic report for a workspace diagnostic result.


## items


```lua
lsp.Diagnostic[]
```


The actual items.

## kind


```lua
"full"
```


A full document diagnostic report.

## resultId


```lua
string?
```


An optional result id. If provided it will
be sent on the next diagnostic request for the
same document.

## uri


```lua
string
```


The URI for which diagnostic information is reported.

## version


```lua
integer|nil
```


The version number for which the diagnostics are reported.
If the document is not marked as open `null` can be provided.


---

# lsp.WorkspaceOptions

Defines workspace specific capabilities of the server.


## fileOperations


```lua
(lsp.FileOperationOptions)?
```


The server is interested in notifications/requests for operations on files.


## textDocumentContent


```lua
(lsp.TextDocumentContentOptions|lsp.TextDocumentContentRegistrationOptions)?
```


The server supports the `workspace/textDocumentContent` request.


## workspaceFolders


```lua
(lsp.WorkspaceFoldersServerCapabilities)?
```


The server supports workspace folder.



---

# lsp.WorkspaceSymbol

A special workspace symbol that supports locations without a range.

See also SymbolInformation.


## containerName


```lua
string?
```


The name of the symbol containing this symbol. This information is for
user interface purposes (e.g. to render a qualifier in the user interface
if necessary). It can't be used to re-infer a hierarchy for the document
symbols.

## data


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


A data entry field that is preserved on a workspace symbol between a
workspace symbol request and a workspace symbol resolve request.

## kind


```lua
1|10|11|12|13...(+21)
```


The kind of this symbol.

## location


```lua
lsp.Location|lsp.LocationUriOnly
```


The location of the symbol. Whether a server is allowed to
return a location without a range depends on the client
capability `workspace.symbol.resolveSupport`.

See SymbolInformation#location for more details.

## name


```lua
string
```


The name of this symbol.

## tags


```lua
1[]?
```


Tags for this symbol.



---

# lsp.WorkspaceSymbolClientCapabilities

Client capabilities for a {@link WorkspaceSymbolRequest}.

## dynamicRegistration


```lua
boolean?
```


Symbol request supports dynamic registration.

## resolveSupport


```lua
(lsp.ClientSymbolResolveOptions)?
```


The client support partial workspace symbols. The client will send the
request `workspaceSymbol/resolve` to the server to resolve additional
properties.


## symbolKind


```lua
(lsp.ClientSymbolKindOptions)?
```


Specific capabilities for the `SymbolKind` in the `workspace/symbol` request.

## tagSupport


```lua
(lsp.ClientSymbolTagOptions)?
```


The client supports tags on `SymbolInformation`.
Clients supporting tags have to handle unknown tags gracefully.



---

# lsp.WorkspaceSymbolOptions

Server capabilities for a {@link WorkspaceSymbolRequest}.

## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for a workspace symbol.


## workDoneProgress


```lua
boolean?
```



---

# lsp.WorkspaceSymbolParams

The parameters of a {@link WorkspaceSymbolRequest}.

## partialResultToken


```lua
(string|integer)?
```


An optional token that a server can use to report partial results (e.g. streaming) to
the client.

## query


```lua
string
```


A query string to filter symbols by. Clients may send an empty
string here to request all symbols.

The `query`-parameter should be interpreted in a *relaxed way* as editors
will apply their own highlighting and scoring on the results. A good rule
of thumb is to match case-insensitive and to simply check that the
characters of *query* appear in their order in a candidate symbol.
Servers shouldn't use prefix, substring, or similar strict matching.

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp.WorkspaceSymbolRegistrationOptions

Registration options for a {@link WorkspaceSymbolRequest}.

## resolveProvider


```lua
boolean?
```


The server provides support to resolve additional
information for a workspace symbol.


## workDoneProgress


```lua
boolean?
```



---

# lsp.WorkspaceUnchangedDocumentDiagnosticReport

An unchanged document diagnostic report for a workspace diagnostic result.


## kind


```lua
"unchanged"
```


A document diagnostic report indicating
no changes to the last result. A server can
only return `unchanged` if result ids are
provided.

## resultId


```lua
string
```


A result id which will be sent on the next
diagnostic request for the same document.

## uri


```lua
string
```


The URI for which diagnostic information is reported.

## version


```lua
integer|nil
```


The version number for which the diagnostics are reported.
If the document is not marked as open `null` can be provided.


---

# lsp._InitializeParams

The initialize parameters

## capabilities


```lua
lsp.ClientCapabilities
```


The capabilities provided by the client (editor or tool)

## clientInfo


```lua
(lsp.ClientInfo)?
```


Information about the client


## initializationOptions


```lua
(boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```


User provided initialization options.

## locale


```lua
string?
```


The locale the client is currently showing the user interface
in. This must not necessarily be the locale of the operating
system.

Uses IETF language tags as the value's syntax
(See https://en.wikipedia.org/wiki/IETF_language_tag)


## processId


```lua
integer|nil
```


The process Id of the parent process that started
the server.

Is `null` if the process has not been started by another process.
If the parent process is not alive then the server should exit.

## rootPath


```lua
(string|nil)?
```


The rootPath of the workspace. Is null
if no folder is open.


## rootUri


```lua
string|nil
```


The rootUri of the workspace. Is null if no
folder is open. If both `rootPath` and `rootUri` are set
`rootUri` wins.


## trace


```lua
("messages"|"off"|"verbose")?
```


The initial trace setting. If omitted trace is disabled ('off').

## workDoneToken


```lua
(string|integer)?
```


An optional token that a server can use to report work done progress.


---

# lsp._anonym1.range


---

# lsp._anonym2.range


---

# lsp.diagnostic.bufstate

## enabled


```lua
boolean
```

Whether inlay hints are enabled for this buffer


---

# lsp.null


---

# nvim.TutorExtmarks


---

# nvim.TutorMetadata

## expect


```lua
table<string, string|-1>
```


---

# uinteger


---

# vim


```lua
table
```


---

# vim.Diagnostic

 [diagnostic-structure](file:///usr/local/share/nvim/runtime/lua/vim)

 Diagnostics use the same indexing as the rest of the Nvim API (i.e. 0-based
 rows and columns). |api-indexing|

## _tags


```lua
{ deprecated: boolean, unnecessary: boolean }?
```


## bufnr


```lua
integer?
```


 Buffer number

## code


```lua
(string|integer)?
```


 The diagnostic code

## col


```lua
integer
```


 The starting column of the diagnostic (0-indexed)

## end_col


```lua
integer?
```


 The final column of the diagnostic (0-indexed)

## end_lnum


```lua
integer?
```


 The final line of the diagnostic (0-indexed)

## lnum


```lua
integer
```


 The starting line of the diagnostic (0-indexed)

## message


```lua
string
```


 The diagnostic text

## namespace


```lua
integer?
```


## severity


```lua
(vim.diagnostic.Severity)?
```


 The severity of the diagnostic |vim.diagnostic.severity|

## source


```lua
string?
```


 The source of the diagnostic

## user_data


```lua
any
```

arbitrary data plugins can add


---

# vim.F


```lua
table
```


---

# vim.NIL


---

# vim.NIL


```lua
unknown
```


---

# vim.Option

 @nodoc

## append


```lua
(method) vim.Option:append(value: string)
```

 Append a value to string-style options. See |:set+=|

 These are equivalent:

 ```lua
 vim.opt.formatoptions:append('j')
 vim.opt.formatoptions = vim.opt.formatoptions + 'j'
 ```

@*param* `value` — Value to append

 luacheck: no unused

## get


```lua
(method) vim.Option:get()
  -> value: boolean|string|integer|nil
```

 Returns a Lua-representation of the option. Boolean, number and string
 values will be returned in exactly the same fashion.

 For values that are comma-separated lists, an array will be returned with
 the values as entries in the array:

 ```lua
 vim.cmd [[set wildignore=*.pyc,*.o]]

 vim.print(vim.opt.wildignore:get())
 -- { "*.pyc", "*.o", }

 for _, ignore_pattern in ipairs(vim.opt.wildignore:get()) do
     print("Will ignore:", ignore_pattern)
 end
 -- Will ignore: *.pyc
 -- Will ignore: *.o
 ```

 For values that are comma-separated maps, a table will be returned with
 the names as keys and the values as entries:

 ```lua
 vim.cmd [[set listchars=space:_,tab:>~]]

 vim.print(vim.opt.listchars:get())
 --  { space = "_", tab = ">~", }

 for char, representation in pairs(vim.opt.listchars:get()) do
     print(char, "=>", representation)
 end
 ```

 For values that are lists of flags, a set will be returned with the flags
 as keys and `true` as entries.

 ```lua
 vim.cmd [[set formatoptions=njtcroql]]

 vim.print(vim.opt.formatoptions:get())
 -- { n = true, j = true, c = true, ... }

 local format_opts = vim.opt.formatoptions:get()
 if format_opts.j then
     print("J is enabled!")
 end
 ```

@*return* `value` — of option

## prepend


```lua
(method) vim.Option:prepend(value: string)
```

 Prepend a value to string-style options. See |:set^=|

 These are equivalent:

 ```lua
 vim.opt.wildignore:prepend('*.o')
 vim.opt.wildignore = vim.opt.wildignore ^ '*.o'
 ```

@*param* `value` — Value to prepend

 luacheck: no unused

## remove


```lua
(method) vim.Option:remove(value: string)
```

 Remove a value from string-style options. See |:set-=|

 These are equivalent:

 ```lua
 vim.opt.wildignore:remove('*.pyc')
 vim.opt.wildignore = vim.opt.wildignore - '*.pyc'
 ```

@*param* `value` — Value to remove

 luacheck: no unused


---

# vim.Ringbuf

## _idx_read


```lua
integer
```

## _idx_write


```lua
integer
```

## _items


```lua
table[]
```

## _size


```lua
integer
```

## clear


```lua
function vim.Ringbuf.clear(self: vim.Ringbuf)
```

 Clear all items

## peek


```lua
function vim.Ringbuf.peek(self: vim.Ringbuf)
  -> <T>?
```

 Returns the first unread item without removing it

## pop


```lua
function vim.Ringbuf.pop(self: vim.Ringbuf)
  -> <T>?
```

 Removes and returns the first unread item

## push


```lua
function vim.Ringbuf.push(self: vim.Ringbuf, item: <T>)
```

 Adds an item, overriding the oldest item if the buffer is full.


---

# vim.SystemCompleted

## code


```lua
integer
```

## signal


```lua
integer
```

## stderr


```lua
string?
```

## stdout


```lua
string?
```


---

# vim.SystemObj

## _state


```lua
vim.SystemState
```

## _timeout


```lua
(method) vim.SystemObj:_timeout(signal?: vim.SystemSig)
```

## cmd


```lua
string[]
```

## is_closing


```lua
fun(self: vim.SystemObj):boolean
```

## kill


```lua
fun(self: vim.SystemObj, signal: string|integer)
```

## pid


```lua
integer
```

## wait


```lua
fun(self: vim.SystemObj, timeout?: integer):vim.SystemCompleted
```

## write


```lua
fun(self: vim.SystemObj, data?: string|string[])
```


---

# vim.SystemOpts

## clear_env


```lua
boolean?
```

## cwd


```lua
string?
```

## detach


```lua
boolean?
```

## env


```lua
table<string, string|number>?
```

## stderr


```lua
(fun(err?: string, data?: string)|false)?
```

## stdin


```lua
(string|string[]|true)?
```

## stdout


```lua
(fun(err?: string, data?: string)|false)?
```

## text


```lua
boolean?
```

## timeout


```lua
integer?
```

Timeout in ms


---

# vim.SystemSig


---

# vim.SystemState

## cmd


```lua
string[]
```

## done


```lua
(boolean|'timeout')?
```

## handle


```lua
(uv.uv_process_t)?
```

## pid


```lua
integer?
```

## result


```lua
(vim.SystemCompleted)?
```

## stderr


```lua
(uv.uv_stream_t)?
```

## stderr_data


```lua
string[]?
```

## stdin


```lua
(uv.uv_stream_t)?
```

## stdout


```lua
(uv.uv_stream_t)?
```

## stdout_data


```lua
string[]?
```

## timeout


```lua
integer?
```

## timer


```lua
(uv.uv_timer_t)?
```


---

# vim.Version

## [1]


```lua
number
```

## [2]


```lua
number
```

## [3]


```lua
number
```

## __eq


```lua
(method) vim.Version:__eq(other: vim.Version)
  -> boolean
```

## __index


```lua
vim.Version
```

## __le


```lua
(method) vim.Version:__le(other: vim.Version)
  -> boolean
```

## __lt


```lua
(method) vim.Version:__lt(other: vim.Version)
  -> boolean
```

## __newindex


```lua
(method) vim.Version:__newindex(key: any, value: any)
```

## __tostring


```lua
(method) vim.Version:__tostring()
  -> string
```

## build


```lua
string?
```

## major


```lua
number
```

## minor


```lua
number
```

## patch


```lua
number
```

## prerelease


```lua
string?
```


---

# vim.VersionRange

## from


```lua
vim.Version
```

## has


```lua
(method) vim.VersionRange:has(version: string|vim.Version)
  -> boolean
```

## to


```lua
(vim.Version)?
```


---

# vim._comment.Parts

## left


```lua
string
```

Left part of comment

## right


```lua
string
```

Right part of comment


---

# vim._create_ts_parser


```lua
function vim._create_ts_parser(lang: string)
  -> TSParser
```


---

# vim._create_ts_querycursor


```lua
function vim._create_ts_querycursor(node: TSNode, query: TSQuery, start?: integer, stop?: integer, opts?: { max_start_depth: integer, match_limit: integer })
  -> TSQueryCursor
```


---

# vim._cs_remote


```lua
function vim._cs_remote(rcid: any, server_addr: string, connect_error: string, args: any)
  -> table
```


---

# vim._defer_deprecated_module


```lua
function vim._defer_deprecated_module(old_name: string, new_name: string)
  -> table
```


---

# vim._defer_require


```lua
function vim._defer_require(root: string, mod: <T>)
  -> <T>
```


---

# vim._ensure_list


```lua
function vim._ensure_list(x?: <T>|elem_or_list<T>)
  -> <T>[]
```


---

# vim._expand_pat


```lua
function vim._expand_pat(pat: string, env: any)
  -> any[]
  2. integer
```


---

# vim._expand_pat_get_parts


```lua
function vim._expand_pat_get_parts(lua_string: string)
  -> (string|string[])[]
  2. integer
```


---

# vim._extra


```lua
table
```


---

# vim._inspector.Filter

## extmarks


```lua
boolean|"all"
```


 Include extmarks. When `all`, then extmarks without a `hl_group` will also be included.
 (default: true)

## semantic_tokens


```lua
boolean
```


 Include semantic token highlights.
 (default: true)

## syntax


```lua
boolean
```

 @inlinedoc

 Include syntax based highlight groups.
 (default: `true`)

## treesitter


```lua
boolean
```


 Include treesitter based highlight groups.
 (default: `true`)


---

# vim._list_insert


```lua
function vim._list_insert(t: any[], first: integer, last: integer, v: any)
```


---

# vim._list_remove


```lua
function vim._list_remove(t: any[], first: integer, last: integer)
```


---

# vim._load_package


```lua
function vim._load_package(name: string)
  -> function|nil
```


---

# vim._on_key


```lua
function vim._on_key(buf: any, typed_buf: any)
  -> boolean
```


---

# vim._option.Info

 @nodoc

## allows_duplicates


```lua
boolean
```

## commalist


```lua
boolean
```

## default


```lua
boolean|string|integer
```

## flaglist


```lua
boolean
```

## global_local


```lua
boolean
```

## last_set_chan


```lua
integer
```

## last_set_linenr


```lua
integer
```

## last_set_sid


```lua
integer
```

## metatype


```lua
'array'|'boolean'|'map'|'number'|'set'...(+1)
```

## name


```lua
string
```

## scope


```lua
'buf'|'global'|'win'
```

## shortname


```lua
string
```

## type


```lua
'boolean'|'number'|'string'
```

## was_set


```lua
boolean
```


---

# vim._os_proc_children


```lua
function vim._os_proc_children(ppid: any)
  -> table
```


---

# vim._os_proc_info


```lua
function vim._os_proc_info(pid: any)
  -> table
```


---

# vim._print


```lua
function vim._print(inspect_strings: boolean, ...any)
  -> unknown
```


---

# vim._resolve_bufnr


```lua
function vim._resolve_bufnr(bufnr?: integer)
  -> integer
```


---

# vim._so_trails


```lua
string[]
```


---

# vim._submodules


```lua
table
```


---

# vim._truncated_echo_once


```lua
function vim._truncated_echo_once(msg: any)
  -> boolean
```


---

# vim._ts_add_language_from_object


```lua
function vim._ts_add_language_from_object(path: string, lang: string, symbol_name?: string)
```


---

# vim._ts_add_language_from_wasm


```lua
function vim._ts_add_language_from_wasm(path: string, lang: string)
```


---

# vim._ts_get_language_version


```lua
function vim._ts_get_language_version()
  -> integer
```


---

# vim._ts_get_minimum_language_version


```lua
function vim._ts_get_minimum_language_version()
  -> integer
```


---

# vim._ts_inspect_language


```lua
function vim._ts_inspect_language(lang: string)
  -> TSLangInfo
```


---

# vim._ts_parse_query


```lua
function vim._ts_parse_query(lang: string, query: string)
  -> TSQuery
```


---

# vim._watch


```lua
table
```


---

# vim._watch.Callback


---

# vim._watch.FileChangeType


---

# vim._watch.Opts

## debounce


```lua
integer?
```

ms

## exclude_pattern


```lua
(vim.lpeg.Pattern)?
```


 An |lpeg| pattern. Only changes to files and directories whose full path does
 not match the pattern will be reported. Matches against both files and
 directories. When nil, matches nothing.

## include_pattern


```lua
(vim.lpeg.Pattern)?
```


 An |lpeg| pattern. Only changes to files whose full paths match the pattern
 will be reported. Only matches against non-directoriess, all directories will
 be watched for new potentially-matching files. exclude_pattern can be used to
 filter out directories. When nil, matches any file name.


---

# vim._watch.watch.Opts

## debounce


```lua
integer?
```

ms

## exclude_pattern


```lua
(vim.lpeg.Pattern)?
```


 An |lpeg| pattern. Only changes to files and directories whose full path does
 not match the pattern will be reported. Matches against both files and
 directories. When nil, matches nothing.

## include_pattern


```lua
(vim.lpeg.Pattern)?
```


 An |lpeg| pattern. Only changes to files whose full paths match the pattern
 will be reported. Only matches against non-directoriess, all directories will
 be watched for new potentially-matching files. exclude_pattern can be used to
 filter out directories. When nil, matches any file name.

## uvflags


```lua
(uv.fs_event_start.flags)?
```


---

# vim._with


```lua
function vim._with(context: vim.context.mods, f: function)
  -> any
```


---

# vim.api


```lua
table
```


---

# vim.api.keyset.buf_attach

## on_bytes


```lua
fun(_: "bytes", bufnr: integer, changedtick: integer, start_row: integer, start_col: integer, start_byte: integer, old_end_row: integer, old_end_col: integer, old_end_byte: integer, new_end_row: integer, new_end_col: integer, new_end_byte: integer):boolean??
```

## on_changedtick


```lua
fun(_: "changedtick", bufnr: integer, changedtick: integer)?
```

## on_detach


```lua
fun(_: "detach", bufnr: integer)?
```

## on_lines


```lua
fun(_: "lines", bufnr: integer, changedtick: integer, first: integer, last_old: integer, last_new: integer, byte_count: integer, deleted_codepoints?: integer, deleted_codeunits?: integer):boolean??
```

## on_reload


```lua
fun(_: "reload", bufnr: integer)?
```

## preview


```lua
boolean?
```

## utf_sizes


```lua
boolean?
```


---

# vim.api.keyset.buf_delete

## force


```lua
boolean?
```

## unload


```lua
boolean?
```


---

# vim.api.keyset.clear_autocmds

## buffer


```lua
integer?
```

## event


```lua
(string|string[])?
```

## group


```lua
(string|integer)?
```

## pattern


```lua
(string|string[])?
```


---

# vim.api.keyset.cmd

## addr


```lua
string?
```

## args


```lua
string[]?
```

## bang


```lua
boolean?
```

## cmd


```lua
string?
```

## count


```lua
integer?
```

## magic


```lua
table<string, any>?
```

## mods


```lua
table<string, any>?
```

## nargs


```lua
(string|integer)?
```

## nextcmd


```lua
string?
```

## range


```lua
any[]?
```

## reg


```lua
string?
```


---

# vim.api.keyset.cmd_magic

## bar


```lua
boolean?
```

## file


```lua
boolean?
```


---

# vim.api.keyset.cmd_mods

## browse


```lua
boolean?
```

## confirm


```lua
boolean?
```

## emsg_silent


```lua
boolean?
```

## filter


```lua
table<string, any>?
```

## hide


```lua
boolean?
```

## horizontal


```lua
boolean?
```

## keepalt


```lua
boolean?
```

## keepjumps


```lua
boolean?
```

## keepmarks


```lua
boolean?
```

## keeppatterns


```lua
boolean?
```

## lockmarks


```lua
boolean?
```

## noautocmd


```lua
boolean?
```

## noswapfile


```lua
boolean?
```

## sandbox


```lua
boolean?
```

## silent


```lua
boolean?
```

## split


```lua
string?
```

## tab


```lua
integer?
```

## unsilent


```lua
boolean?
```

## verbose


```lua
integer?
```

## vertical


```lua
boolean?
```


---

# vim.api.keyset.cmd_mods_filter

## force


```lua
boolean?
```

## pattern


```lua
string?
```


---

# vim.api.keyset.cmd_opts

## output


```lua
boolean?
```


---

# vim.api.keyset.command_info

## addr


```lua
string?
```

## bang


```lua
boolean
```

## bar


```lua
boolean
```

## complete


```lua
string?
```

## complete_arg


```lua
string?
```

## count


```lua
string?
```

## definition


```lua
string
```

## keepscript


```lua
boolean
```

## name


```lua
string
```

## nargs


```lua
string
```

## preview


```lua
boolean
```

## range


```lua
string?
```

## register


```lua
boolean
```

## script_id


```lua
integer
```


---

# vim.api.keyset.complete_set

## info


```lua
string?
```


---

# vim.api.keyset.context

## types


```lua
string[]?
```


---

# vim.api.keyset.create_augroup

## clear


```lua
boolean?
```


---

# vim.api.keyset.create_autocmd

## buffer


```lua
integer?
```

## callback


```lua
(string|fun(args: vim.api.keyset.create_autocmd.callback_args):boolean?)?
```

## command


```lua
string?
```

## desc


```lua
string?
```

## group


```lua
(string|integer)?
```

## nested


```lua
boolean?
```

## once


```lua
boolean?
```

## pattern


```lua
(string|string[])?
```


---

# vim.api.keyset.create_autocmd.callback_args

## buf


```lua
integer
```

expanded value of <abuf>

## data


```lua
any
```

arbitrary data passed from |nvim_exec_autocmds()|                       *event-data*

## event


```lua
string
```

name of the triggered event |autocmd-events|

## file


```lua
string
```

expanded value of <afile>

## group


```lua
integer?
```

autocommand group id, if any

## id


```lua
integer
```

autocommand id

## match


```lua
string
```

expanded value of <amatch>


---

# vim.api.keyset.create_user_command.command_args

## args


```lua
string
```


 The args passed to the command, if any <args>

## bang


```lua
boolean
```


 "true" if the command was executed with a ! modifier <bang>

## count


```lua
integer
```


 Any count supplied <count>

## fargs


```lua
string[]
```


 The args split by unescaped whitespace
 (when more than one argument is allowed), if any <f-args>

## line1


```lua
integer
```


 The starting line of the command range <line1>

## line2


```lua
integer
```


 The final line of the command range <line2>

## mods


```lua
string
```

 Command modifiers, if any <mods>

## name


```lua
string
```

Command name

## nargs


```lua
string
```


 Number of arguments |:command-nargs|

## range


```lua
integer
```


 The number of items in the command range: 0, 1, or 2 <range>

## reg


```lua
string
```

 The optional register, if specified <reg>

## smods


```lua
table
```


 Command modifiers in a structured format. Has the same structure as the
 "mods" key of |nvim_parse_cmd()|.


---

# vim.api.keyset.echo_opts

## err


```lua
boolean?
```

## verbose


```lua
boolean?
```


---

# vim.api.keyset.empty


---

# vim.api.keyset.eval_statusline

## fillchar


```lua
string?
```

## highlights


```lua
boolean?
```

## maxwidth


```lua
integer?
```

## use_statuscol_lnum


```lua
integer?
```

## use_tabline


```lua
boolean?
```

## use_winbar


```lua
boolean?
```

## winid


```lua
integer?
```


---

# vim.api.keyset.exec_autocmds

## buffer


```lua
integer?
```

## data


```lua
any
```

## group


```lua
(string|integer)?
```

## modeline


```lua
boolean?
```

## pattern


```lua
(string|string[])?
```


---

# vim.api.keyset.exec_opts

## output


```lua
boolean?
```


---

# vim.api.keyset.extmark_details

## conceal


```lua
boolean?
```


## cursorline_hl_group


```lua
string?
```

## end_col


```lua
integer?
```

## end_right_gravity


```lua
integer?
```

## end_row


```lua
integer?
```


## hl_eol


```lua
boolean?
```

## hl_group


```lua
string?
```


## hl_mode


```lua
string?
```

## invalid


```lua
true?
```

## invalidate


```lua
true?
```

## line_hl_group


```lua
string?
```

## ns_id


```lua
integer
```

## number_hl_group


```lua
string?
```

## priority


```lua
integer?
```


## right_gravity


```lua
boolean
```

## sign_hl_group


```lua
string?
```

## sign_name


```lua
string?
```

## sign_text


```lua
string?
```


## spell


```lua
boolean?
```

## ui_watched


```lua
boolean?
```

## undo_restore


```lua
false?
```


## url


```lua
string?
```

## virt_lines


```lua
[string, string][][]?
```


## virt_lines_above


```lua
boolean?
```

## virt_lines_leftcol


```lua
boolean?
```

## virt_text


```lua
[string, string][]?
```


## virt_text_hide


```lua
boolean?
```

## virt_text_pos


```lua
string?
```

## virt_text_repeat_linebreak


```lua
boolean?
```

## virt_text_win_col


```lua
integer?
```


---

# vim.api.keyset.get_autocmds

## buffer


```lua
(integer|integer[])?
```

## event


```lua
(string|string[])?
```

## group


```lua
(string|integer)?
```

## id


```lua
integer?
```

## pattern


```lua
(string|string[])?
```


---

# vim.api.keyset.get_autocmds.ret

## buffer


```lua
integer?
```

## buflocal


```lua
boolean?
```

## callback


```lua
function?
```

## command


```lua
string?
```

## desc


```lua
string?
```

## event


```lua
string?
```

## group


```lua
integer?
```

## group_name


```lua
integer?
```

## id


```lua
integer?
```

## once


```lua
boolean?
```

## pattern


```lua
string?
```


---

# vim.api.keyset.get_commands

## builtin


```lua
boolean?
```


---

# vim.api.keyset.get_extmark

## details


```lua
boolean?
```

## hl_name


```lua
boolean?
```


---

# vim.api.keyset.get_extmark_item

## [1]


```lua
integer
```

extmark_id

## [2]


```lua
integer
```

row

## [3]


```lua
integer
```

col

## [4]


```lua
(vim.api.keyset.extmark_details)?
```


---

# vim.api.keyset.get_extmark_item_by_id

## [1]


```lua
integer
```

row

## [2]


```lua
integer
```

col

## [3]


```lua
(vim.api.keyset.extmark_details)?
```


---

# vim.api.keyset.get_extmarks

## details


```lua
boolean?
```

## hl_name


```lua
boolean?
```

## limit


```lua
integer?
```

## overlap


```lua
boolean?
```

## type


```lua
string?
```


---

# vim.api.keyset.get_highlight

## create


```lua
boolean?
```

## id


```lua
integer?
```

## link


```lua
boolean?
```

## name


```lua
string?
```


---

# vim.api.keyset.get_hl_info

## altfont


```lua
true?
```

## bg


```lua
integer?
```

## blend


```lua
integer?
```

## bold


```lua
true?
```

## cterm


```lua
(vim.api.keyset.hl_info.cterm)?
```

## ctermbg


```lua
integer?
```

## ctermfg


```lua
integer?
```

## default


```lua
true?
```

## fg


```lua
integer?
```

## italic


```lua
true?
```

## link


```lua
string?
```

## nocombine


```lua
true?
```

## reverse


```lua
true?
```

## sp


```lua
integer?
```

## standout


```lua
true?
```

## strikethrough


```lua
true?
```

## undercurl


```lua
true?
```

## underdashed


```lua
true?
```

## underdotted


```lua
true?
```

## underdouble


```lua
true?
```

## underline


```lua
true?
```


---

# vim.api.keyset.get_keymap

## abbr


```lua
(0|1)?
```

## buffer


```lua
(0|1)?
```

## callback


```lua
function?
```

## desc


```lua
string?
```

## expr


```lua
(0|1)?
```

## lhs


```lua
string?
```

## lhsraw


```lua
string?
```

## lhsrawalt


```lua
string?
```

## lnum


```lua
integer?
```

## mode


```lua
string?
```

## mode_bits


```lua
integer?
```

## noremap


```lua
(0|1)?
```

## nowait


```lua
(0|1)?
```

## rhs


```lua
string?
```

## script


```lua
(0|1)?
```

## scriptversion


```lua
integer?
```

## sid


```lua
integer?
```

## silent


```lua
(0|1)?
```


---

# vim.api.keyset.get_mark

## [1]


```lua
integer
```

row

## [2]


```lua
integer
```

col

## [3]


```lua
integer
```

buffer

## [4]


```lua
string
```

buffername


---

# vim.api.keyset.get_mode

## blocking


```lua
boolean
```

## mode


```lua
string
```


---

# vim.api.keyset.get_ns

## winid


```lua
integer?
```


---

# vim.api.keyset.get_option_info

## allows_duplicates


```lua
boolean
```

## commalist


```lua
boolean
```

## default


```lua
boolean|string|integer
```

## flaglist


```lua
boolean
```

## global_local


```lua
boolean
```

## last_set_chan


```lua
integer
```

## last_set_linenr


```lua
integer
```

## last_set_sid


```lua
integer
```

## name


```lua
string
```

## scope


```lua
'buf'|'global'|'win'
```

## shortname


```lua
string
```

## type


```lua
'boolean'|'number'|'string'
```

## was_set


```lua
boolean
```


---

# vim.api.keyset.highlight

## altfont


```lua
boolean?
```

## background


```lua
(string|integer)?
```

## bg


```lua
(string|integer)?
```

## bg_indexed


```lua
boolean?
```

## blend


```lua
integer?
```

## bold


```lua
boolean?
```

## cterm


```lua
(string|integer)?
```

## ctermbg


```lua
(string|integer)?
```

## ctermfg


```lua
(string|integer)?
```

## default


```lua
boolean?
```

## fallback


```lua
boolean?
```

## fg


```lua
(string|integer)?
```

## fg_indexed


```lua
boolean?
```

## force


```lua
boolean?
```

## foreground


```lua
(string|integer)?
```

## global_link


```lua
(string|integer)?
```

## italic


```lua
boolean?
```

## link


```lua
(string|integer)?
```

## nocombine


```lua
boolean?
```

## reverse


```lua
boolean?
```

## sp


```lua
(string|integer)?
```

## special


```lua
(string|integer)?
```

## standout


```lua
boolean?
```

## strikethrough


```lua
boolean?
```

## undercurl


```lua
boolean?
```

## underdashed


```lua
boolean?
```

## underdotted


```lua
boolean?
```

## underdouble


```lua
boolean?
```

## underline


```lua
boolean?
```

## url


```lua
string?
```


---

# vim.api.keyset.highlight_cterm

## altfont


```lua
boolean?
```

## bold


```lua
boolean?
```

## italic


```lua
boolean?
```

## nocombine


```lua
boolean?
```

## reverse


```lua
boolean?
```

## standout


```lua
boolean?
```

## strikethrough


```lua
boolean?
```

## undercurl


```lua
boolean?
```

## underdashed


```lua
boolean?
```

## underdotted


```lua
boolean?
```

## underdouble


```lua
boolean?
```

## underline


```lua
boolean?
```


---

# vim.api.keyset.hl_info.base

## altfont


```lua
true?
```

## bold


```lua
true?
```

## ctermbg


```lua
integer?
```

## ctermfg


```lua
integer?
```

## italic


```lua
true?
```

## nocombine


```lua
true?
```

## reverse


```lua
true?
```

## standout


```lua
true?
```

## strikethrough


```lua
true?
```

## undercurl


```lua
true?
```

## underdashed


```lua
true?
```

## underdotted


```lua
true?
```

## underdouble


```lua
true?
```

## underline


```lua
true?
```


---

# vim.api.keyset.hl_info.cterm

## altfont


```lua
true?
```

## background


```lua
integer?
```

## bold


```lua
true?
```

## ctermbg


```lua
integer?
```

## ctermfg


```lua
integer?
```

## foreground


```lua
integer?
```

## italic


```lua
true?
```

## nocombine


```lua
true?
```

## reverse


```lua
true?
```

## standout


```lua
true?
```

## strikethrough


```lua
true?
```

## undercurl


```lua
true?
```

## underdashed


```lua
true?
```

## underdotted


```lua
true?
```

## underdouble


```lua
true?
```

## underline


```lua
true?
```


---

# vim.api.keyset.keymap

## callback


```lua
function?
```

## desc


```lua
string?
```

## expr


```lua
boolean?
```

## noremap


```lua
boolean?
```

## nowait


```lua
boolean?
```

## replace_keycodes


```lua
boolean?
```

## script


```lua
boolean?
```

## silent


```lua
boolean?
```

## unique


```lua
boolean?
```


---

# vim.api.keyset.ns_opts

## wins


```lua
any[]?
```


---

# vim.api.keyset.open_term

## force_crlf


```lua
boolean?
```

## on_input


```lua
fun(_: "input", term: integer, bufnr: integer, data: any)?
```


---

# vim.api.keyset.option

## buf


```lua
integer?
```

## filetype


```lua
string?
```

## scope


```lua
string?
```

## win


```lua
integer?
```


---

# vim.api.keyset.parse_cmd

## addr


```lua
'?'|'arg'|'buf'|'line'|'load'...(+4)
```

## args


```lua
string[]
```

## bang


```lua
boolean
```

## cmd


```lua
string
```

## count


```lua
integer?
```

## magic


```lua
{ bar: boolean, file: boolean }
```

## mods


```lua
vim.api.keyset.parse_cmd.mods
```

## nargs


```lua
'*'|'+'|'0'|'1'|'?'
```

## nextcmd


```lua
string
```

## range


```lua
integer[]?
```

## reg


```lua
string?
```


---

# vim.api.keyset.parse_cmd.mods

## browse


```lua
boolean
```

## confirm


```lua
boolean
```

## emsg_silent


```lua
boolean
```

## filter


```lua
{ force: boolean, pattern: string }
```

## hide


```lua
boolean
```

## horizontal


```lua
boolean
```

## keepalt


```lua
boolean
```

## keepjumps


```lua
boolean
```

## keepmarks


```lua
boolean
```

## keeppatterns


```lua
boolean
```

## lockmarks


```lua
boolean
```

## noautocmd


```lua
boolean
```

## noswapfile


```lua
boolean
```

## sandbox


```lua
boolean
```

## silent


```lua
boolean
```

## split


```lua
''|'aboveleft'|'belowright'|'botright'|'topleft'
```

## tab


```lua
integer
```

## unsilent


```lua
boolean
```

## verbose


```lua
integer
```

## vertical


```lua
boolean
```


---

# vim.api.keyset.redraw

## buf


```lua
integer?
```

## cursor


```lua
boolean?
```

## flush


```lua
boolean?
```

## range


```lua
any[]?
```

## statuscolumn


```lua
boolean?
```

## statusline


```lua
boolean?
```

## tabline


```lua
boolean?
```

## valid


```lua
boolean?
```

## win


```lua
integer?
```

## winbar


```lua
boolean?
```


---

# vim.api.keyset.runtime

## do_source


```lua
boolean?
```

## is_lua


```lua
boolean?
```


---

# vim.api.keyset.set_decoration_provider

## _on_conceal_line


```lua
fun(_: "conceal_line")?
```

## _on_hl_def


```lua
fun(_: "hl_def")?
```

## _on_spell_nav


```lua
fun(_: "spell_nav")?
```

## on_buf


```lua
fun(_: "buf", bufnr: integer, tick: integer)?
```

## on_end


```lua
fun(_: "end", tick: integer)?
```

## on_line


```lua
fun(_: "line", winid: integer, bufnr: integer, row: integer):boolean??
```

## on_start


```lua
fun(_: "start", tick: integer):boolean??
```

## on_win


```lua
fun(_: "win", winid: integer, bufnr: integer, toprow: integer, botrow: integer):boolean??
```


---

# vim.api.keyset.set_extmark

## conceal


```lua
string?
```

## conceal_lines


```lua
string?
```

## cursorline_hl_group


```lua
(string|integer)?
```

## end_col


```lua
integer?
```

## end_line


```lua
integer?
```

## end_right_gravity


```lua
boolean?
```

## end_row


```lua
integer?
```

## ephemeral


```lua
boolean?
```

## hl_eol


```lua
boolean?
```

## hl_group


```lua
any
```

## hl_mode


```lua
string?
```

## id


```lua
integer?
```

## invalidate


```lua
boolean?
```

## line_hl_group


```lua
(string|integer)?
```

## number_hl_group


```lua
(string|integer)?
```

## priority


```lua
integer?
```

## right_gravity


```lua
boolean?
```

## scoped


```lua
boolean?
```

## sign_hl_group


```lua
(string|integer)?
```

## sign_text


```lua
string?
```

## spell


```lua
boolean?
```

## strict


```lua
boolean?
```

## ui_watched


```lua
boolean?
```

## undo_restore


```lua
boolean?
```

## url


```lua
string?
```

## virt_lines


```lua
any[]?
```

## virt_lines_above


```lua
boolean?
```

## virt_lines_leftcol


```lua
boolean?
```

## virt_lines_overflow


```lua
string?
```

## virt_text


```lua
any[]?
```

## virt_text_hide


```lua
boolean?
```

## virt_text_pos


```lua
string?
```

## virt_text_repeat_linebreak


```lua
boolean?
```

## virt_text_win_col


```lua
integer?
```


---

# vim.api.keyset.set_hl_info

## altfont


```lua
true?
```

## bg


```lua
(string|integer)?
```

## blend


```lua
integer?
```

## bold


```lua
true?
```

## cterm


```lua
(vim.api.keyset.hl_info.cterm)?
```

## ctermbg


```lua
integer?
```

## ctermfg


```lua
integer?
```

## default


```lua
true?
```

## fg


```lua
(string|integer)?
```

## force


```lua
true?
```

## italic


```lua
true?
```

## link


```lua
string?
```

## nocombine


```lua
true?
```

## reverse


```lua
true?
```

## sp


```lua
(string|integer)?
```

## standout


```lua
true?
```

## strikethrough


```lua
true?
```

## undercurl


```lua
true?
```

## underdashed


```lua
true?
```

## underdotted


```lua
true?
```

## underdouble


```lua
true?
```

## underline


```lua
true?
```


---

# vim.api.keyset.user_command

## addr


```lua
any
```

## bang


```lua
boolean?
```

## bar


```lua
boolean?
```

## complete


```lua
any
```

## count


```lua
any
```

## desc


```lua
any
```

## force


```lua
boolean?
```

## keepscript


```lua
boolean?
```

## nargs


```lua
any
```

## preview


```lua
any
```

## range


```lua
any
```

## register


```lua
boolean?
```


---

# vim.api.keyset.win_config

## anchor


```lua
('NE'|'NW'|'SE'|'SW')?
```

## border


```lua
('double'|'none'|'rounded'|'shadow'|'single'...(+2))?
```

## bufpos


```lua
integer[]?
```

## col


```lua
number?
```

## external


```lua
boolean?
```

## fixed


```lua
boolean?
```

## focusable


```lua
boolean?
```

## footer


```lua
any
```

## footer_pos


```lua
('center'|'left'|'right')?
```

## height


```lua
integer?
```

## hide


```lua
boolean?
```

## mouse


```lua
boolean?
```

## noautocmd


```lua
boolean?
```

## relative


```lua
('cursor'|'editor'|'laststatus'|'mouse'|'tabline'...(+1))?
```

## row


```lua
number?
```

## split


```lua
('above'|'below'|'left'|'right')?
```

## style


```lua
'minimal'?
```

## title


```lua
any
```

## title_pos


```lua
('center'|'left'|'right')?
```

## vertical


```lua
boolean?
```

## width


```lua
integer?
```

## win


```lua
integer?
```

## zindex


```lua
integer?
```


---

# vim.api.keyset.win_text_height

## end_row


```lua
integer?
```

## end_vcol


```lua
integer?
```

## start_row


```lua
integer?
```

## start_vcol


```lua
integer?
```


---

# vim.api.keyset.xdl_diff

## algorithm


```lua
string?
```

## ctxlen


```lua
integer?
```

## ignore_blank_lines


```lua
boolean?
```

## ignore_cr_at_eol


```lua
boolean?
```

## ignore_whitespace


```lua
boolean?
```

## ignore_whitespace_change


```lua
boolean?
```

## ignore_whitespace_change_at_eol


```lua
boolean?
```

## indent_heuristic


```lua
boolean?
```

## interhunkctxlen


```lua
integer?
```

## linematch


```lua
(boolean|integer)?
```

## on_hunk


```lua
fun(start_a: integer, count_a: integer, start_b: integer, count_b: integer):integer??
```

## result_type


```lua
string?
```


---

# vim.api.nvim__buf_debug_extmarks


```lua
function vim.api.nvim__buf_debug_extmarks(buffer: integer, keys: boolean, dot: boolean)
  -> string
```


---

# vim.api.nvim__buf_stats


```lua
function vim.api.nvim__buf_stats(buffer: integer)
  -> table<string, any>
```


---

# vim.api.nvim__complete_set


```lua
function vim.api.nvim__complete_set(index: integer, opts: vim.api.keyset.complete_set)
  -> table<string, any>
```


---

# vim.api.nvim__get_lib_dir


```lua
function vim.api.nvim__get_lib_dir()
  -> string
```


---

# vim.api.nvim__get_runtime


```lua
function vim.api.nvim__get_runtime(pat: any[], all: boolean, opts: vim.api.keyset.runtime)
  -> string[]
```


---

# vim.api.nvim__id


```lua
function vim.api.nvim__id(obj: any)
  -> any
```


---

# vim.api.nvim__id_array


```lua
function vim.api.nvim__id_array(arr: any[])
  -> any[]
```


---

# vim.api.nvim__id_dict


```lua
function vim.api.nvim__id_dict(dct: table<string, any>)
  -> table<string, any>
```


---

# vim.api.nvim__id_float


```lua
function vim.api.nvim__id_float(flt: number)
  -> number
```


---

# vim.api.nvim__inspect_cell


```lua
function vim.api.nvim__inspect_cell(grid: integer, row: integer, col: integer)
  -> any[]
```


---

# vim.api.nvim__invalidate_glyph_cache


```lua
function vim.api.nvim__invalidate_glyph_cache()
```


---

# vim.api.nvim__ns_get


```lua
function vim.api.nvim__ns_get(ns_id: integer)
  -> vim.api.keyset.ns_opts
```


---

# vim.api.nvim__ns_set


```lua
function vim.api.nvim__ns_set(ns_id: integer, opts: vim.api.keyset.ns_opts)
```


---

# vim.api.nvim__redraw


```lua
function vim.api.nvim__redraw(opts: vim.api.keyset.redraw)
```


---

# vim.api.nvim__runtime_inspect


```lua
function vim.api.nvim__runtime_inspect()
  -> any[]
```


---

# vim.api.nvim__screenshot


```lua
function vim.api.nvim__screenshot(path: string)
```


---

# vim.api.nvim__stats


```lua
function vim.api.nvim__stats()
  -> table<string, any>
```


---

# vim.api.nvim__unpack


```lua
function vim.api.nvim__unpack(str: string)
  -> any
```


---

# vim.api.nvim_buf_add_highlight


```lua
function vim.api.nvim_buf_add_highlight(buffer: integer, ns_id: integer, hl_group: string, line: integer, col_start: integer, col_end: integer)
  -> integer
```


---

# vim.api.nvim_buf_attach


```lua
function vim.api.nvim_buf_attach(buffer: integer, send_buffer: boolean, opts: vim.api.keyset.buf_attach)
  -> boolean
```


---

# vim.api.nvim_buf_call


```lua
function vim.api.nvim_buf_call(buffer: integer, fun: function)
  -> any
```


---

# vim.api.nvim_buf_clear_highlight


```lua
function vim.api.nvim_buf_clear_highlight(buffer: integer, ns_id: integer, line_start: integer, line_end: integer)
```


---

# vim.api.nvim_buf_clear_namespace


```lua
function vim.api.nvim_buf_clear_namespace(buffer: integer, ns_id: integer, line_start: integer, line_end: integer)
```


---

# vim.api.nvim_buf_create_user_command


```lua
function vim.api.nvim_buf_create_user_command(buffer: integer, name: string, command: any, opts: vim.api.keyset.user_command)
```


---

# vim.api.nvim_buf_del_extmark


```lua
function vim.api.nvim_buf_del_extmark(buffer: integer, ns_id: integer, id: integer)
  -> boolean
```


---

# vim.api.nvim_buf_del_keymap


```lua
function vim.api.nvim_buf_del_keymap(buffer: integer, mode: string, lhs: string)
```


---

# vim.api.nvim_buf_del_mark


```lua
function vim.api.nvim_buf_del_mark(buffer: integer, name: string)
  -> boolean
```


---

# vim.api.nvim_buf_del_user_command


```lua
function vim.api.nvim_buf_del_user_command(buffer: integer, name: string)
```


---

# vim.api.nvim_buf_del_var


```lua
function vim.api.nvim_buf_del_var(buffer: integer, name: string)
```


---

# vim.api.nvim_buf_delete


```lua
function vim.api.nvim_buf_delete(buffer: integer, opts: vim.api.keyset.buf_delete)
```


---

# vim.api.nvim_buf_get_changedtick


```lua
function vim.api.nvim_buf_get_changedtick(buffer: integer)
  -> integer
```


---

# vim.api.nvim_buf_get_commands


```lua
function vim.api.nvim_buf_get_commands(buffer: integer, opts: vim.api.keyset.get_commands)
  -> table<string, any>
```


---

# vim.api.nvim_buf_get_extmark_by_id


```lua
function vim.api.nvim_buf_get_extmark_by_id(buffer: integer, ns_id: integer, id: integer, opts: vim.api.keyset.get_extmark)
  -> vim.api.keyset.get_extmark_item_by_id
```


---

# vim.api.nvim_buf_get_extmarks


```lua
function vim.api.nvim_buf_get_extmarks(buffer: integer, ns_id: integer, start: any, end_: any, opts: vim.api.keyset.get_extmarks)
  -> vim.api.keyset.get_extmark_item[]
```


---

# vim.api.nvim_buf_get_keymap


```lua
function vim.api.nvim_buf_get_keymap(buffer: integer, mode: string)
  -> vim.api.keyset.get_keymap[]
```


---

# vim.api.nvim_buf_get_lines


```lua
function vim.api.nvim_buf_get_lines(buffer: integer, start: integer, end_: integer, strict_indexing: boolean)
  -> string[]
```


---

# vim.api.nvim_buf_get_mark


```lua
function vim.api.nvim_buf_get_mark(buffer: integer, name: string)
  -> integer[]
```


---

# vim.api.nvim_buf_get_name


```lua
function vim.api.nvim_buf_get_name(buffer: integer)
  -> string
```


---

# vim.api.nvim_buf_get_number


```lua
function vim.api.nvim_buf_get_number(buffer: integer)
  -> integer
```


---

# vim.api.nvim_buf_get_offset


```lua
function vim.api.nvim_buf_get_offset(buffer: integer, index: integer)
  -> integer
```


---

# vim.api.nvim_buf_get_option


```lua
function vim.api.nvim_buf_get_option(buffer: integer, name: string)
  -> any
```


---

# vim.api.nvim_buf_get_text


```lua
function vim.api.nvim_buf_get_text(buffer: integer, start_row: integer, start_col: integer, end_row: integer, end_col: integer, opts: vim.api.keyset.empty)
  -> string[]
```


---

# vim.api.nvim_buf_get_var


```lua
function vim.api.nvim_buf_get_var(buffer: integer, name: string)
  -> any
```


---

# vim.api.nvim_buf_is_loaded


```lua
function vim.api.nvim_buf_is_loaded(buffer: integer)
  -> boolean
```


---

# vim.api.nvim_buf_is_valid


```lua
function vim.api.nvim_buf_is_valid(buffer: integer)
  -> boolean
```


---

# vim.api.nvim_buf_line_count


```lua
function vim.api.nvim_buf_line_count(buffer: integer)
  -> integer
```


---

# vim.api.nvim_buf_set_extmark


```lua
function vim.api.nvim_buf_set_extmark(buffer: integer, ns_id: integer, line: integer, col: integer, opts: vim.api.keyset.set_extmark)
  -> integer
```


---

# vim.api.nvim_buf_set_keymap


```lua
function vim.api.nvim_buf_set_keymap(buffer: integer, mode: string, lhs: string, rhs: string, opts: vim.api.keyset.keymap)
```


---

# vim.api.nvim_buf_set_lines


```lua
function vim.api.nvim_buf_set_lines(buffer: integer, start: integer, end_: integer, strict_indexing: boolean, replacement: string[])
```


---

# vim.api.nvim_buf_set_mark


```lua
function vim.api.nvim_buf_set_mark(buffer: integer, name: string, line: integer, col: integer, opts: vim.api.keyset.empty)
  -> boolean
```


---

# vim.api.nvim_buf_set_name


```lua
function vim.api.nvim_buf_set_name(buffer: integer, name: string)
```


---

# vim.api.nvim_buf_set_option


```lua
function vim.api.nvim_buf_set_option(buffer: integer, name: string, value: any)
```


---

# vim.api.nvim_buf_set_text


```lua
function vim.api.nvim_buf_set_text(buffer: integer, start_row: integer, start_col: integer, end_row: integer, end_col: integer, replacement: string[])
```


---

# vim.api.nvim_buf_set_var


```lua
function vim.api.nvim_buf_set_var(buffer: integer, name: string, value: any)
```


---

# vim.api.nvim_buf_set_virtual_text


```lua
function vim.api.nvim_buf_set_virtual_text(buffer: integer, src_id: integer, line: integer, chunks: any[], opts: vim.api.keyset.empty)
  -> integer
```


---

# vim.api.nvim_call_dict_function


```lua
function vim.api.nvim_call_dict_function(dict: any, fn: string, args: any[])
  -> any
```


---

# vim.api.nvim_call_function


```lua
function vim.api.nvim_call_function(fn: string, args: any[])
  -> any
```


---

# vim.api.nvim_chan_send


```lua
function vim.api.nvim_chan_send(chan: integer, data: string)
```


---

# vim.api.nvim_clear_autocmds


```lua
function vim.api.nvim_clear_autocmds(opts: vim.api.keyset.clear_autocmds)
```


---

# vim.api.nvim_cmd


```lua
function vim.api.nvim_cmd(cmd: vim.api.keyset.cmd, opts: vim.api.keyset.cmd_opts)
  -> string
```


---

# vim.api.nvim_command


```lua
function vim.api.nvim_command(command: string)
```


---

# vim.api.nvim_command_output


```lua
function vim.api.nvim_command_output(command: string)
  -> string
```


---

# vim.api.nvim_create_augroup


```lua
function vim.api.nvim_create_augroup(name: string, opts: vim.api.keyset.create_augroup)
  -> integer
```


---

# vim.api.nvim_create_autocmd


```lua
function vim.api.nvim_create_autocmd(event: any, opts: vim.api.keyset.create_autocmd)
  -> integer
```


---

# vim.api.nvim_create_buf


```lua
function vim.api.nvim_create_buf(listed: boolean, scratch: boolean)
  -> integer
```


---

# vim.api.nvim_create_namespace


```lua
function vim.api.nvim_create_namespace(name: string)
  -> integer
```


---

# vim.api.nvim_create_user_command


```lua
function vim.api.nvim_create_user_command(name: string, command: string|fun(args: vim.api.keyset.create_user_command.command_args), opts: vim.api.keyset.user_command)
```


---

# vim.api.nvim_del_augroup_by_id


```lua
function vim.api.nvim_del_augroup_by_id(id: integer)
```


---

# vim.api.nvim_del_augroup_by_name


```lua
function vim.api.nvim_del_augroup_by_name(name: string)
```


---

# vim.api.nvim_del_autocmd


```lua
function vim.api.nvim_del_autocmd(id: integer)
```


---

# vim.api.nvim_del_current_line


```lua
function vim.api.nvim_del_current_line()
```


---

# vim.api.nvim_del_keymap


```lua
function vim.api.nvim_del_keymap(mode: string, lhs: string)
```


---

# vim.api.nvim_del_mark


```lua
function vim.api.nvim_del_mark(name: string)
  -> boolean
```


---

# vim.api.nvim_del_user_command


```lua
function vim.api.nvim_del_user_command(name: string)
```


---

# vim.api.nvim_del_var


```lua
function vim.api.nvim_del_var(name: string)
```


---

# vim.api.nvim_echo


```lua
function vim.api.nvim_echo(chunks: any[], history: boolean, opts: vim.api.keyset.echo_opts)
```


---

# vim.api.nvim_err_write


```lua
function vim.api.nvim_err_write(str: string)
```


---

# vim.api.nvim_err_writeln


```lua
function vim.api.nvim_err_writeln(str: string)
```


---

# vim.api.nvim_eval


```lua
function vim.api.nvim_eval(expr: string)
  -> any
```


---

# vim.api.nvim_eval_statusline


```lua
function vim.api.nvim_eval_statusline(str: string, opts: vim.api.keyset.eval_statusline)
  -> these: table<string|Dict>|with
```


---

# vim.api.nvim_exec


```lua
function vim.api.nvim_exec(src: string, output: boolean)
  -> string
```


---

# vim.api.nvim_exec2


```lua
function vim.api.nvim_exec2(src: string, opts: vim.api.keyset.exec_opts)
  -> table<string, any>
```


---

# vim.api.nvim_exec_autocmds


```lua
function vim.api.nvim_exec_autocmds(event: any, opts: vim.api.keyset.exec_autocmds)
```


---

# vim.api.nvim_feedkeys


```lua
function vim.api.nvim_feedkeys(keys: string, mode: string, escape_ks: boolean)
```


---

# vim.api.nvim_get_all_options_info


```lua
function vim.api.nvim_get_all_options_info()
  -> table<string, any>
```


---

# vim.api.nvim_get_autocmds


```lua
function vim.api.nvim_get_autocmds(opts: vim.api.keyset.get_autocmds)
  -> vim.api.keyset.get_autocmds.ret[]
```


---

# vim.api.nvim_get_chan_info


```lua
function vim.api.nvim_get_chan_info(chan: integer)
  -> table<string, any>
```


---

# vim.api.nvim_get_color_by_name


```lua
function vim.api.nvim_get_color_by_name(name: string)
  -> integer
```


---

# vim.api.nvim_get_color_map


```lua
function vim.api.nvim_get_color_map()
  -> table<string, integer>
```


---

# vim.api.nvim_get_commands


```lua
function vim.api.nvim_get_commands(opts: vim.api.keyset.get_commands)
  -> table<string, any>
```


---

# vim.api.nvim_get_context


```lua
function vim.api.nvim_get_context(opts: vim.api.keyset.context)
  -> table<string, any>
```


---

# vim.api.nvim_get_current_buf


```lua
function vim.api.nvim_get_current_buf()
  -> integer
```


---

# vim.api.nvim_get_current_line


```lua
function vim.api.nvim_get_current_line()
  -> string
```


---

# vim.api.nvim_get_current_tabpage


```lua
function vim.api.nvim_get_current_tabpage()
  -> integer
```


---

# vim.api.nvim_get_current_win


```lua
function vim.api.nvim_get_current_win()
  -> integer
```


---

# vim.api.nvim_get_hl


```lua
function vim.api.nvim_get_hl(ns_id: integer, opts: vim.api.keyset.get_highlight)
  -> vim.api.keyset.get_hl_info
```


---

# vim.api.nvim_get_hl_by_id


```lua
function vim.api.nvim_get_hl_by_id(hl_id: integer, rgb: boolean)
  -> table<string, any>
```


---

# vim.api.nvim_get_hl_by_name


```lua
function vim.api.nvim_get_hl_by_name(name: string, rgb: boolean)
  -> table<string, any>
```


---

# vim.api.nvim_get_hl_id_by_name


```lua
function vim.api.nvim_get_hl_id_by_name(name: string)
  -> integer
```


---

# vim.api.nvim_get_hl_ns


```lua
function vim.api.nvim_get_hl_ns(opts: vim.api.keyset.get_ns)
  -> integer
```


---

# vim.api.nvim_get_keymap


```lua
function vim.api.nvim_get_keymap(mode: string)
  -> vim.api.keyset.get_keymap[]
```


---

# vim.api.nvim_get_mark


```lua
function vim.api.nvim_get_mark(name: string, opts: vim.api.keyset.empty)
  -> vim.api.keyset.get_mark
```


---

# vim.api.nvim_get_mode


```lua
function vim.api.nvim_get_mode()
  -> vim.api.keyset.get_mode
```


---

# vim.api.nvim_get_namespaces


```lua
function vim.api.nvim_get_namespaces()
  -> table<string, integer>
```


---

# vim.api.nvim_get_option


```lua
function vim.api.nvim_get_option(name: string)
  -> any
```


---

# vim.api.nvim_get_option_info


```lua
function vim.api.nvim_get_option_info(name: string)
  -> vim.api.keyset.get_option_info
```


---

# vim.api.nvim_get_option_info2


```lua
function vim.api.nvim_get_option_info2(name: string, opts: vim.api.keyset.option)
  -> vim.api.keyset.get_option_info
```


---

# vim.api.nvim_get_option_value


```lua
function vim.api.nvim_get_option_value(name: string, opts: vim.api.keyset.option)
  -> any
```


---

# vim.api.nvim_get_proc


```lua
function vim.api.nvim_get_proc(pid: integer)
  -> any
```


---

# vim.api.nvim_get_proc_children


```lua
function vim.api.nvim_get_proc_children(pid: integer)
  -> any[]
```


---

# vim.api.nvim_get_runtime_file


```lua
function vim.api.nvim_get_runtime_file(name: string, all: boolean)
  -> string[]
```


---

# vim.api.nvim_get_var


```lua
function vim.api.nvim_get_var(name: string)
  -> any
```


---

# vim.api.nvim_get_vvar


```lua
function vim.api.nvim_get_vvar(name: string)
  -> any
```


---

# vim.api.nvim_input


```lua
function vim.api.nvim_input(keys: string)
  -> integer
```


---

# vim.api.nvim_input_mouse


```lua
function vim.api.nvim_input_mouse(button: string, action: string, modifier: string, grid: integer, row: integer, col: integer)
```


---

# vim.api.nvim_list_bufs


```lua
function vim.api.nvim_list_bufs()
  -> integer[]
```


---

# vim.api.nvim_list_chans


```lua
function vim.api.nvim_list_chans()
  -> any[]
```


---

# vim.api.nvim_list_runtime_paths


```lua
function vim.api.nvim_list_runtime_paths()
  -> string[]
```


---

# vim.api.nvim_list_tabpages


```lua
function vim.api.nvim_list_tabpages()
  -> integer[]
```


---

# vim.api.nvim_list_uis


```lua
function vim.api.nvim_list_uis()
  -> any[]
```


---

# vim.api.nvim_list_wins


```lua
function vim.api.nvim_list_wins()
  -> integer[]
```


---

# vim.api.nvim_load_context


```lua
function vim.api.nvim_load_context(dict: table<string, any>)
  -> any
```


---

# vim.api.nvim_notify


```lua
function vim.api.nvim_notify(msg: string, log_level: integer, opts: table<string, any>)
  -> any
```


---

# vim.api.nvim_open_term


```lua
function vim.api.nvim_open_term(buffer: integer, opts: vim.api.keyset.open_term)
  -> integer
```


---

# vim.api.nvim_open_win


```lua
function vim.api.nvim_open_win(buffer: integer, enter: boolean, config: vim.api.keyset.win_config)
  -> integer
```


---

# vim.api.nvim_out_write


```lua
function vim.api.nvim_out_write(str: string)
```


---

# vim.api.nvim_parse_cmd


```lua
function vim.api.nvim_parse_cmd(str: string, opts: vim.api.keyset.empty)
  -> vim.api.keyset.parse_cmd
```


---

# vim.api.nvim_parse_expression


```lua
function vim.api.nvim_parse_expression(expr: string, flags: string, highlight: boolean)
  -> table<string, any>
```


---

# vim.api.nvim_paste


```lua
function vim.api.nvim_paste(data: string, crlf: boolean, phase: integer)
  -> boolean
```


---

# vim.api.nvim_put


```lua
function vim.api.nvim_put(lines: string[], type: string, after: boolean, follow: boolean)
```


---

# vim.api.nvim_replace_termcodes


```lua
function vim.api.nvim_replace_termcodes(str: string, from_part: boolean, do_lt: boolean, special: boolean)
  -> string
```


---

# vim.api.nvim_select_popupmenu_item


```lua
function vim.api.nvim_select_popupmenu_item(item: integer, insert: boolean, finish: boolean, opts: vim.api.keyset.empty)
```


---

# vim.api.nvim_set_current_buf


```lua
function vim.api.nvim_set_current_buf(buffer: integer)
```


---

# vim.api.nvim_set_current_dir


```lua
function vim.api.nvim_set_current_dir(dir: string)
```


---

# vim.api.nvim_set_current_line


```lua
function vim.api.nvim_set_current_line(line: string)
```


---

# vim.api.nvim_set_current_tabpage


```lua
function vim.api.nvim_set_current_tabpage(tabpage: integer)
```


---

# vim.api.nvim_set_current_win


```lua
function vim.api.nvim_set_current_win(window: integer)
```


---

# vim.api.nvim_set_decoration_provider


```lua
function vim.api.nvim_set_decoration_provider(ns_id: integer, opts: vim.api.keyset.set_decoration_provider)
```


---

# vim.api.nvim_set_hl


```lua
function vim.api.nvim_set_hl(ns_id: integer, name: string, val: vim.api.keyset.highlight)
```


---

# vim.api.nvim_set_hl_ns


```lua
function vim.api.nvim_set_hl_ns(ns_id: integer)
```


---

# vim.api.nvim_set_hl_ns_fast


```lua
function vim.api.nvim_set_hl_ns_fast(ns_id: integer)
```


---

# vim.api.nvim_set_keymap


```lua
function vim.api.nvim_set_keymap(mode: string, lhs: string, rhs: string, opts: vim.api.keyset.keymap)
```


---

# vim.api.nvim_set_option


```lua
function vim.api.nvim_set_option(name: string, value: any)
```


---

# vim.api.nvim_set_option_value


```lua
function vim.api.nvim_set_option_value(name: string, value: any, opts: vim.api.keyset.option)
```


---

# vim.api.nvim_set_var


```lua
function vim.api.nvim_set_var(name: string, value: any)
```


---

# vim.api.nvim_set_vvar


```lua
function vim.api.nvim_set_vvar(name: string, value: any)
```


---

# vim.api.nvim_strwidth


```lua
function vim.api.nvim_strwidth(text: string)
  -> integer
```


---

# vim.api.nvim_tabpage_del_var


```lua
function vim.api.nvim_tabpage_del_var(tabpage: integer, name: string)
```


---

# vim.api.nvim_tabpage_get_number


```lua
function vim.api.nvim_tabpage_get_number(tabpage: integer)
  -> integer
```


---

# vim.api.nvim_tabpage_get_var


```lua
function vim.api.nvim_tabpage_get_var(tabpage: integer, name: string)
  -> any
```


---

# vim.api.nvim_tabpage_get_win


```lua
function vim.api.nvim_tabpage_get_win(tabpage: integer)
  -> integer
```


---

# vim.api.nvim_tabpage_is_valid


```lua
function vim.api.nvim_tabpage_is_valid(tabpage: integer)
  -> boolean
```


---

# vim.api.nvim_tabpage_list_wins


```lua
function vim.api.nvim_tabpage_list_wins(tabpage: integer)
  -> integer[]
```


---

# vim.api.nvim_tabpage_set_var


```lua
function vim.api.nvim_tabpage_set_var(tabpage: integer, name: string, value: any)
```


---

# vim.api.nvim_tabpage_set_win


```lua
function vim.api.nvim_tabpage_set_win(tabpage: integer, win: integer)
```


---

# vim.api.nvim_win_call


```lua
function vim.api.nvim_win_call(window: integer, fun: function)
  -> any
```


---

# vim.api.nvim_win_close


```lua
function vim.api.nvim_win_close(window: integer, force: boolean)
```


---

# vim.api.nvim_win_del_var


```lua
function vim.api.nvim_win_del_var(window: integer, name: string)
```


---

# vim.api.nvim_win_get_buf


```lua
function vim.api.nvim_win_get_buf(window: integer)
  -> integer
```


---

# vim.api.nvim_win_get_config


```lua
function vim.api.nvim_win_get_config(window: integer)
  -> vim.api.keyset.win_config
```


---

# vim.api.nvim_win_get_cursor


```lua
function vim.api.nvim_win_get_cursor(window: integer)
  -> integer[]
```


---

# vim.api.nvim_win_get_height


```lua
function vim.api.nvim_win_get_height(window: integer)
  -> integer
```


---

# vim.api.nvim_win_get_number


```lua
function vim.api.nvim_win_get_number(window: integer)
  -> integer
```


---

# vim.api.nvim_win_get_option


```lua
function vim.api.nvim_win_get_option(window: integer, name: string)
  -> any
```


---

# vim.api.nvim_win_get_position


```lua
function vim.api.nvim_win_get_position(window: integer)
  -> integer[]
```


---

# vim.api.nvim_win_get_tabpage


```lua
function vim.api.nvim_win_get_tabpage(window: integer)
  -> integer
```


---

# vim.api.nvim_win_get_var


```lua
function vim.api.nvim_win_get_var(window: integer, name: string)
  -> any
```


---

# vim.api.nvim_win_get_width


```lua
function vim.api.nvim_win_get_width(window: integer)
  -> integer
```


---

# vim.api.nvim_win_hide


```lua
function vim.api.nvim_win_hide(window: integer)
```


---

# vim.api.nvim_win_is_valid


```lua
function vim.api.nvim_win_is_valid(window: integer)
  -> boolean
```


---

# vim.api.nvim_win_set_buf


```lua
function vim.api.nvim_win_set_buf(window: integer, buffer: integer)
```


---

# vim.api.nvim_win_set_config


```lua
function vim.api.nvim_win_set_config(window: integer, config: vim.api.keyset.win_config)
```


---

# vim.api.nvim_win_set_cursor


```lua
function vim.api.nvim_win_set_cursor(window: integer, pos: integer[])
```


---

# vim.api.nvim_win_set_height


```lua
function vim.api.nvim_win_set_height(window: integer, height: integer)
```


---

# vim.api.nvim_win_set_hl_ns


```lua
function vim.api.nvim_win_set_hl_ns(window: integer, ns_id: integer)
```


---

# vim.api.nvim_win_set_option


```lua
function vim.api.nvim_win_set_option(window: integer, name: string, value: any)
```


---

# vim.api.nvim_win_set_var


```lua
function vim.api.nvim_win_set_var(window: integer, name: string, value: any)
```


---

# vim.api.nvim_win_set_width


```lua
function vim.api.nvim_win_set_width(window: integer, width: integer)
```


---

# vim.api.nvim_win_text_height


```lua
function vim.api.nvim_win_text_height(window: integer, opts: vim.api.keyset.win_text_height)
  -> table<string, any>
```


---

# vim.b


```lua
vim.var_accessor
```


---

# vim.b.make_microsoft


```lua
nil
```


```lua
integer
```


---

# vim.b.man_sect


```lua
string
```


```lua
string
```


```lua
string
```


---

# vim.b.tutor_extmarks


```lua
any
```


```lua
table
```


```lua
table<string, string>
```


---

# vim.b.tutor_metadata


```lua
any
```


---

# vim.b.undo_ftplugin


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


---

# vim.base64.decode


```lua
function vim.base64.decode(str: string)
  -> string
```


---

# vim.base64.encode


```lua
function vim.base64.encode(str: string)
  -> string
```


---

# vim.bo

## [integer]


```lua
vim.bo
```

## ai


```lua
boolean
```

## ar


```lua
boolean
```

## autoindent


```lua
boolean
```

## autoread


```lua
boolean
```

## backupcopy


```lua
string
```

## bh


```lua
string|''|'delete'|'hide'|'unload'...(+1)
```

## bin


```lua
boolean
```

## binary


```lua
boolean
```

## bkc


```lua
string
```

## bl


```lua
boolean
```

## bomb


```lua
boolean
```

## bt


```lua
string|''|'acwrite'|'help'|'nofile'...(+4)
```

## bufhidden


```lua
string
```

## buflisted


```lua
boolean
```

## buftype


```lua
string
```

## cfu


```lua
string
```

## channel


```lua
integer
```

## ci


```lua
boolean
```

## cin


```lua
boolean
```

## cindent


```lua
boolean
```

## cink


```lua
string
```

## cinkeys


```lua
string
```

## cino


```lua
string
```

## cinoptions


```lua
string
```

## cinscopedecls


```lua
string
```

## cinsd


```lua
string
```

## cinw


```lua
string
```

## cinwords


```lua
string
```

## cms


```lua
string
```

## com


```lua
string
```

## comments


```lua
string
```

## commentstring


```lua
string
```

## complete


```lua
string
```

## completefunc


```lua
string
```

## completeopt


```lua
string
```

## completeslash


```lua
''|'backslash'|'slash'
```

## copyindent


```lua
boolean
```

## cot


```lua
string
```

## cpt


```lua
string
```

## csl


```lua
''|'backslash'|'slash'
```

## def


```lua
string
```

## define


```lua
string
```

## dict


```lua
string
```

## dictionary


```lua
string
```

## efm


```lua
string
```

## endoffile


```lua
boolean
```

## endofline


```lua
boolean
```

## eof


```lua
boolean
```

## eol


```lua
boolean
```

## ep


```lua
string
```

## equalprg


```lua
string
```

## errorformat


```lua
string
```

## et


```lua
boolean
```

## expandtab


```lua
boolean
```

## fenc


```lua
string
```

## fex


```lua
string
```

## ff


```lua
'dos'|'mac'|'unix'
```

## ffu


```lua
string
```

## fileencoding


```lua
string
```

## fileformat


```lua
'dos'|'mac'|'unix'
```

## filetype


```lua
string
```

## findfunc


```lua
string
```

## fixendofline


```lua
boolean
```

## fixeol


```lua
boolean
```

## flp


```lua
string
```

## fo


```lua
string
```

## formatexpr


```lua
string
```

## formatlistpat


```lua
string
```

## formatoptions


```lua
string
```

## formatprg


```lua
string
```

## fp


```lua
string
```

## ft


```lua
string
```

## gp


```lua
string
```

## grepprg


```lua
string
```

## imi


```lua
integer
```

## iminsert


```lua
integer
```

## ims


```lua
integer
```

## imsearch


```lua
integer
```

## inc


```lua
string
```

## include


```lua
string
```

## includeexpr


```lua
string
```

## inde


```lua
string
```

## indentexpr


```lua
string
```

## indentkeys


```lua
string
```

## indk


```lua
string
```

## inex


```lua
string
```

## inf


```lua
boolean
```

## infercase


```lua
boolean
```

## isk


```lua
string
```

## iskeyword


```lua
string
```

## keymap


```lua
string
```

## keywordprg


```lua
string
```

## kmp


```lua
string
```

## kp


```lua
string
```

## lisp


```lua
boolean
```

## lispoptions


```lua
string
```

## lispwords


```lua
string
```

## lop


```lua
string
```

## lw


```lua
string
```

## ma


```lua
boolean
```

## makeencoding


```lua
string
```

## makeprg


```lua
string
```

## matchpairs


```lua
string
```

## menc


```lua
string
```

## ml


```lua
boolean
```

## mod


```lua
boolean
```

## modeline


```lua
boolean
```

## modifiable


```lua
boolean
```

## modified


```lua
boolean
```

## mp


```lua
string
```

## mps


```lua
string
```

## nf


```lua
string
```

## nrformats


```lua
string
```

## ofu


```lua
string
```

## omnifunc


```lua
string
```

## pa


```lua
string
```

## path


```lua
string
```

## pi


```lua
boolean
```

## preserveindent


```lua
boolean
```

## qe


```lua
string
```

## quoteescape


```lua
string
```

## readonly


```lua
boolean
```

## ro


```lua
boolean
```

## scbk


```lua
integer
```

## scrollback


```lua
integer
```

## shiftwidth


```lua
integer
```

## si


```lua
boolean
```

## smartindent


```lua
boolean
```

## smc


```lua
integer
```

## softtabstop


```lua
integer
```

## spc


```lua
string
```

## spellcapcheck


```lua
string
```

## spellfile


```lua
string
```

## spelllang


```lua
string
```

## spelloptions


```lua
string
```

## spf


```lua
string
```

## spl


```lua
string
```

## spo


```lua
string
```

## sts


```lua
integer
```

## sua


```lua
string
```

## suffixesadd


```lua
string
```

## sw


```lua
integer
```

## swapfile


```lua
boolean
```

## swf


```lua
boolean
```

## syn


```lua
string
```

## synmaxcol


```lua
integer
```

## syntax


```lua
string
```

## tabstop


```lua
integer
```

## tag


```lua
string
```

## tagcase


```lua
'followic'|'followscs'|'ignore'|'match'|'smart'
```

## tagfunc


```lua
string
```

## tags


```lua
string
```

## tc


```lua
'followic'|'followscs'|'ignore'|'match'|'smart'
```

## textwidth


```lua
integer
```

## tfu


```lua
string
```

## thesaurus


```lua
string
```

## thesaurusfunc


```lua
string
```

## ts


```lua
integer
```

## tsr


```lua
string
```

## tsrfu


```lua
string
```

## tw


```lua
integer
```

## udf


```lua
boolean
```

## ul


```lua
integer
```

## undofile


```lua
boolean
```

## undolevels


```lua
integer
```

## varsofttabstop


```lua
string
```

## vartabstop


```lua
string
```

## vsts


```lua
string
```

## vts


```lua
string
```

## wm


```lua
integer
```

## wrapmargin


```lua
integer
```


---

# vim.bo


```lua
table|vim.bo
```


```lua
table
```


---

# vim.bo.ai


```lua
boolean
```


---

# vim.bo.ar


```lua
boolean
```


---

# vim.bo.autoindent


```lua
boolean
```


---

# vim.bo.autoread


```lua
boolean
```


---

# vim.bo.backupcopy


```lua
string
```


---

# vim.bo.bh


```lua
string|''|'delete'|'hide'|'unload'...(+1)
```


---

# vim.bo.bin


```lua
boolean
```


---

# vim.bo.binary


```lua
boolean
```


---

# vim.bo.bkc


```lua
string
```


---

# vim.bo.bl


```lua
boolean
```


---

# vim.bo.bomb


```lua
boolean
```


---

# vim.bo.bt


```lua
string|''|'acwrite'|'help'|'nofile'...(+4)
```


---

# vim.bo.bufhidden


```lua
string
```


```lua
''|'delete'|'hide'|'unload'|'wipe'
```


---

# vim.bo.buflisted


```lua
boolean
```


---

# vim.bo.buftype


```lua
string
```


```lua
''|'acwrite'|'help'|'nofile'|'nowrite'...(+3)
```


---

# vim.bo.cfu


```lua
string
```


---

# vim.bo.channel


```lua
integer
```


---

# vim.bo.ci


```lua
boolean
```


---

# vim.bo.cin


```lua
boolean
```


---

# vim.bo.cindent


```lua
boolean
```


---

# vim.bo.cink


```lua
string
```


---

# vim.bo.cinkeys


```lua
string
```


---

# vim.bo.cino


```lua
string
```


---

# vim.bo.cinoptions


```lua
string
```


---

# vim.bo.cinscopedecls


```lua
string
```


---

# vim.bo.cinsd


```lua
string
```


---

# vim.bo.cinw


```lua
string
```


---

# vim.bo.cinwords


```lua
string
```


---

# vim.bo.cms


```lua
string
```


---

# vim.bo.com


```lua
string
```


---

# vim.bo.comments


```lua
string
```


---

# vim.bo.commentstring


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


---

# vim.bo.complete


```lua
string
```


---

# vim.bo.completefunc


```lua
string
```


---

# vim.bo.completeopt


```lua
string
```


---

# vim.bo.completeslash


```lua
''|'backslash'|'slash'
```


---

# vim.bo.copyindent


```lua
boolean
```


---

# vim.bo.cot


```lua
string
```


---

# vim.bo.cpt


```lua
string
```


---

# vim.bo.csl


```lua
''|'backslash'|'slash'
```


---

# vim.bo.def


```lua
string
```


---

# vim.bo.define


```lua
string
```


```lua
string
```


---

# vim.bo.dict


```lua
string
```


---

# vim.bo.dictionary


```lua
string
```


---

# vim.bo.efm


```lua
string
```


---

# vim.bo.endoffile


```lua
boolean
```


---

# vim.bo.endofline


```lua
boolean
```


---

# vim.bo.eof


```lua
boolean
```


---

# vim.bo.eol


```lua
boolean
```


---

# vim.bo.ep


```lua
string
```


---

# vim.bo.equalprg


```lua
string
```


---

# vim.bo.errorformat


```lua
string
```


---

# vim.bo.et


```lua
boolean
```


---

# vim.bo.expandtab


```lua
boolean
```


---

# vim.bo.fenc


```lua
string
```


---

# vim.bo.fex


```lua
string
```


---

# vim.bo.ff


```lua
'dos'|'mac'|'unix'
```


---

# vim.bo.ffu


```lua
string
```


---

# vim.bo.fileencoding


```lua
string
```


---

# vim.bo.fileformat


```lua
'dos'|'mac'|'unix'
```


---

# vim.bo.filetype


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


```lua
string
```


---

# vim.bo.findfunc


```lua
string
```


---

# vim.bo.fixendofline


```lua
boolean
```


---

# vim.bo.fixeol


```lua
boolean
```


---

# vim.bo.flp


```lua
string
```


---

# vim.bo.fo


```lua
string
```


---

# vim.bo.formatexpr


```lua
string
```


---

# vim.bo.formatlistpat


```lua
string
```


---

# vim.bo.formatoptions


```lua
string
```


---

# vim.bo.formatprg


```lua
string
```


---

# vim.bo.fp


```lua
string
```


---

# vim.bo.ft


```lua
string
```


---

# vim.bo.gp


```lua
string
```


---

# vim.bo.grepprg


```lua
string
```


---

# vim.bo.imi


```lua
integer
```


---

# vim.bo.iminsert


```lua
integer
```


---

# vim.bo.ims


```lua
integer
```


---

# vim.bo.imsearch


```lua
integer
```


---

# vim.bo.inc


```lua
string
```


---

# vim.bo.include


```lua
string
```


```lua
string
```


---

# vim.bo.includeexpr


```lua
string
```


```lua
string
```


---

# vim.bo.inde


```lua
string
```


---

# vim.bo.indentexpr


```lua
string
```


---

# vim.bo.indentkeys


```lua
string
```


---

# vim.bo.indk


```lua
string
```


---

# vim.bo.inex


```lua
string
```


---

# vim.bo.inf


```lua
boolean
```


---

# vim.bo.infercase


```lua
boolean
```


---

# vim.bo.isk


```lua
string
```


---

# vim.bo.iskeyword


```lua
string
```


---

# vim.bo.keymap


```lua
string
```


---

# vim.bo.keywordprg


```lua
string
```


---

# vim.bo.kmp


```lua
string
```


---

# vim.bo.kp


```lua
string
```


---

# vim.bo.lisp


```lua
boolean
```


---

# vim.bo.lispoptions


```lua
string
```


---

# vim.bo.lispwords


```lua
string
```


---

# vim.bo.lop


```lua
string
```


---

# vim.bo.lw


```lua
string
```


---

# vim.bo.ma


```lua
boolean
```


---

# vim.bo.makeencoding


```lua
string
```


---

# vim.bo.makeprg


```lua
string
```


---

# vim.bo.matchpairs


```lua
string
```


---

# vim.bo.menc


```lua
string
```


---

# vim.bo.ml


```lua
boolean
```


---

# vim.bo.mod


```lua
boolean
```


---

# vim.bo.modeline


```lua
boolean
```


---

# vim.bo.modifiable


```lua
boolean
```


```lua
boolean
```


```lua
boolean
```


```lua
boolean
```


```lua
boolean
```


---

# vim.bo.modified


```lua
boolean
```


```lua
boolean
```


```lua
boolean
```


---

# vim.bo.mp


```lua
string
```


---

# vim.bo.mps


```lua
string
```


---

# vim.bo.nf


```lua
string
```


---

# vim.bo.nrformats


```lua
string
```


---

# vim.bo.ofu


```lua
string
```


---

# vim.bo.omnifunc


```lua
string
```


```lua
string
```


```lua
string
```


---

# vim.bo.pa


```lua
string
```


---

# vim.bo.path


```lua
string
```


---

# vim.bo.pi


```lua
boolean
```


---

# vim.bo.preserveindent


```lua
boolean
```


---

# vim.bo.qe


```lua
string
```


---

# vim.bo.quoteescape


```lua
string
```


---

# vim.bo.readonly


```lua
boolean
```


```lua
boolean
```


```lua
boolean
```


---

# vim.bo.ro


```lua
boolean
```


---

# vim.bo.scbk


```lua
integer
```


---

# vim.bo.scrollback


```lua
integer
```


---

# vim.bo.shiftwidth


```lua
integer
```


---

# vim.bo.si


```lua
boolean
```


---

# vim.bo.smartindent


```lua
boolean
```


---

# vim.bo.smc


```lua
integer
```


---

# vim.bo.softtabstop


```lua
integer
```


---

# vim.bo.spc


```lua
string
```


---

# vim.bo.spellcapcheck


```lua
string
```


---

# vim.bo.spellfile


```lua
string
```


---

# vim.bo.spelllang


```lua
string
```


---

# vim.bo.spelloptions


```lua
string
```


---

# vim.bo.spf


```lua
string
```


---

# vim.bo.spl


```lua
string
```


---

# vim.bo.spo


```lua
string
```


---

# vim.bo.sts


```lua
integer
```


---

# vim.bo.sua


```lua
string
```


---

# vim.bo.suffixesadd


```lua
string
```


---

# vim.bo.sw


```lua
integer
```


---

# vim.bo.swapfile


```lua
boolean
```


```lua
boolean
```


```lua
boolean
```


---

# vim.bo.swf


```lua
boolean
```


---

# vim.bo.syn


```lua
string
```


---

# vim.bo.synmaxcol


```lua
integer
```


---

# vim.bo.syntax


```lua
string
```


---

# vim.bo.tabstop


```lua
integer
```


---

# vim.bo.tag


```lua
string
```


---

# vim.bo.tagcase


```lua
'followic'|'followscs'|'ignore'|'match'|'smart'
```


---

# vim.bo.tagfunc


```lua
string
```


---

# vim.bo.tags


```lua
string
```


---

# vim.bo.tc


```lua
'followic'|'followscs'|'ignore'|'match'|'smart'
```


---

# vim.bo.textwidth


```lua
integer
```


---

# vim.bo.tfu


```lua
string
```


---

# vim.bo.thesaurus


```lua
string
```


---

# vim.bo.thesaurusfunc


```lua
string
```


---

# vim.bo.ts


```lua
integer
```


---

# vim.bo.tsr


```lua
string
```


---

# vim.bo.tsrfu


```lua
string
```


---

# vim.bo.tw


```lua
integer
```


---

# vim.bo.udf


```lua
boolean
```


---

# vim.bo.ul


```lua
integer
```


---

# vim.bo.undofile


```lua
boolean
```


---

# vim.bo.undolevels


```lua
integer
```


---

# vim.bo.varsofttabstop


```lua
string
```


---

# vim.bo.vartabstop


```lua
string
```


---

# vim.bo.vsts


```lua
string
```


---

# vim.bo.vts


```lua
string
```


---

# vim.bo.wm


```lua
integer
```


---

# vim.bo.wrapmargin


```lua
integer
```


---

# vim.call


```lua
function vim.call(func: string, ...any)
  -> any
```


---

# vim.cmd


```lua
table
```


---

# vim.context.mods

 @nodoc

## bo


```lua
table<string, any>?
```

## buf


```lua
integer?
```

## emsg_silent


```lua
boolean?
```

## env


```lua
table<string, any>?
```

## go


```lua
table<string, any>?
```

## hide


```lua
boolean?
```

## keepalt


```lua
boolean?
```

## keepjumps


```lua
boolean?
```

## keepmarks


```lua
boolean?
```

## keeppatterns


```lua
boolean?
```

## lockmarks


```lua
boolean?
```

## noautocmd


```lua
boolean?
```

## o


```lua
table<string, any>?
```

## sandbox


```lua
boolean?
```

## silent


```lua
boolean?
```

## unsilent


```lua
boolean?
```

## win


```lua
integer?
```

## wo


```lua
table<string, any>?
```


---

# vim.context.state

 @nodoc

## bo


```lua
table<string, any>?
```

## env


```lua
table<string, any>?
```

## go


```lua
table<string, any>?
```

## wo


```lua
table<string, any>?
```


---

# vim.deep_equal


```lua
function vim.deep_equal(a: any, b: any)
  -> boolean
```


---

# vim.deepcopy


```lua
function vim.deepcopy(orig: <T:table>, noref?: boolean)
  -> Table: <T:table>
```


---

# vim.defaulttable


```lua
function vim.defaulttable(createfn?: fun(key: any):any)
  -> table
```


---

# vim.defer_fn


```lua
function vim.defer_fn(fn: function, timeout: integer)
  -> timer: table
```


---

# vim.deprecate


```lua
function vim.deprecate(name: string, alternative: string|nil, version: string, plugin: string|nil, backtrace: boolean|nil)
  -> string|nil
```


---

# vim.diagnostic


```lua
table
```


---

# vim.diagnostic.Filter

 TODO: inherit from `vim.diagnostic.Opts`, implement its fields.
 Optional filters |kwargs|, or `nil` for all.

## bufnr


```lua
integer?
```


 Buffer number, or 0 for current buffer, or `nil` for all buffers.

## ns_id


```lua
integer?
```

 @inlinedoc

 Diagnostic namespace, or `nil` for all.


---

# vim.diagnostic.GetOpts

 A table with the following keys:

## lnum


```lua
integer?
```


 Limit diagnostics to those spanning the specified line number.

## namespace


```lua
(integer|integer[])?
```


 Limit diagnostics to one or more namespaces.

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 See |diagnostic-severity|.


---

# vim.diagnostic.Handler

## hide


```lua
fun(namespace: integer, bufnr: integer)?
```

## show


```lua
fun(namespace: integer, bufnr: integer, diagnostics: vim.Diagnostic[], opts?: vim.diagnostic.OptsResolved)?
```


---

# vim.diagnostic.JumpOpts

 Configuration table with the keys listed below. Some parameters can have their default values
 changed with |vim.diagnostic.config()|.

## _highest


```lua
boolean?
```


 Go to the diagnostic with the highest severity.
 (default: `false`)

## count


```lua
integer?
```


 The number of diagnostics to move by, starting from {pos}. A positive
 integer moves forward by {count} diagnostics, while a negative integer moves
 backward by {count} diagnostics. Mutually exclusive with {diagnostic}.

## diagnostic


```lua
(vim.Diagnostic)?
```


 The diagnostic to jump to. Mutually exclusive with {count}, {namespace},
 and {severity}.

## float


```lua
(boolean|vim.diagnostic.Opts.Float)?
```


 If `true`, call |vim.diagnostic.open_float()| after moving.
 If a table, pass the table as the {opts} parameter to |vim.diagnostic.open_float()|.
 Unless overridden, the float will show diagnostics at the new cursor
 position (as if "cursor" were passed to the "scope" option).
 (default: `false`)

## lnum


```lua
integer?
```


 Limit diagnostics to those spanning the specified line number.

## namespace


```lua
(integer|integer[])?
```


 Limit diagnostics to one or more namespaces.

## pos


```lua
[integer, integer]?
```


 Cursor position as a `(row, col)` tuple. See |nvim_win_get_cursor()|. Used
 to find the nearest diagnostic when {count} is used. Only used when {count}
 is non-nil. Default is the current cursor position.

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 See |diagnostic-severity|.

## winid


```lua
integer?
```


 Window ID
 (default: `0`)

## wrap


```lua
boolean?
```


 Whether to loop around file or not. Similar to 'wrapscan'.
 (default: `true`)


---

# vim.diagnostic.NS

## disabled


```lua
boolean?
```

## name


```lua
string
```

## opts


```lua
vim.diagnostic.Opts
```

 Many of the configuration options below accept one of the following:
 - `false`: Disable this feature
 - `true`: Enable this feature, use default settings.
 - `table`: Enable this feature with overrides. Use an empty table to use default values.
 - `function`: Function with signature (namespace, bufnr) that returns any of the above.

## user_data


```lua
table
```


---

# vim.diagnostic.Opts

 Many of the configuration options below accept one of the following:
 - `false`: Disable this feature
 - `true`: Enable this feature, use default settings.
 - `table`: Enable this feature with overrides. Use an empty table to use default values.
 - `function`: Function with signature (namespace, bufnr) that returns any of the above.

## float


```lua
(boolean|fun(namespace: integer, bufnr: integer):vim.diagnostic.Opts.Float|vim.diagnostic.Opts.Float)?
```


 Options for floating windows. See |vim.diagnostic.Opts.Float|.

## jump


```lua
(vim.diagnostic.Opts.Jump)?
```


 Default values for |vim.diagnostic.jump()|. See |vim.diagnostic.Opts.Jump|.

## severity_sort


```lua
(boolean|{ reverse: boolean })?
```


 Sort diagnostics by severity. This affects the order in which signs,
 virtual text, and highlights are displayed. When true, higher severities are
 displayed before lower severities (e.g. ERROR is displayed before WARN).
 Options:
   - {reverse}? (boolean) Reverse sort order
 (default: `false`)

## signs


```lua
(boolean|fun(namespace: integer, bufnr: integer):vim.diagnostic.Opts.Signs|vim.diagnostic.Opts.Signs)?
```


 Use signs for diagnostics |diagnostic-signs|.
 (default: `true`)

## underline


```lua
(boolean|fun(namespace: integer, bufnr: integer):vim.diagnostic.Opts.Underline|vim.diagnostic.Opts.Underline)?
```


 Use underline for diagnostics.
 (default: `true`)

## update_in_insert


```lua
boolean?
```


 Update diagnostics in Insert mode
 (if `false`, diagnostics are updated on |InsertLeave|)
 (default: `false`)

## virtual_lines


```lua
(boolean|fun(namespace: integer, bufnr: integer):vim.diagnostic.Opts.VirtualLines|vim.diagnostic.Opts.VirtualLines)?
```


 Use virtual lines for diagnostics.
 (default: `false`)

## virtual_text


```lua
(boolean|fun(namespace: integer, bufnr: integer):vim.diagnostic.Opts.VirtualText|vim.diagnostic.Opts.VirtualText)?
```


 Use virtual text for diagnostics. If multiple diagnostics are set for a
 namespace, one prefix per diagnostic + the last diagnostic message are
 shown.
 (default: `false`)


---

# vim.diagnostic.Opts.Float

## border


```lua
string?
```

see |nvim_open_win()|.

## bufnr


```lua
integer?
```


 Buffer number to show diagnostics from.
 (default: current buffer)

## focus_id


```lua
string?
```


## format


```lua
fun(diagnostic: vim.Diagnostic):string??
```


 A function that takes a diagnostic as input and returns a string or nil.
 If the return value is nil, the diagnostic is not displayed by the handler.
 Else the output text is used to display the diagnostic.
 Overrides the setting from |vim.diagnostic.config()|.

## header


```lua
(string|[string, any])?
```


 String to use as the header for the floating window. If a table, it is
 interpreted as a `[text, hl_group]` tuple.
 Overrides the setting from |vim.diagnostic.config()|.

## namespace


```lua
integer?
```


 Limit diagnostics to the given namespace

## pos


```lua
(integer|[integer, integer])?
```


 If {scope} is "line" or "cursor", use this position rather than the cursor
 position. If a number, interpreted as a line number; otherwise, a
 (row, col) tuple.

## prefix


```lua
(string|table|fun(diagnostic: vim.Diagnostic, i: integer, total: integer):string, string)?
```


 Prefix each diagnostic in the floating window:
 - If a `function`, {i} is the index of the diagnostic being evaluated and
   {total} is the total number of diagnostics displayed in the window. The
   function should return a `string` which is prepended to each diagnostic
   in the window as well as an (optional) highlight group which will be
   used to highlight the prefix.
 - If a `table`, it is interpreted as a `[text, hl_group]` tuple as
   in |nvim_echo()|
 - If a `string`, it is prepended to each diagnostic in the window with no
   highlight.
 Overrides the setting from |vim.diagnostic.config()|.

## scope


```lua
('b'|'buffer'|'c'|'cursor'|'l'...(+1))?
```


 Show diagnostics from the whole buffer (`buffer"`, the current cursor line
 (`line`), or the current cursor position (`cursor`). Shorthand versions
 are also accepted (`c` for `cursor`, `l` for `line`, `b` for `buffer`).
 (default: `line`)

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 See |diagnostic-severity|.
 Overrides the setting from |vim.diagnostic.config()|.

## severity_sort


```lua
(boolean|{ reverse: boolean })?
```


 Sort diagnostics by severity.
 Overrides the setting from |vim.diagnostic.config()|.
 (default: `false`)

## source


```lua
(boolean|'if_many')?
```


 Include the diagnostic source in the message.
 Use "if_many" to only show sources if there is more than one source of
 diagnostics in the buffer. Otherwise, any truthy value means to always show
 the diagnostic source.
 Overrides the setting from |vim.diagnostic.config()|.

## suffix


```lua
(string|table|fun(diagnostic: vim.Diagnostic, i: integer, total: integer):string, string)?
```


 Same as {prefix}, but appends the text to the diagnostic instead of
 prepending it.
 Overrides the setting from |vim.diagnostic.config()|.


---

# vim.diagnostic.Opts.Jump

## _highest


```lua
boolean?
```


 Default value of the {_highest} parameter of |vim.diagnostic.jump()|.

## float


```lua
(boolean|vim.diagnostic.Opts.Float)?
```


 Default value of the {float} parameter of |vim.diagnostic.jump()|.
 (default: false)

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 Default value of the {severity} parameter of |vim.diagnostic.jump()|.

## wrap


```lua
boolean?
```


 Default value of the {wrap} parameter of |vim.diagnostic.jump()|.
 (default: true)


---

# vim.diagnostic.Opts.Signs

## linehl


```lua
table<vim.diagnostic.Severity, string>?
```


 A table mapping |diagnostic-severity| to the highlight group used for the
 whole line the sign is placed in.

## numhl


```lua
table<vim.diagnostic.Severity, string>?
```


 A table mapping |diagnostic-severity| to the highlight group used for the
 line number where the sign is placed.

## priority


```lua
integer?
```


 Base priority to use for signs. When {severity_sort} is used, the priority
 of a sign is adjusted based on its severity.
 Otherwise, all signs use the same priority.
 (default: `10`)

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 Only show signs for diagnostics matching the given
 severity |diagnostic-severity|

## text


```lua
table<vim.diagnostic.Severity, string>?
```


 A table mapping |diagnostic-severity| to the sign text to display in the
 sign column. The default is to use `"E"`, `"W"`, `"I"`, and `"H"` for errors,
 warnings, information, and hints, respectively. Example:
 ```lua
 vim.diagnostic.config({
   signs = { text = { [vim.diagnostic.severity.ERROR] = 'E', ... } }
 })
 ```


---

# vim.diagnostic.Opts.Underline

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 Only underline diagnostics matching the given
 severity |diagnostic-severity|.


---

# vim.diagnostic.Opts.VirtualLines

## current_line


```lua
boolean?
```


 Only show diagnostics for the current line.
 (default: `false`)

## format


```lua
fun(diagnostic: vim.Diagnostic):string??
```


 A function that takes a diagnostic as input and returns a string or nil.
 If the return value is nil, the diagnostic is not displayed by the handler.
 Else the output text is used to display the diagnostic.

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 Only show virtual lines for diagnostics matching the given
 severity |diagnostic-severity|


---

# vim.diagnostic.Opts.VirtualText

## current_line


```lua
boolean?
```


 Only show diagnostics for the current line.
 (default `false`)

## format


```lua
fun(diagnostic: vim.Diagnostic):string??
```


 If not nil, the return value is the text used to display the diagnostic. Example:
 ```lua
 function(diagnostic)
   if diagnostic.severity == vim.diagnostic.severity.ERROR then
     return string.format("E: %s", diagnostic.message)
   end
   return diagnostic.message
 end
 ```
 If the return value is nil, the diagnostic is not displayed by the handler.

## hl_mode


```lua
('blend'|'combine'|'replace')?
```


 See |nvim_buf_set_extmark()|.

## prefix


```lua
(string|fun(diagnostic: vim.Diagnostic, i: integer, total: integer):string)?
```


 Prepend diagnostic message with prefix. If a `function`, {i} is the index
 of the diagnostic being evaluated, and {total} is the total number of
 diagnostics for the line. This can be used to render diagnostic symbols
 or error codes.

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 Only show virtual text for diagnostics matching the given
 severity |diagnostic-severity|

## source


```lua
(boolean|"if_many")?
```


 Include the diagnostic source in virtual text. Use `'if_many'` to only
 show sources if there is more than one diagnostic source in the buffer.
 Otherwise, any truthy value means to always show the diagnostic source.

## spacing


```lua
integer?
```


 Amount of empty spaces inserted at the beginning of the virtual text.

## suffix


```lua
(string|fun(diagnostic: vim.Diagnostic):string)?
```


 Append diagnostic message with suffix.
 This can be used to render an LSP diagnostic error code.

## virt_text


```lua
[string, any][]?
```


 See |nvim_buf_set_extmark()|.

## virt_text_hide


```lua
boolean?
```


 See |nvim_buf_set_extmark()|.

## virt_text_pos


```lua
('eol'|'eol_right_align'|'inline'|'overlay'|'right_align')?
```


 See |nvim_buf_set_extmark()|.

## virt_text_win_col


```lua
integer?
```


 See |nvim_buf_set_extmark()|.


---

# vim.diagnostic.OptsResolved

## float


```lua
vim.diagnostic.Opts.Float
```

## severity_sort


```lua
{ reverse: boolean }
```

## signs


```lua
vim.diagnostic.Opts.Signs
```

## underline


```lua
vim.diagnostic.Opts.Underline
```

## update_in_insert


```lua
boolean
```

## virtual_lines


```lua
vim.diagnostic.Opts.VirtualLines
```

## virtual_text


```lua
vim.diagnostic.Opts.VirtualText
```


---

# vim.diagnostic.Severity

 @nodoc


---

# vim.diagnostic.SeverityFilter

 See |diagnostic-severity| and |vim.diagnostic.get()|


---

# vim.diagnostic.SeverityInt


---

# vim.diagnostic._extmark

## [1]


```lua
integer
```

id

## [2]


```lua
integer
```

start

## [3]


```lua
integer
```

end

## [4]


```lua
table
```

details


---

# vim.diagnostic.setloclist.Opts

Configuration table with the following keys:

## namespace


```lua
integer?
```

 @inlinedoc

 Only add diagnostics from the given namespace.

## open


```lua
boolean?
```


 Open the location list after setting.
 (default: `true`)

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 See |diagnostic-severity|.

## title


```lua
string?
```


 Title of the location list. Defaults to "Diagnostics".

## winnr


```lua
integer?
```


 Window number to set location list for.
 (default: `0`)


---

# vim.diagnostic.setqflist.Opts

 Configuration table with the following keys:

## namespace


```lua
integer?
```

 @inlinedoc

 Only add diagnostics from the given namespace.

## open


```lua
boolean?
```


 Open quickfix list after setting.
 (default: `true`)

## severity


```lua
(vim.diagnostic.Severity|vim.diagnostic.Severity[]|{ min: vim.diagnostic.Severity, max: vim.diagnostic.Severity })?
```


 See |diagnostic-severity|.

## title


```lua
string?
```


 Title of quickfix list. Defaults to "Diagnostics". If there's already a quickfix list with this
 title, it's updated. If not, a new quickfix list is created.


---

# vim.diff


```lua
function vim.diff(a: string, b: string, opts?: vim.diff.Opts)
  -> (string|integer[][])?
```


---

# vim.diff.Opts

 Optional parameters:

## algorithm


```lua
('histogram'|'minimal'|'myers'|'patience')?
```


 Diff algorithm to use. Values:
   - `myers`: the default algorithm
   - `minimal`: spend extra time to generate the smallest possible diff
   - `patience`: patience diff algorithm
   - `histogram`: histogram diff algorithm
 (default: `'myers'`)

## ctxlen


```lua
integer?
```

Context length

## ignore_blank_lines


```lua
boolean?
```

Ignore blank lines

## ignore_cr_at_eol


```lua
boolean?
```

Ignore carriage return at end-of-line

## ignore_whitespace


```lua
boolean?
```

Ignore whitespace

## ignore_whitespace_change


```lua
boolean?
```

Ignore whitespace change

## ignore_whitespace_change_at_eol


```lua
boolean?
```

Ignore whitespace change at end-of-line.

## indent_heuristic


```lua
boolean?
```

Use the indent heuristic for the internal diff library.

## interhunkctxlen


```lua
integer?
```

Inter hunk context length

## linematch


```lua
(boolean|integer)?
```


 Run linematch on the resulting hunks from xdiff. When integer, only hunks
 upto this size in lines are run through linematch.
 Requires `result_type = indices`, ignored otherwise.

## on_hunk


```lua
fun(start_a: integer, count_a: integer, start_b: integer, count_b: integer):integer??
```

 @inlinedoc

 Invoked for each hunk in the diff. Return a negative number
 to cancel the callback for any remaining hunks.
 Arguments:
   - `start_a` (`integer`): Start line of hunk in {a}.
   - `count_a` (`integer`): Hunk size in {a}.
   - `start_b` (`integer`): Start line of hunk in {b}.
   - `count_b` (`integer`): Hunk size in {b}.

## result_type


```lua
('indices'|'unified')?
```


 Form of the returned diff:
   - `unified`: String in unified format.
   - `indices`: Array of hunk locations.
 Note: This option is ignored if `on_hunk` is used.
 (default: `'unified'`)


---

# vim.empty_dict


```lua
function vim.empty_dict()
  -> table
```


```lua
function vim.empty_dict()
  -> table
```


---

# vim.endswith


```lua
function vim.endswith(s: string, suffix: string)
  -> boolean
```


---

# vim.env


```lua
table
```


---

# vim.env.MANSECT


```lua
nil
```


```lua
string
```


---

# vim.filetype


```lua
table
```


---

# vim.filetype.add.filetypes

## extension


```lua
table<string, string|[string|fun(path: string, bufnr: integer, ...any):string?, fun(b: integer)?, { priority: number }]|fun(path: string, bufnr: integer, ...any):string?, fun(b: integer)?>?
```

## filename


```lua
table<string, string|[string|fun(path: string, bufnr: integer, ...any):string?, fun(b: integer)?, { priority: number }]|fun(path: string, bufnr: integer, ...any):string?, fun(b: integer)?>?
```

## pattern


```lua
table<string, string|[string|fun(path: string, bufnr: integer, ...any):string?, fun(b: integer)?, { priority: number }]|fun(path: string, bufnr: integer, ...any):string?, fun(b: integer)?>?
```

 @inlinedoc


---

# vim.filetype.mapfn


---

# vim.filetype.mapopts


---

# vim.filetype.mapping


---

# vim.filetype.mapping.sorted

## [1]


```lua
string
```

parent pattern

## [2]


```lua
string
```

pattern

## [3]


```lua
string|fun(path: string, bufnr: integer, ...any):string?, fun(b: integer)?
```

## [4]


```lua
integer
```

priority


---

# vim.filetype.mapping.value


---

# vim.filetype.maptbl


---

# vim.filetype.match.args

## buf


```lua
integer?
```

 @inlinedoc

 Buffer number to use for matching. Mutually exclusive with {contents}

## contents


```lua
string[]?
```


 An array of lines representing file contents to use for
 matching. Can be used with {filename}. Mutually exclusive
 with {buf}.

## filename


```lua
string?
```


 Filename to use for matching. When {buf} is given,
 defaults to the filename of the given buffer number. The
 file need not actually exist in the filesystem. When used
 without {buf} only the name of the file is used for
 filetype matching. This may result in failure to detect
 the filetype in cases where the filename alone is not
 enough to disambiguate the filetype.


---

# vim.filetype.pattern_cache

 Lookup table/cache for patterns


---

# vim.fn


```lua
table
```


---

# vim.fn.abs


```lua
function table.abs(expr: number)
  -> number
```


---

# vim.fn.acos


```lua
function table.acos(expr: number)
  -> number
```


---

# vim.fn.add


```lua
function table.add(object: any, expr: any)
  -> any
```


---

# vim.fn.and


```lua
function (expr: number, expr1: number)
  -> integer
```


---

# vim.fn.api_info


```lua
function table.api_info()
  -> table
```


---

# vim.fn.append


```lua
function table.append(lnum: string|integer, text: string|string[])
  -> 0|1
```


---

# vim.fn.appendbufline


```lua
function table.appendbufline(buf: string|integer, lnum: integer, text: string)
  -> 0|1
```


---

# vim.fn.argc


```lua
function table.argc(winid?: integer)
  -> integer
```


---

# vim.fn.argidx


```lua
function table.argidx()
  -> integer
```


---

# vim.fn.arglistid


```lua
function table.arglistid(winnr?: integer, tabnr?: integer)
  -> integer
```


---

# vim.fn.argv


```lua
function table.argv(nr?: integer, winid?: integer)
  -> string|string[]
```


---

# vim.fn.asin


```lua
function table.asin(expr: any)
  -> number
```


---

# vim.fn.assert_beeps


```lua
function table.assert_beeps(cmd: string)
  -> 0|1
```


---

# vim.fn.assert_equal


```lua
function table.assert_equal(expected: any, actual: any, msg?: any)
  -> 0|1
```


---

# vim.fn.assert_equalfile


```lua
function table.assert_equalfile(fname_one: string, fname_two: string)
  -> 0|1
```


---

# vim.fn.assert_exception


```lua
function table.assert_exception(error: any, msg?: any)
  -> 0|1
```


---

# vim.fn.assert_fails


```lua
function table.assert_fails(cmd: string, error?: any, msg?: any, lnum?: integer, context?: any)
  -> 0|1
```


---

# vim.fn.assert_false


```lua
function table.assert_false(actual: any, msg?: any)
  -> 0|1
```


---

# vim.fn.assert_inrange


```lua
function table.assert_inrange(lower: number, upper: number, actual: number, msg?: string)
  -> 0|1
```


---

# vim.fn.assert_match


```lua
function table.assert_match(pattern: string, actual: string, msg?: string)
  -> 0|1
```


---

# vim.fn.assert_nobeep


```lua
function table.assert_nobeep(cmd: string)
  -> 0|1
```


---

# vim.fn.assert_notequal


```lua
function table.assert_notequal(expected: any, actual: any, msg?: any)
  -> 0|1
```


---

# vim.fn.assert_notmatch


```lua
function table.assert_notmatch(pattern: string, actual: string, msg?: string)
  -> 0|1
```


---

# vim.fn.assert_report


```lua
function table.assert_report(msg: string)
  -> 0|1
```


---

# vim.fn.assert_true


```lua
function table.assert_true(actual: any, msg?: string)
  -> 0|1
```


---

# vim.fn.atan


```lua
function table.atan(expr: number)
  -> number
```


---

# vim.fn.atan2


```lua
function table.atan2(expr1: number, expr2: number)
  -> number
```


---

# vim.fn.blob2list


```lua
function table.blob2list(blob: any)
  -> any[]
```


---

# vim.fn.browse


```lua
function table.browse(save: any, title: string, initdir: string, default: string)
  -> 0|1
```


---

# vim.fn.browsedir


```lua
function table.browsedir(title: string, initdir: string)
  -> 0|1
```


---

# vim.fn.bufadd


```lua
function table.bufadd(name: string)
  -> integer
```


---

# vim.fn.bufexists


```lua
function table.bufexists(buf: any)
  -> 0|1
```


---

# vim.fn.buffer_exists


```lua
function table.buffer_exists(...any)
  -> 0|1
```


---

# vim.fn.buffer_name


```lua
function table.buffer_name(...any)
  -> string
```


---

# vim.fn.buffer_number


```lua
function table.buffer_number(...any)
  -> integer
```


---

# vim.fn.buflisted


```lua
function table.buflisted(buf: any)
  -> 0|1
```


---

# vim.fn.bufload


```lua
function table.bufload(buf: any)
```


---

# vim.fn.bufloaded


```lua
function table.bufloaded(buf: any)
  -> 0|1
```


---

# vim.fn.bufname


```lua
function table.bufname(buf?: string|integer)
  -> string
```


---

# vim.fn.bufnr


```lua
function table.bufnr(buf?: string|integer, create?: any)
  -> integer
```


---

# vim.fn.bufwinid


```lua
function table.bufwinid(buf: any)
  -> integer
```


---

# vim.fn.bufwinnr


```lua
function table.bufwinnr(buf: any)
  -> integer
```


---

# vim.fn.byte2line


```lua
function table.byte2line(byte: any)
  -> integer
```


---

# vim.fn.byteidx


```lua
function table.byteidx(expr: any, nr: integer, utf16?: any)
  -> integer
```


---

# vim.fn.byteidxcomp


```lua
function table.byteidxcomp(expr: any, nr: integer, utf16?: any)
  -> integer
```


---

# vim.fn.call


```lua
function table.call(func: any, arglist: any, dict?: any)
  -> any
```


---

# vim.fn.ceil


```lua
function table.ceil(expr: number)
  -> number
```


---

# vim.fn.chanclose


```lua
function table.chanclose(id: integer, stream?: string)
  -> 0|1
```


---

# vim.fn.changenr


```lua
function table.changenr()
  -> integer
```


---

# vim.fn.chansend


```lua
function table.chansend(id: number, data: string|string[])
  -> 0|1
```


---

# vim.fn.char2nr


```lua
function table.char2nr(string: string, utf8?: any)
  -> 0|1
```


---

# vim.fn.charclass


```lua
function table.charclass(string: string)
  -> 'other'|0|1|2|3
```


---

# vim.fn.charcol


```lua
function table.charcol(expr: string|any[], winid?: integer)
  -> integer
```


---

# vim.fn.charidx


```lua
function table.charidx(string: string, idx: integer, countcc?: boolean, utf16?: boolean)
  -> integer
```


---

# vim.fn.chdir


```lua
function table.chdir(dir: string)
  -> string
```


---

# vim.fn.cindent


```lua
function table.cindent(lnum: string|integer)
  -> integer
```


---

# vim.fn.clearmatches


```lua
function table.clearmatches(win?: integer)
```


---

# vim.fn.col


```lua
function table.col(expr: string|any[], winid?: integer)
  -> integer
```


---

# vim.fn.complete


```lua
function table.complete(startcol: integer, matches: any[])
```


---

# vim.fn.complete_add


```lua
function table.complete_add(expr: any)
  -> 0|1|2
```


---

# vim.fn.complete_check


```lua
function table.complete_check()
  -> 0|1
```


---

# vim.fn.complete_info


```lua
function table.complete_info(what?: any[])
  -> table
```


---

# vim.fn.confirm


```lua
function table.confirm(msg: string, choices?: string, default?: integer, type?: string)
  -> integer
```


---

# vim.fn.copy


```lua
function table.copy(expr: <T>)
  -> <T>
```


---

# vim.fn.cos


```lua
function table.cos(expr: number)
  -> number
```


---

# vim.fn.cosh


```lua
function table.cosh(expr: number)
  -> number
```


---

# vim.fn.count


```lua
function table.count(comp: string|table|any[], expr: any, ic?: boolean, start?: integer)
  -> integer
```


---

# vim.fn.ctxget


```lua
function table.ctxget(index?: integer)
  -> table
```


---

# vim.fn.ctxpop


```lua
function table.ctxpop()
  -> any
```


---

# vim.fn.ctxpush


```lua
function table.ctxpush(types?: string[])
  -> any
```


---

# vim.fn.ctxset


```lua
function table.ctxset(context: table, index?: integer)
  -> integer
```


---

# vim.fn.ctxsize


```lua
function table.ctxsize()
  -> any
```


---

# vim.fn.cursor


```lua
function table.cursor(lnum: string|integer, col?: integer, off?: integer)
  -> any
```


```lua
function table.cursor(list: integer[])
  -> any
```


---

# vim.fn.debugbreak


```lua
function table.debugbreak(pid: integer)
  -> any
```


---

# vim.fn.deepcopy


```lua
function table.deepcopy(expr: <T>, noref?: boolean)
  -> <T>
```


---

# vim.fn.delete


```lua
function table.delete(fname: string, flags?: string)
  -> integer
```


---

# vim.fn.deletebufline


```lua
function table.deletebufline(buf: string|integer, first: string|integer, last?: string|integer)
  -> any
```


---

# vim.fn.dictwatcheradd


```lua
function table.dictwatcheradd(dict: table, pattern: string, callback: function)
  -> any
```


---

# vim.fn.dictwatcherdel


```lua
function table.dictwatcherdel(dict: any, pattern: string, callback: function)
  -> any
```


---

# vim.fn.did_filetype


```lua
function table.did_filetype()
  -> integer
```


---

# vim.fn.diff_filler


```lua
function table.diff_filler(lnum: string|integer)
  -> integer
```


---

# vim.fn.diff_hlID


```lua
function table.diff_hlID(lnum: string|integer, col: integer)
  -> any
```


---

# vim.fn.digraph_get


```lua
function table.digraph_get(chars: string)
  -> string
```


---

# vim.fn.digraph_getlist


```lua
function table.digraph_getlist(listall?: boolean)
  -> string[][]
```


---

# vim.fn.digraph_set


```lua
function table.digraph_set(chars: string, digraph: string)
  -> any
```


---

# vim.fn.digraph_setlist


```lua
function table.digraph_setlist(digraphlist: table<integer, string[]>)
  -> any
```


---

# vim.fn.empty


```lua
function table.empty(expr: any)
  -> integer
```


---

# vim.fn.environ


```lua
function table.environ()
  -> any
```


---

# vim.fn.escape


```lua
function table.escape(string: string, chars: string)
  -> string
```


---

# vim.fn.eval


```lua
function table.eval(string: string)
  -> any
```


---

# vim.fn.eventhandler


```lua
function table.eventhandler()
  -> any
```


---

# vim.fn.executable


```lua
function table.executable(expr: string)
  -> 0|1
```


---

# vim.fn.execute


```lua
function table.execute(command: string|string[], silent?: ''|'silent!'|'silent')
  -> string
```


---

# vim.fn.exepath


```lua
function table.exepath(expr: string)
  -> string
```


---

# vim.fn.exists


```lua
function table.exists(expr: string)
  -> 0|1
```


---

# vim.fn.exp


```lua
function table.exp(expr: number)
  -> any
```


---

# vim.fn.expand


```lua
function table.expand(string: string, nosuf?: boolean, list?: false)
  -> string
```


```lua
function table.expand(string: string, nosuf: boolean, list: string|number|table|true)
  -> string|string[]
```


---

# vim.fn.expandcmd


```lua
function table.expandcmd(string: string, options?: table)
  -> any
```


---

# vim.fn.extend


```lua
function table.extend(expr1: table, expr2: table, expr3?: table)
  -> any
```


---

# vim.fn.extendnew


```lua
function table.extendnew(expr1: table, expr2: table, expr3?: table)
  -> any
```


---

# vim.fn.feedkeys


```lua
function table.feedkeys(string: string, mode?: string)
  -> any
```


---

# vim.fn.file_readable


```lua
function table.file_readable(file: string)
  -> any
```


---

# vim.fn.filecopy


```lua
function table.filecopy(from: string, to: string)
  -> 0|1
```


---

# vim.fn.filereadable


```lua
function table.filereadable(file: string)
  -> 0|1
```


---

# vim.fn.filewritable


```lua
function table.filewritable(file: string)
  -> 0|1
```


---

# vim.fn.filter


```lua
function table.filter(expr1: string|table, expr2: string|function)
  -> any
```


---

# vim.fn.finddir


```lua
function table.finddir(name: string, path?: string, count?: integer)
  -> any
```


---

# vim.fn.findfile


```lua
function table.findfile(name: string, path?: string, count?: any)
  -> any
```


---

# vim.fn.flatten


```lua
function table.flatten(list: any[], maxdepth?: integer)
  -> 0|any[]
```


---

# vim.fn.flattennew


```lua
function table.flattennew(list: any[], maxdepth?: integer)
  -> 0|any[]
```


---

# vim.fn.float2nr


```lua
function table.float2nr(expr: number)
  -> any
```


---

# vim.fn.floor


```lua
function table.floor(expr: number)
  -> any
```


---

# vim.fn.fmod


```lua
function table.fmod(expr1: number, expr2: number)
  -> any
```


---

# vim.fn.fnameescape


```lua
function table.fnameescape(string: string)
  -> string
```


---

# vim.fn.fnamemodify


```lua
function table.fnamemodify(fname: string, mods: string)
  -> string
```


---

# vim.fn.foldclosed


```lua
function table.foldclosed(lnum: string|integer)
  -> integer
```


---

# vim.fn.foldclosedend


```lua
function table.foldclosedend(lnum: string|integer)
  -> integer
```


---

# vim.fn.foldlevel


```lua
function table.foldlevel(lnum: string|integer)
  -> integer
```


---

# vim.fn.foldtext


```lua
function table.foldtext()
  -> string
```


---

# vim.fn.foldtextresult


```lua
function table.foldtextresult(lnum: string|integer)
  -> string
```


---

# vim.fn.foreach


```lua
function table.foreach(expr1: string|table, expr2: string|function)
  -> string|table
```


---

# vim.fn.fullcommand


```lua
function table.fullcommand(name: string)
  -> string
```


---

# vim.fn.funcref


```lua
function table.funcref(name: string, arglist?: any, dict?: any)
  -> any
```


---

# vim.fn.function


```lua
function (name: string, arglist?: any, dict?: any)
  -> any
```


---

# vim.fn.garbagecollect


```lua
function table.garbagecollect(atexit?: boolean)
  -> any
```


---

# vim.fn.get


```lua
function table.get(list: any[], idx: integer, default?: any)
  -> any
```


```lua
function table.get(blob: string, idx: integer, default?: any)
  -> any
```


```lua
function table.get(dict: table<string, any>, key: string, default?: any)
  -> any
```


```lua
function table.get(func: function, what: string)
  -> any
```


---

# vim.fn.getbufinfo


```lua
function table.getbufinfo(buf?: string|integer)
  -> vim.fn.getbufinfo.ret.item[]
```


```lua
function table.getbufinfo(dict?: vim.fn.getbufinfo.dict)
  -> vim.fn.getbufinfo.ret.item[]
```


---

# vim.fn.getbufinfo.dict

## buflisted


```lua
(0|1)?
```

## bufloaded


```lua
(0|1)?
```

## bufmodified


```lua
(0|1)?
```


---

# vim.fn.getbufinfo.ret.item

## bufnr


```lua
integer
```

## changed


```lua
0|1
```

## changedtick


```lua
integer
```

## hidden


```lua
0|1
```

## lastused


```lua
integer
```

## linecount


```lua
integer
```

## listed


```lua
0|1
```

## lnum


```lua
integer
```

## loaded


```lua
0|1
```

## name


```lua
string
```

## signs


```lua
vim.fn.sign[]
```

## variables


```lua
table<string, any>
```

## windows


```lua
integer[]
```


---

# vim.fn.getbufline


```lua
function table.getbufline(buf: string|integer, lnum: integer, end_?: integer)
  -> string[]
```


---

# vim.fn.getbufoneline


```lua
function table.getbufoneline(buf: string|integer, lnum: integer)
  -> string
```


---

# vim.fn.getbufvar


```lua
function table.getbufvar(buf: string|integer, varname: string, def?: any)
  -> any
```


---

# vim.fn.getcellwidths


```lua
function table.getcellwidths()
  -> any
```


---

# vim.fn.getchangelist


```lua
function table.getchangelist(buf?: string|integer)
  -> table[]
```


---

# vim.fn.getchar


```lua
function table.getchar(expr?: -1|0|1, opts?: table)
  -> string|integer
```


---

# vim.fn.getcharmod


```lua
function table.getcharmod()
  -> integer
```


---

# vim.fn.getcharpos


```lua
function table.getcharpos(expr: string)
  -> integer[]
```


---

# vim.fn.getcharsearch


```lua
function table.getcharsearch()
  -> table
```


---

# vim.fn.getcharstr


```lua
function table.getcharstr(expr?: -1|0|1, opts?: table)
  -> string
```


---

# vim.fn.getcmdcomplpat


```lua
function table.getcmdcomplpat()
  -> string
```


---

# vim.fn.getcmdcompltype


```lua
function table.getcmdcompltype()
  -> string
```


---

# vim.fn.getcmdline


```lua
function table.getcmdline()
  -> string
```


---

# vim.fn.getcmdpos


```lua
function table.getcmdpos()
  -> integer
```


---

# vim.fn.getcmdprompt


```lua
function table.getcmdprompt()
  -> string
```


---

# vim.fn.getcmdscreenpos


```lua
function table.getcmdscreenpos()
  -> integer
```


---

# vim.fn.getcmdtype


```lua
function table.getcmdtype()
  -> '-'|'/'|':'|'='|'>'...(+2)
```


---

# vim.fn.getcmdwintype


```lua
function table.getcmdwintype()
  -> '-'|'/'|':'|'='|'>'...(+2)
```


---

# vim.fn.getcompletion


```lua
function table.getcompletion(pat: string, type: string, filtered?: boolean)
  -> string[]
```


---

# vim.fn.getcurpos


```lua
function table.getcurpos(winid?: integer)
  -> any
```


---

# vim.fn.getcursorcharpos


```lua
function table.getcursorcharpos(winid?: integer)
  -> any
```


---

# vim.fn.getcwd


```lua
function table.getcwd(winnr?: integer, tabnr?: integer)
  -> string
```


---

# vim.fn.getenv


```lua
function table.getenv(name: string)
  -> string
```


---

# vim.fn.getfontname


```lua
function table.getfontname(name?: string)
  -> string
```


---

# vim.fn.getfperm


```lua
function table.getfperm(fname: string)
  -> string
```


---

# vim.fn.getfsize


```lua
function table.getfsize(fname: string)
  -> integer
```


---

# vim.fn.getftime


```lua
function table.getftime(fname: string)
  -> integer
```


---

# vim.fn.getftype


```lua
function table.getftype(fname: string)
  -> 'bdev'|'cdev'|'dir'|'fifo'|'file'...(+3)
```


---

# vim.fn.getjumplist


```lua
function table.getjumplist(winnr?: integer, tabnr?: integer)
  -> [vim.fn.getjumplist.ret.item[], integer]
```


---

# vim.fn.getjumplist.ret


---

# vim.fn.getjumplist.ret.item

## bufnr


```lua
integer
```

## col


```lua
integer
```

## coladd


```lua
integer
```

## filename


```lua
string?
```

## lnum


```lua
integer
```


---

# vim.fn.getline


```lua
function table.getline(lnum: string|integer, end_?: false)
  -> string
```


```lua
function table.getline(lnum: string|integer, end_: string|number|table|true)
  -> string|string[]
```


---

# vim.fn.getloclist


```lua
function table.getloclist(nr: integer, what?: table)
  -> any
```


---

# vim.fn.getmarklist


```lua
function table.getmarklist(buf?: integer)
  -> vim.fn.getmarklist.ret.item[]
```


---

# vim.fn.getmarklist.ret.item

## file


```lua
string
```

## mark


```lua
string
```

## pos


```lua
[integer, integer, integer, integer]
```


---

# vim.fn.getmatches


```lua
function table.getmatches(win?: integer)
  -> any
```


---

# vim.fn.getmousepos


```lua
function table.getmousepos()
  -> vim.fn.getmousepos.ret
```


---

# vim.fn.getmousepos.ret

## column


```lua
integer
```

## line


```lua
integer
```

## screencol


```lua
integer
```

## screenrow


```lua
integer
```

## wincol


```lua
integer
```

## winid


```lua
integer
```

## winrow


```lua
integer
```


---

# vim.fn.getpid


```lua
function table.getpid()
  -> integer
```


---

# vim.fn.getpos


```lua
function table.getpos(expr: string)
  -> integer[]
```


---

# vim.fn.getqflist


```lua
function table.getqflist(what?: table)
  -> any
```


---

# vim.fn.getreg


```lua
function table.getreg(regname?: string, list?: false)
  -> string
```


```lua
function table.getreg(regname: string, list: string|number|table|true)
  -> string|string[]
```


---

# vim.fn.getreginfo


```lua
function table.getreginfo(regname?: string)
  -> table
```


---

# vim.fn.getregion


```lua
function table.getregion(pos1: table, pos2: table, opts?: table)
  -> string[]
```


---

# vim.fn.getregionpos


```lua
function table.getregionpos(pos1: table, pos2: table, opts?: table)
  -> integer[][][]
```


---

# vim.fn.getregtype


```lua
function table.getregtype(regname?: string)
  -> string
```


---

# vim.fn.getscriptinfo


```lua
function table.getscriptinfo(opts?: table)
  -> vim.fn.getscriptinfo.ret[]
```


---

# vim.fn.getscriptinfo.ret

## autoload


```lua
false
```

## functions


```lua
string[]?
```

## name


```lua
string
```

## sid


```lua
string
```

## variables


```lua
table<string, any>?
```

## version


```lua
1
```


---

# vim.fn.getstacktrace


```lua
function table.getstacktrace()
  -> table[]
```


---

# vim.fn.gettabinfo


```lua
function table.gettabinfo(tabnr?: integer)
  -> any
```


---

# vim.fn.gettabvar


```lua
function table.gettabvar(tabnr: integer, varname: string, def?: any)
  -> any
```


---

# vim.fn.gettabwinvar


```lua
function table.gettabwinvar(tabnr: integer, winnr: integer, varname: string, def?: any)
  -> any
```


---

# vim.fn.gettagstack


```lua
function table.gettagstack(winnr?: integer)
  -> any
```


---

# vim.fn.gettext


```lua
function table.gettext(text: string)
  -> string
```


---

# vim.fn.getwininfo


```lua
function table.getwininfo(winid?: integer)
  -> vim.fn.getwininfo.ret.item[]
```


---

# vim.fn.getwininfo.ret.item

## botline


```lua
integer
```

## bufnr


```lua
integer
```

## height


```lua
integer
```

## loclist


```lua
integer
```

## quickfix


```lua
integer
```

## tabnr


```lua
integer
```

## terminal


```lua
integer
```

## textoff


```lua
integer
```

## topline


```lua
integer
```

## variables


```lua
table<string, any>
```

## width


```lua
integer
```

## winbar


```lua
integer
```

## wincol


```lua
integer
```

## winid


```lua
integer
```

## winnr


```lua
integer
```

## winrow


```lua
integer
```


---

# vim.fn.getwinpos


```lua
function table.getwinpos(timeout?: integer)
  -> any
```


---

# vim.fn.getwinposx


```lua
function table.getwinposx()
  -> integer
```


---

# vim.fn.getwinposy


```lua
function table.getwinposy()
  -> integer
```


---

# vim.fn.getwinvar


```lua
function table.getwinvar(winnr: integer, varname: string, def?: any)
  -> any
```


---

# vim.fn.glob


```lua
function table.glob(expr: string, nosuf?: boolean, list?: boolean, alllinks?: boolean)
  -> any
```


---

# vim.fn.glob2regpat


```lua
function table.glob2regpat(string: string)
  -> string
```


---

# vim.fn.globpath


```lua
function table.globpath(path: string, expr: string, nosuf?: boolean, list?: boolean, allinks?: boolean)
  -> any
```


---

# vim.fn.has


```lua
function table.has(feature: string)
  -> 0|1
```


---

# vim.fn.has_key


```lua
function table.has_key(dict: table, key: string)
  -> 0|1
```


---

# vim.fn.haslocaldir


```lua
function table.haslocaldir(winnr?: integer, tabnr?: integer)
  -> 0|1
```


---

# vim.fn.hasmapto


```lua
function table.hasmapto(what: any, mode?: string, abbr?: boolean)
  -> 0|1
```


---

# vim.fn.highlightID


```lua
function table.highlightID(name: string)
  -> any
```


---

# vim.fn.highlight_exists


```lua
function table.highlight_exists(name: string)
  -> any
```


---

# vim.fn.histadd


```lua
function table.histadd(history: string, item: any)
  -> 0|1
```


---

# vim.fn.histdel


```lua
function table.histdel(history: string, item?: any)
  -> 0|1
```


---

# vim.fn.histget


```lua
function table.histget(history: string, index?: string|integer)
  -> string
```


---

# vim.fn.histnr


```lua
function table.histnr(history: string)
  -> integer
```


---

# vim.fn.hlID


```lua
function table.hlID(name: string)
  -> integer
```


---

# vim.fn.hlexists


```lua
function table.hlexists(name: string)
  -> 0|1
```


---

# vim.fn.hostname


```lua
function table.hostname()
  -> string
```


---

# vim.fn.iconv


```lua
function table.iconv(string: string, from: string, to: string)
  -> string
```


---

# vim.fn.id


```lua
function table.id(expr: any)
  -> string
```


---

# vim.fn.indent


```lua
function table.indent(lnum: string|integer)
  -> integer
```


---

# vim.fn.index


```lua
function table.index(object: any, expr: any, start?: integer, ic?: boolean)
  -> integer
```


---

# vim.fn.indexof


```lua
function table.indexof(object: any, expr: any, opts?: table)
  -> integer
```


---

# vim.fn.input


```lua
function table.input(prompt: string, text?: string, completion?: string)
  -> string
```


```lua
function table.input(opts: table)
  -> string
```


---

# vim.fn.inputdialog


```lua
function table.inputdialog(...any)
  -> any
```


---

# vim.fn.inputlist


```lua
function table.inputlist(textlist: string[])
  -> any
```


---

# vim.fn.inputrestore


```lua
function table.inputrestore()
  -> integer
```


---

# vim.fn.inputsave


```lua
function table.inputsave()
  -> integer
```


---

# vim.fn.inputsecret


```lua
function table.inputsecret(prompt: string, text?: string)
  -> string
```


---

# vim.fn.insert


```lua
function table.insert(object: any, item: any, idx?: integer)
  -> any
```


---

# vim.fn.interrupt


```lua
function table.interrupt()
  -> any
```


---

# vim.fn.invert


```lua
function table.invert(expr: integer)
  -> integer
```


---

# vim.fn.isabsolutepath


```lua
function table.isabsolutepath(path: string)
  -> 0|1
```


---

# vim.fn.isdirectory


```lua
function table.isdirectory(directory: string)
  -> 0|1
```


---

# vim.fn.isinf


```lua
function table.isinf(expr: number)
  -> -1|0|1
```


---

# vim.fn.islocked


```lua
function table.islocked(expr: any)
  -> 0|1
```


---

# vim.fn.isnan


```lua
function table.isnan(expr: number)
  -> 0|1
```


---

# vim.fn.items


```lua
function table.items(dict: table)
  -> any
```


---

# vim.fn.jobclose


```lua
function table.jobclose(...any)
  -> any
```


---

# vim.fn.jobpid


```lua
function table.jobpid(job: integer)
  -> integer
```


---

# vim.fn.jobresize


```lua
function table.jobresize(job: integer, width: integer, height: integer)
  -> any
```


---

# vim.fn.jobsend


```lua
function table.jobsend(...any)
  -> any
```


---

# vim.fn.jobstart


```lua
function table.jobstart(cmd: string|string[], opts?: table)
  -> integer
```


---

# vim.fn.jobstop


```lua
function table.jobstop(id: integer)
  -> integer
```


---

# vim.fn.jobwait


```lua
function table.jobwait(jobs: integer[], timeout?: integer)
  -> integer[]
```


---

# vim.fn.join


```lua
function table.join(list: any[], sep?: string)
  -> string
```


---

# vim.fn.json_decode


```lua
function table.json_decode(expr: any)
  -> any
```


---

# vim.fn.json_encode


```lua
function table.json_encode(expr: any)
  -> string
```


---

# vim.fn.keys


```lua
function table.keys(dict: table)
  -> string[]
```


---

# vim.fn.keytrans


```lua
function table.keytrans(string: string)
  -> string
```


---

# vim.fn.last_buffer_nr


```lua
function table.last_buffer_nr()
  -> any
```


---

# vim.fn.len


```lua
function table.len(expr: any[])
  -> integer
```


---

# vim.fn.libcall


```lua
function table.libcall(libname: string, funcname: string, argument: any)
  -> any
```


---

# vim.fn.libcallnr


```lua
function table.libcallnr(libname: string, funcname: string, argument: any)
  -> any
```


---

# vim.fn.line


```lua
function table.line(expr: string|integer[], winid?: integer)
  -> integer
```


---

# vim.fn.line2byte


```lua
function table.line2byte(lnum: string|integer)
  -> integer
```


---

# vim.fn.lispindent


```lua
function table.lispindent(lnum: string|integer)
  -> integer
```


---

# vim.fn.list2blob


```lua
function table.list2blob(list: any[])
  -> string
```


---

# vim.fn.list2str


```lua
function table.list2str(list: any[], utf8?: boolean)
  -> string
```


---

# vim.fn.localtime


```lua
function table.localtime()
  -> integer
```


---

# vim.fn.log


```lua
function table.log(expr: number)
  -> number
```


---

# vim.fn.log10


```lua
function table.log10(expr: number)
  -> number
```


---

# vim.fn.map


```lua
function table.map(expr1: string|table|any[], expr2: string|function)
  -> any
```


---

# vim.fn.maparg


```lua
function table.maparg(name: string, mode?: string, abbr?: boolean, dict?: false)
  -> string
```


```lua
function table.maparg(name: string, mode: string, abbr: boolean, dict: true)
  -> table<string, any>
```


---

# vim.fn.mapcheck


```lua
function table.mapcheck(name: string, mode?: string, abbr?: boolean)
  -> any
```


---

# vim.fn.maplist


```lua
function table.maplist(abbr?: 0|1)
  -> table[]
```


---

# vim.fn.mapnew


```lua
function table.mapnew(expr1: any, expr2: any)
  -> any
```


---

# vim.fn.mapset


```lua
function table.mapset(mode: string, abbr?: boolean, dict?: table<string, any>)
  -> any
```


```lua
function table.mapset(dict: table<string, any>)
  -> any
```


---

# vim.fn.match


```lua
function table.match(expr: string|any[], pat: string, start?: integer, count?: integer)
  -> any
```


---

# vim.fn.matchadd


```lua
function table.matchadd(group: string|integer, pattern: string, priority?: integer, id?: integer, dict?: string)
  -> any
```


---

# vim.fn.matchaddpos


```lua
function table.matchaddpos(group: string|integer, pos: any[], priority?: integer, id?: integer, dict?: string)
  -> any
```


---

# vim.fn.matcharg


```lua
function table.matcharg(nr: integer)
  -> any
```


---

# vim.fn.matchbufline


```lua
function table.matchbufline(buf: string|integer, pat: string, lnum: string|integer, end_: string|integer, dict?: table)
  -> any
```


---

# vim.fn.matchdelete


```lua
function table.matchdelete(id: integer, win?: integer)
  -> any
```


---

# vim.fn.matchend


```lua
function table.matchend(expr: any, pat: string, start?: integer, count?: integer)
  -> any
```


---

# vim.fn.matchfuzzy


```lua
function table.matchfuzzy(list: any[], str: string, dict?: table)
  -> any
```


---

# vim.fn.matchfuzzypos


```lua
function table.matchfuzzypos(list: any[], str: string, dict?: table)
  -> any
```


---

# vim.fn.matchlist


```lua
function table.matchlist(expr: any, pat: string, start?: integer, count?: integer)
  -> any
```


---

# vim.fn.matchstr


```lua
function table.matchstr(expr: any, pat: string, start?: integer, count?: integer)
  -> any
```


---

# vim.fn.matchstrlist


```lua
function table.matchstrlist(list: string[], pat: string, dict?: table)
  -> any
```


---

# vim.fn.matchstrpos


```lua
function table.matchstrpos(expr: any, pat: string, start?: integer, count?: integer)
  -> any
```


---

# vim.fn.max


```lua
function table.max(expr: any)
  -> number
```


---

# vim.fn.menu_get


```lua
function table.menu_get(path: string, modes?: string)
  -> any
```


---

# vim.fn.menu_info


```lua
function table.menu_info(name: string, mode?: string)
  -> any
```


---

# vim.fn.min


```lua
function table.min(expr: any)
  -> number
```


---

# vim.fn.mkdir


```lua
function table.mkdir(name: string, flags?: string, prot?: string)
  -> integer
```


---

# vim.fn.mode


```lua
function table.mode(expr?: any)
  -> any
```


---

# vim.fn.msgpackdump


```lua
function table.msgpackdump(list: any, type?: any)
  -> any
```


---

# vim.fn.msgpackparse


```lua
function table.msgpackparse(data: any)
  -> any
```


---

# vim.fn.nextnonblank


```lua
function table.nextnonblank(lnum: string|integer)
  -> integer
```


---

# vim.fn.nr2char


```lua
function table.nr2char(expr: integer, utf8?: boolean)
  -> string
```


---

# vim.fn.or


```lua
function (expr: number, expr1: number)
  -> any
```


---

# vim.fn.pathshorten


```lua
function table.pathshorten(path: string, len?: integer)
  -> string
```


---

# vim.fn.perleval


```lua
function table.perleval(expr: any)
  -> any
```


---

# vim.fn.pow


```lua
function table.pow(x: number, y: number)
  -> number
```


---

# vim.fn.prevnonblank


```lua
function table.prevnonblank(lnum: string|integer)
  -> integer
```


---

# vim.fn.printf


```lua
function table.printf(fmt: string, expr1?: any)
  -> string
```


---

# vim.fn.prompt_getprompt


```lua
function table.prompt_getprompt(buf: string|integer)
  -> any
```


---

# vim.fn.prompt_setcallback


```lua
function table.prompt_setcallback(buf: string|integer, expr: string|function)
  -> any
```


---

# vim.fn.prompt_setinterrupt


```lua
function table.prompt_setinterrupt(buf: string|integer, expr: string|function)
  -> any
```


---

# vim.fn.prompt_setprompt


```lua
function table.prompt_setprompt(buf: string|integer, text: string)
  -> any
```


---

# vim.fn.pum_getpos


```lua
function table.pum_getpos()
  -> any
```


---

# vim.fn.pumvisible


```lua
function table.pumvisible()
  -> any
```


---

# vim.fn.py3eval


```lua
function table.py3eval(expr: any)
  -> any
```


---

# vim.fn.pyeval


```lua
function table.pyeval(expr: any)
  -> any
```


---

# vim.fn.pyxeval


```lua
function table.pyxeval(expr: any)
  -> any
```


---

# vim.fn.rand


```lua
function table.rand(expr?: number)
  -> any
```


---

# vim.fn.range


```lua
function table.range(expr: any, max?: integer, stride?: integer)
  -> any
```


---

# vim.fn.readblob


```lua
function table.readblob(fname: string, offset?: integer, size?: integer)
  -> any
```


---

# vim.fn.readdir


```lua
function table.readdir(directory: string, expr?: integer)
  -> any
```


---

# vim.fn.readfile


```lua
function table.readfile(fname: string, type?: string, max?: integer)
  -> any
```


---

# vim.fn.reduce


```lua
function table.reduce(object: any, func: fun(accumulator: <T>, current: any):any, initial?: any)
  -> <T>
```


---

# vim.fn.reg_executing


```lua
function table.reg_executing()
  -> any
```


---

# vim.fn.reg_recorded


```lua
function table.reg_recorded()
  -> any
```


---

# vim.fn.reg_recording


```lua
function table.reg_recording()
  -> any
```


---

# vim.fn.reltime


```lua
function table.reltime()
  -> any
```


```lua
function table.reltime(start?: any)
  -> any
```


```lua
function table.reltime(start?: any, end_?: any)
  -> any
```


---

# vim.fn.reltimefloat


```lua
function table.reltimefloat(time: any)
  -> any
```


---

# vim.fn.reltimestr


```lua
function table.reltimestr(time: any)
  -> any
```


---

# vim.fn.remove


```lua
function table.remove(list: any, idx: integer)
  -> any
```


```lua
function table.remove(list: any[], idx: integer, end_?: integer)
  -> any
```


```lua
function table.remove(blob: any, idx: integer)
  -> any
```


```lua
function table.remove(blob: any, idx: integer, end_?: integer)
  -> any
```


```lua
function table.remove(dict: any, key: string)
  -> any
```


---

# vim.fn.rename


```lua
function table.rename(from: string, to: string)
  -> integer
```


---

# vim.fn.repeat


```lua
function (expr: any, count: integer)
  -> any
```


---

# vim.fn.resolve


```lua
function table.resolve(filename: string)
  -> string
```


---

# vim.fn.reverse


```lua
function table.reverse(object: <T>[])
  -> <T>[]
```


---

# vim.fn.round


```lua
function table.round(expr: number)
  -> number
```


---

# vim.fn.rpcnotify


```lua
function table.rpcnotify(channel: integer, event: string, ...any)
  -> integer
```


---

# vim.fn.rpcrequest


```lua
function table.rpcrequest(channel: integer, method: string, ...any)
  -> any
```


---

# vim.fn.rpcstart


```lua
function table.rpcstart(prog: string, argv?: any)
  -> any
```


---

# vim.fn.rpcstop


```lua
function table.rpcstop(...any)
  -> any
```


---

# vim.fn.rubyeval


```lua
function table.rubyeval(expr: any)
  -> any
```


---

# vim.fn.screenattr


```lua
function table.screenattr(row: integer, col: integer)
  -> integer
```


---

# vim.fn.screenchar


```lua
function table.screenchar(row: integer, col: integer)
  -> integer
```


---

# vim.fn.screenchars


```lua
function table.screenchars(row: integer, col: integer)
  -> integer[]
```


---

# vim.fn.screencol


```lua
function table.screencol()
  -> integer[]
```


---

# vim.fn.screenpos


```lua
function table.screenpos(winid: integer, lnum: integer, col: integer)
  -> any
```


---

# vim.fn.screenrow


```lua
function table.screenrow()
  -> integer
```


---

# vim.fn.screenstring


```lua
function table.screenstring(row: integer, col: integer)
  -> string
```


---

# vim.fn.search


```lua
function table.search(pattern: string, flags?: string, stopline?: integer, timeout?: integer, skip?: string|function)
  -> integer
```


---

# vim.fn.searchcount


```lua
function table.searchcount(options?: table)
  -> any
```


---

# vim.fn.searchdecl


```lua
function table.searchdecl(name: string, global?: boolean, thisblock?: boolean)
  -> any
```


---

# vim.fn.searchpair


```lua
function table.searchpair(start: string, middle: string, end_: string, flags?: string, skip?: string|function, stopline?: integer, timeout?: integer)
  -> integer
```


---

# vim.fn.searchpairpos


```lua
function table.searchpairpos(start: string, middle: string, end_: string, flags?: string, skip?: string|function, stopline?: integer, timeout?: integer)
  -> [integer, integer]
```


---

# vim.fn.searchpos


```lua
function table.searchpos(pattern: string, flags?: string, stopline?: integer, timeout?: integer, skip?: string|function)
  -> any
```


---

# vim.fn.serverlist


```lua
function table.serverlist()
  -> string[]
```


---

# vim.fn.serverstart


```lua
function table.serverstart(address?: string)
  -> string
```


---

# vim.fn.serverstop


```lua
function table.serverstop(address: string)
  -> integer
```


---

# vim.fn.setbufline


```lua
function table.setbufline(buf: string|integer, lnum: integer, text: string|string[])
  -> integer
```


---

# vim.fn.setbufvar


```lua
function table.setbufvar(buf: string|integer, varname: string, val: any)
  -> any
```


---

# vim.fn.setcellwidths


```lua
function table.setcellwidths(list: any[])
  -> any
```


---

# vim.fn.setcharpos


```lua
function table.setcharpos(expr: string, list: integer[])
  -> any
```


---

# vim.fn.setcharsearch


```lua
function table.setcharsearch(dict: string)
  -> any
```


---

# vim.fn.setcmdline


```lua
function table.setcmdline(str: string, pos?: integer)
  -> integer
```


---

# vim.fn.setcmdpos


```lua
function table.setcmdpos(pos: integer)
  -> any
```


---

# vim.fn.setcursorcharpos


```lua
function table.setcursorcharpos(lnum: string|integer, col?: integer, off?: integer)
  -> any
```


```lua
function table.setcursorcharpos(list: integer[])
  -> any
```


---

# vim.fn.setenv


```lua
function table.setenv(name: string, val: string)
  -> any
```


---

# vim.fn.setfperm


```lua
function table.setfperm(fname: string, mode: string)
  -> any
```


---

# vim.fn.setline


```lua
function table.setline(lnum: string|integer, text: any)
  -> any
```


---

# vim.fn.setloclist


```lua
function table.setloclist(nr: integer, list: any, action?: string, what?: table)
  -> any
```


---

# vim.fn.setmatches


```lua
function table.setmatches(list: any, win?: integer)
  -> any
```


---

# vim.fn.setpos


```lua
function table.setpos(expr: string, list: integer[])
  -> any
```


---

# vim.fn.setqflist


```lua
function table.setqflist(list: vim.quickfix.entry[], action?: string, what?: vim.fn.setqflist.what)
  -> integer
```


---

# vim.fn.setqflist.what

## context


```lua
table?
```


 quickfix list context. See |quickfix-context|

## efm


```lua
string?
```


 errorformat to use when parsing text from
 "lines". If this is not present, then the
 'errorformat' option value is used.
 See |quickfix-parse|

## id


```lua
integer?
```


 quickfix list identifier |quickfix-ID|

## idx


```lua
integer?
```

 index of the current entry in the quickfix
 list specified by "id" or "nr". If set to '$',
 then the last entry in the list is set as the
 current entry. See |quickfix-index|

## items


```lua
vim.quickfix.entry[]?
```


 list of quickfix entries. Same as the {list}
 argument.

## lines


```lua
string[]?
```


 use 'errorformat' to parse a list of lines and
 add the resulting entries to the quickfix list
 {nr} or {id}. Only a |List| value is supported.
 See |quickfix-parse|

## nr


```lua
integer?
```


 list number in the quickfix stack; zero
 means the current quickfix list and "$" means
 the last quickfix list.

## quickfixtextfunc


```lua
function?
```


 function to get the text to display in the
 quickfix window. The value can be the name of
 a function or a funcref or a lambda. Refer
 to |quickfix-window-function| for an explanation
 of how to write the function and an example.

## title


```lua
string?
```


 quickfix list title text. See |quickfix-title|


---

# vim.fn.setreg


```lua
function table.setreg(regname: string, value: any, options?: string)
  -> any
```


---

# vim.fn.settabvar


```lua
function table.settabvar(tabnr: integer, varname: string, val: any)
  -> any
```


---

# vim.fn.settabwinvar


```lua
function table.settabwinvar(tabnr: integer, winnr: integer, varname: string, val: any)
  -> any
```


---

# vim.fn.settagstack


```lua
function table.settagstack(nr: integer, dict: any, action?: string)
  -> any
```


---

# vim.fn.setwinvar


```lua
function table.setwinvar(nr: integer, varname: string, val: any)
  -> any
```


---

# vim.fn.sha256


```lua
function table.sha256(string: string)
  -> string
```


---

# vim.fn.shellescape


```lua
function table.shellescape(string: string, special?: boolean)
  -> string
```


---

# vim.fn.shiftwidth


```lua
function table.shiftwidth(col?: integer)
  -> integer
```


---

# vim.fn.sign

## group


```lua
string
```

## id


```lua
integer
```

## lnum


```lua
integer
```

## name


```lua
string
```

## priority


```lua
integer
```


---

# vim.fn.sign_define


```lua
function table.sign_define(name: string, dict?: vim.fn.sign_define.dict)
  -> -1|0
```


```lua
function table.sign_define(list: vim.fn.sign_define.dict[])
  -> (-1|0)[]
```


---

# vim.fn.sign_define.dict

## culhl


```lua
string?
```

## icon


```lua
string?
```

## linehl


```lua
string?
```

## numhl


```lua
string?
```

## text


```lua
string
```

## texthl


```lua
string?
```


---

# vim.fn.sign_getdefined


```lua
function table.sign_getdefined(name?: string)
  -> vim.fn.sign_getdefined.ret.item[]
```


---

# vim.fn.sign_getdefined.ret.item

## culhl


```lua
string?
```

## icon


```lua
string?
```

## linehl


```lua
string?
```

## name


```lua
string
```

## numhl


```lua
string?
```

## text


```lua
string
```

## texthl


```lua
string?
```


---

# vim.fn.sign_getplaced


```lua
function table.sign_getplaced(buf?: string|integer, dict?: vim.fn.sign_getplaced.dict)
  -> vim.fn.sign_getplaced.ret.item[]
```


---

# vim.fn.sign_getplaced.dict

## group


```lua
string?
```

## id


```lua
integer?
```

## lnum


```lua
(string|integer)?
```


---

# vim.fn.sign_getplaced.ret.item

## bufnr


```lua
integer
```

## signs


```lua
vim.fn.sign[]
```


---

# vim.fn.sign_jump


```lua
function table.sign_jump(id: integer, group: string, buf: string|integer)
  -> integer
```


---

# vim.fn.sign_place


```lua
function table.sign_place(id: integer, group: string, name: string, buf: string|integer, dict?: vim.fn.sign_place.dict)
  -> integer
```


---

# vim.fn.sign_place.dict

## lnum


```lua
(string|integer)?
```

## priority


```lua
integer?
```


---

# vim.fn.sign_placelist


```lua
function table.sign_placelist(list: vim.fn.sign_placelist.list.item[])
  -> integer[]
```


---

# vim.fn.sign_placelist.list.item

## buffer


```lua
string|integer
```

## group


```lua
string?
```

## id


```lua
integer?
```

## lnum


```lua
(string|integer)?
```

## name


```lua
string
```

## priority


```lua
integer?
```


---

# vim.fn.sign_undefine


```lua
function table.sign_undefine(name?: string)
  -> -1|0
```


```lua
function table.sign_undefine(list?: string[])
  -> integer[]
```


---

# vim.fn.sign_unplace


```lua
function table.sign_unplace(group: string, dict?: vim.fn.sign_unplace.dict)
  -> -1|0
```


---

# vim.fn.sign_unplace.dict

## buffer


```lua
(string|integer)?
```

## id


```lua
integer?
```


---

# vim.fn.sign_unplacelist


```lua
function table.sign_unplacelist(list: vim.fn.sign_unplacelist.list.item)
  -> (-1|0)[]
```


---

# vim.fn.sign_unplacelist.list.item

## buffer


```lua
(string|integer)?
```

## group


```lua
string?
```

## id


```lua
integer?
```


---

# vim.fn.simplify


```lua
function table.simplify(filename: string)
  -> string
```


---

# vim.fn.sin


```lua
function table.sin(expr: number)
  -> number
```


---

# vim.fn.sinh


```lua
function table.sinh(expr: number)
  -> any
```


---

# vim.fn.slice


```lua
function table.slice(expr: any, start: integer, end_?: integer)
  -> any
```


---

# vim.fn.sockconnect


```lua
function table.sockconnect(mode: string, address: string, opts?: table)
  -> any
```


---

# vim.fn.sort


```lua
function table.sort(list: <T>[], how?: string|function, dict?: any)
  -> <T>[]
```


---

# vim.fn.soundfold


```lua
function table.soundfold(word: string)
  -> string
```


---

# vim.fn.spellbadword


```lua
function table.spellbadword(sentence?: string)
  -> any
```


---

# vim.fn.spellsuggest


```lua
function table.spellsuggest(word: string, max?: integer, capital?: boolean)
  -> string[]
```


---

# vim.fn.split


```lua
function table.split(string: string, pattern?: string, keepempty?: boolean)
  -> string[]
```


---

# vim.fn.sqrt


```lua
function table.sqrt(expr: number)
  -> any
```


---

# vim.fn.srand


```lua
function table.srand(expr?: number)
  -> any
```


---

# vim.fn.state


```lua
function table.state(what?: string)
  -> any
```


---

# vim.fn.stdioopen


```lua
function table.stdioopen(opts: table)
  -> any
```


---

# vim.fn.stdpath


```lua
function table.stdpath(what: 'cache'|'config'|'config_dirs'|'data'|'data_dirs'...(+3))
  -> string|string[]
```


```lua
function table.stdpath(what: 'cache'|'config'|'data'|'log'|'run'...(+1))
  -> string
```


```lua
function table.stdpath(what: 'config_dirs'|'data_dirs')
  -> string[]
```


---

# vim.fn.str2float


```lua
function table.str2float(string: string, quoted?: boolean)
  -> any
```


---

# vim.fn.str2list


```lua
function table.str2list(string: string, utf8?: boolean)
  -> any
```


---

# vim.fn.str2nr


```lua
function table.str2nr(string: string, base?: integer)
  -> any
```


---

# vim.fn.strcharlen


```lua
function table.strcharlen(string: string)
  -> any
```


---

# vim.fn.strcharpart


```lua
function table.strcharpart(src: string, start: integer, len?: integer, skipcc?: boolean)
  -> any
```


---

# vim.fn.strchars


```lua
function table.strchars(string: string, skipcc?: boolean)
  -> integer
```


---

# vim.fn.strdisplaywidth


```lua
function table.strdisplaywidth(string: string, col?: integer)
  -> integer
```


---

# vim.fn.strftime


```lua
function table.strftime(format: string, time?: number)
  -> string
```


---

# vim.fn.strgetchar


```lua
function table.strgetchar(str: string, index: integer)
  -> integer
```


---

# vim.fn.stridx


```lua
function table.stridx(haystack: string, needle: string, start?: integer)
  -> integer
```


---

# vim.fn.string


```lua
function table.string(expr: any)
  -> string
```


---

# vim.fn.strlen


```lua
function table.strlen(string: string)
  -> integer
```


---

# vim.fn.strpart


```lua
function table.strpart(src: string, start: integer, len?: integer, chars?: 0|1)
  -> string
```


---

# vim.fn.strptime


```lua
function table.strptime(format: string, timestring: string)
  -> integer
```


---

# vim.fn.strridx


```lua
function table.strridx(haystack: string, needle: string, start?: integer)
  -> integer
```


---

# vim.fn.strtrans


```lua
function table.strtrans(string: string)
  -> string
```


---

# vim.fn.strutf16len


```lua
function table.strutf16len(string: string, countcc?: 0|1)
  -> integer
```


---

# vim.fn.strwidth


```lua
function table.strwidth(string: string)
  -> integer
```


---

# vim.fn.submatch


```lua
function table.submatch(nr: integer, list?: any)
  -> string
```


```lua
function table.submatch(nr: integer, list: integer)
  -> string|string[]
```


---

# vim.fn.substitute


```lua
function table.substitute(string: string, pat: string, sub: string, flags: string)
  -> string
```


---

# vim.fn.swapfilelist


```lua
function table.swapfilelist()
  -> string[]
```


---

# vim.fn.swapinfo


```lua
function table.swapinfo(fname: string)
  -> any
```


---

# vim.fn.swapname


```lua
function table.swapname(buf: string|integer)
  -> string
```


---

# vim.fn.synID


```lua
function table.synID(lnum: string|integer, col: integer, trans: 0|1)
  -> integer
```


---

# vim.fn.synIDattr


```lua
function table.synIDattr(synID: integer, what: string, mode?: string)
  -> string
```


---

# vim.fn.synIDtrans


```lua
function table.synIDtrans(synID: integer)
  -> integer
```


---

# vim.fn.synconcealed


```lua
function table.synconcealed(lnum: string|integer, col: integer)
  -> [integer, string, integer]
```


---

# vim.fn.synstack


```lua
function table.synstack(lnum: string|integer, col: integer)
  -> integer[]
```


---

# vim.fn.system


```lua
function table.system(cmd: string|string[], input?: string|integer|string[])
  -> string
```


---

# vim.fn.systemlist


```lua
function table.systemlist(cmd: string|string[], input?: string|integer|string[], keepempty?: integer)
  -> string[]
```


---

# vim.fn.tabpagebuflist


```lua
function table.tabpagebuflist(arg?: integer)
  -> any
```


---

# vim.fn.tabpagenr


```lua
function table.tabpagenr(arg?: '#'|'$')
  -> integer
```


---

# vim.fn.tabpagewinnr


```lua
function table.tabpagewinnr(tabarg: integer, arg?: '#'|'$')
  -> integer
```


---

# vim.fn.tagfiles


```lua
function table.tagfiles()
  -> string[]
```


---

# vim.fn.taglist


```lua
function table.taglist(expr: any, filename?: string)
  -> any
```


---

# vim.fn.tan


```lua
function table.tan(expr: number)
  -> number
```


---

# vim.fn.tanh


```lua
function table.tanh(expr: number)
  -> number
```


---

# vim.fn.tempname


```lua
function table.tempname()
  -> string
```


---

# vim.fn.termopen


```lua
function table.termopen(cmd: string|string[], opts?: table)
  -> integer
```


---

# vim.fn.timer_info


```lua
function table.timer_info(id?: integer)
  -> any
```


---

# vim.fn.timer_pause


```lua
function table.timer_pause(timer: integer, paused: boolean)
  -> any
```


---

# vim.fn.timer_start


```lua
function table.timer_start(time: number, callback: string|function, options?: table)
  -> any
```


---

# vim.fn.timer_stop


```lua
function table.timer_stop(timer: integer)
  -> any
```


---

# vim.fn.timer_stopall


```lua
function table.timer_stopall()
  -> any
```


---

# vim.fn.tolower


```lua
function table.tolower(expr: string)
  -> string
```


---

# vim.fn.toupper


```lua
function table.toupper(expr: string)
  -> string
```


---

# vim.fn.tr


```lua
function table.tr(src: string, fromstr: string, tostr: string)
  -> string
```


---

# vim.fn.trim


```lua
function table.trim(text: string, mask?: string, dir?: 0|1|2)
  -> string
```


---

# vim.fn.trunc


```lua
function table.trunc(expr: number)
  -> integer
```


---

# vim.fn.type


```lua
function table.type(expr: any)
  -> integer
```


---

# vim.fn.undofile


```lua
function table.undofile(name: string)
  -> string
```


---

# vim.fn.undotree


```lua
function table.undotree(buf?: string|integer)
  -> vim.fn.undotree.ret
```


---

# vim.fn.undotree.entry

## alt


```lua
vim.fn.undotree.entry[]?
```


 Alternate entry.  This is again a List of undo
 blocks.  Each item may again have an "alt"
 item.

## curhead


```lua
integer?
```


 Only appears in the item that is the last one
 that was undone.  This marks the current
 position in the undo tree, the block that will
 be used by a redo command.  When nothing was
 undone after the last change this item will
 not appear anywhere.

## newhead


```lua
integer?
```


 Only appears in the item that is the last one
 that was added.  This marks the last change
 and where further changes will be added.

## save


```lua
integer?
```


 Only appears on the last block before a file
 write.  The number is the write count.  The
 first write has number 1, the last one the
 "save_last" mentioned above.

## seq


```lua
integer
```


 Undo sequence number.  Same as what appears in
 \|:undolist|.

## time


```lua
integer
```


 Timestamp when the change happened.  Use
 \|strftime()| to convert to something readable.


---

# vim.fn.undotree.ret

## entries


```lua
vim.fn.undotree.entry[]
```


 A list of dictionaries with information about
 undo blocks.

## save_cur


```lua
integer
```


 Number of the current position in the undo
 tree.

## save_last


```lua
integer
```


 Number of the last file write.  Zero when no
 write yet.

## seq_cur


```lua
integer
```


 The sequence number of the current position in
 the undo tree.  This differs from "seq_last"
 when some changes were undone.

## seq_last


```lua
integer
```


 The highest undo sequence number used.

## synced


```lua
integer
```


 Non-zero when the last undo block was synced.
 This happens when waiting from input from the
 user.  See |undo-blocks|.

## time_cur


```lua
integer
```


 Time last used for |:earlier| and related
 commands.  Use |strftime()| to convert to
 something readable.


---

# vim.fn.uniq


```lua
function table.uniq(list: any, func?: any, dict?: any)
  -> 0|any[]
```


---

# vim.fn.utf16idx


```lua
function table.utf16idx(string: string, idx: integer, countcc?: boolean, charidx?: boolean)
  -> integer
```


---

# vim.fn.values


```lua
function table.values(dict: any)
  -> any
```


---

# vim.fn.virtcol


```lua
function table.virtcol(expr: string|any[], list?: boolean, winid?: integer)
  -> any
```


---

# vim.fn.virtcol2col


```lua
function table.virtcol2col(winid: integer, lnum: integer, col: integer)
  -> integer
```


---

# vim.fn.visualmode


```lua
function table.visualmode(expr?: boolean)
  -> string
```


---

# vim.fn.wait


```lua
function table.wait(timeout: integer, condition: any, interval?: number)
  -> any
```


---

# vim.fn.wildmenumode


```lua
function table.wildmenumode()
  -> any
```


---

# vim.fn.win_execute


```lua
function table.win_execute(id: integer, command: string, silent?: boolean)
  -> any
```


---

# vim.fn.win_findbuf


```lua
function table.win_findbuf(bufnr: integer)
  -> integer[]
```


---

# vim.fn.win_getid


```lua
function table.win_getid(win?: integer, tab?: integer)
  -> integer
```


---

# vim.fn.win_gettype


```lua
function table.win_gettype(nr?: integer)
  -> ''|'autocmd'|'command'|'loclist'|'popup'...(+3)
```


---

# vim.fn.win_gotoid


```lua
function table.win_gotoid(expr: integer)
  -> 0|1
```


---

# vim.fn.win_id2tabwin


```lua
function table.win_id2tabwin(expr: integer)
  -> any
```


---

# vim.fn.win_id2win


```lua
function table.win_id2win(expr: integer)
  -> integer
```


---

# vim.fn.win_move_separator


```lua
function table.win_move_separator(nr: integer, offset: integer)
  -> any
```


---

# vim.fn.win_move_statusline


```lua
function table.win_move_statusline(nr: integer, offset: integer)
  -> any
```


---

# vim.fn.win_screenpos


```lua
function table.win_screenpos(nr: integer)
  -> any
```


---

# vim.fn.win_splitmove


```lua
function table.win_splitmove(nr: integer, target: integer, options?: table)
  -> any
```


---

# vim.fn.winbufnr


```lua
function table.winbufnr(nr: integer)
  -> integer
```


---

# vim.fn.wincol


```lua
function table.wincol()
  -> integer
```


---

# vim.fn.windowsversion


```lua
function table.windowsversion()
  -> string
```


---

# vim.fn.winheight


```lua
function table.winheight(nr: integer)
  -> integer
```


---

# vim.fn.winlayout


```lua
function table.winlayout(tabnr?: integer)
  -> vim.fn.winlayout.branch|vim.fn.winlayout.empty|vim.fn.winlayout.leaf
```


---

# vim.fn.winlayout.branch

## [1]


```lua
"col"|"row"
```

Node type

## [2]


```lua
(vim.fn.winlayout.branch|vim.fn.winlayout.leaf)[]
```

children


---

# vim.fn.winlayout.empty


---

# vim.fn.winlayout.leaf

## [1]


```lua
"leaf"
```

Node type

## [2]


```lua
integer
```

winid


---

# vim.fn.winlayout.ret


---

# vim.fn.winline


```lua
function table.winline()
  -> integer
```


---

# vim.fn.winnr


```lua
function table.winnr(arg?: string|integer)
  -> integer
```


---

# vim.fn.winrestcmd


```lua
function table.winrestcmd()
  -> string
```


---

# vim.fn.winrestview


```lua
function table.winrestview(dict: vim.fn.winrestview.dict)
  -> any
```


---

# vim.fn.winrestview.dict

## col


```lua
integer?
```

## coladd


```lua
integer?
```

## curswant


```lua
integer?
```

## leftcol


```lua
integer?
```

## lnum


```lua
integer?
```

## skipcol


```lua
integer?
```

## topfill


```lua
integer?
```

## topline


```lua
integer?
```


---

# vim.fn.winsaveview


```lua
function table.winsaveview()
  -> vim.fn.winsaveview.ret
```


---

# vim.fn.winsaveview.ret

## col


```lua
integer
```

## coladd


```lua
integer
```

## curswant


```lua
integer
```

## leftcol


```lua
integer
```

## lnum


```lua
integer
```

## skipcol


```lua
integer
```

## topfill


```lua
integer
```

## topline


```lua
integer
```


---

# vim.fn.winwidth


```lua
function table.winwidth(nr: integer)
  -> integer
```


---

# vim.fn.wordcount


```lua
function table.wordcount()
  -> any
```


---

# vim.fn.writefile


```lua
function table.writefile(object: any, fname: string, flags?: string)
  -> any
```


---

# vim.fn.xor


```lua
function table.xor(expr: integer, expr1: integer)
  -> integer
```


---

# vim.fs


```lua
table
```


---

# vim.fs.dir.Opts

## depth


```lua
integer?
```

 @inlinedoc

 How deep the traverse.
 (default: `1`)

## follow


```lua
boolean?
```


 Follow symbolic links.
 (default: `false`)

## skip


```lua
(fun(dir_name: string):boolean)?
```


 Predicate to control traversal.
 Return false to stop searching the current directory.
 Only useful when depth > 1
 Return an iterator over the items located in {path}


---

# vim.fs.find.Opts

## follow


```lua
boolean?
```


 Follow symbolic links.
 (default: `false`)

## limit


```lua
number?
```


 Stop the search after finding this many matches.
 Use `math.huge` to place no limit on the number of matches.
 (default: `1`)

## path


```lua
string?
```

 @inlinedoc

 Path to begin searching from. If
 omitted, the |current-directory| is used.

## stop


```lua
string?
```


 Stop searching when this directory is reached.
 The directory itself is not searched.

## type


```lua
string?
```


 Find only items of the given type.
 If omitted, all items that match {names} are included.

## upward


```lua
boolean?
```


 Search upward through parent directories.
 Otherwise, search through child directories (recursively).
 (default: `false`)


---

# vim.fs.normalize.Opts

## _fast


```lua
boolean?
```


## expand_env


```lua
boolean?
```

 @inlinedoc

 Expand environment variables.
 (default: `true`)

## win


```lua
boolean?
```


 Path is a Windows path.
 (default: `true` in Windows, `false` otherwise)


---

# vim.fs.rm.Opts

## force


```lua
boolean?
```


 Ignore nonexistent files and arguments

## recursive


```lua
boolean?
```

 @inlinedoc

 Remove directories and their contents recursively


---

# vim.func


```lua
table
```


---

# vim.func.MemoObj


---

# vim.funcref


```lua
function vim.funcref(viml_func_name: any)
  -> unknown
```


---

# vim.g


```lua
vim.var_accessor
```


---

# vim.g.colors_name


```lua
string
```


---

# vim.g.did_load_filetypes


```lua
integer
```


---

# vim.g.ft_ignore_pat


```lua
string
```


---

# vim.g.loaded_2html_plugin


```lua
boolean
```


---

# vim.g.loaded_man


```lua
boolean
```


---

# vim.g.mapleader


```lua
string
```


---

# vim.g.termfeatures


```lua
TermFeatures
```


```lua
TermFeatures
```


```lua
TermFeatures
```


---

# vim.glob


```lua
table
```


---

# vim.go


```lua
table
```


---

# vim.go.acd


```lua
boolean
```


---

# vim.go.allowrevins


```lua
boolean
```


---

# vim.go.ambiwidth


```lua
'double'|'single'
```


---

# vim.go.ambw


```lua
'double'|'single'
```


---

# vim.go.ar


```lua
boolean
```


---

# vim.go.arabicshape


```lua
boolean
```


---

# vim.go.ari


```lua
boolean
```


---

# vim.go.arshape


```lua
boolean
```


---

# vim.go.autochdir


```lua
boolean
```


---

# vim.go.autoread


```lua
boolean
```


---

# vim.go.autowrite


```lua
boolean
```


---

# vim.go.autowriteall


```lua
boolean
```


---

# vim.go.aw


```lua
boolean
```


---

# vim.go.awa


```lua
boolean
```


---

# vim.go.background


```lua
'dark'|'light'
```


---

# vim.go.backspace


```lua
string
```


---

# vim.go.backup


```lua
boolean
```


---

# vim.go.backupcopy


```lua
string
```


---

# vim.go.backupdir


```lua
string
```


---

# vim.go.backupext


```lua
string
```


---

# vim.go.backupskip


```lua
string
```


---

# vim.go.bdir


```lua
string
```


---

# vim.go.belloff


```lua
string
```


---

# vim.go.bex


```lua
string
```


---

# vim.go.bg


```lua
'dark'|'light'
```


---

# vim.go.bk


```lua
boolean
```


---

# vim.go.bkc


```lua
string
```


---

# vim.go.bo


```lua
string
```


---

# vim.go.breakat


```lua
string
```


---

# vim.go.brk


```lua
string
```


---

# vim.go.bs


```lua
string
```


---

# vim.go.bsk


```lua
string
```


---

# vim.go.casemap


```lua
string
```


---

# vim.go.cb


```lua
string
```


---

# vim.go.ccv


```lua
string
```


---

# vim.go.cd


```lua
string
```


---

# vim.go.cdh


```lua
boolean
```


---

# vim.go.cdhome


```lua
boolean
```


---

# vim.go.cdpath


```lua
string
```


---

# vim.go.cedit


```lua
string
```


---

# vim.go.cf


```lua
boolean
```


---

# vim.go.ch


```lua
integer
```


---

# vim.go.charconvert


```lua
string
```


---

# vim.go.cia


```lua
string
```


---

# vim.go.clipboard


```lua
string
```


---

# vim.go.cmdheight


```lua
integer
```


---

# vim.go.cmdwinheight


```lua
integer
```


---

# vim.go.cmp


```lua
string
```


---

# vim.go.co


```lua
integer
```


---

# vim.go.columns


```lua
integer
```


---

# vim.go.completeitemalign


```lua
string
```


---

# vim.go.completeopt


```lua
string
```


---

# vim.go.confirm


```lua
boolean
```


---

# vim.go.cot


```lua
string
```


---

# vim.go.cpo


```lua
string
```


---

# vim.go.cpoptions


```lua
string
```


---

# vim.go.cwh


```lua
integer
```


---

# vim.go.debug


```lua
string
```


---

# vim.go.deco


```lua
boolean
```


---

# vim.go.def


```lua
string
```


---

# vim.go.define


```lua
string
```


---

# vim.go.delcombine


```lua
boolean
```


---

# vim.go.dex


```lua
string
```


---

# vim.go.dg


```lua
boolean
```


---

# vim.go.dict


```lua
string
```


---

# vim.go.dictionary


```lua
string
```


---

# vim.go.diffexpr


```lua
string
```


---

# vim.go.diffopt


```lua
string
```


---

# vim.go.digraph


```lua
boolean
```


---

# vim.go.dip


```lua
string
```


---

# vim.go.dir


```lua
string
```


---

# vim.go.directory


```lua
string
```


---

# vim.go.display


```lua
string
```


---

# vim.go.dy


```lua
string
```


---

# vim.go.ea


```lua
boolean
```


---

# vim.go.ead


```lua
'both'|'hor'|'ver'
```


---

# vim.go.eadirection


```lua
'both'|'hor'|'ver'
```


---

# vim.go.eb


```lua
boolean
```


---

# vim.go.ef


```lua
string
```


---

# vim.go.efm


```lua
string
```


---

# vim.go.ei


```lua
string
```


---

# vim.go.emo


```lua
boolean
```


---

# vim.go.emoji


```lua
boolean
```


---

# vim.go.enc


```lua
string
```


---

# vim.go.encoding


```lua
string
```


---

# vim.go.ep


```lua
string
```


---

# vim.go.equalalways


```lua
boolean
```


---

# vim.go.equalprg


```lua
string
```


---

# vim.go.errorbells


```lua
boolean
```


---

# vim.go.errorfile


```lua
string
```


---

# vim.go.errorformat


```lua
string
```


---

# vim.go.eventignore


```lua
string
```


---

# vim.go.ex


```lua
boolean
```


---

# vim.go.exrc


```lua
boolean
```


---

# vim.go.fcl


```lua
string
```


---

# vim.go.fcs


```lua
string
```


---

# vim.go.fdls


```lua
integer
```


---

# vim.go.fdo


```lua
string
```


---

# vim.go.fencs


```lua
string
```


---

# vim.go.ffs


```lua
string
```


---

# vim.go.ffu


```lua
string
```


---

# vim.go.fic


```lua
boolean
```


---

# vim.go.fileencodings


```lua
string
```


---

# vim.go.fileformats


```lua
string
```


---

# vim.go.fileignorecase


```lua
boolean
```


---

# vim.go.fillchars


```lua
string
```


---

# vim.go.findfunc


```lua
string
```


---

# vim.go.foldclose


```lua
string
```


---

# vim.go.foldlevelstart


```lua
integer
```


---

# vim.go.foldopen


```lua
string
```


---

# vim.go.formatprg


```lua
string
```


---

# vim.go.fp


```lua
string
```


---

# vim.go.fs


```lua
boolean
```


---

# vim.go.fsync


```lua
boolean
```


---

# vim.go.gcr


```lua
string
```


---

# vim.go.gd


```lua
boolean
```


---

# vim.go.gdefault


```lua
boolean
```


---

# vim.go.gfm


```lua
string
```


---

# vim.go.gfn


```lua
string
```


---

# vim.go.gfw


```lua
string
```


---

# vim.go.gp


```lua
string
```


---

# vim.go.grepformat


```lua
string
```


---

# vim.go.grepprg


```lua
string
```


---

# vim.go.guicursor


```lua
string
```


---

# vim.go.guifont


```lua
string
```


---

# vim.go.guifontwide


```lua
string
```


---

# vim.go.helpfile


```lua
string
```


---

# vim.go.helpheight


```lua
integer
```


---

# vim.go.helplang


```lua
string
```


---

# vim.go.hf


```lua
string
```


---

# vim.go.hh


```lua
integer
```


---

# vim.go.hi


```lua
integer
```


---

# vim.go.hid


```lua
boolean
```


---

# vim.go.hidden


```lua
boolean
```


---

# vim.go.history


```lua
integer
```


---

# vim.go.hlg


```lua
string
```


---

# vim.go.hls


```lua
boolean
```


---

# vim.go.hlsearch


```lua
boolean
```


---

# vim.go.ic


```lua
boolean
```


---

# vim.go.icm


```lua
''|'nosplit'|'split'
```


---

# vim.go.icon


```lua
boolean
```


---

# vim.go.iconstring


```lua
string
```


---

# vim.go.ignorecase


```lua
boolean
```


---

# vim.go.inc


```lua
string
```


---

# vim.go.inccommand


```lua
''|'nosplit'|'split'
```


---

# vim.go.include


```lua
string
```


---

# vim.go.incsearch


```lua
boolean
```


---

# vim.go.is


```lua
boolean
```


---

# vim.go.isf


```lua
string
```


---

# vim.go.isfname


```lua
string
```


---

# vim.go.isi


```lua
string
```


---

# vim.go.isident


```lua
string
```


---

# vim.go.isp


```lua
string
```


---

# vim.go.isprint


```lua
string
```


---

# vim.go.joinspaces


```lua
boolean
```


---

# vim.go.jop


```lua
string
```


---

# vim.go.js


```lua
boolean
```


---

# vim.go.jumpoptions


```lua
string
```


---

# vim.go.keymodel


```lua
string
```


---

# vim.go.keywordprg


```lua
string
```


---

# vim.go.km


```lua
string
```


---

# vim.go.kp


```lua
string
```


---

# vim.go.langmap


```lua
string
```


---

# vim.go.langmenu


```lua
string
```


---

# vim.go.langremap


```lua
boolean
```


---

# vim.go.laststatus


```lua
integer
```


---

# vim.go.lazyredraw


```lua
boolean
```


---

# vim.go.lcs


```lua
string
```


---

# vim.go.lines


```lua
integer
```


---

# vim.go.linespace


```lua
integer
```


---

# vim.go.lispwords


```lua
string
```


---

# vim.go.listchars


```lua
string
```


---

# vim.go.lm


```lua
string
```


---

# vim.go.lmap


```lua
string
```


---

# vim.go.loadplugins


```lua
boolean
```


---

# vim.go.lpl


```lua
boolean
```


---

# vim.go.lrm


```lua
boolean
```


---

# vim.go.ls


```lua
integer
```


---

# vim.go.lsp


```lua
integer
```


---

# vim.go.lw


```lua
string
```


---

# vim.go.lz


```lua
boolean
```


---

# vim.go.magic


```lua
boolean
```


---

# vim.go.makeef


```lua
string
```


---

# vim.go.makeencoding


```lua
string
```


---

# vim.go.makeprg


```lua
string
```


---

# vim.go.mat


```lua
integer
```


---

# vim.go.matchtime


```lua
integer
```


---

# vim.go.maxfuncdepth


```lua
integer
```


---

# vim.go.maxmapdepth


```lua
integer
```


---

# vim.go.maxmempattern


```lua
integer
```


---

# vim.go.mef


```lua
string
```


---

# vim.go.menc


```lua
string
```


---

# vim.go.menuitems


```lua
integer
```


---

# vim.go.messagesopt


```lua
string
```


---

# vim.go.mfd


```lua
integer
```


---

# vim.go.mh


```lua
boolean
```


---

# vim.go.mis


```lua
integer
```


---

# vim.go.mkspellmem


```lua
string
```


---

# vim.go.mle


```lua
boolean
```


---

# vim.go.mls


```lua
integer
```


---

# vim.go.mmd


```lua
integer
```


---

# vim.go.mmp


```lua
integer
```


---

# vim.go.modelineexpr


```lua
boolean
```


---

# vim.go.modelines


```lua
integer
```


---

# vim.go.mopt


```lua
string
```


---

# vim.go.more


```lua
boolean
```


---

# vim.go.mouse


```lua
string
```


---

# vim.go.mousef


```lua
boolean
```


---

# vim.go.mousefocus


```lua
boolean
```


---

# vim.go.mousehide


```lua
boolean
```


---

# vim.go.mousem


```lua
'extend'|'popup'|'popup_setpos'
```


---

# vim.go.mousemev


```lua
boolean
```


---

# vim.go.mousemodel


```lua
'extend'|'popup'|'popup_setpos'
```


---

# vim.go.mousemoveevent


```lua
boolean
```


---

# vim.go.mousescroll


```lua
string
```


---

# vim.go.mouset


```lua
integer
```


---

# vim.go.mousetime


```lua
integer
```


---

# vim.go.mp


```lua
string
```


---

# vim.go.msm


```lua
string
```


---

# vim.go.operatorfunc


```lua
string
```


```lua
string
```


```lua
string
```


---

# vim.go.opfunc


```lua
string
```


---

# vim.go.pa


```lua
string
```


---

# vim.go.packpath


```lua
string
```


---

# vim.go.para


```lua
string
```


---

# vim.go.paragraphs


```lua
string
```


---

# vim.go.patchexpr


```lua
string
```


---

# vim.go.patchmode


```lua
string
```


---

# vim.go.path


```lua
string
```


---

# vim.go.pb


```lua
integer
```


---

# vim.go.pex


```lua
string
```


---

# vim.go.ph


```lua
integer
```


---

# vim.go.pm


```lua
string
```


---

# vim.go.pp


```lua
string
```


---

# vim.go.previewheight


```lua
integer
```


---

# vim.go.pumblend


```lua
integer
```


---

# vim.go.pumheight


```lua
integer
```


---

# vim.go.pumwidth


```lua
integer
```


---

# vim.go.pvh


```lua
integer
```


---

# vim.go.pw


```lua
integer
```


---

# vim.go.pyx


```lua
integer
```


---

# vim.go.pyxversion


```lua
integer
```


---

# vim.go.qftf


```lua
string
```


---

# vim.go.quickfixtextfunc


```lua
string
```


---

# vim.go.rdb


```lua
string
```


---

# vim.go.rdt


```lua
integer
```


---

# vim.go.re


```lua
integer
```


---

# vim.go.redrawdebug


```lua
string
```


---

# vim.go.redrawtime


```lua
integer
```


---

# vim.go.regexpengine


```lua
integer
```


---

# vim.go.report


```lua
integer
```


---

# vim.go.revins


```lua
boolean
```


---

# vim.go.ri


```lua
boolean
```


---

# vim.go.rtp


```lua
string
```


---

# vim.go.ru


```lua
boolean
```


---

# vim.go.ruf


```lua
string
```


---

# vim.go.ruler


```lua
boolean
```


---

# vim.go.rulerformat


```lua
string
```


---

# vim.go.runtimepath


```lua
string
```


---

# vim.go.sb


```lua
boolean
```


---

# vim.go.sbo


```lua
string
```


---

# vim.go.sbr


```lua
string
```


---

# vim.go.sc


```lua
boolean
```


---

# vim.go.scrolljump


```lua
integer
```


---

# vim.go.scrolloff


```lua
integer
```


---

# vim.go.scrollopt


```lua
string
```


---

# vim.go.scs


```lua
boolean
```


---

# vim.go.sd


```lua
string
```


---

# vim.go.sdf


```lua
string
```


---

# vim.go.sect


```lua
string
```


---

# vim.go.sections


```lua
string
```


---

# vim.go.sel


```lua
'exclusive'|'inclusive'|'old'
```


---

# vim.go.selection


```lua
'exclusive'|'inclusive'|'old'
```


---

# vim.go.selectmode


```lua
string
```


---

# vim.go.sessionoptions


```lua
string
```


---

# vim.go.sft


```lua
boolean
```


---

# vim.go.sh


```lua
string
```


---

# vim.go.shada


```lua
string
```


---

# vim.go.shadafile


```lua
string
```


---

# vim.go.shcf


```lua
string
```


---

# vim.go.shell


```lua
string
```


---

# vim.go.shellcmdflag


```lua
string
```


---

# vim.go.shellpipe


```lua
string
```


---

# vim.go.shellquote


```lua
string
```


---

# vim.go.shellredir


```lua
string
```


---

# vim.go.shellslash


```lua
boolean
```


---

# vim.go.shelltemp


```lua
boolean
```


---

# vim.go.shellxescape


```lua
string
```


---

# vim.go.shellxquote


```lua
string
```


---

# vim.go.shiftround


```lua
boolean
```


---

# vim.go.shm


```lua
string
```


---

# vim.go.shortmess


```lua
string
```


---

# vim.go.showbreak


```lua
string
```


---

# vim.go.showcmd


```lua
boolean
```


---

# vim.go.showcmdloc


```lua
'last'|'statusline'|'tabline'
```


---

# vim.go.showfulltag


```lua
boolean
```


---

# vim.go.showmatch


```lua
boolean
```


---

# vim.go.showmode


```lua
boolean
```


---

# vim.go.showtabline


```lua
integer
```


---

# vim.go.shq


```lua
string
```


---

# vim.go.sidescroll


```lua
integer
```


---

# vim.go.sidescrolloff


```lua
integer
```


---

# vim.go.siso


```lua
integer
```


---

# vim.go.sj


```lua
integer
```


---

# vim.go.slm


```lua
string
```


---

# vim.go.sloc


```lua
'last'|'statusline'|'tabline'
```


---

# vim.go.sm


```lua
boolean
```


---

# vim.go.smartcase


```lua
boolean
```


---

# vim.go.smarttab


```lua
boolean
```


---

# vim.go.smd


```lua
boolean
```


---

# vim.go.so


```lua
integer
```


---

# vim.go.sol


```lua
boolean
```


---

# vim.go.sp


```lua
string
```


---

# vim.go.spellsuggest


```lua
string
```


---

# vim.go.spk


```lua
'cursor'|'screen'|'topline'
```


---

# vim.go.splitbelow


```lua
boolean
```


---

# vim.go.splitkeep


```lua
'cursor'|'screen'|'topline'
```


---

# vim.go.splitright


```lua
boolean
```


---

# vim.go.spr


```lua
boolean
```


---

# vim.go.sps


```lua
string
```


---

# vim.go.sr


```lua
boolean
```


---

# vim.go.srr


```lua
string
```


---

# vim.go.ss


```lua
integer
```


---

# vim.go.ssl


```lua
boolean
```


---

# vim.go.ssop


```lua
string
```


---

# vim.go.sta


```lua
boolean
```


---

# vim.go.stal


```lua
integer
```


---

# vim.go.startofline


```lua
boolean
```


---

# vim.go.statusline


```lua
string
```


---

# vim.go.stl


```lua
string
```


---

# vim.go.stmp


```lua
boolean
```


---

# vim.go.su


```lua
string
```


---

# vim.go.suffixes


```lua
string
```


---

# vim.go.swb


```lua
string
```


---

# vim.go.switchbuf


```lua
string
```


---

# vim.go.sxe


```lua
string
```


---

# vim.go.sxq


```lua
string
```


---

# vim.go.tabclose


```lua
string
```


---

# vim.go.tabline


```lua
string
```


---

# vim.go.tabpagemax


```lua
integer
```


---

# vim.go.tag


```lua
string
```


---

# vim.go.tagbsearch


```lua
boolean
```


---

# vim.go.tagcase


```lua
'followic'|'followscs'|'ignore'|'match'|'smart'
```


---

# vim.go.taglength


```lua
integer
```


---

# vim.go.tagrelative


```lua
boolean
```


---

# vim.go.tags


```lua
string
```


---

# vim.go.tagstack


```lua
boolean
```


---

# vim.go.tal


```lua
string
```


---

# vim.go.tbidi


```lua
boolean
```


---

# vim.go.tbs


```lua
boolean
```


---

# vim.go.tc


```lua
'followic'|'followscs'|'ignore'|'match'|'smart'
```


---

# vim.go.tcl


```lua
string
```


---

# vim.go.termbidi


```lua
boolean
```


---

# vim.go.termguicolors


```lua
boolean
```


---

# vim.go.termpastefilter


```lua
string
```


---

# vim.go.termsync


```lua
boolean
```


---

# vim.go.tgc


```lua
boolean
```


---

# vim.go.tgst


```lua
boolean
```


---

# vim.go.thesaurus


```lua
string
```


---

# vim.go.thesaurusfunc


```lua
string
```


---

# vim.go.tildeop


```lua
boolean
```


---

# vim.go.timeout


```lua
boolean
```


---

# vim.go.timeoutlen


```lua
integer
```


---

# vim.go.title


```lua
boolean
```


---

# vim.go.titlelen


```lua
integer
```


---

# vim.go.titleold


```lua
string
```


---

# vim.go.titlestring


```lua
string
```


---

# vim.go.tl


```lua
integer
```


---

# vim.go.tm


```lua
integer
```


---

# vim.go.to


```lua
boolean
```


---

# vim.go.top


```lua
boolean
```


---

# vim.go.tpf


```lua
string
```


---

# vim.go.tpm


```lua
integer
```


---

# vim.go.tr


```lua
boolean
```


---

# vim.go.tsr


```lua
string
```


---

# vim.go.tsrfu


```lua
string
```


---

# vim.go.ttimeout


```lua
boolean
```


---

# vim.go.ttimeoutlen


```lua
integer
```


---

# vim.go.ttm


```lua
integer
```


---

# vim.go.uc


```lua
integer
```


---

# vim.go.udir


```lua
string
```


---

# vim.go.ul


```lua
integer
```


---

# vim.go.undodir


```lua
string
```


---

# vim.go.undolevels


```lua
integer
```


---

# vim.go.undoreload


```lua
integer
```


---

# vim.go.updatecount


```lua
integer
```


---

# vim.go.updatetime


```lua
integer
```


---

# vim.go.ur


```lua
integer
```


---

# vim.go.ut


```lua
integer
```


---

# vim.go.vb


```lua
boolean
```


---

# vim.go.vbs


```lua
integer
```


---

# vim.go.vdir


```lua
string
```


---

# vim.go.ve


```lua
string
```


---

# vim.go.verbose


```lua
integer
```


---

# vim.go.verbosefile


```lua
string
```


---

# vim.go.vfile


```lua
string
```


---

# vim.go.viewdir


```lua
string
```


---

# vim.go.viewoptions


```lua
string
```


---

# vim.go.virtualedit


```lua
string
```


---

# vim.go.visualbell


```lua
boolean
```


---

# vim.go.vop


```lua
string
```


---

# vim.go.wa


```lua
boolean
```


---

# vim.go.wak


```lua
'menu'|'no'|'yes'
```


---

# vim.go.warn


```lua
boolean
```


---

# vim.go.wb


```lua
boolean
```


---

# vim.go.wbr


```lua
string
```


---

# vim.go.wc


```lua
integer
```


---

# vim.go.wcm


```lua
integer
```


---

# vim.go.wd


```lua
integer
```


---

# vim.go.wh


```lua
integer
```


---

# vim.go.whichwrap


```lua
string
```


---

# vim.go.wi


```lua
integer
```


---

# vim.go.wic


```lua
boolean
```


---

# vim.go.wig


```lua
string
```


---

# vim.go.wildchar


```lua
integer
```


---

# vim.go.wildcharm


```lua
integer
```


---

# vim.go.wildignore


```lua
string
```


---

# vim.go.wildignorecase


```lua
boolean
```


---

# vim.go.wildmenu


```lua
boolean
```


---

# vim.go.wildmode


```lua
string
```


---

# vim.go.wildoptions


```lua
string
```


---

# vim.go.wim


```lua
string
```


---

# vim.go.winaltkeys


```lua
'menu'|'no'|'yes'
```


---

# vim.go.winbar


```lua
string
```


---

# vim.go.winborder


```lua
''|'bold'|'double'|'none'|'rounded'...(+3)
```


---

# vim.go.window


```lua
integer
```


---

# vim.go.winheight


```lua
integer
```


---

# vim.go.winminheight


```lua
integer
```


---

# vim.go.winminwidth


```lua
integer
```


---

# vim.go.winwidth


```lua
integer
```


---

# vim.go.wiw


```lua
integer
```


---

# vim.go.wmh


```lua
integer
```


---

# vim.go.wmnu


```lua
boolean
```


---

# vim.go.wmw


```lua
integer
```


---

# vim.go.wop


```lua
string
```


---

# vim.go.wrapscan


```lua
boolean
```


---

# vim.go.write


```lua
boolean
```


---

# vim.go.writeany


```lua
boolean
```


---

# vim.go.writebackup


```lua
boolean
```


---

# vim.go.writedelay


```lua
integer
```


---

# vim.go.ws


```lua
boolean
```


---

# vim.go.ww


```lua
string
```


---

# vim.gsplit


```lua
function vim.gsplit(s: string, sep: string, opts?: vim.gsplit.Opts)
  -> fun():string?
```


---

# vim.gsplit.Opts

## plain


```lua
boolean?
```

 @inlinedoc

 Use `sep` literally (as in string.find).

## trimempty


```lua
boolean?
```


 Discard empty segments at start and end of the sequence.


---

# vim.health


```lua
table
```


---

# vim.highlight


```lua
table
```


---

# vim.hl


```lua
table
```


---

# vim.hl.range.Opts

## inclusive


```lua
boolean?
```


 Indicates whether the range is end-inclusive
 (default: `false`)

## priority


```lua
integer?
```


 Highlight priority
 (default: `vim.hl.priorities.user`)

## regtype


```lua
string?
```

 @inlinedoc

 Type of range. See [getregtype()]
 (default: `'v'` i.e. charwise)

## timeout


```lua
integer?
```


 Time in ms before highlight is cleared
 (default: -1 no timeout)


---

# vim.iconv


```lua
function vim.iconv(str: string, from: string, to: string, opts: any)
  -> string?
```


---

# vim.in_fast_event


```lua
function vim.in_fast_event()
```


---

# vim.inspect


```lua
fun(x: any, opts?: vim.inspect.Opts):string
```


---

# vim.inspect.Opts

 @nodoc

## depth


```lua
integer?
```

## newline


```lua
string?
```

## process


```lua
(fun(item: any, path: string[]):any)?
```


---

# vim.inspect_pos


```lua
function vim.inspect_pos(bufnr?: integer, row?: integer, col?: integer, filter?: vim._inspector.Filter)
  -> { treesitter: table, syntax: table, extmarks: table, semantic_tokens: table, buffer: integer, col: integer, row: integer }
```


---

# vim.is_callable


```lua
function vim.is_callable(f: any)
  -> boolean
```


---

# vim.isarray


```lua
function vim.isarray(t?: table)
  -> boolean
```


---

# vim.islist


```lua
function vim.islist(t?: table)
  -> boolean
```


---

# vim.iter


```lua
IterMod
```


---

# vim.json


```lua
table
```


---

# vim.json.decode


```lua
function vim.json.decode(str: string, opts?: table<string, any>)
  -> any
```


---

# vim.json.encode


```lua
function vim.json.encode(obj: any, opts?: table<string, any>)
  -> string
```


---

# vim.keycode


```lua
function vim.keycode(str: string)
  -> string
```


---

# vim.keymap


```lua
table
```


---

# vim.keymap.del.Opts

## buffer


```lua
(boolean|integer)?
```

 @inlinedoc

 Remove a mapping from the given buffer.
 When `0` or `true`, use the current buffer.


---

# vim.keymap.set.Opts

 Table of |:map-arguments|.
 Same as |nvim_set_keymap()| {opts}, except:
 - {replace_keycodes} defaults to `true` if "expr" is `true`.

 Also accepts:

## buffer


```lua
(boolean|integer)?
```

 @inlinedoc

 Creates buffer-local mapping, `0` or `true` for current buffer.

## callback


```lua
function?
```

## desc


```lua
string?
```

## expr


```lua
boolean?
```

## noremap


```lua
boolean?
```

## nowait


```lua
boolean?
```

## remap


```lua
boolean?
```


 Make the mapping recursive. Inverse of {noremap}.
 (Default: `false`)

## replace_keycodes


```lua
boolean?
```

## script


```lua
boolean?
```

## silent


```lua
boolean?
```

## unique


```lua
boolean?
```


---

# vim.list_contains


```lua
function vim.list_contains(t: table, value: any)
  -> boolean
```


---

# vim.list_extend


```lua
function vim.list_extend(dst: <T:table>, src: table, start?: integer, finish?: integer)
  -> dst: <T:table>
```


---

# vim.list_slice


```lua
function vim.list_slice(list: <T>[], start: integer|nil, finish: integer|nil)
  -> Copy: <T>[]
```


---

# vim.loader


```lua
table
```


---

# vim.loader.CacheEntry


---

# vim.loader.CacheHash


---

# vim.loader.ModuleInfo

## modname


```lua
string
```


 Name of the module

## modpath


```lua
string
```

 @inlinedoc

 Path of the module

## stat


```lua
(uv.fs_stat.result)?
```


 The fs_stat of the module path. Won't be returned for `modname="*"`


---

# vim.loader.Stats


---

# vim.loader._profile.Opts

## loaders


```lua
boolean?
```

Add profiling to the loaders


---

# vim.loader.find.Opts

## all


```lua
boolean?
```


 Search for all matches.
 (default: `false`)

## paths


```lua
string[]?
```


 Extra paths to search for modname
 (default: `{}`)

## patterns


```lua
string[]?
```


 List of patterns to use when searching for modules.
 A pattern is a string added to the basename of the Lua module being searched.
 (default: `{"/init.lua", ".lua"}`)

## rtp


```lua
boolean?
```

 @inlinedoc

 Search for modname in the runtime path.
 (default: `true`)


---

# vim.log


```lua
table
```


---

# vim.log.levels


---

# vim.loop


```lua
uv
```


---

# vim.lpeg


```lua
table
```


---

# vim.lpeg.B


```lua
function vim.lpeg.B(pattern: boolean|string|integer|table|vim.lpeg.Pattern)
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.C


```lua
function vim.lpeg.C(patt: boolean|string|integer|function|table...(+1))
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Capture


---

# vim.lpeg.Carg


```lua
function vim.lpeg.Carg(n: integer)
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Cb


```lua
function vim.lpeg.Cb(name: any)
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Cc


```lua
function vim.lpeg.Cc(...any)
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Cf


```lua
function vim.lpeg.Cf(patt: boolean|string|integer|function|table...(+1), func: fun(acc: any, newvalue: any))
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Cg


```lua
function vim.lpeg.Cg(patt: boolean|string|integer|function|table...(+1), name?: string)
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Cmt


```lua
function vim.lpeg.Cmt(patt: boolean|string|integer|function|table...(+1), fn: fun(s: string, i: integer, ...any):(position: boolean|integer, ...any))
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Cp


```lua
function vim.lpeg.Cp()
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Cs


```lua
function vim.lpeg.Cs(patt: boolean|string|integer|function|table...(+1))
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Ct


```lua
function vim.lpeg.Ct(patt: boolean|string|integer|function|table...(+1))
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Locale

 @nodoc

## alnum


```lua
userdata
```

## alpha


```lua
userdata
```

## cntrl


```lua
userdata
```

## digit


```lua
userdata
```

## graph


```lua
userdata
```

## lower


```lua
userdata
```

## print


```lua
userdata
```

## punct


```lua
userdata
```

## space


```lua
userdata
```

## upper


```lua
userdata
```

## xdigit


```lua
userdata
```


---

# vim.lpeg.P


```lua
function vim.lpeg.P(value: boolean|string|integer|function|table...(+1))
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.Pattern

 @nodoc

## match


```lua
(method) vim.lpeg.Pattern:match(subject: string, init?: integer, ...any)
  -> ...any
```

 Matches the given `pattern` against the `subject` string. If the match succeeds, returns the
 index in the subject of the first character after the match, or the captured values (if the
 pattern captured any value). An optional numeric argument `init` makes the match start at
 that position in the subject string. As usual in Lua libraries, a negative value counts from the end.
 Unlike typical pattern-matching functions, `match` works only in anchored mode; that is, it tries
 to match the pattern with a prefix of the given subject string (at position `init`), not with
 an arbitrary substring of the subject. So, if we want to find a pattern anywhere in a string,
 we must either write a loop in Lua or write a pattern that matches anywhere.

 Example:

 ```lua
 local pattern = lpeg.R('az') ^ 1 * -1
 assert(pattern:match('hello') == 6)
 assert(lpeg.match(pattern, 'hello') == 6)
 assert(pattern:match('1 hello') == nil)
 ```


---

# vim.lpeg.R


```lua
function vim.lpeg.R(...string)
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.S


```lua
function vim.lpeg.S(string: string)
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.V


```lua
function vim.lpeg.V(v: boolean|string|number|function|table...(+3))
  -> vim.lpeg.Pattern
```


---

# vim.lpeg.locale


```lua
function vim.lpeg.locale(tab?: table)
  -> vim.lpeg.Locale
```


---

# vim.lpeg.match


```lua
function vim.lpeg.match(pattern: boolean|string|integer|function|table...(+1), subject: string, init?: integer, ...any)
  -> ...any
```


---

# vim.lpeg.setmaxstack


```lua
function vim.lpeg.setmaxstack(max: integer)
```


---

# vim.lpeg.type


```lua
function vim.lpeg.type(value: boolean|string|integer|function|table...(+1))
  -> "pattern"|nil
```


---

# vim.lpeg.version


```lua
function vim.lpeg.version()
  -> string
```


---

# vim.lsp


```lua
table
```


---

# vim.lsp.CTBufferState


## last_flush


```lua
number|nil
```

uv.hrtime of the last flush/didChange-notification

## lines


```lua
string[]
```

snapshot of buffer lines from last didChange

## lines_tmp


```lua
string[]
```

## name


```lua
string
```

name of the buffer

## needs_flush


```lua
boolean
```

true if buffer updates haven't been sent to clients/servers yet

## pending_changes


```lua
table[]
```

List of debounced changes in incremental sync mode

## refs


```lua
integer
```

how many clients are using this group

## timer


```lua
(uv.uv_timer_t)?
```

uv_timer


---

# vim.lsp.CTGroup

 LSP has 3 different sync modes:
   - None (Servers will read the files themselves when needed)
   - Full (Client sends the full buffer content on updates)
   - Incremental (Client sends only the changed parts)

 Changes are tracked per buffer.
 A buffer can have multiple clients attached and each client needs to send the changes
 To minimize the amount of changesets to compute, computation is grouped:

   None: One group for all clients
   Full: One group for all clients
   Incremental: One group per `position_encoding`

 Sending changes can be debounced per buffer. To simplify the implementation the
 smallest debounce interval is used and we don't group clients by different intervals.


## position_encoding


```lua
"utf-16"|"utf-32"|"utf-8"
```

## sync_kind


```lua
integer
```

TextDocumentSyncKind, considers config.flags.allow_incremental_sync


---

# vim.lsp.CTGroupState


## buffers


```lua
table<integer, vim.lsp.CTBufferState>
```

## clients


```lua
table<integer, vim.lsp.Client>
```

clients using this state. {client_id, client}

## debounce


```lua
integer
```

debounce duration in ms


---

# vim.lsp.Client

## __index


```lua
vim.lsp.Client
```

## _add_workspace_folder


```lua
(method) vim.lsp.Client:_add_workspace_folder(dir?: string)
```

 Add a directory to the workspace folders.

## _before_init_cb


```lua
fun(params: lsp.InitializeParams, config: vim.lsp.ClientConfig)?
```

## _get_language_id


```lua
(method) vim.lsp.Client:_get_language_id(bufnr: any)
  -> string
```

## _get_registration


```lua
(method) vim.lsp.Client:_get_registration(method: string, bufnr?: integer)
  -> (lsp.Registration)?
```

## _get_registration_options


```lua
(method) vim.lsp.Client:_get_registration_options(method: string, bufnr?: integer)
  -> (boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1))?
```

 Get options for a method that is registered dynamically.

## _graceful_shutdown_failed


```lua
true?
```



 Track this so that we can escalate automatically if we've already tried a
 graceful shutdown

## _log_prefix


```lua
string
```

## _notification


```lua
(method) vim.lsp.Client:_notification(method: string, params: table)
```

 Handles a notification sent by an LSP server by invoking the
 corresponding handler.

@*param* `method` — LSP method name

@*param* `params` — The parameters for that method.

## _on_attach_cbs


```lua
fun(client: vim.lsp.Client, bufnr: integer)[]
```

## _on_error


```lua
(method) vim.lsp.Client:_on_error(code: integer, err: any)
```

 Invoked when the client operation throws an error.

@*param* `code` — Error code

@*param* `err` — Other arguments may be passed depending on the error kind

 `vim.lsp.rpc.client_errors[code]` to get a human-friendly name.
See: ~vim.lsp.rpc.client_errors~ for possible errors. Use

## _on_error_cb


```lua
fun(code: integer, err: string)?
```

## _on_exit


```lua
(method) vim.lsp.Client:_on_exit(code: integer, signal: integer)
```

 Invoked on client exit.

@*param* `code` — ) exit code of the process

@*param* `signal` — the signal used to terminate (if any)

## _on_exit_cbs


```lua
fun(code: integer, signal: integer, client_id: integer)[]
```

## _on_init_cbs


```lua
fun(client: vim.lsp.Client, init_result: lsp.InitializeResult)[]
```

## _process_request


```lua
(method) vim.lsp.Client:_process_request(id: integer, req_type: 'cancel'|'complete'|'pending', bufnr?: integer, method?: string)
```

@*param* `bufnr` — (only required for req_type='pending')

@*param* `method` — (only required for req_type='pending')

```lua
req_type:
    | 'pending'
    | 'complete'
    | 'cancel'
```

## _register


```lua
(method) vim.lsp.Client:_register(registrations: lsp.Registration[])
```

## _register_dynamic


```lua
(method) vim.lsp.Client:_register_dynamic(registrations: lsp.Registration[])
```

## _remove_workspace_folder


```lua
(method) vim.lsp.Client:_remove_workspace_folder(dir?: string)
```

 Remove a directory to the workspace folders.

## _resolve_handler


```lua
(method) vim.lsp.Client:_resolve_handler(method: string)
  -> handler: (fun(err?: lsp.ResponseError, result: any, context: lsp.HandlerContext, config?: table):...unknown)?
```

 Returns the handler associated with an LSP method.
 Returns the default handler if the user hasn't set a custom one.

@*param* `method` — LSP method name

@*return* `handler` — for the given method, if defined, or the default from |vim.lsp.handlers|

## _run_callbacks


```lua
(method) vim.lsp.Client:_run_callbacks(cbs: function[], error_id: integer, ...any)
```

## _server_request


```lua
(method) vim.lsp.Client:_server_request(method: string, params: table)
  -> result: any
  2. error: lsp.ResponseError
```

 Handles a request from an LSP server by invoking the corresponding handler.

@*param* `method` — LSP method name

@*param* `params` — The parameters for that method

@*return* `result`

@*return* `error` — code and message set in case an exception happens during the request.

## _supports_registration


```lua
(method) vim.lsp.Client:_supports_registration(method: string)
  -> false
```

 Get options for a method that is registered dynamically.

## _text_document_did_open_handler


```lua
(method) vim.lsp.Client:_text_document_did_open_handler(bufnr: integer)
```

 Default handler for the 'textDocument/didOpen' LSP notification.

@*param* `bufnr` — Number of the buffer, or 0 for current

## _trace


```lua
'messages'|'off'|'verbose'
```


 The initial trace setting. If omitted trace is disabled ("off").
 trace = "off" | "messages" | "verbose";

## _unregister


```lua
(method) vim.lsp.Client:_unregister(unregistrations: lsp.Unregistration[])
```

## _unregister_dynamic


```lua
(method) vim.lsp.Client:_unregister_dynamic(unregistrations: lsp.Unregistration[])
```

## attached_buffers


```lua
table<integer, true>
```


## cancel_request


```lua
(method) vim.lsp.Client:cancel_request(id: integer)
  -> status: boolean
```

 Cancels a request with a given request id.

@*param* `id` — id of request to cancel

@*return* `status` — indicating if the notification was successful.

 @see |Client:notify()|

## capabilities


```lua
lsp.ClientCapabilities
```


 Capabilities provided by the client (editor or tool), at startup.

## commands


```lua
table<string, fun(command: lsp.Command, ctx: table)>
```


 Client commands. See [vim.lsp.ClientConfig].

## config


```lua
vim.lsp.ClientConfig
```


 Copy of the config passed to |vim.lsp.start()|.

## create


```lua
function vim.lsp.Client.create(config: vim.lsp.ClientConfig)
  -> (vim.lsp.Client)?
```

 @nodoc

## dynamic_capabilities


```lua
lsp.DynamicCapabilities
```


 Capabilities provided at runtime (after startup).

## exec_cmd


```lua
(method) vim.lsp.Client:exec_cmd(command: lsp.Command, context?: { bufnr: integer }, handler?: fun(err?: lsp.ResponseError, result: any, context: lsp.HandlerContext, config?: table):...unknown)
```

 Execute a lsp command, either via client command function (if available)
 or via workspace/executeCommand (if supported by the server)

@*param* `handler` — only called if a server command

## flags


```lua
vim.lsp.Client.Flags
```


 A table with flags for the client. The current (experimental) flags are:

## get_language_id


```lua
fun(bufnr: integer, filetype: string):string
```


 See [vim.lsp.ClientConfig].

## handlers


```lua
table<string, fun(err?: lsp.ResponseError, result: any, context: lsp.HandlerContext, config?: table):...unknown>
```


 See [vim.lsp.ClientConfig].

## id


```lua
integer
```


 The id allocated to the client.

## initialize


```lua
(method) vim.lsp.Client:initialize()
```

 @nodoc

## initialized


```lua
true?
```


## is_stopped


```lua
(method) vim.lsp.Client:is_stopped()
  -> boolean
```

 Checks whether a client is stopped.

@*return* — true if client is stopped or in the process of being

 stopped; false otherwise

## name


```lua
string
```


 See [vim.lsp.ClientConfig].

## notify


```lua
(method) vim.lsp.Client:notify(method: string, params?: table)
  -> status: boolean
```

 Sends a notification to an LSP server.

@*param* `method` — LSP method name.

@*param* `params` — LSP request params.

@*return* `status` — indicating if the notification was successful.

                        If it is false, then the client has shutdown.

## offset_encoding


```lua
string
```


 See [vim.lsp.ClientConfig].

## on_attach


```lua
(method) vim.lsp.Client:on_attach(bufnr: integer)
```

 Runs the on_attach function from the client's config if it was defined.
 Useful for buffer-local setup.

@*param* `bufnr` — Buffer number

## progress


```lua
vim.lsp.Client.Progress
```


 A ring buffer (|vim.ringbuf()|) containing progress messages
 sent by the server.

## registrations


```lua
table<string, lsp.Registration[]>
```


## request


```lua
(method) vim.lsp.Client:request(method: string, params?: table, handler?: fun(err?: lsp.ResponseError, result: any, context: lsp.HandlerContext, config?: table):...unknown, bufnr?: integer)
  -> status: boolean
  2. request_id: integer?
```

 Sends a request to the server.

 This is a thin wrapper around {client.rpc.request} with some additional
 checks for capabilities and handler availability.

@*param* `method` — LSP method name.

@*param* `params` — LSP request params.

@*param* `handler` — Response |lsp-handler| for this method.

@*param* `bufnr` — (default: 0) Buffer handle, or 0 for current.

@*return* `status` — indicates whether the request was successful.

     If it is `false`, then it will always be `false` (the client has shutdown).

@*return* `request_id` — Can be used with |Client:cancel_request()|.

                             `nil` is request failed.
 to cancel the-request.
 @see |vim.lsp.buf_request_all()|

## request_sync


```lua
(method) vim.lsp.Client:request_sync(method: string, params: table, timeout_ms?: integer, bufnr?: integer)
  -> { err: (lsp.ResponseError)?, result: any }?
  2. err: string?
```

 Sends a request to the server and synchronously waits for the response.

 This is a wrapper around |Client:request()|

@*param* `method` — LSP method name.

@*param* `params` — LSP request params.

@*param* `timeout_ms` — Maximum time in milliseconds to wait for

                                a result. Defaults to 1000

@*param* `bufnr` — (default: 0) Buffer handle, or 0 for current.

@*return* — `result` and `err` from the |lsp-handler|.

                 `nil` is the request was unsuccessful

@*return* `err` — On timeout, cancel or error, where `err` is a

                 string describing the failure reason.
 @see |vim.lsp.buf_request_sync()|

## requests


```lua
table<integer, { type: string, bufnr: integer, method: string }?>
```


 The current pending requests in flight to the server. Entries are key-value
 pairs with the key being the request id while the value is a table with
 `type`, `bufnr`, and `method` key-value pairs. `type` is either "pending"
 for an active request, or "cancel" for a cancel request. It will be
 "complete" ephemerally while executing |LspRequest| autocmds when replies
 are received from the server.

## root_dir


```lua
string?
```


 See [vim.lsp.ClientConfig].

## rpc


```lua
vim.lsp.rpc.PublicClient
```


 RPC client object, for low level interaction with the client.
 See |vim.lsp.rpc.start()|.

## server_capabilities


```lua
(lsp.ServerCapabilities)?
```


 Response from the server sent on `initialize` describing the server's capabilities.

## server_info


```lua
(lsp.ServerInfo)?
```


 Response from the server sent on `initialize` describing server information (e.g. version).

## settings


```lua
table<string, lsp.LSPAny>
```


 See [vim.lsp.ClientConfig].

## stop


```lua
(method) vim.lsp.Client:stop(force?: boolean)
```

 Stops a client, optionally with force.

 By default, it will just request the server to shutdown without force. If
 you request to stop a client which has previously been requested to
 shutdown, it will automatically escalate and force shutdown.

## supports_method


```lua
(method) vim.lsp.Client:supports_method(method: string, bufnr?: integer)
  -> boolean
```

 Checks if a client supports a given method.
 Always returns true for unknown off-spec methods.

 Note: Some language server capabilities can be file specific.

## workspace_folders


```lua
lsp.WorkspaceFolder[]?
```


 See [vim.lsp.ClientConfig].

## write_error


```lua
(method) vim.lsp.Client:write_error(code: integer, err: any)
```

 Logs the given error to the LSP log and to the error buffer.

@*param* `code` — Error code

@*param* `err` — Error arguments


---

# vim.lsp.Client.Flags

## allow_incremental_sync


```lua
boolean?
```

 @inlinedoc

 Allow using incremental sync for buffer edits
 (default: `true`)

## debounce_text_changes


```lua
integer
```


 Debounce `didChange` notifications to the server by the given number in milliseconds.
 No debounce occurs if `nil`.
 (default: `150`)

## exit_timeout


```lua
integer|false
```


 Milliseconds to wait for server to exit cleanly after sending the
 "shutdown" request before sending kill -15. If set to false, nvim exits
 immediately after sending the "shutdown" request to the server.
 (default: `false`)


---

# vim.lsp.Client.Progress

## _idx_read


```lua
integer
```

## _idx_write


```lua
integer
```

## _items


```lua
table[]
```

## _size


```lua
integer
```

## clear


```lua
function vim.Ringbuf.clear(self: vim.Ringbuf)
```

 Clear all items

## peek


```lua
function vim.Ringbuf.peek(self: vim.Ringbuf)
  -> <T>?
```

 Returns the first unread item without removing it

## pending


```lua
table<string|integer, boolean|string|number|boolean|string|number|table<string, lsp.LSPAny>|table<string, lsp.LSPAny>[]...(+1)[]|table<string, lsp.LSPAny>...(+1)>
```

## pop


```lua
function vim.Ringbuf.pop(self: vim.Ringbuf)
  -> <T>?
```

 Removes and returns the first unread item

## push


```lua
function vim.Ringbuf.push(self: vim.Ringbuf, item: <T>)
```

 Adds an item, overriding the oldest item if the buffer is full.


---

# vim.lsp.ClientConfig

## before_init


```lua
fun(params: lsp.InitializeParams, config: vim.lsp.ClientConfig)?
```


 Callback invoked before the LSP "initialize" phase, where `params` contains the parameters
 being sent to the server and `config` is the config that was passed to |vim.lsp.start()|.
 You can use this to modify parameters before they are sent.

## capabilities


```lua
(lsp.ClientCapabilities)?
```

 Map overriding the default capabilities defined by |vim.lsp.protocol.make_client_capabilities()|,
 passed to the language server on initialization. Hint: use make_client_capabilities() and modify
 its result.
 - Note: To send an empty dictionary use |vim.empty_dict()|, else it will be encoded as an
   array.

## cmd


```lua
fun(dispatchers: vim.lsp.rpc.Dispatchers, config: vim.lsp.ClientConfig):vim.lsp.rpc.PublicClient|string[]
```


 Command `string[]` that launches the language server (treated as in |jobstart()|, must be
 absolute or on `$PATH`, shell constructs like "~" are not expanded), or function that creates an
 RPC client. Function receives a `dispatchers` table and the resolved `config`, and must return
 a table with member functions `request`, `notify`, `is_closing` and `terminate`.
 See |vim.lsp.rpc.request()|, |vim.lsp.rpc.notify()|.
 For TCP there is a builtin RPC client factory: |vim.lsp.rpc.connect()|

## cmd_cwd


```lua
string?
```


 Directory to launch the `cmd` process. Not related to `root_dir`.
 (default: cwd)

## cmd_env


```lua
table?
```


 Environment variables passed to the LSP process on spawn. Non-string values are coerced to
 string.
 Example:
 ```lua
 { PORT = 8080; HOST = '0.0.0.0'; }
 ```

## commands


```lua
table<string, fun(command: lsp.Command, ctx: table)>?
```


 Client commands. Map of command names to user-defined functions. Commands passed to `start()`
 take precedence over the global command registry. Each key must be a unique command name, and
 the value is a function which is called if any LSP action (code action, code lenses, …) triggers
 the command.

## detached


```lua
boolean?
```


 Daemonize the server process so that it runs in a separate process group from Nvim.
 Nvim will shutdown the process on exit, but if Nvim fails to exit cleanly this could leave
 behind orphaned server processes.
 (default: `true`)

## flags


```lua
(vim.lsp.Client.Flags)?
```


 A table with flags for the client. The current (experimental) flags are:

## get_language_id


```lua
(fun(bufnr: integer, filetype: string):string)?
```


 Language ID as string. Defaults to the buffer filetype.

## handlers


```lua
table<string, function>?
```


 Map of LSP method names to |lsp-handler|s.

## init_options


```lua
table<string, lsp.LSPAny>?
```


 Values to pass in the initialization request as `initializationOptions`. See `initialize` in
 the LSP spec.

## name


```lua
string?
```


 Name in logs and user messages.
 (default: client-id)

## offset_encoding


```lua
('utf-16'|'utf-32'|'utf-8')?
```


 Called "position encoding" in LSP spec. The encoding that the LSP server expects, used for
 communication. Not validated. Can be modified in `on_init` before text is sent to the server.

## on_attach


```lua
elem_or_list<fun(client: vim.lsp.Client, bufnr: integer)>?
```


 Callback invoked when client attaches to a buffer.

## on_error


```lua
fun(code: integer, err: string)?
```


 Callback invoked when the client operation throws an error. `code` is a number describing the error.
 Other arguments may be passed depending on the error kind.  See `vim.lsp.rpc.client_errors`
 for possible errors. Use `vim.lsp.rpc.client_errors[code]` to get human-friendly name.

## on_exit


```lua
elem_or_list<fun(code: integer, signal: integer, client_id: integer)>?
```

 Callback invoked on client exit.
   - code: exit code of the process
   - signal: number describing the signal used to terminate (if any)
   - client_id: client handle

## on_init


```lua
elem_or_list<fun(client: vim.lsp.Client, init_result: lsp.InitializeResult)>?
```


 Callback invoked after LSP "initialize", where `result` is a table of `capabilities` and
 anything else the server may send. For example, clangd sends `init_result.offsetEncoding` if
 `capabilities.offsetEncoding` was sent to it. You can only modify the `client.offset_encoding`
 here before any notifications are sent.

## root_dir


```lua
string?
```


 Directory where the LSP server will base its workspaceFolders, rootUri, and rootPath on initialization.

## settings


```lua
table<string, lsp.LSPAny>?
```


 Map of language server-specific settings, decided by the client. Sent to the LS if requested via
 `workspace/configuration`. Keys are case-sensitive.

## trace


```lua
('messages'|'off'|'verbose')?
```


 Passed directly to the language server in the initialize request. Invalid/empty values will
 (default: "off")

## workspace_folders


```lua
lsp.WorkspaceFolder[]?
```


 List of workspace folders passed to the language server. For backwards compatibility rootUri and
 rootPath are derived from the first workspace folder in this list. Can be `null` if the client
 supports workspace folders but none are configured. See `workspaceFolders` in LSP spec.

## workspace_required


```lua
boolean?
```


 Server requires a workspace (no "single file" support). Note: Without
 a workspace, cross-file features (navigation, hover) may or may not work depending on the
 language server, even if the server doesn't require a workspace.
 (default: `false`)


---

# vim.lsp.CodeActionResultEntry

## ctx


```lua
lsp.HandlerContext
```

## error


```lua
(lsp.ResponseError)?
```

## result


```lua
(lsp.CodeAction|lsp.Command)[]?
```


---

# vim.lsp.CompletionResult


---

# vim.lsp.Config

## before_init


```lua
fun(params: lsp.InitializeParams, config: vim.lsp.ClientConfig)?
```


 Callback invoked before the LSP "initialize" phase, where `params` contains the parameters
 being sent to the server and `config` is the config that was passed to |vim.lsp.start()|.
 You can use this to modify parameters before they are sent.

## capabilities


```lua
(lsp.ClientCapabilities)?
```

 Map overriding the default capabilities defined by |vim.lsp.protocol.make_client_capabilities()|,
 passed to the language server on initialization. Hint: use make_client_capabilities() and modify
 its result.
 - Note: To send an empty dictionary use |vim.empty_dict()|, else it will be encoded as an
   array.

## cmd


```lua
(fun(dispatchers: vim.lsp.rpc.Dispatchers, config: vim.lsp.ClientConfig):vim.lsp.rpc.PublicClient|string[])?
```


 See `cmd` in [vim.lsp.ClientConfig].
 See also `reuse_client` to dynamically decide (per-buffer) when `cmd` should be re-invoked.

## cmd_cwd


```lua
string?
```


 Directory to launch the `cmd` process. Not related to `root_dir`.
 (default: cwd)

## cmd_env


```lua
table?
```


 Environment variables passed to the LSP process on spawn. Non-string values are coerced to
 string.
 Example:
 ```lua
 { PORT = 8080; HOST = '0.0.0.0'; }
 ```

## commands


```lua
table<string, fun(command: lsp.Command, ctx: table)>?
```


 Client commands. Map of command names to user-defined functions. Commands passed to `start()`
 take precedence over the global command registry. Each key must be a unique command name, and
 the value is a function which is called if any LSP action (code action, code lenses, …) triggers
 the command.

## detached


```lua
boolean?
```


 Daemonize the server process so that it runs in a separate process group from Nvim.
 Nvim will shutdown the process on exit, but if Nvim fails to exit cleanly this could leave
 behind orphaned server processes.
 (default: `true`)

## filetypes


```lua
string[]?
```


 Filetypes the client will attach to, if activated by `vim.lsp.enable()`. If not provided, the
 client will attach to all filetypes.

## flags


```lua
(vim.lsp.Client.Flags)?
```


 A table with flags for the client. The current (experimental) flags are:

## get_language_id


```lua
(fun(bufnr: integer, filetype: string):string)?
```


 Language ID as string. Defaults to the buffer filetype.

## handlers


```lua
table<string, function>?
```


 Map of LSP method names to |lsp-handler|s.

## init_options


```lua
table<string, lsp.LSPAny>?
```


 Values to pass in the initialization request as `initializationOptions`. See `initialize` in
 the LSP spec.

## name


```lua
string?
```


 Name in logs and user messages.
 (default: client-id)

## offset_encoding


```lua
('utf-16'|'utf-32'|'utf-8')?
```


 Called "position encoding" in LSP spec. The encoding that the LSP server expects, used for
 communication. Not validated. Can be modified in `on_init` before text is sent to the server.

## on_attach


```lua
elem_or_list<fun(client: vim.lsp.Client, bufnr: integer)>?
```


 Callback invoked when client attaches to a buffer.

## on_error


```lua
fun(code: integer, err: string)?
```


 Callback invoked when the client operation throws an error. `code` is a number describing the error.
 Other arguments may be passed depending on the error kind.  See `vim.lsp.rpc.client_errors`
 for possible errors. Use `vim.lsp.rpc.client_errors[code]` to get human-friendly name.

## on_exit


```lua
elem_or_list<fun(code: integer, signal: integer, client_id: integer)>?
```

 Callback invoked on client exit.
   - code: exit code of the process
   - signal: number describing the signal used to terminate (if any)
   - client_id: client handle

## on_init


```lua
elem_or_list<fun(client: vim.lsp.Client, init_result: lsp.InitializeResult)>?
```


 Callback invoked after LSP "initialize", where `result` is a table of `capabilities` and
 anything else the server may send. For example, clangd sends `init_result.offsetEncoding` if
 `capabilities.offsetEncoding` was sent to it. You can only modify the `client.offset_encoding`
 here before any notifications are sent.

## reuse_client


```lua
(fun(client: vim.lsp.Client, config: vim.lsp.ClientConfig):boolean)?
```


 Predicate which decides if a client should be re-used. Used on all running clients. The default
 implementation re-uses a client if name and root_dir matches.

## root_dir


```lua
(string|fun(bufnr: integer, on_dir: fun(root_dir?: string)))?
```


 [lsp-root_dir()](file:///usr/local/share/nvim/runtime/lua/vim)
 Decides the workspace root: the directory where the LSP server will base its workspaceFolders,
 rootUri, and rootPath on initialization. The function form must call the `on_dir` callback to
 provide the root dir, or LSP will not be activated for the buffer. Thus a `root_dir()` function
 can dynamically decide per-buffer whether to activate (or skip) LSP.
 See example at |vim.lsp.enable()|.

## root_markers


```lua
(string|string[])[]?
```

 [lsp-root_markers](file:///usr/local/share/nvim/runtime/lua/vim)
 Filename(s) (".git/", "package.json", …) used to decide the workspace root. Unused if `root_dir`
 is defined. The list order decides priority. To indicate "equal priority", specify names in
 a nested list `{ { 'a.txt', 'b.lua' }, ... }`.

 For each item, Nvim will search upwards (from the buffer file) for that marker, or list of
 markers; search stops at the first directory containing that marker, and the directory is used
 as the root dir (workspace folder).

 Example: Find the first ancestor directory containing file or directory "stylua.toml"; if not
 found, find the first ancestor containing ".git":
 ```lua
   root_markers = { 'stylua.toml', '.git' }
 ```

 Example: Find the first ancestor directory containing EITHER "stylua.toml" or ".luarc.json"; if
 not found, find the first ancestor containing ".git":
 ```lua
   root_markers = { { 'stylua.toml', '.luarc.json' }, '.git' }
 ```


## settings


```lua
table<string, lsp.LSPAny>?
```


 Map of language server-specific settings, decided by the client. Sent to the LS if requested via
 `workspace/configuration`. Keys are case-sensitive.

## trace


```lua
('messages'|'off'|'verbose')?
```


 Passed directly to the language server in the initialize request. Invalid/empty values will
 (default: "off")

## workspace_folders


```lua
lsp.WorkspaceFolder[]?
```


 List of workspace folders passed to the language server. For backwards compatibility rootUri and
 rootPath are derived from the first workspace folder in this list. Can be `null` if the client
 supports workspace folders but none are configured. See `workspaceFolders` in LSP spec.

## workspace_required


```lua
boolean?
```


 Server requires a workspace (no "single file" support). Note: Without
 a workspace, cross-file features (navigation, hover) may or may not work depending on the
 language server, even if the server doesn't require a workspace.
 (default: `false`)


---

# vim.lsp.ListOpts

## loclist


```lua
boolean?
```

 Whether to use the |location-list| or the |quickfix| list in the default handler.
 ```lua
 vim.lsp.buf.definition({ loclist = true })
 vim.lsp.buf.references(nil, { loclist = false })
 ```

## on_list


```lua
fun(t: vim.lsp.LocationOpts.OnList)?
```


 list-handler replacing the default handler.
 Called for any non-empty result.
 This table can be used with |setqflist()| or |setloclist()|. E.g.:
 ```lua
 local function on_list(options)
   vim.fn.setqflist({}, ' ', options)
   vim.cmd.cfirst()
 end

 vim.lsp.buf.definition({ on_list = on_list })
 vim.lsp.buf.references(nil, { on_list = on_list })
 ```


---

# vim.lsp.LocationOpts

## loclist


```lua
boolean?
```

 Whether to use the |location-list| or the |quickfix| list in the default handler.
 ```lua
 vim.lsp.buf.definition({ loclist = true })
 vim.lsp.buf.references(nil, { loclist = false })
 ```

## on_list


```lua
fun(t: vim.lsp.LocationOpts.OnList)?
```


 list-handler replacing the default handler.
 Called for any non-empty result.
 This table can be used with |setqflist()| or |setloclist()|. E.g.:
 ```lua
 local function on_list(options)
   vim.fn.setqflist({}, ' ', options)
   vim.cmd.cfirst()
 end

 vim.lsp.buf.definition({ on_list = on_list })
 vim.lsp.buf.references(nil, { on_list = on_list })
 ```

## reuse_win


```lua
boolean?
```


 Jump to existing window if buffer is already open.


---

# vim.lsp.LocationOpts.OnList

## context


```lua
{ bufnr: integer, method: string }?
```

Subset of `ctx` from |lsp-handler|.

## items


```lua
table[]
```

Structured like |setqflist-what|

## title


```lua
string?
```

Title for the list.


---

# vim.lsp.buf.code_action.Opts

## apply


```lua
boolean?
```


 When set to `true`, and there is just one remaining action
 (after filtering), the action is applied without user query.

## context


```lua
(lsp.CodeActionContext)?
```

 @inlinedoc

 Corresponds to `CodeActionContext` of the LSP specification:
   - {diagnostics}? (`table`) LSP `Diagnostic[]`. Inferred from the current
     position if not provided.
   - {only}? (`table`) List of LSP `CodeActionKind`s used to filter the code actions.
     Most language servers support values like `refactor`
     or `quickfix`.
   - {triggerKind}? (`integer`) The reason why code actions were requested.

## filter


```lua
(fun(x: lsp.CodeAction|lsp.Command):boolean)?
```


 Predicate taking an `CodeAction` and returning a boolean.

## range


```lua
{ start: integer[], end: integer[] }?
```


 Range for which code actions should be requested.
 If in visual mode this defaults to the active selection.
 Table must contain `start` and `end` keys with {row,col} tuples
 using mark-like indexing. See |api-indexing|


---

# vim.lsp.buf.format.Opts

## async


```lua
boolean?
```


 If true the method won't block.
 Editing the buffer while formatting asynchronous can lead to unexpected
 changes.
 (Default: false)

## bufnr


```lua
integer?
```


 Restrict formatting to the clients attached to the given buffer.
 (default: current buffer)

## filter


```lua
fun(client: vim.lsp.Client):boolean??
```


 Predicate used to filter clients. Receives a client as argument and must
 return a boolean. Clients matching the predicate are included. Example:
 ```lua
 -- Never request typescript-language-server for formatting
 vim.lsp.buf.format {
   filter = function(client) return client.name ~= "ts_ls" end
 }
 ```

## formatting_options


```lua
table?
```

 @inlinedoc

 Can be used to specify FormattingOptions. Some unspecified options will be
 automatically derived from the current Nvim options.
 See https://microsoft.github.io/language-server-protocol/specification/#formattingOptions

## id


```lua
integer?
```


 Restrict formatting to the client with ID (client.id) matching this field.

## name


```lua
string?
```


 Restrict formatting to the client with name (client.name) matching this field.

## range


```lua
({ start: [integer, integer], end: [integer, integer] }|{ start: [integer, integer], end: [integer, integer] }[])?
```


 Range to format.
 Table must contain `start` and `end` keys with {row,col} tuples using
 (1,0) indexing.
 Can also be a list of tables that contain `start` and `end` keys as described above,
 in which case `textDocument/rangesFormatting` support is required.
 (Default: current selection in visual mode, `nil` in other modes,
 formatting the full buffer)

## timeout_ms


```lua
integer?
```


 Time in milliseconds to block for formatting requests. No effect if async=true.
 (default: `1000`)


---

# vim.lsp.buf.hover.Opts

## _update_win


```lua
integer?
```


## anchor_bias


```lua
('above'|'auto'|'below')?
```


 Adjusts placement relative to cursor.
 - "auto": place window based on which side of the cursor has more lines
 - "above": place the window above the cursor unless there are not enough lines
   to display the full window height.
 - "below": place the window below the cursor unless there are not enough lines
   to display the full window height.
 (default: `'auto'`)

## border


```lua
(string|(string|[string, string])[])?
```

override `border`

## close_events


```lua
table?
```


 List of events that closes the floating window

## focus


```lua
boolean?
```


 If `true`, and if {focusable} is also `true`, focus an existing floating
 window with the same {focus_id}
 (default: `true`)

## focus_id


```lua
string?
```


 If a popup with this id is opened, then focus it

## focusable


```lua
boolean?
```


 Make float focusable.
 (default: `true`)

## height


```lua
integer?
```


 Height of floating window

## max_height


```lua
integer?
```


 Maximal height of floating window

## max_width


```lua
integer?
```


 Maximal width of floating window

## offset_x


```lua
integer?
```


 offset to add to `col`

## offset_y


```lua
integer?
```


 offset to add to `row`

## relative


```lua
('cursor'|'editor'|'mouse')?
```


 (default: `'cursor'`)

## silent


```lua
boolean?
```

## title


```lua
(string|[string, string][])?
```

## title_pos


```lua
('center'|'left'|'right')?
```

## width


```lua
integer?
```


 Width of floating window

## wrap


```lua
boolean?
```


 Wrap long lines
 (default: `true`)

## wrap_at


```lua
integer?
```


 Character to wrap at for computing height when wrap is enabled

## zindex


```lua
integer?
```

override `zindex`, defaults to 50


---

# vim.lsp.buf.rename.Opts

## bufnr


```lua
integer?
```


 (default: current buffer)

## filter


```lua
fun(client: vim.lsp.Client):boolean??
```

 @inlinedoc

 Predicate used to filter clients. Receives a client as argument and
 must return a boolean. Clients matching the predicate are included.

## name


```lua
string?
```


 Restrict clients used for rename to ones where client.name matches
 this field.


---

# vim.lsp.buf.signature_help.Opts

## _update_win


```lua
integer?
```


## anchor_bias


```lua
('above'|'auto'|'below')?
```


 Adjusts placement relative to cursor.
 - "auto": place window based on which side of the cursor has more lines
 - "above": place the window above the cursor unless there are not enough lines
   to display the full window height.
 - "below": place the window below the cursor unless there are not enough lines
   to display the full window height.
 (default: `'auto'`)

## border


```lua
(string|(string|[string, string])[])?
```

override `border`

## close_events


```lua
table?
```


 List of events that closes the floating window

## focus


```lua
boolean?
```


 If `true`, and if {focusable} is also `true`, focus an existing floating
 window with the same {focus_id}
 (default: `true`)

## focus_id


```lua
string?
```


 If a popup with this id is opened, then focus it

## focusable


```lua
boolean?
```


 Make float focusable.
 (default: `true`)

## height


```lua
integer?
```


 Height of floating window

## max_height


```lua
integer?
```


 Maximal height of floating window

## max_width


```lua
integer?
```


 Maximal width of floating window

## offset_x


```lua
integer?
```


 offset to add to `col`

## offset_y


```lua
integer?
```


 offset to add to `row`

## relative


```lua
('cursor'|'editor'|'mouse')?
```


 (default: `'cursor'`)

## silent


```lua
boolean?
```

## title


```lua
(string|[string, string][])?
```

## title_pos


```lua
('center'|'left'|'right')?
```

## width


```lua
integer?
```


 Width of floating window

## wrap


```lua
boolean?
```


 Wrap long lines
 (default: `true`)

## wrap_at


```lua
integer?
```


 Character to wrap at for computing height when wrap is enabled

## zindex


```lua
integer?
```

override `zindex`, defaults to 50


---

# vim.lsp.client.before_init_cb


---

# vim.lsp.client.on_attach_cb


---

# vim.lsp.client.on_exit_cb


---

# vim.lsp.client.on_init_cb


---

# vim.lsp.codelens.refresh.Opts

## bufnr


```lua
integer?
```

filter by buffer. All buffers if nil, 0 for current buffer


---

# vim.lsp.completion.BufHandle

 @nodoc

## clients


```lua
table<integer, vim.lsp.Client>
```

## convert


```lua
(fun(item: lsp.CompletionItem):table)?
```

## triggers


```lua
table<string, vim.lsp.Client[]>
```


---

# vim.lsp.completion.BufferOpts

 @inlinedoc

## autotrigger


```lua
boolean?
```

(default: false) When true, completion triggers automatically based on the server's `triggerCharacters`.

## convert


```lua
(fun(item: lsp.CompletionItem):table)?
```

Transforms an LSP CompletionItem to |complete-items|.


---

# vim.lsp.completion.Context

 @nodoc

## cancel_pending


```lua
(method) vim.lsp.completion.Context:cancel_pending()
```

 @nodoc

## reset


```lua
(method) vim.lsp.completion.Context:reset()
```

 @nodoc


---

# vim.lsp.completion.get.Opts

 @inlinedoc

## ctx


```lua
(lsp.CompletionContext)?
```

Completion context. Defaults to a trigger kind of `invoked`.


---

# vim.lsp.config

 @nodoc

## [string]


```lua
vim.lsp.Config
```

## _configs


```lua
table<string, vim.lsp.Config>
```


---

# vim.lsp.folding_range.BufState

## client_ranges


```lua
table<integer, lsp.FoldingRange[]?>
```


 Never use this directly, `renew()` the cached foldinfo
 then use on demand via `row_*` fields.

 Index In the form of client_id -> ranges

## row_kinds


```lua
table<integer, table<"comment"|"imports"|"region", true?>?>
```

>

## row_level


```lua
table<integer, [integer, ("<"|">")?]?>
```


 Index in the form of row -> [foldlevel, mark]

## row_text


```lua
table<integer, string?>
```


 Index in the form of start_row -> collapsed_text

## version


```lua
integer?
```



---

# vim.lsp.formatexpr.Opts

## timeout_ms


```lua
integer
```

 @inlinedoc

 The timeout period for the formatting request.
 (default: 500ms).


---

# vim.lsp.get_clients.Filter

 Key-value pairs used to filter the returned clients.

## _uninitialized


```lua
boolean?
```


 Also return uninitialized clients.

## bufnr


```lua
integer?
```


 Only return clients attached to this buffer

## id


```lua
integer?
```

 @inlinedoc

 Only return clients with the given id

## method


```lua
string?
```


 Only return clients supporting the given method

## name


```lua
string?
```


 Only return clients with the given name


---

# vim.lsp.inlay_hint.bufstate

## applied


```lua
table<integer, integer>
```

Last version of hints applied to this line

## client_hints


```lua
table<integer, table<integer, lsp.InlayHint[]>>?
```

client_id -> (lnum -> hints)

## enabled


```lua
boolean
```

Whether inlay hints are enabled for this scope

## version


```lua
integer?
```


---

# vim.lsp.inlay_hint.enable.Filter

 Optional filters |kwargs|, or `nil` for all.

## bufnr


```lua
integer?
```

 @inlinedoc
 Buffer number, or 0 for current buffer, or nil for all.


---

# vim.lsp.inlay_hint.get.Filter

 Optional filters |kwargs|:

## bufnr


```lua
integer?
```

 @inlinedoc

## range


```lua
(lsp.Range)?
```

A range in a text document expressed as (zero-based) start and end positions.

If you want to specify a range that contains a line including the line ending
character(s) then use an end position denoting the start of the next line.
For example:
```ts
{
    start: { line: 5, character: 23 }
    end : { line 6, character : 0 }
}
```


---

# vim.lsp.inlay_hint.get.ret

## bufnr


```lua
integer
```

 @inlinedoc

## client_id


```lua
integer
```

## inlay_hint


```lua
lsp.InlayHint
```

Inlay hint information.



---

# vim.lsp.inlay_hint.globalstate

## enabled


```lua
boolean
```

Whether inlay hints are enabled for this scope


---

# vim.lsp.protocol

 Protocol for the Microsoft Language Server Protocol (mslsp)

## Methods


```lua
table
```

 LSP method names.
See: ~https~ ://microsoft.github.io/language-server-protocol/specification/#metaModel

## _request_name_to_capability


```lua
table
```

 stylua: ignore start
 Generated by gen_lsp.lua, keep at end of file.
 Maps method names to the required server capability

## make_client_capabilities


```lua
function vim.lsp.protocol.make_client_capabilities()
  -> lsp.ClientCapabilities
```

 Gets a new ClientCapabilities object describing the LSP client
 capabilities.

## resolve_capabilities


```lua
function vim.lsp.protocol.resolve_capabilities(server_capabilities: table)
  -> lsp.ServerCapabilities|nil
```

 Creates a normalized object describing LSP server capabilities.

@*param* `server_capabilities` — Table of capabilities supported by the server

@*return* — : Normalized table of capabilities


---

# vim.lsp.protocol.Method


---

# vim.lsp.protocol.Method.ClientToServer

 Generated by gen_lsp.lua, keep at end of file.


---

# vim.lsp.protocol.Method.ServerToClient


---

# vim.lsp.protocol.Methods

 Generated by gen_lsp.lua, keep at end of file.


---

# vim.lsp.protocol.constants


---

# vim.lsp.rpc.Client

## dispatchers


```lua
vim.lsp.rpc.Dispatchers
```

 Dispatchers for LSP message types.

## encode_and_send


```lua
(method) vim.lsp.rpc.Client:encode_and_send(payload: any)
  -> boolean
```

## handle_body


```lua
(method) vim.lsp.rpc.Client:handle_body(body: string)
```

## message_callbacks


```lua
table<integer, function>
```

dict of message_id to callback

## message_index


```lua
integer
```

## notify


```lua
(method) vim.lsp.rpc.Client:notify(method: string, params: any)
  -> boolean
```

 Sends a notification to the LSP server.

@*param* `method` — The invoked LSP method

@*param* `params` — Parameters for the invoked LSP method

@*return* — `true` if notification could be sent, `false` if not

## notify_reply_callbacks


```lua
table<integer, function>
```

dict of message_id to callback

## on_error


```lua
(method) vim.lsp.rpc.Client:on_error(errkind: integer, ...any)
```

## pcall_handler


```lua
(method) vim.lsp.rpc.Client:pcall_handler(errkind: integer, status: boolean, head: any, ...any)
  -> status: boolean
  2. head: any
  3. ...any
```

## request


```lua
(method) vim.lsp.rpc.Client:request(method: string, params?: table, callback: fun(err?: lsp.ResponseError, result: any), notify_reply_callback?: fun(message_id: integer))
  -> success: boolean
  2. message_id: integer?
```

 Sends a request to the LSP server and runs {callback} upon response. |vim.lsp.rpc.request()|

@*param* `method` — The invoked LSP method

@*param* `params` — Parameters for the invoked LSP method

@*param* `callback` — Callback to invoke

@*param* `notify_reply_callback` — Callback to invoke as soon as a request is no longer pending

@*return* `success` — `true` if request could be sent, `false` if not

@*return* `message_id` — if request could be sent, `nil` if not

## send_response


```lua
(method) vim.lsp.rpc.Client:send_response(request_id: any, err: any, result: any)
  -> boolean
```

 sends an error object to the remote LSP process.

## transport


```lua
vim.lsp.rpc.Transport
```

## try_call


```lua
(method) vim.lsp.rpc.Client:try_call(errkind: integer, fn: function, ...any)
  -> status: boolean
  2. head: any
  3. ...any
```


---

# vim.lsp.rpc.Dispatchers

 Dispatchers for LSP message types.

## notification


```lua
fun(method: string, params: table)
```

 @inlinedoc

## on_error


```lua
fun(code: integer, err: any)
```

## on_exit


```lua
fun(code: integer, signal: integer)
```

## server_request


```lua
fun(method: string, params: table):any, (lsp.ResponseError)?
```


---

# vim.lsp.rpc.ExtraSpawnParams

 Additional context for the LSP server process.

## cwd


```lua
string?
```

Working directory for the LSP server process

## detached


```lua
boolean?
```

Detach the LSP server process from the current process

## env


```lua
table<string, string>?
```

Additional environment variables for LSP server process. See |vim.system()|


---

# vim.lsp.rpc.Headers

## content_length


```lua
integer
```


---

# vim.lsp.rpc.PublicClient

 Client RPC object

## is_closing


```lua
fun():boolean
```


 Indicates if the RPC is closing.

## notify


```lua
fun(method: string, params: any):boolean
```


 See [vim.lsp.rpc.notify()]

## request


```lua
fun(method: string, params?: table, callback: fun(err?: lsp.ResponseError, result: any), notify_reply_callback?: fun(message_id: integer)):boolean, integer?
```


 See [vim.lsp.rpc.request()]

## terminate


```lua
fun()
```


 Terminates the RPC client.


---

# vim.lsp.rpc.Transport

## is_closing


```lua
fun(self: vim.lsp.rpc.Transport):boolean
```

## terminate


```lua
fun(self: vim.lsp.rpc.Transport)
```

## write


```lua
fun(self: vim.lsp.rpc.Transport, msg: string)
```


---

# vim.lsp.rpc.Transport.Connect

## closing


```lua
boolean
```

## connected


```lua
boolean
```

 Connect returns a PublicClient synchronously so the caller
 can immediately send messages before the connection is established
 -> Need to buffer them until that happens

## handle


```lua
(uv.uv_pipe_t|uv.uv_tcp_t)?
```

## is_closing


```lua
fun(self: vim.lsp.rpc.Transport):boolean
```

## msgbuf


```lua
vim.Ringbuf
```

## new


```lua
fun():vim.lsp.rpc.Transport.Connect
```

## on_exit


```lua
fun(code: integer, signal: integer)?
```

## terminate


```lua
fun(self: vim.lsp.rpc.Transport)
```

## write


```lua
fun(self: vim.lsp.rpc.Transport, msg: string)
```


---

# vim.lsp.rpc.Transport.Run

## is_closing


```lua
fun(self: vim.lsp.rpc.Transport):boolean
```

## new


```lua
fun():vim.lsp.rpc.Transport.Run
```

## sysobj


```lua
(vim.SystemObj)?
```

## terminate


```lua
fun(self: vim.lsp.rpc.Transport)
```

## write


```lua
fun(self: vim.lsp.rpc.Transport, msg: string)
```


---

# vim.lsp.semantic_tokens.highlight_token.Opts

## priority


```lua
integer?
```

 @inlinedoc

 Priority for the applied extmark.
 (Default: `vim.hl.priorities.semantic_tokens + 3`)


---

# vim.lsp.start.Opts

## _root_markers


```lua
(string|string[])[]?
```


## attach


```lua
boolean?
```


 Whether to attach the client to a buffer (default true).
 If set to `false`, `reuse_client` and `bufnr` will be ignored.

## bufnr


```lua
integer?
```


 Buffer handle to attach to if starting or re-using a client (0 for current).

## reuse_client


```lua
(fun(client: vim.lsp.Client, config: vim.lsp.ClientConfig):boolean)?
```

 @inlinedoc

 Predicate used to decide if a client should be re-used. Used on all
 running clients. The default implementation re-uses a client if it has the
 same name and if the given workspace folders (or root_dir) are all included
 in the client's workspace folders.

## silent


```lua
boolean?
```


 Suppress error reporting if the LSP server fails to start (default false).


---

# vim.lsp.sync.Range

## byte_idx


```lua
integer
```

## char_idx


```lua
integer
```

## line_idx


```lua
integer
```


---

# vim.lsp.util._cancel_requests.Filter

## bufnr


```lua
integer?
```

## clients


```lua
vim.lsp.Client[]?
```

## method


```lua
string?
```

## type


```lua
string?
```


---

# vim.lsp.util._normalize_markdown.Opts

## width


```lua
integer
```

Thematic breaks are expanded to this size. Defaults to 80.


---

# vim.lsp.util._refresh.Opts

## bufnr


```lua
integer?
```

Buffer to refresh (default: 0)

## client_id


```lua
integer?
```

Client ID to refresh (default: all clients)

## only_visible


```lua
boolean?
```

Whether to only refresh for the visible regions of the buffer (default: false)


---

# vim.lsp.util.open_floating_preview.Opts

## _update_win


```lua
integer?
```


## anchor_bias


```lua
('above'|'auto'|'below')?
```


 Adjusts placement relative to cursor.
 - "auto": place window based on which side of the cursor has more lines
 - "above": place the window above the cursor unless there are not enough lines
   to display the full window height.
 - "below": place the window below the cursor unless there are not enough lines
   to display the full window height.
 (default: `'auto'`)

## border


```lua
(string|(string|[string, string])[])?
```

override `border`

## close_events


```lua
table?
```


 List of events that closes the floating window

## focus


```lua
boolean?
```


 If `true`, and if {focusable} is also `true`, focus an existing floating
 window with the same {focus_id}
 (default: `true`)

## focus_id


```lua
string?
```


 If a popup with this id is opened, then focus it

## focusable


```lua
boolean?
```


 Make float focusable.
 (default: `true`)

## height


```lua
integer?
```


 Height of floating window

## max_height


```lua
integer?
```


 Maximal height of floating window

## max_width


```lua
integer?
```


 Maximal width of floating window

## offset_x


```lua
integer?
```


 offset to add to `col`

## offset_y


```lua
integer?
```


 offset to add to `row`

## relative


```lua
('cursor'|'editor'|'mouse')?
```


 (default: `'cursor'`)

## title


```lua
(string|[string, string][])?
```

## title_pos


```lua
('center'|'left'|'right')?
```

## width


```lua
integer?
```


 Width of floating window

## wrap


```lua
boolean?
```


 Wrap long lines
 (default: `true`)

## wrap_at


```lua
integer?
```


 Character to wrap at for computing height when wrap is enabled

## zindex


```lua
integer?
```

override `zindex`, defaults to 50


---

# vim.lsp.util.rename.Opts

## ignoreIfExists


```lua
boolean?
```

## overwrite


```lua
boolean?
```

 @inlinedoc


---

# vim.lsp.util.show_document.Opts

## focus


```lua
boolean?
```


 Whether to focus/jump to location if possible.
 (defaults: true)

## reuse_win


```lua
boolean?
```

 @inlinedoc

 Jump to existing window if buffer is already open.


---

# vim.lua_omnifunc


```lua
function vim.lua_omnifunc(find_start: 0|1, _: any)
  -> integer|any[]
```


---

# vim.mpack.decode


```lua
function vim.mpack.decode(str: string)
  -> any
```


---

# vim.mpack.encode


```lua
function vim.mpack.encode(obj: any)
  -> string
```


---

# vim.notify


```lua
function vim.notify(msg: string, level: integer|nil, opts: table|nil)
```


---

# vim.notify_once


```lua
function vim.notify_once(msg: string, level: integer|nil, opts: table|nil)
  -> true: boolean
```


---

# vim.o


```lua
table
```


---

# vim.o.acd


```lua
boolean
```


---

# vim.o.ai


```lua
boolean
```


---

# vim.o.allowrevins


```lua
boolean
```


---

# vim.o.ambiwidth


```lua
string
```


---

# vim.o.ambw


```lua
'double'|'single'
```


---

# vim.o.ar


```lua
boolean
```


---

# vim.o.arab


```lua
boolean
```


---

# vim.o.arabic


```lua
boolean
```


---

# vim.o.arabicshape


```lua
boolean
```


---

# vim.o.ari


```lua
boolean
```


---

# vim.o.arshape


```lua
boolean
```


---

# vim.o.autochdir


```lua
boolean
```


---

# vim.o.autoindent


```lua
boolean
```


---

# vim.o.autoread


```lua
boolean
```


---

# vim.o.autowrite


```lua
boolean
```


---

# vim.o.autowriteall


```lua
boolean
```


---

# vim.o.aw


```lua
boolean
```


---

# vim.o.awa


```lua
boolean
```


---

# vim.o.background


```lua
string
```


---

# vim.o.backspace


```lua
string
```


---

# vim.o.backup


```lua
boolean
```


---

# vim.o.backupcopy


```lua
string
```


---

# vim.o.backupdir


```lua
string
```


---

# vim.o.backupext


```lua
string
```


---

# vim.o.backupskip


```lua
string
```


---

# vim.o.bdir


```lua
string
```


---

# vim.o.belloff


```lua
string
```


---

# vim.o.bex


```lua
string
```


---

# vim.o.bg


```lua
'dark'|'light'
```


---

# vim.o.bh


```lua
''|'delete'|'hide'|'unload'|'wipe'
```


---

# vim.o.bin


```lua
boolean
```


---

# vim.o.binary


```lua
boolean
```


---

# vim.o.bk


```lua
boolean
```


---

# vim.o.bkc


```lua
string
```


---

# vim.o.bl


```lua
boolean
```


---

# vim.o.bo


```lua
string
```


---

# vim.o.bomb


```lua
boolean
```


---

# vim.o.breakat


```lua
string
```


---

# vim.o.breakindent


```lua
boolean
```


---

# vim.o.breakindentopt


```lua
string
```


---

# vim.o.bri


```lua
boolean
```


---

# vim.o.briopt


```lua
string
```


---

# vim.o.brk


```lua
string
```


---

# vim.o.bs


```lua
string
```


---

# vim.o.bsk


```lua
string
```


---

# vim.o.bt


```lua
''|'acwrite'|'help'|'nofile'|'nowrite'...(+3)
```


---

# vim.o.bufhidden


```lua
string
```


---

# vim.o.buflisted


```lua
boolean
```


---

# vim.o.buftype


```lua
string
```


---

# vim.o.casemap


```lua
string
```


---

# vim.o.cb


```lua
string
```


---

# vim.o.cc


```lua
string
```


---

# vim.o.ccv


```lua
string
```


---

# vim.o.cd


```lua
string
```


---

# vim.o.cdh


```lua
boolean
```


---

# vim.o.cdhome


```lua
boolean
```


---

# vim.o.cdpath


```lua
string
```


---

# vim.o.cedit


```lua
string
```


---

# vim.o.cf


```lua
boolean
```


---

# vim.o.cfu


```lua
string
```


---

# vim.o.ch


```lua
integer
```


---

# vim.o.channel


```lua
integer
```


---

# vim.o.charconvert


```lua
string
```


---

# vim.o.ci


```lua
boolean
```


---

# vim.o.cia


```lua
string
```


---

# vim.o.cin


```lua
boolean
```


---

# vim.o.cindent


```lua
boolean
```


---

# vim.o.cink


```lua
string
```


---

# vim.o.cinkeys


```lua
string
```


---

# vim.o.cino


```lua
string
```


---

# vim.o.cinoptions


```lua
string
```


---

# vim.o.cinscopedecls


```lua
string
```


---

# vim.o.cinsd


```lua
string
```


---

# vim.o.cinw


```lua
string
```


---

# vim.o.cinwords


```lua
string
```


---

# vim.o.clipboard


```lua
string
```


```lua
string
```


---

# vim.o.cmdheight


```lua
integer
```


---

# vim.o.cmdwinheight


```lua
integer
```


---

# vim.o.cmp


```lua
string
```


---

# vim.o.cms


```lua
string
```


---

# vim.o.co


```lua
integer
```


---

# vim.o.cocu


```lua
string
```


---

# vim.o.cole


```lua
integer
```


---

# vim.o.colorcolumn


```lua
string
```


---

# vim.o.columns


```lua
integer
```


---

# vim.o.com


```lua
string
```


---

# vim.o.comments


```lua
string
```


---

# vim.o.commentstring


```lua
string
```


---

# vim.o.complete


```lua
string
```


---

# vim.o.completefunc


```lua
string
```


---

# vim.o.completeitemalign


```lua
string
```


---

# vim.o.completeopt


```lua
string
```


---

# vim.o.completeslash


```lua
string
```


---

# vim.o.concealcursor


```lua
string
```


---

# vim.o.conceallevel


```lua
integer
```


---

# vim.o.confirm


```lua
boolean
```


```lua
boolean
```


---

# vim.o.copyindent


```lua
boolean
```


---

# vim.o.cot


```lua
string
```


---

# vim.o.cpo


```lua
string
```


---

# vim.o.cpoptions


```lua
string
```


---

# vim.o.cpt


```lua
string
```


---

# vim.o.crb


```lua
boolean
```


---

# vim.o.csl


```lua
''|'backslash'|'slash'
```


---

# vim.o.cuc


```lua
boolean
```


---

# vim.o.cul


```lua
boolean
```


---

# vim.o.culopt


```lua
string
```


---

# vim.o.cursorbind


```lua
boolean
```


---

# vim.o.cursorcolumn


```lua
boolean
```


---

# vim.o.cursorline


```lua
boolean
```


```lua
boolean
```


---

# vim.o.cursorlineopt


```lua
string
```


---

# vim.o.cwh


```lua
integer
```


---

# vim.o.debug


```lua
string
```


---

# vim.o.deco


```lua
boolean
```


---

# vim.o.def


```lua
string
```


---

# vim.o.define


```lua
string
```


---

# vim.o.delcombine


```lua
boolean
```


---

# vim.o.dex


```lua
string
```


---

# vim.o.dg


```lua
boolean
```


---

# vim.o.dict


```lua
string
```


---

# vim.o.dictionary


```lua
string
```


---

# vim.o.diff


```lua
boolean
```


---

# vim.o.diffexpr


```lua
string
```


---

# vim.o.diffopt


```lua
string
```


---

# vim.o.digraph


```lua
boolean
```


---

# vim.o.dip


```lua
string
```


---

# vim.o.dir


```lua
string
```


---

# vim.o.directory


```lua
string
```


---

# vim.o.display


```lua
string
```


---

# vim.o.dy


```lua
string
```


---

# vim.o.ea


```lua
boolean
```


---

# vim.o.ead


```lua
'both'|'hor'|'ver'
```


---

# vim.o.eadirection


```lua
string
```


---

# vim.o.eb


```lua
boolean
```


---

# vim.o.ef


```lua
string
```


---

# vim.o.efm


```lua
string
```


---

# vim.o.ei


```lua
string
```


---

# vim.o.eiw


```lua
string
```


---

# vim.o.emo


```lua
boolean
```


---

# vim.o.emoji


```lua
boolean
```


---

# vim.o.enc


```lua
string
```


---

# vim.o.encoding


```lua
string
```


---

# vim.o.endoffile


```lua
boolean
```


---

# vim.o.endofline


```lua
boolean
```


---

# vim.o.eof


```lua
boolean
```


---

# vim.o.eol


```lua
boolean
```


---

# vim.o.ep


```lua
string
```


---

# vim.o.equalalways


```lua
boolean
```


---

# vim.o.equalprg


```lua
string
```


---

# vim.o.errorbells


```lua
boolean
```


---

# vim.o.errorfile


```lua
string
```


---

# vim.o.errorformat


```lua
string
```


---

# vim.o.et


```lua
boolean
```


---

# vim.o.eventignore


```lua
string
```


---

# vim.o.eventignorewin


```lua
string
```


---

# vim.o.ex


```lua
boolean
```


---

# vim.o.expandtab


```lua
boolean
```


---

# vim.o.exrc


```lua
boolean
```


---

# vim.o.fcl


```lua
string
```


---

# vim.o.fcs


```lua
string
```


---

# vim.o.fdc


```lua
'0'|'1'|'2'|'3'|'4'...(+15)
```


---

# vim.o.fde


```lua
string
```


---

# vim.o.fdi


```lua
string
```


---

# vim.o.fdl


```lua
integer
```


---

# vim.o.fdls


```lua
integer
```


---

# vim.o.fdm


```lua
'diff'|'expr'|'indent'|'manual'|'marker'...(+1)
```


---

# vim.o.fdn


```lua
integer
```


---

# vim.o.fdo


```lua
string
```


---

# vim.o.fdt


```lua
string
```


---

# vim.o.fen


```lua
boolean
```


---

# vim.o.fenc


```lua
string
```


---

# vim.o.fencs


```lua
string
```


---

# vim.o.fex


```lua
string
```


---

# vim.o.ff


```lua
'dos'|'mac'|'unix'
```


---

# vim.o.ffs


```lua
string
```


---

# vim.o.ffu


```lua
string
```


---

# vim.o.fic


```lua
boolean
```


---

# vim.o.fileencoding


```lua
string
```


---

# vim.o.fileencodings


```lua
string
```


---

# vim.o.fileformat


```lua
string
```


---

# vim.o.fileformats


```lua
string
```


---

# vim.o.fileignorecase


```lua
boolean
```


---

# vim.o.filetype


```lua
string
```


---

# vim.o.fillchars


```lua
string
```


---

# vim.o.findfunc


```lua
string
```


---

# vim.o.fixendofline


```lua
boolean
```


---

# vim.o.fixeol


```lua
boolean
```


---

# vim.o.flp


```lua
string
```


---

# vim.o.fml


```lua
integer
```


---

# vim.o.fmr


```lua
string
```


---

# vim.o.fo


```lua
string
```


---

# vim.o.foldclose


```lua
string
```


---

# vim.o.foldcolumn


```lua
string
```


---

# vim.o.foldenable


```lua
boolean
```


---

# vim.o.foldexpr


```lua
string
```


---

# vim.o.foldignore


```lua
string
```


---

# vim.o.foldlevel


```lua
integer
```


---

# vim.o.foldlevelstart


```lua
integer
```


---

# vim.o.foldmarker


```lua
string
```


---

# vim.o.foldmethod


```lua
string
```


---

# vim.o.foldminlines


```lua
integer
```


---

# vim.o.foldnestmax


```lua
integer
```


---

# vim.o.foldopen


```lua
string
```


---

# vim.o.foldtext


```lua
string
```


---

# vim.o.formatexpr


```lua
string
```


---

# vim.o.formatlistpat


```lua
string
```


---

# vim.o.formatoptions


```lua
string
```


---

# vim.o.formatprg


```lua
string
```


---

# vim.o.fp


```lua
string
```


---

# vim.o.fs


```lua
boolean
```


---

# vim.o.fsync


```lua
boolean
```


---

# vim.o.ft


```lua
string
```


---

# vim.o.gcr


```lua
string
```


---

# vim.o.gd


```lua
boolean
```


---

# vim.o.gdefault


```lua
boolean
```


---

# vim.o.gfm


```lua
string
```


---

# vim.o.gfn


```lua
string
```


---

# vim.o.gfw


```lua
string
```


---

# vim.o.gp


```lua
string
```


---

# vim.o.grepformat


```lua
string
```


```lua
string
```


---

# vim.o.grepprg


```lua
string
```


```lua
string
```


---

# vim.o.guicursor


```lua
string
```


---

# vim.o.guifont


```lua
string
```


---

# vim.o.guifontwide


```lua
string
```


---

# vim.o.helpfile


```lua
string
```


---

# vim.o.helpheight


```lua
integer
```


---

# vim.o.helplang


```lua
string
```


---

# vim.o.hf


```lua
string
```


---

# vim.o.hh


```lua
integer
```


---

# vim.o.hi


```lua
integer
```


---

# vim.o.hid


```lua
boolean
```


---

# vim.o.hidden


```lua
boolean
```


---

# vim.o.history


```lua
integer
```


---

# vim.o.hlg


```lua
string
```


---

# vim.o.hls


```lua
boolean
```


---

# vim.o.hlsearch


```lua
boolean
```


---

# vim.o.ic


```lua
boolean
```


---

# vim.o.icm


```lua
''|'nosplit'|'split'
```


---

# vim.o.icon


```lua
boolean
```


---

# vim.o.iconstring


```lua
string
```


---

# vim.o.ignorecase


```lua
boolean
```


```lua
boolean
```


---

# vim.o.imi


```lua
integer
```


---

# vim.o.iminsert


```lua
integer
```


---

# vim.o.ims


```lua
integer
```


---

# vim.o.imsearch


```lua
integer
```


---

# vim.o.inc


```lua
string
```


---

# vim.o.inccommand


```lua
string
```


---

# vim.o.include


```lua
string
```


---

# vim.o.includeexpr


```lua
string
```


---

# vim.o.incsearch


```lua
boolean
```


---

# vim.o.inde


```lua
string
```


---

# vim.o.indentexpr


```lua
string
```


---

# vim.o.indentkeys


```lua
string
```


---

# vim.o.indk


```lua
string
```


---

# vim.o.inex


```lua
string
```


---

# vim.o.inf


```lua
boolean
```


---

# vim.o.infercase


```lua
boolean
```


---

# vim.o.is


```lua
boolean
```


---

# vim.o.isf


```lua
string
```


---

# vim.o.isfname


```lua
string
```


---

# vim.o.isi


```lua
string
```


---

# vim.o.isident


```lua
string
```


---

# vim.o.isk


```lua
string
```


---

# vim.o.iskeyword


```lua
string
```


---

# vim.o.isp


```lua
string
```


---

# vim.o.isprint


```lua
string
```


---

# vim.o.joinspaces


```lua
boolean
```


---

# vim.o.jop


```lua
string
```


---

# vim.o.js


```lua
boolean
```


---

# vim.o.jumpoptions


```lua
string
```


---

# vim.o.keymap


```lua
string
```


---

# vim.o.keymodel


```lua
string
```


---

# vim.o.keywordprg


```lua
string
```


---

# vim.o.km


```lua
string
```


---

# vim.o.kmp


```lua
string
```


---

# vim.o.kp


```lua
string
```


---

# vim.o.langmap


```lua
string
```


---

# vim.o.langmenu


```lua
string
```


---

# vim.o.langremap


```lua
boolean
```


---

# vim.o.laststatus


```lua
integer
```


---

# vim.o.lazyredraw


```lua
boolean
```


---

# vim.o.lbr


```lua
boolean
```


---

# vim.o.lcs


```lua
string
```


---

# vim.o.linebreak


```lua
boolean
```


---

# vim.o.lines


```lua
integer
```


---

# vim.o.linespace


```lua
integer
```


---

# vim.o.lisp


```lua
boolean
```


---

# vim.o.lispoptions


```lua
string
```


---

# vim.o.lispwords


```lua
string
```


---

# vim.o.list


```lua
boolean
```


```lua
boolean
```


---

# vim.o.listchars


```lua
string
```


---

# vim.o.lm


```lua
string
```


---

# vim.o.lmap


```lua
string
```


---

# vim.o.loadplugins


```lua
boolean
```


---

# vim.o.lop


```lua
string
```


---

# vim.o.lpl


```lua
boolean
```


---

# vim.o.lrm


```lua
boolean
```


---

# vim.o.ls


```lua
integer
```


---

# vim.o.lsp


```lua
integer
```


---

# vim.o.lw


```lua
string
```


---

# vim.o.lz


```lua
boolean
```


---

# vim.o.ma


```lua
boolean
```


---

# vim.o.magic


```lua
boolean
```


---

# vim.o.makeef


```lua
string
```


---

# vim.o.makeencoding


```lua
string
```


---

# vim.o.makeprg


```lua
string
```


---

# vim.o.mat


```lua
integer
```


---

# vim.o.matchpairs


```lua
string
```


---

# vim.o.matchtime


```lua
integer
```


---

# vim.o.maxfuncdepth


```lua
integer
```


---

# vim.o.maxmapdepth


```lua
integer
```


---

# vim.o.maxmempattern


```lua
integer
```


---

# vim.o.mef


```lua
string
```


---

# vim.o.menc


```lua
string
```


---

# vim.o.menuitems


```lua
integer
```


---

# vim.o.messagesopt


```lua
string
```


---

# vim.o.mfd


```lua
integer
```


---

# vim.o.mh


```lua
boolean
```


---

# vim.o.mis


```lua
integer
```


---

# vim.o.mkspellmem


```lua
string
```


---

# vim.o.ml


```lua
boolean
```


---

# vim.o.mle


```lua
boolean
```


---

# vim.o.mls


```lua
integer
```


---

# vim.o.mmd


```lua
integer
```


---

# vim.o.mmp


```lua
integer
```


---

# vim.o.mod


```lua
boolean
```


---

# vim.o.modeline


```lua
boolean
```


---

# vim.o.modelineexpr


```lua
boolean
```


---

# vim.o.modelines


```lua
integer
```


---

# vim.o.modifiable


```lua
boolean
```


---

# vim.o.modified


```lua
boolean
```


---

# vim.o.mopt


```lua
string
```


---

# vim.o.more


```lua
boolean
```


---

# vim.o.mouse


```lua
string
```


---

# vim.o.mousef


```lua
boolean
```


---

# vim.o.mousefocus


```lua
boolean
```


---

# vim.o.mousehide


```lua
boolean
```


---

# vim.o.mousem


```lua
'extend'|'popup'|'popup_setpos'
```


---

# vim.o.mousemev


```lua
boolean
```


---

# vim.o.mousemodel


```lua
string
```


---

# vim.o.mousemoveevent


```lua
boolean
```


---

# vim.o.mousescroll


```lua
string
```


---

# vim.o.mouset


```lua
integer
```


---

# vim.o.mousetime


```lua
integer
```


---

# vim.o.mp


```lua
string
```


---

# vim.o.mps


```lua
string
```


---

# vim.o.msm


```lua
string
```


---

# vim.o.nf


```lua
string
```


---

# vim.o.nrformats


```lua
string
```


---

# vim.o.nu


```lua
boolean
```


---

# vim.o.number


```lua
boolean
```


```lua
boolean
```


---

# vim.o.numberwidth


```lua
integer
```


---

# vim.o.nuw


```lua
integer
```


---

# vim.o.ofu


```lua
string
```


---

# vim.o.omnifunc


```lua
string
```


---

# vim.o.operatorfunc


```lua
string
```


```lua
string
```


---

# vim.o.opfunc


```lua
string
```


---

# vim.o.pa


```lua
string
```


---

# vim.o.packpath


```lua
string
```


---

# vim.o.para


```lua
string
```


---

# vim.o.paragraphs


```lua
string
```


---

# vim.o.patchexpr


```lua
string
```


---

# vim.o.patchmode


```lua
string
```


---

# vim.o.path


```lua
string
```


---

# vim.o.pb


```lua
integer
```


---

# vim.o.pex


```lua
string
```


---

# vim.o.ph


```lua
integer
```


---

# vim.o.pi


```lua
boolean
```


---

# vim.o.pm


```lua
string
```


---

# vim.o.pp


```lua
string
```


---

# vim.o.preserveindent


```lua
boolean
```


---

# vim.o.previewheight


```lua
integer
```


---

# vim.o.previewwindow


```lua
boolean
```


---

# vim.o.pumblend


```lua
integer
```


---

# vim.o.pumheight


```lua
integer
```


---

# vim.o.pumwidth


```lua
integer
```


---

# vim.o.pvh


```lua
integer
```


---

# vim.o.pvw


```lua
boolean
```


---

# vim.o.pw


```lua
integer
```


---

# vim.o.pyx


```lua
integer
```


---

# vim.o.pyxversion


```lua
integer
```


---

# vim.o.qe


```lua
string
```


---

# vim.o.qftf


```lua
string
```


---

# vim.o.quickfixtextfunc


```lua
string
```


---

# vim.o.quoteescape


```lua
string
```


---

# vim.o.rdb


```lua
string
```


---

# vim.o.rdt


```lua
integer
```


---

# vim.o.re


```lua
integer
```


---

# vim.o.readonly


```lua
boolean
```


---

# vim.o.redrawdebug


```lua
string
```


---

# vim.o.redrawtime


```lua
integer
```


---

# vim.o.regexpengine


```lua
integer
```


---

# vim.o.relativenumber


```lua
boolean
```


```lua
boolean
```


---

# vim.o.report


```lua
integer
```


---

# vim.o.revins


```lua
boolean
```


---

# vim.o.ri


```lua
boolean
```


---

# vim.o.rightleft


```lua
boolean
```


---

# vim.o.rightleftcmd


```lua
string
```


---

# vim.o.rl


```lua
boolean
```


---

# vim.o.rlc


```lua
string
```


---

# vim.o.rnu


```lua
boolean
```


---

# vim.o.ro


```lua
boolean
```


---

# vim.o.rtp


```lua
string
```


---

# vim.o.ru


```lua
boolean
```


---

# vim.o.ruf


```lua
string
```


---

# vim.o.ruler


```lua
boolean
```


---

# vim.o.rulerformat


```lua
string
```


---

# vim.o.runtimepath


```lua
string
```


---

# vim.o.sb


```lua
boolean
```


---

# vim.o.sbo


```lua
string
```


---

# vim.o.sbr


```lua
string
```


---

# vim.o.sc


```lua
boolean
```


---

# vim.o.scb


```lua
boolean
```


---

# vim.o.scbk


```lua
integer
```


---

# vim.o.scl


```lua
'auto'|'auto:1'|'auto:2'|'auto:3'|'auto:4'...(+17)
```


---

# vim.o.scr


```lua
integer
```


---

# vim.o.scroll


```lua
integer
```


---

# vim.o.scrollback


```lua
integer
```


---

# vim.o.scrollbind


```lua
boolean
```


---

# vim.o.scrolljump


```lua
integer
```


---

# vim.o.scrolloff


```lua
integer
```


```lua
integer
```


---

# vim.o.scrollopt


```lua
string
```


---

# vim.o.scs


```lua
boolean
```


---

# vim.o.sd


```lua
string
```


---

# vim.o.sdf


```lua
string
```


---

# vim.o.sect


```lua
string
```


---

# vim.o.sections


```lua
string
```


---

# vim.o.sel


```lua
'exclusive'|'inclusive'|'old'
```


---

# vim.o.selection


```lua
string
```


---

# vim.o.selectmode


```lua
string
```


---

# vim.o.sessionoptions


```lua
string
```


---

# vim.o.sft


```lua
boolean
```


---

# vim.o.sh


```lua
string
```


---

# vim.o.shada


```lua
string
```


---

# vim.o.shadafile


```lua
string
```


---

# vim.o.shcf


```lua
string
```


---

# vim.o.shell


```lua
string
```


---

# vim.o.shellcmdflag


```lua
string
```


---

# vim.o.shellpipe


```lua
string
```


---

# vim.o.shellquote


```lua
string
```


---

# vim.o.shellredir


```lua
string
```


---

# vim.o.shellslash


```lua
boolean
```


---

# vim.o.shelltemp


```lua
boolean
```


---

# vim.o.shellxescape


```lua
string
```


---

# vim.o.shellxquote


```lua
string
```


---

# vim.o.shiftround


```lua
boolean
```


---

# vim.o.shiftwidth


```lua
integer
```


---

# vim.o.shm


```lua
string
```


---

# vim.o.shortmess


```lua
string
```


---

# vim.o.showbreak


```lua
string
```


---

# vim.o.showcmd


```lua
boolean
```


---

# vim.o.showcmdloc


```lua
string
```


---

# vim.o.showfulltag


```lua
boolean
```


---

# vim.o.showmatch


```lua
boolean
```


---

# vim.o.showmode


```lua
boolean
```


---

# vim.o.showtabline


```lua
integer
```


---

# vim.o.shq


```lua
string
```


---

# vim.o.si


```lua
boolean
```


---

# vim.o.sidescroll


```lua
integer
```


---

# vim.o.sidescrolloff


```lua
integer
```


---

# vim.o.signcolumn


```lua
string
```


---

# vim.o.siso


```lua
integer
```


---

# vim.o.sj


```lua
integer
```


---

# vim.o.slm


```lua
string
```


---

# vim.o.sloc


```lua
'last'|'statusline'|'tabline'
```


---

# vim.o.sm


```lua
boolean
```


---

# vim.o.smartcase


```lua
boolean
```


```lua
boolean
```


---

# vim.o.smartindent


```lua
boolean
```


---

# vim.o.smarttab


```lua
boolean
```


---

# vim.o.smc


```lua
integer
```


---

# vim.o.smd


```lua
boolean
```


---

# vim.o.smoothscroll


```lua
boolean
```


---

# vim.o.sms


```lua
boolean
```


---

# vim.o.so


```lua
integer
```


---

# vim.o.softtabstop


```lua
integer
```


---

# vim.o.sol


```lua
boolean
```


---

# vim.o.sp


```lua
string
```


---

# vim.o.spc


```lua
string
```


---

# vim.o.spell


```lua
boolean
```


---

# vim.o.spellcapcheck


```lua
string
```


---

# vim.o.spellfile


```lua
string
```


---

# vim.o.spelllang


```lua
string
```


---

# vim.o.spelloptions


```lua
string
```


---

# vim.o.spellsuggest


```lua
string
```


---

# vim.o.spf


```lua
string
```


---

# vim.o.spk


```lua
'cursor'|'screen'|'topline'
```


---

# vim.o.spl


```lua
string
```


---

# vim.o.splitbelow


```lua
boolean
```


---

# vim.o.splitkeep


```lua
string
```


---

# vim.o.splitright


```lua
boolean
```


---

# vim.o.spo


```lua
string
```


---

# vim.o.spr


```lua
boolean
```


---

# vim.o.sps


```lua
string
```


---

# vim.o.sr


```lua
boolean
```


---

# vim.o.srr


```lua
string
```


---

# vim.o.ss


```lua
integer
```


---

# vim.o.ssl


```lua
boolean
```


---

# vim.o.ssop


```lua
string
```


---

# vim.o.sta


```lua
boolean
```


---

# vim.o.stal


```lua
integer
```


---

# vim.o.startofline


```lua
boolean
```


---

# vim.o.statuscolumn


```lua
string
```


---

# vim.o.statusline


```lua
string
```


---

# vim.o.stc


```lua
string
```


---

# vim.o.stl


```lua
string
```


---

# vim.o.stmp


```lua
boolean
```


---

# vim.o.sts


```lua
integer
```


---

# vim.o.su


```lua
string
```


---

# vim.o.sua


```lua
string
```


---

# vim.o.suffixes


```lua
string
```


---

# vim.o.suffixesadd


```lua
string
```


---

# vim.o.sw


```lua
integer
```


---

# vim.o.swapfile


```lua
boolean
```


---

# vim.o.swb


```lua
string
```


---

# vim.o.swf


```lua
boolean
```


---

# vim.o.switchbuf


```lua
string
```


---

# vim.o.sxe


```lua
string
```


---

# vim.o.sxq


```lua
string
```


---

# vim.o.syn


```lua
string
```


---

# vim.o.synmaxcol


```lua
integer
```


---

# vim.o.syntax


```lua
string
```


---

# vim.o.tabclose


```lua
string
```


---

# vim.o.tabline


```lua
string
```


---

# vim.o.tabpagemax


```lua
integer
```


---

# vim.o.tabstop


```lua
integer
```


---

# vim.o.tag


```lua
string
```


---

# vim.o.tagbsearch


```lua
boolean
```


---

# vim.o.tagcase


```lua
string
```


---

# vim.o.tagfunc


```lua
string
```


---

# vim.o.taglength


```lua
integer
```


---

# vim.o.tagrelative


```lua
boolean
```


---

# vim.o.tags


```lua
string
```


---

# vim.o.tagstack


```lua
boolean
```


---

# vim.o.tal


```lua
string
```


---

# vim.o.tbidi


```lua
boolean
```


---

# vim.o.tbs


```lua
boolean
```


---

# vim.o.tc


```lua
'followic'|'followscs'|'ignore'|'match'|'smart'
```


---

# vim.o.tcl


```lua
string
```


---

# vim.o.termbidi


```lua
boolean
```


---

# vim.o.termguicolors


```lua
boolean
```


---

# vim.o.termpastefilter


```lua
string
```


---

# vim.o.termsync


```lua
boolean
```


---

# vim.o.textwidth


```lua
integer
```


---

# vim.o.tfu


```lua
string
```


---

# vim.o.tgc


```lua
boolean
```


---

# vim.o.tgst


```lua
boolean
```


---

# vim.o.thesaurus


```lua
string
```


---

# vim.o.thesaurusfunc


```lua
string
```


---

# vim.o.tildeop


```lua
boolean
```


---

# vim.o.timeout


```lua
boolean
```


---

# vim.o.timeoutlen


```lua
integer
```


---

# vim.o.title


```lua
boolean
```


---

# vim.o.titlelen


```lua
integer
```


---

# vim.o.titleold


```lua
string
```


---

# vim.o.titlestring


```lua
string
```


---

# vim.o.tl


```lua
integer
```


---

# vim.o.tm


```lua
integer
```


---

# vim.o.to


```lua
boolean
```


---

# vim.o.top


```lua
boolean
```


---

# vim.o.tpf


```lua
string
```


---

# vim.o.tpm


```lua
integer
```


---

# vim.o.tr


```lua
boolean
```


---

# vim.o.ts


```lua
integer
```


---

# vim.o.tsr


```lua
string
```


---

# vim.o.tsrfu


```lua
string
```


---

# vim.o.ttimeout


```lua
boolean
```


---

# vim.o.ttimeoutlen


```lua
integer
```


---

# vim.o.ttm


```lua
integer
```


---

# vim.o.tw


```lua
integer
```


---

# vim.o.uc


```lua
integer
```


---

# vim.o.udf


```lua
boolean
```


---

# vim.o.udir


```lua
string
```


---

# vim.o.ul


```lua
integer
```


---

# vim.o.undodir


```lua
string
```


---

# vim.o.undofile


```lua
boolean
```


---

# vim.o.undolevels


```lua
integer
```


---

# vim.o.undoreload


```lua
integer
```


---

# vim.o.updatecount


```lua
integer
```


---

# vim.o.updatetime


```lua
integer
```


---

# vim.o.ur


```lua
integer
```


---

# vim.o.ut


```lua
integer
```


---

# vim.o.varsofttabstop


```lua
string
```


---

# vim.o.vartabstop


```lua
string
```


---

# vim.o.vb


```lua
boolean
```


---

# vim.o.vbs


```lua
integer
```


---

# vim.o.vdir


```lua
string
```


---

# vim.o.ve


```lua
string
```


---

# vim.o.verbose


```lua
integer
```


---

# vim.o.verbosefile


```lua
string
```


---

# vim.o.vfile


```lua
string
```


---

# vim.o.viewdir


```lua
string
```


---

# vim.o.viewoptions


```lua
string
```


---

# vim.o.virtualedit


```lua
string
```


---

# vim.o.visualbell


```lua
boolean
```


---

# vim.o.vop


```lua
string
```


---

# vim.o.vsts


```lua
string
```


---

# vim.o.vts


```lua
string
```


---

# vim.o.wa


```lua
boolean
```


---

# vim.o.wak


```lua
'menu'|'no'|'yes'
```


---

# vim.o.warn


```lua
boolean
```


---

# vim.o.wb


```lua
boolean
```


---

# vim.o.wbr


```lua
string
```


---

# vim.o.wc


```lua
integer
```


---

# vim.o.wcm


```lua
integer
```


---

# vim.o.wd


```lua
integer
```


---

# vim.o.wfb


```lua
boolean
```


---

# vim.o.wfh


```lua
boolean
```


---

# vim.o.wfw


```lua
boolean
```


---

# vim.o.wh


```lua
integer
```


---

# vim.o.whichwrap


```lua
string
```


---

# vim.o.wi


```lua
integer
```


---

# vim.o.wic


```lua
boolean
```


---

# vim.o.wig


```lua
string
```


---

# vim.o.wildchar


```lua
integer
```


---

# vim.o.wildcharm


```lua
integer
```


---

# vim.o.wildignore


```lua
string
```


---

# vim.o.wildignorecase


```lua
boolean
```


---

# vim.o.wildmenu


```lua
boolean
```


---

# vim.o.wildmode


```lua
string
```


---

# vim.o.wildoptions


```lua
string
```


---

# vim.o.wim


```lua
string
```


---

# vim.o.winaltkeys


```lua
string
```


---

# vim.o.winbar


```lua
string
```


---

# vim.o.winbl


```lua
integer
```


---

# vim.o.winblend


```lua
integer
```


---

# vim.o.winborder


```lua
string
```


---

# vim.o.window


```lua
integer
```


---

# vim.o.winfixbuf


```lua
boolean
```


---

# vim.o.winfixheight


```lua
boolean
```


---

# vim.o.winfixwidth


```lua
boolean
```


---

# vim.o.winheight


```lua
integer
```


---

# vim.o.winhighlight


```lua
string
```


---

# vim.o.winhl


```lua
string
```


---

# vim.o.winminheight


```lua
integer
```


---

# vim.o.winminwidth


```lua
integer
```


---

# vim.o.winwidth


```lua
integer
```


---

# vim.o.wiw


```lua
integer
```


---

# vim.o.wm


```lua
integer
```


---

# vim.o.wmh


```lua
integer
```


---

# vim.o.wmnu


```lua
boolean
```


---

# vim.o.wmw


```lua
integer
```


---

# vim.o.wop


```lua
string
```


---

# vim.o.wrap


```lua
boolean
```


---

# vim.o.wrapmargin


```lua
integer
```


---

# vim.o.wrapscan


```lua
boolean
```


---

# vim.o.write


```lua
boolean
```


---

# vim.o.writeany


```lua
boolean
```


---

# vim.o.writebackup


```lua
boolean
```


---

# vim.o.writedelay


```lua
integer
```


---

# vim.o.ws


```lua
boolean
```


---

# vim.o.ww


```lua
string
```


---

# vim.on_key


```lua
function vim.on_key(fn: fun(key: string, typed: string):string?|nil, ns_id?: integer, opts?: table)
  -> Namespace: integer
```


---

# vim.opt


```lua
table
```


---

# vim.opt_global


```lua
table
```


---

# vim.opt_local


```lua
table
```


---

# vim.paste


```lua
function vim.paste(lines: string[], phase: -1|1|2|3)
  -> result: boolean
```


---

# vim.pesc


```lua
function vim.pesc(s: string)
  -> string
```


---

# vim.print


```lua
function vim.print(...any)
  -> any
```


---

# vim.provider


```lua
table
```


---

# vim.quickfix.entry

## bufnr


```lua
integer?
```

 buffer number; must be the number of a valid buffer

## col


```lua
integer?
```


 column number

## end_col


```lua
integer?
```


 end column, if the item spans multiple columns

## end_lnum


```lua
integer?
```


 end of lines, if the item spans multiple lines

## filename


```lua
string?
```


 name of a file; only used when "bufnr" is not
 present or it is invalid.

## lnum


```lua
integer?
```


 line number in the file

## module


```lua
string?
```


 name of a module; if given it will be used in
 quickfix error window instead of the filename.

## nr


```lua
integer?
```


 error number

## pattern


```lua
string?
```


 search pattern used to locate the error

## text


```lua
string?
```


 description of the error

## type


```lua
string?
```


 single-character error type, 'E', 'W', etc.

## user_data


```lua
any
```


 custom data associated with the item, can be
 any type.

## valid


```lua
boolean?
```


 recognized error message

## vcol


```lua
integer?
```


 when non-zero: "col" is visual column
 when zero: "col" is byte index


---

# vim.re


```lua
table
```


---

# vim.re.compile


```lua
function vim.re.compile(string: string, defs?: table)
  -> vim.lpeg.Pattern
```


---

# vim.re.find


```lua
function vim.re.find(subject: string, pattern: string|vim.lpeg.Pattern, init?: integer)
  -> integer|nil
  2. integer|nil
```


---

# vim.re.gsub


```lua
function vim.re.gsub(subject: string, pattern: string|vim.lpeg.Pattern, replacement: string)
  -> string
```


---

# vim.re.match


```lua
function vim.re.match(subject: string, pattern: string|vim.lpeg.Pattern, init?: integer)
  -> integer|vim.lpeg.Pattern|nil
```


---

# vim.re.updatelocale


```lua
function vim.re.updatelocale()
```


---

# vim.regex


```lua
function vim.regex(re: string)
  -> vim.regex
```


---

# vim.regex

 @nodoc

## match_line


```lua
(method) vim.regex:match_line(bufnr: integer, line_idx: integer, start?: integer, end_?: integer)
  -> integer?
  2. integer?
```

 Matches line at `line_idx` (zero-based) in buffer `bufnr`. Match is restricted to byte index
 range `start` and `end_` if given, otherwise see |regex:match_str()|. Returned byte indices are
 relative to `start` if given.

@*return* — match start (byte index) relative to `start`, or `nil` if no match

@*return* — match end (byte index) relative to `start`, or `nil` if no match

## match_str


```lua
(method) vim.regex:match_str(str: string)
  -> integer?
  2. integer?
```

 Matches string `str` against this regex. To match the string precisely, surround the regex with
 "^" and "$". Returns the byte indices for the start and end of the match, or `nil` if there is
 no match. Because any integer is "truthy", `regex:match_str()` can be directly used as
 a condition in an if-statement.

@*return* — match start (byte index), or `nil` if no match

@*return* — match end (byte index), or `nil` if no match


---

# vim.region


```lua
function vim.region(bufnr: integer, pos1: string|integer[], pos2: string|integer[], regtype: string, inclusive: boolean)
  -> region: table
```


---

# vim.ringbuf


```lua
function vim.ringbuf(size: integer)
  -> ringbuf: vim.Ringbuf
```


---

# vim.rpcnotify


```lua
function vim.rpcnotify(channel: integer, method: string, ...any)
```


---

# vim.rpcrequest


```lua
function vim.rpcrequest(channel: integer, method: string, ...any)
```


---

# vim.schedule


```lua
function vim.schedule(fn: fun())
```


---

# vim.schedule_wrap


```lua
function vim.schedule_wrap(fn: function)
  -> function
```


---

# vim.secure


```lua
table
```


---

# vim.show_pos


```lua
function vim.show_pos(bufnr?: integer, row?: integer, col?: integer, filter?: vim._inspector.Filter)
```


---

# vim.snippet


```lua
table
```


---

# vim.snippet.ActiveFilter

## direction


```lua
-1|1
```

Navigation direction. -1 for previous, 1 for next.


---

# vim.snippet.ChoiceData


---

# vim.snippet.Direction


---

# vim.snippet.FormatData


---

# vim.snippet.Node


---

# vim.snippet.PlaceholderData


---

# vim.snippet.Session

## bufnr


```lua
integer
```

## current_tabstop


```lua
vim.snippet.Tabstop
```

## extmark_id


```lua
integer
```

## get_dest_index


```lua
(method) vim.snippet.Session:get_dest_index(direction: -1|1)
  -> integer?
```

 Returns the destination tabstop index when jumping in the given direction.


```lua
direction:
    | -1
    | 1
```

## new


```lua
function vim.snippet.Session.new(bufnr: integer, snippet_extmark: integer, tabstop_data: table<integer, { range: Range4, choices: string[] }[]>)
  -> vim.snippet.Session
```

 Creates a new snippet session in the current buffer.

## set_group_gravity


```lua
(method) vim.snippet.Session:set_group_gravity(index: integer, right_gravity: boolean)
```

 Sets the right gravity of the tabstop group with the given index.

## shift_tab_keymaps


```lua
{ i: table<string, any>?, s: table<string, any>? }
```

## tab_keymaps


```lua
{ i: table<string, any>?, s: table<string, any>? }
```

## tabstops


```lua
table<integer, vim.snippet.Tabstop[]>
```


---

# vim.snippet.SnippetData


---

# vim.snippet.Tabstop

## bufnr


```lua
integer
```

## choices


```lua
string[]?
```

## extmark_id


```lua
integer
```

## get_range


```lua
(method) vim.snippet.Tabstop:get_range()
  -> Range4
```

 Returns the tabstop's range.

## get_text


```lua
(method) vim.snippet.Tabstop:get_text()
  -> string
```

 Returns the text spanned by the tabstop.

## index


```lua
integer
```

## new


```lua
function vim.snippet.Tabstop.new(index: integer, bufnr: integer, range: Range4, choices?: string[])
  -> vim.snippet.Tabstop
```

 Creates a new tabstop.

## set_right_gravity


```lua
(method) vim.snippet.Tabstop:set_right_gravity(right_gravity: boolean)
```

 Sets the right gravity of the tabstop's extmark.

## set_text


```lua
(method) vim.snippet.Tabstop:set_text(text: string)
```

 Sets the tabstop's text.


---

# vim.snippet.TabstopData


---

# vim.snippet.TextData


---

# vim.snippet.Type


---

# vim.snippet.VariableData


---

# vim.spairs


```lua
function vim.spairs(t: <T:table>)
  -> fun(table: table<K, V>, index?: <K>):<K>, <V>
  2. <T:table>
```


---

# vim.spell.check


```lua
function vim.spell.check(str: string)
  -> [string, 'bad'|'caps'|'local'|'rare', integer][]
```


---

# vim.split


```lua
function vim.split(s: string, sep: string, opts?: vim.gsplit.Opts)
  -> string[]
```


---

# vim.startswith


```lua
function vim.startswith(s: string, prefix: string)
  -> boolean
```


---

# vim.str_byteindex


```lua
function vim.str_byteindex(s: string, encoding: "utf-16"|"utf-32"|"utf-8", index: integer, strict_indexing?: boolean)
  -> integer
```


---

# vim.str_utf_end


```lua
function vim.str_utf_end(str: string, index: integer)
  -> integer
```


---

# vim.str_utf_pos


```lua
function vim.str_utf_pos(str: string)
  -> integer[]
```


---

# vim.str_utf_start


```lua
function vim.str_utf_start(str: string, index: integer)
  -> integer
```


---

# vim.str_utfindex


```lua
function vim.str_utfindex(s: string, encoding: "utf-16"|"utf-32"|"utf-8", index?: integer, strict_indexing?: boolean)
  -> integer
```


---

# vim.stricmp


```lua
function vim.stricmp(a: string, b: string)
  -> -1|0|1
```


---

# vim.system


```lua
function vim.system(cmd: string[], opts?: vim.SystemOpts, on_exit?: fun(out: vim.SystemCompleted))
  -> Object: vim.SystemObj
```


---

# vim.t


```lua
vim.var_accessor
```


---

# vim.tbl_add_reverse_lookup


```lua
function vim.tbl_add_reverse_lookup(o: table)
  -> o: table
```


---

# vim.tbl_contains


```lua
function vim.tbl_contains(t: table, value: any, opts?: vim.tbl_contains.Opts)
  -> boolean
```


---

# vim.tbl_contains.Opts

## predicate


```lua
boolean?
```

 @inlinedoc

 `value` is a function reference to be checked (default false)


---

# vim.tbl_count


```lua
function vim.tbl_count(t: table)
  -> integer
```


---

# vim.tbl_deep_extend


```lua
function vim.tbl_deep_extend(behavior: 'error'|'force'|'keep', ...<T2:table>)
  -> <T1:table>|<T2:table>
```


---

# vim.tbl_extend


```lua
function vim.tbl_extend(behavior: 'error'|'force'|'keep', ...table)
  -> table
```


---

# vim.tbl_filter


```lua
function vim.tbl_filter(func: fun(value: <T>):boolean, t: table<any, T>)
  -> <T>[]
```


---

# vim.tbl_flatten


```lua
function vim.tbl_flatten(t: table)
  -> Flattened: table
```


---

# vim.tbl_get


```lua
function vim.tbl_get(o: table, ...any)
  -> any
```


---

# vim.tbl_isempty


```lua
function vim.tbl_isempty(t: table)
  -> boolean
```


---

# vim.tbl_islist


```lua
function vim.tbl_islist(t: any)
  -> boolean
```


---

# vim.tbl_keys


```lua
function vim.tbl_keys(t: table<T, any>)
  -> <T>[]
```


---

# vim.tbl_map


```lua
function vim.tbl_map(func: fun(value: <T>):any, t: table<any, T>)
  -> table
```


---

# vim.tbl_values


```lua
function vim.tbl_values(t: table<any, T>)
  -> <T>[]
```


---

# vim.text


```lua
table
```


---

# vim.tohtml.cell

## [1]


```lua
integer[]
```

start

## [2]


```lua
integer[]
```

close

## [3]


```lua
any[][]
```

virt_text

## [4]


```lua
any[][]
```

overlay_text


---

# vim.tohtml.line

## [integer]


```lua
(vim.tohtml.cell)?
```

(integer: (1-index, exclusive))

## hide


```lua
boolean?
```

## pre_text


```lua
string[][]
```

## virt_lines


```lua
{ [integer]: [string, integer][] }
```


---

# vim.tohtml.opt

## font


```lua
(string|string[])?
```


 Fonts to use.
 (default: `guifont`)

## number_lines


```lua
boolean?
```


 Show line numbers.
 (default: `false`)

## range


```lua
integer[]?
```


 Range of rows to use.
 (default: entire buffer)

## title


```lua
(string|false)?
```

 @inlinedoc

 Title tag to set in the generated HTML code.
 (default: buffer name)

## width


```lua
integer?
```


 Width used for items which are either right aligned or repeat a character
 infinitely.
 (default: 'textwidth' if non-zero or window width otherwise)


---

# vim.tohtml.state

## background


```lua
string
```

## bufnr


```lua
integer
```

## conf


```lua
vim.tohtml.opt
```

## end_


```lua
integer
```

## font


```lua
string
```

## foreground


```lua
string
```

## highlights_name


```lua
table<integer, string>
```

## opt


```lua
vim.wo
```

## start


```lua
integer
```

## style


```lua
vim.tohtml.styletable
```

## tabstop


```lua
string|false
```

## title


```lua
string|false
```

## width


```lua
integer
```

## winid


```lua
integer
```


---

# vim.tohtml.state.global

## background


```lua
string
```

## conf


```lua
vim.tohtml.opt
```

## font


```lua
string
```

## foreground


```lua
string
```

## highlights_name


```lua
table<integer, string>
```

## title


```lua
string|false
```


---

# vim.tohtml.styletable

## [integer]


```lua
vim.tohtml.line
```

(integer: (1-index, exclusive))


---

# vim.treesitter


```lua
table
```


---

# vim.treesitter.LanguageTree

## __index


```lua
vim.treesitter.LanguageTree
```

## _add_injections


```lua
(method) vim.treesitter.LanguageTree:_add_injections(injections_by_lang: table<string, Range6[][]>)
```

## _async_parse


```lua
(method) vim.treesitter.LanguageTree:_async_parse(range?: boolean|Range2|Range4|Range6, on_parse: fun(err?: string, trees?: table<integer, TSTree>))
  -> trees: table<integer, TSTree>?
```

 Run an asynchronous parse, calling {on_parse} when complete.

@*return* `trees` — the list of parsed trees, if parsing completed synchronously

## _callbacks


```lua
table<'bytes'|'changedtree'|'child_added'|'child_removed'|'detach', function[]>
```

Callback handlers

## _callbacks_rec


```lua
table<'bytes'|'changedtree'|'child_added'|'child_removed'|'detach', function[]>
```

Callback handlers (recursive)

## _cb_queues


```lua
table<string, fun(err?: string, trees?: table<integer, TSTree>)[]>
```

Table of callback queues, keyed by each region for which the callbacks should be run

## _children


```lua
table<string, vim.treesitter.LanguageTree>
```

Injected languages

## _do_callback


```lua
(method) vim.treesitter.LanguageTree:_do_callback(cb_name: 'bytes'|'changedtree'|'child_added'|'child_removed'|'detach', ...any)
```

```lua
cb_name:
    | 'changedtree'
    | 'bytes'
    | 'detach'
    | 'child_added'
    | 'child_removed'
```

## _edit


```lua
(method) vim.treesitter.LanguageTree:_edit(start_byte: any, end_byte_old: any, end_byte_new: any, start_row: any, start_col: any, end_row_old: any, end_col_old: any, end_row_new: any, end_col_new: any)
```

## _get_injection


```lua
(method) vim.treesitter.LanguageTree:_get_injection(match: table<integer, TSNode[]>, metadata: vim.treesitter.query.TSMetadata)
  -> string?
  2. boolean
  3. Range6[]
```

 Extract injections according to:
 https://tree-sitter.github.io/tree-sitter/syntax-highlighting#language-injection

## _get_injections


```lua
(method) vim.treesitter.LanguageTree:_get_injections(range: Range2|Range4|Range6|true, thread_state: { timeout: integer? })
  -> table<string, Range6[][]>
```

 Gets language injection regions by language.

 This is where most of the injection processing occurs.

 TODO: Allow for an offset predicate to tailor the injection range
       instead of using the entire nodes range.

```lua
range:
    | true
```

## _injection_query


```lua
vim.treesitter.Query
```

Queries defining injected languages

## _is_entirely_valid


```lua
boolean
```

Whether the entire tree (excluding children) is valid.

## _iter_regions


```lua
(method) vim.treesitter.LanguageTree:_iter_regions(fn: fun(index: integer, region: Range6[]):boolean)
```

Iterate through all the regions. fn returns a boolean to indicate if the
region is valid or not.

## _lang


```lua
string
```

Language name

## _log


```lua
(method) vim.treesitter.LanguageTree:_log(...any)
```

## _logfile


```lua
file*?
```




[View documents](command:extension.lua.doc?["en-us/54/manual.html/pdf-file"])


## _logger


```lua
fun(logtype: string, msg: string)?
```

## _num_regions


```lua
integer
```

The total number of regions. Since _regions can have holes, we cannot simply read this value from #_regions.

## _num_valid_regions


```lua
integer
```

Number of valid regions

## _on_bytes


```lua
(method) vim.treesitter.LanguageTree:_on_bytes(bufnr: integer, changed_tick: integer, start_row: integer, start_col: integer, start_byte: integer, old_row: integer, old_col: integer, old_byte: integer, new_row: integer, new_col: integer, new_byte: integer)
```

## _on_detach


```lua
(method) vim.treesitter.LanguageTree:_on_detach(...any)
```

## _on_reload


```lua
(method) vim.treesitter.LanguageTree:_on_reload()
```

## _opts


```lua
table
```

Options

## _parent


```lua
(vim.treesitter.LanguageTree)?
```

Parent LanguageTree

## _parse


```lua
(method) vim.treesitter.LanguageTree:_parse(range: boolean|Range2|Range4|Range6|nil, thread_state: { timeout: integer? })
  -> trees: table<integer, TSTree>
  2. finished: boolean
```

## _parse_regions


```lua
(method) vim.treesitter.LanguageTree:_parse_regions(range?: boolean|Range2|Range4|Range6, thread_state: { timeout: integer? })
  -> changes: Range6[]
  2. no_regions_parsed: integer
  3. total_parse_time: number
```

## _parser


```lua
TSParser
```

Parser for language

## _processed_injection_range


```lua
(Range2|Range4|Range6)?
```

Range for which injections have been processed

## _push_async_callback


```lua
(method) vim.treesitter.LanguageTree:_push_async_callback(range?: boolean|Range2|Range4|Range6, callback: fun(err?: string, trees?: table<integer, TSTree>))
```

## _ranges_being_parsed


```lua
table<string, boolean>
```

Table of regions for which the tree is currently running an async parse

## _regions


```lua
table<integer, Range6[]>?
```

## _run_async_callbacks


```lua
(method) vim.treesitter.LanguageTree:_run_async_callbacks(range?: boolean|Range2|Range4|Range6, err?: string, trees?: table<integer, TSTree>)
```

## _set_logger


```lua
(method) vim.treesitter.LanguageTree:_set_logger()
```

## _source


```lua
string|integer
```

Buffer or string to parse

## _subtract_time


```lua
(method) vim.treesitter.LanguageTree:_subtract_time(thread_state: { timeout: integer? }, time: integer)
```

## _trees


```lua
table<integer, TSTree>
```

Reference to parsed tree (one for each language).

## _valid_regions


```lua
table<integer, true>
```

Set of valid region IDs.

## add_child


```lua
(method) vim.treesitter.LanguageTree:add_child(lang: string)
  -> injected: vim.treesitter.LanguageTree
```

 Adds a child language to this |LanguageTree|.

 If the language already exists as a child, it will first be removed.

@*param* `lang` — Language to add.

## children


```lua
(method) vim.treesitter.LanguageTree:children()
  -> table<string, vim.treesitter.LanguageTree>
```

 Returns a map of language to child tree.

## contains


```lua
(method) vim.treesitter.LanguageTree:contains(range: Range4)
  -> boolean
```

 Determines whether {range} is contained in the |LanguageTree|.

## destroy


```lua
(method) vim.treesitter.LanguageTree:destroy()
```

 Destroys this |LanguageTree| and all its children.

 Any cleanup logic should be performed here.

 Note: This DOES NOT remove this tree from a parent. Instead,
 `remove_child` must be called on the parent to remove it.

## for_each_tree


```lua
(method) vim.treesitter.LanguageTree:for_each_tree(fn: fun(tree: TSTree, ltree: vim.treesitter.LanguageTree))
```

 Invokes the callback for each |LanguageTree| recursively.

 Note: This includes the invoking tree's child trees as well.

## included_regions


```lua
(method) vim.treesitter.LanguageTree:included_regions()
  -> table<integer, Range6[]>
```

Gets the set of included regions managed by this LanguageTree. This can be different from the
regions set by injection query, because a partial |LanguageTree:parse()| drops the regions
outside the requested range.
Each list represents a range in the form of
{ {start_row}, {start_col}, {start_bytes}, {end_row}, {end_col}, {end_bytes} }.

## invalidate


```lua
(method) vim.treesitter.LanguageTree:invalidate(reload: boolean|nil)
```

 Invalidates this parser and its children.

 Should only be called when the tracked state of the LanguageTree is not valid against the parse
 tree in treesitter. Doesn't clear filesystem cache. Called often, so needs to be fast.

## is_valid


```lua
(method) vim.treesitter.LanguageTree:is_valid(exclude_children?: boolean, range?: Range2|Range4|Range6)
  -> boolean
```

 Returns whether this LanguageTree is valid, i.e., |LanguageTree:trees()| reflects the latest
 state of the source. If invalid, user should call |LanguageTree:parse()|.

@*param* `exclude_children` — whether to ignore the validity of children (default `false`)

@*param* `range` — range to check for validity

## lang


```lua
(method) vim.treesitter.LanguageTree:lang()
  -> string
```

 Gets the language of this tree node.

## language_for_range


```lua
(method) vim.treesitter.LanguageTree:language_for_range(range: Range4)
  -> tree: vim.treesitter.LanguageTree
```

 Gets the appropriate language that contains {range}.

@*return* `tree` — Managing {range}

## named_node_for_range


```lua
(method) vim.treesitter.LanguageTree:named_node_for_range(range: Range4, opts?: vim.treesitter.LanguageTree.tree_for_range.Opts)
  -> TSNode?
```

 Gets the smallest named node that contains {range}.

## new


```lua
function vim.treesitter.LanguageTree.new(source: string|integer, lang: string, opts?: vim.treesitter.LanguageTree.new.Opts)
  -> parser: vim.treesitter.LanguageTree
```

 @nodoc

 LanguageTree contains a tree of parsers: the root treesitter parser for {lang} and any
 "injected" language parsers, which themselves may inject other languages, recursively.

@*param* `source` — Buffer or text string to parse

@*param* `lang` — Root language of this tree

@*return* `parser` — object

## node_for_range


```lua
(method) vim.treesitter.LanguageTree:node_for_range(range: Range4, opts?: vim.treesitter.LanguageTree.tree_for_range.Opts)
  -> TSNode?
```

 Gets the smallest node that contains {range}.

## parent


```lua
(method) vim.treesitter.LanguageTree:parent()
  -> (vim.treesitter.LanguageTree)?
```

Returns the parent tree. `nil` for the root tree.

## parse


```lua
(method) vim.treesitter.LanguageTree:parse(range: boolean|Range2|Range4|Range6|nil, on_parse?: fun(err?: string, trees?: table<integer, TSTree>))
  -> table<integer, TSTree>?
```

 Recursively parse all regions in the language tree using |treesitter-parsers|
 for the corresponding languages and run injection queries on the parsed trees
 to determine whether child trees should be created and parsed.

 Any region with empty range (`{}`, typically only the root tree) is always parsed;
 otherwise (typically injections) only if it intersects {range} (or if {range} is `true`).

@*param* `range` — : Parse this range in the parser's source.

     Set to `true` to run a complete parse of the source (Note: Can be slow!)
     Set to `false|nil` to only parse regions with empty ranges (typically
     only the root tree without injections).

@*param* `on_parse` — Function invoked when parsing completes.

     When provided and `vim.g._ts_force_sync_parsing` is not set, parsing will run
     asynchronously. The first argument to the function is a string representing the error type,
     in case of a failure (currently only possible for timeouts). The second argument is the list
     of trees returned by the parse (upon success), or `nil` if the parse timed out (determined
     by 'redrawtime').

     If parsing was still able to finish synchronously (within 3ms), `parse()` returns the list
     of trees. Otherwise, it returns `nil`.

## register_cbs


```lua
(method) vim.treesitter.LanguageTree:register_cbs(cbs: table<'on_bytes'|'on_changedtree'|'on_child_added'|'on_child_removed'|'on_detach', function>, recursive?: boolean)
```

 Registers callbacks for the [LanguageTree].

@*param* `cbs` — An [nvim_buf_attach()]-like table argument with the following handlers:

           - `on_bytes` : see [nvim_buf_attach()].
           - `on_changedtree` : a callback that will be called every time the tree has syntactical changes.
              It will be passed two arguments: a table of the ranges (as node ranges) that
              changed and the changed tree.
           - `on_child_added` : emitted when a child is added to the tree.
           - `on_child_removed` : emitted when a child is removed from the tree.
           - `on_detach` : emitted when the buffer is detached, see [nvim_buf_detach_event].
              Takes one argument, the number of the buffer.

@*param* `recursive` — Apply callbacks recursively for all children. Any new children will

                           also inherit the callbacks.

## remove_child


```lua
(method) vim.treesitter.LanguageTree:remove_child(lang: string)
```

 Removes a child language from this |LanguageTree|.

@*param* `lang` — Language to remove.

## set_included_regions


```lua
(method) vim.treesitter.LanguageTree:set_included_regions(new_regions: (Range4|Range6|TSNode)[][])
```

 Sets the included regions that should be parsed by this |LanguageTree|.
 A region is a set of nodes and/or ranges that will be parsed in the same context.

 For example, `{ { node1 }, { node2} }` contains two separate regions.
 They will be parsed by the parser in two different contexts, thus resulting
 in two separate trees.

 On the other hand, `{ { node1, node2 } }` is a single region consisting of
 two nodes. This will be parsed by the parser in a single context, thus resulting
 in a single tree.

 This allows for embedded languages to be parsed together across different
 nodes, which is useful for templating languages like ERB and EJS.

@*param* `new_regions` — List of regions this tree should manage and parse.

## source


```lua
(method) vim.treesitter.LanguageTree:source()
  -> string|integer
```

 Returns the source content of the language tree (bufnr or string).

## tree_for_range


```lua
(method) vim.treesitter.LanguageTree:tree_for_range(range: Range4, opts?: vim.treesitter.LanguageTree.tree_for_range.Opts)
  -> TSTree?
```

 Gets the tree that contains {range}.

## trees


```lua
(method) vim.treesitter.LanguageTree:trees()
  -> table<integer, TSTree>
```

 Returns all trees of the regions parsed by this parser.
 Does not include child languages.
 The result is list-like if
 * this LanguageTree is the root, in which case the result is empty or a singleton list; or
 * the root LanguageTree is fully parsed.


---

# vim.treesitter.LanguageTree.new.Opts

Optional arguments:

## injections


```lua
table<string, string>?
```

## queries


```lua
table<string, string>?
```

Deprecated


---

# vim.treesitter.LanguageTree.tree_for_range.Opts

## ignore_injections


```lua
boolean?
```

 @inlinedoc

 Ignore injected languages
 (default: `true`)


---

# vim.treesitter.ParseError


---

# vim.treesitter.Query

Parsed query, see |vim.treesitter.query.parse()|


## __index


```lua
vim.treesitter.Query
```

Parsed query, see |vim.treesitter.query.parse()|


## _apply_directives


```lua
(method) vim.treesitter.Query:_apply_directives(directives: (string|integer)[][], pattern_i: integer, captures: table<integer, TSNode[]>, source: string|integer)
  -> metadata: vim.treesitter.query.TSMetadata
```

## _match_predicates


```lua
(method) vim.treesitter.Query:_match_predicates(predicates: vim.treesitter.query.ProcessedPredicate[], pattern_i: integer, captures: table<integer, TSNode[]>, source: string|integer)
  -> whether: boolean
```

@*return* `whether` — the predicates match

## _process_patterns


```lua
(method) vim.treesitter.Query:_process_patterns()
```

 Splits the query patterns into predicates and directives.

## _processed_patterns


```lua
table<integer, vim.treesitter.query.ProcessedPattern>
```

## captures


```lua
string[]
```

list of (unique) capture names defined in query

## has_combined_injections


```lua
boolean
```

whether the query contains combined injections

## has_conceal_line


```lua
boolean
```

whether the query sets conceal_lines metadata

## info


```lua
vim.treesitter.QueryInfo
```

query context (e.g. captures, predicates, directives)

## iter_captures


```lua
(method) vim.treesitter.Query:iter_captures(node: TSNode, source: string|integer, start?: integer, stop?: integer, opts?: table)
  -> fun(end_line: integer|nil):integer, TSNode, vim.treesitter.query.TSMetadata, TSQueryMatch, TSTree
```

 Iterates over all captures from all matches in {node}.

 {source} is required if the query contains predicates; then the caller
 must ensure to use a freshly parsed tree consistent with the current
 text of the buffer (if relevant). {start} and {stop} can be used to limit
 matches inside a row range (this is typically used with root node
 as the {node}, i.e., to get syntax highlight matches in the current
 viewport). When omitted, the {start} and {stop} row values are used from the given node.

 The iterator returns four values:
 1. the numeric id identifying the capture
 2. the captured node
 3. metadata from any directives processing the match
 4. the match itself

 Example: how to get captures by name:
 ```lua
 for id, node, metadata, match in query:iter_captures(tree:root(), bufnr, first, last) do
   local name = query.captures[id] -- name of the capture in the query
   -- typically useful info about the node:
   local type = node:type() -- type of the captured node
   local row1, col1, row2, col2 = node:range() -- range of the capture
   -- ... use the info here ...
 end
 ```

@*param* `node` — under which the search will occur

@*param* `source` — Source buffer or string to extract text from

@*param* `start` — Starting line for the search. Defaults to `node:start()`.

@*param* `stop` — Stopping line for the search (end-exclusive). Defaults to `node:end_()`.

@*param* `opts` — Optional keyword arguments:

   - max_start_depth (integer) if non-zero, sets the maximum start depth
     for each match. This is used to prevent traversing too deep into a tree.
   - match_limit (integer) Set the maximum number of in-progress matches (Default: 256).

@*return* — :

        capture id, capture node, metadata, match, tree

## iter_matches


```lua
(method) vim.treesitter.Query:iter_matches(node: TSNode, source: string|integer, start?: integer, stop?: integer, opts?: table)
  -> fun():integer, table<integer, TSNode[]>, vim.treesitter.query.TSMetadata, TSTree
```

 Iterates the matches of self on a given range.

 Iterate over all matches within a {node}. The arguments are the same as for
 |Query:iter_captures()| but the iterated values are different: an (1-based)
 index of the pattern in the query, a table mapping capture indices to a list
 of nodes, and metadata from any directives processing the match.

 Example:

 ```lua
 for pattern, match, metadata in cquery:iter_matches(tree:root(), bufnr, 0, -1) do
   for id, nodes in pairs(match) do
     local name = query.captures[id]
     for _, node in ipairs(nodes) do
       -- `node` was captured by the `name` capture in the match

       local node_data = metadata[id] -- Node level metadata
       -- ... use the info here ...
     end
   end
 end
 ```


@*param* `node` — under which the search will occur

@*param* `source` — Source buffer or string to search

@*param* `start` — Starting line for the search. Defaults to `node:start()`.

@*param* `stop` — Stopping line for the search (end-exclusive). Defaults to `node:end_()`.

@*param* `opts` — Optional keyword arguments:

   - max_start_depth (integer) if non-zero, sets the maximum start depth
     for each match. This is used to prevent traversing too deep into a tree.
   - match_limit (integer) Set the maximum number of in-progress matches (Default: 256).
 - all (boolean) When `false` (default `true`), the returned table maps capture IDs to a single
   (last) node instead of the full list of matching nodes. This option is only for backward
   compatibility and will be removed in a future release.

@*return* — : pattern id, match, metadata, tree

## lang


```lua
string
```

parser language name

## new


```lua
function vim.treesitter.Query.new(lang: string, ts_query: TSQuery)
  -> vim.treesitter.Query
```

See: ~vim.treesitter.query.parse~

## query


```lua
TSQuery
```

userdata query object


---

# vim.treesitter.QueryInfo

Information for Query, see |vim.treesitter.query.parse()|

## captures


```lua
string[]
```


List of (unique) capture names defined in query.

## patterns


```lua
table<integer, (string|integer)[][]>
```


Contains information about predicates and directives.
Key is pattern id, and value is list of predicates or directives defined in the pattern.
A predicate or directive is a list of (integer|string); integer represents `capture_id`, and
string represents (literal) arguments to predicate/directive. See |treesitter-predicates|
and |treesitter-directives| for more details.


---

# vim.treesitter.dev.Injection

## lang


```lua
string
```

Source language of this injection

## root


```lua
TSNode
```

Root node of the injection


---

# vim.treesitter.dev.Node

## depth


```lua
integer
```

Depth of this node in the tree

## field


```lua
string?
```

Node field

## lang


```lua
string
```

Source language of this node

## node


```lua
TSNode
```

Treesitter node

## text


```lua
string?
```

Text displayed in the inspector for this node. Not computed until the


---

# vim.treesitter.dev.TSTreeView

## __index


```lua
vim.treesitter.dev.TSTreeView
```

## draw


```lua
(method) vim.treesitter.dev.TSTreeView:draw(bufnr: integer)
```

 Write the contents of this View into {bufnr}.

 Calling this function computes the text that is displayed for each node.

@*param* `bufnr` — Buffer number to write into.

## get


```lua
(method) vim.treesitter.dev.TSTreeView:get(i: integer)
  -> vim.treesitter.dev.Node
```

 Get node {i} from this View.

 The node number is dependent on whether or not anonymous nodes are displayed.

@*param* `i` — Node number to get

## iter


```lua
(method) vim.treesitter.dev.TSTreeView:iter()
  -> Iterator: fun():integer, vim.treesitter.dev.Node
  2. table
  3. integer
```

 Iterate over all of the nodes in this View.

@*return* `Iterator` — over all nodes in this View

@*return*

@*return*

## named


```lua
vim.treesitter.dev.Node[]
```

## new


```lua
(method) vim.treesitter.dev.TSTreeView:new(bufnr: integer, lang: string|nil)
  -> vim.treesitter.dev.TSTreeView|nil
  2. Error: string|nil
```

 Create a new treesitter view.

@*param* `bufnr` — Source buffer number

@*param* `lang` — Language of source buffer


@*return*

@*return* `Error` — message, if any

## nodes


```lua
vim.treesitter.dev.Node[]
```

## ns


```lua
integer
```

API namespace

## opts


```lua
vim.treesitter.dev.TSTreeViewOpts
```


---

# vim.treesitter.dev.TSTreeViewOpts

## anon


```lua
boolean
```

If true, display anonymous nodes.

## indent


```lua
number
```

Number of spaces to indent nested lines.

## lang


```lua
boolean
```

If true, display the language alongside each node.


---

# vim.treesitter.dev.inspect_tree.Opts

## bufnr


```lua
integer?
```


 Buffer to draw the tree into. If omitted, a new buffer is created.

## command


```lua
string?
```


 Vimscript command to create the window. Default value is "60vnew".
 Only used when {winid} is nil.

## lang


```lua
string?
```

 @inlinedoc

 The language of the source buffer. If omitted, the filetype of the source
 buffer is used.

## title


```lua
string|fun(bufnr: integer):string|nil
```


 Title of the window. If a function, it accepts the buffer number of the
 source buffer as its only argument and should return a string.

## winid


```lua
integer?
```


 Window id to display the tree buffer in. If omitted, a new window is
 created with {command}.


---

# vim.treesitter.get_node.Opts

 Optional keyword arguments:

## bufnr


```lua
integer?
```

 @inlinedoc

 Buffer number (nil or 0 for current buffer)

## ignore_injections


```lua
boolean?
```


 Ignore injected languages (default true)

## include_anonymous


```lua
boolean?
```


 Include anonymous nodes (default false)

## lang


```lua
string?
```


 Parser language. (default: from buffer filetype)

## pos


```lua
[integer, integer]?
```


 0-indexed (row, col) tuple. Defaults to cursor position in the
 current window. Required if {bufnr} is not the current buffer


---

# vim.treesitter.highlighter

## __index


```lua
vim.treesitter.highlighter
```

## _conceal_checked


```lua
table<integer, boolean>
```

## _conceal_line


```lua
boolean?
```

## _highlight_states


```lua
table<integer, vim.treesitter.highlighter.State[]>
```

 A map from window ID to highlight states.
 This state is kept during rendering across each line update.

## _on_buf


```lua
function vim.treesitter.highlighter._on_buf(_: any, buf: any)
```

 Clear conceal_lines marks whenever we redraw for a buffer change. Marks are
 added back as either the _conceal_line or on_win callback comes across them.

## _on_conceal_line


```lua
function vim.treesitter.highlighter._on_conceal_line(_: any, win: integer, buf: integer, row: integer)
```

## _on_line


```lua
function vim.treesitter.highlighter._on_line(_: any, win: integer, buf: integer, line: integer, _: any)
```

## _on_spell_nav


```lua
function vim.treesitter.highlighter._on_spell_nav(_: any, win: integer, buf: integer, srow: integer, _: any, erow: integer, _: any)
```

## _on_win


```lua
function vim.treesitter.highlighter._on_win(_: any, win: any, buf: integer, topline: integer, botline: integer)
  -> boolean
```

## _queries


```lua
table<string, vim.treesitter.highlighter.Query>
```

## active


```lua
table<integer, vim.treesitter.highlighter>
```

## bufnr


```lua
integer
```

## destroy


```lua
(method) vim.treesitter.highlighter:destroy()
```

 @nodoc
 Removes all internal references to the highlighter

## for_each_highlight_state


```lua
(method) vim.treesitter.highlighter:for_each_highlight_state(win: integer, fn: fun(state: vim.treesitter.highlighter.State))
```

## get_query


```lua
(method) vim.treesitter.highlighter:get_query(lang: string)
  -> vim.treesitter.highlighter.Query
```

 Gets the query used for @param lang

@*param* `lang` — Language used by the highlighter.

## new


```lua
function vim.treesitter.highlighter.new(tree: vim.treesitter.LanguageTree, opts: table|nil)
  -> Created: vim.treesitter.highlighter
```


 Creates a highlighter for `tree`.

@*param* `tree` — parser object to use for highlighting

@*param* `opts` — Configuration of the highlighter:

           - queries table overwrite queries used by the highlighter

@*return* `Created` — highlighter object

## on_changedtree


```lua
(method) vim.treesitter.highlighter:on_changedtree(changes: Range6[])
```

## on_detach


```lua
(method) vim.treesitter.highlighter:on_detach()
```

## orig_spelloptions


```lua
string
```

## parsing


```lua
table<integer, boolean>
```

 A map from window ID to whether we are currently parsing that window asynchronously

## prepare_highlight_states


```lua
(method) vim.treesitter.highlighter:prepare_highlight_states(win: integer, srow: integer, erow: integer)
```

@*param* `erow` — exclusive

## redraw_count


```lua
integer
```

## tree


```lua
vim.treesitter.LanguageTree
```


---

# vim.treesitter.highlighter.Iter


---

# vim.treesitter.highlighter.Query

## __index


```lua
vim.treesitter.highlighter.Query
```

## _query


```lua
(vim.treesitter.Query)?
```

Parsed query, see |vim.treesitter.query.parse()|


## get_hl_from_capture


```lua
(method) vim.treesitter.highlighter.Query:get_hl_from_capture(capture: integer)
  -> integer?
```

## hl_cache


```lua
table<integer, integer>
```

## lang


```lua
string
```

## new


```lua
function vim.treesitter.highlighter.Query.new(lang: string, query_string?: string)
  -> vim.treesitter.highlighter.Query
```

## query


```lua
(method) vim.treesitter.highlighter.Query:query()
  -> vim.treesitter.Query
```


---

# vim.treesitter.highlighter.State

## highlighter_query


```lua
vim.treesitter.highlighter.Query
```

## iter


```lua
(fun(end_line: integer|nil):integer, TSNode, vim.treesitter.query.TSMetadata, TSQueryMatch)?
```

## next_row


```lua
integer
```

## prev_marks


```lua
{ start_line: integer, start_col: integer, opts: vim.api.keyset.set_extmark }[]
```

## tstree


```lua
TSTree
```


---

# vim.treesitter.language.add.Opts

## path


```lua
string?
```


Optional path the parser is located at

## symbol_name


```lua
string?
```


Internal symbol name for the language to load


---

# vim.treesitter.languagetree.Injection


---

# vim.treesitter.languagetree.InjectionElem

## combined


```lua
boolean
```

## regions


```lua
Range6[][]
```


---

# vim.treesitter.query.ProcessedDirective


---

# vim.treesitter.query.ProcessedPattern

## directives


```lua
(string|integer)[][]
```

## predicates


```lua
vim.treesitter.query.ProcessedPredicate[]
```


---

# vim.treesitter.query.ProcessedPredicate

## [1]


```lua
string
```

predicate name

## [2]


```lua
boolean
```

should match

## [3]


```lua
(string|integer)[]
```

the original predicate


---

# vim.treesitter.query.TSMetadata

## [integer]


```lua
(vim.treesitter.query.TSMetadata)?
```

## [string]


```lua
(string|integer)?
```

## bo.commentstring


```lua
string?
```

## conceal


```lua
string?
```

## range


```lua
(Range2|Range4|Range6)?
```


---

# vim.treesitter.query.add_predicate.Opts

## all


```lua
boolean?
```


 Use the correct implementation of the match table where capture IDs map to
 a list of nodes instead of a single node. Defaults to true. This option will
 be removed in a future release.

## force


```lua
boolean?
```

 @inlinedoc

 Override an existing predicate of the same name


---

# vim.treesitter.query.lint.Opts

 Optional keyword arguments:

## clear


```lua
boolean
```


 Just clear current lint errors

## langs


```lua
(string|string[])?
```

 @inlinedoc

 Language(s) to use for checking the query.
 If multiple languages are specified, queries are validated for all of them


---

# vim.trim


```lua
function vim.trim(s: string)
  -> String: string
```


---

# vim.trust.opts

## action


```lua
'allow'|'deny'|'remove'
```

 @inlinedoc

 - `'allow'` to add a file to the trust database and trust it,
 - `'deny'` to add a file to the trust database and deny it,
 - `'remove'` to remove file from the trust database

## bufnr


```lua
integer?
```

 Buffer number to update. Mutually exclusive with {path}.

## path


```lua
string?
```


 Path to a file to update. Mutually exclusive with {bufnr}.
 Cannot be used when {action} is "allow".


---

# vim.ui


```lua
table
```


---

# vim.ui_attach


```lua
function vim.ui_attach(ns: integer, options: table<string, any>, callback: fun())
```


---

# vim.ui_detach


```lua
function vim.ui_detach(ns: integer)
```


---

# vim.uri_from_bufnr


```lua
function
```


---

# vim.uri_from_fname


```lua
function
```


---

# vim.uri_to_bufnr


```lua
function
```


---

# vim.uri_to_fname


```lua
function
```


---

# vim.uv


```lua
unknown
```


---

# vim.v


```lua
vim.v
```


```lua
unknown
```


---

# vim.v

## argv


```lua
unknown
```

 The command line arguments Vim was invoked with.  This is a
 list of strings.  The first item is the Vim command.
 See `v:progpath` for the command with full path.

## char


```lua
unknown
```

 Argument for evaluating 'formatexpr' and used for the typed
 character when using <expr> in an abbreviation `:map-<expr>`.
 It is also used by the `InsertCharPre` and `InsertEnter` events.

## charconvert_from


```lua
unknown
```

 The name of the character encoding of a file to be converted.
 Only valid while evaluating the 'charconvert' option.

## charconvert_to


```lua
unknown
```

 The name of the character encoding of a file after conversion.
 Only valid while evaluating the 'charconvert' option.

## cmdarg


```lua
unknown
```

 The extra arguments ("++p", "++enc=", "++ff=") given to a file
 read/write command.  This is set before an autocommand event
 for a file read/write command is triggered.  There is a
 leading space to make it possible to append this variable
 directly after the read/write command. Note: "+cmd" isn't
 included here, because it will be executed anyway.

## cmdbang


```lua
unknown
```

 Set like v:cmdarg for a file read/write command.  When a "!"
 was used the value is 1, otherwise it is 0.  Note that this
 can only be used in autocommands.  For user commands `<bang>`
 can be used.

## collate


```lua
unknown
```

 The current locale setting for collation order of the runtime
 environment.  This allows Vim scripts to be aware of the
 current locale encoding.  Technical: it's the value of
 LC_COLLATE.  When not using a locale the value is "C".
 This variable can not be set directly, use the `:language`
 command.
 See `multi-lang`.

## completed_item


```lua
unknown
```

 Dictionary containing the `complete-items` for the most
 recently completed word after `CompleteDone`.  Empty if the
 completion failed, or after leaving and re-entering insert
 mode.
 Note: Plugins can modify the value to emulate the builtin
 `CompleteDone` event behavior.

## count


```lua
unknown
```

 The count given for the last Normal mode command.  Can be used
 to get the count before a mapping.  Read-only.  Example:

 ```vim
   :map _x :<C-U>echo "the count is " .. v:count<CR>
 ```

 Note: The <C-U> is required to remove the line range that you
 get when typing ':' after a count.
 When there are two counts, as in "3d2w", they are multiplied,
 just like what happens in the command, "d6w" for the example.
 Also used for evaluating the 'formatexpr' option.

## count1


```lua
unknown
```

 Just like "v:count", but defaults to one when no count is
 used.

## ctype


```lua
unknown
```

 The current locale setting for characters of the runtime
 environment.  This allows Vim scripts to be aware of the
 current locale encoding.  Technical: it's the value of
 LC_CTYPE.  When not using a locale the value is "C".
 This variable can not be set directly, use the `:language`
 command.
 See `multi-lang`.

## dying


```lua
unknown
```

 Normally zero.  When a deadly signal is caught it's set to
 one.  When multiple signals are caught the number increases.
 Can be used in an autocommand to check if Vim didn't
 terminate normally.
 Example:

 ```vim
   :au VimLeave * if v:dying | echo "\nAAAAaaaarrrggghhhh!!!\n" | endif
 ```

 Note: if another deadly signal is caught when v:dying is one,
 VimLeave autocommands will not be executed.

## echospace


```lua
unknown
```

 Number of screen cells that can be used for an `:echo` message
 in the last screen line before causing the `hit-enter-prompt`.
 Depends on 'showcmd', 'ruler' and 'columns'.  You need to
 check 'cmdheight' for whether there are full-width lines
 available above the last line.

## errmsg


```lua
unknown
```

 Last given error message.
 Modifiable (can be set).
 Example:

 ```vim
   let v:errmsg = ""
   silent! next
   if v:errmsg != ""
     " ... handle error
 ```

## errors


```lua
unknown
```

 Errors found by assert functions, such as `assert_true()`.
 This is a list of strings.
 The assert functions append an item when an assert fails.
 The return value indicates this: a one is returned if an item
 was added to v:errors, otherwise zero is returned.
 To remove old results make it empty:

 ```vim
   let v:errors = []
 ```

 If v:errors is set to anything but a list it is made an empty
 list by the assert function.

## event


```lua
unknown
```

 Dictionary of event data for the current `autocommand`.  Valid
 only during the event lifetime; storing or passing v:event is
 invalid!  Copy it instead:

 ```vim
   au TextYankPost * let g:foo = deepcopy(v:event)
 ```

 Keys vary by event; see the documentation for the specific
 event, e.g. `DirChanged` or `TextYankPost`.
   KEY              DESCRIPTION ~
   abort            Whether the event triggered during
                    an aborting condition (e.g. `c_Esc` or
                    `c_CTRL-C` for `CmdlineLeave`).
   chan             `channel-id`
   changed_window   Is `v:true` if the event fired while
                    changing window  (or tab) on `DirChanged`.
   cmdlevel         Level of cmdline.
   cmdtype          Type of cmdline, `cmdline-char`.
   col              Column count of popup menu on `CompleteChanged`,
                    relative to screen.
   complete_type    See `complete_info_mode`
   complete_word    The selected word, or empty if completion
                    was abandoned/discarded.
   completed_item   Current selected item on `CompleteChanged`,
                    or `{}` if no item selected.
   cwd              Current working directory.
   height           Height of popup menu on `CompleteChanged`
   inclusive        Motion is `inclusive`, else exclusive.
   info             Dict of arbitrary event data.
   operator         Current `operator`.  Also set for Ex
                    commands (unlike `v:operator`). For
                    example if `TextYankPost` is triggered
                    by the `:yank` Ex command then
                    `v:event.operator` is "y".
   reason           `CompleteDone` reason.
   regcontents      Text stored in the register as a
                    `readfile()`-style list of lines.
   regname          Requested register (e.g "x" for "xyy), or
                    empty string for an unnamed operation.
   regtype          Type of register as returned by
                    `getregtype()`.
   row              Row count of popup menu on `CompleteChanged`,
                    relative to screen.
   scope            Event-specific scope name.
   scrollbar        `v:true` if popup menu has a scrollbar, or
                    `v:false` if not.
   size             Total number of completion items on
                    `CompleteChanged`.
   status           Job status or exit code, -1 means "unknown". `TermClose`
   visual           Selection is visual (as opposed to e.g. a motion range).
   width            Width of popup menu on `CompleteChanged`
   windows          List of window IDs that changed on `WinResized`

## exception


```lua
unknown
```

 The value of the exception most recently caught and not
 finished.  See also `v:stacktrace`, `v:throwpoint`, and
 `throw-variables`.
 Example:

 ```vim
   try
     throw "oops"
   catch /.*/
     echo "caught " .. v:exception
   endtry
 ```

 Output: "caught oops".

## exiting


```lua
unknown
```

 Exit code, or `v:null` before invoking the `VimLeavePre`
 and `VimLeave` autocmds.  See `:q`, `:x` and `:cquit`.
 Example:

 ```vim
   :au VimLeave * echo "Exit value is " .. v:exiting
 ```

## fcs_choice


```lua
unknown
```

 What should happen after a `FileChangedShell` event was
 triggered.  Can be used in an autocommand to tell Vim what to
 do with the affected buffer:
   reload  Reload the buffer (does not work if
           the file was deleted).
   edit    Reload the buffer and detect the
           values for options such as
           'fileformat', 'fileencoding', 'binary'
           (does not work if the file was
           deleted).
   ask     Ask the user what to do, as if there
           was no autocommand.  Except that when
           only the timestamp changed nothing
           will happen.
   <empty> Nothing, the autocommand should do
           everything that needs to be done.
 The default is empty.  If another (invalid) value is used then
 Vim behaves like it is empty, there is no warning message.

## fcs_reason


```lua
unknown
```

 The reason why the `FileChangedShell` event was triggered.
 Can be used in an autocommand to decide what to do and/or what
 to set v:fcs_choice to.  Possible values:
   deleted   file no longer exists
   conflict  file contents, mode or timestamp was
             changed and buffer is modified
   changed   file contents has changed
   mode      mode of file changed
   time      only file timestamp changed

## fname


```lua
unknown
```

 When evaluating 'includeexpr': the file name that was
 detected.  Empty otherwise.

## fname_diff


```lua
unknown
```

 The name of the diff (patch) file.  Only valid while
 evaluating 'patchexpr'.

## fname_in


```lua
unknown
```

 The name of the input file.  Valid while evaluating:
   option         used for ~
   'charconvert'  file to be converted
   'diffexpr'     original file
   'patchexpr'    original file
 And set to the swap file name for `SwapExists`.

## fname_new


```lua
unknown
```

 The name of the new version of the file.  Only valid while
 evaluating 'diffexpr'.

## fname_out


```lua
unknown
```

 The name of the output file.  Only valid while
 evaluating:
   option           used for ~
   'charconvert'    resulting converted file [1]
   'diffexpr'       output of diff
   'patchexpr'      resulting patched file
 [1] When doing conversion for a write command (e.g., ":w
 file") it will be equal to v:fname_in.  When doing conversion
 for a read command (e.g., ":e file") it will be a temporary
 file and different from v:fname_in.

## folddashes


```lua
unknown
```

 Used for 'foldtext': dashes representing foldlevel of a closed
 fold.
 Read-only in the `sandbox`. `fold-foldtext`

## foldend


```lua
unknown
```

 Used for 'foldtext': last line of closed fold.
 Read-only in the `sandbox`. `fold-foldtext`

## foldlevel


```lua
unknown
```

 Used for 'foldtext': foldlevel of closed fold.
 Read-only in the `sandbox`. `fold-foldtext`

## foldstart


```lua
unknown
```

 Used for 'foldtext': first line of closed fold.
 Read-only in the `sandbox`. `fold-foldtext`

## hlsearch


```lua
unknown
```

 Variable that indicates whether search highlighting is on.
 Setting it makes sense only if 'hlsearch' is enabled. Setting
 this variable to zero acts like the `:nohlsearch` command,
 setting it to one acts like

 ```vim
   let &hlsearch = &hlsearch
 ```

 Note that the value is restored when returning from a
 function. `function-search-undo`.

## insertmode


```lua
unknown
```

 Used for the `InsertEnter` and `InsertChange` autocommand
 events.  Values:
   i    Insert mode
   r    Replace mode
   v    Virtual Replace mode

## key


```lua
unknown
```

 Key of the current item of a `Dictionary`.  Only valid while
 evaluating the expression used with `map()` and `filter()`.
 Read-only.

## lang


```lua
unknown
```

 The current locale setting for messages of the runtime
 environment.  This allows Vim scripts to be aware of the
 current language.  Technical: it's the value of LC_MESSAGES.
 The value is system dependent.
 This variable can not be set directly, use the `:language`
 command.
 It can be different from `v:ctype` when messages are desired
 in a different language than what is used for character
 encoding.  See `multi-lang`.

## lc_time


```lua
unknown
```

 The current locale setting for time messages of the runtime
 environment.  This allows Vim scripts to be aware of the
 current language.  Technical: it's the value of LC_TIME.
 This variable can not be set directly, use the `:language`
 command.  See `multi-lang`.

## lnum


```lua
unknown
```

 Line number for the 'foldexpr' `fold-expr`, 'formatexpr',
 'indentexpr' and 'statuscolumn' expressions, tab page number
 for 'guitablabel' and 'guitabtooltip'.  Only valid while one of
 these expressions is being evaluated.  Read-only when in the
 `sandbox`.

## lua


```lua
unknown
```

 Prefix for calling Lua functions from expressions.
 See `v:lua-call` for more information.

## maxcol


```lua
unknown
```

 Maximum line length.  Depending on where it is used it can be
 screen columns, characters or bytes.  The value currently is
 2147483647 on all systems.

## mouse_col


```lua
unknown
```

 Column number for a mouse click obtained with `getchar()`.
 This is the screen column number, like with `virtcol()`.  The
 value is zero when there was no mouse button click.

## mouse_lnum


```lua
unknown
```

 Line number for a mouse click obtained with `getchar()`.
 This is the text line number, not the screen line number.  The
 value is zero when there was no mouse button click.

## mouse_win


```lua
unknown
```

 Window number for a mouse click obtained with `getchar()`.
 First window has number 1, like with `winnr()`.  The value is
 zero when there was no mouse button click.

## mouse_winid


```lua
unknown
```

 `window-ID` for a mouse click obtained with `getchar()`.
 The value is zero when there was no mouse button click.

## msgpack_types


```lua
unknown
```

 Dictionary containing msgpack types used by `msgpackparse()`
 and `msgpackdump()`. All types inside dictionary are fixed
 (not editable) empty lists. To check whether some list is one
 of msgpack types, use `is` operator.

## null


```lua
unknown
```

 Special value used to put "null" in JSON and NIL in msgpack.
 See `json_encode()`.  This value is converted to "v:null" when
 used as a String (e.g. in `expr5` with string concatenation
 operator) and to zero when used as a Number (e.g. in `expr5`
 or `expr7` when used with numeric operators). Read-only.
 In some places `v:null` can be used for a List, Dict, etc.
 that is not set.  That is slightly different than an empty
 List, Dict, etc.

## numbermax


```lua
unknown
```

 Maximum value of a number.

## numbermin


```lua
unknown
```

 Minimum value of a number (negative).

## numbersize


```lua
unknown
```

 Number of bits in a Number.  This is normally 64, but on some
 systems it may be 32.

## oldfiles


```lua
unknown
```

 List of file names that is loaded from the `shada` file on
 startup.  These are the files that Vim remembers marks for.
 The length of the List is limited by the ' argument of the
 'shada' option (default is 100).
 When the `shada` file is not used the List is empty.
 Also see `:oldfiles` and `c_#<`.
 The List can be modified, but this has no effect on what is
 stored in the `shada` file later.  If you use values other
 than String this will cause trouble.

## operator


```lua
unknown
```

 The last operator given in Normal mode.  This is a single
 character except for commands starting with <g> or <z>,
 in which case it is two characters.  Best used alongside
 `v:prevcount` and `v:register`.  Useful if you want to cancel
 Operator-pending mode and then use the operator, e.g.:

 ```vim
   :omap O <Esc>:call MyMotion(v:operator)<CR>
 ```

 The value remains set until another operator is entered, thus
 don't expect it to be empty.
 v:operator is not set for `:delete`, `:yank` or other Ex
 commands.
 Read-only.

## option_command


```lua
unknown
```

 Command used to set the option. Valid while executing an
 `OptionSet` autocommand.
   value        option was set via ~
   "setlocal"   `:setlocal` or `:let l:xxx`
   "setglobal"  `:setglobal` or `:let g:xxx`
   "set"        `:set` or `:let`
   "modeline"   `modeline`

## option_new


```lua
unknown
```

 New value of the option. Valid while executing an `OptionSet`
 autocommand.

## option_old


```lua
unknown
```

 Old value of the option. Valid while executing an `OptionSet`
 autocommand. Depending on the command used for setting and the
 kind of option this is either the local old value or the
 global old value.

## option_oldglobal


```lua
unknown
```

 Old global value of the option. Valid while executing an
 `OptionSet` autocommand.

## option_oldlocal


```lua
unknown
```

 Old local value of the option. Valid while executing an
 `OptionSet` autocommand.

## option_type


```lua
unknown
```

 Scope of the set command. Valid while executing an
 `OptionSet` autocommand. Can be either "global" or "local"

## prevcount


```lua
unknown
```

 The count given for the last but one Normal mode command.
 This is the v:count value of the previous command.  Useful if
 you want to cancel Visual or Operator-pending mode and then
 use the count, e.g.:

 ```vim
   :vmap % <Esc>:call MyFilter(v:prevcount)<CR>
 ```

 Read-only.

## profiling


```lua
unknown
```

 Normally zero.  Set to one after using ":profile start".
 See `profiling`.

## progname


```lua
unknown
```

 The name by which Nvim was invoked (with path removed).
 Read-only.

## progpath


```lua
unknown
```

 Absolute path to the current running Nvim.
 Read-only.

## register


```lua
unknown
```

 The name of the register in effect for the current normal mode
 command (regardless of whether that command actually used a
 register).  Or for the currently executing normal mode mapping
 (use this in custom commands that take a register).
 If none is supplied it is the default register '"', unless
 'clipboard' contains "unnamed" or "unnamedplus", then it is
 "*" or '+'.
 Also see `getreg()` and `setreg()`

## relnum


```lua
unknown
```

 Relative line number for the 'statuscolumn' expression.
 Read-only.

## scrollstart


```lua
unknown
```

 String describing the script or function that caused the
 screen to scroll up.  It's only set when it is empty, thus the
 first reason is remembered.  It is set to "Unknown" for a
 typed command.
 This can be used to find out why your script causes the
 hit-enter prompt.

## searchforward


```lua
any
```

## servername


```lua
unknown
```

 Primary listen-address of Nvim, the first item returned by
 `serverlist()`. Usually this is the named pipe created by Nvim
 at `startup` or given by `--listen` (or the deprecated
 `$NVIM_LISTEN_ADDRESS` env var).

 See also `serverstart()` `serverstop()`.
 Read-only.

                                                      *$NVIM*
 $NVIM is set to v:servername by `terminal` and `jobstart()`,
 and is thus a hint that the current environment is a child
 (direct subprocess) of Nvim.

 Example: a child Nvim process can detect and make requests to
 its parent Nvim:

 ```lua

   if vim.env.NVIM then
     local ok, chan = pcall(vim.fn.sockconnect, 'pipe', vim.env.NVIM, {rpc=true})
     if ok and chan then
       local client = vim.api.nvim_get_chan_info(chan).client
       local rv = vim.rpcrequest(chan, 'nvim_exec_lua', [[return ... + 1]], { 41 })
       vim.print(('got "%s" from parent Nvim'):format(rv))
     end
   end
 ```

## shell_error


```lua
unknown
```

 Result of the last shell command.  When non-zero, the last
 shell command had an error.  When zero, there was no problem.
 This only works when the shell returns the error code to Vim.
 The value -1 is often used when the command could not be
 executed.  Read-only.
 Example:

 ```vim
   !mv foo bar
   if v:shell_error
     echo 'could not rename "foo" to "bar"!'
   endif
 ```

## stacktrace


```lua
unknown
```

 The stack trace of the exception most recently caught and
 not finished.  Refer to `getstacktrace()` for the structure of
 stack trace.  See also `v:exception`, `v:throwpoint`, and
 `throw-variables`.

## statusmsg


```lua
unknown
```

 Last given status message.
 Modifiable (can be set).

## stderr


```lua
unknown
```

 `channel-id` corresponding to stderr. The value is always 2;
 use this variable to make your code more descriptive.
 Unlike stdin and stdout (see `stdioopen()`), stderr is always
 open for writing. Example:

 ```vim
 :call chansend(v:stderr, "error: toaster empty\n")
 ```

## swapchoice


```lua
string
```

 Show the prompt.

## swapcommand


```lua
unknown
```

 Normal mode command to be executed after a file has been
 opened.  Can be used for a `SwapExists` autocommand to have
 another Vim open the file and jump to the right place.  For
 example, when jumping to a tag the value is ":tag tagname\r".
 For ":edit +cmd file" the value is ":cmd\r".

## swapname


```lua
unknown
```

 Name of the swapfile found.
 Only valid during `SwapExists` event.
 Read-only.

## t_blob


```lua
unknown
```

 Value of `Blob` type.  Read-only.  See: `type()`

## t_bool


```lua
unknown
```

 Value of `Boolean` type.  Read-only.  See: `type()`

## t_dict


```lua
unknown
```

 Value of `Dictionary` type.  Read-only.  See: `type()`

## t_float


```lua
unknown
```

 Value of `Float` type.  Read-only.  See: `type()`

## t_func


```lua
unknown
```

 Value of `Funcref` type.  Read-only.  See: `type()`

## t_list


```lua
unknown
```

 Value of `List` type.  Read-only.  See: `type()`

## t_number


```lua
unknown
```

 Value of `Number` type.  Read-only.  See: `type()`

## t_string


```lua
unknown
```

 Value of `String` type.  Read-only.  See: `type()`

## termrequest


```lua
unknown
```

 The value of the most recent OSC, DCS or APC control sequence
 sent from a process running in the embedded `terminal`.
 This can be read in a `TermRequest` event handler to respond
 to queries from embedded applications.

## termresponse


```lua
unknown
```

 The value of the most recent OSC or DCS control sequence
 received by Nvim from the terminal. This can be read in a
 `TermResponse` event handler after querying the terminal using
 another escape sequence.

## testing


```lua
unknown
```

 Must be set before using `test_garbagecollect_now()`.

## this_session


```lua
unknown
```

 Full filename of the last loaded or saved session file.
 Empty when no session file has been saved.  See `:mksession`.
 Modifiable (can be set).

## throwpoint


```lua
unknown
```

 The point where the exception most recently caught and not
 finished was thrown.  Not set when commands are typed.  See
 also `v:exception`, `v:stacktrace`, and `throw-variables`.
 Example:

 ```vim
   try
     throw "oops"
   catch /.*/
     echo "Exception from" v:throwpoint
   endtry
 ```

 Output: "Exception from test.vim, line 2"

## val


```lua
unknown
```

 Value of the current item of a `List` or `Dictionary`.  Only
 valid while evaluating the expression used with `map()` and
 `filter()`.  Read-only.

## version


```lua
unknown
```

 Vim version number: major version times 100 plus minor
 version.  Vim 5.0 is 500, Vim 5.1 is 501.
 Read-only.
 Use `has()` to check the Nvim (not Vim) version:

 ```vim
   :if has("nvim-0.2.1")
 ```

## vim_did_enter


```lua
unknown
```

 0 during startup, 1 just before `VimEnter`.
 Read-only.

## virtnum


```lua
unknown
```

 Virtual line number for the 'statuscolumn' expression.
 Negative when drawing the status column for virtual lines, zero
 when drawing an actual buffer line, and positive when drawing
 the wrapped part of a buffer line.
 Read-only.

## warningmsg


```lua
unknown
```

 Last given warning message.
 Modifiable (can be set).

## windowid


```lua
unknown
```

 Application-specific window "handle" which may be set by any
 attached UI. Defaults to zero.
 Note: For Nvim `windows` use `winnr()` or `win_getid()`, see
 `window-ID`.


---

# vim.v.argv


```lua
unknown
```


---

# vim.v.char


```lua
unknown
```


---

# vim.v.charconvert_from


```lua
unknown
```


---

# vim.v.charconvert_to


```lua
unknown
```


---

# vim.v.cmdarg


```lua
unknown
```


---

# vim.v.cmdbang


```lua
unknown
```


---

# vim.v.collate


```lua
unknown
```


---

# vim.v.completed_item

## abbr


```lua
string?
```

 abbreviation of "word"; when not empty it is used in the menu instead of "word"

## abbr_hlgroup


```lua
string?
```

 an additional highlight group whose attributes are combined
 with |hl-PmenuSel| and |hl-Pmenu| or |hl-PmenuMatchSel| and |hl-PmenuMatch|
 highlight attributes in the popup menu to apply cterm and gui properties
 (with higher priority) like strikethrough to the completion items abbreviation

## dup


```lua
integer?
```

 when non-zero this match will be added even when an item with the same word
 is already present.

## empty


```lua
integer?
```

 when non-zero this match will be added even when it is an empty string

## equal


```lua
integer?
```

 when non-zero, always treat this item to be equal when comparing. Which
 means, "equal=1" disables filtering of this item.

## icase


```lua
integer?
```

 when non-zero case is to be ignored when comparing items to be equal; when
 omitted zero is used, thus items that only differ in case are added

## info


```lua
string?
```

 more information about the item, can be displayed in a preview window

## kind


```lua
string?
```

single letter indicating the type of completion

## kind_hlgroup


```lua
string?
```

 an additional highlight group specifically for setting the highlight
 attributes of the completion kind. When this field is present, it will
 override the |hl-PmenuKind| highlight group, allowing for the customization
 of ctermfg and guifg properties for the completion kind

## menu


```lua
string?
```

 extra text for the popup menu, displayed after "word" or "abbr"

## user_data


```lua
any
```

 custom data which is associated with the item and available
 in |v:completed_item|; it can be any type; defaults to an empty string

## word


```lua
string?
```

the text that will be inserted, mandatory


---

# vim.v.completed_item


```lua
unknown
```


---

# vim.v.count


```lua
unknown
```


---

# vim.v.count1


```lua
unknown
```


---

# vim.v.ctype


```lua
unknown
```


---

# vim.v.dying


```lua
unknown
```


---

# vim.v.echospace


```lua
unknown
```


---

# vim.v.errmsg


```lua
unknown
```


---

# vim.v.errors


```lua
unknown
```


---

# vim.v.event


```lua
unknown
```


---

# vim.v.event

## abort


```lua
boolean?
```

 Whether the event triggered during an aborting condition (e.g. |c_Esc| or
 |c_CTRL-C| for |CmdlineLeave|).

## chan


```lua
integer?
```

See |channel-id|

## changed_window


```lua
boolean?
```

 Is |v:true| if the event fired while changing window  (or tab) on |DirChanged|.

## cmdlevel


```lua
integer?
```

Level of cmdline.

## cmdtype


```lua
string?
```

Type of cmdline, |cmdline-char|.

## col


```lua
integer?
```

 Col count of popup menu on |CompleteChanged|, relative to screen.

## complete_type


```lua
string?
```

See |complete_info_mode|

## completed_item


```lua
(vim.v.completed_item)?
```

## cwd


```lua
string?
```

Current working directory.

## height


```lua
integer?
```

 Current selected complete item on |CompleteChanged|, Is `{}` when no
 complete item selected.

## inclusive


```lua
boolean?
```

Motion is |inclusive|, else exclusive.

## info


```lua
table?
```

Dict of arbitrary event data.

## operator


```lua
string?
```

 Current |operator|. Also set for Ex commands (unlike |v:operator|). For
 example if |TextYankPost| is triggered by the |:yank| Ex command then
 `v:event.operator` is "y".

## reason


```lua
string?
```

Reason for completion being done. |CompleteDone|

## regcontents


```lua
string?
```

 Text stored in the register as a |readfile()|-style list of lines.

## regname


```lua
string?
```

 Requested register (e.g "x" for "xyy) or the empty string for an unnamed operation.

## regtype


```lua
string?
```

Type of register as returned by |getregtype()|.

## row


```lua
integer?
```

Width of popup menu on |CompleteChanged|

## scope


```lua
string?
```

Event-specific scope name.

## scrollbar


```lua
boolean?
```

 Is |v:true| if popup menu have scrollbar, or |v:false| if not.

## size


```lua
integer?
```

Total number of completion items on |CompleteChanged|.

## status


```lua
boolean?
```

Job status or exit code, -1 means "unknown". |TermClose|

## visual


```lua
boolean?
```

Selection is visual (as opposed to, e.g., via motion).

## width


```lua
integer?
```

Height of popup menu on |CompleteChanged|

## windows


```lua
integer[]?
```

 List of window IDs that changed on |WinResized|


---

# vim.v.exception


```lua
unknown
```


---

# vim.v.exiting


```lua
unknown
```


---

# vim.v.false


```lua
unknown
```


---

# vim.v.fcs_choice


```lua
unknown
```


---

# vim.v.fcs_reason


```lua
unknown
```


---

# vim.v.fname


```lua
unknown
```


---

# vim.v.fname_diff


```lua
unknown
```


---

# vim.v.fname_in


```lua
unknown
```


---

# vim.v.fname_new


```lua
unknown
```


---

# vim.v.fname_out


```lua
unknown
```


---

# vim.v.folddashes


```lua
unknown
```


---

# vim.v.foldend


```lua
unknown
```


---

# vim.v.foldlevel


```lua
unknown
```


---

# vim.v.foldstart


```lua
unknown
```


---

# vim.v.hlsearch


```lua
unknown
```


---

# vim.v.insertmode


```lua
unknown
```


---

# vim.v.key


```lua
unknown
```


---

# vim.v.lang


```lua
unknown
```


---

# vim.v.lc_time


```lua
unknown
```


---

# vim.v.lnum


```lua
unknown
```


---

# vim.v.lua


```lua
unknown
```


---

# vim.v.maxcol


```lua
unknown
```


---

# vim.v.mouse_col


```lua
unknown
```


---

# vim.v.mouse_lnum


```lua
unknown
```


---

# vim.v.mouse_win


```lua
unknown
```


---

# vim.v.mouse_winid


```lua
unknown
```


---

# vim.v.msgpack_types


```lua
unknown
```


---

# vim.v.null


```lua
unknown
```


---

# vim.v.numbermax


```lua
unknown
```


---

# vim.v.numbermin


```lua
unknown
```


---

# vim.v.numbersize


```lua
unknown
```


---

# vim.v.oldfiles


```lua
unknown
```


---

# vim.v.operator


```lua
unknown
```


---

# vim.v.option_command


```lua
unknown
```


---

# vim.v.option_new


```lua
unknown
```


---

# vim.v.option_old


```lua
unknown
```


---

# vim.v.option_oldglobal


```lua
unknown
```


---

# vim.v.option_oldlocal


```lua
unknown
```


---

# vim.v.option_type


```lua
unknown
```


---

# vim.v.prevcount


```lua
unknown
```


---

# vim.v.profiling


```lua
unknown
```


---

# vim.v.progname


```lua
unknown
```


---

# vim.v.progpath


```lua
unknown
```


---

# vim.v.register


```lua
unknown
```


---

# vim.v.relnum


```lua
unknown
```


---

# vim.v.scrollstart


```lua
unknown
```


---

# vim.v.searchforward


```lua
any
```


```lua
unknown
```


---

# vim.v.servername


```lua
unknown
```


---

# vim.v.shell_error


```lua
unknown
```


---

# vim.v.stacktrace


```lua
unknown
```


---

# vim.v.statusmsg


```lua
unknown
```


---

# vim.v.stderr


```lua
unknown
```


---

# vim.v.swapchoice


```lua
string
```


```lua
string
```


```lua
unknown
```


---

# vim.v.swapcommand


```lua
unknown
```


---

# vim.v.swapname


```lua
unknown
```


---

# vim.v.t_blob


```lua
unknown
```


---

# vim.v.t_bool


```lua
unknown
```


---

# vim.v.t_dict


```lua
unknown
```


---

# vim.v.t_float


```lua
unknown
```


---

# vim.v.t_func


```lua
unknown
```


---

# vim.v.t_list


```lua
unknown
```


---

# vim.v.t_number


```lua
unknown
```


---

# vim.v.t_string


```lua
unknown
```


---

# vim.v.termrequest


```lua
unknown
```


---

# vim.v.termresponse


```lua
unknown
```


---

# vim.v.testing


```lua
unknown
```


---

# vim.v.this_session


```lua
unknown
```


---

# vim.v.throwpoint


```lua
unknown
```


---

# vim.v.true


```lua
unknown
```


---

# vim.v.val


```lua
unknown
```


---

# vim.v.version


```lua
unknown
```


---

# vim.v.vim_did_enter


```lua
unknown
```


---

# vim.v.virtnum


```lua
unknown
```


---

# vim.v.warningmsg


```lua
unknown
```


---

# vim.v.windowid


```lua
unknown
```


---

# vim.validate


```lua
function vim.validate(name: string, value: any, validator: "boolean"|"function"|"nil"|"number"|"string"...(+6), optional?: boolean, message?: string)
```


---

# vim.validate.Spec

 @nodoc

## [1]


```lua
any
```

Argument value

## [2]


```lua
"boolean"|"function"|"nil"|"number"|"string"...(+6)
```

Argument validator

## [3]


```lua
(boolean|string)?
```

Optional flag or error message


---

# vim.validate.Validator


---

# vim.var_accessor

## [integer]


```lua
vim.var_accessor
```

## [string]


```lua
any
```


---

# vim.version


```lua
table
```


---

# vim.w


```lua
vim.var_accessor
```


---

# vim.w.qf_toc


```lua
string
```


```lua
string
```


---

# vim.wait


```lua
function vim.wait(time: integer, callback?: fun():boolean, interval?: integer, fast_only?: boolean)
  -> boolean
  2. -1|-2|nil
```


---

# vim.wo


```lua
table|vim.wo
```


```lua
table
```


---

# vim.wo

## [integer]


```lua
vim.wo
```

## arab


```lua
boolean
```

## arabic


```lua
boolean
```

## breakindent


```lua
boolean
```

## breakindentopt


```lua
string
```

## bri


```lua
boolean
```

## briopt


```lua
string
```

## cc


```lua
string
```

## cocu


```lua
string
```

## cole


```lua
integer
```

## colorcolumn


```lua
string
```

## concealcursor


```lua
string
```

## conceallevel


```lua
integer
```

## crb


```lua
boolean
```

## cuc


```lua
boolean
```

## cul


```lua
boolean
```

## culopt


```lua
string
```

## cursorbind


```lua
boolean
```

## cursorcolumn


```lua
boolean
```

## cursorline


```lua
boolean
```

## cursorlineopt


```lua
string
```

## diff


```lua
boolean
```

## eiw


```lua
string
```

## eventignorewin


```lua
string
```

## fcs


```lua
string
```

## fdc


```lua
'0'|'1'|'2'|'3'|'4'...(+15)
```

## fde


```lua
string
```

## fdi


```lua
string
```

## fdl


```lua
integer
```

## fdm


```lua
'diff'|'expr'|'indent'|'manual'|'marker'...(+1)
```

## fdn


```lua
integer
```

## fdt


```lua
string
```

## fen


```lua
boolean
```

## fillchars


```lua
string
```

## fml


```lua
integer
```

## fmr


```lua
string
```

## foldcolumn


```lua
'0'|'1'|'2'|'3'|'4'...(+15)
```

## foldenable


```lua
boolean
```

## foldexpr


```lua
string
```

## foldignore


```lua
string
```

## foldlevel


```lua
integer
```

## foldmarker


```lua
string
```

## foldmethod


```lua
'diff'|'expr'|'indent'|'manual'|'marker'...(+1)
```

## foldminlines


```lua
integer
```

## foldnestmax


```lua
integer
```

## foldtext


```lua
string
```

## lbr


```lua
boolean
```

## lcs


```lua
string
```

## linebreak


```lua
boolean
```

## list


```lua
boolean
```

## listchars


```lua
string
```

## nu


```lua
boolean
```

## number


```lua
boolean
```

## numberwidth


```lua
integer
```

## nuw


```lua
integer
```

## previewwindow


```lua
boolean
```

## pvw


```lua
boolean
```

## relativenumber


```lua
boolean
```

## rightleft


```lua
boolean
```

## rightleftcmd


```lua
string
```

## rl


```lua
boolean
```

## rlc


```lua
string
```

## rnu


```lua
boolean
```

## sbr


```lua
string
```

## scb


```lua
boolean
```

## scl


```lua
'auto'|'auto:1'|'auto:2'|'auto:3'|'auto:4'...(+17)
```

## scr


```lua
integer
```

## scroll


```lua
integer
```

## scrollbind


```lua
boolean
```

## scrolloff


```lua
integer
```

## showbreak


```lua
string
```

## sidescrolloff


```lua
integer
```

## signcolumn


```lua
'auto'|'auto:1'|'auto:2'|'auto:3'|'auto:4'...(+17)
```

## siso


```lua
integer
```

## smoothscroll


```lua
boolean
```

## sms


```lua
boolean
```

## so


```lua
integer
```

## spell


```lua
boolean
```

## statuscolumn


```lua
string
```

## statusline


```lua
string
```

## stc


```lua
string
```

## stl


```lua
string
```

## ve


```lua
string
```

## virtualedit


```lua
string
```

## wbr


```lua
string
```

## wfb


```lua
boolean
```

## wfh


```lua
boolean
```

## wfw


```lua
boolean
```

## winbar


```lua
string
```

## winbl


```lua
integer
```

## winblend


```lua
integer
```

## winfixbuf


```lua
boolean
```

## winfixheight


```lua
boolean
```

## winfixwidth


```lua
boolean
```

## winhighlight


```lua
string
```

## winhl


```lua
string
```

## wrap


```lua
boolean
```


---

# vim.wo.0.0.foldcolumn


```lua
string
```


---

# vim.wo.0.0.foldexpr


```lua
string
```


---

# vim.wo.0.0.list


```lua
boolean
```


---

# vim.wo.0.0.number


```lua
boolean
```


---

# vim.wo.0.0.relativenumber


```lua
boolean
```


---

# vim.wo.0.0.signcolumn


```lua
string
```


---

# vim.wo.0.0.winhighlight


```lua
string
```


---

# vim.wo.0.0.wrap


```lua
boolean
```


---

# vim.wo.arab


```lua
boolean
```


---

# vim.wo.arabic


```lua
boolean
```


---

# vim.wo.breakindent


```lua
boolean
```


---

# vim.wo.breakindentopt


```lua
string
```


---

# vim.wo.bri


```lua
boolean
```


---

# vim.wo.briopt


```lua
string
```


---

# vim.wo.cc


```lua
string
```


---

# vim.wo.cocu


```lua
string
```


---

# vim.wo.cole


```lua
integer
```


---

# vim.wo.colorcolumn


```lua
string
```


---

# vim.wo.concealcursor


```lua
string
```


---

# vim.wo.conceallevel


```lua
integer
```


---

# vim.wo.crb


```lua
boolean
```


---

# vim.wo.cuc


```lua
boolean
```


---

# vim.wo.cul


```lua
boolean
```


---

# vim.wo.culopt


```lua
string
```


---

# vim.wo.cursorbind


```lua
boolean
```


---

# vim.wo.cursorcolumn


```lua
boolean
```


---

# vim.wo.cursorline


```lua
boolean
```


---

# vim.wo.cursorlineopt


```lua
string
```


---

# vim.wo.diff


```lua
boolean
```


---

# vim.wo.eiw


```lua
string
```


---

# vim.wo.eventignorewin


```lua
string
```


---

# vim.wo.fcs


```lua
string
```


---

# vim.wo.fdc


```lua
'0'|'1'|'2'|'3'|'4'...(+15)
```


---

# vim.wo.fde


```lua
string
```


---

# vim.wo.fdi


```lua
string
```


---

# vim.wo.fdl


```lua
integer
```


---

# vim.wo.fdm


```lua
'diff'|'expr'|'indent'|'manual'|'marker'...(+1)
```


---

# vim.wo.fdn


```lua
integer
```


---

# vim.wo.fdt


```lua
string
```


---

# vim.wo.fen


```lua
boolean
```


---

# vim.wo.fillchars


```lua
string
```


---

# vim.wo.fml


```lua
integer
```


---

# vim.wo.fmr


```lua
string
```


---

# vim.wo.foldcolumn


```lua
'0'|'1'|'2'|'3'|'4'...(+15)
```


---

# vim.wo.foldenable


```lua
boolean
```


---

# vim.wo.foldexpr


```lua
string
```


---

# vim.wo.foldignore


```lua
string
```


---

# vim.wo.foldlevel


```lua
integer
```


---

# vim.wo.foldmarker


```lua
string
```


---

# vim.wo.foldmethod


```lua
'diff'|'expr'|'indent'|'manual'|'marker'...(+1)
```


---

# vim.wo.foldminlines


```lua
integer
```


---

# vim.wo.foldnestmax


```lua
integer
```


---

# vim.wo.foldtext


```lua
string
```


---

# vim.wo.lbr


```lua
boolean
```


---

# vim.wo.lcs


```lua
string
```


---

# vim.wo.linebreak


```lua
boolean
```


---

# vim.wo.list


```lua
boolean
```


---

# vim.wo.listchars


```lua
string
```


---

# vim.wo.nu


```lua
boolean
```


---

# vim.wo.number


```lua
boolean
```


---

# vim.wo.numberwidth


```lua
integer
```


---

# vim.wo.nuw


```lua
integer
```


---

# vim.wo.previewwindow


```lua
boolean
```


---

# vim.wo.pvw


```lua
boolean
```


---

# vim.wo.relativenumber


```lua
boolean
```


---

# vim.wo.rightleft


```lua
boolean
```


---

# vim.wo.rightleftcmd


```lua
string
```


---

# vim.wo.rl


```lua
boolean
```


---

# vim.wo.rlc


```lua
string
```


---

# vim.wo.rnu


```lua
boolean
```


---

# vim.wo.sbr


```lua
string
```


---

# vim.wo.scb


```lua
boolean
```


---

# vim.wo.scl


```lua
'auto'|'auto:1'|'auto:2'|'auto:3'|'auto:4'...(+17)
```


---

# vim.wo.scr


```lua
integer
```


---

# vim.wo.scroll


```lua
integer
```


---

# vim.wo.scrollbind


```lua
boolean
```


---

# vim.wo.scrolloff


```lua
integer
```


---

# vim.wo.showbreak


```lua
string
```


---

# vim.wo.sidescrolloff


```lua
integer
```


---

# vim.wo.signcolumn


```lua
'auto'|'auto:1'|'auto:2'|'auto:3'|'auto:4'...(+17)
```


---

# vim.wo.siso


```lua
integer
```


---

# vim.wo.smoothscroll


```lua
boolean
```


---

# vim.wo.sms


```lua
boolean
```


---

# vim.wo.so


```lua
integer
```


---

# vim.wo.spell


```lua
boolean
```


---

# vim.wo.statuscolumn


```lua
string
```


---

# vim.wo.statusline


```lua
string
```


---

# vim.wo.stc


```lua
string
```


---

# vim.wo.stl


```lua
string
```


---

# vim.wo.ve


```lua
string
```


---

# vim.wo.virtualedit


```lua
string
```


---

# vim.wo.wbr


```lua
string
```


---

# vim.wo.wfb


```lua
boolean
```


---

# vim.wo.wfh


```lua
boolean
```


---

# vim.wo.wfw


```lua
boolean
```


---

# vim.wo.winbar


```lua
string
```


---

# vim.wo.winbl


```lua
integer
```


---

# vim.wo.winblend


```lua
integer
```


---

# vim.wo.winfixbuf


```lua
boolean
```


---

# vim.wo.winfixheight


```lua
boolean
```


---

# vim.wo.winfixwidth


```lua
boolean
```


---

# vim.wo.winhighlight


```lua
string
```


---

# vim.wo.winhl


```lua
string
```


---

# vim.wo.wrap


```lua
boolean
```