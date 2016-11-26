# whawty.groups storage schema

The whawty.groups store consists of a directory structure which stores
group membership and user meta data in files and symlinks.
The directory structure looks like this (path names are relative to the
base directory:

    users/
      username        ; yaml file which contains user meta data:
      nicoo           ;   firstname, lastname, mail, ....
      equinox         ;
      fredl           ;
    groups/
      groupa/
        equinox       ; symlink to user file in users directory
        nicoo         ; symlink to user file in users directory
      groupb/
        groupa        ; symlink to group directory
        fredl         ; symlink to user file in users directory

User and Group names must only contain the following characters:

     [-_.@A-Za-z0-9]

A whawty.groups agent must be configurable to add an implicit group for
any username. Only the user which has the same name as this group should
be a member of that group.
This membership information must not be stored in the directory but be
generated for group membership queries.

Obviously the storage schema allows to store membership loops. An compliant
agent must be resilient against these situations and in case it detects a
loop print a warning to it's log.
