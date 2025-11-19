CREATE USER healthcheck WITH ENCRYPTED PASSWORD 'healthcheck_pass';

GRANT CONNECT ON DATABASE setetes TO healthcheck;
GRANT USAGE ON SCHEMA public TO healthcheck;
GRANT SELECT ON pg_stat_database TO healthcheck;

CREATE OR REPLACE FUNCTION public.health_check()
    RETURNS TEXT AS
$$
BEGIN
    RETURN 'OK';
END;
$$ LANGUAGE plpgsql;

GRANT EXECUTE ON FUNCTION public.health_check() TO healthcheck;

SELECT 'Health check user created successfully' AS status;