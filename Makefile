# Copyright 2015 Measures authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# lists all available targets
list:
	@sh -c "$(MAKE) -p __no_targets__ | \
		awk -F':' '/^[a-zA-Z0-9][^\$$#\/\\t=]*:([^=]|$$)/ { \
			split(\$$1,A,/ /);for(i in A)print A[i] \
		}' | \grep -v '__\$$' | \grep -v 'make\[' | \grep -v 'Makefile' | sort"

# required for list
__no_targets__:

__test:
	@go test ./... $(verbosity) $(cover)

__count:
	@go test -test.v -covermode=count -coverprofile=cover.out

__func:
	@go tool cover -func=cover.out

__html:
	@go tool cover -html=cover.out -o cover.html
	@read -p "Open report? [y/N] " -n 1; \
	if [[ $$REPLY = y ]]; then \
		open cover.html; \
	fi

test:
	@$(MAKE) __test verbosity="-test.v" cover=""

cover:
	@$(MAKE) __test verbosity="" cover="-cover -coverprofile=cover.out"

func-cover: cover __func

html-cover: cover __html

func-count: __count __func

html-count: __count __html
