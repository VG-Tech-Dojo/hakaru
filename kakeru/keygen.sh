#!/bin/sh
touch $HOME/.ssh/known_hosts
ssh-keygen -R $1
ssh-keyscan -H $1 >> $HOME/.ssh/known_hosts
