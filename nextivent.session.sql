SELECT "evento"."evento_id","evento"."organizador_id","evento"."categoria_id","evento"."titulo","evento"."descripcion","evento"."lugar","evento"."evento_estado","evento"."cant_me_gusta","evento"."cant_no_interesa","evento"."cant_vendido_total","evento"."imagen_descripcion","evento"."imagen_portada","evento"."video_presentacion","evento"."imagen_escenario","evento"."total_recaudado","evento"."estado","evento"."usuario_creacion","evento"."fecha_creacion","evento"."usuario_modificacion","evento"."fecha_modificacion", f.fecha_evento FROM "evento" 
JOIN evento_fecha ef ON ef.evento_id = evento.evento_id 
JOIN fecha f ON f.fecha_id = ef.fecha_id 
WHERE f.fecha_evento >= CURRENT_DATE 
AND evento.evento_estado = 1 
AND evento.estado = 1 
AND f.fecha_evento = (
    SELECT MIN(f2.fecha_evento)
    FROM fecha f2
    JOIN evento_fecha ef2 ON ef2.fecha_id = f2.fecha_id
    WHERE ef2.evento_id = evento.evento_id
)
ORDER BY ((evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) DESC;



SELECT evento_id, fecha.fecha_evento FROM evento_fecha JOIN fecha ON fecha.fecha_id = evento_fecha_id;


SELECT f.fecha_evento as fecha_, CURRENT_DATE as hoy, GREATEST(1, (f.fecha_evento::date - CURRENT_DATE)) as dif, ((evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) AS score FROM fecha f
JOIN evento_fecha ef on ef.fecha_id = f.fecha_id
JOIN evento ON evento.evento_id = ef.evento_id 
WHERE ef.evento_id = 9
ORDER BY score DESC;

SELECT "evento"."evento_id","evento"."organizador_id","evento"."categoria_id","evento"."titulo","evento"."descripcion","evento"."lugar","evento"."evento_estado","evento"."cant_me_gusta","evento"."cant_no_interesa","evento"."cant_vendido_total","evento"."imagen_descripcion","evento"."imagen_portada","evento"."video_presentacion","evento"."imagen_escenario","evento"."total_recaudado","evento"."estado","evento"."usuario_creacion","evento"."fecha_creacion","evento"."usuario_modificacion","evento"."fecha_modificacion" FROM "evento" JOIN evento_fecha ef ON ef.evento_id = evento.evento_id JOIN fecha f ON f.fecha_id = ef.fecha_id WHERE f.fecha_evento >= CURRENT_DATE AND evento.evento_estado = 1 AND evento.estado = 1 AND f.fecha_evento = ( SELECT MIN(f2.fecha_evento)
                FROM fecha f2
                JOIN evento_fecha ef2 ON ef2.fecha_id = f2.fecha_id
                WHERE ef2.evento_id = evento.evento_id) ORDER BY ((2*evento.cant_me_gusta - evento.cant_no_interesa) / GREATEST(1, (f.fecha_evento::date - CURRENT_DATE))) DESC ;