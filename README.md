Go Points
=========

Webservice that collects points from users and exposes friends points in a feed.

DISCLAIMER: there are no security considerations built into this example app.

[![Build Status](https://travis-ci.org/martinvium/gopoints.png)](https://travis-ci.org/martinvium/gopoints)

Installation
------------

    CREATE KEYSPACE gopoints WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

    CREATE TABLE playthroughs (
      user_id text,
      timestamp timeuuid,
      points int,
      PRIMARY KEY (user_id, timestamp)
    );

    CREATE TABLE users (
      user_id text,
      maxpoints int,
      PRIMARY KEY (user_id)
    );

    CREATE TABLE friends (
      user_id text,
      friend_id text,
      maxpoints int,
      PRIMARY KEY (user_id, friend_id)
    );

Dummy Data
----------

    INSERT INTO users (user_id, maxpoints) VALUES ('player_a', 0);
    INSERT INTO users (user_id, maxpoints) VALUES ('player_b', 0);
    INSERT INTO users (user_id, maxpoints) VALUES ('player_c', 0);

    INSERT INTO friends (user_id, friend_id, maxpoints) VALUES ('player_a', 'player_b', 0);
    INSERT INTO friends (user_id, friend_id, maxpoints) VALUES ('player_a', 'player_c', 0);
    INSERT INTO friends (user_id, friend_id, maxpoints) VALUES ('player_b', 'player_a', 0);
    INSERT INTO friends (user_id, friend_id, maxpoints) VALUES ('player_b', 'player_c', 0);
    INSERT INTO friends (user_id, friend_id, maxpoints) VALUES ('player_c', 'player_b', 0);

API
---

Register a playthrough

    POST /playthroughs
    Content-Type: application/json

    {
      "user_id": "abcd",
      "points": 1234
    }

Get friends list