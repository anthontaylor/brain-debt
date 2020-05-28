CREATE TABLE topic (
       user_guid UUID,
       topics MAP<UUID,TEXT>,
       PRIMARY KEY(user_guid)
);
