Aria2go
=======
Go bindings for libaria2.

Usage
===
* If you are using `dep`, add the flowing:
```toml
required = [
	anbillon.com/aria2go
]
```
or you can just use：
```shell
$ go get anbillon.com/aria2go
```
* Prepare libaria2, be sure you have installed `gcc`:
```shell
$ cd path/to/aria2go
$ ./prepare.sh
```
or if you want to compile to arm, be sure you have installed `gnueabihf` cross compiler::
```shell
$ ./prepare-arm.sh
```
* Now you can check the example to run.


License
=======

    Copyright (C) 2018 Anbillon Team
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
         http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

