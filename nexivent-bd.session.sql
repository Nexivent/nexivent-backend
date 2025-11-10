-- Genera lista de tablas (en psql puedes copiar el resultado)
SELECT format('%I.%I', schemaname, tablename)
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY 1;
-- Luego:
TRUNCATE TABLE public.categoria,
public.evento,
public.cupon,
public.orden_de_compra,
public.metodo_de_pago,
public.ticket,
public.sector,
public.tipo_de_ticket,
public.perfil_de_persona,
public.usuario RESTART IDENTITY CASCADE;
DO $$
DECLARE r RECORD;
v_max bigint;
BEGIN FOR r IN
SELECT n.nspname AS schemaname,
    s.relname AS sequencename,
    t.relname AS tablename,
    a.attname AS columnname
FROM pg_class s
    JOIN pg_namespace n ON n.oid = s.relnamespace
    JOIN pg_depend d ON d.objid = s.oid
    AND d.deptype = 'a' -- owned-by
    JOIN pg_class t ON t.oid = d.refobjid
    JOIN pg_attribute a ON a.attrelid = t.oid
    AND a.attnum = d.refobjsubid
WHERE s.relkind = 'S' -- sequences
    AND n.nspname = 'public' -- <-- cambia si usas otro schema
    LOOP EXECUTE format(
        'SELECT COALESCE(MAX(%I), 0) FROM %I.%I',
        r.columnname,
        r.schemaname,
        r.tablename
    ) INTO v_max;
-- Deja la secuencia en v_max (para que el próximo nextval sea v_max+1)
EXECUTE format(
    'SELECT setval(%L, %s, true)',
    r.schemaname || '.' || r.sequencename,
    v_max
);
RAISE NOTICE 'Sincronizada: %.% (tabla %.%, columna %) => max=%',
r.schemaname,
r.sequencename,
r.schemaname,
r.tablename,
r.columnname,
v_max;
END LOOP;
END $$;