# OPW

This tool was born from a hatred of having to fumble for 1password on applications
that don't support it. (Namely, in my specific case, OpenVPN).

It is a simple wrapper around the 1password CLI tool `op`, it didn't look like the
1password go libraries were open source, so I had to wrap the CLI.

## Overview

The core functionality is the ability to assign an alias to a set of particular items
(logins) in any of your vaults, and easily capture the password from them.

> opw get -p vpn

This command will use the alias `vpn` (set in the configuration file `~/.opw.yaml`)
to grab the JSON blob from the `op` tool (calling `op get item $uuid`) and output
the password, which could easily be piped to `pbcopy` or used in an `alfred` script.

I've created an alfred workflow and made it available in the `alfred` directory.

## Setup

Install the application with `go` (tested on 1.13):

> go get -u github.com/shawncatz/opw

Run the following command to generate your configuration:

> opw init

Edit the file `~/.opw.yaml` and set the required values.

```
subdomain: name # subdomain of your 1password account
email: me@example.com # your email address
cache: /your/home/.opw.cache # where you want to store your cached session
aliases: # set of aliases (key: value)
  nickname: uuid # find uuid by running 'opw list'

# For both passphrase and secret, you specify how to obtain the secret
# rather than specifying it directly.
# if the string starts with 'file:' then it will load from the file
# if the string starts with 'keychain:' then it will load from the keychain
# the keyring library supports MacOS, Linux (d-bus), and Windows
# for keychain, the value is specified as 'keychain:service:account'
# using the keychain is more secure

# The secret currently isn't used
#secret: file:/path/to/file/containing/secret
secret: keychain:opw-secret:subdomain

#passphrase: file:/path/to/file/containing/passphrase
passphrase: keychain:opw-passphrase:subdomain # keychain entry
```

## Usage

To use the aliases, you must first look up the `uuid`, you can use the `list` command
to view a list of all items in your vaults.

> opw list

After finding the `uuid` of the item for which you wish to create an alias, add an
entry in the aliases configuration, like below:

```
aliases:
  vpn: uuid
``` 

You should now be able to acquire the password value of that item by using the alias
name. The output of the following commands should be the same.

> opw get vpn
> opw get uuid 
