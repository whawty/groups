# whawty.groups storage schema

The whawty.groups store consists of a directory structure which stores
group membership and user meta data in files and symlinks.
The directory structure looks like this (path names are relative to the
base directory:

    users/
      username        ; yaml file which contains user meta data:
      nicoo           ;   first name, last name, mail, ....
      equinox         ;
      fredl           ;
    groups/
      groupa/
        _meta.yaml    ; yaml file which contains the groups meta data
        equinox       ; symlink to user file in users directory
        nicoo         ; symlink to user file in users directory
      groupb/
        _meta.yaml    ;
        groupa        ; symlink to group directory
        fredl         ; symlink to user file in users directory

A whawty.groups agent must use the following regular expressing to match for
valid user and group names:

     [A-Za-z0-9][-_.@A-Za-z0-9]*

Also user and group names share a name space. This means if there is a user
called `foo` there *must not* be a group called `foo` and vice-versa.

A whawty.groups agent must be configurable to add an implicit group for
any user name. Only the user which has the same name as this group should
be a member of that group.
This membership information must not be stored in the directory but be
generated for group membership queries.

Obviously the storage schema allows to store membership loops. An compliant
agent must be resilient against these situations and in case it detects a
loop print a warning to it's log.
