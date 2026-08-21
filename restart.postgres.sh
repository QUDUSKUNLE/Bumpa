#!/bin/bash

rm -f /usr/local/var/postgresql@14/postmaster.pid

brew services start postgresql@14
