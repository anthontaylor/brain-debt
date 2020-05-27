CREATE KEYSPACE IF NOT EXISTS brain_debt
       WITH REPLICATION = {
       'class' : 'SimpleStrategy',
       'replication_factor' : 1
};
