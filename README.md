# Daykon
## What is Daykon?
Daykon is implementation in-memory Redis cache app

TCP client and server 

Operators: GET, SET, DEL, KEYS, SAVE, BACKUP

## Getting started
Open server.exe and client.exe and go on with operators

Or use terminal with
```
src/server> go run server.go
```
and
```
src/client> go run client.go
```
Also you can run server_test.go:
```
src/server> go test
```
___


__GET "keyName"__
```
daykon> SET mykey "Hello, World"
OK
daykon> GET mykey
"Hello, World"
daykon> SET mykey2 "Hi" 100s
OK
daykon> GET mykey2
"Hi" 1m39.771775231s
```
___

__SET "keyName" "value" "timeDuration"__

Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
Such as "300ms", "1.5h" or "2h45m"

Space is unvalid symbol in SET

```
daykon> SET mykey "Hello"
OK
daykon> GET mykey
"Hello"
daykon> SET anotherkey "some_string_in_20_seconds_winthout_spaces" 20s
OK
daykon> SET anotherkey2 {key1:"1",key2:"2"}

```
You can save lists also, because Daykon keep your data in bytes
___
__DEL "keyName1" "keyName2" ...__

Return the number of keys that were removed
```
daykon> SET key1 "1"
OK
daykon> SET key2 "2"
OK
daykon> DEL key1 key2
2
```
___
__KEYS "pattern"__

```
daykon> KEYS .*name.*
1) name
2) firstname
daykon> KEYS .
1) name
2) firstname
3) age
4) ag6
daykon> KEYS ag\d
1) ag6
```
|Symbol|Point|
|:----:|:---|
|h`.`llo|Matches hllo and heeeello|
|h`[ab]`llo|Matches hello and hallo, but not hillo|
|h`[^e]`llo|Matches hallo, hbllo, ... but not hello|
|h`[a-b]`llo|Matches hallo and hbllo|
|h`\d`llo|Matches h1llo and h2llo, but not hello|
|h`\D`llo|Matches hello and hallo, but not h2llo|
___
Number of symbol X
|Symbol|Point|
|:----:|:---|
|x*|Zero or more x, prefers more
|x*?|Zero or more x, prefers less (not greedy)
|x+|One or more x, prefers more
|x+?|One or more x, prefers less (not greedy)
|x?|Zero or one x, prefers one
|x??|Zero or one x, prefers zero
|x{n}| x n times
___

__SAVE "fileName"__

You can save all your data in "fileName".bin file:
```
daykon> SAVE main
OK
```
 __BACKUP "fileName"__

 And read .bin file:
 ```
daykon> BACKUP main
OK
```
___
